package algorithm_ngram

import (
	"io"
	utils_collections "kirieshki/running-archiver/utils/collections"
	utils_math "kirieshki/running-archiver/utils/math"
	"os"
	"sort"
)

const ngramLen = 2

func ConstructNgramDict(sourceFile *os.File) map[string]int {
	bufLen := 100
	res := make(map[string]int)
	allocBuff := make([]byte, bufLen)
	for {
		read, err := sourceFile.Read(allocBuff)
		if err == io.EOF {
			break
		}
		buff := allocBuff[:read]
		for i := range buff {
			var ngram string
			for o := 0; o < ngramLen; o++ {
				var char string
				if o+i < len(buff) {
					char = string(buff[o+i])
				} else {
					char = ""
				}
				ngram += char
			}
			res[ngram] += 1
		}
	}
	return res
}

type GroupUnit struct {
	Probability int
	Id          int
}

func SplitGroup(group []GroupUnit) int {
	if len(group) < 3 {
		return 1
	}
	pointer := 1
	sU := utils_math.Sum(utils_collections.Map(group[0:pointer], func(g GroupUnit) int { return g.Probability })...)
	sD := utils_math.Sum(utils_collections.Map(group[pointer:], func(g GroupUnit) int { return g.Probability })...)
	diff := utils_math.Abs(float64(sD - sU))
	for utils_math.Abs(float64(sU+group[pointer].Probability-sD+group[pointer+1].Probability)) < diff {
		sU += group[pointer].Probability
		sD -= group[pointer].Probability
		diff = utils_math.Abs(float64(sD - sU))
		pointer += 1
		if pointer == len(group)-1 {
			break
		}
	}
	return pointer
}

func CreateShannonCode(sourceProbabilities []GroupUnit, code map[int]string, key_offset int) {
	if len(sourceProbabilities) == 1 {
		return
	}
	pointer := SplitGroup(sourceProbabilities)
	for k := 0; k < pointer; k++ {
		var char string
		if _, ok := code[k+key_offset]; ok {
			char = "0"
		} else {
			char = "1"
		}
		code[k+key_offset] = code[k+key_offset] + char
	}

	for k := pointer; k < len(sourceProbabilities); k++ {
		var char string
		if _, ok := code[k+key_offset]; ok {
			char = "1"
		} else {
			char = "0"
		}
		code[k+key_offset] = code[k+key_offset] + char
	}
	CreateShannonCode(sourceProbabilities[0:pointer], code, key_offset+0)
	CreateShannonCode(sourceProbabilities[pointer:], code, key_offset+pointer)
}

func Compress(sourceFile *os.File, targetFileName string, cleanUp bool) error {
	type kvp struct {
		ngram string
		count int
		id    int
	}
	var ngrams = ConstructNgramDict(sourceFile)
	var probList = make([]kvp, 0)
	for k, v := range ngrams {
		probList = append(probList, kvp{ngram: k, count: v})
	}

	sort.Slice(probList, func(i, j int) bool { return probList[i].count > probList[j].count })

	ngramId := 0
	for i := range probList {
		probList[i].id = ngramId
		ngramId++
	}

	code := make(map[int]string)
	CreateShannonCode(
		utils_collections.Map(probList, func(s kvp) GroupUnit { return GroupUnit{Probability: s.count, Id: s.id} }),
		code,
		0,
	)
	ngramToCode := make(map[string]string)
	for i := range probList {
		ngramToCode[probList[i].ngram] = code[probList[i].id]
	}
	sourceFile, _ = os.Open(sourceFile.Name())
	WriteCompressedFile(targetFileName, sourceFile, ngramToCode, cleanUp)
	return nil
}

func WriteCompressedFile(targetFileName string, sourceFile *os.File, code map[string]string, cleanUp bool) {
	targetFile, _ := os.Create(targetFileName)
	tempFileName := ".temp." + targetFileName
	tempFile, _ := os.Create(tempFileName)
	targetFile.Write(encodeCode(code))
	allocBuff := make([]byte, 100)
	for {
		read, err := sourceFile.Read(allocBuff)
		if err == io.EOF {
			break
		}
		buff := allocBuff[:read]
		for i := 0; i < read; {
			var ngram string
			for o := 0; o < ngramLen; o++ {
				var char string
				if o+i < len(buff) {
					char = string(buff[o+i])
				} else {
					char = ""
				}
				ngram += char
			}
			i += ngramLen
			toWrite := code[ngram]
			tempFile.Write([]byte(toWrite))
		}
	}

	tempFile, _ = os.Open(tempFileName)

	var toWrite byte
	offset := 0
	for {
		read, err := tempFile.Read(allocBuff)
		if err == io.EOF {
			break
		}
		buff := allocBuff[:read]
		toWriteBuff := make([]byte, 0)
		for _, b := range buff {
			offset++
			if offset%9 == 0 {
				offset = 1
				toWriteBuff = append(toWriteBuff, toWrite)
				if len(toWriteBuff) == 100 {
					targetFile.Write(toWriteBuff)
					toWriteBuff = make([]byte, 0)
				}
				toWrite = 0
			}
			toWrite <<= 1
			if b == '1' {
				toWrite |= 1
			} else {
				toWrite |= 0
			}
		}
		targetFile.Write(toWriteBuff)
	}
	targetFile.Write([]byte{toWrite << (8 - byte(offset)), byte(offset)})

	tempFile.Close()
	sourceFile.Close()
	targetFile.Close()
	removeFileIfNeed(tempFileName, cleanUp)
}

func encodeCode(code map[string]string) []byte {
	res := make([]byte, 0)
	for k, v := range code {
		res = append(res, []byte(k+"="+v)...)
		res = append(res, byte(';'))
	}
	res = append(res, ';')
	return res
}

func restoreCode(sourceFile *os.File) map[string]string {
	allocBuff := make([]byte, 1)
	res := make(map[string]string)

	var ngram string
	var value string
	isReadingNgram := true
	for {
		read, err := sourceFile.Read(allocBuff)
		if err == io.EOF {
			break
		}
		buff := allocBuff[:read]
		for _, b := range buff {
			bStr := string(b)
			if bStr == "=" {
				isReadingNgram = false
				continue
			}

			if bStr == ";" && ngram == "" {
				return res
			}

			if bStr == ";" {
				isReadingNgram = true
				res[value] = ngram
				value = ""
				ngram = ""
				continue
			}

			if isReadingNgram {
				ngram += bStr
			} else {
				value += bStr
			}
		}
	}
	return res
}

func Decompress(sourceFile *os.File, targetFileName string, cleanUp bool) error {
	targetFile, err := os.Create(targetFileName)
	if err != nil {
		return err
	}
	tempFileName := ".temp." + targetFileName
	tempFile, err := os.Create(tempFileName)
	if err != nil {
		removeFileIfNeed(tempFileName, cleanUp)
		return err
	}
	table := restoreCode(sourceFile)
	allocBuff := make([]byte, 100)
	var buff []byte
	var totalWrote int64
	for {
		read, err := sourceFile.Read(allocBuff)
		if err == io.EOF {
			break
		}
		buff = allocBuff[:read]
		for _, b := range buff {
			counter := 0
			var toWrite string
			for {
				if counter > 7 {
					break
				}

				firstBit := (b & (1 << 7)) >> 7
				if firstBit == 1 {
					toWrite += "1"
				} else {
					toWrite += "0"
				}

				b <<= 1
				counter++
			}
			tempFile.WriteString(toWrite)
			totalWrote += int64(len(toWrite))
		}
	}
	offset := int(buff[len(buff)-1])

	tempFile, _ = os.Open(tempFileName)

	allocBuff = make([]byte, 1)
	var totalRead int64
	var key string
	for {
		read, err := tempFile.Read(allocBuff)
		if err == io.EOF {
			break
		}
		if totalRead+int64(read) > totalWrote-int64(16-offset) {
			read -= int(totalRead + int64(read) - totalWrote + int64(16-offset))
		}
		totalRead += int64(read)
		buff := allocBuff[:read]

		for _, b := range buff {
			if b == '1' {
				key += "1"
			} else {
				key += "0"
			}
		}

		if ngram, ok := table[key]; ok {
			targetFile.WriteString(ngram)
			key = ""
		}
	}
	removeFileIfNeed(tempFileName, cleanUp)
	return nil
}

func removeFileIfNeed(filename string, need bool) error {
	if need {
		return os.Remove(filename)
	}
	return nil
}
