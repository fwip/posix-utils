package locale

//go:generate goyacc localedef.y

import (
	"fmt"
	"strings"

	"unicode/utf8"

	"github.com/bbuck/go-lexer"
)

const upperChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ_"
const lowerChars = "abcdefghijklmnopqrstuvwxyz_"
const digits = "0123456789"
const whiteSpace = " \t"
const newLine = "\n\r"

type decode struct {
	typ  lexer.TokenType
	grab grabFunc
}

var escapeChar = '\\'

// Ctype mappings
var ctypeMap = map[string]decode{
	"upper":     {UPPER_STR, grabCharList},
	"lower":     {LOWER_STR, grabCharList},
	"alpha":     {ALPHA_STR, grabCharList},
	"digit":     {DIGIT_STR, grabCharList},
	"alnum":     {ALNUM_STR, grabCharList},
	"space":     {SPACE_STR, grabCharList},
	"cntrl":     {CNTRL_STR, grabCharList},
	"punct":     {PUNCT_STR, grabCharList},
	"graph":     {GRAPH_STR, grabCharList},
	"print":     {PRINT_STR, grabCharList},
	"xdigit":    {XDIGIT_STR, grabCharList},
	"blank":     {BLANK_STR, grabCharList},
	"charclass": {CHARCLASS_STR, grabStringList},
	"toupper":   {TOUPPER_STR, grabCharPairList},
	"tolower":   {TOLOWER_STR, grabCharPairList},
}

// Time mappings
var timeMap = map[string]decode{
	"abday":       {ABDAY_STR, grabStringList},
	"day":         {DAY_STR, grabStringList},
	"abmon":       {ABMON_STR, grabStringList},
	"mon":         {MON_STR, grabStringList},
	"d_t_fmt":     {D_T_FMT_STR, grabString},
	"d_fmt":       {D_FMT_STR, grabString},
	"t_fmt":       {T_FMT_STR, grabString},
	"am_pm":       {AM_PM_STR, grabStringList},
	"t_fmt_ampm":  {T_FMT_AMPM_STR, grabString},
	"era":         {ERA_STR, grabStringList},
	"era_d_fmt":   {ERA_D_FMT_STR, grabString},
	"era_t_fmt":   {ERA_T_FMT_STR, grabString},
	"era_d_t_fmt": {ERA_D_T_FMT_STR, grabString},
	"alt_digits":  {ALT_DIGITS_STR, grabStringList},
}

// Numeric mappings
var numericMap = map[string]decode{
	"decimal_point": {DECIMAL_POINT_STR, grabString},
	"thousands_sep": {THOUSANDS_SEP_STR, grabString},
	"grouping":      {GROUPING_STR, grabInt},
}

// Monetary mappings
var monetaryMap = map[string]decode{
	"int_curr_symbol":    {INT_CURR_SYMBOL_STR, grabString},
	"currency_symbol":    {CURRENCY_SYMBOL_STR, grabString},
	"mon_decimal_point":  {MON_DECIMAL_POINT_STR, grabString},
	"mon_thousands_sep":  {MON_THOUSANDS_SEP_STR, grabString},
	"mon_grouping":       {MON_GROUPING_STR, grabIntList},
	"positive_sign":      {POSITIVE_SIGN_STR, grabString},
	"negative_sign":      {NEGATIVE_SIGN_STR, grabString},
	"int_frac_digits":    {INT_FRAC_DIGITS_STR, grabInt},
	"frac_digits":        {FRAC_DIGITS_STR, grabInt},
	"p_cs_precedes":      {P_CS_PRECEDES_STR, grabInt},
	"p_sep_by_space":     {P_SEP_BY_SPACE_STR, grabInt},
	"n_cs_precedes":      {N_CS_PRECEDES_STR, grabInt},
	"n_sep_by_space":     {N_SEP_BY_SPACE_STR, grabInt},
	"p_sign_posn":        {P_SIGN_POSN_STR, grabInt},
	"n_sign_posn":        {N_SIGN_POSN_STR, grabInt},
	"int_p_cs_precedes":  {INT_P_CS_PRECEDES_STR, grabInt},
	"int_p_sep_by_space": {INT_P_SEP_BY_SPACE_STR, grabInt},
	"int_n_cs_precedes":  {INT_N_CS_PRECEDES_STR, grabInt},
	"int_n_sep_by_space": {INT_N_SEP_BY_SPACE_STR, grabInt},
	"int_p_sign_posn":    {INT_P_SIGN_POSN_STR, grabInt},
	"int_n_sign_posn":    {INT_N_SIGN_POSN_STR, grabInt},
}

// Message mappings
var msgMap = map[string]decode{
	"yesexpr": {YESEXPR_STR, grabMsg},
	"noexpr":  {NOEXPR_STR, grabMsg},
}

// WHITESPACE is a pseudo-token that the lexer doesn't want
const WHITESPACE = 123987654

// Lexer does some lexing, y'all
type Lexer struct {
	lexer.L
	def Def
}

// Lex does the lex
func (l *Lexer) Lex(lval *yySymType) int {
	tok, done := l.NextToken()
	if done {
		return int(lexer.EOFRune)
	}
	// Ignore whitespace and try again
	if tok.Type == WHITESPACE {
		return l.Lex(lval)
	}
	lval.val = tok.Value
	return int(tok.Type)
}

// sInit is the starting function for most lexers
func sInit(l *lexer.L) lexer.StateFunc {
	// We don't do escape/comment char processing here because it's preprocessed
	return sCategory
}

func sCategory(l *lexer.L) lexer.StateFunc {
	var next lexer.StateFunc
	var typ lexer.TokenType

	l.Take(upperChars + "_")
	header := l.Current()
	switch header {
	case "LC_CTYPE":
		typ = LC_CTYPE_STR
		next = getsGenericCategory(header, typ, ctypeMap)
	case "LC_COLLATE":
		typ = LC_COLLATE_STR
		next = sCollateCategory
	case "LC_MONETARY":
		typ = LC_MONETARY_STR
		next = getsGenericCategory(header, typ, monetaryMap)
	case "LC_NUMERIC":
		typ = LC_NUMERIC_STR
		next = getsGenericCategory(header, typ, numericMap)
	case "LC_TIME":
		typ = LC_TIME_STR
		next = getsGenericCategory(header, typ, timeMap)
	case "LC_MESSAGES":
		typ = LC_MESSAGES_STR
		next = getsGenericCategory(header, typ, msgMap)
	}
	if next == nil {
		return nil
	}
	l.Emit(typ)
	l.Take(whiteSpace)
	must(grabNewline, l)
	return next
}

// A grabFunc will consume and Emit tokens.
// Returns true if any tokens are emitted, returns false otherwise
type grabFunc func(l *lexer.L) (ok bool)

// must should be used with 'grab' functions - if it doesn't grab it, it dies
func must(grabber grabFunc, l *lexer.L) {
	ok := grabber(l)
	if !ok {
		l.Error("oops")
	}
}

func grabWhitespace(l *lexer.L) bool {
	l.Take(whiteSpace)
	if l.Current() == "" {
		return false
	}
	l.Emit(WHITESPACE)
	return true
}
func grabNewline(l *lexer.L) bool {
	l.Take(newLine)
	if l.Current() == "" {
		return false
	}
	l.Emit(EOL)
	return true
}
func grabChar(l *lexer.L) bool {
	if l.Next() != '<' {
		l.Rewind()
		return false
	}
	l.Take(upperChars + lowerChars + "-" + "_" + digits)
	if l.Next() != '>' {
		for range l.Current() {
			l.Rewind()
		}
		return false
	}
	l.Emit(CHARSYMBOL)
	return true
}

func grabCharList(l *lexer.L) (ok bool) {
	for grabChar(l) {
		ok = true
		if l.Next() == ';' {
			l.Emit(';')
		} else {
			l.Rewind()
			return true
		}
	}
	return ok
}
func grabCharPair(l *lexer.L) (ok bool) {
	if l.Next() != '(' {
		l.Rewind()
		return false
	}
	l.Emit('(')
	if !grabChar(l) {
		return false
	}
	if l.Next() != (',') {
		l.Rewind()
		return false
	}
	l.Emit(',')
	if !grabChar(l) {
		return false
	}
	if l.Next() != (')') {
		l.Rewind()
		return false
	}
	l.Emit(')')
	return true
}

func grabCharPairList(l *lexer.L) (ok bool) {
	for grabCharPair(l) {
		ok = true
		if l.Next() == ';' {
			l.Emit(';')
		} else {
			l.Rewind()
			return true
		}
	}
	return ok
}

func grabIntList(l *lexer.L) (ok bool) {
	for grabInt(l) {
		ok = true
		if l.Next() == ';' {
			l.Emit(';')
		} else {
			l.Rewind()
			return true
		}
	}
	return ok
}
func grabInt(l *lexer.L) bool {
	l.Take("-1234567890")
	n := l.Current()
	if n == "" {
		return false
	}
	if n == "-" || (len(n) > 2 && strings.ContainsRune(n[1:], '-')) {
		l.Error(fmt.Sprintf("expected integer, got %s", n))
	}
	l.Emit(NUMBER)
	return true
}

// Handle escaped characters in strings differently
func grabEscapeInString(l *lexer.L) (ok bool) {
	r := l.Next()
	switch r {
	case ',', ';', '<', '>', escapeChar:
		l.Emit(CHAR)
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		// octal - slurp up two more characters
		l.Next()
		l.Next()
		l.Emit(OCTAL_CHAR)
	case 'd':
		// decimal - slurp up 2 or 3 characters
		l.Next()
		l.Next()
		if !strings.ContainsRune(digits, l.Next()) {
			l.Rewind()
		}
		l.Emit(DECIMAL_CHAR)

	case 'x':
		// hex - slurp up 2 characters
		l.Next()
		l.Next()
		l.Emit(HEX_CHAR)
	default:
		l.Error(fmt.Sprintf("%q not valid after escape char %q", r, escapeChar))
		return false
	}
	return true
}

func grabString(l *lexer.L) (ok bool) {
	escapeChar := '\\' // TODO
	if l.Peek() != '"' {
		return false
	}
	l.Next()
	if l.Peek() == '"' {
		l.Next()
		l.Emit(TWO_DOUBLE_QUOTE_STR)
		return true
	}
	l.Emit('"')
	escapeNext := false
	for {
		r := l.Next()
		// TODO: Escape character handling is likely broken at the moment.
		if escapeNext {
			l.Rewind()
			grabEscapeInString(l)
			escapeNext = false
			continue
		}
		if r == '"' {
			l.Emit('"')
			return true
		}
		if r == '\n' {
			l.Error(fmt.Sprintf("Aww jeez. Got a newline before an end quote"))
			return false
		}
		if r == escapeChar {
			escapeNext = true
		} else {
			l.Emit(CHAR)
		}
	}
}
func grabStringList(l *lexer.L) (ok bool) {
	for grabString(l) {
		ok = true
		if l.Next() == ';' {
			l.Emit(';')
		} else {
			l.Rewind()
			return true
		}
	}
	return ok
}
func grabMsg(l *lexer.L) (ok bool) {
	if l.Next() != '"' {
		l.Rewind()
		return false
	}
	l.Emit('"')
	for l.Next() != '"' {
	}
	l.Rewind()
	l.Emit(EXTENDED_REG_EXP)
	l.Next()
	l.Emit('"')
	return true
}

func sCollateCategory(l *lexer.L) lexer.StateFunc {
	grabWhitespace(l)
	l.Take(upperChars + lowerChars + "_")
	word := l.Current()
	switch word {
	case "END":
		l.Emit(END_STR)
		must(grabWhitespace, l)
		l.Take(upperChars)
		if l.Current() != "LC_COLLATE" {
			l.Error(fmt.Sprintf("wrong category end: expected %s, got %s", "LC_COLLATE", l.Current()))
		}
		l.Emit(LC_COLLATE_STR)
		grabNewline(l)
		return sCategory
	case "order_start":
		l.Emit(ORDER_START_STR)
		return sCollateOrderStart
	case "collating-element":
		l.Emit(COLLATING_ELEMENT_STR)
		must(grabWhitespace, l)
		must(grabChar, l) // TODO: Is this right?
		must(grabWhitespace, l)
		l.Take("from")
		if l.Current() != "from" {
			l.Error("Expected 'from' after 'collating-element'")
		}
		l.Emit(FROM_STR)
		must(grabWhitespace, l)
		must(grabString, l)
		grabWhitespace(l)
		must(grabNewline, l)
	case "collating-symbol":
		l.Emit(COLLATING_SYMBOL_STR)
		must(grabWhitespace, l)
		must(grabString, l)
		grabWhitespace(l)
		must(grabNewline, l)
	}

	return sCollateCategory
}

func sCollateOrderStart(l *lexer.L) lexer.StateFunc {
	grabWhitespace(l)
	l.Take(lowerChars)
	switch l.Current() {
	case "":
		must(grabNewline, l)
		return sCollateOrder
	case "forward":
		l.Emit(FORWARD_STR)
	case "backward":
		l.Emit(BACKWARD_STR)
	case "position":
		l.Emit(POSITION_STR)
	}
	return sCollateOrderStart
}

func sCollateOrder(l *lexer.L) lexer.StateFunc {
	if !grabChar(l) {
		l.Take(lowerChars)
		if l.Current() != "order_end" {
			l.Error("what's going on here lol")
		}
		l.Emit(ORDER_END_STR)
		must(grabNewline, l)
		return sCollateCategory
	}
	must(grabNewline, l)
	return sCollateOrder
}

func getsGenericCategory(category string, tok lexer.TokenType, defs map[string]decode) lexer.StateFunc {

	var f lexer.StateFunc
	f = func(l *lexer.L) lexer.StateFunc {
		grabWhitespace(l)
		l.Take(upperChars + lowerChars + "_")
		word := l.Current()
		if word == "END" { // End the category
			l.Emit(END_STR)
			must(grabWhitespace, l)
			l.Take(upperChars)
			if l.Current() != category {
				l.Error("expected: END " + category + ", but got: END " + l.Current())
			}
			l.Emit(tok)
			grabWhitespace(l)
			must(grabNewline, l)
			return sCategory
		}

		// Handle basic keyword & option pair
		if typ, ok := defs[word]; ok {
			l.Emit(typ.typ)
			must(grabWhitespace, l)
			must(typ.grab, l)
			must(grabNewline, l)
		} else {
			l.Error(fmt.Sprintf("unexpected keyword '%s'", word))
		}
		return f
	}
	return f
}

func checkSpecialWord(line, word string, dest *rune) bool {
	wordLen := len(word)
	if len(line) >= wordLen && line[:wordLen] == word {
		words := strings.Fields(line)
		if len(words) == 2 {
			*dest, _ = utf8.DecodeRuneInString(words[1])
			return true
		}
	}
	return false
}

// Returns a string without comments
func removeComments(input string) string {
	var out strings.Builder
	lines := strings.Split(input, "\n")
	commentChar := '#'
	escapeChar = '\\' // FIXME: Uses global escapeChar
	// TODO: Resolve comment / escape-joined lines interaction
	for i := 0; i < len(lines); i++ {
		l := lines[i]

		if c, _ := utf8.DecodeRuneInString(l); c == commentChar {
			continue
		}
		// Check for escape/comment_char redefinition
		checkSpecialWord(l, "comment_char", &commentChar)
		if !checkSpecialWord(l, "escape_char", &escapeChar) {

			// Join newlines if they end with escape character
			for len(l) > 0 {
				lastRune, size := utf8.DecodeLastRuneInString(l)
				if lastRune == utf8.RuneError {
					panic(fmt.Errorf("couldn't read last rune of '%s'", l))
				}
				if lastRune != escapeChar {
					break
				}
				i++
				l = l[:len(l)-size] + strings.TrimSpace(lines[i])
			}
			_, err := out.WriteString(l)
			if err != nil {
				panic(err)
			}
			_, err = out.WriteRune('\n')
			if err != nil {
				panic(err)
			}
		}
	}
	return out.String()
}

// NewLexer creates a new lexer, ready to read from input
func NewLexer(input string) *Lexer {
	lite := removeComments(input)
	l := Lexer{
		*lexer.New(lite, sInit),
		Def{},
	}
	l.Start()
	return &l
}
