package tools

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func ReadFileTxt(filepath string) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("Failed to read file %s: %v", filepath, err)
	}

	defer file.Close()

	var value []string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		value = append(value, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Error reading file %s: %v", filepath, err)
	}

	return value, nil
}

func ReadFileJson(filePath string) (interface{}, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Coba unmarshal sebagai array of generic maps (map[string]interface{})
	var dataArray []map[string]interface{}
	if err := json.Unmarshal(byteValue, &dataArray); err == nil {
		return dataArray, nil
	}

	// Jika gagal, coba unmarshal sebagai generic map
	var dataObject map[string]interface{}
	if err := json.Unmarshal(byteValue, &dataObject); err == nil {
		return dataObject, nil
	}

	return nil, fmt.Errorf("failed to unmarshal JSON from file %s", filePath)
}
