package ed

//go:generate peg ed.peg

import (
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
)

const endSymbol rune = 1114112

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	rulefirst
	rulee
	rulecmd
	rulechangeTextCmd
	ruleaddTextCmd
	rulemarkCmd
	ruledestCmd
	rulereadCmd
	rulewriteCmd
	ruleshellCmd
	rulenullCmd
	ruletext
	ruletextTerm
	rulerangeCmd
	rulerange
	ruleaddrCmd
	rulestartAddr
	ruleendAddr
	ruleaddrO
	ruleaddr
	ruleliteralAddr
	rulemarkAddr
	ruleregexAddr
	ruleregexReverseAddr
	rulebasic_regex
	ruleback_regex
	rulebareCmd
	ruleoffset
	ruleparamCmd
	ruleparamC
	ruleparam
	ruleaddrC
	rulechangeTextC
	ruleaddTextC
	rulerangeC
	ruledestC
	rulenewLine
	rulesp
	ruleAction0
	ruleAction1
	rulePegText
	ruleAction2
	ruleAction3
	ruleAction4
	ruleAction5
	ruleAction6
	ruleAction7
	ruleAction8
	ruleAction9
	ruleAction10
	ruleAction11
	ruleAction12
	ruleAction13
	ruleAction14
	ruleAction15
	ruleAction16
	ruleAction17
	ruleAction18
	ruleAction19
	ruleAction20
	ruleAction21
	ruleAction22
	ruleAction23
	ruleAction24
	ruleAction25
	ruleAction26
	ruleAction27
	ruleAction28
	ruleAction29
	ruleAction30
	ruleAction31
	ruleAction32
	ruleAction33
	ruleAction34
	ruleAction35
	ruleAction36
	ruleAction37
	ruleAction38
)

var rul3s = [...]string{
	"Unknown",
	"first",
	"e",
	"cmd",
	"changeTextCmd",
	"addTextCmd",
	"markCmd",
	"destCmd",
	"readCmd",
	"writeCmd",
	"shellCmd",
	"nullCmd",
	"text",
	"textTerm",
	"rangeCmd",
	"range",
	"addrCmd",
	"startAddr",
	"endAddr",
	"addrO",
	"addr",
	"literalAddr",
	"markAddr",
	"regexAddr",
	"regexReverseAddr",
	"basic_regex",
	"back_regex",
	"bareCmd",
	"offset",
	"paramCmd",
	"paramC",
	"param",
	"addrC",
	"changeTextC",
	"addTextC",
	"rangeC",
	"destC",
	"newLine",
	"sp",
	"Action0",
	"Action1",
	"PegText",
	"Action2",
	"Action3",
	"Action4",
	"Action5",
	"Action6",
	"Action7",
	"Action8",
	"Action9",
	"Action10",
	"Action11",
	"Action12",
	"Action13",
	"Action14",
	"Action15",
	"Action16",
	"Action17",
	"Action18",
	"Action19",
	"Action20",
	"Action21",
	"Action22",
	"Action23",
	"Action24",
	"Action25",
	"Action26",
	"Action27",
	"Action28",
	"Action29",
	"Action30",
	"Action31",
	"Action32",
	"Action33",
	"Action34",
	"Action35",
	"Action36",
	"Action37",
	"Action38",
}

type token32 struct {
	pegRule
	begin, end uint32
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v", rul3s[t.pegRule], t.begin, t.end)
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(w io.Writer, pretty bool, buffer string) {
	var print func(node *node32, depth int)
	print = func(node *node32, depth int) {
		for node != nil {
			for c := 0; c < depth; c++ {
				fmt.Fprintf(w, " ")
			}
			rule := rul3s[node.pegRule]
			quote := strconv.Quote(string(([]rune(buffer)[node.begin:node.end])))
			if !pretty {
				fmt.Fprintf(w, "%v %v\n", rule, quote)
			} else {
				fmt.Fprintf(w, "\x1B[34m%v\x1B[m %v\n", rule, quote)
			}
			if node.up != nil {
				print(node.up, depth+1)
			}
			node = node.next
		}
	}
	print(node, 0)
}

func (node *node32) Print(w io.Writer, buffer string) {
	node.print(w, false, buffer)
}

func (node *node32) PrettyPrint(w io.Writer, buffer string) {
	node.print(w, true, buffer)
}

type tokens32 struct {
	tree []token32
}

func (t *tokens32) Trim(length uint32) {
	t.tree = t.tree[:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) AST() *node32 {
	type element struct {
		node *node32
		down *element
	}
	tokens := t.Tokens()
	var stack *element
	for _, token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	if stack != nil {
		return stack.node
	}
	return nil
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	t.AST().Print(os.Stdout, buffer)
}

func (t *tokens32) WriteSyntaxTree(w io.Writer, buffer string) {
	t.AST().Print(w, buffer)
}

func (t *tokens32) PrettyPrintSyntaxTree(buffer string) {
	t.AST().PrettyPrint(os.Stdout, buffer)
}

func (t *tokens32) Add(rule pegRule, begin, end, index uint32) {
	if tree := t.tree; int(index) >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		copy(expanded, tree)
		t.tree = expanded
	}
	t.tree[index] = token32{
		pegRule: rule,
		begin:   begin,
		end:     end,
	}
}

func (t *tokens32) Tokens() []token32 {
	return t.tree
}

type Parser struct {
	Out     chan<- command
	curCmd  command
	curAddr address

	Buffer string
	buffer []rune
	rules  [79]func() bool
	parse  func(rule ...int) error
	reset  func()
	Pretty bool
	tokens32
}

func (p *Parser) Parse(rule ...int) error {
	return p.parse(rule...)
}

func (p *Parser) Reset() {
	p.reset()
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer []rune, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer {
		if c == '\n' {
			line, symbol = line+1, 0
		} else {
			symbol++
		}
		if i == positions[j] {
			translations[positions[j]] = textPosition{line, symbol}
			for j++; j < length; j++ {
				if i != positions[j] {
					continue search
				}
			}
			break search
		}
	}

	return translations
}

type parseError struct {
	p   *Parser
	max token32
}

func (e *parseError) Error() string {
	tokens, err := []token32{e.max}, "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.buffer, positions)
	format := "parse error near %v (line %v symbol %v - line %v symbol %v):\n%v\n"
	if e.p.Pretty {
		format = "parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n"
	}
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		err += fmt.Sprintf(format,
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			strconv.Quote(string(e.p.buffer[begin:end])))
	}

	return err
}

func (p *Parser) PrintSyntaxTree() {
	if p.Pretty {
		p.tokens32.PrettyPrintSyntaxTree(p.Buffer)
	} else {
		p.tokens32.PrintSyntaxTree(p.Buffer)
	}
}

func (p *Parser) WriteSyntaxTree(w io.Writer) {
	p.tokens32.WriteSyntaxTree(w, p.Buffer)
}

func (p *Parser) Execute() {
	buffer, _buffer, text, begin, end := p.Buffer, p.buffer, "", 0, 0
	for _, token := range p.Tokens() {
		switch token.pegRule {

		case rulePegText:
			begin, end = int(token.begin), int(token.end)
			text = string(_buffer[begin:end])

		case ruleAction0:
			close(p.Out)
		case ruleAction1:

			p.Out <- p.curCmd
			p.curCmd = command{}

		case ruleAction2:

			p.curCmd.typ = ctmark
			p.curCmd.params = []string{buffer[begin:end]}

		case ruleAction3:
			p.curCmd.dest = p.curAddr
			p.curAddr = address{}
		case ruleAction4:

			p.curCmd.typ = ctread
			p.curCmd.text = buffer[begin:end]

		case ruleAction5:

			p.curCmd.typ = ctwrite
			p.curCmd.text = buffer[begin:end]

		case ruleAction6:

			p.curCmd.typ = ctshell
			p.curCmd.text = buffer[begin:end]

		case ruleAction7:
			p.curCmd.text = buffer[begin:end]
		case ruleAction8:
			p.curCmd.start.text = buffer[begin:end]
		case ruleAction9:
			p.curCmd.start = p.curAddr
			p.curAddr = address{}
		case ruleAction10:
			p.curCmd.end = p.curAddr
			p.curAddr = address{}
		case ruleAction11:
			p.curAddr.text = buffer[begin:end]
		case ruleAction12:
			p.curAddr.typ = lCurrent
		case ruleAction13:
			p.curAddr.typ = lLast
		case ruleAction14:
			p.curAddr.typ = lNum
		case ruleAction15:
			p.curAddr.typ = lMark
		case ruleAction16:
			p.curAddr.typ = lRegex
		case ruleAction17:
			p.curAddr.typ = lRegexReverse
		case ruleAction18:
			p.curCmd.typ = cthelp
		case ruleAction19:
			p.curCmd.typ = cthelpMode
		case ruleAction20:
			p.curCmd.typ = ctprompt
		case ruleAction21:
			p.curCmd.typ = ctquit
		case ruleAction22:
			p.curCmd.typ = ctquitForce
		case ruleAction23:
			p.curCmd.typ = ctundo
		case ruleAction24:
			p.curCmd.params = []string{buffer[begin:end]}
		case ruleAction25:
			p.curCmd.typ = ctedit
		case ruleAction26:
			p.curCmd.typ = cteditForce
		case ruleAction27:
			p.curCmd.typ = ctfilename
		case ruleAction28:
			p.curCmd.typ = ctlineNumber
		case ruleAction29:
			p.curCmd.typ = ctchange
		case ruleAction30:
			p.curCmd.typ = ctappend
		case ruleAction31:
			p.curCmd.typ = ctinsert
		case ruleAction32:
			p.curCmd.typ = ctdelete
		case ruleAction33:
			p.curCmd.typ = ctjoin
		case ruleAction34:
			p.curCmd.typ = ctlist
		case ruleAction35:
			p.curCmd.typ = ctnumber
		case ruleAction36:
			p.curCmd.typ = ctprint
		case ruleAction37:
			p.curCmd.typ = ctmove
		case ruleAction38:
			p.curCmd.typ = ctcopy

		}
	}
	_, _, _, _, _ = buffer, _buffer, text, begin, end
}

func (p *Parser) Init() {
	var (
		max                  token32
		position, tokenIndex uint32
		buffer               []rune
	)
	p.reset = func() {
		max = token32{}
		position, tokenIndex = 0, 0

		p.buffer = []rune(p.Buffer)
		if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != endSymbol {
			p.buffer = append(p.buffer, endSymbol)
		}
		buffer = p.buffer
	}
	p.reset()

	_rules := p.rules
	tree := tokens32{tree: make([]token32, math.MaxInt16)}
	p.parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokens32 = tree
		if matches {
			p.Trim(tokenIndex)
			return nil
		}
		return &parseError{p, max}
	}

	add := func(rule pegRule, begin uint32) {
		tree.Add(rule, begin, position, tokenIndex)
		tokenIndex++
		if begin != position && position > max.end {
			max = token32{rule, begin, position}
		}
	}

	matchDot := func() bool {
		if buffer[position] != endSymbol {
			position++
			return true
		}
		return false
	}

	/*matchChar := func(c byte) bool {
		if buffer[position] == c {
			position++
			return true
		}
		return false
	}*/

	/*matchRange := func(lower byte, upper byte) bool {
		if c := buffer[position]; c >= lower && c <= upper {
			position++
			return true
		}
		return false
	}*/

	_rules = [...]func() bool{
		nil,
		/* 0 first <- <(e* !. Action0)> */
		func() bool {
			position0, tokenIndex0 := position, tokenIndex
			{
				position1 := position
			l2:
				{
					position3, tokenIndex3 := position, tokenIndex
					if !_rules[rulee]() {
						goto l3
					}
					goto l2
				l3:
					position, tokenIndex = position3, tokenIndex3
				}
				{
					position4, tokenIndex4 := position, tokenIndex
					if !matchDot() {
						goto l4
					}
					goto l0
				l4:
					position, tokenIndex = position4, tokenIndex4
				}
				if !_rules[ruleAction0]() {
					goto l0
				}
				add(rulefirst, position1)
			}
			return true
		l0:
			position, tokenIndex = position0, tokenIndex0
			return false
		},
		/* 1 e <- <(cmd newLine Action1)> */
		func() bool {
			position5, tokenIndex5 := position, tokenIndex
			{
				position6 := position
				if !_rules[rulecmd]() {
					goto l5
				}
				if !_rules[rulenewLine]() {
					goto l5
				}
				if !_rules[ruleAction1]() {
					goto l5
				}
				add(rulee, position6)
			}
			return true
		l5:
			position, tokenIndex = position5, tokenIndex5
			return false
		},
		/* 2 cmd <- <(bareCmd / paramCmd / rangeCmd / addrCmd / changeTextCmd / addTextCmd / markCmd / destCmd / readCmd / writeCmd / shellCmd / nullCmd)> */
		func() bool {
			position7, tokenIndex7 := position, tokenIndex
			{
				position8 := position
				{
					position9, tokenIndex9 := position, tokenIndex
					if !_rules[rulebareCmd]() {
						goto l10
					}
					goto l9
				l10:
					position, tokenIndex = position9, tokenIndex9
					if !_rules[ruleparamCmd]() {
						goto l11
					}
					goto l9
				l11:
					position, tokenIndex = position9, tokenIndex9
					if !_rules[rulerangeCmd]() {
						goto l12
					}
					goto l9
				l12:
					position, tokenIndex = position9, tokenIndex9
					if !_rules[ruleaddrCmd]() {
						goto l13
					}
					goto l9
				l13:
					position, tokenIndex = position9, tokenIndex9
					if !_rules[rulechangeTextCmd]() {
						goto l14
					}
					goto l9
				l14:
					position, tokenIndex = position9, tokenIndex9
					if !_rules[ruleaddTextCmd]() {
						goto l15
					}
					goto l9
				l15:
					position, tokenIndex = position9, tokenIndex9
					if !_rules[rulemarkCmd]() {
						goto l16
					}
					goto l9
				l16:
					position, tokenIndex = position9, tokenIndex9
					if !_rules[ruledestCmd]() {
						goto l17
					}
					goto l9
				l17:
					position, tokenIndex = position9, tokenIndex9
					if !_rules[rulereadCmd]() {
						goto l18
					}
					goto l9
				l18:
					position, tokenIndex = position9, tokenIndex9
					if !_rules[rulewriteCmd]() {
						goto l19
					}
					goto l9
				l19:
					position, tokenIndex = position9, tokenIndex9
					if !_rules[ruleshellCmd]() {
						goto l20
					}
					goto l9
				l20:
					position, tokenIndex = position9, tokenIndex9
					if !_rules[rulenullCmd]() {
						goto l7
					}
				}
			l9:
				add(rulecmd, position8)
			}
			return true
		l7:
			position, tokenIndex = position7, tokenIndex7
			return false
		},
		/* 3 changeTextCmd <- <(range? changeTextC '\n' text)> */
		func() bool {
			position21, tokenIndex21 := position, tokenIndex
			{
				position22 := position
				{
					position23, tokenIndex23 := position, tokenIndex
					if !_rules[rulerange]() {
						goto l23
					}
					goto l24
				l23:
					position, tokenIndex = position23, tokenIndex23
				}
			l24:
				if !_rules[rulechangeTextC]() {
					goto l21
				}
				if buffer[position] != rune('\n') {
					goto l21
				}
				position++
				if !_rules[ruletext]() {
					goto l21
				}
				add(rulechangeTextCmd, position22)
			}
			return true
		l21:
			position, tokenIndex = position21, tokenIndex21
			return false
		},
		/* 4 addTextCmd <- <(startAddr? addTextC '\n' text)> */
		func() bool {
			position25, tokenIndex25 := position, tokenIndex
			{
				position26 := position
				{
					position27, tokenIndex27 := position, tokenIndex
					if !_rules[rulestartAddr]() {
						goto l27
					}
					goto l28
				l27:
					position, tokenIndex = position27, tokenIndex27
				}
			l28:
				if !_rules[ruleaddTextC]() {
					goto l25
				}
				if buffer[position] != rune('\n') {
					goto l25
				}
				position++
				if !_rules[ruletext]() {
					goto l25
				}
				add(ruleaddTextCmd, position26)
			}
			return true
		l25:
			position, tokenIndex = position25, tokenIndex25
			return false
		},
		/* 5 markCmd <- <(startAddr? 'k' <[a-z]> Action2)> */
		func() bool {
			position29, tokenIndex29 := position, tokenIndex
			{
				position30 := position
				{
					position31, tokenIndex31 := position, tokenIndex
					if !_rules[rulestartAddr]() {
						goto l31
					}
					goto l32
				l31:
					position, tokenIndex = position31, tokenIndex31
				}
			l32:
				if buffer[position] != rune('k') {
					goto l29
				}
				position++
				{
					position33 := position
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l29
					}
					position++
					add(rulePegText, position33)
				}
				if !_rules[ruleAction2]() {
					goto l29
				}
				add(rulemarkCmd, position30)
			}
			return true
		l29:
			position, tokenIndex = position29, tokenIndex29
			return false
		},
		/* 6 destCmd <- <(range? destC <addrO> Action3)> */
		func() bool {
			position34, tokenIndex34 := position, tokenIndex
			{
				position35 := position
				{
					position36, tokenIndex36 := position, tokenIndex
					if !_rules[rulerange]() {
						goto l36
					}
					goto l37
				l36:
					position, tokenIndex = position36, tokenIndex36
				}
			l37:
				if !_rules[ruledestC]() {
					goto l34
				}
				{
					position38 := position
					if !_rules[ruleaddrO]() {
						goto l34
					}
					add(rulePegText, position38)
				}
				if !_rules[ruleAction3]() {
					goto l34
				}
				add(ruledestCmd, position35)
			}
			return true
		l34:
			position, tokenIndex = position34, tokenIndex34
			return false
		},
		/* 7 readCmd <- <(startAddr? 'r' sp <param> Action4)> */
		func() bool {
			position39, tokenIndex39 := position, tokenIndex
			{
				position40 := position
				{
					position41, tokenIndex41 := position, tokenIndex
					if !_rules[rulestartAddr]() {
						goto l41
					}
					goto l42
				l41:
					position, tokenIndex = position41, tokenIndex41
				}
			l42:
				if buffer[position] != rune('r') {
					goto l39
				}
				position++
				if !_rules[rulesp]() {
					goto l39
				}
				{
					position43 := position
					if !_rules[ruleparam]() {
						goto l39
					}
					add(rulePegText, position43)
				}
				if !_rules[ruleAction4]() {
					goto l39
				}
				add(rulereadCmd, position40)
			}
			return true
		l39:
			position, tokenIndex = position39, tokenIndex39
			return false
		},
		/* 8 writeCmd <- <(range? 'w' sp <param> Action5)> */
		func() bool {
			position44, tokenIndex44 := position, tokenIndex
			{
				position45 := position
				{
					position46, tokenIndex46 := position, tokenIndex
					if !_rules[rulerange]() {
						goto l46
					}
					goto l47
				l46:
					position, tokenIndex = position46, tokenIndex46
				}
			l47:
				if buffer[position] != rune('w') {
					goto l44
				}
				position++
				if !_rules[rulesp]() {
					goto l44
				}
				{
					position48 := position
					if !_rules[ruleparam]() {
						goto l44
					}
					add(rulePegText, position48)
				}
				if !_rules[ruleAction5]() {
					goto l44
				}
				add(rulewriteCmd, position45)
			}
			return true
		l44:
			position, tokenIndex = position44, tokenIndex44
			return false
		},
		/* 9 shellCmd <- <('!' <param> Action6)> */
		func() bool {
			position49, tokenIndex49 := position, tokenIndex
			{
				position50 := position
				if buffer[position] != rune('!') {
					goto l49
				}
				position++
				{
					position51 := position
					if !_rules[ruleparam]() {
						goto l49
					}
					add(rulePegText, position51)
				}
				if !_rules[ruleAction6]() {
					goto l49
				}
				add(ruleshellCmd, position50)
			}
			return true
		l49:
			position, tokenIndex = position49, tokenIndex49
			return false
		},
		/* 10 nullCmd <- <startAddr?> */
		func() bool {
			{
				position53 := position
				{
					position54, tokenIndex54 := position, tokenIndex
					if !_rules[rulestartAddr]() {
						goto l54
					}
					goto l55
				l54:
					position, tokenIndex = position54, tokenIndex54
				}
			l55:
				add(rulenullCmd, position53)
			}
			return true
		},
		/* 11 text <- <(<(!textTerm .)*> textTerm Action7)> */
		func() bool {
			position56, tokenIndex56 := position, tokenIndex
			{
				position57 := position
				{
					position58 := position
				l59:
					{
						position60, tokenIndex60 := position, tokenIndex
						{
							position61, tokenIndex61 := position, tokenIndex
							if !_rules[ruletextTerm]() {
								goto l61
							}
							goto l60
						l61:
							position, tokenIndex = position61, tokenIndex61
						}
						if !matchDot() {
							goto l60
						}
						goto l59
					l60:
						position, tokenIndex = position60, tokenIndex60
					}
					add(rulePegText, position58)
				}
				if !_rules[ruletextTerm]() {
					goto l56
				}
				if !_rules[ruleAction7]() {
					goto l56
				}
				add(ruletext, position57)
			}
			return true
		l56:
			position, tokenIndex = position56, tokenIndex56
			return false
		},
		/* 12 textTerm <- <('\n' '.')> */
		func() bool {
			position62, tokenIndex62 := position, tokenIndex
			{
				position63 := position
				if buffer[position] != rune('\n') {
					goto l62
				}
				position++
				if buffer[position] != rune('.') {
					goto l62
				}
				position++
				add(ruletextTerm, position63)
			}
			return true
		l62:
			position, tokenIndex = position62, tokenIndex62
			return false
		},
		/* 13 rangeCmd <- <(range? rangeC)> */
		func() bool {
			position64, tokenIndex64 := position, tokenIndex
			{
				position65 := position
				{
					position66, tokenIndex66 := position, tokenIndex
					if !_rules[rulerange]() {
						goto l66
					}
					goto l67
				l66:
					position, tokenIndex = position66, tokenIndex66
				}
			l67:
				if !_rules[rulerangeC]() {
					goto l64
				}
				add(rulerangeCmd, position65)
			}
			return true
		l64:
			position, tokenIndex = position64, tokenIndex64
			return false
		},
		/* 14 range <- <((startAddr ',' endAddr) / (startAddr ',') / (',' endAddr) / (startAddr ';' endAddr) / (startAddr ';') / (';' endAddr) / startAddr)> */
		func() bool {
			position68, tokenIndex68 := position, tokenIndex
			{
				position69 := position
				{
					position70, tokenIndex70 := position, tokenIndex
					if !_rules[rulestartAddr]() {
						goto l71
					}
					if buffer[position] != rune(',') {
						goto l71
					}
					position++
					if !_rules[ruleendAddr]() {
						goto l71
					}
					goto l70
				l71:
					position, tokenIndex = position70, tokenIndex70
					if !_rules[rulestartAddr]() {
						goto l72
					}
					if buffer[position] != rune(',') {
						goto l72
					}
					position++
					goto l70
				l72:
					position, tokenIndex = position70, tokenIndex70
					if buffer[position] != rune(',') {
						goto l73
					}
					position++
					if !_rules[ruleendAddr]() {
						goto l73
					}
					goto l70
				l73:
					position, tokenIndex = position70, tokenIndex70
					if !_rules[rulestartAddr]() {
						goto l74
					}
					if buffer[position] != rune(';') {
						goto l74
					}
					position++
					if !_rules[ruleendAddr]() {
						goto l74
					}
					goto l70
				l74:
					position, tokenIndex = position70, tokenIndex70
					if !_rules[rulestartAddr]() {
						goto l75
					}
					if buffer[position] != rune(';') {
						goto l75
					}
					position++
					goto l70
				l75:
					position, tokenIndex = position70, tokenIndex70
					if buffer[position] != rune(';') {
						goto l76
					}
					position++
					if !_rules[ruleendAddr]() {
						goto l76
					}
					goto l70
				l76:
					position, tokenIndex = position70, tokenIndex70
					if !_rules[rulestartAddr]() {
						goto l68
					}
				}
			l70:
				add(rulerange, position69)
			}
			return true
		l68:
			position, tokenIndex = position68, tokenIndex68
			return false
		},
		/* 15 addrCmd <- <((<startAddr> addrC Action8) / addrC)> */
		func() bool {
			position77, tokenIndex77 := position, tokenIndex
			{
				position78 := position
				{
					position79, tokenIndex79 := position, tokenIndex
					{
						position81 := position
						if !_rules[rulestartAddr]() {
							goto l80
						}
						add(rulePegText, position81)
					}
					if !_rules[ruleaddrC]() {
						goto l80
					}
					if !_rules[ruleAction8]() {
						goto l80
					}
					goto l79
				l80:
					position, tokenIndex = position79, tokenIndex79
					if !_rules[ruleaddrC]() {
						goto l77
					}
				}
			l79:
				add(ruleaddrCmd, position78)
			}
			return true
		l77:
			position, tokenIndex = position77, tokenIndex77
			return false
		},
		/* 16 startAddr <- <(addrO Action9)> */
		func() bool {
			position82, tokenIndex82 := position, tokenIndex
			{
				position83 := position
				if !_rules[ruleaddrO]() {
					goto l82
				}
				if !_rules[ruleAction9]() {
					goto l82
				}
				add(rulestartAddr, position83)
			}
			return true
		l82:
			position, tokenIndex = position82, tokenIndex82
			return false
		},
		/* 17 endAddr <- <(addrO Action10)> */
		func() bool {
			position84, tokenIndex84 := position, tokenIndex
			{
				position85 := position
				if !_rules[ruleaddrO]() {
					goto l84
				}
				if !_rules[ruleAction10]() {
					goto l84
				}
				add(ruleendAddr, position85)
			}
			return true
		l84:
			position, tokenIndex = position84, tokenIndex84
			return false
		},
		/* 18 addrO <- <(<(addr offset?)> Action11)> */
		func() bool {
			position86, tokenIndex86 := position, tokenIndex
			{
				position87 := position
				{
					position88 := position
					if !_rules[ruleaddr]() {
						goto l86
					}
					{
						position89, tokenIndex89 := position, tokenIndex
						if !_rules[ruleoffset]() {
							goto l89
						}
						goto l90
					l89:
						position, tokenIndex = position89, tokenIndex89
					}
				l90:
					add(rulePegText, position88)
				}
				if !_rules[ruleAction11]() {
					goto l86
				}
				add(ruleaddrO, position87)
			}
			return true
		l86:
			position, tokenIndex = position86, tokenIndex86
			return false
		},
		/* 19 addr <- <(literalAddr / markAddr / regexAddr / regexReverseAddr / ('.' Action12) / ('$' Action13))> */
		func() bool {
			position91, tokenIndex91 := position, tokenIndex
			{
				position92 := position
				{
					position93, tokenIndex93 := position, tokenIndex
					if !_rules[ruleliteralAddr]() {
						goto l94
					}
					goto l93
				l94:
					position, tokenIndex = position93, tokenIndex93
					if !_rules[rulemarkAddr]() {
						goto l95
					}
					goto l93
				l95:
					position, tokenIndex = position93, tokenIndex93
					if !_rules[ruleregexAddr]() {
						goto l96
					}
					goto l93
				l96:
					position, tokenIndex = position93, tokenIndex93
					if !_rules[ruleregexReverseAddr]() {
						goto l97
					}
					goto l93
				l97:
					position, tokenIndex = position93, tokenIndex93
					if buffer[position] != rune('.') {
						goto l98
					}
					position++
					if !_rules[ruleAction12]() {
						goto l98
					}
					goto l93
				l98:
					position, tokenIndex = position93, tokenIndex93
					if buffer[position] != rune('$') {
						goto l91
					}
					position++
					if !_rules[ruleAction13]() {
						goto l91
					}
				}
			l93:
				add(ruleaddr, position92)
			}
			return true
		l91:
			position, tokenIndex = position91, tokenIndex91
			return false
		},
		/* 20 literalAddr <- <(<[0-9]+> Action14)> */
		func() bool {
			position99, tokenIndex99 := position, tokenIndex
			{
				position100 := position
				{
					position101 := position
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l99
					}
					position++
				l102:
					{
						position103, tokenIndex103 := position, tokenIndex
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l103
						}
						position++
						goto l102
					l103:
						position, tokenIndex = position103, tokenIndex103
					}
					add(rulePegText, position101)
				}
				if !_rules[ruleAction14]() {
					goto l99
				}
				add(ruleliteralAddr, position100)
			}
			return true
		l99:
			position, tokenIndex = position99, tokenIndex99
			return false
		},
		/* 21 markAddr <- <('\'' [a-z] Action15)> */
		func() bool {
			position104, tokenIndex104 := position, tokenIndex
			{
				position105 := position
				if buffer[position] != rune('\'') {
					goto l104
				}
				position++
				if c := buffer[position]; c < rune('a') || c > rune('z') {
					goto l104
				}
				position++
				if !_rules[ruleAction15]() {
					goto l104
				}
				add(rulemarkAddr, position105)
			}
			return true
		l104:
			position, tokenIndex = position104, tokenIndex104
			return false
		},
		/* 22 regexAddr <- <('/' basic_regex '/' Action16)> */
		func() bool {
			position106, tokenIndex106 := position, tokenIndex
			{
				position107 := position
				if buffer[position] != rune('/') {
					goto l106
				}
				position++
				if !_rules[rulebasic_regex]() {
					goto l106
				}
				if buffer[position] != rune('/') {
					goto l106
				}
				position++
				if !_rules[ruleAction16]() {
					goto l106
				}
				add(ruleregexAddr, position107)
			}
			return true
		l106:
			position, tokenIndex = position106, tokenIndex106
			return false
		},
		/* 23 regexReverseAddr <- <('?' back_regex '?' Action17)> */
		func() bool {
			position108, tokenIndex108 := position, tokenIndex
			{
				position109 := position
				if buffer[position] != rune('?') {
					goto l108
				}
				position++
				if !_rules[ruleback_regex]() {
					goto l108
				}
				if buffer[position] != rune('?') {
					goto l108
				}
				position++
				if !_rules[ruleAction17]() {
					goto l108
				}
				add(ruleregexReverseAddr, position109)
			}
			return true
		l108:
			position, tokenIndex = position108, tokenIndex108
			return false
		},
		/* 24 basic_regex <- <(('\\' '/') / (!('\n' / '/') .))+> */
		func() bool {
			position110, tokenIndex110 := position, tokenIndex
			{
				position111 := position
				{
					position114, tokenIndex114 := position, tokenIndex
					if buffer[position] != rune('\\') {
						goto l115
					}
					position++
					if buffer[position] != rune('/') {
						goto l115
					}
					position++
					goto l114
				l115:
					position, tokenIndex = position114, tokenIndex114
					{
						position116, tokenIndex116 := position, tokenIndex
						{
							position117, tokenIndex117 := position, tokenIndex
							if buffer[position] != rune('\n') {
								goto l118
							}
							position++
							goto l117
						l118:
							position, tokenIndex = position117, tokenIndex117
							if buffer[position] != rune('/') {
								goto l116
							}
							position++
						}
					l117:
						goto l110
					l116:
						position, tokenIndex = position116, tokenIndex116
					}
					if !matchDot() {
						goto l110
					}
				}
			l114:
			l112:
				{
					position113, tokenIndex113 := position, tokenIndex
					{
						position119, tokenIndex119 := position, tokenIndex
						if buffer[position] != rune('\\') {
							goto l120
						}
						position++
						if buffer[position] != rune('/') {
							goto l120
						}
						position++
						goto l119
					l120:
						position, tokenIndex = position119, tokenIndex119
						{
							position121, tokenIndex121 := position, tokenIndex
							{
								position122, tokenIndex122 := position, tokenIndex
								if buffer[position] != rune('\n') {
									goto l123
								}
								position++
								goto l122
							l123:
								position, tokenIndex = position122, tokenIndex122
								if buffer[position] != rune('/') {
									goto l121
								}
								position++
							}
						l122:
							goto l113
						l121:
							position, tokenIndex = position121, tokenIndex121
						}
						if !matchDot() {
							goto l113
						}
					}
				l119:
					goto l112
				l113:
					position, tokenIndex = position113, tokenIndex113
				}
				add(rulebasic_regex, position111)
			}
			return true
		l110:
			position, tokenIndex = position110, tokenIndex110
			return false
		},
		/* 25 back_regex <- <(('\\' '?') / (!('\n' / '?') .))+> */
		func() bool {
			position124, tokenIndex124 := position, tokenIndex
			{
				position125 := position
				{
					position128, tokenIndex128 := position, tokenIndex
					if buffer[position] != rune('\\') {
						goto l129
					}
					position++
					if buffer[position] != rune('?') {
						goto l129
					}
					position++
					goto l128
				l129:
					position, tokenIndex = position128, tokenIndex128
					{
						position130, tokenIndex130 := position, tokenIndex
						{
							position131, tokenIndex131 := position, tokenIndex
							if buffer[position] != rune('\n') {
								goto l132
							}
							position++
							goto l131
						l132:
							position, tokenIndex = position131, tokenIndex131
							if buffer[position] != rune('?') {
								goto l130
							}
							position++
						}
					l131:
						goto l124
					l130:
						position, tokenIndex = position130, tokenIndex130
					}
					if !matchDot() {
						goto l124
					}
				}
			l128:
			l126:
				{
					position127, tokenIndex127 := position, tokenIndex
					{
						position133, tokenIndex133 := position, tokenIndex
						if buffer[position] != rune('\\') {
							goto l134
						}
						position++
						if buffer[position] != rune('?') {
							goto l134
						}
						position++
						goto l133
					l134:
						position, tokenIndex = position133, tokenIndex133
						{
							position135, tokenIndex135 := position, tokenIndex
							{
								position136, tokenIndex136 := position, tokenIndex
								if buffer[position] != rune('\n') {
									goto l137
								}
								position++
								goto l136
							l137:
								position, tokenIndex = position136, tokenIndex136
								if buffer[position] != rune('?') {
									goto l135
								}
								position++
							}
						l136:
							goto l127
						l135:
							position, tokenIndex = position135, tokenIndex135
						}
						if !matchDot() {
							goto l127
						}
					}
				l133:
					goto l126
				l127:
					position, tokenIndex = position127, tokenIndex127
				}
				add(ruleback_regex, position125)
			}
			return true
		l124:
			position, tokenIndex = position124, tokenIndex124
			return false
		},
		/* 26 bareCmd <- <(('h' Action18) / ('H' Action19) / ('P' Action20) / ('q' Action21) / ('Q' Action22) / ('u' Action23))> */
		func() bool {
			position138, tokenIndex138 := position, tokenIndex
			{
				position139 := position
				{
					position140, tokenIndex140 := position, tokenIndex
					if buffer[position] != rune('h') {
						goto l141
					}
					position++
					if !_rules[ruleAction18]() {
						goto l141
					}
					goto l140
				l141:
					position, tokenIndex = position140, tokenIndex140
					if buffer[position] != rune('H') {
						goto l142
					}
					position++
					if !_rules[ruleAction19]() {
						goto l142
					}
					goto l140
				l142:
					position, tokenIndex = position140, tokenIndex140
					if buffer[position] != rune('P') {
						goto l143
					}
					position++
					if !_rules[ruleAction20]() {
						goto l143
					}
					goto l140
				l143:
					position, tokenIndex = position140, tokenIndex140
					if buffer[position] != rune('q') {
						goto l144
					}
					position++
					if !_rules[ruleAction21]() {
						goto l144
					}
					goto l140
				l144:
					position, tokenIndex = position140, tokenIndex140
					if buffer[position] != rune('Q') {
						goto l145
					}
					position++
					if !_rules[ruleAction22]() {
						goto l145
					}
					goto l140
				l145:
					position, tokenIndex = position140, tokenIndex140
					if buffer[position] != rune('u') {
						goto l138
					}
					position++
					if !_rules[ruleAction23]() {
						goto l138
					}
				}
			l140:
				add(rulebareCmd, position139)
			}
			return true
		l138:
			position, tokenIndex = position138, tokenIndex138
			return false
		},
		/* 27 offset <- <(('+' / '-') [0-9]*)> */
		func() bool {
			position146, tokenIndex146 := position, tokenIndex
			{
				position147 := position
				{
					position148, tokenIndex148 := position, tokenIndex
					if buffer[position] != rune('+') {
						goto l149
					}
					position++
					goto l148
				l149:
					position, tokenIndex = position148, tokenIndex148
					if buffer[position] != rune('-') {
						goto l146
					}
					position++
				}
			l148:
			l150:
				{
					position151, tokenIndex151 := position, tokenIndex
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l151
					}
					position++
					goto l150
				l151:
					position, tokenIndex = position151, tokenIndex151
				}
				add(ruleoffset, position147)
			}
			return true
		l146:
			position, tokenIndex = position146, tokenIndex146
			return false
		},
		/* 28 paramCmd <- <((paramC sp <param> Action24) / paramC)> */
		func() bool {
			position152, tokenIndex152 := position, tokenIndex
			{
				position153 := position
				{
					position154, tokenIndex154 := position, tokenIndex
					if !_rules[ruleparamC]() {
						goto l155
					}
					if !_rules[rulesp]() {
						goto l155
					}
					{
						position156 := position
						if !_rules[ruleparam]() {
							goto l155
						}
						add(rulePegText, position156)
					}
					if !_rules[ruleAction24]() {
						goto l155
					}
					goto l154
				l155:
					position, tokenIndex = position154, tokenIndex154
					if !_rules[ruleparamC]() {
						goto l152
					}
				}
			l154:
				add(ruleparamCmd, position153)
			}
			return true
		l152:
			position, tokenIndex = position152, tokenIndex152
			return false
		},
		/* 29 paramC <- <(('e' Action25) / ('E' Action26) / ('f' Action27))> */
		func() bool {
			position157, tokenIndex157 := position, tokenIndex
			{
				position158 := position
				{
					position159, tokenIndex159 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l160
					}
					position++
					if !_rules[ruleAction25]() {
						goto l160
					}
					goto l159
				l160:
					position, tokenIndex = position159, tokenIndex159
					if buffer[position] != rune('E') {
						goto l161
					}
					position++
					if !_rules[ruleAction26]() {
						goto l161
					}
					goto l159
				l161:
					position, tokenIndex = position159, tokenIndex159
					if buffer[position] != rune('f') {
						goto l157
					}
					position++
					if !_rules[ruleAction27]() {
						goto l157
					}
				}
			l159:
				add(ruleparamC, position158)
			}
			return true
		l157:
			position, tokenIndex = position157, tokenIndex157
			return false
		},
		/* 30 param <- <(!'\n' .)+> */
		func() bool {
			position162, tokenIndex162 := position, tokenIndex
			{
				position163 := position
				{
					position166, tokenIndex166 := position, tokenIndex
					if buffer[position] != rune('\n') {
						goto l166
					}
					position++
					goto l162
				l166:
					position, tokenIndex = position166, tokenIndex166
				}
				if !matchDot() {
					goto l162
				}
			l164:
				{
					position165, tokenIndex165 := position, tokenIndex
					{
						position167, tokenIndex167 := position, tokenIndex
						if buffer[position] != rune('\n') {
							goto l167
						}
						position++
						goto l165
					l167:
						position, tokenIndex = position167, tokenIndex167
					}
					if !matchDot() {
						goto l165
					}
					goto l164
				l165:
					position, tokenIndex = position165, tokenIndex165
				}
				add(ruleparam, position163)
			}
			return true
		l162:
			position, tokenIndex = position162, tokenIndex162
			return false
		},
		/* 31 addrC <- <('=' Action28)> */
		func() bool {
			position168, tokenIndex168 := position, tokenIndex
			{
				position169 := position
				if buffer[position] != rune('=') {
					goto l168
				}
				position++
				if !_rules[ruleAction28]() {
					goto l168
				}
				add(ruleaddrC, position169)
			}
			return true
		l168:
			position, tokenIndex = position168, tokenIndex168
			return false
		},
		/* 32 changeTextC <- <('c' Action29)> */
		func() bool {
			position170, tokenIndex170 := position, tokenIndex
			{
				position171 := position
				if buffer[position] != rune('c') {
					goto l170
				}
				position++
				if !_rules[ruleAction29]() {
					goto l170
				}
				add(rulechangeTextC, position171)
			}
			return true
		l170:
			position, tokenIndex = position170, tokenIndex170
			return false
		},
		/* 33 addTextC <- <(('a' Action30) / ('i' Action31))> */
		func() bool {
			position172, tokenIndex172 := position, tokenIndex
			{
				position173 := position
				{
					position174, tokenIndex174 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l175
					}
					position++
					if !_rules[ruleAction30]() {
						goto l175
					}
					goto l174
				l175:
					position, tokenIndex = position174, tokenIndex174
					if buffer[position] != rune('i') {
						goto l172
					}
					position++
					if !_rules[ruleAction31]() {
						goto l172
					}
				}
			l174:
				add(ruleaddTextC, position173)
			}
			return true
		l172:
			position, tokenIndex = position172, tokenIndex172
			return false
		},
		/* 34 rangeC <- <(('d' Action32) / ('j' Action33) / ('l' Action34) / ('n' Action35) / ('p' Action36))> */
		func() bool {
			position176, tokenIndex176 := position, tokenIndex
			{
				position177 := position
				{
					position178, tokenIndex178 := position, tokenIndex
					if buffer[position] != rune('d') {
						goto l179
					}
					position++
					if !_rules[ruleAction32]() {
						goto l179
					}
					goto l178
				l179:
					position, tokenIndex = position178, tokenIndex178
					if buffer[position] != rune('j') {
						goto l180
					}
					position++
					if !_rules[ruleAction33]() {
						goto l180
					}
					goto l178
				l180:
					position, tokenIndex = position178, tokenIndex178
					if buffer[position] != rune('l') {
						goto l181
					}
					position++
					if !_rules[ruleAction34]() {
						goto l181
					}
					goto l178
				l181:
					position, tokenIndex = position178, tokenIndex178
					if buffer[position] != rune('n') {
						goto l182
					}
					position++
					if !_rules[ruleAction35]() {
						goto l182
					}
					goto l178
				l182:
					position, tokenIndex = position178, tokenIndex178
					if buffer[position] != rune('p') {
						goto l176
					}
					position++
					if !_rules[ruleAction36]() {
						goto l176
					}
				}
			l178:
				add(rulerangeC, position177)
			}
			return true
		l176:
			position, tokenIndex = position176, tokenIndex176
			return false
		},
		/* 35 destC <- <(('m' Action37) / ('t' Action38))> */
		func() bool {
			position183, tokenIndex183 := position, tokenIndex
			{
				position184 := position
				{
					position185, tokenIndex185 := position, tokenIndex
					if buffer[position] != rune('m') {
						goto l186
					}
					position++
					if !_rules[ruleAction37]() {
						goto l186
					}
					goto l185
				l186:
					position, tokenIndex = position185, tokenIndex185
					if buffer[position] != rune('t') {
						goto l183
					}
					position++
					if !_rules[ruleAction38]() {
						goto l183
					}
				}
			l185:
				add(ruledestC, position184)
			}
			return true
		l183:
			position, tokenIndex = position183, tokenIndex183
			return false
		},
		/* 36 newLine <- <'\n'> */
		func() bool {
			position187, tokenIndex187 := position, tokenIndex
			{
				position188 := position
				if buffer[position] != rune('\n') {
					goto l187
				}
				position++
				add(rulenewLine, position188)
			}
			return true
		l187:
			position, tokenIndex = position187, tokenIndex187
			return false
		},
		/* 37 sp <- <(' ' / '\t')+> */
		func() bool {
			position189, tokenIndex189 := position, tokenIndex
			{
				position190 := position
				{
					position193, tokenIndex193 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l194
					}
					position++
					goto l193
				l194:
					position, tokenIndex = position193, tokenIndex193
					if buffer[position] != rune('\t') {
						goto l189
					}
					position++
				}
			l193:
			l191:
				{
					position192, tokenIndex192 := position, tokenIndex
					{
						position195, tokenIndex195 := position, tokenIndex
						if buffer[position] != rune(' ') {
							goto l196
						}
						position++
						goto l195
					l196:
						position, tokenIndex = position195, tokenIndex195
						if buffer[position] != rune('\t') {
							goto l192
						}
						position++
					}
				l195:
					goto l191
				l192:
					position, tokenIndex = position192, tokenIndex192
				}
				add(rulesp, position190)
			}
			return true
		l189:
			position, tokenIndex = position189, tokenIndex189
			return false
		},
		/* 39 Action0 <- <{ close(p.Out)}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 40 Action1 <- <{
		  p.Out <- p.curCmd
		  p.curCmd = command{}
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		nil,
		/* 42 Action2 <- <{
		  p.curCmd.typ = ctmark
		  p.curCmd.params = []string{buffer[begin:end]}
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 43 Action3 <- <{p.curCmd.dest = p.curAddr ; p.curAddr = address{}}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 44 Action4 <- <{
		  p.curCmd.typ = ctread
		  p.curCmd.text = buffer[begin:end]
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 45 Action5 <- <{
		  p.curCmd.typ = ctwrite
		  p.curCmd.text = buffer[begin:end]
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 46 Action6 <- <{
		  p.curCmd.typ = ctshell
		  p.curCmd.text = buffer[begin:end]
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 47 Action7 <- <{p.curCmd.text = buffer[begin:end]}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 48 Action8 <- <{p.curCmd.start.text = buffer[begin:end]}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 49 Action9 <- <{p.curCmd.start = p.curAddr; p.curAddr = address{}}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 50 Action10 <- <{p.curCmd.end = p.curAddr; p.curAddr = address{}}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 51 Action11 <- <{p.curAddr.text = buffer[begin:end]}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 52 Action12 <- <{p.curAddr.typ = lCurrent}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 53 Action13 <- <{p.curAddr.typ = lLast}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 54 Action14 <- <{p.curAddr.typ = lNum}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 55 Action15 <- <{ p.curAddr.typ = lMark }> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 56 Action16 <- <{p.curAddr.typ = lRegex}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 57 Action17 <- <{p.curAddr.typ = lRegexReverse}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 58 Action18 <- <{p.curCmd.typ = cthelp}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 59 Action19 <- <{p.curCmd.typ = cthelpMode}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 60 Action20 <- <{p.curCmd.typ = ctprompt}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 61 Action21 <- <{p.curCmd.typ = ctquit}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 62 Action22 <- <{p.curCmd.typ = ctquitForce}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 63 Action23 <- <{p.curCmd.typ = ctundo}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 64 Action24 <- <{ p.curCmd.params = []string{buffer[begin:end]}}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 65 Action25 <- <{p.curCmd.typ = ctedit}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 66 Action26 <- <{p.curCmd.typ = cteditForce}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 67 Action27 <- <{p.curCmd.typ = ctfilename}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 68 Action28 <- <{p.curCmd.typ = ctlineNumber}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 69 Action29 <- <{p.curCmd.typ = ctchange}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 70 Action30 <- <{p.curCmd.typ = ctappend}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 71 Action31 <- <{p.curCmd.typ = ctinsert}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 72 Action32 <- <{p.curCmd.typ = ctdelete}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 73 Action33 <- <{p.curCmd.typ = ctjoin}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 74 Action34 <- <{p.curCmd.typ = ctlist}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 75 Action35 <- <{p.curCmd.typ = ctnumber}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 76 Action36 <- <{p.curCmd.typ = ctprint}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 77 Action37 <- <{p.curCmd.typ = ctmove}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 78 Action38 <- <{p.curCmd.typ = ctcopy}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
	}
	p.rules = _rules
}
