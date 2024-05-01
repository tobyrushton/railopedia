package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var ISO8601Layout = "2006-01-02T15:04:05Z0700"

func RoundToNextFive(num int) int {
	for num%5 != 0 {
		num++
	}
	return num
}

func PriceToFloat(price string) float32 {
	price64, _ := strconv.ParseFloat(price, 32)
	return float32(price64)
}

func SanitisePrice(price string) string {
	temp := strings.Replace(price, "£", "", -1)
	return strings.Replace(temp, "+", "", -1)
}

func GetTime(timeString string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05Z0700", timeString)
}

func RoundToNextFifteen(num int) int {
	for num%15 != 0 {
		num++
	}
	return num
}

func SplitString(str string, key string) (string, string) {
	split := strings.Split(str, key)
	return split[0], split[1]
}

func ParseTime(journeyTime time.Time, requestedTime time.Time) time.Time {
	if journeyTime.Before(requestedTime) {
		return journeyTime.AddDate(0, 0, 1)
	}
	return journeyTime
}

func HourStringToISO(hour string, day time.Time) string {
	timeValue := fmt.Sprintf("%d-%02d-%02d %s:%02d", day.Year(), day.Month(), day.Day(), hour, 0)
	departTime, _ := time.Parse("2006-01-02 15:04:05", timeValue)
	ISOTime := ParseTime(departTime, day).Format(ISO8601Layout)
	return ISOTime
}
