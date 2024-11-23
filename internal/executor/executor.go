package executor

import (
	"fmt"
	"os"
	"time"

	"go1brc/internal/solution1"
	"go1brc/internal/solution2"
	"go1brc/internal/solution3"
)

type Solution struct {
	name   string
	rows   int
	cities int
}

var (
	small   = Solution{rows: 12, name: "small", cities: 10}
	tenmil  = Solution{rows: 10_000_000, name: "tenmils", cities: 413}
	billion = Solution{rows: 1_000_000_000, name: "billion", cities: 413}
)

var (
	SolutionFn   = solution3.Execute
	SolutionCase = small
)

func ExecuteSolution(solution Solution) (string, int) {
	fmt.Fprintln(os.Stdout, []any{solution1.Execute, solution2.Execute, solution3.Execute}...)

	file, _ := os.Open(fmt.Sprintf("resources/%v.txt", solution.name))
	defer file.Close()

	fmt.Printf("Name: %v, Rows: %v\n", solution.name, solution.rows)

	t1 := time.Now()
	str, cities := SolutionFn(file, solution.rows)
	if cities != solution.cities {
		fmt.Printf("-Wrong solution!, got %v cities, when expected %v\n", cities, solution.cities)
	}
	fmt.Printf("Cities: %v, Expected: %v\n", cities, solution.cities)
	fmt.Printf("Execution time: %v\n", time.Since(t1))

	return str, cities
}
