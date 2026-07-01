package api

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/UGRORF/go-todo/pkg/db"
)

const (
	dateFormat     = "20060102"
	maxDaysInWeek  = 7
	maxDaysInMonth = 31
	maxDaysForD    = 400
	monthsInYear   = 12
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", nil
	}

	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", err
	}

	s := strings.Split(repeat, " ")

	switch s[0] {
	case "d":
		if len(s) != 2 {
			return "", errors.New("invalid repeat, must be in format 'd n', where n - count days for next date")
		}
		days, err := strconv.Atoi(s[1])
		if err != nil {
			return "", err
		}
		if days < 1 || days > maxDaysForD {
			return "", errors.New("invalid repeat, must be in range 0-400")
		}

		return advanceByInterval(date, now, 0, 0, days), nil

	case "y":
		if len(s) != 1 {
			return "", errors.New("invalid repeat, must be in format 'y'")
		}

		return advanceByInterval(date, now, 1, 0, 0), nil

	case "w":
		if len(s) != 2 {
			return "", errors.New("invalid repeat, must be in format 'w n,m', where n,m - number days for the week")
		}

		numberDayOfWeeks := strings.Split(s[1], ",")
		return findNextWeekDay(date, now, numberDayOfWeeks)

	case "m":
		if len(s) < 2 || len(s) > 3 {
			return "", errors.New("invalid repeat, must be in format 'm n,p [x,y]', where n,p - number days of the month, x,y - number of month (optional)")
		}
		numberDayOfMonths := strings.Split(s[1], ",")
		var numberMonthsOfYear []string
		if len(s) == 3 {
			numberMonthsOfYear = strings.Split(s[2], ",")
		}

		return findNextMonthDay(date, now, numberDayOfMonths, numberMonthsOfYear)

	default:
		return "", fmt.Errorf("unsupported repeat format: %s", s[0])
	}
}

func findNextWeekDay(date time.Time, now time.Time, numberDayOfWeeks []string) (string, error) {
	if len(numberDayOfWeeks) == 0 {
		return "", fmt.Errorf("no numberDayOfWeeks provided")
	}
	var allowedWeekDays [maxDaysInWeek + 1]bool
	for _, day := range numberDayOfWeeks {
		dayIndex, err := strconv.Atoi(strings.TrimSpace(day))
		if err != nil {
			return "", err
		}
		if dayIndex < 1 || dayIndex > maxDaysInWeek {
			return "", fmt.Errorf("invalid numberDayOfWeeks: %s", day)
		}

		var goDay time.Weekday
		if dayIndex == 7 {
			goDay = time.Sunday
		} else {
			goDay = time.Weekday(dayIndex)
		}
		allowedWeekDays[goDay] = true
	}
	date = date.AddDate(0, 0, 1)
	for !date.After(now) || !allowedWeekDays[date.Weekday()] {
		date = date.AddDate(0, 0, 1)
	}

	return date.Format(dateFormat), nil
}

func findNextMonthDay(date, now time.Time, numberDayOfMonths, numberMonthsOfYear []string) (string, error) {
	if len(numberDayOfMonths) == 0 {
		return "", fmt.Errorf("numberDayOfMonths is empty")
	}

	allowedDayOfMonth := make(map[int]bool)
	allowedMonthOfYear := make(map[int]bool)

	if len(numberMonthsOfYear) == 0 {
		for i := 1; i <= monthsInYear; i++ {
			allowedMonthOfYear[i] = true
		}
	} else {
		for i := range numberMonthsOfYear {
			month, err := strconv.Atoi(numberMonthsOfYear[i])
			if err != nil {
				return "", err
			}
			if month < 1 || month > monthsInYear {
				return "", fmt.Errorf("month must be between 1 and %d", monthsInYear)
			}
			allowedMonthOfYear[month] = true
		}
	}

	for i := range numberDayOfMonths {
		day, err := strconv.Atoi(numberDayOfMonths[i])
		if err != nil {
			return "", err
		}
		if day < -2 || day > maxDaysInMonth || day == 0 {
			return "", fmt.Errorf("day must be between -2 to %d", maxDaysInMonth)
		}

		allowedDayOfMonth[day] = true
	}
	if len(allowedDayOfMonth) == 0 {
		return "", fmt.Errorf("no valid days provided")
	}

	if len(allowedMonthOfYear) == 0 {
		return "", fmt.Errorf("no valid months provided")
	}

	for {
		date = date.AddDate(0, 0, 1)

		if !date.After(now) {
			continue
		}

		if !allowedMonthOfYear[int(date.Month())] {
			continue
		}
		day := date.Day()
		daysInCurrentMonth := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()

		isAllowed := allowedDayOfMonth[date.Day()] || (allowedDayOfMonth[-1] && day == daysInCurrentMonth) || (allowedDayOfMonth[-2] && day == daysInCurrentMonth-1)
		if isAllowed {
			return date.Format(dateFormat), nil
		}
	}
}

func advanceByInterval(date, now time.Time, years, months, days int) string {
	date = date.AddDate(years, months, days)
	for !date.After(now) {
		date = date.AddDate(years, months, days)
	}
	return date.Format(dateFormat)
}

func CheckDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(dateFormat)
		return nil
	}

	t, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("invalid date format: must be YYYYMMDD")
	}

	if t.After(now) || t.Format(dateFormat) == now.Format(dateFormat) {
		return nil
	}

	if task.Repeat != "" {
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("invalid repeat rule: %w", err)
		}
		if next == "" {
			return fmt.Errorf("failed to calculate next date")
		}
		task.Date = next
		return nil
	}

	task.Date = now.Format(dateFormat)
	return nil
}
