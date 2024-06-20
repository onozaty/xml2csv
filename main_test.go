package main

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/onozaty/go-customcsv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun_File(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	inputPath := "testdata/rss.xml"

	mappingPath := createFile(t, temp, "mapping.json", `
	{
		"rowsPath": "//item",
		"columns": [
			{
				"header": "title",
				"valuePath": "/title"
			},
			{
				"header": "link",
				"valuePath": "/link"
			}
		]
	}`)

	outputPath := filepath.Join(temp, "output.csv")
	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			"-i", inputPath,
			"-m", mappingPath,
			"-o", outputPath,
		},
		out,
	)

	// ASSERT
	require.Equal(t, OK, exitCode)
	require.Empty(t, out.String())

	result := readString(t, outputPath)
	expect := joinRows(
		"title,link",
		"RSS Tutorial,https://www.w3schools.com/xml/xml_rss.asp",
		"XML Tutorial,https://www.w3schools.com/xml",
	)

	assert.Equal(t, expect, result)
}

func TestRun_URL(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	inputPath := "https://github.com/onozaty/xml2csv/raw/master/testdata/rss.xml"

	mappingPath := createFile(t, temp, "mapping.json", `
	{
		"rowsPath": "//item",
		"columns": [
			{
				"header": "title",
				"valuePath": "/title"
			},
			{
				"header": "link",
				"valuePath": "/link"
			}
		]
	}`)

	outputPath := filepath.Join(temp, "output.csv")
	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			"-i", inputPath,
			"-m", mappingPath,
			"-o", outputPath,
		},
		out,
	)

	// ASSERT
	require.Equal(t, OK, exitCode)
	require.Empty(t, out.String())

	result := readString(t, outputPath)
	expect := joinRows(
		"title,link",
		"RSS Tutorial,https://www.w3schools.com/xml/xml_rss.asp",
		"XML Tutorial,https://www.w3schools.com/xml",
	)

	assert.Equal(t, expect, result)
}

func TestRun_Dir(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	inputPath := "testdata/junit" // ディレクトリ指定

	mappingPath := createFile(t, temp, "mapping.json", `
	{
		"rowsPath": "//testcase",
		"columns": [
			{
				"header": "classname",
				"valuePath": "/@classname"
			},
			{
				"header": "name",
				"valuePath": "/@name"
			},
			{
				"header": "success",
				"valuePath": "not(/*)",
				"useEvaluate": true
			}
		]
	}`)

	outputPath := filepath.Join(temp, "output.csv")
	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			"-i", inputPath,
			"-m", mappingPath,
			"-o", outputPath,
		},
		out,
	)

	// ASSERT
	require.Equal(t, OK, exitCode)
	require.Empty(t, out.String())

	result := readString(t, outputPath)
	expect := joinRows(
		"classname,name,success",
		"com.github.onozaty.junit.xml2csv.TestCase1,test1,true",
		"com.github.onozaty.junit.xml2csv.TestCase1,test2,false",
		"com.github.onozaty.junit.xml2csv.TestCase1,test3,false",
		"com.github.onozaty.junit.xml2csv.TestCase1,test4,false",
		"com.github.onozaty.junit.xml2csv.TestCase1,test5,true",
		"com.github.onozaty.junit.xml2csv.TestCase2,test1,true",
		"com.github.onozaty.junit.xml2csv.TestCase2,test2,true",
	)

	assert.Equal(t, expect, result)
}

func TestRun_WithBom(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	inputPath := "testdata/rss.xml"

	mappingPath := createFile(t, temp, "mapping.json", `
	{
		"rowsPath": "//item",
		"columns": [
			{
				"header": "title",
				"valuePath": "/title"
			},
			{
				"header": "link",
				"valuePath": "/link"
			}
		]
	}`)

	outputPath := filepath.Join(temp, "output.csv")
	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			"-i", inputPath,
			"-m", mappingPath,
			"-o", outputPath,
			"-b",
		},
		out,
	)

	// ASSERT
	require.Equal(t, OK, exitCode)
	require.Empty(t, out.String())

	result := readString(t, outputPath)
	expect := joinRows(
		"\uFEFFtitle,link",
		"RSS Tutorial,https://www.w3schools.com/xml/xml_rss.asp",
		"XML Tutorial,https://www.w3schools.com/xml",
	)

	assert.Equal(t, expect, result)
}

func TestRun_CommandParseFailed(t *testing.T) {

	// ARRANGE
	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			"-a", // 存在しないフラグ
		},
		out,
	)

	// ASSERT
	require.Equal(t, NG, exitCode)

	expect := `xml2csv vdev (none)

Usage: xml2csv [flags]

Flags
  -i, --input string     XML input file path or directory or url
  -m, --mapping string   XML to CSV mapping file path or url
  -o, --output string    CSV output file path
  -b, --bom              CSV with BOM
  -h, --help             Help

unknown shorthand flag: 'a' in -a
`
	assert.Equal(t, expect, out.String())
}

func TestRun_Help(t *testing.T) {

	// ARRANGE
	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			"-h",
		},
		out,
	)

	// ASSERT
	require.Equal(t, OK, exitCode)

	expect := `xml2csv vdev (none)

Usage: xml2csv [flags]

Flags
  -i, --input string     XML input file path or directory or url
  -m, --mapping string   XML to CSV mapping file path or url
  -o, --output string    CSV output file path
  -b, --bom              CSV with BOM
  -h, --help             Help

`
	assert.Equal(t, expect, out.String())
}

func TestRun_NoneInput(t *testing.T) {

	// ARRANGE
	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			"-m", "xxx",
			"-o", "yyy",
		},
		out,
	)

	// ASSERT
	require.Equal(t, NG, exitCode)

	expect := `xml2csv vdev (none)

Usage: xml2csv [flags]

Flags
  -i, --input string     XML input file path or directory or url
  -m, --mapping string   XML to CSV mapping file path or url
  -o, --output string    CSV output file path
  -b, --bom              CSV with BOM
  -h, --help             Help

`
	assert.Equal(t, expect, out.String())
}

func TestRun_NoneMapping(t *testing.T) {

	// ARRANGE
	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			"-i", "xxx",
			"-o", "yyy",
		},
		out,
	)

	// ASSERT
	require.Equal(t, NG, exitCode)

	expect := `xml2csv vdev (none)

Usage: xml2csv [flags]

Flags
  -i, --input string     XML input file path or directory or url
  -m, --mapping string   XML to CSV mapping file path or url
  -o, --output string    CSV output file path
  -b, --bom              CSV with BOM
  -h, --help             Help

`
	assert.Equal(t, expect, out.String())
}

func TestRun_NoneOutput(t *testing.T) {

	// ARRANGE
	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			"-i", "xxx",
			"-m", "yyy",
		},
		out,
	)

	// ASSERT
	require.Equal(t, NG, exitCode)

	expect := `xml2csv vdev (none)

Usage: xml2csv [flags]

Flags
  -i, --input string     XML input file path or directory or url
  -m, --mapping string   XML to CSV mapping file path or url
  -o, --output string    CSV output file path
  -b, --bom              CSV with BOM
  -h, --help             Help

`
	assert.Equal(t, expect, out.String())
}

func TestRun_InputFileNotFound(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	// 作成はしない
	inputPath := filepath.Join(temp, "input.xml")

	mappingPath := createFile(t, temp, "mapping.json", `
	{
		"rowsPath": "//item",
		"columns": [
			{
				"header": "title",
				"valuePath": "/title"
			},
			{
				"header": "link",
				"valuePath": "/link"
			}
		]
	}`)

	outputPath := filepath.Join(temp, "output.csv")
	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			"-i", inputPath,
			"-m", mappingPath,
			"-o", outputPath,
		},
		out,
	)

	// ASSERT
	require.Equal(t, NG, exitCode)

	expect := inputPath + " is not found\n"
	assert.Equal(t, expect, out.String())
}

func TestRun_MappingFileNotFound(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	inputPath := "testdata/rss.xml"

	// 作成はしない
	mappingPath := filepath.Join(temp, "mapping.json")

	outputPath := filepath.Join(temp, "output.csv")
	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			"-i", inputPath,
			"-m", mappingPath,
			"-o", outputPath,
		},
		out,
	)

	// ASSERT
	require.Equal(t, NG, exitCode)

	expect := mappingPath + " is not found\n"
	assert.Equal(t, expect, out.String())
}

func TestRun_OutputFileDirNotFound(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	inputPath := "testdata/rss.xml"

	mappingPath := createFile(t, temp, "mapping.json", `
	{
		"rowsPath": "//item",
		"columns": [
			{
				"header": "title",
				"valuePath": "/title"
			},
			{
				"header": "link",
				"valuePath": "/link"
			}
		]
	}`)

	// 存在しないディレクトリを親に指定
	outputPath := filepath.Join(temp, "xxx", "output.csv")
	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			"-i", inputPath,
			"-m", mappingPath,
			"-o", outputPath,
		},
		out,
	)

	// ASSERT
	require.Equal(t, NG, exitCode)

	// OSによってエラーメッセージが異なるのでファイル名部分だけチェック
	expect := "open " + outputPath
	assert.Contains(t, out.String(), expect)
}

func TestRun_InvalidXPath_RowPath(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	inputPath := "testdata/rss.xml"

	mappingPath := createFile(t, temp, "mapping.json", `
	{
		"rowsPath": "item[",
		"columns": [
			{
				"header": "title",
				"valuePath": "/title"
			},
			{
				"header": "link",
				"valuePath": "/link"
			}
		]
	}`)

	outputPath := filepath.Join(temp, "output.csv")
	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			"-i", inputPath,
			"-m", mappingPath,
			"-o", outputPath,
		},
		out,
	)

	// ASSERT
	require.Equal(t, NG, exitCode)

	expect := "xpath 'item[' is failed: invalid streamElementXPath 'item[', err: expression must evaluate to a node-set\n"
	assert.Equal(t, expect, out.String())
}

func TestRun_InvalidXPath_ValuePath(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	inputPath := "testdata/rss.xml"

	mappingPath := createFile(t, temp, "mapping.json", `
	{
		"rowsPath": "//item",
		"columns": [
			{
				"header": "title",
				"valuePath": "/title["
			},
			{
				"header": "link",
				"valuePath": "/link"
			}
		]
	}`)

	outputPath := filepath.Join(temp, "output.csv")
	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			"-i", inputPath,
			"-m", mappingPath,
			"-o", outputPath,
		},
		out,
	)

	// ASSERT
	require.Equal(t, NG, exitCode)

	expect := "xpath '/title[' is failed: expression must evaluate to a node-set\n"
	assert.Equal(t, expect, out.String())
}

func TestRun_InvalidXPath_ValuePath_UseEvaluate(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	inputPath := "testdata/rss.xml"

	mappingPath := createFile(t, temp, "mapping.json", `
	{
		"rowsPath": "//item",
		"columns": [
			{
				"header": "title",
				"valuePath": "/title"
			},
			{
				"header": "link",
				"valuePath": "boolean(/link",
				"useEvaluate": true
			}
		]
	}`)

	outputPath := filepath.Join(temp, "output.csv")
	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			"-i", inputPath,
			"-m", mappingPath,
			"-o", outputPath,
		},
		out,
	)

	// ASSERT
	require.Equal(t, NG, exitCode)

	expect := "xpath 'boolean(/link' is failed: boolean(/link has an invalid token\n"
	assert.Equal(t, expect, out.String())
}

func TestRun_InvalidXML(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	inputPath := createFile(t, temp, "input.xml", `
	<root>
		<item id="1">
			<name>name1</name>
			<value>value1</value>
		</item>
	`)

	mappingPath := createFile(t, temp, "mapping.json", `
	{
		"rowsPath": "//item",
		"columns": [
			{
				"header": "name",
				"valuePath": "/name"
			},
			{
				"header": "value",
				"valuePath": "/value"
			}
		]
	}`)

	outputPath := filepath.Join(temp, "output.csv")
	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			"-i", inputPath,
			"-m", mappingPath,
			"-o", outputPath,
		},
		out,
	)

	// ASSERT
	require.Equal(t, NG, exitCode)

	expect := inputPath + " is failed: XML syntax error on line 7: unexpected EOF\n"
	assert.Equal(t, expect, out.String())
}

func TestRun_InvalidMappingJson(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	inputPath := "testdata/rss.xml"

	mappingPath := createFile(t, temp, "mapping.json", `
	{
		"rowsPath": "//item",
	}`)

	outputPath := filepath.Join(temp, "output.csv")
	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			"-i", inputPath,
			"-m", mappingPath,
			"-o", outputPath,
		},
		out,
	)

	// ASSERT
	require.Equal(t, NG, exitCode)

	expect := "invalid mapping format: invalid character '}' looking for beginning of object key string\n"
	assert.Equal(t, expect, out.String())
}

func TestConvertOne(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	inputPath := createFile(t, temp, "test.xml", `<root>
	<item id="1">
		<name>name1</name>
		<value>value1</value>
	</item>
	<item id="2">
		<name>name2</name>
		<value>value2,xx</value>
	</item>
	<item id="3">
		<name>name3</name>
	</item>
	</root>`)

	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	csv := customcsv.NewWriter(writer)

	mapping := Mapping{
		RowsPath: "//item",
		Columns: []Column{
			{Header: "id", ValuePath: "/@id"},
			{Header: "name", ValuePath: "/name"},
			{Header: "value", ValuePath: "/value"},
			{Header: "has value", ValuePath: "boolean(/value)", UseEvaluate: true},
		},
	}

	// ACT
	err := convertOne(inputPath, &mapping, csv)
	csv.Flush()

	// ASSERT
	require.NoError(t, err)

	result := b.String()

	expect := joinRows(
		"1,name1,value1,true",
		"2,name2,\"value2,xx\",true",
		"3,name3,,false",
	)

	assert.Equal(t, expect, result)
}

func TestLoadMapping_File(t *testing.T) {

	// ARRANGE/ACT
	result, err := loadMapping("mapping/rss.json")

	// ASSERT
	require.NoError(t, err)

	expect := &Mapping{
		RowsPath: "//item",
		Columns: []Column{
			{Header: "title", ValuePath: "/title"},
			{Header: "link", ValuePath: "/link"},
			{Header: "description", ValuePath: "/description"},
		},
	}

	assert.Equal(t, expect, result)
}

func TestLoadMapping_URL(t *testing.T) {

	// ARRANGE/ACT
	result, err := loadMapping("https://github.com/onozaty/xml2csv/raw/master/mapping/rss.json")

	// ASSERT
	require.NoError(t, err)

	expect := &Mapping{
		RowsPath: "//item",
		Columns: []Column{
			{Header: "title", ValuePath: "/title"},
			{Header: "link", ValuePath: "/link"},
			{Header: "description", ValuePath: "/description"},
		},
	}

	assert.Equal(t, expect, result)
}

func TestFindXML_Dir(t *testing.T) {

	// ARRANGE/ACT
	result, err := findXML("testdata/junit")

	// ASSERT
	require.NoError(t, err)

	expect := []string{
		filepath.Join("testdata", "junit", "TestCase1.xml"),
		filepath.Join("testdata", "junit", "TestCase2.xml"),
	}

	assert.Equal(t, expect, result)
}

func TestFindXML_Dir_Nest(t *testing.T) {

	// ARRANGE/ACT
	result, err := findXML("testdata")

	// ASSERT
	require.NoError(t, err)

	expect := []string{filepath.Join("testdata", "rss.xml")}

	assert.Equal(t, expect, result)
}

func TestFindXML_File(t *testing.T) {

	// ARRANGE/ACT
	result, err := findXML("testdata/rss.xml")

	// ASSERT
	require.NoError(t, err)

	expect := []string{"testdata/rss.xml"}

	assert.Equal(t, expect, result)
}

func TestFindXML_URL(t *testing.T) {

	// ARRANGE/ACT
	result, err := findXML("https://github.com/onozaty/xml2csv/raw/master/testdata/rss.xml")

	// ASSERT
	require.NoError(t, err)

	expect := []string{"https://github.com/onozaty/xml2csv/raw/master/testdata/rss.xml"}

	assert.Equal(t, expect, result)
}

func createFile(t *testing.T, dir string, name string, content string) string {

	file, err := os.Create(filepath.Join(dir, name))
	if err != nil {
		t.Fatal("craete file failed\n", err)
	}

	_, err = file.Write([]byte(content))
	if err != nil {
		t.Fatal("write file failed\n", err)
	}

	err = file.Close()
	if err != nil {
		t.Fatal("write file failed\n", err)
	}

	return file.Name()
}

func readBytes(t *testing.T, name string) []byte {

	bo, err := os.ReadFile(name)
	if err != nil {
		t.Fatal("read failed\n", err)
	}

	return bo
}

func readString(t *testing.T, name string) string {

	bo := readBytes(t, name)
	return string(bo)
}

func joinRows(rows ...string) string {
	return strings.Join(rows, "\r\n") + "\r\n"
}
