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
		/* 1 e <- <(sp* cmd sp* newLine Action1)> */
		func() bool {
			position5, tokenIndex5 := position, tokenIndex
			{
				position6 := position
			l7:
				{
					position8, tokenIndex8 := position, tokenIndex
					if !_rules[rulesp]() {
						goto l8
					}
					goto l7
				l8:
					position, tokenIndex = position8, tokenIndex8
				}
				if !_rules[rulecmd]() {
					goto l5
				}
			l9:
				{
					position10, tokenIndex10 := position, tokenIndex
					if !_rules[rulesp]() {
						goto l10
					}
					goto l9
				l10:
					position, tokenIndex = position10, tokenIndex10
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
			position11, tokenIndex11 := position, tokenIndex
			{
				position12 := position
				{
					position13, tokenIndex13 := position, tokenIndex
					if !_rules[rulebareCmd]() {
						goto l14
					}
					goto l13
				l14:
					position, tokenIndex = position13, tokenIndex13
					if !_rules[ruleparamCmd]() {
						goto l15
					}
					goto l13
				l15:
					position, tokenIndex = position13, tokenIndex13
					if !_rules[rulerangeCmd]() {
						goto l16
					}
					goto l13
				l16:
					position, tokenIndex = position13, tokenIndex13
					if !_rules[ruleaddrCmd]() {
						goto l17
					}
					goto l13
				l17:
					position, tokenIndex = position13, tokenIndex13
					if !_rules[rulechangeTextCmd]() {
						goto l18
					}
					goto l13
				l18:
					position, tokenIndex = position13, tokenIndex13
					if !_rules[ruleaddTextCmd]() {
						goto l19
					}
					goto l13
				l19:
					position, tokenIndex = position13, tokenIndex13
					if !_rules[rulemarkCmd]() {
						goto l20
					}
					goto l13
				l20:
					position, tokenIndex = position13, tokenIndex13
					if !_rules[ruledestCmd]() {
						goto l21
					}
					goto l13
				l21:
					position, tokenIndex = position13, tokenIndex13
					if !_rules[rulereadCmd]() {
						goto l22
					}
					goto l13
				l22:
					position, tokenIndex = position13, tokenIndex13
					if !_rules[rulewriteCmd]() {
						goto l23
					}
					goto l13
				l23:
					position, tokenIndex = position13, tokenIndex13
					if !_rules[ruleshellCmd]() {
						goto l24
					}
					goto l13
				l24:
					position, tokenIndex = position13, tokenIndex13
					if !_rules[rulenullCmd]() {
						goto l11
					}
				}
			l13:
				add(rulecmd, position12)
			}
			return true
		l11:
			position, tokenIndex = position11, tokenIndex11
			return false
		},
		/* 3 changeTextCmd <- <(range? changeTextC newLine text)> */
		func() bool {
			position25, tokenIndex25 := position, tokenIndex
			{
				position26 := position
				{
					position27, tokenIndex27 := position, tokenIndex
					if !_rules[rulerange]() {
						goto l27
					}
					goto l28
				l27:
					position, tokenIndex = position27, tokenIndex27
				}
			l28:
				if !_rules[rulechangeTextC]() {
					goto l25
				}
				if !_rules[rulenewLine]() {
					goto l25
				}
				if !_rules[ruletext]() {
					goto l25
				}
				add(rulechangeTextCmd, position26)
			}
			return true
		l25:
			position, tokenIndex = position25, tokenIndex25
			return false
		},
		/* 4 addTextCmd <- <(startAddr? addTextC newLine text)> */
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
				if !_rules[ruleaddTextC]() {
					goto l29
				}
				if !_rules[rulenewLine]() {
					goto l29
				}
				if !_rules[ruletext]() {
					goto l29
				}
				add(ruleaddTextCmd, position30)
			}
			return true
		l29:
			position, tokenIndex = position29, tokenIndex29
			return false
		},
		/* 5 markCmd <- <(startAddr? 'k' <[a-z]> Action2)> */
		func() bool {
			position33, tokenIndex33 := position, tokenIndex
			{
				position34 := position
				{
					position35, tokenIndex35 := position, tokenIndex
					if !_rules[rulestartAddr]() {
						goto l35
					}
					goto l36
				l35:
					position, tokenIndex = position35, tokenIndex35
				}
			l36:
				if buffer[position] != rune('k') {
					goto l33
				}
				position++
				{
					position37 := position
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l33
					}
					position++
					add(rulePegText, position37)
				}
				if !_rules[ruleAction2]() {
					goto l33
				}
				add(rulemarkCmd, position34)
			}
			return true
		l33:
			position, tokenIndex = position33, tokenIndex33
			return false
		},
		/* 6 destCmd <- <(range? destC <addrO> Action3)> */
		func() bool {
			position38, tokenIndex38 := position, tokenIndex
			{
				position39 := position
				{
					position40, tokenIndex40 := position, tokenIndex
					if !_rules[rulerange]() {
						goto l40
					}
					goto l41
				l40:
					position, tokenIndex = position40, tokenIndex40
				}
			l41:
				if !_rules[ruledestC]() {
					goto l38
				}
				{
					position42 := position
					if !_rules[ruleaddrO]() {
						goto l38
					}
					add(rulePegText, position42)
				}
				if !_rules[ruleAction3]() {
					goto l38
				}
				add(ruledestCmd, position39)
			}
			return true
		l38:
			position, tokenIndex = position38, tokenIndex38
			return false
		},
		/* 7 readCmd <- <(startAddr? 'r' sp <param> Action4)> */
		func() bool {
			position43, tokenIndex43 := position, tokenIndex
			{
				position44 := position
				{
					position45, tokenIndex45 := position, tokenIndex
					if !_rules[rulestartAddr]() {
						goto l45
					}
					goto l46
				l45:
					position, tokenIndex = position45, tokenIndex45
				}
			l46:
				if buffer[position] != rune('r') {
					goto l43
				}
				position++
				if !_rules[rulesp]() {
					goto l43
				}
				{
					position47 := position
					if !_rules[ruleparam]() {
						goto l43
					}
					add(rulePegText, position47)
				}
				if !_rules[ruleAction4]() {
					goto l43
				}
				add(rulereadCmd, position44)
			}
			return true
		l43:
			position, tokenIndex = position43, tokenIndex43
			return false
		},
		/* 8 writeCmd <- <((range? 'w' sp <param> Action5) / (range? 'w' Action6))> */
		func() bool {
			position48, tokenIndex48 := position, tokenIndex
			{
				position49 := position
				{
					position50, tokenIndex50 := position, tokenIndex
					{
						position52, tokenIndex52 := position, tokenIndex
						if !_rules[rulerange]() {
							goto l52
						}
						goto l53
					l52:
						position, tokenIndex = position52, tokenIndex52
					}
				l53:
					if buffer[position] != rune('w') {
						goto l51
					}
					position++
					if !_rules[rulesp]() {
						goto l51
					}
					{
						position54 := position
						if !_rules[ruleparam]() {
							goto l51
						}
						add(rulePegText, position54)
					}
					if !_rules[ruleAction5]() {
						goto l51
					}
					goto l50
				l51:
					position, tokenIndex = position50, tokenIndex50
					{
						position55, tokenIndex55 := position, tokenIndex
						if !_rules[rulerange]() {
							goto l55
						}
						goto l56
					l55:
						position, tokenIndex = position55, tokenIndex55
					}
				l56:
					if buffer[position] != rune('w') {
						goto l48
					}
					position++
					if !_rules[ruleAction6]() {
						goto l48
					}
				}
			l50:
				add(rulewriteCmd, position49)
			}
			return true
		l48:
			position, tokenIndex = position48, tokenIndex48
			return false
		},
		/* 9 shellCmd <- <('!' <param> Action7)> */
		func() bool {
			position57, tokenIndex57 := position, tokenIndex
			{
				position58 := position
				if buffer[position] != rune('!') {
					goto l57
				}
				position++
				{
					position59 := position
					if !_rules[ruleparam]() {
						goto l57
					}
					add(rulePegText, position59)
				}
				if !_rules[ruleAction7]() {
					goto l57
				}
				add(ruleshellCmd, position58)
			}
			return true
		l57:
			position, tokenIndex = position57, tokenIndex57
			return false
		},
		/* 10 nullCmd <- <startAddr?> */
		func() bool {
			{
				position61 := position
				{
					position62, tokenIndex62 := position, tokenIndex
					if !_rules[rulestartAddr]() {
						goto l62
					}
					goto l63
				l62:
					position, tokenIndex = position62, tokenIndex62
				}
			l63:
				add(rulenullCmd, position61)
			}
			return true
		},
		/* 11 text <- <(<(!textTerm .)*> textTerm Action8)> */
		func() bool {
			position64, tokenIndex64 := position, tokenIndex
			{
				position65 := position
				{
					position66 := position
				l67:
					{
						position68, tokenIndex68 := position, tokenIndex
						{
							position69, tokenIndex69 := position, tokenIndex
							if !_rules[ruletextTerm]() {
								goto l69
							}
							goto l68
						l69:
							position, tokenIndex = position69, tokenIndex69
						}
						if !matchDot() {
							goto l68
						}
						goto l67
					l68:
						position, tokenIndex = position68, tokenIndex68
					}
					add(rulePegText, position66)
				}
				if !_rules[ruletextTerm]() {
					goto l64
				}
				if !_rules[ruleAction8]() {
					goto l64
				}
				add(ruletext, position65)
			}
			return true
		l64:
			position, tokenIndex = position64, tokenIndex64
			return false
		},
		/* 12 textTerm <- <('\n' '.')> */
		func() bool {
			position70, tokenIndex70 := position, tokenIndex
			{
				position71 := position
				if buffer[position] != rune('\n') {
					goto l70
				}
				position++
				if buffer[position] != rune('.') {
					goto l70
				}
				position++
				add(ruletextTerm, position71)
			}
			return true
		l70:
			position, tokenIndex = position70, tokenIndex70
			return false
		},
		/* 13 rangeCmd <- <(range? sp* rangeC)> */
		func() bool {
			position72, tokenIndex72 := position, tokenIndex
			{
				position73 := position
				{
					position74, tokenIndex74 := position, tokenIndex
					if !_rules[rulerange]() {
						goto l74
					}
					goto l75
				l74:
					position, tokenIndex = position74, tokenIndex74
				}
			l75:
			l76:
				{
					position77, tokenIndex77 := position, tokenIndex
					if !_rules[rulesp]() {
						goto l77
					}
					goto l76
				l77:
					position, tokenIndex = position77, tokenIndex77
				}
				if !_rules[rulerangeC]() {
					goto l72
				}
				add(rulerangeCmd, position73)
			}
			return true
		l72:
			position, tokenIndex = position72, tokenIndex72
			return false
		},
		/* 14 range <- <((startAddr ',' endAddr) / (startAddr ',' Action9) / (',' endAddr Action10) / (startAddr ';' endAddr) / (startAddr ';' Action11) / (';' endAddr Action12) / (startAddr Action13) / (',' Action14) / (';' Action15))> */
		func() bool {
			position78, tokenIndex78 := position, tokenIndex
			{
				position79 := position
				{
					position80, tokenIndex80 := position, tokenIndex
					if !_rules[rulestartAddr]() {
						goto l81
					}
					if buffer[position] != rune(',') {
						goto l81
					}
					position++
					if !_rules[ruleendAddr]() {
						goto l81
					}
					goto l80
				l81:
					position, tokenIndex = position80, tokenIndex80
					if !_rules[rulestartAddr]() {
						goto l82
					}
					if buffer[position] != rune(',') {
						goto l82
					}
					position++
					if !_rules[ruleAction9]() {
						goto l82
					}
					goto l80
				l82:
					position, tokenIndex = position80, tokenIndex80
					if buffer[position] != rune(',') {
						goto l83
					}
					position++
					if !_rules[ruleendAddr]() {
						goto l83
					}
					if !_rules[ruleAction10]() {
						goto l83
					}
					goto l80
				l83:
					position, tokenIndex = position80, tokenIndex80
					if !_rules[rulestartAddr]() {
						goto l84
					}
					if buffer[position] != rune(';') {
						goto l84
					}
					position++
					if !_rules[ruleendAddr]() {
						goto l84
					}
					goto l80
				l84:
					position, tokenIndex = position80, tokenIndex80
					if !_rules[rulestartAddr]() {
						goto l85
					}
					if buffer[position] != rune(';') {
						goto l85
					}
					position++
					if !_rules[ruleAction11]() {
						goto l85
					}
					goto l80
				l85:
					position, tokenIndex = position80, tokenIndex80
					if buffer[position] != rune(';') {
						goto l86
					}
					position++
					if !_rules[ruleendAddr]() {
						goto l86
					}
					if !_rules[ruleAction12]() {
						goto l86
					}
					goto l80
				l86:
					position, tokenIndex = position80, tokenIndex80
					if !_rules[rulestartAddr]() {
						goto l87
					}
					if !_rules[ruleAction13]() {
						goto l87
					}
					goto l80
				l87:
					position, tokenIndex = position80, tokenIndex80
					if buffer[position] != rune(',') {
						goto l88
					}
					position++
					if !_rules[ruleAction14]() {
						goto l88
					}
					goto l80
				l88:
					position, tokenIndex = position80, tokenIndex80
					if buffer[position] != rune(';') {
						goto l78
					}
					position++
					if !_rules[ruleAction15]() {
						goto l78
					}
				}
			l80:
				add(rulerange, position79)
			}
			return true
		l78:
			position, tokenIndex = position78, tokenIndex78
			return false
		},
		/* 15 addrCmd <- <((<startAddr> addrC Action16) / addrC)> */
		func() bool {
			position89, tokenIndex89 := position, tokenIndex
			{
				position90 := position
				{
					position91, tokenIndex91 := position, tokenIndex
					{
						position93 := position
						if !_rules[rulestartAddr]() {
							goto l92
						}
						add(rulePegText, position93)
					}
					if !_rules[ruleaddrC]() {
						goto l92
					}
					if !_rules[ruleAction16]() {
						goto l92
					}
					goto l91
				l92:
					position, tokenIndex = position91, tokenIndex91
					if !_rules[ruleaddrC]() {
						goto l89
					}
				}
			l91:
				add(ruleaddrCmd, position90)
			}
			return true
		l89:
			position, tokenIndex = position89, tokenIndex89
			return false
		},
		/* 16 startAddr <- <(sp* addrO sp* Action17)> */
		func() bool {
			position94, tokenIndex94 := position, tokenIndex
			{
				position95 := position
			l96:
				{
					position97, tokenIndex97 := position, tokenIndex
					if !_rules[rulesp]() {
						goto l97
					}
					goto l96
				l97:
					position, tokenIndex = position97, tokenIndex97
				}
				if !_rules[ruleaddrO]() {
					goto l94
				}
			l98:
				{
					position99, tokenIndex99 := position, tokenIndex
					if !_rules[rulesp]() {
						goto l99
					}
					goto l98
				l99:
					position, tokenIndex = position99, tokenIndex99
				}
				if !_rules[ruleAction17]() {
					goto l94
				}
				add(rulestartAddr, position95)
			}
			return true
		l94:
			position, tokenIndex = position94, tokenIndex94
			return false
		},
		/* 17 endAddr <- <(sp* addrO sp* Action18)> */
		func() bool {
			position100, tokenIndex100 := position, tokenIndex
			{
				position101 := position
			l102:
				{
					position103, tokenIndex103 := position, tokenIndex
					if !_rules[rulesp]() {
						goto l103
					}
					goto l102
				l103:
					position, tokenIndex = position103, tokenIndex103
				}
				if !_rules[ruleaddrO]() {
					goto l100
				}
			l104:
				{
					position105, tokenIndex105 := position, tokenIndex
					if !_rules[rulesp]() {
						goto l105
					}
					goto l104
				l105:
					position, tokenIndex = position105, tokenIndex105
				}
				if !_rules[ruleAction18]() {
					goto l100
				}
				add(ruleendAddr, position101)
			}
			return true
		l100:
			position, tokenIndex = position100, tokenIndex100
			return false
		},
		/* 18 addrO <- <(<(addr offset?)> Action19)> */
		func() bool {
			position106, tokenIndex106 := position, tokenIndex
			{
				position107 := position
				{
					position108 := position
					if !_rules[ruleaddr]() {
						goto l106
					}
					{
						position109, tokenIndex109 := position, tokenIndex
						if !_rules[ruleoffset]() {
							goto l109
						}
						goto l110
					l109:
						position, tokenIndex = position109, tokenIndex109
					}
				l110:
					add(rulePegText, position108)
				}
				if !_rules[ruleAction19]() {
					goto l106
				}
				add(ruleaddrO, position107)
			}
			return true
		l106:
			position, tokenIndex = position106, tokenIndex106
			return false
		},
		/* 19 addr <- <(literalAddr / markAddr / regexAddr / regexReverseAddr / ('.' Action20) / ('$' Action21))> */
		func() bool {
			position111, tokenIndex111 := position, tokenIndex
			{
				position112 := position
				{
					position113, tokenIndex113 := position, tokenIndex
					if !_rules[ruleliteralAddr]() {
						goto l114
					}
					goto l113
				l114:
					position, tokenIndex = position113, tokenIndex113
					if !_rules[rulemarkAddr]() {
						goto l115
					}
					goto l113
				l115:
					position, tokenIndex = position113, tokenIndex113
					if !_rules[ruleregexAddr]() {
						goto l116
					}
					goto l113
				l116:
					position, tokenIndex = position113, tokenIndex113
					if !_rules[ruleregexReverseAddr]() {
						goto l117
					}
					goto l113
				l117:
					position, tokenIndex = position113, tokenIndex113
					if buffer[position] != rune('.') {
						goto l118
					}
					position++
					if !_rules[ruleAction20]() {
						goto l118
					}
					goto l113
				l118:
					position, tokenIndex = position113, tokenIndex113
					if buffer[position] != rune('$') {
						goto l111
					}
					position++
					if !_rules[ruleAction21]() {
						goto l111
					}
				}
			l113:
				add(ruleaddr, position112)
			}
			return true
		l111:
			position, tokenIndex = position111, tokenIndex111
			return false
		},
		/* 20 literalAddr <- <(<[0-9]+> Action22)> */
		func() bool {
			position119, tokenIndex119 := position, tokenIndex
			{
				position120 := position
				{
					position121 := position
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l119
					}
					position++
				l122:
					{
						position123, tokenIndex123 := position, tokenIndex
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l123
						}
						position++
						goto l122
					l123:
						position, tokenIndex = position123, tokenIndex123
					}
					add(rulePegText, position121)
				}
				if !_rules[ruleAction22]() {
					goto l119
				}
				add(ruleliteralAddr, position120)
			}
			return true
		l119:
			position, tokenIndex = position119, tokenIndex119
			return false
		},
		/* 21 markAddr <- <('\'' [a-z] Action23)> */
		func() bool {
			position124, tokenIndex124 := position, tokenIndex
			{
				position125 := position
				if buffer[position] != rune('\'') {
					goto l124
				}
				position++
				if c := buffer[position]; c < rune('a') || c > rune('z') {
					goto l124
				}
				position++
				if !_rules[ruleAction23]() {
					goto l124
				}
				add(rulemarkAddr, position125)
			}
			return true
		l124:
			position, tokenIndex = position124, tokenIndex124
			return false
		},
		/* 22 regexAddr <- <('/' basic_regex '/' Action24)> */
		func() bool {
			position126, tokenIndex126 := position, tokenIndex
			{
				position127 := position
				if buffer[position] != rune('/') {
					goto l126
				}
				position++
				if !_rules[rulebasic_regex]() {
					goto l126
				}
				if buffer[position] != rune('/') {
					goto l126
				}
				position++
				if !_rules[ruleAction24]() {
					goto l126
				}
				add(ruleregexAddr, position127)
			}
			return true
		l126:
			position, tokenIndex = position126, tokenIndex126
			return false
		},
		/* 23 regexReverseAddr <- <('?' back_regex '?' Action25)> */
		func() bool {
			position128, tokenIndex128 := position, tokenIndex
			{
				position129 := position
				if buffer[position] != rune('?') {
					goto l128
				}
				position++
				if !_rules[ruleback_regex]() {
					goto l128
				}
				if buffer[position] != rune('?') {
					goto l128
				}
				position++
				if !_rules[ruleAction25]() {
					goto l128
				}
				add(ruleregexReverseAddr, position129)
			}
			return true
		l128:
			position, tokenIndex = position128, tokenIndex128
			return false
		},
		/* 24 basic_regex <- <(('\\' '/') / (!('\n' / '/') .))+> */
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
					if buffer[position] != rune('/') {
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
							if buffer[position] != rune('/') {
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
						if buffer[position] != rune('/') {
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
								if buffer[position] != rune('/') {
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
				add(rulebasic_regex, position131)
			}
			return true
		l130:
			position, tokenIndex = position130, tokenIndex130
			return false
		},
		/* 25 back_regex <- <(('\\' '?') / (!('\n' / '?') .))+> */
		func() bool {
			position144, tokenIndex144 := position, tokenIndex
			{
				position145 := position
				{
					position148, tokenIndex148 := position, tokenIndex
					if buffer[position] != rune('\\') {
						goto l149
					}
					position++
					if buffer[position] != rune('?') {
						goto l149
					}
					position++
					goto l148
				l149:
					position, tokenIndex = position148, tokenIndex148
					{
						position150, tokenIndex150 := position, tokenIndex
						{
							position151, tokenIndex151 := position, tokenIndex
							if buffer[position] != rune('\n') {
								goto l152
							}
							position++
							goto l151
						l152:
							position, tokenIndex = position151, tokenIndex151
							if buffer[position] != rune('?') {
								goto l150
							}
							position++
						}
					l151:
						goto l144
					l150:
						position, tokenIndex = position150, tokenIndex150
					}
					if !matchDot() {
						goto l144
					}
				}
			l148:
			l146:
				{
					position147, tokenIndex147 := position, tokenIndex
					{
						position153, tokenIndex153 := position, tokenIndex
						if buffer[position] != rune('\\') {
							goto l154
						}
						position++
						if buffer[position] != rune('?') {
							goto l154
						}
						position++
						goto l153
					l154:
						position, tokenIndex = position153, tokenIndex153
						{
							position155, tokenIndex155 := position, tokenIndex
							{
								position156, tokenIndex156 := position, tokenIndex
								if buffer[position] != rune('\n') {
									goto l157
								}
								position++
								goto l156
							l157:
								position, tokenIndex = position156, tokenIndex156
								if buffer[position] != rune('?') {
									goto l155
								}
								position++
							}
						l156:
							goto l147
						l155:
							position, tokenIndex = position155, tokenIndex155
						}
						if !matchDot() {
							goto l147
						}
					}
				l153:
					goto l146
				l147:
					position, tokenIndex = position147, tokenIndex147
				}
				add(ruleback_regex, position145)
			}
			return true
		l144:
			position, tokenIndex = position144, tokenIndex144
			return false
		},
		/* 26 bareCmd <- <(('h' Action26) / ('H' Action27) / ('P' Action28) / ('q' Action29) / ('Q' Action30) / ('u' Action31))> */
		func() bool {
			position158, tokenIndex158 := position, tokenIndex
			{
				position159 := position
				{
					position160, tokenIndex160 := position, tokenIndex
					if buffer[position] != rune('h') {
						goto l161
					}
					position++
					if !_rules[ruleAction26]() {
						goto l161
					}
					goto l160
				l161:
					position, tokenIndex = position160, tokenIndex160
					if buffer[position] != rune('H') {
						goto l162
					}
					position++
					if !_rules[ruleAction27]() {
						goto l162
					}
					goto l160
				l162:
					position, tokenIndex = position160, tokenIndex160
					if buffer[position] != rune('P') {
						goto l163
					}
					position++
					if !_rules[ruleAction28]() {
						goto l163
					}
					goto l160
				l163:
					position, tokenIndex = position160, tokenIndex160
					if buffer[position] != rune('q') {
						goto l164
					}
					position++
					if !_rules[ruleAction29]() {
						goto l164
					}
					goto l160
				l164:
					position, tokenIndex = position160, tokenIndex160
					if buffer[position] != rune('Q') {
						goto l165
					}
					position++
					if !_rules[ruleAction30]() {
						goto l165
					}
					goto l160
				l165:
					position, tokenIndex = position160, tokenIndex160
					if buffer[position] != rune('u') {
						goto l158
					}
					position++
					if !_rules[ruleAction31]() {
						goto l158
					}
				}
			l160:
				add(rulebareCmd, position159)
			}
			return true
		l158:
			position, tokenIndex = position158, tokenIndex158
			return false
		},
		/* 27 offset <- <(('+' / '-') [0-9]*)> */
		func() bool {
			position166, tokenIndex166 := position, tokenIndex
			{
				position167 := position
				{
					position168, tokenIndex168 := position, tokenIndex
					if buffer[position] != rune('+') {
						goto l169
					}
					position++
					goto l168
				l169:
					position, tokenIndex = position168, tokenIndex168
					if buffer[position] != rune('-') {
						goto l166
					}
					position++
				}
			l168:
			l170:
				{
					position171, tokenIndex171 := position, tokenIndex
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l171
					}
					position++
					goto l170
				l171:
					position, tokenIndex = position171, tokenIndex171
				}
				add(ruleoffset, position167)
			}
			return true
		l166:
			position, tokenIndex = position166, tokenIndex166
			return false
		},
		/* 28 paramCmd <- <((paramC sp <param> Action32) / paramC)> */
		func() bool {
			position172, tokenIndex172 := position, tokenIndex
			{
				position173 := position
				{
					position174, tokenIndex174 := position, tokenIndex
					if !_rules[ruleparamC]() {
						goto l175
					}
					if !_rules[rulesp]() {
						goto l175
					}
					{
						position176 := position
						if !_rules[ruleparam]() {
							goto l175
						}
						add(rulePegText, position176)
					}
					if !_rules[ruleAction32]() {
						goto l175
					}
					goto l174
				l175:
					position, tokenIndex = position174, tokenIndex174
					if !_rules[ruleparamC]() {
						goto l172
					}
				}
			l174:
				add(ruleparamCmd, position173)
			}
			return true
		l172:
			position, tokenIndex = position172, tokenIndex172
			return false
		},
		/* 29 paramC <- <(('e' Action33) / ('E' Action34) / ('f' Action35))> */
		func() bool {
			position177, tokenIndex177 := position, tokenIndex
			{
				position178 := position
				{
					position179, tokenIndex179 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l180
					}
					position++
					if !_rules[ruleAction33]() {
						goto l180
					}
					goto l179
				l180:
					position, tokenIndex = position179, tokenIndex179
					if buffer[position] != rune('E') {
						goto l181
					}
					position++
					if !_rules[ruleAction34]() {
						goto l181
					}
					goto l179
				l181:
					position, tokenIndex = position179, tokenIndex179
					if buffer[position] != rune('f') {
						goto l177
					}
					position++
					if !_rules[ruleAction35]() {
						goto l177
					}
				}
			l179:
				add(ruleparamC, position178)
			}
			return true
		l177:
			position, tokenIndex = position177, tokenIndex177
			return false
		},
		/* 30 param <- <(!'\n' .)+> */
		func() bool {
			position182, tokenIndex182 := position, tokenIndex
			{
				position183 := position
				{
					position186, tokenIndex186 := position, tokenIndex
					if buffer[position] != rune('\n') {
						goto l186
					}
					position++
					goto l182
				l186:
					position, tokenIndex = position186, tokenIndex186
				}
				if !matchDot() {
					goto l182
				}
			l184:
				{
					position185, tokenIndex185 := position, tokenIndex
					{
						position187, tokenIndex187 := position, tokenIndex
						if buffer[position] != rune('\n') {
							goto l187
						}
						position++
						goto l185
					l187:
						position, tokenIndex = position187, tokenIndex187
					}
					if !matchDot() {
						goto l185
					}
					goto l184
				l185:
					position, tokenIndex = position185, tokenIndex185
				}
				add(ruleparam, position183)
			}
			return true
		l182:
			position, tokenIndex = position182, tokenIndex182
			return false
		},
		/* 31 addrC <- <('=' Action36)> */
		func() bool {
			position188, tokenIndex188 := position, tokenIndex
			{
				position189 := position
				if buffer[position] != rune('=') {
					goto l188
				}
				position++
				if !_rules[ruleAction36]() {
					goto l188
				}
				add(ruleaddrC, position189)
			}
			return true
		l188:
			position, tokenIndex = position188, tokenIndex188
			return false
		},
		/* 32 changeTextC <- <('c' Action37)> */
		func() bool {
			position190, tokenIndex190 := position, tokenIndex
			{
				position191 := position
				if buffer[position] != rune('c') {
					goto l190
				}
				position++
				if !_rules[ruleAction37]() {
					goto l190
				}
				add(rulechangeTextC, position191)
			}
			return true
		l190:
			position, tokenIndex = position190, tokenIndex190
			return false
		},
		/* 33 addTextC <- <(('a' Action38) / ('i' Action39))> */
		func() bool {
			position192, tokenIndex192 := position, tokenIndex
			{
				position193 := position
				{
					position194, tokenIndex194 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l195
					}
					position++
					if !_rules[ruleAction38]() {
						goto l195
					}
					goto l194
				l195:
					position, tokenIndex = position194, tokenIndex194
					if buffer[position] != rune('i') {
						goto l192
					}
					position++
					if !_rules[ruleAction39]() {
						goto l192
					}
				}
			l194:
				add(ruleaddTextC, position193)
			}
			return true
		l192:
			position, tokenIndex = position192, tokenIndex192
			return false
		},
		/* 34 rangeC <- <(('d' Action40) / ('j' Action41) / ('l' Action42) / ('n' Action43) / ('p' Action44))> */
		func() bool {
			position196, tokenIndex196 := position, tokenIndex
			{
				position197 := position
				{
					position198, tokenIndex198 := position, tokenIndex
					if buffer[position] != rune('d') {
						goto l199
					}
					position++
					if !_rules[ruleAction40]() {
						goto l199
					}
					goto l198
				l199:
					position, tokenIndex = position198, tokenIndex198
					if buffer[position] != rune('j') {
						goto l200
					}
					position++
					if !_rules[ruleAction41]() {
						goto l200
					}
					goto l198
				l200:
					position, tokenIndex = position198, tokenIndex198
					if buffer[position] != rune('l') {
						goto l201
					}
					position++
					if !_rules[ruleAction42]() {
						goto l201
					}
					goto l198
				l201:
					position, tokenIndex = position198, tokenIndex198
					if buffer[position] != rune('n') {
						goto l202
					}
					position++
					if !_rules[ruleAction43]() {
						goto l202
					}
					goto l198
				l202:
					position, tokenIndex = position198, tokenIndex198
					if buffer[position] != rune('p') {
						goto l196
					}
					position++
					if !_rules[ruleAction44]() {
						goto l196
					}
				}
			l198:
				add(rulerangeC, position197)
			}
			return true
		l196:
			position, tokenIndex = position196, tokenIndex196
			return false
		},
		/* 35 destC <- <(('m' Action45) / ('t' Action46))> */
		func() bool {
			position203, tokenIndex203 := position, tokenIndex
			{
				position204 := position
				{
					position205, tokenIndex205 := position, tokenIndex
					if buffer[position] != rune('m') {
						goto l206
					}
					position++
					if !_rules[ruleAction45]() {
						goto l206
					}
					goto l205
				l206:
					position, tokenIndex = position205, tokenIndex205
					if buffer[position] != rune('t') {
						goto l203
					}
					position++
					if !_rules[ruleAction46]() {
						goto l203
					}
				}
			l205:
				add(ruledestC, position204)
			}
			return true
		l203:
			position, tokenIndex = position203, tokenIndex203
			return false
		},
		/* 36 newLine <- <'\n'> */
		func() bool {
			position207, tokenIndex207 := position, tokenIndex
			{
				position208 := position
				if buffer[position] != rune('\n') {
					goto l207
				}
				position++
				add(rulenewLine, position208)
			}
			return true
		l207:
			position, tokenIndex = position207, tokenIndex207
			return false
		},
		/* 37 sp <- <(' ' / '\t')+> */
		func() bool {
			position209, tokenIndex209 := position, tokenIndex
			{
				position210 := position
				{
					position213, tokenIndex213 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l214
					}
					position++
					goto l213
				l214:
					position, tokenIndex = position213, tokenIndex213
					if buffer[position] != rune('\t') {
						goto l209
					}
					position++
				}
			l213:
			l211:
				{
					position212, tokenIndex212 := position, tokenIndex
					{
						position215, tokenIndex215 := position, tokenIndex
						if buffer[position] != rune(' ') {
							goto l216
						}
						position++
						goto l215
					l216:
						position, tokenIndex = position215, tokenIndex215
						if buffer[position] != rune('\t') {
							goto l212
						}
						position++
					}
				l215:
					goto l211
				l212:
					position, tokenIndex = position212, tokenIndex212
				}
				add(rulesp, position210)
			}
			return true
		l209:
			position, tokenIndex = position209, tokenIndex209
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
