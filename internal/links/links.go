package links

import (
	"log"

	database "github.com/aniruddha2000/hackernews/internal/pkg/db/migrations/mysql"
	"github.com/aniruddha2000/hackernews/internal/users"
)

type Links struct {
	ID      string
	Title   string
	Address string
	User    *users.Users
}

func (link Links) Save() int64 {
	stmt, err := database.Db.Prepare("INSERT INTO Links(Title, Address) VALUES(?,?)")
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(link.Title, link.Address)
	if err != nil {
		log.Fatal(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
	log.Print("Row inserted")
	return id
}

func GetAll() (links []Links) {
	stmt, err := database.Db.Prepare("SELECT ID, Title, Address FROM Links")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var link Links
		err = rows.Scan(&link.ID, &link.Title, &link.Address)
		if err != nil {
			log.Fatal(err)
		}
		links = append(links, link)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return links
}

func Get(id string) Links {
	var link Links
	stmt, err := database.Db.Prepare("SELECT ID, Title, Address FROM Links WHERE ID=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&link.ID, &link.Title, &link.Address)
	if err != nil {
		log.Fatal(err)
	}
	return link
}

func (link Links) Update(id string) int64 {
	stmt, err := database.Db.Prepare("UPDATE Links SET Title=? , Address=? WHERE ID=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(link.Title, link.Address, id)
	if err != nil {
		log.Fatal(err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	return rowsAffected
}

func Delete(id string) int64 {
	stmt, err := database.Db.Prepare("DELETE FROM Links WHERE ID=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		log.Fatal(err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	return rowsAffected
}
