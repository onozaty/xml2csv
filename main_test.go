package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"reflect"
	"strings"
	"testing"

	"github.com/antchfx/xmlquery"
)

func TestConvert(t *testing.T) {

	var b bytes.Buffer
	writer := bufio.NewWriter(&b)

	mapping := Mapping{
		RowsPath: "//item",
		Columns: []Column{
			Column{Header: "title", ValuePath: "/title"},
			Column{Header: "link", ValuePath: "/link"},
		},
	}

	convert([]string{"testdata/rss.xml"}, mapping, writer)

	result := string(b.Bytes())

	expect := "title,link\n" +
		"RSS Tutorial,https://www.w3schools.com/xml/xml_rss.asp\n" +
		"XML Tutorial,https://www.w3schools.com/xml\n"

	if result != expect {
		t.Fatalf("failed test\n%s", result)
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
		panic(err)
	}

	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	csv := csv.NewWriter(writer)

	mapping := Mapping{
		RowsPath: "//item",
		Columns: []Column{
			Column{Header: "id", ValuePath: "/@id"},
			Column{Header: "name", ValuePath: "/name"},
			Column{Header: "value", ValuePath: "/value"},
		},
	}

	convertOne(doc, mapping, csv)

	csv.Flush()

	result := string(b.Bytes())

	expect := "1,name1,value1\n" +
		"2,name2,\"value2,xx\"\n" +
		"3,name3,\n"

	if result != expect {
		t.Fatalf("failed test\n%s", result)
	}
}

func TestLoadMapping(t *testing.T) {

	result := loadMapping("mapping/rss.json")

	expect := Mapping{
		RowsPath: "//item",
		Columns: []Column{
			Column{Header: "title", ValuePath: "/title"},
			Column{Header: "link", ValuePath: "/link"},
			Column{Header: "description", ValuePath: "/description"},
		},
	}

	if !reflect.DeepEqual(result.Columns, expect.Columns) {
		t.Fatalf("failed test\n%s", result)
	}
}
