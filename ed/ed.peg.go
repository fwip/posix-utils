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
	ruleAction39
	ruleAction40
	ruleAction41
	ruleAction42
	ruleAction43
	ruleAction44
	ruleAction45
	ruleAction46
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
	"Action39",
	"Action40",
	"Action41",
	"Action42",
	"Action43",
	"Action44",
	"Action45",
	"Action46",
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
	Out     chan<- Command
	curCmd  Command
	curAddr address

	Buffer string
	buffer []rune
	rules  [87]func() bool
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

		case ruleAction1:

			p.Out <- p.curCmd
			p.curCmd = Command{}

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

			p.curCmd.typ = ctwrite

		case ruleAction7:

			p.curCmd.typ = ctshell
			p.curCmd.text = buffer[begin:end]

		case ruleAction8:
			p.curCmd.text = buffer[begin:end]
			fmt.Println("t", p.curCmd.text)
		case ruleAction9:
			p.curCmd.end = p.curCmd.start
		case ruleAction10:
			p.curCmd.start = aFirst
		case ruleAction11:
			p.curCmd.end = p.curCmd.start
		case ruleAction12:
			p.curCmd.start = aCur
		case ruleAction13:
			p.curCmd.end = p.curCmd.start
		case ruleAction14:
			p.curCmd.start = aFirst
			p.curCmd.end = aLast
		case ruleAction15:
			p.curCmd.start = aCur
			p.curCmd.end = aLast
		case ruleAction16:
			p.curCmd.start.text = buffer[begin:end]
		case ruleAction17:
			p.curCmd.start = p.curAddr
			p.curAddr = address{}
		case ruleAction18:
			p.curCmd.end = p.curAddr
			p.curAddr = address{}
		case ruleAction19:
			p.curAddr.text = buffer[begin:end]
		case ruleAction20:
			p.curAddr.typ = lCurrent
		case ruleAction21:
			p.curAddr.typ = lLast
		case ruleAction22:
			p.curAddr.typ = lNum
		case ruleAction23:
			p.curAddr.typ = lMark
		case ruleAction24:
			p.curAddr.typ = lRegex
		case ruleAction25:
			p.curAddr.typ = lRegexReverse
		case ruleAction26:
			p.curCmd.typ = cthelp
		case ruleAction27:
			p.curCmd.typ = cthelpMode
		case ruleAction28:
			p.curCmd.typ = ctprompt
		case ruleAction29:
			p.curCmd.typ = ctquit
		case ruleAction30:
			p.curCmd.typ = ctquitForce
		case ruleAction31:
			p.curCmd.typ = ctundo
		case ruleAction32:
			p.curCmd.params = []string{buffer[begin:end]}
		case ruleAction33:
			p.curCmd.typ = ctedit
		case ruleAction34:
			p.curCmd.typ = cteditForce
		case ruleAction35:
			p.curCmd.typ = ctfilename
		case ruleAction36:
			p.curCmd.typ = ctlineNumber
		case ruleAction37:
			p.curCmd.typ = ctchange
		case ruleAction38:
			p.curCmd.typ = ctappend
		case ruleAction39:
			p.curCmd.typ = ctinsert
		case ruleAction40:
			p.curCmd.typ = ctdelete
		case ruleAction41:
			p.curCmd.typ = ctjoin
		case ruleAction42:
			p.curCmd.typ = ctlist
		case ruleAction43:
			p.curCmd.typ = ctnumber
		case ruleAction44:
			p.curCmd.typ = ctprint
		case ruleAction45:
			p.curCmd.typ = ctmove
		case ruleAction46:
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
		/* 3 changeTextCmd <- <(range? changeTextC newLine text)> */
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
				if !_rules[rulenewLine]() {
					goto l21
				}
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
		/* 4 addTextCmd <- <(startAddr? addTextC newLine text)> */
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
				if !_rules[rulenewLine]() {
					goto l25
				}
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
		/* 8 writeCmd <- <((range? 'w' sp <param> Action5) / (range? 'w' Action6))> */
		func() bool {
			position44, tokenIndex44 := position, tokenIndex
			{
				position45 := position
				{
					position46, tokenIndex46 := position, tokenIndex
					{
						position48, tokenIndex48 := position, tokenIndex
						if !_rules[rulerange]() {
							goto l48
						}
						goto l49
					l48:
						position, tokenIndex = position48, tokenIndex48
					}
				l49:
					if buffer[position] != rune('w') {
						goto l47
					}
					position++
					if !_rules[rulesp]() {
						goto l47
					}
					{
						position50 := position
						if !_rules[ruleparam]() {
							goto l47
						}
						add(rulePegText, position50)
					}
					if !_rules[ruleAction5]() {
						goto l47
					}
					goto l46
				l47:
					position, tokenIndex = position46, tokenIndex46
					{
						position51, tokenIndex51 := position, tokenIndex
						if !_rules[rulerange]() {
							goto l51
						}
						goto l52
					l51:
						position, tokenIndex = position51, tokenIndex51
					}
				l52:
					if buffer[position] != rune('w') {
						goto l44
					}
					position++
					if !_rules[ruleAction6]() {
						goto l44
					}
				}
			l46:
				add(rulewriteCmd, position45)
			}
			return true
		l44:
			position, tokenIndex = position44, tokenIndex44
			return false
		},
		/* 9 shellCmd <- <('!' <param> Action7)> */
		func() bool {
			position53, tokenIndex53 := position, tokenIndex
			{
				position54 := position
				if buffer[position] != rune('!') {
					goto l53
				}
				position++
				{
					position55 := position
					if !_rules[ruleparam]() {
						goto l53
					}
					add(rulePegText, position55)
				}
				if !_rules[ruleAction7]() {
					goto l53
				}
				add(ruleshellCmd, position54)
			}
			return true
		l53:
			position, tokenIndex = position53, tokenIndex53
			return false
		},
		/* 10 nullCmd <- <startAddr?> */
		func() bool {
			{
				position57 := position
				{
					position58, tokenIndex58 := position, tokenIndex
					if !_rules[rulestartAddr]() {
						goto l58
					}
					goto l59
				l58:
					position, tokenIndex = position58, tokenIndex58
				}
			l59:
				add(rulenullCmd, position57)
			}
			return true
		},
		/* 11 text <- <(<(!textTerm .)*> textTerm Action8)> */
		func() bool {
			position60, tokenIndex60 := position, tokenIndex
			{
				position61 := position
				{
					position62 := position
				l63:
					{
						position64, tokenIndex64 := position, tokenIndex
						{
							position65, tokenIndex65 := position, tokenIndex
							if !_rules[ruletextTerm]() {
								goto l65
							}
							goto l64
						l65:
							position, tokenIndex = position65, tokenIndex65
						}
						if !matchDot() {
							goto l64
						}
						goto l63
					l64:
						position, tokenIndex = position64, tokenIndex64
					}
					add(rulePegText, position62)
				}
				if !_rules[ruletextTerm]() {
					goto l60
				}
				if !_rules[ruleAction8]() {
					goto l60
				}
				add(ruletext, position61)
			}
			return true
		l60:
			position, tokenIndex = position60, tokenIndex60
			return false
		},
		/* 12 textTerm <- <('\n' '.')> */
		func() bool {
			position66, tokenIndex66 := position, tokenIndex
			{
				position67 := position
				if buffer[position] != rune('\n') {
					goto l66
				}
				position++
				if buffer[position] != rune('.') {
					goto l66
				}
				position++
				add(ruletextTerm, position67)
			}
			return true
		l66:
			position, tokenIndex = position66, tokenIndex66
			return false
		},
		/* 13 rangeCmd <- <(range? rangeC)> */
		func() bool {
			position68, tokenIndex68 := position, tokenIndex
			{
				position69 := position
				{
					position70, tokenIndex70 := position, tokenIndex
					if !_rules[rulerange]() {
						goto l70
					}
					goto l71
				l70:
					position, tokenIndex = position70, tokenIndex70
				}
			l71:
				if !_rules[rulerangeC]() {
					goto l68
				}
				add(rulerangeCmd, position69)
			}
			return true
		l68:
			position, tokenIndex = position68, tokenIndex68
			return false
		},
		/* 14 range <- <((startAddr ',' endAddr) / (startAddr ',' Action9) / (',' endAddr Action10) / (startAddr ';' endAddr) / (startAddr ';' Action11) / (';' endAddr Action12) / (startAddr Action13) / (',' Action14) / (';' Action15))> */
		func() bool {
			position72, tokenIndex72 := position, tokenIndex
			{
				position73 := position
				{
					position74, tokenIndex74 := position, tokenIndex
					if !_rules[rulestartAddr]() {
						goto l75
					}
					if buffer[position] != rune(',') {
						goto l75
					}
					position++
					if !_rules[ruleendAddr]() {
						goto l75
					}
					goto l74
				l75:
					position, tokenIndex = position74, tokenIndex74
					if !_rules[rulestartAddr]() {
						goto l76
					}
					if buffer[position] != rune(',') {
						goto l76
					}
					position++
					if !_rules[ruleAction9]() {
						goto l76
					}
					goto l74
				l76:
					position, tokenIndex = position74, tokenIndex74
					if buffer[position] != rune(',') {
						goto l77
					}
					position++
					if !_rules[ruleendAddr]() {
						goto l77
					}
					if !_rules[ruleAction10]() {
						goto l77
					}
					goto l74
				l77:
					position, tokenIndex = position74, tokenIndex74
					if !_rules[rulestartAddr]() {
						goto l78
					}
					if buffer[position] != rune(';') {
						goto l78
					}
					position++
					if !_rules[ruleendAddr]() {
						goto l78
					}
					goto l74
				l78:
					position, tokenIndex = position74, tokenIndex74
					if !_rules[rulestartAddr]() {
						goto l79
					}
					if buffer[position] != rune(';') {
						goto l79
					}
					position++
					if !_rules[ruleAction11]() {
						goto l79
					}
					goto l74
				l79:
					position, tokenIndex = position74, tokenIndex74
					if buffer[position] != rune(';') {
						goto l80
					}
					position++
					if !_rules[ruleendAddr]() {
						goto l80
					}
					if !_rules[ruleAction12]() {
						goto l80
					}
					goto l74
				l80:
					position, tokenIndex = position74, tokenIndex74
					if !_rules[rulestartAddr]() {
						goto l81
					}
					if !_rules[ruleAction13]() {
						goto l81
					}
					goto l74
				l81:
					position, tokenIndex = position74, tokenIndex74
					if buffer[position] != rune(',') {
						goto l82
					}
					position++
					if !_rules[ruleAction14]() {
						goto l82
					}
					goto l74
				l82:
					position, tokenIndex = position74, tokenIndex74
					if buffer[position] != rune(';') {
						goto l72
					}
					position++
					if !_rules[ruleAction15]() {
						goto l72
					}
				}
			l74:
				add(rulerange, position73)
			}
			return true
		l72:
			position, tokenIndex = position72, tokenIndex72
			return false
		},
		/* 15 addrCmd <- <((<startAddr> addrC Action16) / addrC)> */
		func() bool {
			position83, tokenIndex83 := position, tokenIndex
			{
				position84 := position
				{
					position85, tokenIndex85 := position, tokenIndex
					{
						position87 := position
						if !_rules[rulestartAddr]() {
							goto l86
						}
						add(rulePegText, position87)
					}
					if !_rules[ruleaddrC]() {
						goto l86
					}
					if !_rules[ruleAction16]() {
						goto l86
					}
					goto l85
				l86:
					position, tokenIndex = position85, tokenIndex85
					if !_rules[ruleaddrC]() {
						goto l83
					}
				}
			l85:
				add(ruleaddrCmd, position84)
			}
			return true
		l83:
			position, tokenIndex = position83, tokenIndex83
			return false
		},
		/* 16 startAddr <- <(addrO Action17)> */
		func() bool {
			position88, tokenIndex88 := position, tokenIndex
			{
				position89 := position
				if !_rules[ruleaddrO]() {
					goto l88
				}
				if !_rules[ruleAction17]() {
					goto l88
				}
				add(rulestartAddr, position89)
			}
			return true
		l88:
			position, tokenIndex = position88, tokenIndex88
			return false
		},
		/* 17 endAddr <- <(addrO Action18)> */
		func() bool {
			position90, tokenIndex90 := position, tokenIndex
			{
				position91 := position
				if !_rules[ruleaddrO]() {
					goto l90
				}
				if !_rules[ruleAction18]() {
					goto l90
				}
				add(ruleendAddr, position91)
			}
			return true
		l90:
			position, tokenIndex = position90, tokenIndex90
			return false
		},
		/* 18 addrO <- <(<(addr offset?)> Action19)> */
		func() bool {
			position92, tokenIndex92 := position, tokenIndex
			{
				position93 := position
				{
					position94 := position
					if !_rules[ruleaddr]() {
						goto l92
					}
					{
						position95, tokenIndex95 := position, tokenIndex
						if !_rules[ruleoffset]() {
							goto l95
						}
						goto l96
					l95:
						position, tokenIndex = position95, tokenIndex95
					}
				l96:
					add(rulePegText, position94)
				}
				if !_rules[ruleAction19]() {
					goto l92
				}
				add(ruleaddrO, position93)
			}
			return true
		l92:
			position, tokenIndex = position92, tokenIndex92
			return false
		},
		/* 19 addr <- <(literalAddr / markAddr / regexAddr / regexReverseAddr / ('.' Action20) / ('$' Action21))> */
		func() bool {
			position97, tokenIndex97 := position, tokenIndex
			{
				position98 := position
				{
					position99, tokenIndex99 := position, tokenIndex
					if !_rules[ruleliteralAddr]() {
						goto l100
					}
					goto l99
				l100:
					position, tokenIndex = position99, tokenIndex99
					if !_rules[rulemarkAddr]() {
						goto l101
					}
					goto l99
				l101:
					position, tokenIndex = position99, tokenIndex99
					if !_rules[ruleregexAddr]() {
						goto l102
					}
					goto l99
				l102:
					position, tokenIndex = position99, tokenIndex99
					if !_rules[ruleregexReverseAddr]() {
						goto l103
					}
					goto l99
				l103:
					position, tokenIndex = position99, tokenIndex99
					if buffer[position] != rune('.') {
						goto l104
					}
					position++
					if !_rules[ruleAction20]() {
						goto l104
					}
					goto l99
				l104:
					position, tokenIndex = position99, tokenIndex99
					if buffer[position] != rune('$') {
						goto l97
					}
					position++
					if !_rules[ruleAction21]() {
						goto l97
					}
				}
			l99:
				add(ruleaddr, position98)
			}
			return true
		l97:
			position, tokenIndex = position97, tokenIndex97
			return false
		},
		/* 20 literalAddr <- <(<[0-9]+> Action22)> */
		func() bool {
			position105, tokenIndex105 := position, tokenIndex
			{
				position106 := position
				{
					position107 := position
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l105
					}
					position++
				l108:
					{
						position109, tokenIndex109 := position, tokenIndex
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l109
						}
						position++
						goto l108
					l109:
						position, tokenIndex = position109, tokenIndex109
					}
					add(rulePegText, position107)
				}
				if !_rules[ruleAction22]() {
					goto l105
				}
				add(ruleliteralAddr, position106)
			}
			return true
		l105:
			position, tokenIndex = position105, tokenIndex105
			return false
		},
		/* 21 markAddr <- <('\'' [a-z] Action23)> */
		func() bool {
			position110, tokenIndex110 := position, tokenIndex
			{
				position111 := position
				if buffer[position] != rune('\'') {
					goto l110
				}
				position++
				if c := buffer[position]; c < rune('a') || c > rune('z') {
					goto l110
				}
				position++
				if !_rules[ruleAction23]() {
					goto l110
				}
				add(rulemarkAddr, position111)
			}
			return true
		l110:
			position, tokenIndex = position110, tokenIndex110
			return false
		},
		/* 22 regexAddr <- <('/' basic_regex '/' Action24)> */
		func() bool {
			position112, tokenIndex112 := position, tokenIndex
			{
				position113 := position
				if buffer[position] != rune('/') {
					goto l112
				}
				position++
				if !_rules[rulebasic_regex]() {
					goto l112
				}
				if buffer[position] != rune('/') {
					goto l112
				}
				position++
				if !_rules[ruleAction24]() {
					goto l112
				}
				add(ruleregexAddr, position113)
			}
			return true
		l112:
			position, tokenIndex = position112, tokenIndex112
			return false
		},
		/* 23 regexReverseAddr <- <('?' back_regex '?' Action25)> */
		func() bool {
			position114, tokenIndex114 := position, tokenIndex
			{
				position115 := position
				if buffer[position] != rune('?') {
					goto l114
				}
				position++
				if !_rules[ruleback_regex]() {
					goto l114
				}
				if buffer[position] != rune('?') {
					goto l114
				}
				position++
				if !_rules[ruleAction25]() {
					goto l114
				}
				add(ruleregexReverseAddr, position115)
			}
			return true
		l114:
			position, tokenIndex = position114, tokenIndex114
			return false
		},
		/* 24 basic_regex <- <(('\\' '/') / (!('\n' / '/') .))+> */
		func() bool {
			position116, tokenIndex116 := position, tokenIndex
			{
				position117 := position
				{
					position120, tokenIndex120 := position, tokenIndex
					if buffer[position] != rune('\\') {
						goto l121
					}
					position++
					if buffer[position] != rune('/') {
						goto l121
					}
					position++
					goto l120
				l121:
					position, tokenIndex = position120, tokenIndex120
					{
						position122, tokenIndex122 := position, tokenIndex
						{
							position123, tokenIndex123 := position, tokenIndex
							if buffer[position] != rune('\n') {
								goto l124
							}
							position++
							goto l123
						l124:
							position, tokenIndex = position123, tokenIndex123
							if buffer[position] != rune('/') {
								goto l122
							}
							position++
						}
					l123:
						goto l116
					l122:
						position, tokenIndex = position122, tokenIndex122
					}
					if !matchDot() {
						goto l116
					}
				}
			l120:
			l118:
				{
					position119, tokenIndex119 := position, tokenIndex
					{
						position125, tokenIndex125 := position, tokenIndex
						if buffer[position] != rune('\\') {
							goto l126
						}
						position++
						if buffer[position] != rune('/') {
							goto l126
						}
						position++
						goto l125
					l126:
						position, tokenIndex = position125, tokenIndex125
						{
							position127, tokenIndex127 := position, tokenIndex
							{
								position128, tokenIndex128 := position, tokenIndex
								if buffer[position] != rune('\n') {
									goto l129
								}
								position++
								goto l128
							l129:
								position, tokenIndex = position128, tokenIndex128
								if buffer[position] != rune('/') {
									goto l127
								}
								position++
							}
						l128:
							goto l119
						l127:
							position, tokenIndex = position127, tokenIndex127
						}
						if !matchDot() {
							goto l119
						}
					}
				l125:
					goto l118
				l119:
					position, tokenIndex = position119, tokenIndex119
				}
				add(rulebasic_regex, position117)
			}
			return true
		l116:
			position, tokenIndex = position116, tokenIndex116
			return false
		},
		/* 25 back_regex <- <(('\\' '?') / (!('\n' / '?') .))+> */
		func() bool {
			position130, tokenIndex130 := position, tokenIndex
			{
				position131 := position
				{
					position134, tokenIndex134 := position, tokenIndex
					if buffer[position] != rune('\\') {
						goto l135
					}
					position++
					if buffer[position] != rune('?') {
						goto l135
					}
					position++
					goto l134
				l135:
					position, tokenIndex = position134, tokenIndex134
					{
						position136, tokenIndex136 := position, tokenIndex
						{
							position137, tokenIndex137 := position, tokenIndex
							if buffer[position] != rune('\n') {
								goto l138
							}
							position++
							goto l137
						l138:
							position, tokenIndex = position137, tokenIndex137
							if buffer[position] != rune('?') {
								goto l136
							}
							position++
						}
					l137:
						goto l130
					l136:
						position, tokenIndex = position136, tokenIndex136
					}
					if !matchDot() {
						goto l130
					}
				}
			l134:
			l132:
				{
					position133, tokenIndex133 := position, tokenIndex
					{
						position139, tokenIndex139 := position, tokenIndex
						if buffer[position] != rune('\\') {
							goto l140
						}
						position++
						if buffer[position] != rune('?') {
							goto l140
						}
						position++
						goto l139
					l140:
						position, tokenIndex = position139, tokenIndex139
						{
							position141, tokenIndex141 := position, tokenIndex
							{
								position142, tokenIndex142 := position, tokenIndex
								if buffer[position] != rune('\n') {
									goto l143
								}
								position++
								goto l142
							l143:
								position, tokenIndex = position142, tokenIndex142
								if buffer[position] != rune('?') {
									goto l141
								}
								position++
							}
						l142:
							goto l133
						l141:
							position, tokenIndex = position141, tokenIndex141
						}
						if !matchDot() {
							goto l133
						}
					}
				l139:
					goto l132
				l133:
					position, tokenIndex = position133, tokenIndex133
				}
				add(ruleback_regex, position131)
			}
			return true
		l130:
			position, tokenIndex = position130, tokenIndex130
			return false
		},
		/* 26 bareCmd <- <(('h' Action26) / ('H' Action27) / ('P' Action28) / ('q' Action29) / ('Q' Action30) / ('u' Action31))> */
		func() bool {
			position144, tokenIndex144 := position, tokenIndex
			{
				position145 := position
				{
					position146, tokenIndex146 := position, tokenIndex
					if buffer[position] != rune('h') {
						goto l147
					}
					position++
					if !_rules[ruleAction26]() {
						goto l147
					}
					goto l146
				l147:
					position, tokenIndex = position146, tokenIndex146
					if buffer[position] != rune('H') {
						goto l148
					}
					position++
					if !_rules[ruleAction27]() {
						goto l148
					}
					goto l146
				l148:
					position, tokenIndex = position146, tokenIndex146
					if buffer[position] != rune('P') {
						goto l149
					}
					position++
					if !_rules[ruleAction28]() {
						goto l149
					}
					goto l146
				l149:
					position, tokenIndex = position146, tokenIndex146
					if buffer[position] != rune('q') {
						goto l150
					}
					position++
					if !_rules[ruleAction29]() {
						goto l150
					}
					goto l146
				l150:
					position, tokenIndex = position146, tokenIndex146
					if buffer[position] != rune('Q') {
						goto l151
					}
					position++
					if !_rules[ruleAction30]() {
						goto l151
					}
					goto l146
				l151:
					position, tokenIndex = position146, tokenIndex146
					if buffer[position] != rune('u') {
						goto l144
					}
					position++
					if !_rules[ruleAction31]() {
						goto l144
					}
				}
			l146:
				add(rulebareCmd, position145)
			}
			return true
		l144:
			position, tokenIndex = position144, tokenIndex144
			return false
		},
		/* 27 offset <- <(('+' / '-') [0-9]*)> */
		func() bool {
			position152, tokenIndex152 := position, tokenIndex
			{
				position153 := position
				{
					position154, tokenIndex154 := position, tokenIndex
					if buffer[position] != rune('+') {
						goto l155
					}
					position++
					goto l154
				l155:
					position, tokenIndex = position154, tokenIndex154
					if buffer[position] != rune('-') {
						goto l152
					}
					position++
				}
			l154:
			l156:
				{
					position157, tokenIndex157 := position, tokenIndex
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l157
					}
					position++
					goto l156
				l157:
					position, tokenIndex = position157, tokenIndex157
				}
				add(ruleoffset, position153)
			}
			return true
		l152:
			position, tokenIndex = position152, tokenIndex152
			return false
		},
		/* 28 paramCmd <- <((paramC sp <param> Action32) / paramC)> */
		func() bool {
			position158, tokenIndex158 := position, tokenIndex
			{
				position159 := position
				{
					position160, tokenIndex160 := position, tokenIndex
					if !_rules[ruleparamC]() {
						goto l161
					}
					if !_rules[rulesp]() {
						goto l161
					}
					{
						position162 := position
						if !_rules[ruleparam]() {
							goto l161
						}
						add(rulePegText, position162)
					}
					if !_rules[ruleAction32]() {
						goto l161
					}
					goto l160
				l161:
					position, tokenIndex = position160, tokenIndex160
					if !_rules[ruleparamC]() {
						goto l158
					}
				}
			l160:
				add(ruleparamCmd, position159)
			}
			return true
		l158:
			position, tokenIndex = position158, tokenIndex158
			return false
		},
		/* 29 paramC <- <(('e' Action33) / ('E' Action34) / ('f' Action35))> */
		func() bool {
			position163, tokenIndex163 := position, tokenIndex
			{
				position164 := position
				{
					position165, tokenIndex165 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l166
					}
					position++
					if !_rules[ruleAction33]() {
						goto l166
					}
					goto l165
				l166:
					position, tokenIndex = position165, tokenIndex165
					if buffer[position] != rune('E') {
						goto l167
					}
					position++
					if !_rules[ruleAction34]() {
						goto l167
					}
					goto l165
				l167:
					position, tokenIndex = position165, tokenIndex165
					if buffer[position] != rune('f') {
						goto l163
					}
					position++
					if !_rules[ruleAction35]() {
						goto l163
					}
				}
			l165:
				add(ruleparamC, position164)
			}
			return true
		l163:
			position, tokenIndex = position163, tokenIndex163
			return false
		},
		/* 30 param <- <(!'\n' .)+> */
		func() bool {
			position168, tokenIndex168 := position, tokenIndex
			{
				position169 := position
				{
					position172, tokenIndex172 := position, tokenIndex
					if buffer[position] != rune('\n') {
						goto l172
					}
					position++
					goto l168
				l172:
					position, tokenIndex = position172, tokenIndex172
				}
				if !matchDot() {
					goto l168
				}
			l170:
				{
					position171, tokenIndex171 := position, tokenIndex
					{
						position173, tokenIndex173 := position, tokenIndex
						if buffer[position] != rune('\n') {
							goto l173
						}
						position++
						goto l171
					l173:
						position, tokenIndex = position173, tokenIndex173
					}
					if !matchDot() {
						goto l171
					}
					goto l170
				l171:
					position, tokenIndex = position171, tokenIndex171
				}
				add(ruleparam, position169)
			}
			return true
		l168:
			position, tokenIndex = position168, tokenIndex168
			return false
		},
		/* 31 addrC <- <('=' Action36)> */
		func() bool {
			position174, tokenIndex174 := position, tokenIndex
			{
				position175 := position
				if buffer[position] != rune('=') {
					goto l174
				}
				position++
				if !_rules[ruleAction36]() {
					goto l174
				}
				add(ruleaddrC, position175)
			}
			return true
		l174:
			position, tokenIndex = position174, tokenIndex174
			return false
		},
		/* 32 changeTextC <- <('c' Action37)> */
		func() bool {
			position176, tokenIndex176 := position, tokenIndex
			{
				position177 := position
				if buffer[position] != rune('c') {
					goto l176
				}
				position++
				if !_rules[ruleAction37]() {
					goto l176
				}
				add(rulechangeTextC, position177)
			}
			return true
		l176:
			position, tokenIndex = position176, tokenIndex176
			return false
		},
		/* 33 addTextC <- <(('a' Action38) / ('i' Action39))> */
		func() bool {
			position178, tokenIndex178 := position, tokenIndex
			{
				position179 := position
				{
					position180, tokenIndex180 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l181
					}
					position++
					if !_rules[ruleAction38]() {
						goto l181
					}
					goto l180
				l181:
					position, tokenIndex = position180, tokenIndex180
					if buffer[position] != rune('i') {
						goto l178
					}
					position++
					if !_rules[ruleAction39]() {
						goto l178
					}
				}
			l180:
				add(ruleaddTextC, position179)
			}
			return true
		l178:
			position, tokenIndex = position178, tokenIndex178
			return false
		},
		/* 34 rangeC <- <(('d' Action40) / ('j' Action41) / ('l' Action42) / ('n' Action43) / ('p' Action44))> */
		func() bool {
			position182, tokenIndex182 := position, tokenIndex
			{
				position183 := position
				{
					position184, tokenIndex184 := position, tokenIndex
					if buffer[position] != rune('d') {
						goto l185
					}
					position++
					if !_rules[ruleAction40]() {
						goto l185
					}
					goto l184
				l185:
					position, tokenIndex = position184, tokenIndex184
					if buffer[position] != rune('j') {
						goto l186
					}
					position++
					if !_rules[ruleAction41]() {
						goto l186
					}
					goto l184
				l186:
					position, tokenIndex = position184, tokenIndex184
					if buffer[position] != rune('l') {
						goto l187
					}
					position++
					if !_rules[ruleAction42]() {
						goto l187
					}
					goto l184
				l187:
					position, tokenIndex = position184, tokenIndex184
					if buffer[position] != rune('n') {
						goto l188
					}
					position++
					if !_rules[ruleAction43]() {
						goto l188
					}
					goto l184
				l188:
					position, tokenIndex = position184, tokenIndex184
					if buffer[position] != rune('p') {
						goto l182
					}
					position++
					if !_rules[ruleAction44]() {
						goto l182
					}
				}
			l184:
				add(rulerangeC, position183)
			}
			return true
		l182:
			position, tokenIndex = position182, tokenIndex182
			return false
		},
		/* 35 destC <- <(('m' Action45) / ('t' Action46))> */
		func() bool {
			position189, tokenIndex189 := position, tokenIndex
			{
				position190 := position
				{
					position191, tokenIndex191 := position, tokenIndex
					if buffer[position] != rune('m') {
						goto l192
					}
					position++
					if !_rules[ruleAction45]() {
						goto l192
					}
					goto l191
				l192:
					position, tokenIndex = position191, tokenIndex191
					if buffer[position] != rune('t') {
						goto l189
					}
					position++
					if !_rules[ruleAction46]() {
						goto l189
					}
				}
			l191:
				add(ruledestC, position190)
			}
			return true
		l189:
			position, tokenIndex = position189, tokenIndex189
			return false
		},
		/* 36 newLine <- <'\n'> */
		func() bool {
			position193, tokenIndex193 := position, tokenIndex
			{
				position194 := position
				if buffer[position] != rune('\n') {
					goto l193
				}
				position++
				add(rulenewLine, position194)
			}
			return true
		l193:
			position, tokenIndex = position193, tokenIndex193
			return false
		},
		/* 37 sp <- <(' ' / '\t')+> */
		func() bool {
			position195, tokenIndex195 := position, tokenIndex
			{
				position196 := position
				{
					position199, tokenIndex199 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l200
					}
					position++
					goto l199
				l200:
					position, tokenIndex = position199, tokenIndex199
					if buffer[position] != rune('\t') {
						goto l195
					}
					position++
				}
			l199:
			l197:
				{
					position198, tokenIndex198 := position, tokenIndex
					{
						position201, tokenIndex201 := position, tokenIndex
						if buffer[position] != rune(' ') {
							goto l202
						}
						position++
						goto l201
					l202:
						position, tokenIndex = position201, tokenIndex201
						if buffer[position] != rune('\t') {
							goto l198
						}
						position++
					}
				l201:
					goto l197
				l198:
					position, tokenIndex = position198, tokenIndex198
				}
				add(rulesp, position196)
			}
			return true
		l195:
			position, tokenIndex = position195, tokenIndex195
			return false
		},
		/* 39 Action0 <- <{ }> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 40 Action1 <- <{
		  p.Out <- p.curCmd
		  p.curCmd = Command{}
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
		  p.curCmd.typ = ctwrite
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 47 Action7 <- <{
		  p.curCmd.typ = ctshell
		  p.curCmd.text = buffer[begin:end]
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 48 Action8 <- <{p.curCmd.text = buffer[begin:end]; fmt.Println("t", p.curCmd.text)}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 49 Action9 <- <{p.curCmd.end = p.curCmd.start}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 50 Action10 <- <{p.curCmd.start = aFirst}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 51 Action11 <- <{p.curCmd.end = p.curCmd.start}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 52 Action12 <- <{p.curCmd.start = aCur}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 53 Action13 <- <{p.curCmd.end = p.curCmd.start}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 54 Action14 <- <{p.curCmd.start = aFirst; p.curCmd.end = aLast}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 55 Action15 <- <{p.curCmd.start = aCur; p.curCmd.end = aLast}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 56 Action16 <- <{p.curCmd.start.text = buffer[begin:end]}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 57 Action17 <- <{p.curCmd.start = p.curAddr; p.curAddr = address{}}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 58 Action18 <- <{p.curCmd.end = p.curAddr; p.curAddr = address{}}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 59 Action19 <- <{p.curAddr.text = buffer[begin:end]}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 60 Action20 <- <{p.curAddr.typ = lCurrent}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 61 Action21 <- <{p.curAddr.typ = lLast}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 62 Action22 <- <{p.curAddr.typ = lNum}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 63 Action23 <- <{ p.curAddr.typ = lMark }> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 64 Action24 <- <{p.curAddr.typ = lRegex}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 65 Action25 <- <{p.curAddr.typ = lRegexReverse}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 66 Action26 <- <{p.curCmd.typ = cthelp}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 67 Action27 <- <{p.curCmd.typ = cthelpMode}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 68 Action28 <- <{p.curCmd.typ = ctprompt}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 69 Action29 <- <{p.curCmd.typ = ctquit}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 70 Action30 <- <{p.curCmd.typ = ctquitForce}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 71 Action31 <- <{p.curCmd.typ = ctundo}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 72 Action32 <- <{ p.curCmd.params = []string{buffer[begin:end]}}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 73 Action33 <- <{p.curCmd.typ = ctedit}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 74 Action34 <- <{p.curCmd.typ = cteditForce}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 75 Action35 <- <{p.curCmd.typ = ctfilename}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 76 Action36 <- <{p.curCmd.typ = ctlineNumber}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 77 Action37 <- <{p.curCmd.typ = ctchange}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 78 Action38 <- <{p.curCmd.typ = ctappend}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 79 Action39 <- <{p.curCmd.typ = ctinsert}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 80 Action40 <- <{p.curCmd.typ = ctdelete}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 81 Action41 <- <{p.curCmd.typ = ctjoin}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 82 Action42 <- <{p.curCmd.typ = ctlist}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 83 Action43 <- <{p.curCmd.typ = ctnumber}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 84 Action44 <- <{p.curCmd.typ = ctprint}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 85 Action45 <- <{p.curCmd.typ = ctmove}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 86 Action46 <- <{p.curCmd.typ = ctcopy}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
	}
	p.rules = _rules
}
