package cccedictparser

import (
	"testing"
)

type testItem struct {
	Name string
	Test func(*testing.T)
}

func TestParseLine(t *testing.T) {
	tests := []testItem{
		{Name: "parseLine_PinyinInGloss", Test: parseLine_PinyinInGloss},
		{Name: "parseLine_HandlesOneGloss", Test: parseLine_HandlesOneGloss},
		{Name: "parseLine_HandlesManyGloss", Test: parseLine_HandlesManyGloss},
		{Name: "parseLine_Error_MalformedLine", Test: parseLine_Error_MalformedLine},
		{Name: "parseLine_TraditionalMatches", Test: parseLine_TraditionalMatches},
		{Name: "parseLine_SimplifiedMatches", Test: parseLine_SimplifiedMatches},
		{Name: "parseLine_PinyinV1Matches", Test: parseLine_PinyinV1Matches},
		{Name: "parseLine_PinyinV2Matches", Test: parseLine_PinyinV2Matches},
		{Name: "parseLine_FullMatches", Test: parseLine_FullMatches},
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
	}

	if len(out.Gloss) != len(expected) {
		t.Errorf("error when parsing line (%s) -- incorrect gloss number found. Expected (1) was (%d)", line, len(out.Gloss))
	}

	for i, v := range expected {
		if out.Gloss[i] != v {
			t.Errorf("gloss did not match expected. Expected (%s) Actual (%s)", v, out.Gloss[i])
		}
	}
}

type malformedLineCase struct {
	Line string
	// Error can include more, but must include this
	ExpectedErrorMessage string
}

func parseLine_Error_MalformedLine(t *testing.T) {
	line := "浮泛 浮泛 [fu2 fan4] /to float about/(of a feeling) to show on the face/(of speech, friendship etc) shallow/vague/"
	cases := []malformedLineCase{
		{
			ExpectedErrorMessage: "no gloss found",
			Line:                 "浮泛 浮泛 [fu2 fan4] ",
		},
		{
			ExpectedErrorMessage: "no pinyin found",
			Line:                 "浮泛 浮泛 /to float about/(of a feeling) to show on the face/(of speech, friendship etc) shallow/vague/",
		},
		{
			ExpectedErrorMessage: "no traditional/simplified found",
			Line:                 "浮泛 [fu2 fan4] /to float about/(of a feeling) to show on the face/(of speech, friendship etc) shallow/vague/",
		},
		{
			ExpectedErrorMessage: "malformed pinyin v1",
			Line:                 "浮泛 浮泛 [fu2fan4] /to float about/(of a feeling) to show on the face/(of speech, friendship etc) shallow/vague/",
		},
		{
			ExpectedErrorMessage: "malformed pinyin v2 - ambiguity",
			Line:                 "e人 e人 [[eren2]] /(slang) extroverted person/",
		},
		{
			ExpectedErrorMessage: "malformed pinyin v2 - no dots",
			Line:                 "大衛·艾登堡 大卫·艾登堡 [[Da4wei4 · Ai4deng1bao3]] /David Attenborough (1926), British naturalist and broadcaster/",
		},
		{
			ExpectedErrorMessage: "malformed pinyin - no diacritics",
			Line:                 "",
		},
	}
}

func parseLine_TraditionalMatches(t *testing.T) {

}

func parseLine_SimplifiedMatches(t *testing.T) {

}

type expectedPinyinV1 struct {
	Sentence string
	Pinyin   []PinyinV1
}

type expectedPinyinV2 struct {
	Sentence string
	Pinyin   []PinyinV2
}

func parseLine_PinyinV1Matches(t *testing.T) {
	cases := []expectedPinyinV1{
		{
			Sentence: "K人 K人 [K ren2] /(slang) to hit sb; to beat sb/",
			Pinyin: []PinyinV1{
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
			Pinyin: []PinyinV1{
				{
					Sound: "xx",
					Type:  Unknown,
					Tone:  T5,
				},
			},
		},
		{
			Sentence: "大衛·艾登堡 大卫·艾登堡 [[Da4 wei4 · Ai4 deng1 bao3]] /David Attenborough (1926), British naturalist and broadcaster/",
			Pinyin: []PinyinV1{
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
			t.Errorf(err.Error())
			continue
		}

		if parsed.FormatVersion != V1 {
			t.Errorf("expected v1 pinyin for line \"%s\".", v.Sentence)
			continue
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

			if !(pyw.Word[i].Sound == v.Pinyin[i].Sound && pyw.Word[i].Tone == v.Pinyin[i].Tone && pyw.Word[i].Type == v.Pinyin[i].Type) {
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
			Pinyin: []PinyinV2{
				{
					Word: []PinyinV1{
						{
							Sound: "xx",
							Type:  Unknown,
							Tone:  T5,
						},
					},
				},
			},
		},
		{
			Sentence: "打算 打算 [[{e}ren2]] /words/",
			Pinyin: []PinyinV2{
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
			Pinyin: []PinyinV2{
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
			Pinyin: []PinyinV2{
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
			Pinyin: []PinyinV2{
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
			Pinyin: []PinyinV2{
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
			Pinyin: []PinyinV2{
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
			Pinyin: []PinyinV2{
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
			Pinyin: []PinyinV2{
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
			Pinyin: []PinyinV2{
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
			t.Errorf(err.Error())
			continue
		}

		if parsed.FormatVersion != V2 {
			t.Errorf("expected v2 pinyin for line \"%s\".", v.Sentence)
			continue
		}

		//TODO assert correct
	}
}

func parseLine_FullMatches(t *testing.T) {

}
