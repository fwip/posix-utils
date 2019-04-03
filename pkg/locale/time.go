package locale

import (
	"fmt"
	"strconv"
	"strings"
)

// Time defines the interpretation of the conversion specifications supported
// by the date utility and shall affect the behavior of the strftime(),
// wcsftime(), strptime(), and nl_langinfo() functions.
type Time struct {
	// TODO: Should these be fixed-length arrays?
	abday     []string
	day       []string
	abmon     []string
	mon       []string
	dtFmt     string
	dFmt      string
	tFmt      string
	am        string
	pm        string
	tFmtAmPm  string
	eras      []Era
	eraDTFmt  string
	eraDFmt   string
	eraTFmt   string
	altDigits []string
}

// Era defines how years are counted and displayed for each era in a locale
type Era struct {
	reverse   bool
	offset    int
	startDate string // String?? Use go date??
	endDate   string
	name      string
	format    string
}

// NewEra creates an era from an input string, as specified in Chapter 7 of the Base Definitions
// This format is: direction:offset:start_date:end_date:era_name:era_format
// direction:  +/-
// offset:     int
// start_date: yyyy/mm/dd
// end_date:   yyyy/mm/dd
// era_name:   string
// era_format: string
func NewEra(input string) (Era, error) {
	parts := strings.Split(input, ":")
	era := Era{}
	if len(parts) < 6 {
		return Era{}, fmt.Errorf("not enough fields in input: %s", input)
	}

	// Direction
	switch parts[0] {
	case "+":
	case "-":
		era.reverse = true
	default:
		return Era{}, fmt.Errorf("direction '%s' must be '+' or '-'", parts[0])
	}

	// Offset
	offset, err := strconv.Atoi(parts[1])
	if err != nil {
		return Era{}, nil
	}
	era.offset = offset

	// Dates
	// TODO: Check formatting?
	era.startDate = parts[2]
	era.endDate = parts[3]

	// Name
	era.name = parts[4]

	// Format
	era.format = strings.Join(parts[5:], ":")

	return era, nil
}
