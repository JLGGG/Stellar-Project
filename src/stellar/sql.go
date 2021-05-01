package stellar

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func WriteTxToDB(id, pw string) {
	// Start the SQL server before using SQL: mysql.server start
	// Turn off SQL server: mysql.server stop
	// Keep the sql connection string private.
	db, err := sql.Open("mysql", "root:1111@tcp(127.0.0.1:3306)/testdb")
	CheckError(err)
	defer db.Close()

	result, err := db.Exec("INSERT INTO entity_table (ID, PW) VALUES (?, ?)", id, pw)
	CheckError(err)

	n, err := result.RowsAffected()
	log.Printf("%d row inserted\n", n)
}

func WriteFavoriteAccountToDB(id string) {
	db, err := sql.Open("mysql", "root:1111@tcp(127.0.0.1:3306)/testdb")
	CheckError(err)
	defer db.Close()

	result, err := db.Exec("INSERT INTO favorite_table (ID) VALUES (?)", id)
	CheckError(err)

	n, err := result.RowsAffected()
	log.Printf("%d row inserted\n", n)
}

func DeleteFavoriteAccountFromDB(s string) {
	db, err := sql.Open("mysql", "root:1111@tcp(127.0.0.1:3306)/testdb")
	CheckError(err)
	defer db.Close()

	result, err := db.Exec("DELETE from favorite_table where id=(?)", s)
	CheckError(err)

	n, err := result.RowsAffected()
	log.Printf("%d row deleted\n", n)
}

func ReadTxFromDB(pID, pPW []string) (count int) {
	// Keep the sql connection string private.
	db, err := sql.Open("mysql", "root:1111@tcp(127.0.0.1:3306)/testdb")
	CheckError(err)
	defer db.Close()

	rows, err := db.Query("SELECT id, pw FROM entity_table")
	CheckError(err)
	defer rows.Close()

	var id, pw string
	count = 0
	for rows.Next() {
		err := rows.Scan(&id, &pw)
		CheckError(err)
		pID[count] = id
		pPW[count] = pw
		count += 1
	}
	return
}

func ReadFavoriteAccountFromDB(pID []string) (count int) {
	// Keep the sql connection string private.
	db, err := sql.Open("mysql", "root:1111@tcp(127.0.0.1:3306)/testdb")
	CheckError(err)
	defer db.Close()

	rows, err := db.Query("SELECT id FROM favorite_table")
	CheckError(err)
	defer rows.Close()

	var id string
	count = 0
	for rows.Next() {
		err := rows.Scan(&id)
		CheckError(err)
		pID[count] = id
		count += 1
	}
	return
}
