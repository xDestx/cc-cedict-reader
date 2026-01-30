package cccedictparser

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type tone = int

const (
	_ tone = iota
	Yi
	Er
	San
	Si
	Wu
)

type PinyinSyllable struct {
	Romanization string
	Tone         tone
}

type PinyinChunk = []PinyinSyllable

type Ci struct {
	Fantizi  string
	Jiantizi string
	Pinyin   []PinyinChunk
	Gloss    []string
}

const parseLineRegex = `^([\p{Han}，·]+) ([\p{Han}，·]+) \[\[?((?:(?:\w+\d)+\s*,?\s*)+)\]\]? \/([\w\(\)\s,\/\.\--–;\p{Han}，·]+)/$`
const pinyinRegex = `([a-zA-Z]+)([12345])`

func pinyinChunkFromStr(pyChunkStr string) (PinyinChunk, error) {
	r, err := regexp.Compile(pinyinRegex)

	m := r.FindAllStringSubmatch(pyChunkStr, -1)

	if err != nil {
		return []PinyinSyllable{}, err
	}

	pinyinYiGes := []PinyinSyllable{}
	var actualStrBuilder strings.Builder
	for vArr := range m {
		t, err := StringToTone(m[vArr][2])

		if err != nil {
			return []PinyinSyllable{}, fmt.Errorf("error parsing pinyin for (%s) error: %s", pyChunkStr, err.Error())
		}

		newPy := PinyinSyllable{
			Romanization: m[vArr][1],
			Tone:         t,
		}

		pinyinYiGes = append(pinyinYiGes, newPy)

		actualStrBuilder.WriteString(PinyinSyllableString(newPy))
	}

	actualStr := actualStrBuilder.String()

	if pyChunkStr != actualStr {
		return []PinyinSyllable{}, fmt.Errorf("mismatch between parsed and actual value. Expected (%s) Was (%s)", pyChunkStr, actualStr)
	}

	return pinyinYiGes, nil
}

func PinyinFromString(pystr string) ([]PinyinChunk, error) {
	pinyins := strings.Split(pystr, " ")

	py := []PinyinChunk{}

	for i := range pinyins {
		v, err := pinyinChunkFromStr(pinyins[i])

		if err != nil {
			return [][]PinyinSyllable{}, fmt.Errorf("error parsing pinyin for (%s). %s", pystr, err)
		}

		py = append(py, v)
	}

	return py, nil
}

// Traditional Simplified [[pinb1yin1]] /gloss; gloss; .../gloss; gloss; .../
func ParseLine(line string) (Ci, error) {
	r, err := regexp.Compile(parseLineRegex)

	if err != nil {
		return Ci{}, err
	}

	s := r.FindStringSubmatch(line)

	if len(s) == 0 {
		return Ci{}, errors.New("error parsing line: " + line)
	}

	py, err := PinyinFromString(s[3])

	if err != nil {
		return Ci{}, err
	}

	gloss := strings.Split(s[4], "/")

	return Ci{
		Fantizi:  s[1],
		Jiantizi: s[2],
		Pinyin:   py,
		Gloss:    gloss,
	}, nil
}

func ToneToString(t tone) string {
	switch t {
	case Yi:
		return "1"
	case Er:
		return "2"
	case San:
		return "3"
	case Si:
		return "4"
	case Wu:
		return "5"
	default:
		return "?"
	}
}

func StringToTone(s string) (tone, error) {
	switch s {
	case "1":
		return Yi, nil
	case "2":
		return Er, nil
	case "3":
		return San, nil
	case "4":
		return Si, nil
	case "5":
		return Wu, nil
	default:
		return 0, errors.New("unable to identify tone " + s)
	}
}

func PinyinSyllableString(p PinyinSyllable) string {
	return p.Romanization + ToneToString(p.Tone)
}

func PinyinChunkString(py PinyinChunk) string {
	var r strings.Builder

	for i := 0; i < len(py); i++ {
		r.WriteString(PinyinSyllableString(py[i]))
	}

	return r.String()
}

func PinyinString(p []PinyinChunk) string {
	var r strings.Builder
	l := len(p)

	for i := range l {
		r.WriteString(PinyinChunkString(p[i]))
		if (i + 1) < l {
			r.WriteString(" ")
		}
	}

	return r.String()
}
