package util

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func RandomNumeric(size uint8) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if size <= 0 {
		panic(fmt.Sprintf("{ size : %d } must be more than 0", size))
	}
	value := strconv.Itoa(r.Intn(10))
	for ; size > 1; size-- {
		value += strconv.Itoa(r.Intn(10))
	}

	return value
}

func EndOfDay(t time.Time) time.Time {
	return t.Truncate(time.Hour * 24).Add(time.Hour*24 - 1)
}

func EndOfDay2(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 0, t.Location())
}
