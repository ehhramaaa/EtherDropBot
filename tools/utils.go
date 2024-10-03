package tools

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
)

func ExitsRecover() {
	if r := recover(); r != nil {
		Logger("error", fmt.Sprintf("%v", r))
		Logger("info", "Press Enter to exit...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}

func RecoverPanic() {
	if r := recover(); r != nil {
		Logger("error", fmt.Sprintf("%v", r))
	}
}

func RandomNumber(min, max int) int {
	return rand.Intn(max-min) + min
}

func InputTerminal(prompt string) string {
	Logger("input", prompt)

	reader := bufio.NewReader(os.Stdin)

	value, _ := reader.ReadString('\n')

	return strings.TrimSpace(value)
}

func SaveFileToJson(filePath string, data interface{}) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)

	// Encode data ke file
	err = encoder.Encode(data)
	if err != nil {
		return err
	}

	return nil
}

func AppendFileToJson(filePath string, newData interface{}) error {
	var existingData []interface{}

	if _, err := os.Stat(filePath); err == nil {
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&existingData); err != nil {
			if err != io.EOF {
				return fmt.Errorf("error decoding existing data: %v", err)
			}
		}
	}

	existingData = append(existingData, newData)

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(existingData); err != nil {
		return fmt.Errorf("error encoding new data: %v", err)
	}

	return nil
}
