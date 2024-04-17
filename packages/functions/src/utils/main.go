package utils

func RoundToNextFive(num int) int {
	for num%5 != 0 {
		num++
	}
	return num
}
