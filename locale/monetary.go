package locale

// Monetary is the monetary stuff
type Monetary struct {
	intCurrSymbol   string
	currencySymbol  string
	monDecimalPoint string
	monThousandsSep string
	monGrouping     []int
	positiveSign    string
	negativeSign    string
	intFracDigits   int
	fracDigits      int
	pCsPrecedes     int
	pSepBySpace     int
	nCsPrecedes     int
	nSepBySpace     int
	pSignPosn       int
	nSignPosn       int
	intPcsPrecedes  int
	intPsepBySpace  int
	intNcsPrecedes  int
	intNsepBySpace  int
	intPsignPosn    int
	intNsignPosn    int
}

// EmptyMonetary returns a monetary object for which nothing is initialized
// -1 is the official "unset" sentinel value for locale
// TODO: Should we replace these with pointers (to allow nil and thereby provide a sensible empty struct)?
func EmptyMonetary() Monetary {
	return Monetary{
		"",
		"",
		"",
		"",
		nil,
		"",
		"",
		-1,
		-1,
		-1,
		-1,
		-1,
		-1,
		-1,
		-1,
		-1,
		-1,
		-1,
		-1,
		-1,
		-1,
	}
}
