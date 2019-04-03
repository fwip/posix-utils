%{
package locale

import (
  "strconv"
  "strings"
)
var def Def

// TODO: Why don't the tokens have the value I expect in tok.yys?
func assignMonetaryString (tok yySymType, value string) {
  switch tok.val {
  case "int_curr_symbol":
    def.monetary.intCurrSymbol = value
  }
}

func assignMonetaryNumber (tok int, value string) {
  n, err := strconv.Atoi(value)
  if err != nil {
    panic("Can't convert " + value + " to integer")
  }
  switch tok {
  case INT_FRAC_DIGITS_STR:
    def.monetary.intFracDigits = n
  default:
    //panic("NYI" + value)
  }
}

func assignMonetaryGrouping(tok int, value string) {
  vals := strings.Split(value, ";")
  g := make([]int, len(vals))
  var err error
  for i, v := range vals {
    g[i], err = strconv.Atoi(v)
    if err != nil {
      panic("Err: " + err.Error())
    }
  }
  def.monetary.monGrouping = g
}

%}
%union{
  val string
}

%token              LOC_NAME
%token              CHAR
%token              NUMBER
%token              COLLSYMBOL COLLELEMENT
%token              CHARSYMBOL OCTAL_CHAR HEX_CHAR DECIMAL_CHAR
%token              CHARCLASS
%token              ELLIPSIS
%token              EXTENDED_REG_EXP
%token              EOL

%token ESCAPE_CHAR_STR
%token COMMENT_CHAR_STR
%token COPY_STR
%token LC_CTYPE_STR
%token CHARCLASS_STR

%token UPPER_STR
%token LOWER_STR
%token ALPHA_STR
%token DIGIT_STR
%token PUNCT_STR
%token XDIGIT_STR
%token SPACE_STR
%token PRINT_STR
%token GRAPH_STR
%token BLANK_STR
%token CNTRL_STR
%token ALNUM_STR
%token TOUPPER_STR
%token TOLOWER_STR
%token COPY_STR
%token LC_COLLATE_STR
%token COLLATING_SYMBOL_STR
%token COLLATING_ELEMENT_STR
%token ORDER_START_STR
%token ORDER_START_STR
%token FORWARD_STR
%token BACKWARD_STR
%token POSITION_STR
%token UNDEFINED_STR
%token IGNORE_STR
%token END_STR
%token LC_TIME_STR
%token ERA_STR
%token ERA_D_FMT_STR
%token ERA_T_FMT_STR
%token ALT_DIGITS_STR
%token ERA_D_T_FMT_STR
%token AM_PM_STR
%token D_T_FMT_STR
%token D_FMT_STR
%token T_FMT_STR
%token T_FMT_AMPM_STR
%token ABDAY_STR
%token DAY_STR
%token ABMON_STR
%token MON_STR
%token LC_TIME_STR
%token END_STR
%token LC_NUMERIC_STR
%token DECIMAL_POINT_STR
%token THOUSANDS_SEP_STR
%token LC_NUMERIC_STR
%token COPY_STR
%token LC_MONETARY_STR
%token END_STR
%token INT_P_SIGN_POSN_STR
%token INT_N_CS_PRECEDES_STR
%token INT_P_CS_PRECEDES_STR
%token P_SIGN_POSN_STR
%token N_CS_PRECEDES_STR
%token P_CS_PRECEDES_STR
%token INT_FRAC_DIGITS_STR
%token FRAC_DIGITS_STR
%token P_SEP_BY_SPACE_STR
%token N_SEP_BY_SPACE_STR
%token N_SIGN_POSN_STR
%token INT_P_SEP_BY_SPACE_STR
%token INT_N_SEP_BY_SPACE_STR
%token INT_N_SIGN_POSN_STR
%token NEGATIVE_SIGN_STR
%token POSITIVE_SIGN_STR
%token MON_DECIMAL_POINT_STR
%token INT_CURR_SYMBOL_STR
%token CURRENCY_SYMBOL_STR
%token MON_THOUSANDS_SEP_STR
%token LC_MONETARY_STR
%token END_STR
%token LC_MESSAGES_STR
%token NOEXPR_STR
%token YESEXPR_STR
%token LC_MESSAGES_STR
%token END_STR
%token LC_CTYPE_STR
%token FROM_STR
%token ORDER_END_STR
%token END_STR
%token LC_COLLATE_STR
%token COPY_STR
%token COPY_STR
%token MON_GROUPING_STR
%token GROUPING_STR
%token COPY_STR
%token NEGATIVE_ONE_STR
%token TWO_DOUBLE_QUOTE_STR
%start              locale_definition

%%

locale_definition   : global_statements locale_categories
                    |                   locale_categories
                    ;


global_statements   : global_statements symbol_redefine
                    | symbol_redefine
                    ;


symbol_redefine     : ESCAPE_CHAR_STR  CHAR EOL
                    | COMMENT_CHAR_STR CHAR EOL
                    ;


locale_categories   : locale_categories locale_category
                    | locale_category
                    ;


locale_category     : lc_ctype | lc_collate | lc_messages
                    | lc_monetary | lc_numeric | lc_time
                    ;


/* The following grammar rules are common to all categories */


char_list           : char_list char_symbol {$$.val = $1.val + string($2.val)}
                    | char_symbol
                    ;


char_symbol         : CHAR | CHARSYMBOL
                    | OCTAL_CHAR | HEX_CHAR | DECIMAL_CHAR
                    ;


elem_list           : elem_list char_symbol
                    | elem_list COLLSYMBOL
                    | elem_list COLLELEMENT
                    | char_symbol
                    | COLLSYMBOL
                    | COLLELEMENT
                    ;


symb_list           : symb_list COLLSYMBOL
                    | COLLSYMBOL
                    ;


locale_name         : LOC_NAME
                    | '"' LOC_NAME '"'
                    ;


/* The following is the LC_CTYPE category grammar */


lc_ctype            : ctype_hdr ctype_keywords         ctype_tlr
                    | ctype_hdr COPY_STR locale_name EOL ctype_tlr
                    ;


ctype_hdr           : LC_CTYPE_STR EOL
                    ;


ctype_keywords      : ctype_keywords ctype_keyword
                    | ctype_keyword
                    ;


ctype_keyword       : charclass_keyword charclass_list EOL
                    | charconv_keyword charconv_list EOL
                    | CHARCLASS_STR charclass_namelist EOL
                    ;


charclass_namelist  : charclass_namelist ';' CHARCLASS
                    | CHARCLASS
                    ;


charclass_keyword   : UPPER_STR | LOWER_STR | ALPHA_STR | DIGIT_STR
                    | PUNCT_STR | XDIGIT_STR | SPACE_STR | PRINT_STR
                    | GRAPH_STR | BLANK_STR | CNTRL_STR | ALNUM_STR
                    | CHARCLASS
                    ;


charclass_list      : charclass_list ';' char_symbol
                    | charclass_list ';' ELLIPSIS ';' char_symbol
                    | char_symbol
                    ;


charconv_keyword    : TOUPPER_STR
                    | TOLOWER_STR
                    ;


charconv_list       : charconv_list ';' charconv_entry
                    | charconv_entry
                    ;


charconv_entry      : '(' char_symbol ',' char_symbol ')'
                    ;


ctype_tlr           : END_STR LC_CTYPE_STR EOL
                    ;


/* The following is the LC_COLLATE category grammar */


lc_collate          : collate_hdr collate_keywords       collate_tlr
                    | collate_hdr COPY_STR locale_name EOL collate_tlr
                    ;


collate_hdr         : LC_COLLATE_STR EOL
                    ;


collate_keywords    :                order_statements
                    | opt_statements order_statements
                    ;


opt_statements      : opt_statements collating_symbols
                    | opt_statements collating_elements
                    | collating_symbols
                    | collating_elements
                    ;


collating_symbols   : COLLATING_SYMBOL_STR COLLSYMBOL EOL
                    ;


collating_elements  : COLLATING_ELEMENT_STR COLLELEMENT
                    | FROM_STR '"' elem_list '"' EOL
                    ;


order_statements    : order_start collation_order order_end
                    ;


order_start         : ORDER_START_STR EOL
                    | ORDER_START_STR order_opts EOL
                    ;


order_opts          : order_opts ';' order_opt
                    | order_opt
                    ;


order_opt           : order_opt ',' opt_word
                    | opt_word
                    ;


opt_word            : FORWARD_STR | BACKWARD_STR | POSITION_STR
                    ;


collation_order     : collation_order collation_entry
                    | collation_entry
                    ;


collation_entry     : COLLSYMBOL EOL
                    | collation_element weight_list EOL
                    | collation_element             EOL
                    ;


collation_element   : char_symbol
                    | COLLELEMENT
                    | ELLIPSIS
                    | UNDEFINED_STR
                    ;


weight_list         : weight_list ';' weight_symbol
                    | weight_list ';'
                    | weight_symbol
                    ;


weight_symbol       : /* empty */
                    | char_symbol
                    | COLLSYMBOL
                    | '"' elem_list '"'
                    | '"' symb_list '"'
                    | ELLIPSIS
                    | IGNORE_STR
                    ;


order_end           : ORDER_END_STR EOL
                    ;


collate_tlr         : END_STR LC_COLLATE_STR EOL
                    ;


/* The following is the LC_MESSAGES category grammar */


lc_messages         : messages_hdr messages_keywords      messages_tlr
                    | messages_hdr COPY_STR locale_name EOL messages_tlr
                    ;


messages_hdr        : LC_MESSAGES_STR EOL
                    ;


messages_keywords   : messages_keywords messages_keyword
                    | messages_keyword
                    ;


messages_keyword    : YESEXPR_STR '"' EXTENDED_REG_EXP '"' EOL
                    | NOEXPR_STR  '"' EXTENDED_REG_EXP '"' EOL
                    ;


messages_tlr        : END_STR LC_MESSAGES_STR EOL
                    ;


/* The following is the LC_MONETARY category grammar */


lc_monetary         : monetary_hdr monetary_keywords       monetary_tlr
                    | monetary_hdr COPY_STR locale_name EOL  monetary_tlr
                    ;


monetary_hdr        : LC_MONETARY_STR EOL
                    ;


monetary_keywords   : monetary_keywords monetary_keyword
                    | monetary_keyword
                    ;


monetary_keyword    : mon_keyword_string mon_string EOL       {assignMonetaryString($1, $2.val)}
                    | mon_keyword_char NUMBER EOL             {assignMonetaryNumber($1.yys, $2.val)}
                    | mon_keyword_char NEGATIVE_ONE_STR   EOL {assignMonetaryNumber($1.yys, "-1")}
                    | mon_keyword_grouping mon_group_list EOL {assignMonetaryGrouping($1.yys, $2.val)}
                    ;


mon_keyword_string  : INT_CURR_SYMBOL_STR | CURRENCY_SYMBOL_STR
                    | MON_DECIMAL_POINT_STR | MON_THOUSANDS_SEP_STR
                    | POSITIVE_SIGN_STR | NEGATIVE_SIGN_STR
                    ;


mon_string          : '"' char_list '"' { $$ = $2 }
                    | TWO_DOUBLE_QUOTE_STR
                    ;


mon_keyword_char    : INT_FRAC_DIGITS_STR | FRAC_DIGITS_STR
                    | P_CS_PRECEDES_STR | P_SEP_BY_SPACE_STR
                    | N_CS_PRECEDES_STR | N_SEP_BY_SPACE_STR
                    | P_SIGN_POSN_STR | N_SIGN_POSN_STR
                    | INT_P_CS_PRECEDES_STR | INT_P_SEP_BY_SPACE_STR
                    | INT_N_CS_PRECEDES_STR | INT_N_SEP_BY_SPACE_STR
                    | INT_P_SIGN_POSN_STR | INT_N_SIGN_POSN_STR
                    ;


mon_keyword_grouping : MON_GROUPING_STR
                    ;


mon_group_list      : NUMBER
                    | mon_group_list ';' NUMBER { $$.val = $1.val + ";" + $3.val }
                    ;


monetary_tlr        : END_STR LC_MONETARY_STR EOL
                    ;


/* The following is the LC_NUMERIC category grammar */


lc_numeric          : numeric_hdr numeric_keywords       numeric_tlr
                    | numeric_hdr COPY_STR locale_name EOL numeric_tlr
                    ;


numeric_hdr         : LC_NUMERIC_STR EOL
                    ;


numeric_keywords    : numeric_keywords numeric_keyword
                    | numeric_keyword
                    ;


numeric_keyword     : num_keyword_string num_string EOL
                    | num_keyword_grouping num_group_list EOL
                    ;


num_keyword_string  : DECIMAL_POINT_STR
                    | THOUSANDS_SEP_STR
                    ;


num_string          : '"' char_list '"'
                    | TWO_DOUBLE_QUOTE_STR
                    ;


num_keyword_grouping: GROUPING_STR
                    ;


num_group_list      : NUMBER
                    | num_group_list ';' NUMBER
                    ;


numeric_tlr         : END_STR LC_NUMERIC_STR EOL
                    ;


/* The following is the LC_TIME category grammar */


lc_time             : time_hdr time_keywords          time_tlr
                    | time_hdr COPY_STR locale_name EOL time_tlr
                    ;


time_hdr            : LC_TIME_STR EOL
                    ;


time_keywords       : time_keywords time_keyword
                    | time_keyword
                    ;


time_keyword        : time_keyword_name time_list EOL
                    | time_keyword_fmt time_string EOL
                    | time_keyword_opt time_list EOL
                    ;


time_keyword_name   : AM_PM_STR | ABDAY_STR | DAY_STR | ABMON_STR | MON_STR
                    ;


time_keyword_fmt    : D_T_FMT_STR | D_FMT_STR | T_FMT_STR
                    |  T_FMT_AMPM_STR
                    ;


time_keyword_opt    : ERA_STR | ERA_D_FMT_STR | ERA_T_FMT_STR
                    | ERA_D_T_FMT_STR | ALT_DIGITS_STR
                    ;


time_list           : time_list ';' time_string
                    | time_string
                    ;


time_string         : '"' char_list '"'
                    ;


time_tlr            : END_STR LC_TIME_STR EOL
                    ;
