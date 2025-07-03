package api

import (
	"fmt"
	"time"
	"strings"
	"strconv"
	"errors"
	"net/http"
)

func NextDate(now time.Time, dstart, repeat string) (string, error) {
	data, err := time.Parse("20060102", dstart)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	if repeat != "" {
		spl := strings.Split(repeat, " ")
		if spl[0] == "d" {
			num, err := strconv.Atoi(spl[1])
			if err != nil {
				return "", err
			}
			if num > 400 {
				return "", errors.New("error")

				for {
					data = data.AddDate(0, 0, num)
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
			}
		}
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
	nowStr, err := time.Parse("20060102", now)
	if err != nil || now == "" {
		nowStr = time.Now()
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
