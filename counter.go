package main

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type StringWithID struct {
	ID   int    `json:"id" validate:"required,gte=1"`
	Data string `json:"data" validate:"required"`
}

type CountWithID struct {
	ID     int `json:"id"`
	Result int `json:"result"`
}

type CountWithIDAndError struct {
	Count *CountWithID
	Err   error
}

var re = regexp.MustCompile(`^\d+(?:,\d+)*$`)

func CountIncreases(data string) (int, error) {
	correct := re.MatchString(data)
	if !correct {
		return -1, errors.New("invalid input")
	}

	count := 0
	var prev int
	for i, v := range strings.Split(data, ",") {
		int_v, err := strconv.Atoi(v)
		if err != nil {
			return -1, err
		}
		if i > 0 && int_v > prev {
			count++
		}
		prev = int_v
	}

	return count, nil
}

func CountIncreases2(stringWithID *StringWithID, co chan *CountWithIDAndError) {
	count := 0
	id, data := stringWithID.ID, stringWithID.Data
	count, err := CountIncreases(data)
	if err != nil {
		co <- &CountWithIDAndError{nil, err}
	}
	co <- &CountWithIDAndError{&CountWithID{id, count}, nil}
}
