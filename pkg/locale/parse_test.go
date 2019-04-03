package locale

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

var parseTests = []string{
	/*
			``,          // Empty string
			`# comment`, // Comment
			`comment_char $
		$ bakery!`, // Comment char redefinition
	*/
	`LC_MONETARY
int_curr_symbol "USD "
END LC_MONETARY`, // Basic monetary block
	`escape_char /
LC_MONETARY
int_curr_symbol "USD "
END /
LC_MONETARY`, // Escape char redefinition
	`LC_MONETARY
int_curr_symbol     "USD "
currency_symbol     "$"
mon_decimal_point   "."
mon_thousands_sep   ","
mon_grouping        3;\
  3
positive_sign       ""
negative_sign       "-"
int_frac_digits     2
frac_digits         2
p_cs_precedes       1
int_p_sep_by_space  1
p_sep_by_space      0
n_cs_precedes       1
int_n_sep_by_space  1
n_sep_by_space      0
p_sign_posn         1
n_sign_posn         1
END LC_MONETARY `, // Monetary stanza
	`LC_TIME
abday "mon";"tue";"wed";"thu";"fri";"sat";"sun"
am_pm "AM";"PM"
era "+:2:1990/01/01:+*:平成:%EC%Ey年";\
    "+:1:1989/01/08:1989/12/31:平成:%EC元年";\
    "+:2:1927/01/01:1989/01/07:昭和:%EC%Ey年";\
    "+:1:1926/12/25:1926/12/31:昭和:%EC元年";\
    "+:2:1913/01/01:1926/12/24:大正:%EC%Ey年";\
    "+:2:1912/07/30:1912/12/31:大正:%EC元年";\
    "+:6:1873/01/01:1912/07/29:明治:%EC%Ey年";\
    "+:1:0001/01/01:1872/12/31:西暦:%EC%Ey年";\
    "+:1:-0001/12/31:-*:紀元前:%EC%Ey年"
era_d_fmt "%EY%m月%d日"
era_t_fmt "%H時%M分%S秒"
era_d_t_fmt "%EY%m月%d日 %H時%M分%S秒"
alt_digits "〇";"一";"二";"三";"四";"五";"六";"七";"八";"九";"十";"十一";"十二";"十三";"十四";"十五";"十六";"十七";"十八";"十九";"二十";"二十一";"二十二";"二十三";"二十四";"二十五";"二十六";"二十七";"二十八";"二十九";"三十";"三十一";"三十二";"三十三";"三十四";"三十五";"三十六";"三十七";"三十八";"三十九";"四十";"四十一";"四十二";"四十三";"四十四";"四十五";"四十六";"四十七";"四十八";"四十九";"五十";"五十一";"五十二";"五十三";"五十四";"五十五";"五十六";"五十七";"五十八";"五十九";"六十";"六十一";"六十二";"六十三";"六十四";"六十五";"六十六";"六十七";"六十八";"六十九";"七十";"七十一";"七十二";"七十三";"七十四";"七十五";"七十六";"七十七";"七十八";"七十九";"八十";"八十一";"八十二";"八十三";"八十四";"八十五";"八十六";"八十七";"八十八";"八十九";"九十";"九十一";"九十二";"九十三";"九十四";"九十五";"九十六";"九十七";"九十八";"九十九"
END LC_TIME`, // LC_TIME Era test
	`LC_CTYPE
# The following is the minimum POSIX locale LC_CTYPE.
# "alpha" is by definition "upper" and "lower"
# "alnum" is by definition "alpha" and "digit"
# "print" is by definition "alnum", "punct", and the <space>
# "graph" is by definition "alnum" and "punct"
#
upper    <A>;<B>;<C>;<D>;<E>;<F>;<G>;<H>;<I>;<J>;<K>;<L>;<M>;\
         <N>;<O>;<P>;<Q>;<R>;<S>;<T>;<U>;<V>;<W>;<X>;<Y>;<Z>
#
lower    <a>;<b>;<c>;<d>;<e>;<f>;<g>;<h>;<i>;<j>;<k>;<l>;<m>;\
         <n>;<o>;<p>;<q>;<r>;<s>;<t>;<u>;<v>;<w>;<x>;<y>;<z>
#
digit    <zero>;<one>;<two>;<three>;<four>;<five>;<six>;\
         <seven>;<eight>;<nine>
#
space    <tab>;<newline>;<vertical-tab>;<form-feed>;\
         <carriage-return>;<space>
#
cntrl    <alert>;<backspace>;<tab>;<newline>;<vertical-tab>;\
         <form-feed>;<carriage-return>;\
         <NUL>;<SOH>;<STX>;<ETX>;<EOT>;<ENQ>;<ACK>;<SO>;\
         <SI>;<DLE>;<DC1>;<DC2>;<DC3>;<DC4>;<NAK>;<SYN>;\
         <ETB>;<CAN>;<EM>;<SUB>;<ESC>;<IS4>;<IS3>;<IS2>;\
         <IS1>;<DEL>
#
punct    <exclamation-mark>;<quotation-mark>;<number-sign>;\
         <dollar-sign>;<percent-sign>;<ampersand>;<apostrophe>;\
         <left-parenthesis>;<right-parenthesis>;<asterisk>;\
         <plus-sign>;<comma>;<hyphen-minus>;<period>;<slash>;\
         <colon>;<semicolon>;<less-than-sign>;<equals-sign>;\
         <greater-than-sign>;<question-mark>;<commercial-at>;\
         <left-square-bracket>;<backslash>;<right-square-bracket>;\
         <circumflex>;<underscore>;<grave-accent>;<left-curly-bracket>;\
         <vertical-line>;<right-curly-bracket>;<tilde>
#
xdigit   <zero>;<one>;<two>;<three>;<four>;<five>;<six>;<seven>;\
         <eight>;<nine>;<A>;<B>;<C>;<D>;<E>;<F>;<a>;<b>;<c>;<d>;<e>;<f>
#
blank    <space>;<tab>
#
toupper (<a>,<A>);(<b>,<B>);(<c>,<C>);(<d>,<D>);(<e>,<E>);\
        (<f>,<F>);(<g>,<G>);(<h>,<H>);(<i>,<I>);(<j>,<J>);\
        (<k>,<K>);(<l>,<L>);(<m>,<M>);(<n>,<N>);(<o>,<O>);\
        (<p>,<P>);(<q>,<Q>);(<r>,<R>);(<s>,<S>);(<t>,<T>);\
        (<u>,<U>);(<v>,<V>);(<w>,<W>);(<x>,<X>);(<y>,<Y>);(<z>,<Z>)
#
tolower (<A>,<a>);(<B>,<b>);(<C>,<c>);(<D>,<d>);(<E>,<e>);\
        (<F>,<f>);(<G>,<g>);(<H>,<h>);(<I>,<i>);(<J>,<j>);\
        (<K>,<k>);(<L>,<l>);(<M>,<m>);(<N>,<n>);(<O>,<o>);\
        (<P>,<p>);(<Q>,<q>);(<R>,<r>);(<S>,<s>);(<T>,<t>);\
        (<U>,<u>);(<V>,<v>);(<W>,<w>);(<X>,<x>);(<Y>,<y>);(<Z>,<z>)
END LC_CTYPE`, // Ctype stanza
}

var parseFailTests = []string{
	`LC_TIME`, // No end stanza
	`LC_BUTT
END LC_BUTT`, // Fake category
	`LC_TIME
frac_digits 4
END LC_TIME`, // Invalid keyword for category
	`LC_TIME
END LC_TIME\`, // Trailing escape char
}

func TestShouldParse(t *testing.T) {
	for i, test := range parseTests {
		t.Run(fmt.Sprintf("Subtest %d", i), func(t *testing.T) {
			l := NewLexer(test)
			parsed := yyParse(l)
			if parsed != 0 {
				t.Errorf("yyParse returned %d", parsed)
			}
		})
	}
}

func TestParsePOSIX(t *testing.T) {
	f, err := os.Open("POSIX.locale")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		t.Error(err)
	}
	yyErrorVerbose = true
	l := NewLexer(string(bytes))
	parsed := yyParse(l)
	fmt.Println("parsed:", parsed)
	fmt.Printf("def: %#v\n", def)
}

func BenchmarkParsePOSIX(b *testing.B) {
	f, err := os.Open("POSIX.locale")
	if err != nil {
		b.Error(err)
	}
	bytes, err := ioutil.ReadAll(f)
	f.Close()
	if err != nil {
		b.Error(err)
	}
	input := string(bytes)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		yyParse(NewLexer(input))
	}
}
