package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	repeat = strings.TrimSpace(repeat)
	if repeat == "" {
		return "", errors.New("empty repeat")
	}

	start, err := time.Parse(DateLayout, dstart)
	if err != nil {
		return "", err
	}

	now = dateOnly(now)
	start = dateOnly(start)

	parts := strings.Fields(repeat)

	if len(parts) == 1 && parts[0] == "y" {
		date := start.AddDate(1, 0, 0)
		for !afterDate(date, now) {
			date = date.AddDate(1, 0, 0)
		}
		return date.Format(DateLayout), nil
	}

	if len(parts) == 2 && parts[0] == "d" {
		interval, err := strconv.Atoi(parts[1])
		if err != nil || interval < 1 || interval > 400 {
			return "", errors.New("invalid day interval")
		}
		date := start.AddDate(0, 0, interval)
		for !afterDate(date, now) {
			date = date.AddDate(0, 0, interval)
		}
		return date.Format(DateLayout), nil
	}

	return "", errors.New("unsupported repeat format")
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeatStr := r.FormValue("repeat")

	var now time.Time
	var err error

	if strings.TrimSpace(nowStr) == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(DateLayout, nowStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	next, err := NextDate(now, dateStr, repeatStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(next))
}

func dateOnly(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func afterDate(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	if ay != by {
		return ay > by
	}
	if am != bm {
		return am > bm
	}
	return ad > bd
}
