package util

import (
	"testing"
	"time"
)

func BenchmarkEndOfDay1(t *testing.B) {
	for i := 0; i < t.N; i++ {
		_ = EndOfDay(time.Now())
	}

}
func BenchmarkEndOfDay2(t *testing.B) {

	for i := 0; i < t.N; i++ {
		_ = EndOfDay2(time.Now())
	}
}
