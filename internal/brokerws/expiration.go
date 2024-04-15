package brokerws

import (
	"time"
)

func GetExpirationTime(timestamp int64, duration int) (int, int) {
	now := time.Unix(timestamp, 0)
	expDate := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location())

	if int(expDate.Add(time.Minute).Unix())-int(timestamp) > 30 {
		expDate = expDate.Add(time.Minute)
	} else {
		expDate = expDate.Add(2 * time.Minute)
	}

	var exp []int64

	for i := 0; i < 5; i++ {
		exp = append(exp, expDate.Unix())
		expDate = expDate.Add(time.Minute)
	}

	idx := 50
	index := 0
	now = time.Unix(timestamp, 0)
	expDate = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location())

	for index < idx {
		if expDate.Minute()%15 == 0 && (int(expDate.Unix())-int(timestamp)) > 60*5 {
			exp = append(exp, expDate.Unix())
			index++
		}
		expDate = expDate.Add(time.Minute)
	}

	var remaining []int

	for _, t := range exp {
		remaining = append(remaining, int(t)-int(time.Now().Unix()))
	}

	close := make([]int, len(remaining))
	for i, x := range remaining {
		close[i] = abs(x - 60*duration)
	}

	minClose := min(close)
	return int(exp[getIndex(close, minClose)]), int(getIndex(close, minClose))
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func min(arr []int) int {
	if len(arr) == 0 {
		return 0
	}

	minVal := arr[0]
	for _, val := range arr {
		if val < minVal {
			minVal = val
		}
	}
	return minVal
}

func getIndex(arr []int, val int) int {
	for i, v := range arr {
		if v == val {
			return i
		}
	}
	return -1
}
