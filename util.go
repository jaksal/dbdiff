package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// RoutineDelimeter : mysql procedure delemiter
const RoutineDelimeter = "//"

func mysqlEscape(source []byte) string {
	if len(source) == 0 {
		return ""
	}

	dest := make([]byte, 0, 2*len(source))
	var escape byte
	for i := 0; i < len(source); i++ {
		c := source[i]

		escape = 0

		switch c {
		case 0: /* Must be escaped for 'mysql' */
			escape = '0'
		case '\n': /* Must be escaped for logs */
			escape = 'n'
		case '\r':
			escape = 'r'
		case '\\':
			escape = '\\'
		case '\'':
			escape = '\''
		case '"': /* Better safe than sorry */
			escape = '"'
		case '\032': /* This gives problems on Win32 */
			escape = 'Z'
		}

		if escape != 0 {
			dest = append(dest, '\\', escape)
		} else {
			dest = append(dest, c)
		}
	}

	return string(dest)
}

//CheckFilePath :  Check ouptput file path.
func CheckFilePath(name string) string {
	name = strings.Replace(name, "[DATE]", time.Now().Format("20060102"), -1)
	name = strings.Replace(name, "[DATETIME]", time.Now().Format("20060102150405"), -1)

	dir := filepath.Dir(name)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalln("mkdir err", dir, err)
			return ""
		}
	}
	return name
}

/*
\ (역방향 슬래시) ... 따라서 역방향 슬래시를 그대로 표현하려면 역방향 슬래시 기호를 두 번 쓰면 됩니다. \\ 이렇게요.
` (backtick = 제2강세 악센트 기호)
* 별표
_ 언더바
{ } 중괄호
[ ] 대괄호
( ) 소괄호
# 우물정자
+ 플러스기호
- 마이너스 기호
. 마침표
! 느낌표
*/

var r = strings.NewReplacer(
	"\\", "\\\\",
	"`", "\\`",
	"*", "\\*",
	"_", "\\_",
	"{", "\\{",
	"{", "\\}",
	"[", "\\[",
	"]", "\\]",
	"(", "\\(",
	")", "\\)",
	"#", "\\#",
	"+", "\\+",
	"-", "\\-",
	".", "\\.",
	"!", "\\!")

// MDReplace replace md special char
func MDReplace(src string) string {
	ret := strings.ToLower(src)
	ret = strings.Replace(ret, " ", "", -1)
	ret = strings.Replace(ret, "_", "", -1)
	return ret
}

// WikiReplace replace wiki special char
func WikiReplace(src string) string {
	//return r.Replace(src)
	return strings.Replace(src, " ", "", -1)
}

// CommentReplace mysql comment replace
func CommentReplace(src string) string {
	return strings.Replace(src, "'", "\\'", -1)
}
