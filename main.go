package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type jsonTuple struct {
	intKey int
	strKey string
	data   interface{}
}

var flagInput = flag.String("i", "", "Specifies the input file")
var flagOutput = flag.String("o", "", "Specifies the output file")

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

func createFile(outputFile *string) {
	file, err := os.Create(*outputFile)
	if err != nil {
		panic(err)
	}

	defer file.Close()
}

func createWriter() *bufio.Writer {

	file, err := os.OpenFile(*flagOutput, os.O_WRONLY, 0666)
	if os.IsNotExist(err) == true {
		createFile(flagOutput)
		return createWriter()
	}

	dataWriter := bufio.NewWriter(file)
	return dataWriter
}

func writeValue(buffer interface{}, inputString string) {
	// convert to a bufio.Writer type so we can Flush() the buffer
	writeBuffer := buffer.(*bufio.Writer)

	fmt.Fprintf(writeBuffer, inputString)
	fmt.Fprintf(writeBuffer, "\n")

	writeBuffer.Flush()

}

func findValue(strKey string, intKey int, jsonData interface{}, dataWriter interface{}) {

	// use switch cases to destructure the JSON data structure depending on its current type
	switch jsonData.(type) {
	case []interface{}:
		// iterate through values in interface array, calling getValue() on each, preserving the index
		for key, val := range jsonData.([]interface{}) {
			findValue(strKey, key, val, dataWriter)
		}
	case map[string]interface{}:
		// when the data structure is a map, preserve the string primary key value instead of an index
		for key, val := range jsonData.(map[string]interface{}) {
			findValue(key, intKey, val, dataWriter)
		}
	case string:
		// in the case of an array without a primary key, where the pkey value would be "", format our output accordingly
		if strKey == "" {
			fmt.Printf("[%d]  ,  `%v` \n", intKey, jsonData)

			// if an output file was specified, write the value that we found
			if *flagOutput != "" {
				str := fmt.Sprintf("%d, %v", intKey, jsonData)
				writeValue(dataWriter, str)
			}
		} else {
			fmt.Printf("[%d]: %s  ,  `%v` \n", intKey, strKey, jsonData)

			if *flagOutput != "" {
				str := fmt.Sprintf("%d, %s, %v", intKey, strKey, jsonData)
				writeValue(dataWriter, str)
			}
		}
	default:
		// for other types, e.g. float64, bool, int, return this generic format
		fmt.Printf("[%d]: %s  ,  `%v` \n", intKey, strKey, jsonData)

		if *flagOutput != "" {
			str := fmt.Sprintf("%d, %s, %v", intKey, strKey, jsonData)
			writeValue(dataWriter, str)
		}
	}
}

func parseJSON(jsonInput string) {

	// determine the type of the unmarshaled json string
	jsonType := checkDataType(jsonInput)

	// depending on the type, send the data into its corresponding function
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

	// if an output file was specified, create a writer and pass it along to the value finder
	if *flagOutput != "" {
		dataWriter := createWriter()
		findValue("", 0, jsonData, dataWriter)
	} else {
		// without an output specified, pass along the generic os.Stdout writer
		// to satisfy the argument requirement
		findValue("", 0, jsonData, os.Stdout)
	}

}

func parseJSONObject(jsonInput string) {

	// for objects, unmarshal into a map[string]interface{} so we can retrieve primary key values
	var jsonData map[string]interface{}
	jsonBytes := []byte(jsonInput)

	if err := json.Unmarshal(jsonBytes, &jsonData); err != nil {
		panic(err)
	}

	if *flagOutput != "" {
		dataWriter := createWriter()
		findValue("", 0, jsonData, dataWriter)
	} else {
		findValue("", 0, jsonData, os.Stdout)
	}
}

func checkFileExists(filepath string) {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return
	} else if os.IsExist(err) == true {
		fmt.Printf("File %s already exists.", filepath)
		os.Exit(1)
	}
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

func main() {
	// "Usage: main.exe filepath"
	// Optionally, specify an output "main.exe -i filepath -o output.txt"

	if *flagOutput != "" {
		checkFileExists(*flagOutput)
	}

	flag.Parse()

	if len(os.Args) == 1 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// read the file at the path given at runtime and minify it, returning one line of JSON
	jsonLines := readLines(*flagInput)
	jsonLine := compactJSON(jsonLines)

	parseJSON(jsonLine)

}
