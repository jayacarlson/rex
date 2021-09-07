package rex

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"testing"
)

type (
	pt struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
		Z float64 `json:"z"`
	}
	quat struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
		Z float64 `json:"z"`
		W float64 `json:"w"`
	}
	color struct {
		R int `json:"r"`
		G int `json:"g"`
		B int `json:"b"`
	}
	vert struct {
		P pt    `json:"p"`
		C color `json:"c"`
		O quat  `json:"o"`
	}
	subvert struct {
		S [3]vert `json:"verts"`
	}

	simpleObject struct {
		Name        string `json:"name"`
		Location    pt     `json:"location"`
		Orientation quat   `json:"orientation"`
		Color       color  `json:"color"`
		Numbers     []int  `json:"numbers"`
	}

	objects struct {
		Verticies []vert  `json:"verts"`
		SubVerts  subvert `json:"subverts"`
	}
)

func TestRexJSONCleanupNamedPack(t *testing.T) {
	// Only process the 'location' and 'orientation' objects, no other "named" objects
	limitedObjRex := regexp.MustCompile(`((?s).*? +"(?:location|orientation)": {)\n((?sm).+?^) +((?s)}.*)`)
	obj := simpleObject{}
	obj.Name = "Test"
	obj.Numbers = []int{1, 2, 3}

	expected := `{
  "name": "Test",
  "location": { "x": 0, "y": 0, "z": 0 },
  "orientation": { "x": 0, "y": 0, "z": 0, "w": 0 },
  "color": {
    "r": 0,
    "g": 0,
    "b": 0
  },
  "numbers": [
    1,
    2,
    3
  ]
}`
	blue("Running: " + iAm())
	b, _ := json.MarshalIndent(obj, "", "  ")
	text := RexJSONCleanup(string(b), limitedObjRex, PackLines)
	if text != expected {
		green(expected)
		red(text)
		t.Fail()
	}
}

var (
	testText string // json.Marshal generated test object

	// Extract the "verts" array at depth for cleanup, returning a list of unnamed objects
	vertsArrayRex = regexp.MustCompile(`((?sm).*?^ +"verts": \[\n)((?s).*?)((?sm)^ +\].*)`)
	// Do chaining of object following object:  }, {
	objectChainRex = regexp.MustCompile(`((?sm).*?^ *},)\n *((?s){.*)`)
	// Pack lead object in array of objects: [ {
	arrayOfObjRex = regexp.MustCompile(`((?sm).*?"\w+": \[)\n +({)((?s).*)`)
)

func TestVerifyJSONMarshalling(t *testing.T) {
	testObject := objects{}
	testObject.Verticies = make([]vert, 3)

	blue("Running: " + iAm())
	b, _ := json.MarshalIndent(testObject, "", "  ")
	testText = string(b)

	// testText ends up being over 100 lines long, lets not have all that crap here
	expected := "1c90047f4653f48dd6204c3a9e059e83"
	md5Sum := fmt.Sprintf("%x", md5.Sum(b))
	if md5Sum != expected {
		yellow(testText)
		green(expected)
		red(md5Sum)
		t.Fail()
	}
}

func TestRexJSONSimpleArrayOfObjects(t *testing.T) {
	expected := `{
  "verts": [
    {
      "p": { "x": 0, "y": 0, "z": 0 },
      "c": { "r": 0, "g": 0, "b": 0 },
      "o": { "x": 0, "y": 0, "z": 0, "w": 0 }
    },
    {
      "p": { "x": 0, "y": 0, "z": 0 },
      "c": { "r": 0, "g": 0, "b": 0 },
      "o": { "x": 0, "y": 0, "z": 0, "w": 0 }
    },
    {
      "p": { "x": 0, "y": 0, "z": 0 },
      "c": { "r": 0, "g": 0, "b": 0 },
      "o": { "x": 0, "y": 0, "z": 0, "w": 0 }
    }
  ],
  "subverts": {
    "verts": [
      {
        "p": { "x": 0, "y": 0, "z": 0 },
        "c": { "r": 0, "g": 0, "b": 0 },
        "o": { "x": 0, "y": 0, "z": 0, "w": 0 }
      },
      {
        "p": { "x": 0, "y": 0, "z": 0 },
        "c": { "r": 0, "g": 0, "b": 0 },
        "o": { "x": 0, "y": 0, "z": 0, "w": 0 }
      },
      {
        "p": { "x": 0, "y": 0, "z": 0 },
        "c": { "r": 0, "g": 0, "b": 0 },
        "o": { "x": 0, "y": 0, "z": 0, "w": 0 }
      }
    ]
  }
}`
	blue("Running: " + iAm())
	cleanVerts := func(src string) string {
		return RexJSONCleanup(src, UnnamedJSONObjectRex, func(s string) string {
			return "\n" + RexJSONCleanup(s, NamedJSONObjectRex, PackLines)
		})
	}

	text := RexJSONCleanup(testText, vertsArrayRex, cleanVerts)
	if text != expected {
		green(expected)
		red(text)
		t.Fail()
	}
}

func TestRexJSONChainedObjects(t *testing.T) {
	expected := `{
  "verts": [
    {
      "p": { "x": 0, "y": 0, "z": 0 },
      "c": { "r": 0, "g": 0, "b": 0 },
      "o": { "x": 0, "y": 0, "z": 0, "w": 0 }
    }, {
      "p": { "x": 0, "y": 0, "z": 0 },
      "c": { "r": 0, "g": 0, "b": 0 },
      "o": { "x": 0, "y": 0, "z": 0, "w": 0 }
    }, {
      "p": { "x": 0, "y": 0, "z": 0 },
      "c": { "r": 0, "g": 0, "b": 0 },
      "o": { "x": 0, "y": 0, "z": 0, "w": 0 }
    }
  ],
  "subverts": {
    "verts": [
      {
        "p": { "x": 0, "y": 0, "z": 0 },
        "c": { "r": 0, "g": 0, "b": 0 },
        "o": { "x": 0, "y": 0, "z": 0, "w": 0 }
      }, {
        "p": { "x": 0, "y": 0, "z": 0 },
        "c": { "r": 0, "g": 0, "b": 0 },
        "o": { "x": 0, "y": 0, "z": 0, "w": 0 }
      }, {
        "p": { "x": 0, "y": 0, "z": 0 },
        "c": { "r": 0, "g": 0, "b": 0 },
        "o": { "x": 0, "y": 0, "z": 0, "w": 0 }
      }
    ]
  }
}`
	blue("Running: " + iAm())
	cleanVerts := func(src string) string {
		r := RexJSONCleanup(src, UnnamedJSONObjectRex, func(s string) string {
			return "\n" + RexJSONCleanup(s, NamedJSONObjectRex, PackLines)
		})
		return RexReplace(r, objectChainRex, func(x []string, rx *regexp.Regexp, rf RexFunc) string {
			return x[1] + " " + RexReplace(x[2], rx, rf)
		})
	}

	text := RexJSONCleanup(testText, vertsArrayRex, cleanVerts)
	if text != expected {
		green(expected)
		red(text)
		t.Fail()
		return
	}

	expected = `{
  "verts": [ {
      "p": { "x": 0, "y": 0, "z": 0 },
      "c": { "r": 0, "g": 0, "b": 0 },
      "o": { "x": 0, "y": 0, "z": 0, "w": 0 }
    }, {
      "p": { "x": 0, "y": 0, "z": 0 },
      "c": { "r": 0, "g": 0, "b": 0 },
      "o": { "x": 0, "y": 0, "z": 0, "w": 0 }
    }, {
      "p": { "x": 0, "y": 0, "z": 0 },
      "c": { "r": 0, "g": 0, "b": 0 },
      "o": { "x": 0, "y": 0, "z": 0, "w": 0 }
    }
  ],
  "subverts": {
    "verts": [ {
        "p": { "x": 0, "y": 0, "z": 0 },
        "c": { "r": 0, "g": 0, "b": 0 },
        "o": { "x": 0, "y": 0, "z": 0, "w": 0 }
      }, {
        "p": { "x": 0, "y": 0, "z": 0 },
        "c": { "r": 0, "g": 0, "b": 0 },
        "o": { "x": 0, "y": 0, "z": 0, "w": 0 }
      }, {
        "p": { "x": 0, "y": 0, "z": 0 },
        "c": { "r": 0, "g": 0, "b": 0 },
        "o": { "x": 0, "y": 0, "z": 0, "w": 0 }
      }
    ]
  }
}`
	blue("Running: " + iAm() + " ExtraPacking")

	text = RexReplace(text, arrayOfObjRex, func(x []string, rx *regexp.Regexp, rf RexFunc) string {
		return x[1] + " " + x[2] + RexReplace(x[3], rx, rf)
	})
	if text != expected {
		green(expected)
		red(text)
		t.Fail()
	}
}

func TestRexJSONTightPacking(t *testing.T) {
	expected := `{
  "verts": [
    { "p": { "x": 0, "y": 0, "z": 0 },
      "c": { "r": 0, "g": 0, "b": 0 },
      "o": { "x": 0, "y": 0, "z": 0, "w": 0 } },
    { "p": { "x": 0, "y": 0, "z": 0 },
      "c": { "r": 0, "g": 0, "b": 0 },
      "o": { "x": 0, "y": 0, "z": 0, "w": 0 } },
    { "p": { "x": 0, "y": 0, "z": 0 },
      "c": { "r": 0, "g": 0, "b": 0 },
      "o": { "x": 0, "y": 0, "z": 0, "w": 0 } }
  ],
  "subverts": {
    "verts": [
      { "p": { "x": 0, "y": 0, "z": 0 },
        "c": { "r": 0, "g": 0, "b": 0 },
        "o": { "x": 0, "y": 0, "z": 0, "w": 0 } },
      { "p": { "x": 0, "y": 0, "z": 0 },
        "c": { "r": 0, "g": 0, "b": 0 },
        "o": { "x": 0, "y": 0, "z": 0, "w": 0 } },
      { "p": { "x": 0, "y": 0, "z": 0 },
        "c": { "r": 0, "g": 0, "b": 0 },
        "o": { "x": 0, "y": 0, "z": 0, "w": 0 } }
    ]
  }
}`
	blue("Running: " + iAm())
	cleanVerts := func(src string) string {
		// This will only pad out correctly if the json.Marshal uses 2 spaces for indent
		return RexJSONCleanupPost(src, UnnamedJSONObjectRex, func(s string) string {
			return RexJSONCleanup(s, NamedJSONObjectRex, PackLines)
		}, func(i string) string {
			return " " + strings.TrimSpace(i) + " "
		})
	}
	text := RexJSONCleanup(testText, vertsArrayRex, cleanVerts)
	if text != expected {
		green(expected)
		red(text)
		t.Fail()
	}
}

func TestRexJSONJustPackLines(t *testing.T) {
	expected := `{
  "verts": [
    { "p": { "x": 0, "y": 0, "z": 0 }, "c": { "r": 0, "g": 0, "b": 0 }, "o": { "x": 0, "y": 0, "z": 0, "w": 0 } },
    { "p": { "x": 0, "y": 0, "z": 0 }, "c": { "r": 0, "g": 0, "b": 0 }, "o": { "x": 0, "y": 0, "z": 0, "w": 0 } },
    { "p": { "x": 0, "y": 0, "z": 0 }, "c": { "r": 0, "g": 0, "b": 0 }, "o": { "x": 0, "y": 0, "z": 0, "w": 0 } }
  ],
  "subverts": {
    "verts": [
      { "p": { "x": 0, "y": 0, "z": 0 }, "c": { "r": 0, "g": 0, "b": 0 }, "o": { "x": 0, "y": 0, "z": 0, "w": 0 } },
      { "p": { "x": 0, "y": 0, "z": 0 }, "c": { "r": 0, "g": 0, "b": 0 }, "o": { "x": 0, "y": 0, "z": 0, "w": 0 } },
      { "p": { "x": 0, "y": 0, "z": 0 }, "c": { "r": 0, "g": 0, "b": 0 }, "o": { "x": 0, "y": 0, "z": 0, "w": 0 } }
    ]
  }
}`
	blue("Running: " + iAm())
	cleanVerts := func(src string) string {
		return RexJSONCleanup(src, UnnamedJSONObjectRex, func(s string) string {
			return PackLines(RexJSONCleanup(s, NamedJSONObjectRex, PackLines))
		})
	}

	text := RexJSONCleanup(testText, vertsArrayRex, cleanVerts)
	if text != expected {
		green(expected)
		red(text)
		t.Fail()
	}
}

// json.Marshal outputs any empty object or array as {} & [] -- this routine will
// convert those to the multi-line versions the rex functions need; which will later
// convert them back...
func fixEmptyObjectsAndArrays(src string) string {
	namedObjectRex := regexp.MustCompile(`((?sm).*?^)( *)("\w+?": {)((?s)}.*)`)
	// test only uses empty named objects (for now)
	//	namedArrayRex := regexp.MustCompile(`((?sm).*?^)( *)("\w+?": \[)((?s)\].*)`)
	//	unnamedObjectRex := regexp.MustCompile(`((?sm).*?^)( *)({)((?s)}.*)`)
	//	unnamedArrayRex := regexp.MustCompile(`((?sm).*?^)( *)(\[)((?s)\].*)`)

	src = RexReplace(src, namedObjectRex, func(x []string, rx *regexp.Regexp, rf RexFunc) string {
		return x[1] + x[2] + x[3] + "\n" + x[2] + RexReplace(x[4], rx, rf)
	})
	return src
}

func padPackMaxSubObjects(src string) string {
	return "\n" + RexJSONCleanup(src, NamedJSONObjectRex, func(s string) string { return PackLinesMax(s, 35) })
}

func paddedPackMaxSubObjects(src string) string {
	return "\n" + RexJSONCleanup(src, NamedJSONObjectRex, func(s string) string { return PackLinesMax(s, 35) })
}

func padPackObjsDepth2(src string) string {
	return "\n" + RexJSONCleanup(src, NamedJSONObjectRex, paddedPackMaxSubObjects)
}

func TestRexJSONObjectsOfObjects(t *testing.T) {
	// JSON as json.Marshall would output
	source := `{
  "string": "string data",
  "objectOfObjects": {
    "object0": {},
    "object1": {
      "a": 1
    },
    "object4": {
      "a": 1,
      "b": 2,
      "c": 3,
      "d": 4
    },
    "object8": {
      "a": 1,
      "b": 2,
      "c": 3,
      "d": 4,
      "e": 5,
      "f": 6,
      "g": 7,
      "h": 8
    }
  },
  "objectOfSubObjects": {
    "objectA": {
      "subObjectA0": {},
      "subObjectA1": {
        "a": "a"
      }
    },
    "objectB": {
      "subObjectB1": {
        "a": "a"
      },
      "subObjectB0": {}
    },
    "objectC": {
      "subObjectC": {
        "a": "a",
        "b": "b",
        "c": "c",
        "d": "d"
      }
    },
    "objectCL": {
      "subObjectC": {
        "a": "this is a very long string for testing",
        "b": "this is a very long string for testing",
        "c": "this is a very long string for testing",
        "d": "this is a very long string for testing"
      }
    },
    "objectD": {
      "subObjectD": {
        "a": "a",
        "b": "b",
        "c": "c",
        "d": "d",
        "e": "e",
        "f": "f",
        "g": "g",
        "h": "h"
      }
    }
  }
}`

	expected := `{
  "string": "string data",
  "objectOfObjects": {
    "object0": { },
    "object1": { "a": 1 },
    "object4": { "a": 1, "b": 2, "c": 3, "d": 4 },
    "object8": {
      "a": 1, "b": 2, "c": 3, "d": 4,
      "e": 5, "f": 6, "g": 7, "h": 8
    }
  },
  "objectOfSubObjects": {
    "objectA": {
      "subObjectA0": { },
      "subObjectA1": { "a": "a" }
    },
    "objectB": {
      "subObjectB1": { "a": "a" },
      "subObjectB0": { }
    },
    "objectC": {
      "subObjectC": {
        "a": "a", "b": "b", "c": "c",
        "d": "d"
      }
    },
    "objectCL": {
      "subObjectC": {
        "a": "this is a very long string for testing",
        "b": "this is a very long string for testing",
        "c": "this is a very long string for testing",
        "d": "this is a very long string for testing"
      }
    },
    "objectD": {
      "subObjectD": {
        "a": "a", "b": "b", "c": "c",
        "d": "d", "e": "e", "f": "f",
        "g": "g", "h": "h"
      }
    }
  }
}`
	objOfObjRex := regexp.MustCompile(`((?s).*?\n  "objectOfObjects": {)\n((?s).*?)((?s)\n  }.*)`)
	objOfSubObjRex := regexp.MustCompile(`((?s).*?\n  "objectOfSubObjects": {)\n((?s).*?)((?s)\n  }.*)`)

	source = fixEmptyObjectsAndArrays(source)

	blue("Running: " + iAm())
	text := RexJSONCleanup(source, objOfObjRex, padPackMaxSubObjects)
	text = RexJSONCleanup(text, objOfSubObjRex, padPackObjsDepth2)
	text = removeTrailingSpaces(text)
	if text != expected {
		green(expected)
		red(text)
		t.Fail()
	}

	packAllExpect := `{
  "string": "string data",
  "objectOfObjects": {
    "object0": {},
    "object1": { "a": 1 },
    "object4": { "a": 1, "b": 2, "c": 3, "d": 4 },
    "object8": { "a": 1, "b": 2, "c": 3, "d": 4, "e": 5, "f": 6, "g": 7, "h": 8 }
  },
  "objectOfSubObjects": {
    "objectA": {
      "subObjectA0": {},
      "subObjectA1": { "a": "a" }
    },
    "objectB": {
      "subObjectB1": { "a": "a" },
      "subObjectB0": {}
    },
    "objectC": {
      "subObjectC": { "a": "a", "b": "b", "c": "c", "d": "d" }
    },
    "objectCL": {
      "subObjectC": { "a": "this is a very long string for testing", "b": "this is a very long string for testing", "c": "this is a very long string for testing", "d": "this is a very long string for testing" }
    },
    "objectD": {
      "subObjectD": { "a": "a", "b": "b", "c": "c", "d": "d", "e": "e", "f": "f", "g": "g", "h": "h" }
    }
  }
}`
	blue("Running: " + iAm() + " Using packEverything")
	text = RexJSONCleanup(source, UnnamedJSONObjectRex, packEverything)
	text = removeBlankLines(text)
	text = removeTrailingSpaces(text)
	if text != packAllExpect {
		green(packAllExpect)
		red(text)
		t.Fail()
	}
}
