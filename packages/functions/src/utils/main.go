package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func RoundToNextFive(num int) int {
	for num%5 != 0 {
		num++
	}
	return num
}

func PriceToFloat(price string) float32 {
	fmt.Println(price)
	price64, _ := strconv.ParseFloat(price, 32)
	return float32(price64)
}

func SanitisePrice(price string) string {
	temp := strings.Replace(price, "£", "", -1)
	return strings.Replace(temp, "+", "", -1)
}
