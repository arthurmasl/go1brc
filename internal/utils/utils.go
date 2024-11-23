package utils

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"
)

func Profile() func() {
	fmt.Println("Profiling...")

	f, err := os.Create("cpu_profile.prof")
	if err != nil {
		panic(err)
	}

	if err := pprof.StartCPUProfile(f); err != nil {
		f.Close()
		panic(err)
	}

	return func() {
		pprof.StopCPUProfile()
		f.Close()
	}
}

func Perf(message string) func() {
	t1 := time.Now()

	return func() {
		fmt.Println(message, time.Since(t1))
	}
}
