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

// type DataMap map[string]*Data
type DataMap map[string]Data

type Data struct {
	min, max, sum float64
	count         int64
}

func Execute(file *os.File, rows int) (string, int) {
	workers := runtime.GOMAXPROCS(runtime.NumCPU())

	parts, _ := splitFile(file, workers)

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
	_, err := file.Seek(part.offset, io.SeekStart)
	if err != nil {
		panic(err)
	}

	f := io.LimitReader(file, part.size)
	data := make(DataMap)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		name, temperature, _ := strings.Cut(line, ";")
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

type Part struct {
	offset, size int64
}

func splitFile(f *os.File, numParts int) ([]Part, error) {
	const maxLineLength = 100

	st, _ := f.Stat()
	size := st.Size()
	splitSize := size / int64(numParts)

	buf := make([]byte, maxLineLength)
	parts := make([]Part, numParts)
	offset := int64(0)

	for i := 0; i < numParts; i++ {
		if i == numParts-1 {
			if offset < size {
				parts = append(parts, Part{offset, size - offset})
			}
			break
		}

		seekOffset := max(offset+splitSize-maxLineLength, 0)
		_, err := f.Seek(seekOffset, io.SeekStart)
		if err != nil {
			return nil, err
		}
		n, _ := io.ReadFull(f, buf)
		chunk := buf[:n]
		newline := bytes.LastIndexByte(chunk, '\n')
		if newline < 0 {
			return nil, fmt.Errorf("newline not found at offset %d", offset+splitSize-maxLineLength)
		}
		remaining := len(chunk) - newline - 1
		nextOffset := seekOffset + int64(len(chunk)) - int64(remaining)
		parts = append(parts, Part{offset, nextOffset - offset})
		offset = nextOffset
	}
	return parts, nil
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
