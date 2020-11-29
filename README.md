# Ranking Algorithms

[![](https://godoc.org/github.com/alextanhongpin/go-rate?status.svg)](https://godoc.org/github.com/alextanhongpin/go-rate)

## Installation

```
$ go get github.com/alextanhongpin/go-rate
```

## Run

```go
func main () {
  upvotes := 1000
  downvotes := 10
  score := rate.Wilson(upvotes, downvotes)

  // or
  createdAt := time.Now()
  score := rate.Hot(upvotes, downvotes, createdAt)
}
```

## Sorting by Wilson-Score Interval

  Reddit comment ranking is using the Wilson-Score Interval

Wilson-Score Interval formula is displayed below:

![Wilson Score Interval](/assets/wilson-score-interval.png)

- `p-hat` is the fraction of positive votes out of the total votes
- `n` is the total number of upvotes and downvotes

Here's an example of the algorithm written in go:

```golang
func Wilson(upvotes, downvotes int64) float64 {
	n := float64(upvotes + downvotes)
  // if n == 0.0 { return 0 } will return false results
	if upvotes == 0 {
		return 0
	}
	phat := float64(upvotes) / n
	z := float64(1.96) // for 0.95 confidentiality
	lower := (phat + z*z/(2*n) - z*math.Sqrt((phat*(1-phat)+z*z/(4*n))/n)) / (1 + z*z/n)
	return lower
}
```

If upvotes is zero the score will be `0`. Note that if you implement the logic where `upvotes + downvotes = 0`, you might face the issue below:

```
1. upvotes=0 downvotes=100 score=0.016648
2. upvotes=0 downvotes=1 score=0.000000
3. upvotes=0 downvotes=0 score=0.000000
```



The item with 100 downvotes should be placed lower then the items with zero vites. But since the range is only from 0 to 1, there is no such thing as negative scores.

To fix this, we added another additional logic - swap the upvotes with downvotes, and turn the value to negative:

```go
	if upvotes == 0 {
		return -Wilson(downvotes, upvotes)
	}
```
Once corrected, we get the following:

```
1. upvotes=0 downvotes=10 score=-0.722460
2. upvotes=0 downvotes=1 score=-0.206543
3. upvotes=0 downvotes=0 score=0.000000
```

To test how the votes will affect the score, run the `Wilson` algorithm against different values of upvotes and downvotes.

Sample code:
```go
package main

import (
	"fmt"
	"math"
	"sort"
	"time"
)

type Result struct {
	Upvotes   int64
	Downvotes int64
	Score     float64
}

func main() {
	fmt.Println("Wilson Score Ranking")
	fmt.Println("")

	votes := []struct {
		upvotes   int64
		downvotes int64
	}{
		{0, 0},
		{1, 0},
		{0, 1},
		{1, 1},
		{2, 1},
		{3, 1},
		{3, 2},
		{10, 0},
		{0, 10},
		{100, 10},
		{100, 100},
		{0, 100},
	}
	result := make([]Result, len(votes))
	for i, v := range votes {
		score := Wilson(v.upvotes, v.downvotes)
		result[i] = Result{
			Upvotes:   v.upvotes,
			Downvotes: v.downvotes,
			Score:     score,
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Score > result[j].Score
	})
	for idx, r := range result {
		fmt.Printf("%2d. up=%-3d down=%-3d score=%+f\n", idx+1, r.Upvotes, r.Downvotes, r.Score)
	}
}
```

Output:

```
Wilson Score Ranking

 1. up=100 down=10  score=+0.840702
 2. up=10  down=0   score=+0.722460
 3. up=100 down=100 score=+0.431360
 4. up=3   down=1   score=+0.300636
 5. up=3   down=2   score=+0.230720
 6. up=2   down=1   score=+0.207655
 7. up=1   down=0   score=+0.206543
 8. up=1   down=1   score=+0.094529
 9. up=0   down=0   score=+0.000000
10. up=0   down=1   score=-0.206543
11. up=0   down=10  score=-0.722460
12. up=0   down=100 score=-0.963005
```

Here we can conclude several things:

- An item with **1000 upvotes** and **1000 downvotes** is ranked higher than an item with **100 upvotes** and **100 downvotes** - since it has more votes
- An item with **0 upvote** will always have a score of zero
- An item with **100 upvotes** and **0 downvote** will be placed above an item with **100 upvotes** and **100 downvotes**


## Hot Ranking Algorithm

Hot Ranking algorithm is described below:

![Hot Ranking](/assets/hot.png)

This is the same algorithm that is used by Reddit to rank their stories. It takes the account of submission time into the ranking. What this means is:

1. Newer stories will be ranked higher than older
2. The score won't decrease as time goes by, but newer stories will get a higher score than older

The equivalent code written in go:

```golang
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
	seconds := (date.UnixNano() / 1e6 - epoch) / 1e3  - 1134028003
	return round(sign * order + float64(seconds) / 45000.0, 0.5, 7)
}
```

To see the effect of the submission date:
```
log.Println("January 2017:", Hot(1000, 10, time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC)))
log.Println("January 2016:", Hot(1000, 10, time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC)))
```

Output:

```bash
2017/10/30 00:05:49 January 2017: 7763.0133463
2017/10/30 00:05:49 January 2016: 7060.2933463
```

The score returned from the recent submission (January 2017) is higher than that of the one from a year ago.


## References:
1. [How Not To Sort By Average Rating](http://www.evanmiller.org/how-not-to-sort-by-average-rating.html)
2. [How Reddit Ranking Algorithm Works](https://medium.com/hacking-and-gonzo/how-reddit-ranking-algorithms-work-ef111e33d0d9)
