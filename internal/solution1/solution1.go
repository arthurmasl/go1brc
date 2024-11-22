package solution1

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Data struct {
	min, max, sum float64
	count         int64
}

func Execute(file *os.File, rows int) (string, int) {
	data := make(map[string]*Data)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ";")
		name := parts[0]
		temp, _ := strconv.ParseFloat(parts[1], 64)

		d, ok := data[name]
		if !ok {
			data[name] = &Data{
				min:   temp,
				max:   temp,
				sum:   temp,
				count: 1,
			}
		} else {
			if temp < d.min {
				d.min = temp
			}
			if temp > d.max {
				d.max = temp
			}

			d.sum += temp
			d.count++
		}

	}

	index := 0
	strs := make([]string, len(data))

	for key, value := range data {
		mean := value.sum / float64(value.count)
		strs[index] = fmt.Sprintf("%v=%v/%v/%v", key, value.min, mean, value.max)
		index++
	}

	slices.Sort(strs)
	result := fmt.Sprintf("{%v}", strings.Join(strs, ", "))

	return result, len(strs)
}
