package rrcf

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestCreation(t *testing.T) {
	src := rand.NewSource(0)
	rng := rand.New(src)
	d := 3
	x := NewMatrix(d)
	for i := 0; i < 100; i++ {
		row := make([]float64, d)
		for i := 0; i < d; i++ {
			row[i] = rng.Float64()
		}
		x.Data = append(x.Data, row)
	}
	r := New(x)
	fmt.Println(r)
}
