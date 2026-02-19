package cccedictparser

import (
	"strings"
	"testing"
)

type testItem struct {
	Name string
	Test func(*testing.T)
}

type testCase[T any] struct {
	Sentence string
	Expected T
}

func TestParseLine(t *testing.T) {
	tests := []testItem{
		// {Name: "parseLine_PinyinInGloss", Test: parseLine_PinyinInGloss},
		// {Name: "parseLine_HandlesOneGloss", Test: parseLine_HandlesOneGloss},
		// {Name: "parseLine_HandlesManyGloss", Test: parseLine_HandlesManyGloss},
		// {Name: "parseLine_Error_MalformedLine", Test: parseLine_Error_MalformedLine},
		// {Name: "parseLine_TraditionalMatches", Test: parseLine_TraditionalMatches},
		// {Name: "parseLine_SimplifiedMatches", Test: parseLine_SimplifiedMatches},
		// {Name: "parseLine_PinyinV1Matches", Test: parseLine_PinyinV1Matches},
		{Name: "parseLine_PinyinV2Matches", Test: parseLine_PinyinV2Matches},
		// {Name: "parseLine_FullMatches", Test: parseLine_FullMatches},
	}

	for _, v := range tests {
		t.Run(v.Name, v.Test)
	}
}

func parseLine_PinyinInGloss(t *testing.T) {
	line := "㗂 㗂 [sheng3] /variant of 省[sheng3]/tight-lipped/to examine/to watch/to scour (esp. Cantonese)/"
	expected := "variant of 省[sheng3]"

	out, err := ParseLine(line)

	if err != nil {
		t.Errorf("error when parsing line (%s). %s", line, err.Error())
	}

	if len(out.Gloss) == 0 {
		t.Errorf("error when parsing line (%s) -- no gloss found.", line)
		return
	}

	if out.Gloss[0] != expected {
		t.Errorf("gloss did not match expected. Expected (%s) Actual (%s)", expected, out.Gloss[0])
	}
}

func parseLine_HandlesOneGloss(t *testing.T) {
	line := "海嘯 海啸 [hai3 xiao4] /tsunami/"
	expected := "tsunami"

	out, err := ParseLine(line)

	if err != nil {
		t.Errorf("error when parsing line (%s). %s", line, err.Error())
	}

	if len(out.Gloss) != 1 {
		t.Errorf("error when parsing line (%s) -- incorrect gloss number found. Expected (1) was (%d)", line, len(out.Gloss))
		return
	}

	if out.Gloss[0] != expected {
		t.Errorf("gloss did not match expected. Expected (%s) Actual (%s)", expected, out.Gloss[0])
	}
}

func parseLine_HandlesManyGloss(t *testing.T) {
	line := "浮泛 浮泛 [fu2 fan4] /to float about/(of a feeling) to show on the face/(of speech, friendship etc) shallow/vague/"
	expected := []string{
		"to float about",
		"(of a feeling) to show on the face",
		"(of speech, friendship etc) shallow",
		"vague",
	}

	out, err := ParseLine(line)

	if err != nil {
		t.Errorf("error when parsing line (%s). %s", line, err.Error())
		return
	}

	if len(out.Gloss) != len(expected) {
		t.Errorf("error when parsing line (%s) -- incorrect gloss number found. Expected (1) was (%d)", line, len(out.Gloss))
		return
	}

	for i, v := range expected {
		if out.Gloss[i] != v {
			t.Errorf("gloss did not match expected. Expected (%s) Actual (%s)", v, out.Gloss[i])
		}
	}
}

type malformedLineCase = testCase[string]

func parseLine_Error_MalformedLine(t *testing.T) {
	cases := []malformedLineCase{
		{
			Expected: "no gloss found",
			Sentence: "浮泛 浮泛 [fu2 fan4] ",
		},
		{
			Expected: "no pinyin found",
			Sentence: "浮泛 浮泛 /to float about/(of a feeling) to show on the face/(of speech, friendship etc) shallow/vague/",
		},
		{
			Expected: "malformed pinyin (cannot determine version)",
			Sentence: "浮泛 浮泛 [[fu2 fan4] /to float about/(of a feeling) to show on the face/(of speech, friendship etc) shallow/vague/",
		},
		{
			Expected: "malformed pinyin (unrecognized version)",
			Sentence: "浮泛 浮泛 [[[fu2 fan4]]] /to float about/(of a feeling) to show on the face/(of speech, friendship etc) shallow/vague/",
		},
		{
			Expected: "malformed pinyin (cannot determine version)",
			Sentence: "浮泛 浮泛 [fu2 fan4]] /to float about/(of a feeling) to show on the face/(of speech, friendship etc) shallow/vague/",
		},
		{
			Expected: "no traditional/simplified found",
			Sentence: "浮泛 [fu2 fan4] /to float about/(of a feeling) to show on the face/(of speech, friendship etc) shallow/vague/",
		},
		{
			Expected: "malformed pinyin v1",
			Sentence: "浮泛 浮泛 [fu2fan4] /to float about/(of a feeling) to show on the face/(of speech, friendship etc) shallow/vague/",
		},
		{
			Expected: "malformed pinyin v2 - ambiguity",
			Sentence: "e人 e人 [[eren2]] /(slang) extroverted person/",
		},
		{
			Expected: "malformed pinyin v2 - no dots",
			Sentence: "大衛·艾登堡 大卫·艾登堡 [[Da4wei4 · Ai4deng1bao3]] /David Attenborough (1926), British naturalist and broadcaster/",
		},
		{
			Expected: "malformed pinyin - no diacritics",
			Sentence: "",
		},
	}

	for _, v := range cases {
		_, err := ParseLine(v.Sentence)

		if err == nil {
			t.Errorf("expected error for line \"%s\".", v.Sentence)
			continue
		}

		if err.Error() != v.Expected {
			t.Errorf("expected error (%s), actual error (%s)", v.Expected, err.Error())
			continue
		}
	}
}

func parseLine_TraditionalMatches(t *testing.T) {
	cases := []testCase[string]{
		{
			Sentence: "禁酒 禁酒 [jin4 jiu3] /prohibition/ban on alcohol/dry law/",
			Expected: "禁酒",
		},
		{
			Sentence: "航海年表 航海年表 [hang2 hai3 nian2 biao3] /nautical ephemeris/",
			Expected: "航海年表",
		},
		{
			Sentence: "頭髮 头发 [tou2 fa5] /hair (on the head)/",
			Expected: "頭髮",
		},
		{
			Sentence: "顆粒歸倉 颗粒归仓 [ke1 li4 gui1 cang1] /to gather all the harvested grain into the granary; to harvest every single grain/",
			Expected: "顆粒歸倉",
		},
		{
			Sentence: "援軍 援军 [yuan2 jun1] /(military) reinforcements/",
			Expected: "援軍",
		},
	}

	for _, v := range cases {
		parsed, err := ParseLine(v.Sentence)

		if err != nil {
			t.Errorf("error: %s. Line %s", err.Error(), v.Sentence)
			continue
		}

		if parsed.Fantizi != v.Expected {
			t.Errorf("expected %s, actual %s", parsed.Fantizi, v.Expected)
			continue
		}
	}
}

func parseLine_SimplifiedMatches(t *testing.T) {
	cases := []testCase[string]{
		{
			Sentence: "禁酒 禁酒 [jin4 jiu3] /prohibition/ban on alcohol/dry law/",
			Expected: "禁酒",
		},
		{
			Sentence: "航海年表 航海年表 [hang2 hai3 nian2 biao3] /nautical ephemeris/",
			Expected: "航海年表",
		},
		{
			Sentence: "頭髮 头发 [tou2 fa5] /hair (on the head)/",
			Expected: "头发",
		},
		{
			Sentence: "顆粒歸倉 颗粒归仓 [ke1 li4 gui1 cang1] /to gather all the harvested grain into the granary; to harvest every single grain/",
			Expected: "颗粒归仓",
		},
		{
			Sentence: "援軍 援军 [yuan2 jun1] /(military) reinforcements/",
			Expected: "援军",
		},
	}

	for _, v := range cases {
		parsed, err := ParseLine(v.Sentence)

		if err != nil {
			t.Errorf("error: %s. Line %s", err.Error(), v.Sentence)
			continue
		}

		if parsed.Jiantizi != v.Expected {
			t.Errorf("expected %s, actual %s", parsed.Jiantizi, v.Expected)
			continue
		}
	}
}

type expectedPinyinV1 = testCase[[]PinyinV1]

type expectedPinyinV2 = testCase[[]PinyinV2]

func parseLine_PinyinV1Matches(t *testing.T) {
	cases := []expectedPinyinV1{
		{
			Sentence: "K人 K人 [K ren2] /(slang) to hit sb; to beat sb/",
			Expected: []PinyinV1{
				{
					Sound: "K",
					Type:  Alphabet,
					Tone:  None,
				},
				{
					Sound: "ren",
					Type:  Normal,
					Tone:  T2,
				},
			},
		},
		{
			Sentence: "打算 打算 [xx5] /words/",
			Expected: []PinyinV1{
				{
					Sound: "xx",
					Type:  Unknown,
					Tone:  T5,
				},
			},
		},
		{
			Sentence: "大衛·艾登堡 大卫·艾登堡 [[Da4 wei4 · Ai4 deng1 bao3]] /David Attenborough (1926), British naturalist and broadcaster/",
			Expected: []PinyinV1{
				{
					Sound: "Da",
					Type:  Normal,
					Tone:  T4,
				},
				{
					Sound: "wei",
					Type:  Normal,
					Tone:  T4,
				},
				{
					Sound: "·",
					Type:  Special,
					Tone:  None,
				},
				{
					Sound: "Ai",
					Type:  Normal,
					Tone:  T4,
				},
				{
					Sound: "deng",
					Type:  Normal,
					Tone:  T1,
				},
				{
					Sound: "bao",
					Type:  Normal,
					Tone:  T3,
				},
			},
		},
	}

	for _, v := range cases {
		parsed, err := ParseLine(v.Sentence)

		if err != nil {
			t.Errorf("error: %s. Line %s", err.Error(), v.Sentence)
			continue
		}

		if parsed.FormatVersion != V1 {
			t.Errorf("expected v1 pinyin for line \"%s\".", v.Sentence)
			continue
		}

		if len(parsed.Pinyin) != len(v.Expected) {
			t.Errorf("length mismatch for expected vs actual. expected %d, actual %d", len(v.Expected), len(parsed.Pinyin))
		}

		skip := false
		for i := 0; i < len(parsed.Pinyin); i++ {
			if skip {
				break
			}

			pyw := parsed.Pinyin[i]
			if len(pyw.Word) != 1 {
				t.Errorf("expected len of 1 for v1 pinyin on line \"%s\".", v.Sentence)
				skip = true
				continue
			}

			if !(pyw.Word[0].Sound == v.Expected[i].Sound && pyw.Word[0].Tone == v.Expected[i].Tone && pyw.Word[0].Type == v.Expected[i].Type) {
				t.Errorf("failed when checking v1 pinyin. Line \"%s\".", v.Sentence)
				skip = true
				continue
			}
		}
	}
}

func parseLine_PinyinV2Matches(t *testing.T) {
	cases := []expectedPinyinV2{
		{
			Sentence: "打算 打算 [[xx5]] /words/",
			Expected: []PinyinV2{
				{
					Word: []PinyinV1{
						{
							Sound: "xx",
							Type:  Unknown,
							Tone:  None,
						},
					},
				},
			},
		},
		{
			Sentence: "打算 打算 [[{e}ren2]] /words/",
			Expected: []PinyinV2{
				{
					Word: []PinyinV1{
						{
							Sound: "e",
							Type:  Alphabet,
							Tone:  None,
						},
						{
							Sound: "ren",
							Type:  Normal,
							Tone:  T2,
						},
					},
				},
			},
		},
		{
			Sentence: "打算 打算 [[e-ren2]] /words/",
			Expected: []PinyinV2{
				{
					Word: []PinyinV1{
						{
							Sound: "e",
							Type:  Alphabet,
							Tone:  None,
						},
						{
							Sound: "ren",
							Type:  Normal,
							Tone:  T2,
						},
					},
				},
			},
		},
		{
			Sentence: "打算 打算 [[yi4ren2]] /words/",
			Expected: []PinyinV2{
				{
					Word: []PinyinV1{
						{
							Sound: "yi",
							Type:  Normal,
							Tone:  T4,
						},
						{
							Sound: "ren",
							Type:  Normal,
							Tone:  T2,
						},
					},
				},
			},
		},
		{
			Sentence: "打算 打算 [[nu:3]] /words/",
			Expected: []PinyinV2{
				{
					Word: []PinyinV1{
						{
							Sound: "nv",
							Type:  Normal,
							Tone:  T3,
						},
					},
				},
			},
		},
		{
			Sentence: "打算 打算 [[zen3me5 hui2shi4 r5]] /words/",
			Expected: []PinyinV2{
				{
					Word: []PinyinV1{
						{
							Sound: "zen",
							Type:  Normal,
							Tone:  T3,
						},
						{
							Sound: "me",
							Type:  Normal,
							Tone:  T5,
						},
					},
				},
				{
					Word: []PinyinV1{
						{
							Sound: "hui",
							Type:  Normal,
							Tone:  T2,
						},
						{
							Sound: "shi",
							Type:  Normal,
							Tone:  T4,
						},
					},
				},
				{
					Word: []PinyinV1{
						{
							Sound: "r",
							Type:  Normal,
							Tone:  T5,
						},
					},
				},
			},
		},
		{
			Sentence: "K人 K人 [[K ren2]] /(slang) to hit sb; to beat sb/",
			Expected: []PinyinV2{
				{
					Word: []PinyinV1{
						{
							Sound: "K",
							Type:  Alphabet,
							Tone:  None,
						},
					},
				},
				{
					Word: []PinyinV1{
						{
							Sound: "ren",
							Type:  Normal,
							Tone:  T2,
						},
					},
				},
			},
		},
		{
			Sentence: "3Q 3Q [[san1 Q]] /thx/",
			Expected: []PinyinV2{
				{
					Word: []PinyinV1{
						{
							Sound: "san",
							Type:  Normal,
							Tone:  T1,
						},
					},
				},
				{
					Word: []PinyinV1{
						{
							Sound: "Q",
							Type:  Alphabet,
							Tone:  None,
						},
					},
				},
			},
		},
		{
			Sentence: "大衛·艾登堡 大卫·艾登堡 [[Da4wei4 Ai4deng1bao3]] /David Attenborough (1926), British naturalist and broadcaster/",
			Expected: []PinyinV2{
				{
					Word: []PinyinV1{
						{
							Sound: "Da",
							Type:  Normal,
							Tone:  T4,
						},
						{
							Sound: "wei",
							Type:  Normal,
							Tone:  T4,
						},
					},
				},
				{
					Word: []PinyinV1{
						{
							Sound: "Ai",
							Type:  Normal,
							Tone:  T4,
						},
						{
							Sound: "deng",
							Type:  Normal,
							Tone:  T1,
						},
						{
							Sound: "bao",
							Type:  Normal,
							Tone:  T3,
						},
					},
				},
			},
		},
		{
			Sentence: "分久必合，合久必分 分久必合，合久必分 [[fen1jiu3-bi4he2, he2jiu3-bi4fen1]] /lit. that which is long divided must unify, and that which is long unified must divide (proverb, from Romance of the Three Kingdoms 三國演義|三国演义[San1guo2 Yan3yi4])/fig. things are constantly changing/",
			Expected: []PinyinV2{
				{
					Word: []PinyinV1{
						{
							Sound: "fen",
							Type:  Normal,
							Tone:  T1,
						},
						{
							Sound: "jiu",
							Type:  Normal,
							Tone:  T3,
						},
						{
							Sound: "bi",
							Type:  Normal,
							Tone:  T4,
						},
						{
							Sound: "he",
							Type:  Normal,
							Tone:  T2,
						},
					},
				},
				{
					Word: []PinyinV1{
						{
							Sound: "he",
							Type:  Normal,
							Tone:  T2,
						},
						{
							Sound: "jiu",
							Type:  Normal,
							Tone:  T3,
						},
						{
							Sound: "bi",
							Type:  Normal,
							Tone:  T4,
						},
						{
							Sound: "fen",
							Type:  Normal,
							Tone:  T1,
						},
					},
				},
			},
		},
	}
	for _, v := range cases {
		parsed, err := ParseLine(v.Sentence)

		if err != nil {
			t.Errorf("error: %s. Line %s", err.Error(), v.Sentence)
			continue
		}

		if parsed.FormatVersion != V2 {
			t.Errorf("expected v2 pinyin for line \"%s\".", v.Sentence)
			continue
		}

		if len(v.Expected) != len(parsed.Pinyin) {
			t.Errorf("expected pinyin length to be %d, was %d. Ex(%s) Ac(%s), Line: %s", len(v.Expected), len(parsed.Pinyin), pyV2ArrStr(v.Expected), pyV2ArrStr(parsed.Pinyin), v.Sentence)
			continue
		}

		for i := 0; i < len(parsed.Pinyin); i++ {
			if len(v.Expected[i].Word) != len(parsed.Pinyin[i].Word) {
				t.Errorf("expected pinyin word %d length to be %d, was %d. Ex(%s) Ac(%s), Line: %s", i, len(v.Expected[i].Word), len(parsed.Pinyin[i].Word), v.Expected[i].String(), parsed.Pinyin[i].String(), v.Sentence)
				continue
			}

			for j := 0; j < len(parsed.Pinyin[i].Word); j++ {
				if !pyEq(parsed.Pinyin[i].Word[j], v.Expected[i].Word[j]) {
					t.Errorf("expected %s, actual %s", v.Expected[i].Word[j], parsed.Pinyin[i].Word[j])
					i = len(parsed.Pinyin)
					break
				}
			}
		}
	}
}

func pyV2ArrStr(pyv2arr []PinyinV2) string {
	items := []string{}
	for _, v := range pyv2arr {
		items = append(items, v.String())
	}
	return strings.Join(items, ", ")
}

func pyEq(a PinyinV1, b PinyinV1) bool {
	return a.Sound == b.Sound && a.Tone == b.Tone && a.Type == b.Type
}

func parseLine_FullMatches(t *testing.T) {
	cases := []testCase[Ci]{
		{
			Sentence: "損人不利己 损人不利己 [sun3 ren2 bu4 li4 ji3] /to harm others without benefiting oneself (idiom)/",
			Expected: Ci{
				Fantizi:  "損人不利己",
				Jiantizi: "损人不利己",
				Pinyin: []PinyinV2{
					{
						Word: []PinyinV1{
							{
								Sound: "sun",
								Tone:  T3,
								Type:  Normal,
							},
						},
					},
					{
						Word: []PinyinV1{
							{
								Sound: "ren",
								Tone:  T2,
								Type:  Normal,
							},
						},
					},
					{
						Word: []PinyinV1{
							{
								Sound: "bu",
								Tone:  T4,
								Type:  Normal,
							},
						},
					},
					{
						Word: []PinyinV1{
							{
								Sound: "li",
								Tone:  T4,
								Type:  Normal,
							},
						},
					},
					{
						Word: []PinyinV1{
							{
								Sound: "ji",
								Tone:  T3,
								Type:  Normal,
							},
						},
					},
				},
				PinyinRaw:     "sun3 ren2 bu4 li4 ji3",
				Gloss:         []string{"to harm others without benefiting oneself (idiom)"},
				FormatVersion: V1,
			},
		},
		{
			Sentence: "眼觀四面，耳聽八方 眼观四面，耳听八方 [yan3 guan1 si4 mian4 , er3 ting1 ba1 fang1] /lit. the eyes observe all sides and the ears listen in all directions (idiom)/fig. to be observant and alert/",
			Expected: Ci{
				Fantizi:  "眼觀四面，耳聽八方",
				Jiantizi: "眼观四面，耳听八方",
				Pinyin: []PinyinV2{
					{
						Word: []PinyinV1{
							{
								Sound: "yan",
								Tone:  T3,
								Type:  Normal,
							},
						},
					},
					{
						Word: []PinyinV1{
							{
								Sound: "guan",
								Tone:  T1,
								Type:  Normal,
							},
						},
					},
					{
						Word: []PinyinV1{
							{
								Sound: "si",
								Tone:  T4,
								Type:  Normal,
							},
						},
					},
					{
						Word: []PinyinV1{
							{
								Sound: "mian",
								Tone:  T4,
								Type:  Normal,
							},
						},
					},
					{
						Word: []PinyinV1{
							{
								Sound: ",",
								Tone:  None,
								Type:  Special,
							},
						},
					},
					{
						Word: []PinyinV1{
							{
								Sound: "er",
								Tone:  T3,
								Type:  Normal,
							},
						},
					},
					{
						Word: []PinyinV1{
							{
								Sound: "ting",
								Tone:  T1,
								Type:  Normal,
							},
						},
					},
					{
						Word: []PinyinV1{
							{
								Sound: "ba",
								Tone:  T1,
								Type:  Normal,
							},
						},
					},
					{
						Word: []PinyinV1{
							{
								Sound: "fang",
								Tone:  T1,
								Type:  Normal,
							},
						},
					},
				},
				PinyinRaw:     "yan3 guan1 si4 mian4 , er3 ting1 ba1 fang1",
				Gloss:         []string{"lit. the eyes observe all sides and the ears listen in all directions (idiom)", "fig. to be observant and alert"},
				FormatVersion: V1,
			},
		},
		{
			Sentence: "薄命 薄命 [bo2 ming4] /to be born under an unlucky star (usu. of women)/to be born unlucky/",
			Expected: Ci{
				Fantizi:  "薄命",
				Jiantizi: "薄命",
				Pinyin: []PinyinV2{
					{
						Word: []PinyinV1{
							{
								Sound: "bo",
								Tone:  T2,
								Type:  Normal,
							},
						},
					},
					{
						Word: []PinyinV1{
							{
								Sound: "ming",
								Tone:  T4,
								Type:  Normal,
							},
						},
					},
				},
				PinyinRaw:     "bo2 ming4",
				Gloss:         []string{"to be born under an unlucky star (usu. of women)", "to be born unlucky"},
				FormatVersion: V1,
			},
		},
		{
			Sentence: "皮實 皮实 [[pi2shi5]] /(of things) durable/(of people) sturdy; tough/",
			Expected: Ci{
				Fantizi:  "皮實",
				Jiantizi: "皮实",
				Pinyin: []PinyinV2{
					{
						Word: []PinyinV1{
							{
								Sound: "pi",
								Tone:  T2,
								Type:  Normal,
							},
							{
								Sound: "shi",
								Tone:  T5,
								Type:  Normal,
							},
						},
					},
				},
				PinyinRaw:     "pi2shi5",
				Gloss:         []string{"(of things) durable", "(of people) sturdy; tough"},
				FormatVersion: V2,
			},
		},
	}

	for _, v := range cases {
		parsed, err := ParseLine(v.Sentence)

		if err != nil {
			t.Errorf("error: %s. Line %s", err.Error(), v.Sentence)
			continue
		}

		if !ciEq(v.Expected, parsed) {
			t.Errorf("expected %s, got %s", v.Expected, parsed)
			continue
		}
	}
}

func ciEq(a Ci, b Ci) bool {
	if a.Fantizi != b.Fantizi {
		return false
	}
	if a.Jiantizi != b.Jiantizi {
		return false
	}
	if a.PinyinRaw != b.PinyinRaw {
		return false
	}
	if a.FormatVersion != b.FormatVersion {
		return false
	}
	if len(a.Gloss) != len(b.Gloss) {
		return false
	}
	if len(a.Pinyin) != len(b.Pinyin) {
		return false
	}
	for i := 0; i < len(a.Pinyin); i++ {
		pyA := a.Pinyin[i]
		pyB := b.Pinyin[i]
		if len(pyA.Word) != len(pyB.Word) {
			return false
		}
		for j := 0; j < len(pyA.Word); j++ {
			pyAWJ := pyA.Word[j]
			pyBWJ := pyB.Word[j]

			if !pyEq(pyAWJ, pyBWJ) {
				return false
			}
		}
	}
	return true
}
