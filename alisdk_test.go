package alisdk

import (
	"strings"
	"testing"
)

func TestRandomUuid(t *testing.T) {
	s := ""
	for i := 0; i < 1000; i++ {
		u := randomUuid()
		if strings.Contains(s, u) {
			t.Fatal("100个UUID有重复的")
		}
		s += "|" + u + "|"
	}
}
