package main

import (
	"log"
	"testing"
)

func TestParseBalanceStr(t *testing.T) {
	testStr := "{10000.0000000  0.0000000 0.0000000  %!s(uint32=0) %!s(*bool=<nil>) %!s(*bool=<nil>) %!s(*bool=<nil>) {native  }}"
	testStr = ParseBalanceStr(testStr)

	log.Printf("Parsing result: %s", testStr)
}

func TestFileFunc(t *testing.T) {
	fn := "hello.txt"
	if existFile(fn) == false {
		writeFile(fn, []byte("hello world"))
	} else {
		log.Println(string(readFile(fn)))
	}
}
