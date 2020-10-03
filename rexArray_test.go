package rex

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"testing"
)

type arrays struct {
	Numbers0 []int         `json:"n0"`
	Numbers1 []int         `json:"n1"`
	Numbers2 [][]int       `json:"n2"`
	Numbers3 [][][]int     `json:"n3"`
	Numbers4 [][][][]int   `json:"n4"`
	Numbers5 [][][][][]int `json:"n5"`
}

var arrayText string

func TestRexJSONCleanupArraysOfArrays(t *testing.T) {
	numbers := arrays{}
	numbers.Numbers0 = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	numbers.Numbers1 = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	numbers.Numbers2 = [][]int{
		[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}
	numbers.Numbers3 = [][][]int{
		[][]int{
			[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		[][]int{
			[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		[][]int{
			[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
	}
	numbers.Numbers4 = [][][][]int{
		[][][]int{
			[][]int{
				[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			},
			[][]int{
				[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			},
		},
		[][][]int{
			[][]int{
				[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			},
			[][]int{
				[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			},
		},
		[][][]int{
			[][]int{
				[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			},
			[][]int{
				[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			},
		},
	}
	numbers.Numbers5 = [][][][][]int{
		[][][][]int{
			[][][]int{
				[][]int{
					[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
					[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				},
				[][]int{
					[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
					[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				},
			},
			[][][]int{
				[][]int{
					[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
					[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				},
				[][]int{
					[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
					[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				},
			},
		},
		[][][][]int{
			[][][]int{
				[][]int{
					[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
					[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				},
				[][]int{
					[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
					[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				},
			},
			[][][]int{
				[][]int{
					[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
					[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				},
				[][]int{
					[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
					[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				},
			},
		},
	}
	blue("Running: " + iAm())
	b, _ := json.MarshalIndent(numbers, "", "  ")
	arrayText = string(b)

	// arrayText ends up being over 500 lines long, lets not have all that crap here
	expected := "d17b96bf21f37133432b7091664b1967"
	md5Sum := fmt.Sprintf("%x", md5.Sum(b))
	if md5Sum != expected {
		yellow(arrayText)
		green(expected)
		red(md5Sum)
		t.Fail()
	}
}

func TestRexJSONCleanupArrays_PackEverything(t *testing.T) {
	expected := `{
  "n0": [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ],
  "n1": [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ],
  "n2": [
    [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ],
    [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ]
  ],
  "n3": [ [
      [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ],
      [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ]
    ], [
      [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ],
      [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ]
    ], [
      [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ],
      [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ]
  ] ],
  "n4": [ [ [
        [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ],
        [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ]
      ], [
        [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ],
        [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ]
    ] ], [ [
        [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ],
        [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ]
      ], [
        [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ],
        [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ]
    ] ], [ [
        [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ],
        [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ]
      ], [
        [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ],
        [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ]
  ] ] ],
  "n5": [ [ [ [
          [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ],
          [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ]
        ], [
          [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ],
          [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ]
      ] ], [ [
          [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ],
          [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ]
        ], [
          [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ],
          [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ]
    ] ] ], [ [ [
          [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ],
          [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ]
        ], [
          [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ],
          [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ]
      ] ], [ [
          [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ],
          [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ]
        ], [
          [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ],
          [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ]
  ] ] ] ]
}`
	blue("Running: " + iAm())
	// RexJSONCleanup(arrayText, UnnamedJSONObjectRex, packEverything) could be used, but packEverything does the same,
	// but without the removing / adding of lead spacing -- but also then adds a leading '\n' which needs removal
	text := packEverything(arrayText)
	text = removeBlackLines(text[1:])
	text = RexJSONCleanup(text, UnnamedJSONObjectRex, concatArrays)
	if text != expected {
		green(expected)
		red(text)
		t.Fail()
	}
}
