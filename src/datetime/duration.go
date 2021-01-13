package datetime

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type duration int

type Duration interface {
	InMinutes() int
	Add(Duration) Duration
	ToString() string
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func NewDuration(amountHours int, amountMinutes int) Duration {
	return duration(amountHours*60) + duration(amountMinutes)
}

func (d duration) InMinutes() int {
	return int(d)
}

func (d duration) Add(additional Duration) Duration {
	return NewDuration(0, d.InMinutes()+additional.InMinutes())
}

func (d duration) ToString() string {
	if d == 0 {
		return "0m"
	}
	hours := abs(int((int(d) / 60)))
	minutes := abs(int(d) % 60)
	result := ""
	if int(d) < 0 {
		result += "-"
	}
	if hours > 0 {
		result += fmt.Sprintf("%dh", hours)
	}
	if hours > 0 && minutes > 0 {
		result += " "
	}
	if minutes > 0 {
		result += fmt.Sprintf("%dm", minutes)
	}
	return result
}

var durationPattern = regexp.MustCompile(`^\s*(-)?((\d+)h)? *((\d+)m)?\s*$`)

func NewDurationFromString(hhmm string) (Duration, error) {
	match := durationPattern.FindStringSubmatch(hhmm)
	if match == nil {
		return nil, errors.New("MALFORMED_DURATION")
	}
	sign := 1
	if match[1] == "-" {
		sign = -1
	}
	hours, _ := strconv.Atoi(match[3])
	minutes, _ := strconv.Atoi(match[5])
	if minutes > 60 {
		return nil, errors.New("UNREPRESENTABLE_DURATION")
	}
	return NewDuration(sign*hours, sign*minutes), nil
}