package cccedictparser

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type FormatVersion = string

const (
	V1 FormatVersion = "V1"
	V2 FormatVersion = "V2"
)

type Ci struct {
	Fantizi   string
	Jiantizi  string
	Pinyin    []PinyinV2
	PinyinRaw string
	Gloss     []string
	FormatVersion
}

type Tone = uint8

const (
	None Tone = 0
	T1   Tone = 1
	T2   Tone = 2
	T3   Tone = 3
	T4   Tone = 4
	T5   Tone = 5
)

type PinyinType = uint8

const (
	Unknown  PinyinType = 0
	Normal   PinyinType = 1
	Alphabet PinyinType = 2
	Special  PinyinType = 3
)

type PinyinV1 struct {
	Sound string
	Tone  Tone
	Type  PinyinType
}

type PinyinV2 struct {
	Word []PinyinV1
}

const section_traditional = 1
const section_simplified = 3
const section_transition_pinyin = 4
const section_pinyin = 5
const section_transition_gloss = 6
const section_gloss = 7

// true if not out of bounds
// false if out of bounds
func tryPeak(arr []string, index int) (string, bool) {
	if index >= len(arr) {
		return "", false
	}
	return arr[index], true
}

func (p PinyinV1) String() string {
	return fmt.Sprintf("PinyinV1{Sound:\"%s\", Tone: %d, Type: %d}", p.Sound, p.Tone, p.Type)
}

func (p PinyinV2) String() string {
	strs := []string{}
	for _, v := range p.Word {
		strs = append(strs, v.String())
	}
	return fmt.Sprintf("PinyinV2{Word:[%s]}", strings.Join(strs, ", "))
}

func (ci Ci) String() string {
	return fmt.Sprintf("Ci{Fantizi:\"%s\", Jiantizi:\"%s\", Pinyin:%s, PinyinRaw:\"%s\", Gloss:[%s], FormatVersion:%s}", ci.Fantizi, ci.Jiantizi, "", strings.Join(ci.Gloss, ", "), ci.PinyinRaw, ci.FormatVersion)
}

func pinyinV1StrToPinyin(pys string) ([]PinyinV2, error) {
	items := strings.Split(pys, " ")

	pyItems := []PinyinV2{}

	//ci2 shu1
	for _, v := range items {
		//ci2
		runes := []rune{}
		for _, c := range v {
			runes = append(runes, c)
		}

		py, err := getPyV1ForPySegmentRunes(runes)

		if err != nil {
			return []PinyinV2{}, err
		}

		pyItems = append(pyItems, PinyinV2{
			Word: []PinyinV1{py},
		})
	}

	return pyItems, nil
}

func getPyV1ForPySegmentRunes(runes []rune) (PinyinV1, error) {
	hasTone := false
	var tone Tone
	if toneVal, err := strconv.Atoi(string(runes[len(runes)-1])); err == nil && tone <= 5 && tone >= 1 {
		hasTone = true
		tone = uint8(toneVal)
	}

	var sound string
	if !hasTone {
		sound = string(runes)
	} else {
		sound = string(runes[:len(runes)-1])
	}

	isAlphabetic, err := regexp.MatchString(`[a-zA-Z]`, sound)

	var t PinyinType
	if err != nil {
		return PinyinV1{}, errors.New("malformed pinyin v1")
	} else if isAlphabetic && tone != None {
		t = Normal
	} else if isAlphabetic {
		t = Alphabet
	} else {
		t = Special
	}

	py := PinyinV1{
		Sound: sound,
		Tone:  tone,
		Type:  t,
	}

	return py, nil
}

func pinyinV2StrToPinyin(pys string) ([]PinyinV2, error) {
	words := strings.Split(pys, " ")

	v2List := []PinyinV2{}

	//Ping2guo3 shou3ji1
	for _, word := range words {
		wordsForPyV2 := []PinyinV1{}

		//Ping2guo3
		runesBuilder := []rune{}
		pyItems := [][]rune{}
		for _, c := range word {
			runesBuilder = append(runesBuilder, c)
			if _, err := strconv.Atoi(string(c)); err != nil {
				pyItems = append(pyItems, runesBuilder)
				runesBuilder = []rune{}
			}
		}

		for _, pyItem := range pyItems {
			item, err := getPyV1ForPySegmentRunes(pyItem)

			if err != nil {
				return []PinyinV2{}, err
			}

			wordsForPyV2 = append(wordsForPyV2, item)
		}

		v2List = append(v2List, PinyinV2{
			Word: wordsForPyV2,
		})
	}

	return v2List, nil
}

// Traditional Simplified [[pin1yin1]] /gloss; gloss; .../gloss; gloss; .../
func ParseLine(line string) (Ci, error) {
	chars := strings.Split(line, "")

	traditionalDelimit := " "
	simplifiedDelimit := " "
	pinyinEnd := "]"

	glossStart := "/"

	var fantizi string
	var jiantizi string
	var pinyin string
	var gloss []string

	var builder strings.Builder
	currentSection := section_traditional

	pyVersion := V1

	for i, r := range line {

		if currentSection == section_traditional {

			if string(r) == traditionalDelimit {
				fantizi = builder.String()
				builder.Reset()
				currentSection = section_simplified
			} else {
				builder.WriteRune(r)
			}

		} else if currentSection == section_simplified {

			if string(r) == simplifiedDelimit {
				jiantizi = builder.String()
				builder.Reset()
				currentSection = section_transition_pinyin
			} else {
				builder.WriteRune(r)
			}

		} else if currentSection == section_transition_pinyin {

			if next, ok := tryPeak(chars, i+1); ok && string(r) == "[" {
				if next != "[" {
					currentSection = section_pinyin
				}
			}

		} else if currentSection == section_pinyin {
			if string(r) == pinyinEnd {
				pinyin = builder.String()
				builder.Reset()
				currentSection = section_transition_gloss
			} else {
				builder.WriteRune(r)
			}
		} else if currentSection == section_transition_gloss {
			if string(r) == pinyinEnd || string(r) == " " {
				//nothing
			} else if string(r) == glossStart {
				currentSection = section_gloss
			} else {
				return Ci{}, fmt.Errorf("failed to read gloss for line (%s)", line)
			}
		} else if currentSection == section_gloss {
			if string(r) == "/" {
				gloss = append(gloss, builder.String())
				builder.Reset()
			} else {
				builder.WriteRune(r)
			}
		}
	}

	var py []PinyinV2
	var err error
	if pyVersion == V1 {
		py, err = pinyinV1StrToPinyin(pinyin)
	} else {
		py, err = pinyinV2StrToPinyin(pinyin)
	}

	if err != nil {
		return Ci{}, err
	}

	return Ci{
		Fantizi:       fantizi,
		Jiantizi:      jiantizi,
		Pinyin:        py,
		Gloss:         gloss,
		FormatVersion: pyVersion,
	}, nil
}
