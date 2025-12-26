package main

import "time"

type Order struct {
	OrderID            string `json:"order_id"`
	ProductID          string `json:"product_id"`
	Side               string `json:"side"`
	FilledSize         string `json:"filled_size"`
	AverageFilledPrice string `json:"average_filled_price"`
	TotalFees          string `json:"total_fees"`
	LastFillTime       string `json:"last_fill_time"`
	Status             string `json:"status"`
}

type OrdersResponse struct {
	Orders []Order `json:"orders"`
}

type Fill struct {
	EntryID   string `json:"entry_id"`
	TradeID   string `json:"trade_id"`
	OrderID   string `json:"order_id"`
	ProductID string `json:"product_id"`
	Price     string `json:"price"`
	Size      string `json:"size"`
	Commission string `json:"commission"`
	Side      string `json:"side"`
	TradeTime string `json:"trade_time"`
}

type FillsResponse struct {
	Fills []Fill `json:"fills"`
}

type Position struct {
	Size       float64
	EntryPrice float64
	OrderID    string
}

type PositionQueue struct {
	ProductID string
	Queue     []Position
}

type Trade struct {
	Date        time.Time
	ProductID   string
	Side        string
	Size        float64
	Price       float64
	Fees        float64
	RealizedPnL float64
	OrderID     string
}

type DayPnL struct {
	Date       string  `json:"date"`
	DayOfMonth int     `json:"dayOfMonth"`
	DayOfWeek  int     `json:"dayOfWeek"`
	PnL        float64 `json:"pnl"`
	TradeCount int     `json:"tradeCount"`
	HasNotes   bool    `json:"hasNotes"`
}

type WeekPnL struct {
	WeekNumber     int     `json:"weekNumber"`
	WeekPnL        float64 `json:"weekPnL"`
	WeekTradeCount int     `json:"weekTradeCount"`
}

type MonthData struct {
	Month      string    `json:"month"`
	Year       int       `json:"year"`
	MonthlyPnL float64   `json:"monthlyPnL"`
	Days       []DayPnL  `json:"days"`
	Weeks      []WeekPnL `json:"weeks"`
}

type Report map[string]MonthData

type PnLCalculator interface {
	ProcessOrders(orders []Order) ([]Trade, error)
	ProcessFills(fills []Fill) ([]Trade, error)
}
