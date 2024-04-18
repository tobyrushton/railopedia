package utils

import (
	"strconv"
	"strings"
	"time"
)

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
