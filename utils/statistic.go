package utils

import (
	"math"
	"math/big"
	"sort"
)

func Median(values []*big.Float) *big.Float {
	n := len(values)
	if n == 0 {
		return nil
	}

	sort.Slice(values, func(i, j int) bool {
		return values[i].Cmp(values[j]) < 0
	})

	mid := n / 2
	if n%2 == 1 {
		return values[mid]
	}

	sum := new(big.Float).Add(values[mid-1], values[mid])
	return sum.Quo(sum, big.NewFloat(2))
}

func CheckMedian(
	price *big.Float,
	median *big.Float,
	cutoffPercent *big.Float,
) bool {
	bottom := new(big.Float).Sub(median, cutoffPercent)
	top := new(big.Float).Add(median, cutoffPercent)
	return price.Cmp(bottom) >= 0 && price.Cmp(top) <= 0
}

func Average(values []*big.Float) *big.Float {
	sum := big.NewFloat(0)
	count := 0

	for _, v := range values {
		if v != nil {
			sum.Add(sum, v)
			count++
		}
	}

	if count == 0 {
		return big.NewFloat(0)
	}

	avg := new(big.Float).Quo(sum, big.NewFloat(float64(count)))
	i := new(big.Int)
	avg.Int(i)
	return new(big.Float).SetInt(i)
}

func StandardDeviation(values []*big.Float) *big.Float {
	var usable []*big.Float
	for _, v := range values {
		if v != nil {
			usable = append(usable, v)
		}
	}

	n := len(usable)
	if n <= 1 {
		return nil
	}

	sum := big.NewFloat(0)
	for _, v := range usable {
		sum.Add(sum, v)
	}
	mean := new(big.Float).Quo(sum, big.NewFloat(float64(n)))

	varianceSum := big.NewFloat(0)
	for _, v := range usable {
		diff := new(big.Float).Sub(v, mean)
		squared := new(big.Float).Mul(diff, diff)
		varianceSum.Add(varianceSum, squared)
	}

	variance := new(big.Float).Quo(varianceSum, big.NewFloat(float64(n-1)))

	f64, _ := variance.Float64()
	std := math.Sqrt(f64)
	i := new(big.Int)
	s := big.NewFloat(std)
	s.Int(i)
	return new(big.Float).SetInt(i)
}
