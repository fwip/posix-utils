package locale

import "fmt"
import "strings"
import "os"

// Locale holds all the information needed for localization of posix utils
type Locale struct {
	lang     string // Default for unset values
	collate  string
	ctype    string
	messages string
	monetary string
	numeric  string
	time     string
	all      string // Overrides all options beside lang
}

// Def is the settings & stuff in a locale
type Def struct {
	monetary Monetary
	ctype    Ctype
	time     Time
	numeric  Numeric
	collate  Collate
}

// getVal will process overrides or fallbacks as necessary
func (l Locale) getVal(val string) string {
	if l.all != "" {
		return "\"" + l.all + "\""
	}
	if val != "" {
		return val
	}
	return "\"" + l.lang + "\""
}

func (l Locale) String() string {
	var out strings.Builder
	out.WriteString(fmt.Sprintf("LANG=\"%s\"\n", l.lang))
	settings := []struct{ name, val string }{
		{"LC_COLLATE", l.collate},
		{"LC_CTYPE", l.ctype},
		{"LC_MESSAGES", l.messages},
		{"LC_MONETARY", l.monetary},
		{"LC_NUMERIC", l.numeric},
		{"LC_TIME", l.time},
	}

	for _, setting := range settings {
		out.WriteString(fmt.Sprintf("%s=%s\n", setting.name, l.getVal(setting.val)))
	}

	out.WriteString(fmt.Sprintf("LC_ALL=%s\n", l.all))
	return out.String()
}

// FromEnv returns a new Locale based on environment variables
func FromEnv() Locale {
	return Locale{
		lang:     os.Getenv("LANG"),
		collate:  os.Getenv("LC_COLLATE"),
		ctype:    os.Getenv("LC_CTYPE"),
		messages: os.Getenv("LC_MESSAGES"),
		monetary: os.Getenv("LC_MONETARY"),
		numeric:  os.Getenv("LC_NUMERIC"),
		time:     os.Getenv("LC_TIME"),
		all:      os.Getenv("LC_ALL"),
	}
}

func main() {
	fmt.Println("vim-go")
}
