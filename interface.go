package plag

// interface for all types
type pType interface {
	getHash() string
	getShort() string
	getLong() string
	getHelp() string
	getValue() interface{}
	getStatus() bool // set or not

	setDefault()                // set the value to default
	updateIsSet(bool)           // update the isSet value
	stringToValue(string) error // convert string to value
}
