package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
)

func readLines(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func compactJSON(jsonInput []string) string {
	jsonLine := strings.Join(jsonInput, "\n")
	var jsonBuffer *bytes.Buffer = new(bytes.Buffer)
	json.Compact(jsonBuffer, []byte(jsonLine))

	return jsonBuffer.String()
}

func checkDataType(jsonInput string) string {

	var jsonData interface{}

	// decode JSON-encoded collection of bytes into map[string]interface{}
	jsonBytes := []byte(jsonInput)
	if err := json.Unmarshal(jsonBytes, &jsonData); err != nil {
		panic(err)
	}

	return reflect.TypeOf(jsonData).String()
}

func getValue(strKey string, intKey int, jsonData interface{}) {
	switch jsonData.(type) {
	case []interface{}:
		for key, val := range jsonData.([]interface{}) {
			getValue(strKey, key, val)
		}
	case map[string]interface{}:
		for key, val := range jsonData.(map[string]interface{}) {
			getValue(key, intKey, val)
		}
	case string:
		if strKey == "" {
			fmt.Printf("Index:  %d  |  Value: %v \n", intKey, jsonData)
		} else {
			fmt.Printf("Key[%d]: %s  |  Value:  %v \n", intKey, strKey, jsonData)
		}
	default:
		fmt.Printf("Key[%d]: %s  |  Value:  %v \n", intKey, strKey, jsonData)
	}

}

func parseJSON(jsonInput string, argv []string) {

	jsonType := checkDataType(jsonInput)

	switch jsonType {
	case "map[string]interface {}":
		parseJSONObject(jsonInput, argv)
		break
	case "[]interface {}":
		parseJSONArray(jsonInput, argv)
		break
	default:
		fmt.Printf("Received a JSON string of type %s ?", jsonType)
		break
	}

}

func parseJSONArray(jsonInput string, argv []string) {

	var jsonData []interface{}
	jsonBytes := []byte(jsonInput)

	if err := json.Unmarshal(jsonBytes, &jsonData); err != nil {
		panic(err)
	}
	getValue("", 0, jsonData)
}

func parseJSONObject(jsonInput string, argv []string) {

	var jsonData map[string]interface{}
	jsonBytes := []byte(jsonInput)

	if err := json.Unmarshal(jsonBytes, &jsonData); err != nil {
		panic(err)
	}
	getValue("", 0, jsonData)
}

func main() {
	// "Usage: main.exe filepath"

	if len(os.Args) == 1 {
		fmt.Println("Usage: main.exe filepath")
		os.Exit(-1)
	}

	jsonLines := readLines(os.Args[1])
	jsonLine := compactJSON(jsonLines)

	parseJSON(jsonLine, os.Args[2:])
}
