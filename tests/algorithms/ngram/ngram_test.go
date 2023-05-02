package test_ngram

import (
	algorithm_ngram "kirieshki/running-archiver/app/algorithms/ngram"
	utils_collections "kirieshki/running-archiver/utils/collections"
	"os"
	"testing"
)

func TestCreateNgrams(t *testing.T) {
	file, err := os.Open("construct-ngrams-input")
	if err != nil {
		t.Errorf("Cannot open construct-ngrams-input: %v", err)
		return
	}
	res := algorithm_ngram.ConstructNgramDict(file)
	if len(res) != 6 {
		t.Errorf("Found %d ngrams but should be %d", len(res), 6)
	}
}

func TestCreateShannonCode(t *testing.T) {
	p, code := []int{1, 2, 3, 4, 5, 6}, make(map[int]string)
	algorithm_ngram.CreateShannonCode(
		utils_collections.Map(p, func(i int) algorithm_ngram.GroupUnit { return algorithm_ngram.GroupUnit{Id: i, Probability: i} }),
		code,
		0,
	)
	if len(code) != 6 {
		t.Errorf("Found %d code phrases but should be %d", len(code), 6)
	}
}

func TestSplitGroup(t *testing.T) {
	g := []int{2, 4, 5, 6, 7, 9}
	answer := 4
	pointer := algorithm_ngram.SplitGroup(
		utils_collections.Map(g, func(i int) algorithm_ngram.GroupUnit { return algorithm_ngram.GroupUnit{Id: i, Probability: i} }),
	)
	if pointer != answer {
		t.Errorf("Pointer must be %d but actual is %d", answer, pointer)
	}
}

func TestWriteCompressedFile(t *testing.T) {
	sourceFile, err := os.Open("construct-ngrams-input")
	if err != nil {
		t.Errorf("Cannot open construct-ngrams-input: %v", err)
		return
	}
	algorithm_ngram.WriteCompressedFile("write-compressed-test", sourceFile, map[string]string{
		"ab": "1000",
		"bc": "1001",
		"cd": "101",
		"de": "11",
		"ef": "00",
		"f":  "01",
	})
}

func TestDecompressFile(t *testing.T) {
	sourceFile, err := os.Open("decompress-file-test-input")
	if err != nil {
		t.Errorf("Cannot open decompress-file-test-input: %v", err)
		return
	}
	algorithm_ngram.Decompress(sourceFile, "decompress-file-test-result")
}

func TestFullTest(t *testing.T) {
	sourceFile, err := os.Open("full-test-input-file")
	if err != nil {
		t.Errorf("Cannot open full-test-input-file: %v", err)
		return
	}
	err = algorithm_ngram.Compress(sourceFile, "full-test-compressed-file")
	if err != nil {
		t.Errorf("Error when compressing: %v", err)
		return
	}

	sourceFile, err = os.Open("full-test-compressed-file")
	if err != nil {
		t.Errorf("Cannot open full-test-compressed-file: %v", err)
		return
	}
	err = algorithm_ngram.Decompress(sourceFile, "full-test-result-file")
	if err != nil {
		t.Errorf("Error when decompress: %v", err)
		return
	}

	res, err := os.ReadFile("full-test-result-file")
	if err != nil {
		t.Errorf("Cannot read full-test-result-file: %v", err)
		return
	}

	exp, err := os.ReadFile("full-test-input-file")
	if err != nil {
		t.Errorf("Cannot read full-test-input-file: %v", err)
		return
	}

	if string(exp) != string(res) {
		t.Errorf("Expected: '%s' actual: '%s'", exp, res)
	}
}
