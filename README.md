# Simple cc-cedict parser written in Go

This package handles parsing lines in the cc-cedict v1 and v2 formats and structures the data.

# Install

To install the library:

`go get github.com/xDestx/cc-cedict-reader`

To install the command:

`go install github.com/xDestx/cc-cedict-reader/cmd/cccedict-parser`

# Usage

## Library

Parse a line with `ParseLine(line string) (Ci, error)`.

From `examples/examples.go`
```go
package examples

import cccedictparser "github.com/xDestx/cc-cedict-reader"

func ex_read_line() {
	ci, err := cccedictparser.ParseLine("同床異夢 同床异梦 [tong2 chuang2 yi4 meng4] /lit. to share the same bed with different dreams (idiom); ostensible partners with different agendas/strange bedfellows/marital dissension/")

	if err != nil {
		println(err.Error())
	}

	// term
	trad := ci.Fantizi
	simp := ci.Jiantizi
	// structured pinyin
	pinyin := ci.Pinyin
	// raw string pinyin
	pyStr := ci.PinyinRaw
	// gloss
	meaning := ci.Gloss
	// cc cedict version detected
	ver := ci.FormatVersion

	//print as string
	println(ci.String())
	println(trad)
	println(simp)
	for _, v := range pinyin {
		println(v.String())
	}
	println(pyStr)
	for _, v := range meaning {
		println(v)
	}
	println(ver)
}
```

## Command

The command reads from stdin and outputs to stdout.

Example:

`echo "各得其所 各得其所 [ge4 de2 qi2 suo3] /(idiom) each in the correct place; each is provided for/" | cc-cedict-reader > test.txt`

test.txt
```
Ci{Fantizi:"各得其所", Jiantizi:"各得其所", Pinyin:PinyinV2{Word:[PinyinV1{Sound:"ge", Tone: 4, Type: 1}]}, PinyinV2{Word:[PinyinV1{Sound:"de", Tone: 2, Type: 1}]}, PinyinV2{Word:[PinyinV1{Sound:"qi", Tone: 2, Type: 1}]}, PinyinV2{Word:[PinyinV1{Sound:"suo", Tone: 3, Type: 1}]}, PinyinRaw:"ge4 de2 qi2 suo3", Gloss:[(idiom) each in the correct place; each is provided for], FormatVersion:V1}
```
