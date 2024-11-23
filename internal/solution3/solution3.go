package solution3

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"slices"
	"strconv"
	"strings"
)

type Part struct {
	offset, size int64
}

type DataMap map[string]Data

type Data struct {
	min, max, sum float64
	count         int64
}

func Execute(file *os.File, rows int) (string, int) {
	workers := runtime.GOMAXPROCS(runtime.NumCPU())
	parts := splitFile(file, workers)

	resultCh := make(chan DataMap)
	for _, part := range parts {
		go processPart(file, &part, resultCh)
	}

	totals := make(DataMap)
	for range len(parts) {
		result := <-resultCh

		for station, s := range result {
			ts, ok := totals[station]
			if !ok {
				totals[station] = Data{
					min:   s.min,
					max:   s.max,
					sum:   s.sum,
					count: s.count,
				}
				continue
			}

			ts.min = min(ts.min, s.min)
			ts.max = max(ts.max, s.max)
			ts.sum += s.sum
			ts.count += s.count
			totals[station] = ts
		}
	}

	result := parseResult(&totals)
	return result, len(totals)
}

func processPart(file *os.File, part *Part, ch chan DataMap) {
	data := make(DataMap)

	file.Seek(part.offset, io.SeekStart)
	f := io.LimitReader(file, part.size)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		name, temperature, found := strings.Cut(line, ";")
		if !found {
			continue
		}
		temp, _ := strconv.ParseFloat(temperature, 64)

		d, ok := data[name]
		if !ok {
			d.min = temp
			d.max = temp
			d.sum = temp
			d.count++
		} else {
			d.min = min(d.min, temp)
			d.max = max(d.min, temp)
			d.sum += temp
			d.count++
		}

		data[name] = d
	}

	ch <- data
}

func splitFile(file *os.File, numParts int) []Part {
	const maxLineLength = 100

	parts := make([]Part, 0, numParts)
	buf := make([]byte, maxLineLength)

	stat, _ := file.Stat()
	fileSize := stat.Size()
	partSize := fileSize / int64(numParts)
	offset := int64(0)

	for offset < fileSize {
		seekOffset := max(offset+partSize-maxLineLength, 0)
		if seekOffset > fileSize {
			break
		}

		file.Seek(seekOffset, io.SeekStart)
		n, _ := io.ReadFull(file, buf)
		chunk := buf[:n]

		newline := bytes.LastIndexByte(chunk, '\n')
		remaining := len(chunk) - newline - 1
		nextOffset := seekOffset + int64(len(chunk)) - int64(remaining)

		parts = append(parts, Part{offset, nextOffset - offset})
		offset = nextOffset
	}

	return parts
}

func parseResult(data *DataMap) string {
	index := 0
	dataArr := make([]string, len(*data))

	for key, value := range *data {
		mean := value.sum / float64(value.count)
		dataArr[index] = fmt.Sprintf("%v=%.1f/%.1f/%.1f", key, value.min, mean, value.max)
		index++
	}

	slices.Sort(dataArr)
	result := fmt.Sprintf("{%v}", strings.Join(dataArr, ", "))

	return result
}
