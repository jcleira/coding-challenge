package repositories

import (
	"context"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"

	"github.com/jcleira/coding-challenge/internal/domain/aggregates"
)

const (
	// confirmationTimeout is the timeout for the transaction confirmation
	//
	// I'm keeping it low for the sake of providing a quick feedback to the user,
	// and I'd be returning an especific error for the frontend to handle. This
	// is a design decision that requires a proper coordination in the frontend,
	// to inform the user that the transaction has been sent and we're waiting.
	confirmationTimeout = 3 * time.Second

	// confirmationInterval is the interval to check for the transaction
	// confirmation.
	confirmationInterval = 500 * time.Millisecond
)

// Solana defines the dependencies for sending transactions to the Solana
// blockchain.
type Solana struct {
	client *rpc.Client
}

// NewSolana creates a new Solana.
func NewSolana(rpcURL string) *Solana {
	return &Solana{
		client: rpc.New(rpcURL),
	}
}

// SendTransaction sends a transaction to the Solana blockchain, returning the
// transaction signature.
//
// I didn't rely on the solana-go client wait for confirmation, because I was
// not aware that it did existed, I could be using it even though we migh want
// to have a custom implementation of the interval and timeout.
func (s *Solana) SendTransaction(ctx context.Context,
	transaction aggregates.Transaction, wallet aggregates.Wallet) (string, error) {
	fromPublicKey, err := solana.PublicKeyFromBase58(wallet.PublicKey)
	if err != nil {
		return "", fmt.Errorf("error converting string to solana.PublicKey: %w", err)
	}

	toPublicKey, err := solana.PublicKeyFromBase58(transaction.CounterParty)
	if err != nil {
		return "", fmt.Errorf("error converting string to solana.PublicKey: %w", err)
	}

	transferInstruction := system.NewTransferInstruction(
		transaction.AmountLAM,
		fromPublicKey,
		toPublicKey,
	).Build()

	recentBlockhash, err := s.GetRecentBlockhash(ctx)
	if err != nil {
		return "", fmt.Errorf("error getting recent blockhash: %w", err)
	}

	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			transferInstruction,
		},
		recentBlockhash,
		solana.TransactionPayer(fromPublicKey),
	)
	if err != nil {
		return "", fmt.Errorf("error creating transaction: %w", err)
	}

	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			accountFrom := solana.PrivateKey(wallet.PrivateKey)
			if accountFrom.PublicKey().Equals(key) {
				return &accountFrom
			}
			return nil
		},
	)
	if err != nil {
		return "", fmt.Errorf("error signing transaction: %w", err)
	}

	signature, err := s.client.SendTransaction(ctx, tx)
	if err != nil {
		return "", fmt.Errorf("error sending transaction: %w", err)
	}

	ticker := time.NewTicker(confirmationInterval)
	defer ticker.Stop()

	for {
		select {
		case <-time.After(confirmationTimeout):
			return "", aggregates.ErrTransactionConfirmationTimeout
		case <-ticker.C:
			status, err := s.client.GetSignatureStatuses(ctx, false, signature)
			if err != nil {
				return "", fmt.Errorf("error getting signature status: %w", err)
			}

			if status != nil && len(status.Value) > 0 && status.Value[0] != nil {
				switch {
				case status.Value[0].Err != nil:
					return "", fmt.Errorf("error confirming transaction: %v", status.Value[0].Err)
				case status.Value[0].ConfirmationStatus == rpc.ConfirmationStatusConfirmed:
					return signature.String(), nil
				default:
					continue
				}
			}
		}
	}
}

// GetRecentBlockhash gets the recent blockhash from the Solana blockchain.
func (s *Solana) GetRecentBlockhash(ctx context.Context) (solana.Hash, error) {
	recentBlockHash, err := s.client.GetRecentBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return solana.Hash{}, fmt.Errorf("error getting recent blockhash: %w", err)
	}

	if recentBlockHash.Value == nil {
		return solana.Hash{}, fmt.Errorf("error getting recent blockhash result: %w", aggregates.ErrNoRecentBlockHashValue)
	}

	return recentBlockHash.Value.Blockhash, nil
}

func (s *Solana) GetBalance(ctx context.Context, publicKey string) (uint64, error) {
	publicKeySol, err := solana.PublicKeyFromBase58(publicKey)
	if err != nil {
		return 0, fmt.Errorf("error decoding public key: %w", err)
	}

	balance, err := s.client.GetBalance(ctx, publicKeySol, rpc.CommitmentFinalized)
	if err != nil {
		return 0, fmt.Errorf("error getting balance: %w", err)
	}

	return balance.Value, nil
}

// GetTransactions gets the transactions for a given public key.
func (s *Solana) GetTransactions(ctx context.Context, publicKey string) ([]aggregates.Transaction, error) {
	publicKeySol, err := solana.PublicKeyFromBase58(publicKey)
	if err != nil {
		return nil, fmt.Errorf("error decoding public key: %w", err)
	}

	signatures, err := s.client.GetSignaturesForAddress(ctx, publicKeySol)
	if err != nil {
		return nil, fmt.Errorf("error getting transactions: %w", err)
	}

	transactions := make([]aggregates.Transaction, len(signatures))

	for i := range signatures {
		tx, err := s.client.GetTransaction(ctx,
			signatures[i].Signature,
			&rpc.GetTransactionOpts{
				Encoding: solana.EncodingBase64,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("error getting transaction: %w", err)
		}

		parsed, err := tx.Transaction.GetTransaction()
		if err != nil {
			return nil, fmt.Errorf("error parsing transaction: %w", err)
		}

		counterParty, err := getCounterParty(publicKey, parsed.Message.AccountKeys)
		if err != nil {
			return nil, fmt.Errorf("error getting counter party: %w", err)
		}

		amount, err := getAmount(parsed.Message)
		if err != nil {
			return nil, fmt.Errorf("error getting amount: %w", err)
		}

		transactions[i] = aggregates.Transaction{
			BlockTime:    convertToTime(tx.BlockTime),
			Signature:    signatures[i].Signature.String(),
			CounterParty: counterParty,
			AmountLAM:    amount,
		}
	}

	return transactions, nil
}

func convertToTime(unixTime *solana.UnixTimeSeconds) time.Time {
	if unixTime == nil {
		return time.Time{}
	}

	return time.Unix(int64(*unixTime), 0)
}

func getCounterParty(pubKey string, accounts []solana.PublicKey) (string, error) {
	for _, account := range accounts {
		accountString := account.String()
		if accountString != pubKey {
			return accountString, nil
		}
	}

	return "", aggregates.ErrNoCounterParty
}

// getAmount makes the assumption that the transaction is a transfer transaction
// and that the first instruction is the transfer instruction. This is not a
// safe assumption...
//
// TODO: This method doesnt return the proper amount... I've spent quite some time
// trying to figure out how to get the proper amount from the instruction data
// but I'm not sure how to do it.
func getAmount(message solana.Message) (uint64, error) {
	if len(message.Instructions) == 0 {
		return 0, aggregates.ErrNoInstructions
	}

	if len(message.Instructions) > 1 {
		return 0, aggregates.ErrMultipleInstructions
	}

	instruction := message.Instructions[0]

	if instruction.Data[0] != 2 {
		return 0, aggregates.ErrInvalidInstruction
	}

	if len(instruction.Data) < 9 {
		return 0, aggregates.ErrInvalidInstruction
	}

	return binary.LittleEndian.Uint64(instruction.Data[1:9]), nil
}
