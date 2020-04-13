package main

import (
	"fmt"
	"math/rand"

	"github.com/a-h/rrcf"
)

func main() {
	src := rand.NewSource(0)
	rng := rand.New(src)
	d := 3
	x := rrcf.NewMatrix(d)
	for i := 0; i < 100; i++ {
		row := make([]float64, d)
		for i := 0; i < d; i++ {
			row[i] = rng.Float64()
		}
		x.Data = append(x.Data, row)
	}
	r := rrcf.New(x)
	fmt.Println(r.String())
}

// X = np.random.randn(100, 2)
// tree = rrcf.RCTree(X)
// print(tree)
