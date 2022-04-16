package rex

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

const (
	sngSQt = "\u2018"
	sngEQt = "\u2019"
	dblSQt = "\u201C"
	dblEQt = "\u201D"

	wordList = `apple   banana   cherry
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
	list := ""
	RexGather(wordList, lcwordRex, func(x []string, rx *regexp.Regexp, gf GatherFunc) {
		list = list + x[2] + " "
		RexGather(x[3], rx, gf)
	})
	// note trailing space from gather above
	if list != "apple banana cherry date fig grape apple banana cherry date fig grape berry mellon orange berry mellon orange " {
		fmt.Println("Gather failed:", list)
		t.Fail()
	}
}

func TestRexGatherUniqueLC(t *testing.T) {
	var uniqueValues []string
	exists := make(map[string]bool)

	RexGather(wordList, lcwordRex, func(x []string, rx *regexp.Regexp, gf GatherFunc) {
		if _, ok := exists[x[2]]; !ok {
			uniqueValues = append(uniqueValues, x[2])
			exists[x[2]] = true
		}
		RexGather(x[3], rx, gf)
	})
	list := strings.Join(uniqueValues, " ")
	if list != "apple banana cherry date fig grape berry mellon orange" {
		fmt.Println("Gather unique failed:", list)
		t.Fail()
	}
}
