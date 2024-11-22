package main

import (
	"os"
	"runtime/pprof"

	"go1brc/internal/executor"
)

const PROFILE = true

func main() {
	if PROFILE {
		f, err := os.Create("cpu_profile.prof")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if err := pprof.StartCPUProfile(f); err != nil {
			panic(err)
		}
		defer pprof.StopCPUProfile()
	}

	executor.ExecuteSolution(executor.SolutionCase)
}
