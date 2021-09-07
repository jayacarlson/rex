package rex

import (
	"regexp"
	"runtime"
	"strings"

	"github.com/jayacarlson/dbg"
)

func iAm() string {
	pc := make([]uintptr, 4)
	runtime.Callers(2, pc)
	nm := runtime.FuncForPC(pc[0]).Name()
	return nm[strings.LastIndex(nm, ".")+1:]
}

// some output funcions for debugging
func red(src string) {
	if src == "" {
		dbg.Error("EMPTY STRING")
		return
	}
	for len(src) > 0 {
		if i := strings.Index(src, "\n"); i >= 0 {
			dbg.Echo("\033[30;41m" + src[:i] + "\033[0m")
			src = src[i+1:]
		} else {
			dbg.Echo("\033[30;41m" + src + "\033[0m")
			return
		}
	}
}

func green(src string) {
	if src == "" {
		dbg.Error("EMPTY STRING")
		return
	}
	for len(src) > 0 {
		if i := strings.Index(src, "\n"); i >= 0 {
			dbg.Echo("\033[30;42m" + src[:i] + "\033[0m")
			src = src[i+1:]
		} else {
			dbg.Echo("\033[30;42m" + src + "\033[0m")
			return
		}
	}
}

func yellow(src string) {
	if src == "" {
		dbg.Error("EMPTY STRING")
		return
	}
	for len(src) > 0 {
		if i := strings.Index(src, "\n"); i >= 0 {
			dbg.Echo("\033[30;43m" + src[:i] + "\033[0m")
			src = src[i+1:]
		} else {
			dbg.Echo("\033[30;43m" + src + "\033[0m")
			return
		}
	}
}

func blue(src string) {
	if src == "" {
		dbg.Error("EMPTY STRING")
		return
	}
	for len(src) > 0 {
		if i := strings.Index(src, "\n"); i >= 0 {
			dbg.Echo("\033[30;44m" + src[:i] + "\033[0m")
			src = src[i+1:]
		} else {
			dbg.Echo("\033[30;44m" + src + "\033[0m")
			return
		}
	}
}

func magenta(src string) {
	if src == "" {
		dbg.Error("EMPTY STRING")
		return
	}
	for len(src) > 0 {
		if i := strings.Index(src, "\n"); i >= 0 {
			dbg.Echo("\033[30;45m" + src[:i] + "\033[0m")
			src = src[i+1:]
		} else {
			dbg.Echo("\033[30;45m" + src + "\033[0m")
			return
		}
	}
}

var (
	// these rex versions for object/array extract only if the text STARTS with them
	namedJSONObjectRex   = regexp.MustCompile(`(^"\w+": {)\n((?s).*?)((?sm)^}.*)`)
	namedJSONArrayRex    = regexp.MustCompile(`(^"\w+": \[)\n((?s).*?)((?sm)^\].*)`)
	unnamedJSONObjectRex = regexp.MustCompile(`(^{)\n((?s).*?)((?sm)^}.*)`)
	unnamedJSONArrayRex  = regexp.MustCompile(`(^\[)\n((?s).*?)((?sm)^\].*)`)
	blankLinesRex        = regexp.MustCompile(`((?s).*?)\n *((?s)\n.*)`)        // find `\n +\n`
	trailingSpacesRex    = regexp.MustCompile(`((?s).*?) +((?s)\n.*)`)          // find `... +\n`
	extraSpacesRex       = regexp.MustCompile(`((?s).*?[\[{] ) +((?s).*)`)      // find `{|[  +`
	collapsArrayChainRex = regexp.MustCompile(`((?s).*? +\[)((?s)\n.*)`)        // find ` +[\n`
	closeArrayChainRex   = regexp.MustCompile(`((?s).*?)((?m)^ +)\]((?s)\n.*)`) // find `^ +]\n`
	arrayChainRex        = regexp.MustCompile(`((?sm).*?^ *\],)\n *((?s)\[.*)`) // find `],\n +[`
)

func removeBlankLines(src string) string {
	return RexReplace(src, blankLinesRex, func(x []string, rx *regexp.Regexp, rf RexFunc) string {
		return x[1] + RexReplace(x[2], rx, rf)
	})
}

func removeTrailingSpaces(src string) string {
	return RexReplace(src, trailingSpacesRex, func(x []string, rx *regexp.Regexp, rf RexFunc) string {
		return x[1] + RexReplace(x[2], rx, rf)
	})
}

func removeExtraSpaces(src string) string {
	return RexReplace(src, extraSpacesRex, func(x []string, rx *regexp.Regexp, rf RexFunc) string {
		return x[1] + RexReplace(x[2], rx, rf)
	})
}

func collapsArrayChains(src string) (string, string) {
	if src[0] == '\n' {
		var i = 1
		for i < len(src) && src[i] == ' ' {
			i += 1
		}
		if src[i] == '[' && src[i+1] == '\n' {
			brace, more := collapsArrayChains(src[i+1:])
			return " [" + brace, more
		}
	}
	return "", src
}

func closeArrayChains(src string) (string, string) {
	if src[0] == '\n' && len(src) >= 2 {
		var i = 1
		for i < len(src) && src[i] == ' ' {
			i += 1
		}
		if len(src[i:]) > i && src[i] == ']' && (src[i+1] == '\n') || (src[i+1] == ',') {
			brace, more := closeArrayChains(src[i+1:])
			return "] " + brace, more
		}
	}
	return "]", src
}

func concatArrays(src string) string {
	src = RexReplace(src, arrayChainRex, func(x []string, rx *regexp.Regexp, rf RexFunc) string {
		return x[1] + " " + RexReplace(x[2], rx, rf)
	})
	src = RexReplace(src, collapsArrayChainRex, func(x []string, rx *regexp.Regexp, rf RexFunc) string {
		braces, more := collapsArrayChains(x[2])
		return x[1] + braces + RexReplace(more, rx, rf)
	})
	src = RexReplace(src, closeArrayChainRex, func(x []string, rx *regexp.Regexp, rf RexFunc) string {
		braces, more := closeArrayChains(x[3])
		return x[1] + strings.Repeat(" ", len(x[2])-2*(strings.Count(braces, "]")-1)) + braces + RexReplace(more, rx, rf)
	})
	return "\n" + src
}

func getMore(src string) (string, string) {
	// src : "},\n..."  |  "}\n..."  |  "}"
	//     | "],\n..."  |  "]\n..."  |  "]"
	i := strings.Index(src, "\n")
	if i < 0 {
		return src, ""
	}
	return src[:i], src[i:]
}

func autoPackObjAry(sub, more string) (string, string) {
	end, more := getMore(more)
	lead, sub := RemoveJSONPadding(sub)
	return AddJSONPadding(lead, autoPack(sub)) + end, more
}

func autoPack(src string) string {
	result := ""
	for src != "" {
		switch src[0] {
		case '\n':
			return result + "\n" + autoPack(src[1:])
		case '"':
			if x := namedJSONObjectRex.FindStringSubmatch(src); x != nil {
				sub, more := autoPackObjAry(x[2], x[3])
				return result + "\n" + x[1] + sub + autoPack(more)
			} else if x := namedJSONArrayRex.FindStringSubmatch(src); x != nil {
				sub, more := autoPackObjAry(x[2], x[3])
				return result + "\n" + x[1] + sub + autoPack(more)
			}
			// generic entry in object
			if result == "" {
				result = " "
			}
			i := strings.Index(src, "\n")
			if i < 0 {
				return result + src
			}
			result += src[:i] + " "
			src = src[i+1:]
		case '{':
			if x := unnamedJSONObjectRex.FindStringSubmatch(src); x != nil {
				sub, more := autoPackObjAry(x[2], x[3])
				return result + "\n" + x[1] + sub + autoPack(more)
			}
		// error...
		case '[':
			if x := unnamedJSONArrayRex.FindStringSubmatch(src); x != nil {
				sub, more := autoPackObjAry(x[2], x[3])
				return result + "\n" + x[1] + sub + autoPack(more)
			}
		default:
			// better be inside an array of things...
			if result == "" {
				result = " "
			}
			i := strings.Index(src, "\n")
			if i < 0 {
				return result + src
			}
			result += src[:i] + " "
			src = src[i+1:]
		}
	}
	return result + src
}

func packEverything(src string) string {
	src = autoPack(src)
	if src[0] == ' ' {
		src = "\n" + src[1:]
	}
	return src
}
