package utils

import (
	"errors"
	"strconv"
	"strings"
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
