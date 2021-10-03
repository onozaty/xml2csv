package main

import (
	"bufio"
	"bytes"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/antchfx/xmlquery"
	"github.com/onozaty/go-customcsv"
)

func TestConvert_File(t *testing.T) {

	var b bytes.Buffer
	writer := bufio.NewWriter(&b)

	mapping := Mapping{
		RowsPath: "//item",
		Columns: []Column{
			{Header: "title", ValuePath: "/title"},
			{Header: "link", ValuePath: "/link"},
		},
	}

	err := convert([]string{"testdata/rss.xml"}, &mapping, writer)
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := b.String()

	expect := "title,link\r\n" +
		"RSS Tutorial,https://www.w3schools.com/xml/xml_rss.asp\r\n" +
		"XML Tutorial,https://www.w3schools.com/xml\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestConvert_URL(t *testing.T) {

	var b bytes.Buffer
	writer := bufio.NewWriter(&b)

	mapping := Mapping{
		RowsPath: "//item",
		Columns: []Column{
			{Header: "title", ValuePath: "/title"},
			{Header: "link", ValuePath: "/link"},
		},
	}

	err := convert([]string{"https://github.com/onozaty/xml2csv/raw/master/testdata/rss.xml"}, &mapping, writer)
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := b.String()

	expect := "title,link\r\n" +
		"RSS Tutorial,https://www.w3schools.com/xml/xml_rss.asp\r\n" +
		"XML Tutorial,https://www.w3schools.com/xml\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestConvertOne(t *testing.T) {

	xml := `<root>
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
	</root>`

	doc, err := xmlquery.Parse(strings.NewReader(xml))
	if err != nil {
		t.Fatal("failed test\n", err)
	}

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

	convertOne(doc, &mapping, csv)

	csv.Flush()

	result := b.String()

	expect := "1,name1,value1,true\r\n" +
		"2,name2,\"value2,xx\",true\r\n" +
		"3,name3,,false\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestLoadMapping_File(t *testing.T) {

	result, err := loadMapping("mapping/rss.json")
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	expect := &Mapping{
		RowsPath: "//item",
		Columns: []Column{
			{Header: "title", ValuePath: "/title"},
			{Header: "link", ValuePath: "/link"},
			{Header: "description", ValuePath: "/description"},
		},
	}

	if !reflect.DeepEqual(result, expect) {
		t.Fatal("failed test\n", result)
	}
}

func TestLoadMapping_URL(t *testing.T) {

	result, err := loadMapping("https://github.com/onozaty/xml2csv/raw/master/mapping/rss.json")
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	expect := &Mapping{
		RowsPath: "//item",
		Columns: []Column{
			{Header: "title", ValuePath: "/title"},
			{Header: "link", ValuePath: "/link"},
			{Header: "description", ValuePath: "/description"},
		},
	}

	if !reflect.DeepEqual(result, expect) {
		t.Fatal("failed test\n", result)
	}
}

func TestFindXML_Dir(t *testing.T) {

	result, err := findXML("testdata/junit")
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	expect := []string{
		filepath.Join("testdata", "junit", "TestCase1.xml"),
		filepath.Join("testdata", "junit", "TestCase2.xml"),
	}

	if !reflect.DeepEqual(result, expect) {
		t.Fatal("failed test\n", result)
	}
}

func TestFindXML_Dir_Nest(t *testing.T) {

	result, err := findXML("testdata")
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	expect := []string{filepath.Join("testdata", "rss.xml")}

	if !reflect.DeepEqual(result, expect) {
		t.Fatal("failed test\n", result)
	}
}

func TestFindXML_File(t *testing.T) {

	result, err := findXML("testdata/rss.xml")
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	expect := []string{"testdata/rss.xml"}

	if !reflect.DeepEqual(result, expect) {
		t.Fatal("failed test\n", result)
	}
}

func TestFindXML_URL(t *testing.T) {

	result, err := findXML("https://github.com/onozaty/xml2csv/raw/master/testdata/rss.xml")
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	expect := []string{"https://github.com/onozaty/xml2csv/raw/master/testdata/rss.xml"}

	if !reflect.DeepEqual(result, expect) {
		t.Fatal("failed test\n", result)
	}
}
