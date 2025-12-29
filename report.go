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

	startDate := getFirstSunday(firstDay)
	endDate := getFirstSaturday(lastDay)

	days := generateDays(startDate, endDate, dailyTrades, year, month)
	weeks := generateWeeks(days, year, month)
	monthlyPnL := calculateMonthlyPnL(days, year, month)

	return MonthData{
		Month:      month.String(),
		Year:       year,
		MonthlyPnL: roundToTwoDecimals(monthlyPnL),
		Days:       days,
		Weeks:      weeks,
	}
}

func getFirstSunday(date time.Time) time.Time {
	for date.Weekday() != time.Sunday {
		date = date.AddDate(0, 0, -1)
	}
	return date
}

func getFirstSaturday(date time.Time) time.Time {
	for date.Weekday() != time.Saturday {
		date = date.AddDate(0, 0, 1)
	}
	return date
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
			DayOfMonth: current.Day(),
			DayOfWeek:  int(current.Weekday()),
			PnL:        roundToTwoDecimals(pnl),
			TradeCount: len(trades),
			HasNotes:   false,
		})
	}

	return days
}

func generateWeeks(days []DayPnL, targetYear int, targetMonth time.Month) []WeekPnL {
	if len(days) == 0 {
		return []WeekPnL{}
	}

	var weeks []WeekPnL
	weekNum := 1
	weekPnL := 0.0
	weekTradeCount := 0

	for i, day := range days {
		dayDate, _ := time.Parse("2006-01-02", day.Date)

		if dayDate.Year() == targetYear && dayDate.Month() == targetMonth {
			weekPnL += day.PnL
			weekTradeCount += day.TradeCount
		}

		if day.DayOfWeek == 6 || i == len(days)-1 {
			weeks = append(weeks, WeekPnL{
				WeekNumber:     weekNum,
				WeekPnL:        roundToTwoDecimals(weekPnL),
				WeekTradeCount: weekTradeCount,
			})
			weekNum++
			weekPnL = 0.0
			weekTradeCount = 0
		}
	}

	return weeks
}

func calculateMonthlyPnL(days []DayPnL, year int, month time.Month) float64 {
	total := 0.0
	for _, day := range days {
		dayDate, _ := time.Parse("2006-01-02", day.Date)
		if dayDate.Year() == year && dayDate.Month() == month {
			total += day.PnL
		}
	}
	return total
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
