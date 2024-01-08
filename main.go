package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	_ "net/http/pprof"

	"golang.org/x/sync/errgroup"

	"github.com/jcleira/coding-challenge/internal/domain/services"
	"github.com/jcleira/coding-challenge/internal/infra/handlers"
	"github.com/jcleira/coding-challenge/internal/infra/repositories"
)

// I will be using constants for the configuration, but it could be easily
// extended to use a configuration file, or environment variables.
const (
	// solanaRPCURL is the URL of the Solana RPC endpoint
	solanaRPCURL = "https://api.devnet.solana.com"

	// valutPath is the path where the wallets will be stored
	vaultPath = "./tmp/wallets"

	// exchangeURL is the URL of the exchange API
	exchangeURL = "https://api.kraken.com/0/public/Ticker"
)

func main() {
	vault, err := repositories.NewVault(vaultPath)
	if err != nil {
		slog.Error("error initializing vault: ", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	exchange, err := repositories.NewExchange(ctx, exchangeURL)
	if err != nil {
		slog.Error("error initializing exchange: ", err)
		os.Exit(1)
	}

	solana := repositories.NewSolana(solanaRPCURL)

	transactionsGetterHandler := handlers.NewTransactionsGetterHandler(
		services.NewTransactionsGetter(solana, exchange),
	)

	transactionsSenderHandler := handlers.NewTransactionsSenderHandler(
		services.NewTransactionsSender(vault, solana, exchange),
	)

	walletInitializerHandler := handlers.NewWalletInitializerHandler(
		services.NewWalletInitializer(vault),
	)

	walletBalanceGetterHandler := handlers.NewWalletBalanceGetterHandler(
		services.NewWalletBalanceGetter(solana, exchange),
	)

	exchangeRateGetterHandler := handlers.NewExchangeRateGetterHandler(
		services.NewExchangeRateGetter(exchange),
	)

	http.HandleFunc("/init", walletInitializerHandler.Handler())
	http.HandleFunc("/balance", walletBalanceGetterHandler.Handler())
	http.HandleFunc("/exchange_rate", exchangeRateGetterHandler.Handler())
	http.HandleFunc("/send", transactionsSenderHandler.Handler())
	http.HandleFunc("/transactions", transactionsGetterHandler.Handler())

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		if err := http.ListenAndServe(":8888", nil); err != nil {
			slog.Error("error starting server: ", err)
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		slog.Error("error running server: ", err)
		cancel()
		os.Exit(1)
	}

	slog.Info("Server started")
}
