# xml2csv

[![GitHub license](https://img.shields.io/github/license/onozaty/xml2csv)](https://github.com/onozaty/xml2csv/blob/main/LICENSE)
[![Test](https://github.com/onozaty/xml2csv/actions/workflows/test.yaml/badge.svg)](https://github.com/onozaty/xml2csv/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/onozaty/xml2csv/branch/master/graph/badge.svg?token=QSRHZ6IMJF)](https://codecov.io/gh/onozaty/xml2csv)

xml2csv converts XML to CSV.  
You can easily define mappings for converts using XPath.

## Usage

```
$ xml2csv -i input.xml -m mapping.json -o output.csv
```

The arguments are as follows.

```
Usage: xml2csv [flags]

Flags
  -i, --input string       XML input file path or directory or url
  -m, --mapping string     XML to CSV mapping file path or url
  -o, --output string      CSV output file path
  -d, --delimiter string   (optional) CSV output delimiter (e.g. ';' or '\t' for tab) (default ",")
  -b, --bom                (optional) CSV with BOM
  -h, --help               Help
```

### Custom delimiter

Use the `-d` option with `'\t'` to output tab-separated values:

```
xml2csv -i input.xml -m mapping.json -o output.tsv -d '\t'
```

For semicolon-separated values:

```
xml2csv -i input.xml -m mapping.json -o output.csv -d ';'
```

### Using URL

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

### Homebrew (macOS/Linux)

```bash
brew install onozaty/tap/xml2csv
```

### Scoop (Windows)

```bash
scoop bucket add onozaty https://github.com/onozaty/scoop-bucket
scoop install xml2csv
```

### Binary Download

Download the latest binary from [GitHub Releases](https://github.com/onozaty/xml2csv/releases/latest).

## License

MIT

## Author

[onozaty](https://github.com/onozaty)
