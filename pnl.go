package main

import (
	"fmt"
	"sort"
	"strconv"
	"time"
)

type FIFOCalculator struct {
	positions map[string]*PositionQueue
}

func NewFIFOCalculator() *FIFOCalculator {
	return &FIFOCalculator{
		positions: make(map[string]*PositionQueue),
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
	price := parseFloat(order.AverageFilledPrice)
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
		Price:     price,
		Fees:      fees,
		OrderID:   order.OrderID,
	}

	if order.Side == "BUY" {
		trade.RealizedPnL = -fees
		f.addPosition(order.ProductID, size, price, order.OrderID)
	} else if order.Side == "SELL" {
		pnl := f.matchFIFO(order.ProductID, size, price)
		trade.RealizedPnL = pnl - fees
	}

	return trade, nil
}

func (f *FIFOCalculator) addPosition(productID string, size, price float64, orderID string) {
	if f.positions[productID] == nil {
		f.positions[productID] = &PositionQueue{
			ProductID: productID,
			Queue:     []Position{},
		}
	}

	f.positions[productID].Queue = append(f.positions[productID].Queue, Position{
		Size:       size,
		EntryPrice: price,
		OrderID:    orderID,
	})
}

func (f *FIFOCalculator) matchFIFO(productID string, sellSize, sellPrice float64) float64 {
	if f.positions[productID] == nil {
		f.positions[productID] = &PositionQueue{
			ProductID: productID,
			Queue:     []Position{},
		}
	}

	queue := &f.positions[productID].Queue
	remainingSize := sellSize
	totalPnL := 0.0

	for remainingSize > 0 && len(*queue) > 0 {
		position := &(*queue)[0]

		matchSize := min(remainingSize, position.Size)
		pnl := (sellPrice - position.EntryPrice) * matchSize
		totalPnL += pnl

		position.Size -= matchSize
		remainingSize -= matchSize

		if position.Size == 0 {
			*queue = (*queue)[1:]
		}
	}

	if remainingSize > 0 {
		f.addPosition(productID, -remainingSize, sellPrice, "")
	}

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
