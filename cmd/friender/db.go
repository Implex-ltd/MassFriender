package main

import (
	"log"
	"os"
	"strings"
	"sync"

	"path/filepath"

	"github.com/0xF7A4C6/GoCycle"
)

var (
	Inputs = make(map[string]*GoCycle.Cycle)
)

func LoadDataset() error {
	inputDir := "../../assets/input/"

	files, err := os.ReadDir(inputDir)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	assetsMutex := sync.Mutex{}

	for _, file := range files {
		wg.Add(1)
		go func(file os.DirEntry) {
			defer wg.Done()

			if strings.HasSuffix(file.Name(), ".txt") {
				filePath := filepath.Join(inputDir, file.Name())

				c, err := GoCycle.NewFromFile(filePath)
				if err != nil {
					panic(err)
				}

				//c.ClearDuplicates()
				c.RandomiseIndex()

				assetName := strings.Split(file.Name(), ".txt")[0]
				log.Printf("Loaded %v (%s)", len(c.List), assetName)

				assetsMutex.Lock()
				Inputs[assetName] = c
				assetsMutex.Unlock()
			}
		}(file)
	}

	wg.Wait()

	// do not iterate trought tokens
	Inputs["tokens"].WaitForUnlock = false
	Inputs["username"].WaitForUnlock = false

	Unprocessed += len(Inputs["username"].List)

	return nil
}
