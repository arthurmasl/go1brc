package main

import (
	"flag"

	"go1brc/internal/executor"
	"go1brc/internal/utils"
)

func main() {
	profile := flag.Bool("profile", false, "start profiling")
	flag.Parse()

	if *profile {
		stopProfiling := utils.Profile()
		defer stopProfiling()
	}

	executor.ExecuteSolution(executor.SolutionCase)
}
