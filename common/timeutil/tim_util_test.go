package timeutil

import (
	"math/rand"
	"testing"
)

func Test_IntToWeekSartAndEnd(t *testing.T) {
	for i := 0; i < 100; i++ {
		s, e := IntToWeekSartAndEnd(201700 + rand.Int63n(52))
		t.Logf("res:%d,%d", s, e)
	}
}
