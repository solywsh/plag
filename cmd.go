package plag

import (
	"errors"
	cmap "github.com/orcaman/concurrent-map"
	"strings"
)

type Cmd struct {
	name      string             // Name of the command
	help      string             // Help text for the command
	isParse   bool               // Whether the command has been parsed
	flags     cmap.ConcurrentMap // Flags for the command
	hashTable map[string]string  // Hash table for the flags
	splitFlag string             // Flag to split the command line, default is " "
	command   string             // Command line
	params    *params            // Parameters for the command
}

type params struct {
	value    string
	types    paramsTypes // types of the parameter
	next     *params
	previous *params
}

type paramsTypes int

const (
	paramsMain  paramsTypes = iota // command name
	paramsKey                      // key
	paramsValue                    // value
)

func (p *paramsTypes) String() string {
	switch *p {
	case paramsMain:
		return "main"
	case paramsKey:
		return "key"
	case paramsValue:
		return "value"
	}
	return "unknown"
}

func NewCmd(name, help string) *Cmd {
	return &Cmd{
		name:      name,
		help:      help,
		flags:     cmap.New(),
		splitFlag: " ",
		hashTable: make(map[string]string),
	}
}

func (c *Cmd) Set(args ...interface{}) {
	for _, arg := range args {
		switch arg.(type) {
		case pType:
			c.flags.Set(arg.(pType).getHash(), arg)
			c.hashTable[arg.(pType).getShort()] = arg.(pType).getHash()
			c.hashTable[arg.(pType).getLong()] = arg.(pType).getHash()
		}
	}
}

func (c *Cmd) SetSplitFlag(flag string) {
	c.splitFlag = flag
}

func (c *Cmd) Parse(cmdLine string) (isHelp bool, help string, err error) {

	c.params = nil
	for _, v := range c.flags.Items() {
		switch v.(type) {
		case pType:
			v.(pType).setDefault()
			v.(pType).updateIsSet(false)
		}
	}
	c.command = cmdLine
	args := strings.Split(cmdLine, c.splitFlag)
	if len(args) < 0 {
		return false, "", errors.New("invalid command line")
	} else if args[0] != c.name {
		return false, "", errors.New("invalid command name")
	}
	c.formatParams(args)
	c.isParse = true
	return false, "", nil
}

func (c *Cmd) Exist(arg interface{}) bool {
	switch arg.(type) {
	case pType:
		// type interface
		if res, ok := c.flags.Get(arg.(pType).getHash()); ok {
			return res.(pType).getStatus()
		}
	case string:
		// short or long
		if res, ok := c.flags.Get(c.hashTable[arg.(string)]); ok {
			return res.(pType).getStatus()
		}
		// hash
		if res, ok := c.flags.Get(arg.(string)); ok {
			return res.(pType).getStatus()
		}
		return false
	}
	return false
}

func (c *Cmd) Get(arg interface{}) (bool, interface{}) {
	switch arg.(type) {
	case pType:
		// type interface
		if res, ok := c.flags.Get(arg.(pType).getHash()); ok {
			return res.(pType).getStatus(), res.(pType).getValue()
		}
	case string:
		// short or long
		if res, ok := c.flags.Get(c.hashTable[arg.(string)]); ok {
			return res.(pType).getStatus(), res.(pType).getValue()
		}
		// hash
		if res, ok := c.flags.Get(arg.(string)); ok {
			return res.(pType).getStatus(), res.(pType).getValue()
		}
		return false, nil
	}
	return false, nil
}

func (c *Cmd) Clean() *Cmd {
	c.isParse = false
	c.command = ""
	c.params = nil
	for _, i := range c.flags.Items() {
		switch i.(type) {
		case pType:
			i.(pType).setDefault()
			i.(pType).updateIsSet(false)
		}
	}
	return c
}

func (c *Cmd) formatParams(args []string) {
	var pLast *params
	for i, arg := range args {
		if i == 0 && arg == c.name {
			pNew := new(params)
			pNew.value = arg
			pNew.types = paramsMain
			c.params = pNew
			pNew.previous = nil

			pLast = pNew
			continue
		}
		if len(arg) > 0 && arg[0] == '-' && pLast.types != paramsKey {
			pNew := new(params)
			pNew.value = arg[1:]
			// hashTable not exist
			if c.hashTable[pNew.value] == "" {
				continue
			}
			pNew.types = paramsKey
			pNew.previous = pLast
			pLast.next = pNew
			pLast = pNew
			continue
		}

		if len(arg) > 0 && pLast.types == paramsKey {
			pNew := new(params)
			pNew.value = arg
			ht := c.hashTable[pLast.value]
			if res, ok := c.flags.Get(ht); ok {
				switch res.(type) {
				case pType:
					err := res.(pType).stringToValue(pNew.value)
					if err != nil {
						continue
					}
					res.(pType).updateIsSet(true)
				}
			}
			pNew.types = paramsValue
			pNew.previous = pLast
			pLast.next = pNew
			pLast = pNew
			continue
		}
	}
}

func (c *Cmd) parseHelp(cmdLine string) (isHelp bool, help string) {
	if !strings.Contains(cmdLine, "-h") ||
		strings.Contains(cmdLine, "help") {
		return false, ""
	}
	// TODO: parse help
	return false, ""
}
