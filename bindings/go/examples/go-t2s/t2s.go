package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

var m_t2s map[string]string

func IsPathExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func ReadFileToSlice(filePath string, handleFunc func(string) (string, bool)) ([]string, error) {
	if !IsPathExist(filePath) {
		return nil, errors.New(filePath + " not found")
	}
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var fileList []string
	var s string
	fileScanner := bufio.NewScanner(f)
	for fileScanner.Scan() {
		s = fileScanner.Text()
		if s, ok := handleFunc(s); ok {
			fileList = append(fileList, s)
		}
	}
	return fileList, nil
}

func ReadFileToSliceEx(filePath string) ([]string, error) {
	return ReadFileToSlice(filePath, func(s string) (string, bool) {
		return s, true
	})
}

func ReadFile(filePath string) ([]byte, error) {
	bs, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func LoadDict(filePath string) map[string]string {
	m := make(map[string]string)

	ss, _ := ReadFileToSliceEx(filePath)
	for _, s := range ss {
		sps := strings.Split(s, "\t")
		if len(sps) < 2 {
			fmt.Println("split error:", sps)
			continue
		}
		m[sps[0]] = sps[1]
	}

	fmt.Println("LoadDict length:", len(m))
	return m
}

func T2S(ss string) string {
	if len(m_t2s) == 0 {
		return ss
	}
	result := ""
	for _, s := range ss {
		v, ok := m_t2s[string(s)]
		if !ok {
			v = string(s)
		}
		result += v
	}
	return result
}
