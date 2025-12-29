package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

func GenerateReport(trades []Trade) (Report, error) {
	dailyTrades := groupTradesByDate(trades)
	months := extractMonths(trades)

	report := make(Report)

	for _, monthKey := range months {
		year, month := parseMonthKey(monthKey)
		monthData := generateMonthData(year, month, dailyTrades)
		report[monthKey] = monthData
	}

	return report, nil
}

func groupTradesByDate(trades []Trade) map[string][]Trade {
	grouped := make(map[string][]Trade)
	for _, trade := range trades {
		dateKey := trade.Date.Format("2006-01-02")
		grouped[dateKey] = append(grouped[dateKey], trade)
	}
	return grouped
}

func extractMonths(trades []Trade) []string {
	monthSet := make(map[string]bool)
	for _, trade := range trades {
		monthKey := trade.Date.Format("2006-01")
		monthSet[monthKey] = true
	}

	months := make([]string, 0, len(monthSet))
	for month := range monthSet {
		months = append(months, month)
	}

	sort.Slice(months, func(i, j int) bool {
		return months[i] > months[j]
	})

	return months
}

func parseMonthKey(monthKey string) (int, time.Month) {
	t, _ := time.Parse("2006-01", monthKey)
	return t.Year(), t.Month()
}

func generateMonthData(year int, month time.Month, dailyTrades map[string][]Trade) MonthData {
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDay.AddDate(0, 1, -1)

	days := generateDays(firstDay, lastDay, dailyTrades, year, month)

	return MonthData{
		Month: month.String(),
		Year:  year,
		Days:  days,
	}
}

func generateDays(start, end time.Time, dailyTrades map[string][]Trade, targetYear int, targetMonth time.Month) []DayPnL {
	var days []DayPnL

	for current := start; !current.After(end); current = current.AddDate(0, 0, 1) {
		dateKey := current.Format("2006-01-02")
		trades := dailyTrades[dateKey]

		pnl := 0.0
		for _, trade := range trades {
			pnl += trade.RealizedPnL
		}

		days = append(days, DayPnL{
			Date:       dateKey,
			PnL:        roundToTwoDecimals(pnl),
			TradeCount: len(trades),
			HasNotes:   false,
		})
	}

	return days
}


func roundToTwoDecimals(value float64) float64 {
	return float64(int(value*100+0.5)) / 100
}

func WriteReport(report Report, filename string) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling report: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	return nil
}
