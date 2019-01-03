package test

import (
	"io/ioutil"
	"me/bsd/utils"
	"os"
	"testing"
	"github.com/charnlsxy/gotools/forjava"
	"fmt"
)

func TestJson2java(t *testing.T) {
	fd, err := os.Open("xxx")

	if utils.ErrNotNil(err) {
		return
	}

	bytes, err := ioutil.ReadAll(fd)
	if utils.ErrNotNil(err) {
		return
	}

	s, err := forjava.Json2java(bytes, "")
	if utils.ErrNotNil(err) {
		return
	}

	fmt.Println(s)

}
