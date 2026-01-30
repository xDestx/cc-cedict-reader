package cccedictparser

import (
	"errors"
	"regexp"
	"strings"
)

type Ci struct {
	Fantizi  string
	Jiantizi string
	Pinyin   string
	Gloss    []string
}

const pinyinRegex = `\[\[?((?:(?:\w+:?·?\d?)+\s*,?·?\s*)+)\]\]?`

func joinAllStrings(s []string, joinStr string) string {
	var sb strings.Builder
	for i := range s {
		sb.WriteString(s[i])
		sb.WriteString(joinStr)
	}

	return sb.String()
}

// Traditional Simplified [[pinb1yin1]] /gloss; gloss; .../gloss; gloss; .../
func ParseLine(line string) (Ci, error) {
	items := strings.Split(line, " ")
	r, err := regexp.Compile(pinyinRegex)

	pinyinAndGlossStr := joinAllStrings(items[2:], " ")
	pinyinGlossSplit := strings.Split(pinyinAndGlossStr, " /")
	pinyinOnly := pinyinGlossSplit[0]
	fullGloss := joinAllStrings(pinyinGlossSplit[1:], " /")
	fullGlossStripped := fullGloss[0:][:len(fullGloss)-4]
	g := strings.Split(fullGlossStripped, "/")

	if err != nil {
		return Ci{}, err
	}

	pinyinMatch := r.FindStringSubmatch(pinyinOnly)

	if len(pinyinMatch) == 0 {
		return Ci{}, errors.New("error parsing pinyin for line: " + line)
	}

	t := items[0]
	s := items[1]
	p := pinyinMatch[1]

	return Ci{
		Fantizi:  t,
		Jiantizi: s,
		Pinyin:   p,
		Gloss:    g,
	}, nil
}
