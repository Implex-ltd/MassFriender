package main

import (
	"log"
	"os"
	"strings"
	"sync"

	"path/filepath"

	"github.com/0xF7A4C6/GoCycle"
	"github.com/Implex-ltd/friender/internal/console"
)

var (
	Inputs = make(map[string]*GoCycle.Cycle)
)

func LoadDir(inputDir string) error {
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

	return nil
}

func LoadDataset() error {
	for _, path := range []string{
		"../../assets/input/",
		"../../assets/data/",
	} {
		if err := LoadDir(path); err != nil {
			return err
		}
	}
	
	Inputs["tokens"].ClearDuplicates()
	Inputs["tokens"].WaitForUnlock = false
	Inputs["username"].WaitForUnlock = false

	console.Unprocessed += len(Inputs["username"].List)

	return nil
}
