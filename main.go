package rate

import (
	"math"
	"time"
)

// Wilson will return the wilson-score interval based on the number of upvotes and downvotes
func Wilson(upvotes, downvotes int64) float64 {
	if upvotes == 0 {
		return 0
	}
	n := float64(upvotes + downvotes)
	phat := float64(upvotes) / n
	z := float64(1.96) // for 0.95 confidentiality
	lower := (phat + z*z/(2*n) - z*math.Sqrt(phat*(1-phat)+z*z/(4*n))/n) / (1 + z*z/n)
	return lower
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
	return round(sign*order+float64(seconds)/45000.0, 0.5, 7)
}

func round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}
