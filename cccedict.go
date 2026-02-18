package cccedictparser

import (
	"fmt"
	"strings"
)

type FormatVersion = string

const (
	V1 FormatVersion = "V1"
	V2 FormatVersion = "V2"
)

type Ci struct {
	Fantizi  string
	Jiantizi string
	Pinyin   string
	Gloss    []string
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

	return Ci{
		Fantizi:  fantizi,
		Jiantizi: jiantizi,
		Pinyin:   pinyin,
		Gloss:    gloss,
	}, nil
}
