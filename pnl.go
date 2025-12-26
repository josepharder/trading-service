package main

import (
	"fmt"
	"sort"
	"strconv"
	"time"
)

type FIFOCalculator struct {
	positions map[string][]Position
}

func NewFIFOCalculator() *FIFOCalculator {
	return &FIFOCalculator{
		positions: make(map[string][]Position),
	}
}

func (f *FIFOCalculator) ProcessOrders(orders []Order) ([]Trade, error) {
	sortOrdersByTime(orders)

	var trades []Trade

	for _, order := range orders {
		trade, err := f.processOrder(order)
		if err != nil {
			return nil, fmt.Errorf("processing order %s: %w", order.OrderID, err)
		}
		trades = append(trades, trade)
	}

	return trades, nil
}

func (f *FIFOCalculator) ProcessFills(fills []Fill) ([]Trade, error) {
	return nil, fmt.Errorf("fill-level processing not yet implemented")
}

func (f *FIFOCalculator) processOrder(order Order) (Trade, error) {
	size := parseFloat(order.FilledSize)
	value := parseFloat(order.FilledValue)
	fees := parseFloat(order.TotalFees)

	fillTime, err := time.Parse(time.RFC3339, order.LastFillTime)
	if err != nil {
		return Trade{}, fmt.Errorf("parsing time: %w", err)
	}

	trade := Trade{
		Date:      fillTime,
		ProductID: order.ProductID,
		Side:      order.Side,
		Size:      size,
		Price:     value / size,
		Fees:      fees,
		OrderID:   order.OrderID,
	}

	if order.Side == "BUY" {
		trade.RealizedPnL = -fees
		f.addPosition(order.ProductID, size, value, order.OrderID)
	} else if order.Side == "SELL" {
		pnl := f.closePosition(order.ProductID, size, value)
		trade.RealizedPnL = pnl - fees
	}

	return trade, nil
}

func (f *FIFOCalculator) addPosition(productID string, size, value float64, orderID string) {
	if f.positions[productID] == nil {
		f.positions[productID] = []Position{}
	}

	avgPrice := value / size
	f.positions[productID] = append(f.positions[productID], Position{
		Size:       size,
		EntryPrice: avgPrice,
		OrderID:    orderID,
	})
}

func (f *FIFOCalculator) closePosition(productID string, sellSize, sellValue float64) float64 {
	if f.positions[productID] == nil || len(f.positions[productID]) == 0 {
		return 0
	}

	queue := f.positions[productID]
	remainingSize := sellSize
	avgSellPrice := sellValue / sellSize
	totalPnL := 0.0

	newQueue := []Position{}

	for remainingSize > 0 && len(queue) > 0 {
		position := queue[0]
		queue = queue[1:]

		matchSize := min(remainingSize, position.Size)

		pnl := (avgSellPrice - position.EntryPrice) * matchSize
		totalPnL += pnl

		position.Size -= matchSize
		remainingSize -= matchSize

		if position.Size > 0 {
			newQueue = append(newQueue, position)
		}
	}

	f.positions[productID] = append(newQueue, queue...)

	return totalPnL
}

func sortOrdersByTime(orders []Order) {
	sort.Slice(orders, func(i, j int) bool {
		ti, _ := time.Parse(time.RFC3339, orders[i].LastFillTime)
		tj, _ := time.Parse(time.RFC3339, orders[j].LastFillTime)
		return ti.Before(tj)
	})
}

func parseFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
