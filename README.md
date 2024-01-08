# Coding challenge

Implement a Solana payment system that manage wallets, send and retrieve transactions.

## Requirements

This server should implement the following endpoints:

```
// This endpoint should write the private key to disk, so it can be used to
// sign transactions later.
POST /init {} {
  "public_key": "t5LVRrjhJU3g9FegCaL5XzXXcQdX7ZagEFrJ9NnjFeH"
}

POST /exchange_rate {} {
  "sol_eur": 85.87262
}

POST /balance{
  "public_key": "t5LVRrjhJU3g9FegCaL5XzXXcQdX7ZagEFrJ9NnjFeH"
} {
  "balance": "EUR 18.11"
}

POST /send {
  "public_key": "t5LVRrjhJU3g9FegCaL5XzXXcQdX7ZagEFrJ9NnjFeH",
  "to": "4zdNGgAtFsW1cQgHqkiWyRsxaAgxrSRRynnuunxzjxu1",
  "amount": "EUR 5.05"
} {
  "signature": "61gr5CoTKS54CnP11DL9kdrLkvwxhiFBYboFnk5f8FPVsVtsquXmBEZH3ZiSq1Y2M3C7a1o7QCm5yJMJTr3Hvt5h"
}

POST /transasctions {
  "public_key": "t5LVRrjhJU3g9FegCaL5XzXXcQdX7ZagEFrJ9NnjFeH"
} {
  // Sorted newest first, oldest last
  "transactions": [
    {
      // Timestamp should be in RFC3339 format
      "created": "2023-09-01T16:01:05Z",
      // Positive amount indicates a credit - i.e. someone has sent us money
      "amount": "EUR 3",
      "counterparty": "4zdNGgAtFsW1cQgHqkiWyRsxaAgxrSRRynnuunxzjxu1",
      "signature": "MJTr3Hvt5h..."
    },
    {
      "created": "2023-09-01T16:00:05Z",
      // Negative amount indicates a debit - i.e. we've sent someone money
      "amount": "EUR -5.05",
      "counterparty": "4zdNGgAtFsW1cQgHqkiWyRsxaAgxrSRRynnuunxzjxu1",
      "signature": "61gr5CoTKS..."
    }
  ]
}
```

- The program should be written in Go, and we've provided some sample code to get you started. The sample code is just a guide - feel free to change the structure if you like. Please write the code as if it were going to be used in production.
- The server should listen on port 8888.
- The `send` endpoint should block until the transaction gets to the `confirmed` commitment level.
- All Solana operations that take a commitment level should use `confirmed`. I should be able to `send` a transaction and when that call returns, read it using the `transactions` endpoint.
- The `transactions` endpoint should print all transactions on the account, converting the amount in SOL to EUR at the current exchange rate.
- The `transactions` endpoint should be able to list >1000 transactions efficiently and quickly.
- The wallet should use the Solana Devnet (Solana's staging environment).
- We run automated tests against the server, so it's important that the endpoint names, input and output objects match the specification above.
- We'll run your code in a Docker container for security. It would be great if you could check your code still runs with the Dockerfile provided by running:

### Questions

If you've got any questions just email us - we're happy to help.

### Time

We don't want the task to take longer than four hours to complete. If you're running over this time, please stop and we'll be happy to review what you've written.

### The interview

In the interview, we'll talk through the code you've written. We'll be assessing your understanding of the code and how it runs. With that in mind:

- You can use external libraries, but please understand what they're doing
- You can use AI tools like Copilot or ChatGPT, but please make sure you understand the code they've generated

### Exchange rate

Your code will need to convert SOL to EUR and back. There's a free API for getting the exchange rate at `https://api.kraken.com/0/public/Ticker?pair=SOLEUR`. The response is [documented here](https://docs.kraken.com/rest/#tag/Market-Data/operation/getTickerInformation), but we recommend using the field `p[1]`. For example:

```sh
curl "https://api.kraken.com/0/public/Ticker?pair=SOLEUR" | jq '.result.SOLEUR.p[1]'
"84.28463"
```

--- 

# Candidate notes

### 1. Summary
This coding test has been a fun and insightful journey, significantly as it deepened my understanding of Solana's intricacies, a domain I've only recently begun to explore.

Time Allocation:
- Coding: Approximately 4 hours 45 minutes.
- Load testing: 30 minutes.
- Final review and documentation: 45 minutes.

### 2. Design Choices
#### 2.1 Domain Driven Design Architecture
I commenced with a simplistic approach but soon evolved to include a three-layer architecture:
- **HTTP Layer**: Interface to external HTTP requests.
- **Domain Layer**: Encompasses Services and Aggregates, which is pivotal for managing conversions and business logic.
- **Repository Layer**: Includes interfaces with Solana, Kraken Oracle, and a Vault for wallet management.

#### 2.2 Kraken Rate Retrieval
I established a dedicated repository for Kraken, featuring an engine to update currency rates frequently. This subsystem was designed to avoid additional third-party HTTP calls on user requests. Key features include:
- Mandatory initial rate fetch for service initialization, ensuring operations commence with a valid exchange rate.
- Utilization of mutex for state management over channel communication, chosen for its simplicity and effectiveness in this context.

#### 2.3 Solana RPC Integration
A CQRS (Command Query Responsibility Segregation) approach is what I would recommend for production environments, especially for asynchronous handling of RPC calls. This would entail:
- **POST requests**: Implementing commands rather than direct RPC calls, with updates delivered via WebSockets.
- **GET requests**: Leveraging Solana Geyser Postgres plugin for efficient transaction tracking and Solana Geyser Kafka for real-time updates.

Challenges:
- The `GetSignatures` method (and the `GetTransaction` calls) were the primary tool for transaction retrieval but proved inefficient for larger transaction volumes.

### 3. Identified Challenges that I didn't finish
#### 3.1 Incomplete Transaction Amount Retrieval
Due to time constraints, accurately decoding transaction amounts using `solana-go` remained unresolved.

#### 3.2 Suboptimal Transactions Endpoint Performance
Load testing indicated a 99th percentile response time of 829.418ms. While the performance was hampered by RPC throttling on Devnet, this area needs further optimization.

- I was spending up to 45 min trying to get the amount right while fetching ALL the transactions.
- I could not get the proper decoding for the Instruction right with solana-go.
- It may be straightforward, but I didn't find the proper docs or a meaningful resource in the solana-go repo.

Current performance profile after load testing with [vegeta](https://github.com/tsenart/vegeta).
```bash
Requests      [total, rate, throughput]         150, 5.03, 2.72
Duration      [total, attack, wait]             30.54s, 29.8s, 739.843ms
Latencies     [min, mean, 50, 90, 95, 99, max]  43.175ms, 660.859ms, 736.335ms, 739.732ms, 743.434ms, 829.418ms, 879.144ms
Bytes In      [total, mean]                     77754, 518.36
Bytes Out     [total, mean]                     9150, 61.00
Success       [ratio]                           55.33%
Status Codes  [code:count]                      200:83  500:67
```

I also did some pprof, but again without meaningful results due to the RPC throttling:

![pprof](https://i.imgur.com/PYGpi97.jpeg)

### 4. Unimplemented Features (Due to Time Constraints)
#### 4.1 Tests for Solana Repository
Future improvements should include a dedicated interface for Solana methods, allowing for comprehensive testing.

#### 4.2 Integration Testing
Although unit tests were conducted for domain services, comprehensive integration tests incorporating mocked HTTP calls were not completed.

#### 4.3 Enhanced Logging with `slog`
While `slog` was utilized, it lacked proper initialization and dependency management, which could be improved for better logging and error tracking.

