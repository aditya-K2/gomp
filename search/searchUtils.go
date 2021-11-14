package search

import (
	"fmt"
	"math"
	"math/bits"
	"strings"
)

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

func ConvertToBitMap(b string) uint64 {
	b = strings.ToLower(b)
	var a uint64 = 0
	for _, i := range b {
		asc := (int(i) - 97)
		if asc >= 97 && asc <= 122 {
			a |= 1 << asc
		}
	}
	return a
}
func CalculateUpperBound(a, b string) float64 {
	aB := ConvertToBitMap(a)
	bB := ConvertToBitMap(b)
	lengthAf64 := float64(len(a))
	lengthBf64 := float64(len(b))
	m := math.Min(lengthAf64, lengthBf64) - float64(bits.OnesCount64(aB&(^bB)))
	return (1.0 / 3.0) * ((float64(m) / float64(len(a))) + (float64(m) / float64(len(b))) + 1)
}

func min(a, b, c int) int {
	if a > b && b < c {
		return b
	} else if b > a && a < c {
		return a
	} else {
		return c
	}
}

func cost(i, j rune) int {
	if i != j {
		return 1
	} else {
		return 0
	}
}

func GetLevenshteinDistance(a, b string) int {
	c := []rune(a)
	e := []rune(b)
	m, n := len(c), len(e)
	d := make([][]int, m+1)
	for i := range d {
		d[i] = make([]int, n+1)
	}
	for i := range d {
		for j := range d[i] {
			if j == 0 {
				d[i][j] = i
			} else if i == 0 {
				d[i][j] = j
			} else {
				d[i][j] = 0
			}
		}
	}
	for j := 1; j < n+1; j++ {
		for i := 1; i < m+1; i++ {
			d[i][j] = min(d[i-1][j-1]+cost(c[i-1], e[j-1]),
				d[i][j-1]+1,
				d[i-1][j]+1)
		}
	}
	return d[m][n]
}

func main() {
	fmt.Println(GetLevenshteinDistance("cat", "wildcat"))
}
