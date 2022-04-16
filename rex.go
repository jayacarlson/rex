package rex

import (
	"regexp"
	"strings"
)

/*
	Utility functions using regular expressions.

	CleanerFunc: TYPE
		Utility functions that take a string to do any cleanup work, return cleaned string

	RexFunc: TYPE
		Utility functions that take:
			x []string:			String array that rx will generate and that rf can handle
			rx regexp.Regexp:	Rex to extract further data to work on as the []string
			rf RexFunc:			The RexFunc to call again with text for further processing

			The []string can be in whatever format needed for the RexFunc handler

	PackLines: CleanerFunc
		Utility function, mainly as a JSON cleanup function
		Creates a RexFunc that uses RexPlace to do all the work
		Returns a string with all \n replaced by a single ' '
		Any leading spaces are also removed from the source lines, as in:
			`   apple,
			 banana,			=>  ` apple, banana, cherry `
			    cherry`

	PackLinesMax: CleanerFunc
		More complex utility function, mainly as a JSON cleanup function
		As PackLines, but with a maximum length check for the joined lines

	RexReplace:
		Used to recursivly process data for re-formatting using a RexFunc & regexp
		to process the text. The RexFunc usually just calls itself for further
		processing until the regexp fails.  See PackLines for an example.

	RexGather:
		Used to gather up text that matches a regular expression using a GatherFunc & regexp
		to process the text. The GatherFunc usually just calls itself for further
		processing until the regexp fails.  RexGather expects the regexp to generate
		the []string in the following format:
			x[1] == all leading text -- not modified
			x[2] == text matching the regexp
			x[3] == trailing text, passed back into RexGather

		Sample that grabs things and places them in a semi-colon seperated list:

		gathered = RexGather(src, regExp, func(g string, x []string, rx *regexp.Regexp, rf GatherFunc) string {
			return x[1] + "; " + RexGather(x[2], gather, rx, rf)
		})
		-- gathered could then have any trailing "; " removed
		-- gathered could then be searched for unique entries if needed

	RexCleanup:
		Used to recursivly process text for re-formatting using a RexFunc & regexp
		to process the text. The RexFunc usually just calls itself for further
		processing until the regexp fails.  RexCleanup expects the regexp to generate
		the []string in the following format:
			x[1] == all leading text -- not modified
			x[2] == text given to the CleanerFunc
			x[3] == trailing text, passed back into RexCleanup

	RexJSONCleanup:
		Used to recursivly process JSON for re-formatting using a RexFunc & regexp
		to process the JSON. The RexFunc usually just calls itself for further JSON
		processing until the regexp fails.  RexCleanupJSON expects the regexp to
		generate the []string in the following format:
			x[1] == all leading text -- not modified
			x[2] == text given to the CleanerFunc
			x[3] == trailing text, passed back into RexCleanupJSON
		Before calling the CleanerFunc, the text has the leading spaces removed from all
		lines, with the 1st line used as a guide for how many spaces to remove.  The lead
		spaces are then returned to all resulting lines.

	RexJSONCleanupPost:
		Like RexCleanupJSON, but has an additional CleanerFunc to do any post cleanup

	RemoveJSONPadding:
		A routine to remove the indentation of a captured object / array for further processing

	AddJSONPadding:
		A routine to replace the indentation after cleanup

	Here are a series of variables that can be used to capture simple JSON objects and arrays.
	They are designed for use by the CleanerFunc called from RexJSONCleanup.  These regexp do
	not capture the surrounding newlines to give the greatest flexability for cleaners.

	NamedJSONObjectRex: VAR
		Regex to capture simple "named": {object}
	UnnamedJSONObjectRex: VAR
		Regex to capture simple unnamed {object}
	NamedJSONArrayRex: VAR
		Regex to capture simple "named": [array]
	UnnamedJSONArrayRex: VAR
		Regex to capture simple unnamed [array]

	Custom regexps can be used to extract objects and arrays, but care must be taken to match
	the leading { and [ with the correct closing } and ], e.g.
		Named object at a depth of 6	`((?s).*?\n {6}"<NAME>": {\n)((?s).*?)((?s)\n {6}}.*)`
		Unamed array at a depth of 8	`((?s).*?\n {8}\[\n)((?s).+?)((?s)\n {8}\].*)`

	Some suggestions and ideas:
		It is usually better to cascade inward than to start at depth - finding named objects
		and arrays and cleaning them as later cleaners can undu earlier cleanings.
		If you want the same behavior for several named arrays / objects, you can use:
			... "(?:NAME1|NAME2|NAME3...)": ...
		You can modify the behavior of the regexp by including the \n in the different capture
		expressions.
*/

type CleanerFunc func(string) string
type RexFunc func([]string, *regexp.Regexp, RexFunc) string
type GatherFunc func([]string, *regexp.Regexp, GatherFunc)

var (
	// Capture Named objects / arrays
	NamedJSONObjectRex = regexp.MustCompile(`((?sm).*?^"\w+": {)\n((?s).*?)((?sm)^}.*)`)
	NamedJSONArrayRex  = regexp.MustCompile(`((?sm).*?^"\w+": \[)\n((?s).*?)((?sm)^\].*)`)

	// Capture unnamed objects / arrays -- can also be used to start at the JSON root
	UnnamedJSONObjectRex = regexp.MustCompile(`((?sm).*?^{)\n((?s).*?)((?sm)^}.*)`)
	UnnamedJSONArrayRex  = regexp.MustCompile(`((?sm).*?^\[)\n((?s).*?)((?sm)^\].*)`)

	// used internally
	lineRex = regexp.MustCompile("(.*\n)((?s).*)")   // x[1] == 1st line (including \n), x[2] == all the rest
	packRex = regexp.MustCompile(" *(.*)\n((?s).*)") // x[1] == 1st line (leading spaces removed, no \n), x[2] == all the rest
)

func PackLines(src string) string {
	src = RexReplace(src, packRex, func(x []string, rx *regexp.Regexp, rf RexFunc) string {
		return x[1] + " " + RexReplace(x[2], rx, rf)
	})
	return " " + src
}

func PackLinesMax(src string, max int) string {
	cur, result := "", ""
	for x := packRex.FindStringSubmatch(src); x != nil; x = packRex.FindStringSubmatch(x[2]) {
		if cur != "" && len(cur)+len(x[1]) > max {
			result += "\n" + cur
			cur = ""
		}
		cur += x[1] + " "
	} // last x[2] is most likely a bunch of spaces which are tossed
	if result == "" {
		result = " " + cur
	} else {
		result += "\n" + cur + "\n"
	}
	return result
}

func RexReplace(src string, rx *regexp.Regexp, rf RexFunc) string {
	if x := rx.FindStringSubmatch(src); x != nil {
		src = rf(x, rx, rf)
	}
	return src
}

func RexGather(src string, rx *regexp.Regexp, gf GatherFunc) {
	if x := rx.FindStringSubmatch(src); x != nil {
		gf(x, rx, gf)
	}
}

// regexp must extract data as x[1]: lead data  x[2]: SubBlock  x[3]: tail data
func RexCleanup(src string, rx *regexp.Regexp, cf CleanerFunc) string {
	if x := rx.FindStringSubmatch(src); x != nil {
		return x[1] + cf(x[2]) + RexCleanup(x[3], rx, cf)
	}
	return src
}

// regexp must extract data as x[1]: lead data  x[2]: SubBlock  x[3]: tail data
func RexJSONCleanup(src string, rx *regexp.Regexp, cf CleanerFunc) string {
	if x := rx.FindStringSubmatch(src); x != nil {
		lead, sub := RemoveJSONPadding(x[2])
		return x[1] + AddJSONPadding(lead, cf(sub)) + RexJSONCleanup(x[3], rx, cf)
	}
	return src
}

// regexp must extract data as x[1]: lead data  x[2]: SubBlock  x[3]: tail data
func RexJSONCleanupPost(src string, rx *regexp.Regexp, cf CleanerFunc, post CleanerFunc) string {
	if x := rx.FindStringSubmatch(src); x != nil {
		lead, sub := RemoveJSONPadding(x[2])
		return x[1] + post(AddJSONPadding(lead, cf(sub))) + RexJSONCleanupPost(x[3], rx, cf, post)
	}
	return src
}

func RemoveJSONPadding(src string) (string, string) {
	result := ""
	if len(src) == 0 || src[0] != ' ' {
		return result, src
	}
	var l = 0
	for l < len(src) && src[l] == ' ' {
		l += 1
	}
	lead := src[:l]
	for len(src) >= l {
		if i := strings.Index(src, "\n"); i >= 0 {
			result += src[l : i+1]
			src = src[i+1:]
		} else {
			return lead, result + src[l:]
		}
	}
	return lead, result + src
}

func AddJSONPadding(lead, src string) string {
	result := ""
	// check for '\n' at beginning and skip it
	if src != "" && src[0] == '\n' {
		result = "\n"
		src = src[1:]
	}
	for len(src) > 0 {
		if i := strings.Index(src, "\n"); i >= 0 {
			result += lead + src[:i+1]
			src = src[i+1:]
		} else if result == "" { // single line result -- assume packed line
			return src
		} else {
			return result + lead + src // lead needed for multi-line packs
		}
	}
	return result
}
