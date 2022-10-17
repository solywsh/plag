package plag

import (
	"fmt"
	"testing"
)

func TestArgsInt(t *testing.T) {
	a := Int("a", "an", "command a test line", 1)
	c := NewCmd("test", "test help")
	// fmt.Println(a.getHelp())
	fmt.Println(c.Help())
	c.Set(a)
	c.Parse("test -an 2")
	fmt.Println(c.Exist(a))
	fmt.Println(c.Get("a"))
	fmt.Println(c.Clean().Get("an"))

	isHelp, help, err := c.Parse("test -h")
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	if isHelp {
		fmt.Println(help)
	}
}
