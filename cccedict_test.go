package cccedictparser

import "testing"

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

func parseLine_PinyinV1Matches(t *testing.T) {
	cases := []string{
		"K人 K人 [K ren2] /(slang) to hit sb; to beat sb/",
		"打算 打算 [xx5] /words/",
		"大衛·艾登堡 大卫·艾登堡 [[Da4 wei4 · Ai4 deng1 bao3]] /David Attenborough (1926), British naturalist and broadcaster/",
	}
}

func parseLine_PinyinV2Matches(t *testing.T) {
	cases := []string{
		"打算 打算 [[xx5]] /words/",
		"打算 打算 [[{e}ren2]] /words/",
		"打算 打算 [[e-ren2]] /words/",
		"打算 打算 [[yi4ren2]] /words/",
		"打算 打算 [[nu:3]] /words/",
		"打算 打算 [[zen3me5 hui2shi4 r5]] /words/",
		"K人 K人 [[K ren2]] /(slang) to hit sb; to beat sb/",
		"3Q 3Q [[san1 Q]] /thx/",
		"大衛·艾登堡 大卫·艾登堡 [[Da4wei4 Ai4deng1bao3]] /David Attenborough (1926), British naturalist and broadcaster/",
		"分久必合，合久必分 分久必合，合久必分 [[fen1jiu3-bi4he2, he2jiu3-bi4fen1]] /lit. that which is long divided must unify, and that which is long unified must divide (proverb, from Romance of the Three Kingdoms 三國演義|三国演义[San1guo2 Yan3yi4])/fig. things are constantly changing/",
	}
}

func parseLine_FullMatches(t *testing.T) {

}
