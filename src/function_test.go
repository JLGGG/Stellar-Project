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
	currentDirectory()
	file := createFile("test.txt")

	testStr := make([]string, 2)
	testStr[0] = "hello "
	testStr[1] = "world\n"
	writeFile(file, testStr)
	str := string(readFile(file))
	log.Println(str)

}
