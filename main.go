package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Record struct {
	From  string `json:"from"`
	Moved string `json:"moved"`
}

const (
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
)

func main() {
	file, err := os.Open("./json.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	rawCSVdata, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var records []Record

	for _, each := range rawCSVdata {
		record := Record{
			From:  each[0],
			Moved: each[1],
		}
		records = append(records, record)
	}

	outputFile, err := os.Create("a.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer outputFile.Close()

	encoder := json.NewEncoder(outputFile)
	if err := encoder.Encode(records); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	TestRedirects()
}

func TestRedirects() {

	filepaths := []string{"/Users/gabriel.franca/repos/docs/docs/cicd/azion-massive-redirect-ptbr.json", "/Users/gabriel.franca/repos/docs/docs/cicd/azion-massive-redirect-en.json"}

	for _, fpath := range filepaths {
		jsonFile, err := os.Open(fpath)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer jsonFile.Close()

		byteValue, _ := io.ReadAll(jsonFile)

		var redirects []Record
		err = json.Unmarshal(byteValue, &redirects)
		if err != nil {
			fmt.Println(err)
			return
		}

		for i := 0; i < len(redirects); i++ {
			resp, err := http.Get(redirects[i].Moved)
			if err != nil {
				fmt.Println(err)
				continue
			}

			if resp.StatusCode == 404 {
				fmt.Println(colorRed, "Link:", redirects[i].Moved)
			} //else {
			//fmt.Println(colorGreen, "Link:", redirects[i].Moved)

			//}
			resp.Body.Close()
		}

	}

}
