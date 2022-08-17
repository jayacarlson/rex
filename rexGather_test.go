package rex

import (
	"regexp"
	"strings"
	"testing"

	"github.com/jayacarlson/dbg"
	"github.com/jayacarlson/tst"
)

const (
	sngSQt = "\u2018"
	sngEQt = "\u2019"
	dblSQt = "\u201C"
	dblEQt = "\u201D"

	wordtext = `apple   banana   cherry
date fig grape
` + sngSQt + `apple` + sngEQt + ` ` + sngSQt + `banana` + sngEQt + ` ` + sngSQt + `cherry` + sngEQt + `
Apple Banana Cherry
` + dblSQt + `date` + dblEQt + ` ` + dblSQt + `fig` + dblEQt + ` ` + dblSQt + `grape` + dblEQt + `
Date Fig Grape
  berry  mellon  orange 
  (berry), 'mellon';  "orange".
`
)

var lcwordRex = regexp.MustCompile(`(?s)(.*?\b)?([a-z]+)(\b.*)`)

func TestRexGatherLC(t *testing.T) {
	text := ""
	expected := "apple banana cherry date fig grape apple banana cherry date fig grape berry mellon orange berry mellon orange "
	RexGather(wordtext, lcwordRex, func(x []string, rx *regexp.Regexp, gf GatherFunc) {
		text = text + x[2] + " "
		RexGather(x[3], rx, gf)
	})
	// note expected contains trailing space from gather above
	if text != expected {
		tst.Failed(t, dbg.IAm(), "Expected in green, genereted in red")
		tst.AsGreen(expected)
		tst.AsRed(text)
	} else {
		tst.Passed(t, "", dbg.IAm())
	}
}

func TestRexGatherUniqueLC(t *testing.T) {
	expected := "apple banana cherry date fig grape berry mellon orange"
	var uniqueValues []string
	exists := make(map[string]bool)

	RexGather(wordtext, lcwordRex, func(x []string, rx *regexp.Regexp, gf GatherFunc) {
		if _, ok := exists[x[2]]; !ok {
			uniqueValues = append(uniqueValues, x[2])
			exists[x[2]] = true
		}
		RexGather(x[3], rx, gf)
	})
	text := strings.Join(uniqueValues, " ")
	if text != expected {
		tst.Failed(t, dbg.IAm(), "Expected in green, genereted in red")
		tst.AsGreen(expected)
		tst.AsRed(text)
	} else {
		tst.Passed(t, "", dbg.IAm())
	}
}
