# Crypto Match Core

High-performance in-memory matching engine for cryptocurrency exchange, written in Go.

## Features

- **In-Memory Matching**: O(log n) order book operations using red-black tree
- **Multi-Market Support**: Supports multiple trading pairs (BTC/USDT, ETH/BTC, etc.)
- **Price-Time Priority**: Fair order matching algorithm
- **PostgreSQL Persistence**: Slice-based snapshot mechanism with operlog replay
- **Kafka Integration**: Real-time order/deal/balance event publishing
- **gRPC API**: High-performance RPC interface

## Architecture

```
┌──────────────────────────────────────────────────────────┐
│                    API Server                           │
│                  (gRPC + protobuf)                       │
└─────────────────────────┬────────────────────────────────┘
                          │
┌─────────────────────────▼────────────────────────────────┐
│                 Matching Engine                          │
│  ┌──────────────────────────────────────────────────┐    │
│  │              OrderBook (Red-Black Tree)          │    │
│  │    Asks (ascending)  │  Bids (descending)        │    │
│  └──────────────────────────────────────────────────┘    │
│  ┌──────────────────────────────────────────────────┐    │
│  │              Balance Manager (In-Memory)          │    │
│  └──────────────────────────────────────────────────┘    │
└─────────────────────────┬────────────────────────────────┘
                          │
          ┌───────────────┼───────────────┐
          ▼               ▼               ▼
    ┌──────────┐   ┌──────────┐   ┌──────────┐
    │ Storage  │   │ OperLog  │   │ Message  │
    │ (PG)     │   │ (PG)     │   │ (Kafka)  │
    └──────────┘   └──────────┘   └──────────┘
```

## Quick Start

### Prerequisites

- Go 1.18+
- PostgreSQL 13+
- Kafka 2.8+

### Build

```bash
go build ./...
```

### Run

```bash
go run ./cmd/server
```

## Matching Algorithm

The engine uses **price-time priority**:

1. **Bid orders** (buy): Match against asks from lowest to highest price
2. **Ask orders** (sell): Match against bids from highest to lowest price
3. At the same price level, earlier orders have priority

### Fee Structure

- **Taker Fee**: 0.1% (paid by the order that crosses the spread)
- **Maker Rebate**: 0.05% (earned by the order that provides liquidity)

## Persistence

### Slice Mechanism

Every `slice_interval` seconds:
1. Snapshot current order book and balances to PostgreSQL
2. Clean up old slices older than `slice_keeptime`

### Recovery

On startup:
1. Load latest slice from `slice_history`
2. Restore orders and balances
3. Replay operlog entries since last slice

## Testing

```bash
go test ./... -v
```

## License

MIT