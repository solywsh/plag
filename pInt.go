package plag

import (
	uuid "github.com/satori/go.uuid"
	"strconv"
)

// todo all of operation to lock/unlock

type pInt struct {
	hash   string // Hash of the flag
	short  string // Short name of the argument
	long   string // Long name of the argument
	help   string // Help text for the argument
	dValue int    // Default value for the argument
	value  int    // Value of the argument
	isSet  bool   // Whether the argument has been set
}

func Int(short, long, help string, dValue int) *pInt {
	return &pInt{
		hash:   uuid.NewV4().String(),
		short:  short,
		long:   long,
		help:   help,
		value:  dValue,
		dValue: dValue,
		isSet:  false,
	}
}

func (i *pInt) setDefault() {
	i.value = i.dValue
}

func (i *pInt) stringToValue(s string) error {
	v, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	i.value = v
	return nil
}

func (i *pInt) getHash() string {
	return i.hash
}

func (i *pInt) getShort() string {
	return i.short
}

func (i *pInt) getLong() string {
	return i.long
}

func (i *pInt) getHelp() string {
	return "-" + i.short + "\t" + "--" + i.long + "\t" + i.help + "\n"
}

func (i *pInt) getValue() interface{} {
	return i.value
}

func (i *pInt) getStatus() bool {
	return i.isSet
}

func (i *pInt) updateIsSet(set bool) {
	i.isSet = set
}
