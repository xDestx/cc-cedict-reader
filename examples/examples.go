package examples

import cccedictparser "github.com/xDestx/cc-cedict-reader"

func ex_read_line() {
	ci, err := cccedictparser.ParseLine("同床異夢 同床异梦 [tong2 chuang2 yi4 meng4] /lit. to share the same bed with different dreams (idiom); ostensible partners with different agendas/strange bedfellows/marital dissension/")

	if err != nil {
		println(err.Error())
	}

	// term
	trad := ci.Fantizi
	simp := ci.Jiantizi
	// structured pinyin
	pinyin := ci.Pinyin
	// raw string pinyin
	pyStr := ci.PinyinRaw
	// gloss
	meaning := ci.Gloss
	// cc cedict version detected
	ver := ci.FormatVersion

	//print as string
	println(ci.String())
	println(trad)
	println(simp)
	for _, v := range pinyin {
		println(v.String())
	}
	println(pyStr)
	for _, v := range meaning {
		println(v)
	}
	println(ver)
}
