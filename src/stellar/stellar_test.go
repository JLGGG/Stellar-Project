package stellar

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
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
	if ExistFile(fn) == false {
		WriteFile(fn, []byte("hello world"))
	} else {
		log.Println(string(ReadFile(fn)))
	}
}

func TestSQL(t *testing.T) {
	db, err := sql.Open("mysql", "")
	CheckError(err)
	defer db.Close()

	result, err := db.Exec("INSERT INTO entity_table (ID, PW) VALUES (?, ?)", "hello", "world")
	CheckError(err)

	n, err := result.RowsAffected()
	log.Printf("%d row inserted\n", n)
}
