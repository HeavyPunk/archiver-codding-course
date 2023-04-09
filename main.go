package main

import (
	"fmt"
	algorithm_running "kirieshki/running-archiver/app/algorithms/running"
	"os"
)

var compress_algorithms = make(map[string]func(*os.File, string) error)
var decompress_algorithms = make(map[string]func(*os.File, string) error)
var actions = make(map[string]func(string, string))

func main() {
	compress_algorithms["running"] = algorithm_running.Compress
	decompress_algorithms["running"] = algorithm_running.Decompress
	actions["compress"] = compress
	actions["decompress"] = decompress

	if len(os.Args) < 3 {
		fmt.Printf("You should specify at least filename for run archiver")
		return
	}

	fileName := os.Args[2]
	actionName := os.Args[1]
	algorithm := "running"
	if len(os.Args) > 3 {
		algorithm = os.Args[3]
	}
	action, ok := actions[actionName]
	if !ok {
		fmt.Printf("Action not found: %s", actionName)
		return
	}

	action(fileName, algorithm)
}

func compress(fileName string, algorithm string) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Problem when opening %s: %v", fileName, err)
	}

	algo, ok := compress_algorithms[algorithm]

	if !ok {
		fmt.Printf("Algorithm %s not found", algorithm)
		return
	}

	targetFileName := fileName + "." + algorithm
	err = algo(file, targetFileName)

	if err != nil {
		fmt.Printf("Problem when archiving %s: %v", fileName, err)
	}

	fmt.Printf("Created archive %s", targetFileName)
}

func decompress(fileName string, algorithm string) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Problem when opening %s: %v", fileName, err)
	}

	algo, ok := decompress_algorithms[algorithm]

	if !ok {
		fmt.Printf("Algorithm %s not found", algorithm)
		return
	}

	targetFileName := fileName + ".recovered"
	err = algo(file, targetFileName)
	if err != nil {
		fmt.Printf("Problem when decompressing %s: %v", targetFileName, err)
	}

	fmt.Printf("Recovered file %s", fileName)
}