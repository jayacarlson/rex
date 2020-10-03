package rex

import (
	"regexp"
	"testing"
)

var complexSource = `{
  "mbid": "ABCDABCDABCDABCDABCDABCDABC-0",
  "artist": "FooBar & Boo",
  "artist-match": "FooBarBoo",
  "save-artist": "FooBar",
  "album": "The Best Of: FooBar",
  "save-album": "Best Of FooBar",
  "length": "32:47",
  "mbuuid": "12345678-1234-abcd-efab-0123456789ab",
  "freedbID": "ab123456",
  "source": "mbrainz",
  "trackCount": 7,
  "tracks": [
    {
      "title": "All The Way To Foobar",
      "number": 1,
      "trk-stt": "0:00",
      "trk-end": "3:42'659!48",
      "trk-len": "3:42.733",
      "aud-stt": "0:00'244",
      "aud-end": "3:39'227!48",
      "aud-len": "3:38.982"
    },
    {
      "title": "Back To Foobar",
      "artist": "Foo & The Bars",
      "number": 2,
      "trk-stt": "3:42'660",
      "trk-end": "7:56'179!48",
      "trk-len": "4:13.467",
      "aud-stt": "3:42'897",
      "aud-end": "7:53'444!48",
      "aud-len": "4:10.498"
    },
    {
      "title": "Foobar All Night Long",
      "number": 3,
      "trk-stt": "7:56'180",
      "trk-end": "12:18'23!48",
      "trk-len": "4:21.827",
      "aud-stt": "7:56'451",
      "aud-end": "12:15'132!48",
      "aud-len": "4:18.647"
    },
    {
      "title": "Where The Foobar Are You?",
      "number": 4,
      "trk-stt": "12:18'24",
      "trk-end": "16:20'179!48",
      "trk-len": "4:02.173",
      "aud-stt": "12:18'304",
      "aud-end": "16:16'309!48",
      "aud-len": "3:58.007"
    },
    {
      "title": "Kick The Foobar",
      "artist": "Foo & The Bars (f/ Boo)",
      "number": 5,
      "trk-stt": "16:20'180",
      "trk-end": "24:43'899!48",
      "trk-len": "8:23.800",
      "aud-stt": "16:20'414",
      "aud-end": "24:39'163!48",
      "aud-len": "8:18.722",
      "silences": [
        {
          "stt": "22:13'395!32",
          "end": "22:13'742!20"
        },
        {
          "stt": "22:15'102!08",
          "end": "22:16'4!42"
        },
        {
          "stt": "23:37'217!22",
          "end": "23:38'193!28"
        }
      ]
    },
    {
      "title": "When I Reach Foobar, I'm Happy",
      "artist": "FooBarBoo",
      "number": 6,
      "trk-stt": "24:44",
      "trk-end": "27:58'803!48",
      "trk-len": "3:14.893",
      "aud-stt": "24:44'268",
      "aud-end": "27:54'791!48",
      "aud-len": "3:10.582"
    },
    {
      "title": "Merry Foo - Happy Bar",
      "number": 7,
      "trk-stt": "27:58'804",
      "trk-end": "32:46'899!48",
      "trk-len": "4:48.107",
      "aud-stt": "27:59'164",
      "aud-end": "32:43'655!48",
      "aud-len": "4:44.547"
    }
  ]
}`

func silencesCleaner(src string) string {
	return "\n" + RexJSONCleanup(src, UnnamedJSONObjectRex, PackLines)
}

func trackCleaner(src string) string {
	var titleRex = regexp.MustCompile(`((?s).*?"title": )((?s).*)`)
	var artistRex = regexp.MustCompile(`((?s).*?"artist": )((?s).*)`)
	var trkRex = regexp.MustCompile(`((?s).*?trk-stt":)((?s).*?trk-len.+?)(\n(?s).*)`)
	var audRex = regexp.MustCompile(`((?s).*?aud-stt":)((?s).*?aud-len.+?)(\n(?s).*)`)
	src = RexReplace(src, titleRex, func(x []string, rx *regexp.Regexp, rf RexFunc) string {
		return x[1] + "       " + RexReplace(x[2], rx, rf)
	})
	src = RexReplace(src, artistRex, func(x []string, rx *regexp.Regexp, rf RexFunc) string {
		return x[1] + "      " + RexReplace(x[2], rx, rf)
	})
	src = RexCleanup(src, trkRex, PackLines)
	src = RexCleanup(src, audRex, PackLines)
	src = RexJSONCleanup(src, NamedJSONArrayRex, silencesCleaner)
	return src
}

func cleanTrack(src string) string {
	return RexJSONCleanup(src, UnnamedJSONObjectRex, trackCleaner)
}

func TestComplexJSON(t *testing.T) {
	expected := `{
  "mbid": "ABCDABCDABCDABCDABCDABCDABC-0",
  "artist":       "FooBar & Boo",
  "artist-match": "FooBarBoo",
  "save-artist":  "FooBar",
  "album":        "The Best Of: FooBar",
  "save-album":   "Best Of FooBar",
  "length": "32:47",
  "mbuuid": "12345678-1234-abcd-efab-0123456789ab",
  "freedbID": "ab123456",
  "source": "mbrainz",
  "trackCount": 7,
  "tracks": [
    { "title":        "All The Way To Foobar",
      "number": 1,
      "trk-stt": "0:00", "trk-end": "3:42'659!48", "trk-len": "3:42.733",
      "aud-stt": "0:00'244", "aud-end": "3:39'227!48", "aud-len": "3:38.982"
    },
    { "title":        "Back To Foobar",
      "artist":       "Foo & The Bars",
      "number": 2,
      "trk-stt": "3:42'660", "trk-end": "7:56'179!48", "trk-len": "4:13.467",
      "aud-stt": "3:42'897", "aud-end": "7:53'444!48", "aud-len": "4:10.498"
    },
    { "title":        "Foobar All Night Long",
      "number": 3,
      "trk-stt": "7:56'180", "trk-end": "12:18'23!48", "trk-len": "4:21.827",
      "aud-stt": "7:56'451", "aud-end": "12:15'132!48", "aud-len": "4:18.647"
    },
    { "title":        "Where The Foobar Are You?",
      "number": 4,
      "trk-stt": "12:18'24", "trk-end": "16:20'179!48", "trk-len": "4:02.173",
      "aud-stt": "12:18'304", "aud-end": "16:16'309!48", "aud-len": "3:58.007"
    },
    { "title":        "Kick The Foobar",
      "artist":       "Foo & The Bars (f/ Boo)",
      "number": 5,
      "trk-stt": "16:20'180", "trk-end": "24:43'899!48", "trk-len": "8:23.800",
      "aud-stt": "16:20'414", "aud-end": "24:39'163!48", "aud-len": "8:18.722",
      "silences": [
        { "stt": "22:13'395!32", "end": "22:13'742!20" },
        { "stt": "22:15'102!08", "end": "22:16'4!42" },
        { "stt": "23:37'217!22", "end": "23:38'193!28" }
      ]
    },
    { "title":        "When I Reach Foobar, I'm Happy",
      "artist":       "FooBarBoo",
      "number": 6,
      "trk-stt": "24:44", "trk-end": "27:58'803!48", "trk-len": "3:14.893",
      "aud-stt": "24:44'268", "aud-end": "27:54'791!48", "aud-len": "3:10.582"
    },
    { "title":        "Merry Foo - Happy Bar",
      "number": 7,
      "trk-stt": "27:58'804", "trk-end": "32:46'899!48", "trk-len": "4:48.107",
      "aud-stt": "27:59'164", "aud-end": "32:43'655!48", "aud-len": "4:44.547"
    }
  ]
}`
	artistRex := regexp.MustCompile(`((?sm).*?^  "artist": )((?sm).*)`)
	artistMatchRex := regexp.MustCompile(`((?sm).*?^  "artist-match": )((?sm).*)`)
	artistSaveRex := regexp.MustCompile(`((?sm).*?^  "save-artist": )((?sm).*)`)
	albumRex := regexp.MustCompile(`((?sm).*?^  "album": )((?sm).*)`)
	albumMatchRex := regexp.MustCompile(`((?sm).*?^  "album-match": )((?sm).*)`)
	albumSaveRex := regexp.MustCompile(`((?sm).*?^  "save-album": )((?sm).*)`)
	tracksRex := regexp.MustCompile(`((?sm).*?^  "tracks": \[\n)((?s).*?)((?sm)^  \].*)`)

	blue("Running: " + iAm())

	text := RexJSONCleanup(complexSource, tracksRex, cleanTrack)
	text = removeExtraSpaces(text)
	if x := artistRex.FindStringSubmatch(text); x != nil {
		text = x[1] + "      " + x[2]
	}
	if x := artistMatchRex.FindStringSubmatch(text); x != nil {
		text = x[1] + "" + x[2]
	}
	if x := artistSaveRex.FindStringSubmatch(text); x != nil {
		text = x[1] + " " + x[2]
	}
	if x := albumRex.FindStringSubmatch(text); x != nil {
		text = x[1] + "       " + x[2]
	}
	if x := albumMatchRex.FindStringSubmatch(text); x != nil {
		text = x[1] + " " + x[2]
	}
	if x := albumSaveRex.FindStringSubmatch(text); x != nil {
		text = x[1] + "  " + x[2]
	}

	if text != expected {
		green(expected)
		red(text)
		t.Fail()
	}
}
