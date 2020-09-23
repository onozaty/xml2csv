package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/antchfx/xmlquery"
)

// Column カラムの定義
type Column struct {
	Header    string `json:"header"`
	ValuePath string `json:"valuePath"`
}

// Mapping マッピング情報
type Mapping struct {
	RowsPath string   `json:"rowsPath"`
	Columns  []Column `json:"columns"`
}

func main() {

	var xmlPath string
	var mappingPath string
	var csvPath string

	flag.StringVar(&xmlPath, "i", "", "XML input file path or directory")
	flag.StringVar(&mappingPath, "m", "", "XML to CSV mapping file path")
	flag.StringVar(&csvPath, "o", "", "CSV output file path")
	flag.Parse()

	if xmlPath == "" || mappingPath == "" || csvPath == "" {
		flag.Usage()
	}

	mapping := loadMapping(mappingPath)

	csvFile, err := os.Create(csvPath)
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	convert([]string{xmlPath}, mapping, csvFile)
}

func convert(xmlPaths []string, mapping Mapping, writer io.Writer) {

	csvWriter := csv.NewWriter(writer)

	// header
	var headers []string
	for _, column := range mapping.Columns {
		headers = append(headers, column.Header)
	}

	err := csvWriter.Write(headers)
	if err != nil {
		log.Fatal(err)
	}

	// rows
	for _, xmlPath := range xmlPaths {
		doc := loadXML(xmlPath)
		convertOne(doc, mapping, csvWriter)
	}

	csvWriter.Flush()
}

func convertOne(doc *xmlquery.Node, mapping Mapping, csvWriter *csv.Writer) {

	rows, err := xmlquery.QueryAll(doc, mapping.RowsPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, row := range rows {

		var values []string
		for _, column := range mapping.Columns {

			value, err := xmlquery.Query(row, column.ValuePath)
			if err != nil {
				log.Fatal(err)
			}

			if value != nil {
				values = append(values, value.InnerText())
			} else {
				values = append(values, "")
			}
		}

		err := csvWriter.Write(values)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func loadMapping(path string) Mapping {

	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	var mapping Mapping
	err = json.Unmarshal(content, &mapping)
	if err != nil {
		log.Fatal(err)
	}

	return mapping
}

func loadXML(path string) *xmlquery.Node {

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	doc, err := xmlquery.Parse(file)
	if err != nil {
		log.Fatal(err)
	}

	return doc
}
