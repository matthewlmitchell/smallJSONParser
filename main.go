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

	// minify the provided JSON file into a single string
	var jsonBuffer *bytes.Buffer = new(bytes.Buffer)
	json.Compact(jsonBuffer, []byte(jsonLine))

	return jsonBuffer.String()
}

func checkDataType(jsonInput string) string {

	var jsonData interface{}

	// decode the JSON-encoded byte string into a generic interface{}
	jsonBytes := []byte(jsonInput)
	if err := json.Unmarshal(jsonBytes, &jsonData); err != nil {
		panic(err)
	}

	// return the data type of the decoded JSON string
	return reflect.TypeOf(jsonData).String()
}

func getValue(strKey string, intKey int, jsonData interface{}) {

	// when continually destructure the JSON data structure depending on its current type
	switch jsonData.(type) {
	case []interface{}:
		// iterate through values in interface array, calling getValue() on each, preserving the index
		for key, val := range jsonData.([]interface{}) {
			getValue(strKey, key, val)
		}
	case map[string]interface{}:
		// when the data structure is a map, preserve the string primary key value instead of an index
		for key, val := range jsonData.(map[string]interface{}) {
			getValue(key, intKey, val)
		}
	case string:
		// in the case of an array without a primary key, where the pkey value would be "", format our output accordingly
		if strKey == "" {
			fmt.Printf("Index:  %d  |  Value: %v \n", intKey, jsonData)
		} else {
			fmt.Printf("Key[%d]: %s  |  Value:  %v \n", intKey, strKey, jsonData)
		}
	default:
		// for other types, e.g. float64, bool, int, return this generic format
		fmt.Printf("Key[%d]: %s  |  Value:  %v \n", intKey, strKey, jsonData)
	}

}

func parseJSON(jsonInput string) {

	// determine the type of the unmarshaled json string
	jsonType := checkDataType(jsonInput)

	switch jsonType {
	case "map[string]interface {}":
		parseJSONObject(jsonInput)
		break
	case "[]interface {}":
		parseJSONArray(jsonInput)
		break
	default:
		fmt.Printf("Received a JSON string of type %s ?", jsonType)
		break
	}

}

func parseJSONArray(jsonInput string) {

	// for a JSON array, demarshal the JSON byte string into an array of interfaces, []interface{}
	// then proceed with value retrieval
	var jsonData []interface{}
	jsonBytes := []byte(jsonInput)

	if err := json.Unmarshal(jsonBytes, &jsonData); err != nil {
		panic(err)
	}
	getValue("", 0, jsonData)
}

func parseJSONObject(jsonInput string) {

	// for objects, unmarshal into a map[string]interface{} so we can retrieve primary key values
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

	// read the file at the path given at runtime and minify it, returning one line of JSON
	jsonLines := readLines(os.Args[1])
	jsonLine := compactJSON(jsonLines)

	parseJSON(jsonLine)
}
