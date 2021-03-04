# xml2csv

xml2csv converts XML to CSV.  
You can easily define mappings for converts using XPath.

## Usage

```
$ xml2csv -i input.xml -m mapping.json -o output.csv
```

The arguments are as follows.

```
Usage of xml2csv:
  -i, --input string     XML input file path or directory or url
  -m, --mapping string   XML to CSV mapping file path or url
  -o, --output string    CSV output file path
  -b, --bom              CSV with BOM
  -h, --help             Help
```

XML and mapping files can be specified by URL.

```
xml2csv -i https://github.com/onozaty/xml2csv/raw/master/testdata/rss.xml -m https://github.com/onozaty/xml2csv/raw/master/mapping/rss.json -o output.csv
```

## Mapping

The conversion mapping definition is written in JSON.    
Specify the position on the XML with XPath.

```json
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
        },
        {
            "header": "description",
            "valuePath": "/description"
        }
    ]
}
```

* `rowsPath` : XPath to get as a rows.
* `columns` : Definition of each column.
    * `header` : CSV header.
    * `valuePath` : XPath to get as a value.
    * `useEvaluate` : Specify `true` when using an expression with `valuePath`. For example, when using `sum()` or `not()`, `boolean()`.

[antchfx/xpath](https://github.com/antchfx/xpath) is used in xml2csv.  
See below for supported XPath.

* https://github.com/antchfx/xpath#supported-features

Please refer to the sample below.

* https://github.com/onozaty/xml2csv/tree/master/mapping

## Install

You can download the binary from the following.

* https://github.com/onozaty/xml2csv/releases/latest

## License

MIT

## Author

[onozaty](https://github.com/onozaty)
