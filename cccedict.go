package cccedictparser

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/text/unicode/norm"
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
func tryPeak(arr []rune, index int) (rune, bool) {
	if index >= len(arr) {
		return 0, false
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

func pyV2ArrStr(pyv2arr []PinyinV2) string {
	items := []string{}
	for _, v := range pyv2arr {
		items = append(items, v.String())
	}
	return strings.Join(items, ", ")
}

func (ci Ci) String() string {
	return fmt.Sprintf("Ci{Fantizi:\"%s\", Jiantizi:\"%s\", Pinyin:%s, PinyinRaw:\"%s\", Gloss:[%s], FormatVersion:%s}", ci.Fantizi, ci.Jiantizi, pyV2ArrStr(ci.Pinyin), ci.PinyinRaw, strings.Join(ci.Gloss, ", "), ci.FormatVersion)
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

	if len(runes) == 0 {
		return PinyinV1{}, errors.New("no runes provided")
	}

	if string(runes) == "xx5" {
		return PinyinV1{
			Sound: "xx",
			Tone:  T5,
			Type:  Unknown,
		}, nil
	}

	if toneVal, err := strconv.Atoi(string(runes[len(runes)-1])); err == nil && toneVal <= 5 && toneVal >= 1 {
		hasTone = true
		tone = uint8(toneVal)
	}

	var sound string
	if !hasTone {
		sound = string(runes)
	} else {
		sound = string(runes[:len(runes)-1])
	}

	startsWithBrackets := strings.HasPrefix(sound, "{")
	endsWithBrakets := strings.HasSuffix(sound, "}")

	if startsWithBrackets != endsWithBrakets {
		return PinyinV1{}, errors.New("malformed pinyin")
	}

	if ns := strings.TrimSuffix(sound, "-"); ns != "" {
		sound = ns
	}

	isAlphabetic, err := regexp.MatchString(`^[a-zA-Z]+$`, sound)
	if err != nil {
		return PinyinV1{}, errors.New("malformed pinyin v1")
	}

	hasNumber, err := regexp.MatchString(`\d`, sound)

	if err != nil {
		return PinyinV1{}, errors.New("malformed pinyin v1")
	}

	if hasNumber && len(sound) != 1 {
		return PinyinV1{}, errors.New("malformed pinyin v1")
	}

	var t PinyinType
	if isAlphabetic && tone != None {
		t = Normal
	} else if isAlphabetic {
		t = Alphabet
	} else {
		t = Special
	}

	finSound := strings.Replace(sound, "u:", "v", -1)

	py := PinyinV1{
		Sound: finSound,
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
		openBracket := false
		for _, c := range word {
			runesBuilder = append(runesBuilder, c)
			sc := string(c)

			if sc == "{" {
				openBracket = true
			}

			_, err := strconv.Atoi(sc)

			if err == nil || (openBracket && sc == "}") || sc == "-" {
				removeSuffix := 0
				if openBracket || sc == "-" {
					removeSuffix = 1
				}
				removePrefix := 0
				if openBracket {
					removePrefix = 1
				}

				cleanedRunes := runesBuilder[removePrefix : len(runesBuilder)-removeSuffix]
				if len(cleanedRunes) != 0 {
					pyItems = append(pyItems, cleanedRunes)
				} else {
					//likely a special character
					pyItems = append(pyItems, runesBuilder)
				}
				runesBuilder = []rune{}
				openBracket = false
			}
		}

		if len(runesBuilder) != 0 {
			pyItems = append(pyItems, runesBuilder)
			runesBuilder = []rune{}
		}

		for _, pyItem := range pyItems {
			item, err := getPyV1ForPySegmentRunes(pyItem)

			if err != nil {
				return []PinyinV2{}, err
			}

			if item.Sound == `·` {
				return []PinyinV2{}, errors.New("malformed pinyin v2 - no dots")
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
	traditionalDelimit := " "
	simplifiedDelimit := " "
	pinyinStart := "["
	pinyinEnd := "]"

	glossStart := "/"

	var fantizi string
	var jiantizi string
	var pinyin string
	var gloss []string

	var builder strings.Builder
	currentSection := section_traditional

	pyVersion := V1
	pyOpenBracketCount := 0
	pyCloseBracketCount := 0

	lineRunes := []rune{}

	for _, r := range line {
		lineRunes = append(lineRunes, r)
	}

	for i, r := range lineRunes {

		if currentSection == section_traditional {

			if string(r) == pinyinStart {
				return Ci{}, errors.New("found pinyin section before completing traditional section")
			}

			if string(r) == traditionalDelimit {
				fantizi = builder.String()
				builder.Reset()
				currentSection = section_simplified
			} else {
				builder.WriteRune(r)
			}

		} else if currentSection == section_simplified {

			if string(r) == pinyinStart {
				return Ci{}, errors.New("found pinyin section before completing simplified section")
			}

			if string(r) == simplifiedDelimit {
				jiantizi = builder.String()
				builder.Reset()
				currentSection = section_transition_pinyin
			} else {
				builder.WriteRune(r)
			}

		} else if currentSection == section_transition_pinyin {

			if string(r) == glossStart {
				return Ci{}, errors.New("found gloss section before pinyin section")
			}

			if string(r) == pinyinStart {
				pyOpenBracketCount++
			}

			if next, ok := tryPeak(lineRunes, i+1); ok && string(r) == pinyinStart {
				if string(next) != pinyinStart {
					currentSection = section_pinyin
				}
			}

		} else if currentSection == section_pinyin {
			if string(r) == pinyinEnd {
				pyCloseBracketCount++
				pinyin = builder.String()
				builder.Reset()
				currentSection = section_transition_gloss
			} else {
				builder.WriteRune(r)
			}
		} else if currentSection == section_transition_gloss {
			if string(r) == pinyinEnd || string(r) == " " {
				if string(r) == pinyinEnd {
					pyCloseBracketCount++
				}
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

	if pyOpenBracketCount != pyCloseBracketCount {
		return Ci{}, fmt.Errorf("malformed pinyin (cannot determine version) (%d %d)", pyOpenBracketCount, pyCloseBracketCount)
	}

	switch pyOpenBracketCount {
	case 1:
		pyVersion = V1
	case 2:
		pyVersion = V2
	default:
		return Ci{}, errors.New("malformed pinyin (unrecognized version)")
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

	if fantizi == "" || jiantizi == "" {
		return Ci{}, errors.New("no traditional/simplified found")
	}

	if pinyin == "" || len(py) == 0 {
		return Ci{}, errors.New("no pinyin found")
	}

	if len(gloss) == 0 {
		return Ci{}, errors.New("no gloss found")
	}

	for _, v := range pinyin {
		if string(norm.NFD.Bytes([]byte(string(v)))) != string(v) {
			// Really struggling to detect this
			return Ci{}, errors.New("malformed pinyin - no diacritics")
		}
	}

	return Ci{
		Fantizi:       fantizi,
		Jiantizi:      jiantizi,
		Pinyin:        py,
		PinyinRaw:     pinyin,
		Gloss:         gloss,
		FormatVersion: pyVersion,
	}, nil
}
