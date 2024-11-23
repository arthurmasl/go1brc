package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"

	"go1brc/internal/executor"
)

func main() {
	profile := flag.Bool("profile", false, "start profiling")
	flag.Parse()

	if *profile {
		fmt.Println("Profiling...")
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
