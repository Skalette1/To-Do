package utils

import (
	"net/http"
	"strconv"
	"strings"
	"errors"
	"time"
)

func NextDate(now time.Time, dstart, repeat string) (string, error) {
	data, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", err
	}

	if repeat == "" {
		return "", errors.New("repeat is empty")
	}

	spl := strings.Split(repeat, " ")
	if len(spl) == 0 {
		return "", errors.New("invalid repeat format")
	}

	if spl[0] == "d" {
		if len(spl) < 2 {
			return "", errors.New("invalid repeat format for daily rule")
		}
		num, err := strconv.Atoi(spl[1])
		if err != nil {
			return "", err
		}
		if num > 400 {
			return "", errors.New("interval too large")
		}
		for {
			data = data.AddDate(0, 0, num)
			if data.After(now) {
				break
			}
		}
	} else if spl[0] == "w" {
		if len(spl) < 2 {
			return "", errors.New("invalid repeat format for weekly rule")
		}
		num, err := strconv.Atoi(spl[1])
		if err != nil {
			return "", err
		}
		if num > 400 {
			return "", errors.New("interval too large")
		}
		for {
			data = data.AddDate(0, 0, num*7)
			if data.After(now) {
				break
			}
		}
	} else if spl[0] == "y" {
		for {
			data = data.AddDate(1, 0, 0)
			if data.After(now) {
				break
			}
		}
	} else {
		return "", errors.New("unsupported repeat format")
	}

	return data.Format("20060102"), nil
}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	if date == "" || repeat == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
		return
	}

	var nowStr time.Time
	if now == "" {
		nowStr = time.Now()
	} else {
		var err error
		nowStr, err = time.Parse("20060102", now)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid now date format"))
			return
		}
	}

	response, err := NextDate(nowStr, date, repeat)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server Error"))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(response))
}
