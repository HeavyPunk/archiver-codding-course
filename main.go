package main

import (
	"fmt"
	algorithm_basic "kirieshki/running-archiver/app/algorithms/basic"
	algorithm_ngram "kirieshki/running-archiver/app/algorithms/ngram"
	algorithm_running "kirieshki/running-archiver/app/algorithms/running"
	utils_collections "kirieshki/running-archiver/utils/collections"
	"os"
)

var compress_algorithms = make(map[string]func(*os.File, string, bool) error)
var decompress_algorithms = make(map[string]func(*os.File, string, bool) error)
var actions = make(map[string]func(string, string, bool))

func main() {
	compress_algorithms["running"] = algorithm_running.Compress
	decompress_algorithms["running"] = algorithm_running.Decompress
	compress_algorithms["basic"] = algorithm_basic.Compress
	decompress_algorithms["basic"] = algorithm_basic.Decompress
	compress_algorithms["ngram"] = algorithm_ngram.Compress
	decompress_algorithms["ngram"] = algorithm_ngram.Decompress
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
	debugFlag, _ := utils_collections.Find(os.Args, func(arg string) bool { return arg == "--debug" })
	if !ok {
		fmt.Printf("Action not found: %s", actionName)
		return
	}

	action(fileName, algorithm, debugFlag != "")
}

func compress(fileName string, algorithm string, debug bool) {
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
	cleanUp := !debug
	err = algo(file, targetFileName, cleanUp)

	if err != nil {
		fmt.Printf("Problem when archiving %s: %v", fileName, err)
	}

	fmt.Printf("Created archive %s", targetFileName)
}

func decompress(fileName string, algorithm string, debug bool) {
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
	cleanUp := !debug
	err = algo(file, targetFileName, cleanUp)
	if err != nil {
		fmt.Printf("Problem when decompressing %s: %v", targetFileName, err)
	}

	fmt.Printf("Recovered file %s", fileName)
}
