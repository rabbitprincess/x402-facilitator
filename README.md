# x402-facilitator

**x402-facilitator** is a Go-based middleware that settles on-chain payments authorized via the [x402 protocol](https://x402.dev).

## Prerequisites
- Docker
- Docker Compose

## Support Schemes
| Scheme     | Status           | Description                   |
|------------|------------------|-------------------------------|
| EVM       | ✅ Supported      | Ethereum and EVM chains       |
| Solana    | 🚧 Planned        |                               |
| Sui       | 🚧 Planned        |                               |
| Tron      | 🚧 Planned        |                               |

## How to run

### Build binary
```bash
make build
```

### Run x402-facilitator using docker
```bash
docker compose up
```

### Run x402-client
```
Usage:
  client [flags]

Flags:
  -A, --amount string    Amount to send
  -F, --from string      Sender address
  -h, --help             help for x402-client
  -n, --network string   Blockchain network to use (default "base-sepolia")
  -P, --privkey string   Sender private key
  -s, --scheme string    Scheme to use (default "evm")
  -T, --to string        Recipient address
  -t, --token string     token contract for sending (default "USDC")
  -u, --url string       Base URL of the facilitator server (default "http://localhost:9090")
```

## Api Specification
After starting the service, open your browser to:
```
/swagger/index.html
```

## Contributing
We welcome any contributions! Feel free to open issues or submit pull requests at any time.
