package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// GET /api/nextdate?now=...&date=...&repeat=...
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	nowStr := q.Get("now")
	var now time.Time
	var err error
	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(dateLayout, nowStr)
		if err != nil {
			http.Error(w, "invalid now date: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	dstart := q.Get("date")
	if dstart == "" {
		http.Error(w, "missing date parameter", http.StatusBadRequest)
		return
	}

	repeat := q.Get("repeat")
	if strings.TrimSpace(repeat) == "" {
		http.Error(w, "empty repeat rule", http.StatusBadRequest)
		return
	}

	next, err := NextDate(now, dstart, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprint(w, next)
}

// NextDate поддерживает базовые правила d N и y.
func NextDate(now time.Time, dstart, repeat string) (string, error) {
	date, err := time.Parse(dateLayout, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid start date: %w", err)
	}

	afterNow := func(t time.Time) bool {
		y1, m1, d1 := t.Date()
		y2, m2, d2 := now.Date()
		base := time.Date(y2, m2, d2, 0, 0, 0, 0, now.Location())
		cur := time.Date(y1, m1, d1, 0, 0, 0, 0, now.Location())
		return cur.After(base)
	}

	parts := strings.Fields(repeat)
	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid d rule format")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days < 1 || days > 400 {
			return "", fmt.Errorf("invalid d interval")
		}
		for !afterNow(date) {
			date = date.AddDate(0, 0, days)
		}

	case "y":
		if len(parts) != 1 {
			return "", fmt.Errorf("invalid y rule format")
		}
		for !afterNow(date) {
			date = date.AddDate(1, 0, 0)
		}

	default:
		return "", fmt.Errorf("unsupported repeat rule: %q", parts[0])
	}

	return date.Format(dateLayout), nil
}
