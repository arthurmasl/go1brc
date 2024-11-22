package main

import (
	"os"
	"testing"

	"go1brc/internal/solution1"
	"go1brc/internal/solution2"
)

var (
	input, _  = os.Open("./small.txt")
	output, _ = os.ReadFile("./small.out")
)

func Test1BRC(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		result, cities := solution2.Execute(input, small.rows)

		if result != string(output)[:len(output)-1] || cities != small.cities {
			t.Fatalf("test failed")
		}
	})

	defer input.Close()
}

func BenchmarkSolution1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		solution1.Execute(input, small.rows)
	}
}

func BenchmarkSolution2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		solution2.Execute(input, small.rows)
	}
}
