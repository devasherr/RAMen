package server

import (
	"fmt"
	"strings"

	"github.com/Rohit-Dnath/RAMen/internal/store"
)

func (c *conn) cmdCopy(args []string) error {
	if len(args) < 3 || len(args) > 4 {
		return c.wrongArgs("copy")
	}

	src, dst := args[1], args[2]
	replace := len(args) == 4 && strings.ToUpper(args[3]) == "REPLACE"

	v, ok, err := c.s.store.Get(src)
	if err != nil {
		return c.storeErr(err)
	}
	if !ok {
		return c.storeErr(fmt.Errorf("invalid source to copy from"))
	}

	_, ok, err = c.s.store.Get(dst)
	if err != nil {
		return c.storeErr(err)
	}
	if ok && replace == false {
		return c.storeErr(fmt.Errorf("0"))
	}

	if c.s.store.Set(dst, v, store.SetOptions{}) {
		return c.writeSimple("1")
	}
	return c.writeNull()
}

func (c *conn) cmdObjectEncoding(args []string) error {
	if len(args) != 3 {
		return c.wrongArgs("object encoding")
	}

	// currently behaves similar to `type` command because store doesn't use different internal representation
	return c.writeSimple(c.s.store.Type(args[2]))
}
