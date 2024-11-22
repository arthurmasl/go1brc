package solution3

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"slices"
	"strings"
	"sync"
)

type DataMap map[string]*Data

type Data struct {
	min, max, count int32
	sum             int64
}

type Chunks chan []string

const BUFFER_SIZE = 2048 * 2048

func Execute(file *os.File, rows int) (string, int) {
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	workers := runtime.GOMAXPROCS(runtime.NumCPU())
	chunkSize := rows / workers

	data := make(DataMap)
	chunks := make(Chunks, workers)

	fmt.Printf("Workers: %v, Chunk Size: %v\n", workers, chunkSize)

	runWorkers(chunks, workers, &data, &wg, &mu)
	scanLines(chunks, file, chunkSize)

	wg.Wait()

	result := parseResult(&data)
	return result, len(data)
}

func runWorkers(
	chunks Chunks,
	workers int,
	data *DataMap,
	wg *sync.WaitGroup,
	mu *sync.Mutex,
) {
	for range workers {
		wg.Add(1)

		go func() {
			defer wg.Done()

			lines := <-chunks
			processChunk(lines, data, mu)
		}()
	}
}

func scanLines(
	chunks Chunks,
	file *os.File,
	chunkSize int,
) {
	lines := make([]string, 0, chunkSize)

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, BUFFER_SIZE), BUFFER_SIZE)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())

		if len(lines) == chunkSize {
			chunks <- lines
			lines = nil
		}
	}

	if len(lines) > 0 {
		chunks <- lines
	}

	close(chunks)
}

func processChunk(
	lines []string,
	data *DataMap,
	mu *sync.Mutex,
) {
	// mu.Lock()
	// defer mu.Unlock()

	for _, line := range lines {
		name, tempStr, _ := strings.Cut(line, ";")

		negative := false
		index := 0
		if tempStr[0] == '-' {
			index++
			negative = true
		}
		temp := int32(tempStr[index] - '0')
		index++
		if tempStr[index] != '.' {
			temp = temp*10 + int32(tempStr[index]-'0')
		}
		index++
		temp = temp*10 + int32(tempStr[index]-'0')
		if negative {
			temp = -temp
		}

		d := (*data)[name]
		if d == nil {
			(*data)[name] = &Data{
				min:   temp,
				max:   temp,
				sum:   int64(temp),
				count: 1,
			}
		} else {
			d.min = min(d.min, temp)
			d.max = max(d.min, temp)
			d.sum += int64(temp)
			d.count++
		}
	}
}

func parseResult(data *DataMap) string {
	index := 0
	dataArr := make([]string, len(*data))

	for key, value := range *data {
		mean := float64(value.sum) / float64(value.count) / 10
		dataArr[index] = fmt.Sprintf(
			"%v=%.1f/%.1f/%.1f",
			key,
			float64(value.min)/10,
			mean,
			float64(value.max)/10,
		)
		index++
	}

	slices.Sort(dataArr)
	result := fmt.Sprintf("{%v}", strings.Join(dataArr, ", "))

	return result
}
