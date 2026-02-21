package cccedictparser

import (
	"errors"
	"fmt"
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

type LineParser interface {
	ParseLine(line string) (Ci, error)
}

type basicLineParser struct {
	Pym map[string]bool
}

var full_pinyin_list = []string{
	`a`, `ba`, `pa`, `ma`, `fa`, `da`, `ta`, `na`, `la`, `ga`, `ka`, `ha`, `zha`, `cha`,
	`sha`, `za`, `ca`, `sa`, `ai`, `bai`, `pai`, `mai`, `dai`, `tai`, `nai`, `lai`, `gai`, `kai`,
	`hai`, `zhai`, `chai`, `shai`, `zai`, `cai`, `sai`, `an`, `ban`, `pan`, `man`, `fan`, `dan`, `tan`,
	`nan`, `lan`, `gan`, `kan`, `han`, `zhan`, `chan`, `shan`, `ran`, `zan`, `can`, `san`, `ang`, `bang`,
	`pang`, `mang`, `fang`, `dang`, `tang`, `nang`, `lang`, `gang`, `kang`, `hang`, `zhang`, `chang`, `shang`, `rang`,
	`zang`, `cang`, `sang`, `ao`, `bao`, `pao`, `mao`, `dao`, `tao`, `nao`, `lao`, `gao`, `kao`, `hao`,
	`zhao`, `chao`, `shao`, `rao`, `zao`, `cao`, `sao`, `e`, `me`, `de`, `te`, `ne`, `le`, `ge`,
	`ke`, `he`, `zhe`, `che`, `she`, `re`, `ze`, `ce`, `se`, `ei`, `bei`, `pei`, `mei`, `fei`,
	`dei`, `nei`, `lei`, `gei`, `hei`, `shei`, `zei`, `en`, `ben`, `pen`, `men`, `fen`, `den`, `nen`,
	`gen`, `ken`, `hen`, `zhen`, `chen`, `shen`, `ren`, `zen`, `cen`, `sen`, `beng`, `peng`, `meng`, `feng`,
	`deng`, `teng`, `neng`, `leng`, `geng`, `keng`, `heng`, `zheng`, `cheng`, `sheng`, `reng`, `zeng`, `ceng`, `seng`,
	`er`, `yi`, `bi`, `pi`, `mi`, `di`, `ti`, `ni`, `li`, `ji`, `qi`, `xi`, `zhi`, `chi`,
	`shi`, `ri`, `zi`, `ci`, `si`, `ya`, `dia`, `lia`, `jia`, `qia`, `xia`, `yan`, `bian`, `pian`,
	`mian`, `dian`, `tian`, `nian`, `lian`, `jian`, `qian`, `xian`, `yang`, `niang`, `liang`, `jiang`, `qiang`, `xiang`,
	`yao`, `biao`, `piao`, `miao`, `diao`, `tiao`, `niao`, `liao`, `jiao`, `qiao`, `xiao`, `ye`, `bie`, `pie`,
	`mie`, `die`, `tie`, `nie`, `lie`, `jie`, `qie`, `xie`, `yin`, `bin`, `pin`, `min`, `nin`, `lin`,
	`jin`, `qin`, `xin`, `ying`, `bing`, `ping`, `ming`, `ding`, `ting`, `ning`, `ling`, `jing`, `qing`, `xing`,
	`yo`, `yong`, `jiong`, `qiong`, `xiong`, `you`, `miu`, `diu`, `niu`, `liu`, `jiu`, `qiu`, `xiu`, `o`,
	`bo`, `po`, `mo`, `fo`, `lo`, `weng`, `dong`, `tong`, `nong`, `long`, `gong`, `kong`, `hong`, `zhong`,
	`chong`, `rong`, `zong`, `cong`, `song`, `ou`, `pou`, `mou`, `fou`, `dou`, `tou`, `nou`, `lou`, `gou`,
	`kou`, `hou`, `zhou`, `chou`, `shou`, `rou`, `zou`, `cou`, `sou`, `wu`, `bu`, `pu`, `mu`, `fu`,
	`du`, `tu`, `nu`, `lu`, `gu`, `ku`, `hu`, `zhu`, `chu`, `shu`, `ru`, `zu`, `cu`, `su`,
	`wa`, `gua`, `kua`, `hua`, `zhua`, `shua`, `wai`, `guai`, `kuai`, `huai`, `chuai`, `shuai`, `wan`, `duan`,
	`tuan`, `nuan`, `luan`, `guan`, `kuan`, `huan`, `zhuan`, `chuan`, `shuan`, `ruan`, `zuan`, `cuan`, `suan`, `wang`,
	`guang`, `kuang`, `huang`, `zhuang`, `chuang`, `shuang`, `yue`, `nve`, `lve`, `jue`, `que`, `xue`, `wei`, `dui`,
	`tui`, `gui`, `kui`, `hui`, `zhui`, `chui`, `shui`, `rui`, `zui`, `cui`, `sui`, `wen`, `dun`, `tun`,
	`lun`, `gun`, `kun`, `hun`, `zhun`, `chun`, `shun`, `run`, `zun`, `cun`, `sun`, `wo`, `duo`, `tuo`,
	`nuo`, `luo`, `guo`, `kuo`, `huo`, `zhuo`, `chuo`, `shuo`, `ruo`, `zuo`, `cuo`, `suo`, `yu`, `nv`,
	`lv`, `ju`, `qu`, `xu`, `yuan`, `juan`, `quan`, `xuan`, `yun`, `jun`, `qun`, `xun`, `r`,
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
	strs := make([]string, 0, len(p.Word))
	for _, v := range p.Word {
		strs = append(strs, v.String())
	}
	return fmt.Sprintf("PinyinV2{Word:[%s]}", strings.Join(strs, ", "))
}

func pyV2ArrStr(pyv2arr []PinyinV2) string {
	items := make([]string, 0, len(pyv2arr))
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

	pyItems := make([]PinyinV2, 0, len(items))

	runes := make([]rune, 0, len(items))

	//ci2 shu1
	for _, v := range items {
		//ci2
		runes = runes[:0]
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

func getTone(s rune) (Tone, error) {

	switch s {
	case 49:
		return T1, nil
	case 50:
		return T2, nil
	case 51:
		return T3, nil
	case 52:
		return T4, nil
	case 53:
		return T5, nil
	default:
		return None, errors.New("unrecognized tone")
	}
}

func soundIsAlphabetic(str string) bool {
	for _, v := range str {
		if !((v >= 97 && v <= 122) || (v >= 65 && v <= 90)) {
			return false
		}
	}
	return true
}

func soundHasNumber(str string) bool {
	for _, v := range str {
		if v >= 48 && v <= 57 {
			return true
		}
	}
	return false
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

	if toneVal, err := getTone(runes[len(runes)-1]); err == nil {
		hasTone = true
		tone = toneVal
	}

	var sound string
	if !hasTone {
		sound = string(runes)
	} else {
		sound = string(runes[:len(runes)-1])
	}

	sound = strings.Replace(sound, "u:", "v", -1)

	startsWithBrackets := strings.HasPrefix(sound, "{")
	endsWithBrakets := strings.HasSuffix(sound, "}")

	if startsWithBrackets != endsWithBrakets {
		return PinyinV1{}, errors.New("malformed pinyin")
	}

	if ns := strings.TrimSuffix(sound, "-"); ns != "" {
		sound = ns
	}

	isAlphabetic := soundIsAlphabetic(sound)

	hasNumber := soundHasNumber(sound)

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

	py := PinyinV1{
		Sound: sound,
		Tone:  tone,
		Type:  t,
	}

	return py, nil
}

func pinyinV2StrToPinyin(pys string) ([]PinyinV2, error) {
	words := strings.Split(pys, " ")

	v2List := make([]PinyinV2, 0, len(words))

	wordsForPyV2 := make([]PinyinV1, 0, len(words))

	runesBuilder := make([]rune, 0, len(words))
	pyItems := make([][]rune, 0, len(words))

	//Ping2guo3 shou3ji1
	for _, word := range words {
		wordsForPyV2 = wordsForPyV2[:0]

		//Ping2guo3
		runesBuilder := runesBuilder[:0]
		pyItems := pyItems[:0]
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
				runesBuilder = runesBuilder[:0]
				openBracket = false
			}
		}

		if len(runesBuilder) != 0 {
			pyItems = append(pyItems, runesBuilder)
			runesBuilder = runesBuilder[:0]
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

func makePyMap() map[string]bool {
	pym := make(map[string]bool)
	for _, v := range full_pinyin_list {
		pym[v] = true
	}
	return pym
}

func ParseLine(line string) (Ci, error) {
	pym := makePyMap()

	return parseLine(pym, line)
}

func NewLineParser() LineParser {
	pym := makePyMap()

	return basicLineParser{
		Pym: pym,
	}
}

func (blp basicLineParser) ParseLine(line string) (Ci, error) {
	return parseLine(blp.Pym, line)
}

// Traditional Simplified [[pin1yin1]] /gloss; gloss; .../gloss; gloss; .../
func parseLine(pinyinVals map[string]bool, line string) (Ci, error) {
	if strings.HasPrefix(line, "#") {
		return Ci{}, errors.New("comment line")
	}

	if strings.TrimSpace(line) == "" {
		return Ci{}, errors.New("empty line")
	}

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

	lineRunes := make([]rune, 0, len(line))

	for _, r := range line {
		lineRunes = append(lineRunes, r)
	}

	for i, r := range lineRunes {

		if currentSection == section_traditional {

			if string(r) == pinyinStart {
				return Ci{}, fmt.Errorf("found pinyin section before completing traditional section. Line: %s", line)
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
				return Ci{}, fmt.Errorf("found pinyin section before completing simplified section. Line: %s", line)
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
				return Ci{}, fmt.Errorf("found gloss section before pinyin section. Line: %s", line)
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
		return Ci{}, fmt.Errorf("malformed pinyin (cannot determine version) (%d %d). Line: %s", pyOpenBracketCount, pyCloseBracketCount, line)
	}

	switch pyOpenBracketCount {
	case 1:
		pyVersion = V1
	case 2:
		pyVersion = V2
	default:
		return Ci{}, fmt.Errorf("malformed pinyin (unrecognized version). Line: %s", line)
	}

	var py []PinyinV2
	var err error
	if pyVersion == V1 {
		py, err = pinyinV1StrToPinyin(pinyin)
	} else {
		py, err = pinyinV2StrToPinyin(pinyin)
	}

	if err != nil {
		return Ci{}, errors.Join(fmt.Errorf("error on line: %s", line), err)
	}

	if fantizi == "" || jiantizi == "" {
		return Ci{}, fmt.Errorf("no traditional/simplified found. Line: %s", line)
	}

	if pinyin == "" || len(py) == 0 {
		return Ci{}, fmt.Errorf("no pinyin found. Line: %s", line)
	}

	if len(gloss) == 0 {
		return Ci{}, fmt.Errorf("no gloss found. Line: %s", line)
	}

	for _, v := range pinyin {
		if string(norm.NFD.Bytes([]byte(string(v)))) != string(v) {
			// Really struggling to detect this
			return Ci{}, fmt.Errorf("malformed pinyin - no diacritics. Line: %s", line)
		}
	}

	for _, v := range py {
		for _, p := range v.Word {
			if p.Type == Normal && !pinyinVals[strings.ToLower(p.Sound)] {
				return Ci{}, fmt.Errorf("malformed pinyin - unrecognized pinyin value (check for ambiguity). Line: %s", line)
			}
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
