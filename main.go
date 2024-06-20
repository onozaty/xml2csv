package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/onozaty/go-customcsv"

	"github.com/antchfx/xmlquery"
	"github.com/antchfx/xpath"

	flag "github.com/spf13/pflag"
)

var (
	Version = "dev"
	Commit  = "none"
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

const (
	OK int = 0
	NG int = 1
)

func main() {
	exitCode := run(os.Args[1:], os.Stdout)
	os.Exit(exitCode)
}

func run(arguments []string, output io.Writer) int {

	var xmlPath string
	var mappingPath string
	var csvPath string
	var withBom bool
	var help bool

	flagSet := flag.NewFlagSet("xml2csv", flag.ContinueOnError)

	flagSet.StringVarP(&xmlPath, "input", "i", "", "XML input file path or directory or url")
	flagSet.StringVarP(&mappingPath, "mapping", "m", "", "XML to CSV mapping file path or url")
	flagSet.StringVarP(&csvPath, "output", "o", "", "CSV output file path")
	flagSet.BoolVarP(&withBom, "bom", "b", false, "CSV with BOM")
	flagSet.BoolVarP(&help, "help", "h", false, "Help")

	flagSet.SortFlags = false
	flagSet.Usage = func() {
		fmt.Fprintf(output, "xml2csv v%s (%s)\n\n", Version, Commit)
		fmt.Fprint(output, "Usage: xml2csv [flags]\n\nFlags\n")
		flagSet.PrintDefaults()
		fmt.Fprintln(output)
	}
	flagSet.SetOutput(output)

	if err := flagSet.Parse(arguments); err != nil {
		flagSet.Usage()
		fmt.Fprintln(output, err)
		return NG
	}

	if help {
		flagSet.Usage()
		return OK
	}

	if xmlPath == "" || mappingPath == "" || csvPath == "" {
		flagSet.Usage()
		return NG
	}

	mapping, err := loadMapping(mappingPath)
	if err != nil {
		fmt.Fprintln(output, err)
		return NG
	}

	csvFile, err := os.Create(csvPath)
	if err != nil {
		fmt.Fprintln(output, err)
		return NG
	}
	defer csvFile.Close()

	xmlPaths, err := findXML(xmlPath)
	if err != nil {
		fmt.Fprintln(output, err)
		return NG
	}

	if withBom {
		// BOMを付与
		if _, err := csvFile.Write([]byte{0xEF, 0xBB, 0xBF}); err != nil {
			fmt.Fprintln(output, err)
			return NG
		}
	}

	if err := convert(xmlPaths, mapping, csvFile); err != nil {
		fmt.Fprintln(output, err)
		return NG
	}

	return OK
}

func convert(xmlPaths []string, mapping *Mapping, writer io.Writer) error {

	csvWriter := customcsv.NewWriter(writer)

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
		err = convertOne(xmlPath, mapping, csvWriter)
		if err != nil {
			return err
		}
	}

	csvWriter.Flush()

	return nil
}

func convertOne(xmlPath string, mapping *Mapping, csvWriter *customcsv.Writer) error {

	reader, err := open(xmlPath)
	if err != nil {
		return fmt.Errorf("%s is failed: %w", xmlPath, err)
	}
	defer reader.Close()

	parser, err := xmlquery.CreateStreamParser(reader, mapping.RowsPath)
	if err != nil {
		return fmt.Errorf("xpath '%s' is failed: %w", mapping.RowsPath, err)
	}

	for {
		row, err := parser.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("%s is failed: %w", xmlPath, err)
		}

		var values []string
		for _, column := range mapping.Columns {
			value, err := getValue(row, column.ValuePath, column.UseEvaluate)
			if err != nil {
				return err
			}

			values = append(values, value)
		}

		err = csvWriter.Write(values)
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
			return "", fmt.Errorf("xpath '%s' is failed: %w", valuePath, err)
		}

		value := expr.Evaluate(xmlquery.CreateXPathNavigator(row))
		return fmt.Sprint(value), nil
	}

	// Nodeを返す場合
	value, err := xmlquery.Query(row, valuePath)
	if err != nil {
		return "", fmt.Errorf("xpath '%s' is failed: %w", valuePath, err)
	}

	if value == nil {
		return "", nil
	}

	return value.InnerText(), nil
}

func loadMapping(path string) (*Mapping, error) {

	reader, err := open(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var mapping Mapping
	if err := json.Unmarshal(content, &mapping); err != nil {
		return nil, fmt.Errorf("invalid mapping format: %w", err)
	}

	return &mapping, nil
}

func findXML(path string) ([]string, error) {

	if isURL(path) {
		// URL
		return []string{path}, nil
	}

	// URL以外の場合には存在チェック
	if !exist(path) {
		return nil, fmt.Errorf("%s is not found", path)
	}

	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if !fileInfo.IsDir() {
		// ファイル
		return []string{path}, nil
	}

	// ディレクトリの場合、配下のファイルを取得
	fileInfosInDir, err := os.ReadDir(path)
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

func open(path string) (io.ReadCloser, error) {

	if isURL(path) {
		// URL
		resp, err := http.Get(path)
		if err != nil {
			return nil, err
		}

		return resp.Body, nil
	}

	// ファイル
	if !exist(path) {
		return nil, fmt.Errorf("%s is not found", path)
	}

	return os.Open(path)
}

func isURL(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}

func exist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
