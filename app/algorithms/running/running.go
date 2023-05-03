package algorithm_running

import (
	"io"
	"os"
)

func compressArray(arr []byte) ([]byte, error) {
	if len(arr) == 0 {
		return nil, io.EOF
	}

	res := make([]byte, 0)
	curr := arr[0]
	count := 0
	for _, b := range arr {
		if count == 255 {
			res = append(res, curr, byte(count))
			curr = b
			count = 1
		}
		if b == curr {
			count++
		} else {
			res = append(res, curr, byte(count))
			curr = b
			count = 1
		}
	}

	res = append(res, curr, byte(count))
	return res, nil
}

func decompressArray(arr []byte) (res []byte, err error) {
	res = make([]byte, 0)
	for i := 0; i < len(arr); i += 2 {
		count := int(arr[i+1])
		buff := make([]byte, count)
		for j := 0; j < count; j++ {
			buff[j] = arr[i]
		}
		res = append(res, buff...)
	}
	return res, nil
}

func Compress(sourceFile *os.File, targetFileName string, cleanUp bool) error {
	allocBuff := make([]byte, 100)
	targetFile, err := os.Create(targetFileName)
	if err != nil {
		return err
	}
	for {
		read, err := sourceFile.Read(allocBuff)
		if err == io.EOF {
			break
		}
		buff := allocBuff[:read]
		toWrite, err := compressArray(buff)
		if err != nil {
			return err
		}
		_, err = targetFile.Write(toWrite)
		if err != nil {
			return err
		}
	}
	return nil
}

func Decompress(sourceFile *os.File, targetFileName string, cleanUp bool) error {
	allocBuff := make([]byte, 100)
	targetFile, err := os.Create(targetFileName)
	if err != nil {
		return err
	}
	for {
		readed, err := sourceFile.Read(allocBuff)
		if err == io.EOF {
			break
		}
		buff := allocBuff[:readed]
		toWrite, err := decompressArray(buff)
		if err != nil {
			return err
		}
		_, err = targetFile.Write(toWrite)
		if err != nil {
			return err
		}
	}
	return nil
}
