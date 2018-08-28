package test

import (
	"testing"
	"github.com/charnlsxy/gotools/forjava"
	"fmt"
)

func TestJson2java(t *testing.T) {
	str :=``
	r,_ := forjava.Json2java([]byte(str),"")
	fmt.Println(r)
}
