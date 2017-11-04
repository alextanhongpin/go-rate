package rate

import (
	"log"
	"testing"
)

func TestLowerBound(t *testing.T) {
	testTable := []struct {
		upvotes   int64
		downvotes int64
	}{
		{596, 18},
		{2360, 98},
		{1, 0},
		{0, 10},
		{1, 10},
	}
	for _, v := range testTable {
		score := Votes(v.upvotes, v.downvotes)
		score2 := Wilson(v.upvotes, v.downvotes)
		log.Println(score, score2)
	}
}
