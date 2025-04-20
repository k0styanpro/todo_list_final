package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// nextDateHandler обрабатывает GET /api/nextdate?now=...&date=...&repeat=...
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	// Парсим параметр now или берём текущее время
	nowStr := q.Get("now")
	var now time.Time
	var err error
	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(dateLayout, nowStr)
		if err != nil {
			http.Error(w, "invalid now: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	// Обязательный параметр date
	dstart := q.Get("date")
	if dstart == "" {
		http.Error(w, "missing date", http.StatusBadRequest)
		return
	}

	// Повторение может быть пустым — тогда просто возвращаем пустую строку
	repeat := q.Get("repeat")
	if strings.TrimSpace(repeat) == "" {
		fmt.Fprint(w, "")
		return
	}

	// Иначе вычисляем следующую дату
	next, err := NextDate(now, dstart, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprint(w, next)
}

// NextDate возвращает дату следующего выполнения задачи (> now) в формате YYYYMMDD.
// Поддерживаются правила "d N" и "y".
func NextDate(now time.Time, dstart, repeat string) (string, error) {
	start, err := time.Parse(dateLayout, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid start date: %w", err)
	}
	// Сравниваем только по календарным датам
	base := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	parts := strings.Fields(repeat)
	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid d format")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days < 1 || days > 400 {
			return "", fmt.Errorf("invalid d interval")
		}
		next := start.AddDate(0, 0, days)
		for !next.After(base) {
			next = next.AddDate(0, 0, days)
		}
		return next.Format(dateLayout), nil

	case "y":
		if len(parts) != 1 {
			return "", fmt.Errorf("invalid y format")
		}
		next := start.AddDate(1, 0, 0)
		for !next.After(base) {
			next = next.AddDate(1, 0, 0)
		}
		return next.Format(dateLayout), nil

	default:
		return "", fmt.Errorf("unsupported repeat: %q", parts[0])
	}
}
