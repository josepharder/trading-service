# Coinbase PnL Calculator

A Go-based tool for calculating Profit and Loss (PnL) from Coinbase Advanced Trade futures orders using FIFO (First In First Out) accounting.

## Overview

This tool fetches your Coinbase Advanced Trade futures trading history and generates daily PnL reports using FIFO position tracking.

## Features

- **FIFO Position Tracking**: Matches buy and sell orders using First In First Out methodology
- **Futures Contract Support**: Handles different contract multipliers (e.g., ETP: 0.1x, XPP: 500x)
- **Daily PnL Reports**: Generates JSON reports with daily profit/loss by month
- **API Integration**: Fetches data from Coinbase Advanced Trade API
- **Fee Accounting**: Includes trading fees in PnL calculations

## Prerequisites

- Go 1.24.1 or higher
- Coinbase Advanced Trade account with API credentials
- Active or historical futures trading activity

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/coinbase-pnl-calculator.git
cd coinbase-pnl-calculator
```

2. Install dependencies:
```bash
go mod download
```

3. Set up your environment variables:
```bash
cp .env.example .env
```

4. Edit `.env` and add your Coinbase API credentials:
```
KEY_NAME=your_api_key_name
KEY_SECRET=your_api_private_key
```

### Getting Coinbase API Credentials

1. Log in to [Coinbase Advanced Trade](https://www.coinbase.com/advanced-trade)
2. Navigate to Settings → API
3. Create a new API key with the following permissions:
   - View: Trading history
   - Trade: Not required
4. Copy the Key Name and Private Key to your `.env` file

## Usage

Run the calculator:
```bash
go run .
```

The tool will:
1. Fetch all filled futures orders from Coinbase API
2. Process orders chronologically using FIFO methodology
3. Calculate daily PnL for each trading day
4. Generate `pnl_report.json` with monthly breakdowns

## Output Format

The generated `pnl_report.json` contains monthly data structured as:

```json
{
  "2025-08": {
    "month": "August",
    "year": 2025,
    "days": [
      {
        "date": "2025-08-01",
        "pnl": -308.67,
        "tradeCount": 5,
        "hasNotes": false
      }
    ]
  }
}
```

### Report Fields

- **month**: Month name
- **year**: Year
- **days**: Array of daily PnL entries
  - **date**: ISO 8601 date (YYYY-MM-DD)
  - **pnl**: Realized profit/loss for the day (USD)
  - **tradeCount**: Number of trades executed
  - **hasNotes**: Flag for manual notes (currently always false)

## How FIFO Works

The calculator uses FIFO (First In First Out) position tracking:

1. **Buy Orders**: Create positions in a queue with entry price and size
2. **Sell Orders**: Match against oldest positions first
3. **PnL Calculation**: `(sell_price - entry_price) × size - fees`
4. **Partial Fills**: Handles partial position closures correctly

### Example

```
BUY  1.0 BTC @ $50,000  → Position queue: [1.0 @ $50,000]
BUY  0.5 BTC @ $51,000  → Position queue: [1.0 @ $50,000, 0.5 @ $51,000]
SELL 1.2 BTC @ $52,000  → Matches: 1.0 @ $50k + 0.2 @ $51k
                        → PnL: (52k - 50k) × 1.0 + (52k - 51k) × 0.2 - fees
```

## Project Structure

```
.
├── main.go         # Application entry point and orchestration
├── auth.go         # JWT authentication for Coinbase API
├── fetcher.go      # API client for fetching orders
├── pnl.go          # FIFO PnL calculation engine
├── report.go       # Report generation and formatting
├── types.go        # Data structure definitions
├── go.mod          # Go module dependencies
└── .env            # API credentials (not in git)
```

## Development

### Building

```bash
go build -o coinbase-pnl-calculator .
```

### Running the Binary

```bash
./coinbase-pnl-calculator
```

## License

MIT License - See LICENSE file for details

## Disclaimer

This tool is for informational purposes only. Always consult with a tax professional for official tax reporting. The authors are not responsible for any financial decisions made based on this tool's output.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

For issues or questions, please open an issue on GitHub.
