package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func AppendLineInDirectory(directory, fileName, line string) error {
	fullPath := filepath.Join(directory, fileName)

	file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintln(file, line)
	if err != nil {
		return err
	}

	return nil
}

func CalculateAverage(arr []int) float64 {
	if len(arr) == 0 {
		return 0.0
	}

	sum := 0
	for _, num := range arr {
		sum += num
	}

	return float64(sum) / float64(len(arr))
}
