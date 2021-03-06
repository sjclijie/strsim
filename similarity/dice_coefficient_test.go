package similarity

import (
	"fmt"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func Test_DiceCoefficient_CompareAscii(t *testing.T) {
	// ngram = 1
	d := &DiceCoefficient{Ngram: 1}

	for k, v := range []testOneCase{
		{s1: "ivan1", s2: "ivan2", cost: 0.8},
		{s2: "ivan1", s1: "ivan2", cost: 0.8},

		{s1: "love", s2: "love", cost: 1},
	} {
		m := fmt.Sprintf("error case:%d", k)
		assert.Equal(t, d.CompareAscii(v.s1, v.s2), v.cost, m)
		assert.Equal(t, d.l1, len(v.s1), m)
		assert.Equal(t, d.l2, len(v.s2), m)
	}

}

func Test_DiceCoefficient_CompareAscii_NgramOrMore(t *testing.T) {
	// ngram = 2
	d := &DiceCoefficient{Ngram: 2, test: true}
	for k, v := range []testOneCase{
		{s1: "John Smith", s2: "Smith, John D.", cost: 0.7272727272727273, ngram: 2},
		{s2: "John Smith", s1: "Smith, John D.", cost: 0.7272727272727273, ngram: 2},

		{s1: "John Smith", s2: "Smith, John D.", cost: 0.6, ngram: 3},
		{s2: "John Smith", s1: "Smith, John D.", cost: 0.6, ngram: 3},

		{s1: "John Smith", s2: "Smith, John D.", cost: 0.4444444444444444, ngram: 4},
		{s2: "John Smith", s1: "Smith, John D.", cost: 0.4444444444444444, ngram: 4},
	} {
		if v.ngram != 0 {
			d.Ngram = v.ngram
		}

		m := fmt.Sprintf("error case:%d", k)
		assert.Equal(t, d.CompareAscii(v.s1, v.s2), v.cost, m)
		for _, v := range d.key {
			assert.Equal(t, utf8.RuneCountInString(v), d.Ngram, fmt.Sprintf("key is (%s)", v))
		}

		d.key = nil

	}
}

func Test_DiceCoefficient_CompareUtf8(t *testing.T) {
	d := &DiceCoefficient{Ngram: 1}

	for k, v := range []testOneCase{
		{s1: "你好中国", s2: "你好中国", cost: 1},
		{s1: "中文也被称为华文、汉文。中文（汉语）有标准语和方言之分，其标准语即汉语普通话", s2: "方块", cost: 0.05},
		{s1: "加油，来个", s2: "加油，来吧", cost: 0.8},
	} {
		assert.Equal(t, d.CompareUtf8(v.s1, v.s2), v.cost, fmt.Sprintf("error case:%d", k))
		//fmt.Printf("mixed:%d, l1:%d, l2:%d, l1:%d\n", d.mixed, d.l1, d.l2, utf8.RuneCountInString(v.s1))
	}
}

func Test_DiceCoefficient_FindBestMatch(t *testing.T) {
	d := &DiceCoefficient{Ngram: 1}

	for k, v := range []testBestCase{
		{s: "白日依山尽", targets: []string{"白日依山尽", "黄河入海流", "欲穷千里目", "更上一层楼"}, bestIndex: 0},
		{s: "黄河流", targets: []string{"白日依山尽", "黄河入海流", "欲穷千里目", "更上一层楼"}, bestIndex: 1},
		{s: "一层", targets: []string{"白日依山尽", "黄河入海流", "欲穷千里目", "更上一层楼"}, bestIndex: 3},
		{s: "楼", targets: []string{"白日依山尽", "黄河入海流", "欲穷千里目", "更上一层楼"}, bestIndex: 3},
		{s: "山近", targets: []string{"白日依山尽", "黄河入海流", "欲穷千里目", "更上一层楼"}, bestIndex: 0},
		{s: "海刘", targets: []string{"白日依山尽", "黄河入海流", "欲穷千里目", "更上一层楼"}, bestIndex: 1},
	} {
		mr := findBestMatch(v.s, v.targets, d.CompareUtf8)
		assert.Equal(t, mr.BestIndex, v.bestIndex, fmt.Sprintf("error case:%d", k))
	}
}
