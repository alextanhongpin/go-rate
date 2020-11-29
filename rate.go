package rate

import (
	"math"
	"time"
)

// Wilson will return the wilson-score interval based on the number of upvotes and downvotes
func Wilson(upvotes, downvotes int64) float64 {
	if upvotes == 0 && downvotes == 0 {
		return 0.0
	}
	if upvotes == 0 {
		return -Wilson(downvotes, upvotes)
	}
	n := float64(upvotes + downvotes)
	phat := float64(upvotes) / n
	z := float64(1.96) // for 0.95 confidentiality
	lower := (phat + z*z/(2*n) - z*math.Sqrt((phat*(1-phat)+z*z/(4*n))/n)) / (1 + z*z/n)
	return lower
}

// Votes returns the score from sorting (which includes negative scores too)
// http://nbviewer.jupyter.org/github/CamDavidsonPilon/Probabilistic-Programming-and-Bayesian-Methods-for-Hackers/blob/master/Chapter4_TheGreatestTheoremNeverTold/Ch4_LawOfLargeNumbers_PyMC3.ipynb
func Votes(upvotes, downvotes int64) float64 {
	a := float64(1 + upvotes)
	b := float64(1 + downvotes)
	mu := a / (a + b)
	stdErr := 1.65 * math.Sqrt((a*b)/(math.Pow(a+b, 2)*(a+b+1)))
	return mu - stdErr
}

// Stars return the scores for star ratings
func Stars(n, s int64) float64 {
	// s is sum of all the ratings
	// n is the number of users who rated
	a := float64(1 + s)
	b := float64(1 + n - s)
	mu := a / (a + b)
	stdErr := 1.65 * math.Sqrt((a*b)/(math.Pow(a+b, 2)*(a+b+1)))
	return mu - stdErr
}

// Hot will return ranking based on the time it is created
func Hot(upvotes, downvotes int64, date time.Time) float64 {
	s := float64(upvotes - downvotes)
	order := math.Log10(math.Max(math.Abs(s), 1))
	var sign float64
	if s > 0 {
		sign = 1.0
	} else if s < 0 {
		sign = -1.0
	} else {
		sign = 0.0
	}
	epoch := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC).UnixNano() / 1e6
	// epoch_seconds := time.Date(1970,1,14, 3, 0, 28, 3e6, time.UTC).UnixNano() / 1e6
	seconds := (date.UnixNano()/1e6-epoch)/1e3 - 1134028003
	return sign*order + float64(seconds)/45000.0
}
