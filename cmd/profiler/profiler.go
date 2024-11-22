package profiler

import (
	"go1brc/internal/executor"
	"os"
	"runtime/pprof"
)

func main() {
	f, err := os.Create("cpu_profile.prof")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := pprof.StartCPUProfile(f); err != nil {
		panic(err)
	}
	defer pprof.StopCPUProfile()

	executor.ExecuteSolution(executor.SolutionCase)
}
