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

var (
	version = "dev"
	commit  = "none"
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

	if len(commit) > 7 {
		commit = commit[:7]
	}
	fmt.Printf("xml2csv v%s (%s)\n", version, commit)

	var xmlPath string
	var mappingPath string
	var csvPath string
	var withBom bool
	var help bool

	flag.StringVar(&xmlPath, "i", "", "XML input file path or directory")
	flag.StringVar(&mappingPath, "m", "", "XML to CSV mapping file path")
	flag.StringVar(&csvPath, "o", "", "CSV output file path")
	flag.BoolVar(&withBom, "b", false, "CSV with BOM")
	flag.BoolVar(&help, "h", false, "Help")
	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(0)
	}

	if xmlPath == "" || mappingPath == "" || csvPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	mapping, err := loadMapping(mappingPath)
	if err != nil {
		log.Fatal(err)
	}

	csvFile, err := os.Create(csvPath)
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	xmlPaths, err := findXML(xmlPath)
	if err != nil {
		log.Fatal(err)
	}

	if withBom {
		// BOMを付与
		csvFile.Write([]byte{0xEF, 0xBB, 0xBF})
	}

	err = convert(xmlPaths, mapping, csvFile)
	if err != nil {
		log.Fatal(err)
	}
}

func convert(xmlPaths []string, mapping *Mapping, writer io.Writer) error {

	csvWriter := csv.NewWriter(writer)

	// header
	var headers []string
	for _, column := range mapping.Columns {
		headers = append(headers, column.Header)
	}

	err := csvWriter.Write(headers)
	if err != nil {
		return err
	}

	// rows
	for _, xmlPath := range xmlPaths {
		doc, err := parseXML(xmlPath)
		if err != nil {
			return err
		}
		err = convertOne(doc, mapping, csvWriter)
		if err != nil {
			return err
		}
	}

	csvWriter.Flush()

	return nil
}

func convertOne(doc *xmlquery.Node, mapping *Mapping, csvWriter *csv.Writer) error {

	rows, err := xmlquery.QueryAll(doc, mapping.RowsPath)
	if err != nil {
		return err
	}

	for _, row := range rows {

		var values []string
		for _, column := range mapping.Columns {
			value, err := getValue(row, column.ValuePath, column.UseEvaluate)
			if err != nil {
				return err
			}

			values = append(values, value)
		}

		err := csvWriter.Write(values)
		if err != nil {
			return err
		}
	}

	return nil
}

func getValue(row *xmlquery.Node, valuePath string, useEvaluate bool) (string, error) {

	// Node以外を返すような式の場合(count()、boolean()など)
	if useEvaluate {
		expr, err := xpath.Compile(valuePath)
		if err != nil {
			return "", err
		}

		value := expr.Evaluate(xmlquery.CreateXPathNavigator(row))
		return fmt.Sprint(value), nil
	}

	// Nodeを返す場合
	value, err := xmlquery.Query(row, valuePath)
	if err != nil {
		return "", err
	}

	if value == nil {
		return "", nil
	}

	return value.InnerText(), nil
}

func loadMapping(path string) (*Mapping, error) {

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var mapping Mapping
	err = json.Unmarshal(content, &mapping)
	if err != nil {
		return nil, err
	}

	return &mapping, nil
}

func parseXML(path string) (*xmlquery.Node, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	doc, err := xmlquery.Parse(file)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func findXML(path string) ([]string, error) {

	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if !fileInfo.IsDir() {
		// ファイル
		return []string{path}, nil
	}

	// ディレクトリの場合、配下のファイルを取得
	fileInfosInDir, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, fileInfoInDir := range fileInfosInDir {
		if !fileInfoInDir.IsDir() {
			files = append(files, filepath.Join(path, fileInfoInDir.Name()))
		}
	}

	sort.Strings(files)
	return files, nil
}
