package cccedictparser

import "testing"

const benchLine = "分久必合，合久必分 分久必合，合久必分 [[fen1jiu3-bi4he2, he2jiu3-bi4fen1]] /lit. that which is long divided must unify, and that which is long unified must divide (proverb, from Romance of the Three Kingdoms 三國演義|三国演义[San1guo2 Yan3yi4])/fig. things are constantly changing/"

func BenchmarkParseLine(b *testing.B) {
	for b.Loop() {
		ParseLine(benchLine)
	}
}

func BenchmarkLineParserParseLine(b *testing.B) {
	lp := NewLineParser()
	for b.Loop() {
		lp.ParseLine(benchLine)
	}
}
