package ed

// Command is magick
type Command struct {
	typ    cmdType
	start  address
	end    address
	dest   address
	text   string
	params []string
}

type cmdType byte

//
const (
	ctnull cmdType = iota
	ctappend
	ctchange
	ctdelete
	ctedit
	cteditForce
	ctfilename
	ctglobal
	ctinteractive
	cthelp
	cthelpMode
	ctinsert
	ctjoin
	ctmark
	ctlist
	ctmove
	ctnumber
	ctprint
	ctprompt
	ctquit
	ctquitForce
	ctread
	ctsubstitute
	ctcopy
	ctundo
	ctglobalInverse
	ctinteractiveInverse
	ctwrite
	ctlineNumber
	ctshell
)
