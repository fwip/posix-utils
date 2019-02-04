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
	rules  [80]func() bool
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
		case ruleAction9:
			p.curCmd.start.text = buffer[begin:end]
		case ruleAction10:
			p.curCmd.start = p.curAddr
			p.curAddr = address{}
		case ruleAction11:
			p.curCmd.end = p.curAddr
			p.curAddr = address{}
		case ruleAction12:
			p.curAddr.text = buffer[begin:end]
		case ruleAction13:
			p.curAddr.typ = lCurrent
		case ruleAction14:
			p.curAddr.typ = lLast
		case ruleAction15:
			p.curAddr.typ = lNum
		case ruleAction16:
			p.curAddr.typ = lMark
		case ruleAction17:
			p.curAddr.typ = lRegex
		case ruleAction18:
			p.curAddr.typ = lRegexReverse
		case ruleAction19:
			p.curCmd.typ = cthelp
		case ruleAction20:
			p.curCmd.typ = cthelpMode
		case ruleAction21:
			p.curCmd.typ = ctprompt
		case ruleAction22:
			p.curCmd.typ = ctquit
		case ruleAction23:
			p.curCmd.typ = ctquitForce
		case ruleAction24:
			p.curCmd.typ = ctundo
		case ruleAction25:
			p.curCmd.params = []string{buffer[begin:end]}
		case ruleAction26:
			p.curCmd.typ = ctedit
		case ruleAction27:
			p.curCmd.typ = cteditForce
		case ruleAction28:
			p.curCmd.typ = ctfilename
		case ruleAction29:
			p.curCmd.typ = ctlineNumber
		case ruleAction30:
			p.curCmd.typ = ctchange
		case ruleAction31:
			p.curCmd.typ = ctappend
		case ruleAction32:
			p.curCmd.typ = ctinsert
		case ruleAction33:
			p.curCmd.typ = ctdelete
		case ruleAction34:
			p.curCmd.typ = ctjoin
		case ruleAction35:
			p.curCmd.typ = ctlist
		case ruleAction36:
			p.curCmd.typ = ctnumber
		case ruleAction37:
			p.curCmd.typ = ctprint
		case ruleAction38:
			p.curCmd.typ = ctmove
		case ruleAction39:
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
		/* 14 range <- <((startAddr ',' endAddr) / (startAddr ',') / (',' endAddr) / (startAddr ';' endAddr) / (startAddr ';') / (';' endAddr) / startAddr)> */
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
					goto l74
				l80:
					position, tokenIndex = position74, tokenIndex74
					if !_rules[rulestartAddr]() {
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
		/* 15 addrCmd <- <((<startAddr> addrC Action9) / addrC)> */
		func() bool {
			position81, tokenIndex81 := position, tokenIndex
			{
				position82 := position
				{
					position83, tokenIndex83 := position, tokenIndex
					{
						position85 := position
						if !_rules[rulestartAddr]() {
							goto l84
						}
						add(rulePegText, position85)
					}
					if !_rules[ruleaddrC]() {
						goto l84
					}
					if !_rules[ruleAction9]() {
						goto l84
					}
					goto l83
				l84:
					position, tokenIndex = position83, tokenIndex83
					if !_rules[ruleaddrC]() {
						goto l81
					}
				}
			l83:
				add(ruleaddrCmd, position82)
			}
			return true
		l81:
			position, tokenIndex = position81, tokenIndex81
			return false
		},
		/* 16 startAddr <- <(addrO Action10)> */
		func() bool {
			position86, tokenIndex86 := position, tokenIndex
			{
				position87 := position
				if !_rules[ruleaddrO]() {
					goto l86
				}
				if !_rules[ruleAction10]() {
					goto l86
				}
				add(rulestartAddr, position87)
			}
			return true
		l86:
			position, tokenIndex = position86, tokenIndex86
			return false
		},
		/* 17 endAddr <- <(addrO Action11)> */
		func() bool {
			position88, tokenIndex88 := position, tokenIndex
			{
				position89 := position
				if !_rules[ruleaddrO]() {
					goto l88
				}
				if !_rules[ruleAction11]() {
					goto l88
				}
				add(ruleendAddr, position89)
			}
			return true
		l88:
			position, tokenIndex = position88, tokenIndex88
			return false
		},
		/* 18 addrO <- <(<(addr offset?)> Action12)> */
		func() bool {
			position90, tokenIndex90 := position, tokenIndex
			{
				position91 := position
				{
					position92 := position
					if !_rules[ruleaddr]() {
						goto l90
					}
					{
						position93, tokenIndex93 := position, tokenIndex
						if !_rules[ruleoffset]() {
							goto l93
						}
						goto l94
					l93:
						position, tokenIndex = position93, tokenIndex93
					}
				l94:
					add(rulePegText, position92)
				}
				if !_rules[ruleAction12]() {
					goto l90
				}
				add(ruleaddrO, position91)
			}
			return true
		l90:
			position, tokenIndex = position90, tokenIndex90
			return false
		},
		/* 19 addr <- <(literalAddr / markAddr / regexAddr / regexReverseAddr / ('.' Action13) / ('$' Action14))> */
		func() bool {
			position95, tokenIndex95 := position, tokenIndex
			{
				position96 := position
				{
					position97, tokenIndex97 := position, tokenIndex
					if !_rules[ruleliteralAddr]() {
						goto l98
					}
					goto l97
				l98:
					position, tokenIndex = position97, tokenIndex97
					if !_rules[rulemarkAddr]() {
						goto l99
					}
					goto l97
				l99:
					position, tokenIndex = position97, tokenIndex97
					if !_rules[ruleregexAddr]() {
						goto l100
					}
					goto l97
				l100:
					position, tokenIndex = position97, tokenIndex97
					if !_rules[ruleregexReverseAddr]() {
						goto l101
					}
					goto l97
				l101:
					position, tokenIndex = position97, tokenIndex97
					if buffer[position] != rune('.') {
						goto l102
					}
					position++
					if !_rules[ruleAction13]() {
						goto l102
					}
					goto l97
				l102:
					position, tokenIndex = position97, tokenIndex97
					if buffer[position] != rune('$') {
						goto l95
					}
					position++
					if !_rules[ruleAction14]() {
						goto l95
					}
				}
			l97:
				add(ruleaddr, position96)
			}
			return true
		l95:
			position, tokenIndex = position95, tokenIndex95
			return false
		},
		/* 20 literalAddr <- <(<[0-9]+> Action15)> */
		func() bool {
			position103, tokenIndex103 := position, tokenIndex
			{
				position104 := position
				{
					position105 := position
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l103
					}
					position++
				l106:
					{
						position107, tokenIndex107 := position, tokenIndex
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l107
						}
						position++
						goto l106
					l107:
						position, tokenIndex = position107, tokenIndex107
					}
					add(rulePegText, position105)
				}
				if !_rules[ruleAction15]() {
					goto l103
				}
				add(ruleliteralAddr, position104)
			}
			return true
		l103:
			position, tokenIndex = position103, tokenIndex103
			return false
		},
		/* 21 markAddr <- <('\'' [a-z] Action16)> */
		func() bool {
			position108, tokenIndex108 := position, tokenIndex
			{
				position109 := position
				if buffer[position] != rune('\'') {
					goto l108
				}
				position++
				if c := buffer[position]; c < rune('a') || c > rune('z') {
					goto l108
				}
				position++
				if !_rules[ruleAction16]() {
					goto l108
				}
				add(rulemarkAddr, position109)
			}
			return true
		l108:
			position, tokenIndex = position108, tokenIndex108
			return false
		},
		/* 22 regexAddr <- <('/' basic_regex '/' Action17)> */
		func() bool {
			position110, tokenIndex110 := position, tokenIndex
			{
				position111 := position
				if buffer[position] != rune('/') {
					goto l110
				}
				position++
				if !_rules[rulebasic_regex]() {
					goto l110
				}
				if buffer[position] != rune('/') {
					goto l110
				}
				position++
				if !_rules[ruleAction17]() {
					goto l110
				}
				add(ruleregexAddr, position111)
			}
			return true
		l110:
			position, tokenIndex = position110, tokenIndex110
			return false
		},
		/* 23 regexReverseAddr <- <('?' back_regex '?' Action18)> */
		func() bool {
			position112, tokenIndex112 := position, tokenIndex
			{
				position113 := position
				if buffer[position] != rune('?') {
					goto l112
				}
				position++
				if !_rules[ruleback_regex]() {
					goto l112
				}
				if buffer[position] != rune('?') {
					goto l112
				}
				position++
				if !_rules[ruleAction18]() {
					goto l112
				}
				add(ruleregexReverseAddr, position113)
			}
			return true
		l112:
			position, tokenIndex = position112, tokenIndex112
			return false
		},
		/* 24 basic_regex <- <(('\\' '/') / (!('\n' / '/') .))+> */
		func() bool {
			position114, tokenIndex114 := position, tokenIndex
			{
				position115 := position
				{
					position118, tokenIndex118 := position, tokenIndex
					if buffer[position] != rune('\\') {
						goto l119
					}
					position++
					if buffer[position] != rune('/') {
						goto l119
					}
					position++
					goto l118
				l119:
					position, tokenIndex = position118, tokenIndex118
					{
						position120, tokenIndex120 := position, tokenIndex
						{
							position121, tokenIndex121 := position, tokenIndex
							if buffer[position] != rune('\n') {
								goto l122
							}
							position++
							goto l121
						l122:
							position, tokenIndex = position121, tokenIndex121
							if buffer[position] != rune('/') {
								goto l120
							}
							position++
						}
					l121:
						goto l114
					l120:
						position, tokenIndex = position120, tokenIndex120
					}
					if !matchDot() {
						goto l114
					}
				}
			l118:
			l116:
				{
					position117, tokenIndex117 := position, tokenIndex
					{
						position123, tokenIndex123 := position, tokenIndex
						if buffer[position] != rune('\\') {
							goto l124
						}
						position++
						if buffer[position] != rune('/') {
							goto l124
						}
						position++
						goto l123
					l124:
						position, tokenIndex = position123, tokenIndex123
						{
							position125, tokenIndex125 := position, tokenIndex
							{
								position126, tokenIndex126 := position, tokenIndex
								if buffer[position] != rune('\n') {
									goto l127
								}
								position++
								goto l126
							l127:
								position, tokenIndex = position126, tokenIndex126
								if buffer[position] != rune('/') {
									goto l125
								}
								position++
							}
						l126:
							goto l117
						l125:
							position, tokenIndex = position125, tokenIndex125
						}
						if !matchDot() {
							goto l117
						}
					}
				l123:
					goto l116
				l117:
					position, tokenIndex = position117, tokenIndex117
				}
				add(rulebasic_regex, position115)
			}
			return true
		l114:
			position, tokenIndex = position114, tokenIndex114
			return false
		},
		/* 25 back_regex <- <(('\\' '?') / (!('\n' / '?') .))+> */
		func() bool {
			position128, tokenIndex128 := position, tokenIndex
			{
				position129 := position
				{
					position132, tokenIndex132 := position, tokenIndex
					if buffer[position] != rune('\\') {
						goto l133
					}
					position++
					if buffer[position] != rune('?') {
						goto l133
					}
					position++
					goto l132
				l133:
					position, tokenIndex = position132, tokenIndex132
					{
						position134, tokenIndex134 := position, tokenIndex
						{
							position135, tokenIndex135 := position, tokenIndex
							if buffer[position] != rune('\n') {
								goto l136
							}
							position++
							goto l135
						l136:
							position, tokenIndex = position135, tokenIndex135
							if buffer[position] != rune('?') {
								goto l134
							}
							position++
						}
					l135:
						goto l128
					l134:
						position, tokenIndex = position134, tokenIndex134
					}
					if !matchDot() {
						goto l128
					}
				}
			l132:
			l130:
				{
					position131, tokenIndex131 := position, tokenIndex
					{
						position137, tokenIndex137 := position, tokenIndex
						if buffer[position] != rune('\\') {
							goto l138
						}
						position++
						if buffer[position] != rune('?') {
							goto l138
						}
						position++
						goto l137
					l138:
						position, tokenIndex = position137, tokenIndex137
						{
							position139, tokenIndex139 := position, tokenIndex
							{
								position140, tokenIndex140 := position, tokenIndex
								if buffer[position] != rune('\n') {
									goto l141
								}
								position++
								goto l140
							l141:
								position, tokenIndex = position140, tokenIndex140
								if buffer[position] != rune('?') {
									goto l139
								}
								position++
							}
						l140:
							goto l131
						l139:
							position, tokenIndex = position139, tokenIndex139
						}
						if !matchDot() {
							goto l131
						}
					}
				l137:
					goto l130
				l131:
					position, tokenIndex = position131, tokenIndex131
				}
				add(ruleback_regex, position129)
			}
			return true
		l128:
			position, tokenIndex = position128, tokenIndex128
			return false
		},
		/* 26 bareCmd <- <(('h' Action19) / ('H' Action20) / ('P' Action21) / ('q' Action22) / ('Q' Action23) / ('u' Action24))> */
		func() bool {
			position142, tokenIndex142 := position, tokenIndex
			{
				position143 := position
				{
					position144, tokenIndex144 := position, tokenIndex
					if buffer[position] != rune('h') {
						goto l145
					}
					position++
					if !_rules[ruleAction19]() {
						goto l145
					}
					goto l144
				l145:
					position, tokenIndex = position144, tokenIndex144
					if buffer[position] != rune('H') {
						goto l146
					}
					position++
					if !_rules[ruleAction20]() {
						goto l146
					}
					goto l144
				l146:
					position, tokenIndex = position144, tokenIndex144
					if buffer[position] != rune('P') {
						goto l147
					}
					position++
					if !_rules[ruleAction21]() {
						goto l147
					}
					goto l144
				l147:
					position, tokenIndex = position144, tokenIndex144
					if buffer[position] != rune('q') {
						goto l148
					}
					position++
					if !_rules[ruleAction22]() {
						goto l148
					}
					goto l144
				l148:
					position, tokenIndex = position144, tokenIndex144
					if buffer[position] != rune('Q') {
						goto l149
					}
					position++
					if !_rules[ruleAction23]() {
						goto l149
					}
					goto l144
				l149:
					position, tokenIndex = position144, tokenIndex144
					if buffer[position] != rune('u') {
						goto l142
					}
					position++
					if !_rules[ruleAction24]() {
						goto l142
					}
				}
			l144:
				add(rulebareCmd, position143)
			}
			return true
		l142:
			position, tokenIndex = position142, tokenIndex142
			return false
		},
		/* 27 offset <- <(('+' / '-') [0-9]*)> */
		func() bool {
			position150, tokenIndex150 := position, tokenIndex
			{
				position151 := position
				{
					position152, tokenIndex152 := position, tokenIndex
					if buffer[position] != rune('+') {
						goto l153
					}
					position++
					goto l152
				l153:
					position, tokenIndex = position152, tokenIndex152
					if buffer[position] != rune('-') {
						goto l150
					}
					position++
				}
			l152:
			l154:
				{
					position155, tokenIndex155 := position, tokenIndex
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l155
					}
					position++
					goto l154
				l155:
					position, tokenIndex = position155, tokenIndex155
				}
				add(ruleoffset, position151)
			}
			return true
		l150:
			position, tokenIndex = position150, tokenIndex150
			return false
		},
		/* 28 paramCmd <- <((paramC sp <param> Action25) / paramC)> */
		func() bool {
			position156, tokenIndex156 := position, tokenIndex
			{
				position157 := position
				{
					position158, tokenIndex158 := position, tokenIndex
					if !_rules[ruleparamC]() {
						goto l159
					}
					if !_rules[rulesp]() {
						goto l159
					}
					{
						position160 := position
						if !_rules[ruleparam]() {
							goto l159
						}
						add(rulePegText, position160)
					}
					if !_rules[ruleAction25]() {
						goto l159
					}
					goto l158
				l159:
					position, tokenIndex = position158, tokenIndex158
					if !_rules[ruleparamC]() {
						goto l156
					}
				}
			l158:
				add(ruleparamCmd, position157)
			}
			return true
		l156:
			position, tokenIndex = position156, tokenIndex156
			return false
		},
		/* 29 paramC <- <(('e' Action26) / ('E' Action27) / ('f' Action28))> */
		func() bool {
			position161, tokenIndex161 := position, tokenIndex
			{
				position162 := position
				{
					position163, tokenIndex163 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l164
					}
					position++
					if !_rules[ruleAction26]() {
						goto l164
					}
					goto l163
				l164:
					position, tokenIndex = position163, tokenIndex163
					if buffer[position] != rune('E') {
						goto l165
					}
					position++
					if !_rules[ruleAction27]() {
						goto l165
					}
					goto l163
				l165:
					position, tokenIndex = position163, tokenIndex163
					if buffer[position] != rune('f') {
						goto l161
					}
					position++
					if !_rules[ruleAction28]() {
						goto l161
					}
				}
			l163:
				add(ruleparamC, position162)
			}
			return true
		l161:
			position, tokenIndex = position161, tokenIndex161
			return false
		},
		/* 30 param <- <(!'\n' .)+> */
		func() bool {
			position166, tokenIndex166 := position, tokenIndex
			{
				position167 := position
				{
					position170, tokenIndex170 := position, tokenIndex
					if buffer[position] != rune('\n') {
						goto l170
					}
					position++
					goto l166
				l170:
					position, tokenIndex = position170, tokenIndex170
				}
				if !matchDot() {
					goto l166
				}
			l168:
				{
					position169, tokenIndex169 := position, tokenIndex
					{
						position171, tokenIndex171 := position, tokenIndex
						if buffer[position] != rune('\n') {
							goto l171
						}
						position++
						goto l169
					l171:
						position, tokenIndex = position171, tokenIndex171
					}
					if !matchDot() {
						goto l169
					}
					goto l168
				l169:
					position, tokenIndex = position169, tokenIndex169
				}
				add(ruleparam, position167)
			}
			return true
		l166:
			position, tokenIndex = position166, tokenIndex166
			return false
		},
		/* 31 addrC <- <('=' Action29)> */
		func() bool {
			position172, tokenIndex172 := position, tokenIndex
			{
				position173 := position
				if buffer[position] != rune('=') {
					goto l172
				}
				position++
				if !_rules[ruleAction29]() {
					goto l172
				}
				add(ruleaddrC, position173)
			}
			return true
		l172:
			position, tokenIndex = position172, tokenIndex172
			return false
		},
		/* 32 changeTextC <- <('c' Action30)> */
		func() bool {
			position174, tokenIndex174 := position, tokenIndex
			{
				position175 := position
				if buffer[position] != rune('c') {
					goto l174
				}
				position++
				if !_rules[ruleAction30]() {
					goto l174
				}
				add(rulechangeTextC, position175)
			}
			return true
		l174:
			position, tokenIndex = position174, tokenIndex174
			return false
		},
		/* 33 addTextC <- <(('a' Action31) / ('i' Action32))> */
		func() bool {
			position176, tokenIndex176 := position, tokenIndex
			{
				position177 := position
				{
					position178, tokenIndex178 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l179
					}
					position++
					if !_rules[ruleAction31]() {
						goto l179
					}
					goto l178
				l179:
					position, tokenIndex = position178, tokenIndex178
					if buffer[position] != rune('i') {
						goto l176
					}
					position++
					if !_rules[ruleAction32]() {
						goto l176
					}
				}
			l178:
				add(ruleaddTextC, position177)
			}
			return true
		l176:
			position, tokenIndex = position176, tokenIndex176
			return false
		},
		/* 34 rangeC <- <(('d' Action33) / ('j' Action34) / ('l' Action35) / ('n' Action36) / ('p' Action37))> */
		func() bool {
			position180, tokenIndex180 := position, tokenIndex
			{
				position181 := position
				{
					position182, tokenIndex182 := position, tokenIndex
					if buffer[position] != rune('d') {
						goto l183
					}
					position++
					if !_rules[ruleAction33]() {
						goto l183
					}
					goto l182
				l183:
					position, tokenIndex = position182, tokenIndex182
					if buffer[position] != rune('j') {
						goto l184
					}
					position++
					if !_rules[ruleAction34]() {
						goto l184
					}
					goto l182
				l184:
					position, tokenIndex = position182, tokenIndex182
					if buffer[position] != rune('l') {
						goto l185
					}
					position++
					if !_rules[ruleAction35]() {
						goto l185
					}
					goto l182
				l185:
					position, tokenIndex = position182, tokenIndex182
					if buffer[position] != rune('n') {
						goto l186
					}
					position++
					if !_rules[ruleAction36]() {
						goto l186
					}
					goto l182
				l186:
					position, tokenIndex = position182, tokenIndex182
					if buffer[position] != rune('p') {
						goto l180
					}
					position++
					if !_rules[ruleAction37]() {
						goto l180
					}
				}
			l182:
				add(rulerangeC, position181)
			}
			return true
		l180:
			position, tokenIndex = position180, tokenIndex180
			return false
		},
		/* 35 destC <- <(('m' Action38) / ('t' Action39))> */
		func() bool {
			position187, tokenIndex187 := position, tokenIndex
			{
				position188 := position
				{
					position189, tokenIndex189 := position, tokenIndex
					if buffer[position] != rune('m') {
						goto l190
					}
					position++
					if !_rules[ruleAction38]() {
						goto l190
					}
					goto l189
				l190:
					position, tokenIndex = position189, tokenIndex189
					if buffer[position] != rune('t') {
						goto l187
					}
					position++
					if !_rules[ruleAction39]() {
						goto l187
					}
				}
			l189:
				add(ruledestC, position188)
			}
			return true
		l187:
			position, tokenIndex = position187, tokenIndex187
			return false
		},
		/* 36 newLine <- <'\n'> */
		func() bool {
			position191, tokenIndex191 := position, tokenIndex
			{
				position192 := position
				if buffer[position] != rune('\n') {
					goto l191
				}
				position++
				add(rulenewLine, position192)
			}
			return true
		l191:
			position, tokenIndex = position191, tokenIndex191
			return false
		},
		/* 37 sp <- <(' ' / '\t')+> */
		func() bool {
			position193, tokenIndex193 := position, tokenIndex
			{
				position194 := position
				{
					position197, tokenIndex197 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l198
					}
					position++
					goto l197
				l198:
					position, tokenIndex = position197, tokenIndex197
					if buffer[position] != rune('\t') {
						goto l193
					}
					position++
				}
			l197:
			l195:
				{
					position196, tokenIndex196 := position, tokenIndex
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
							goto l196
						}
						position++
					}
				l199:
					goto l195
				l196:
					position, tokenIndex = position196, tokenIndex196
				}
				add(rulesp, position194)
			}
			return true
		l193:
			position, tokenIndex = position193, tokenIndex193
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
		/* 48 Action8 <- <{p.curCmd.text = buffer[begin:end]}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 49 Action9 <- <{p.curCmd.start.text = buffer[begin:end]}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 50 Action10 <- <{p.curCmd.start = p.curAddr; p.curAddr = address{}}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 51 Action11 <- <{p.curCmd.end = p.curAddr; p.curAddr = address{}}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 52 Action12 <- <{p.curAddr.text = buffer[begin:end]}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 53 Action13 <- <{p.curAddr.typ = lCurrent}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 54 Action14 <- <{p.curAddr.typ = lLast}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 55 Action15 <- <{p.curAddr.typ = lNum}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 56 Action16 <- <{ p.curAddr.typ = lMark }> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 57 Action17 <- <{p.curAddr.typ = lRegex}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 58 Action18 <- <{p.curAddr.typ = lRegexReverse}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 59 Action19 <- <{p.curCmd.typ = cthelp}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 60 Action20 <- <{p.curCmd.typ = cthelpMode}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 61 Action21 <- <{p.curCmd.typ = ctprompt}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 62 Action22 <- <{p.curCmd.typ = ctquit}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 63 Action23 <- <{p.curCmd.typ = ctquitForce}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 64 Action24 <- <{p.curCmd.typ = ctundo}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 65 Action25 <- <{ p.curCmd.params = []string{buffer[begin:end]}}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 66 Action26 <- <{p.curCmd.typ = ctedit}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 67 Action27 <- <{p.curCmd.typ = cteditForce}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 68 Action28 <- <{p.curCmd.typ = ctfilename}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 69 Action29 <- <{p.curCmd.typ = ctlineNumber}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 70 Action30 <- <{p.curCmd.typ = ctchange}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 71 Action31 <- <{p.curCmd.typ = ctappend}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 72 Action32 <- <{p.curCmd.typ = ctinsert}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 73 Action33 <- <{p.curCmd.typ = ctdelete}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 74 Action34 <- <{p.curCmd.typ = ctjoin}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 75 Action35 <- <{p.curCmd.typ = ctlist}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 76 Action36 <- <{p.curCmd.typ = ctnumber}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 77 Action37 <- <{p.curCmd.typ = ctprint}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 78 Action38 <- <{p.curCmd.typ = ctmove}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 79 Action39 <- <{p.curCmd.typ = ctcopy}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
	}
	p.rules = _rules
}
