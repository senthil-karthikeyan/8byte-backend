# 8byte Backend

The backend service for the 8byte portfolio tracker, built with Go and the Gin web framework. It provides APIs to serve portfolio data, fetch real-time stock prices from Yahoo Finance, and retrieve fundamental data (PE, EPS) using Google Sheets as a proxy.

## Features

- **Portfolio Management**: Loads portfolio data from a local `portfolio.csv` file.
- **Real-time Prices**: Fetches live stock prices using the Yahoo Finance API with caching (2-minute TTL).
- **Fundamental Data**: Retrieves PE and EPS ratios by leveraging Google Sheets `GOOGLEFINANCE` functions dynamically.
- **WebSocket Support**: (Experimental) WebSocket endpoint for live updates.

## Prerequisites

- **Go**: Version 1.20 or later.
- **Google Cloud Service Account**: Required for accessing the Google Sheets API.

## Installation

1.  **Clone the repository** (if not already done).

2.  **Navigate to the backend directory**:

    ```bash
    cd backend
    ```

3.  **Install dependencies**:
    ```bash
    go mod download
    ```

## Configuration

Set the following environment variables. You can export them in your shell or use a `.env` file manager if configured (note: the current code reads directly from `os.Getenv`).

| Variable                      | Description                                                |
| :---------------------------- | :--------------------------------------------------------- |
| `SPREADSHEET_ID`              | The ID of the Google Sheet used for fetching fundamentals. |
| `GOOGLE_SERVICE_ACCOUNT_JSON` | The full JSON content of your Google Service Account key.  |

## Usage

Start the server:

```bash
go run main.go
```

The server will start on port `8080`.

## API Endpoints

### `GET /portfolio`

Returns the calculated portfolio data including:

- Purchase details (quantity, price, investment)
- Current Market Price (CMP) and Market Value
- Gain/Loss calculations
- Fundamental data (PE, EPS)

### `GET /ws`

WebSocket endpoint for streaming portfolio updates.
