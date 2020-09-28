package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/antchfx/xmlquery"
	"github.com/antchfx/xpath"
)

// Column カラムの定義
type Column struct {
	Header      string `json:"header"`
	ValuePath   string `json:"valuePath"`
	UseEvaluate bool   `json:"useEvaluate"`
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
	var withBom bool

	flag.StringVar(&xmlPath, "i", "", "XML input file path or directory")
	flag.StringVar(&mappingPath, "m", "", "XML to CSV mapping file path")
	flag.StringVar(&csvPath, "o", "", "CSV output file path")
	flag.BoolVar(&withBom, "bom", false, "CSV with BOM")
	flag.Parse()

	if xmlPath == "" || mappingPath == "" || csvPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	mapping := loadMapping(mappingPath)

	csvFile, err := os.Create(csvPath)
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	xmlPaths := findXML(xmlPath)

	if withBom {
		// BOMを付与
		csvFile.Write([]byte{0xEF, 0xBB, 0xBF})
	}

	convert(xmlPaths, mapping, csvFile)
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
		doc := parseXML(xmlPath)
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
			value := getValue(row, column.ValuePath, column.UseEvaluate)
			values = append(values, value)
		}

		err := csvWriter.Write(values)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getValue(row *xmlquery.Node, valuePath string, useEvaluate bool) string {

	// Node以外を返すような式の場合(count()、boolean()など)
	if useEvaluate {
		expr, err := xpath.Compile(valuePath)
		if err != nil {
			log.Fatal(err)
		}

		value := expr.Evaluate(xmlquery.CreateXPathNavigator(row))
		return fmt.Sprint(value)
	}

	// Nodeを返す場合
	value, err := xmlquery.Query(row, valuePath)
	if err != nil {
		log.Fatal(err)
	}

	if value == nil {
		return ""
	}

	return value.InnerText()
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

func parseXML(path string) *xmlquery.Node {

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

func findXML(path string) []string {

	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}

	if !fileInfo.IsDir() {
		// ファイル
		return []string{path}
	}

	// ディレクトリの場合、配下のファイルを取得
	fileInfosInDir, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	var files []string
	for _, fileInfoInDir := range fileInfosInDir {
		if !fileInfoInDir.IsDir() {
			files = append(files, filepath.Join(path, fileInfoInDir.Name()))
		}
	}

	sort.Strings(files)
	return files
}
