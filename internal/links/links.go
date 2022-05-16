package links

import (
	"log"

	database "github.com/aniruddha2000/hackernews/internal/pkg/db/migrations/mysql"
	"github.com/aniruddha2000/hackernews/internal/users"
)

type Link struct {
	ID      string
	Title   string
	Address string
	User    *users.User
}

func (link *Link) Save() int64 {
	stmt, err := database.Db.Prepare("INSERT INTO Links(Title,Address, UserID) VALUES(?,?,?)")
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(link.Title, link.Address, link.User.ID)
	if err != nil {
		log.Fatal(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal("Error:", err.Error())
	}

	log.Print("Row inserted!")
	return id
}

func GetAll() []Link {
	const query = `
	SELECT L.ID, L.Title, L.Address, L.UserID, U.Username
	FROM Links L
	INNER JOIN Users U on L.UserID = U.ID
	`

	stmt, err := database.Db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var (
		links    []Link
		username string
		id       string
	)

	for rows.Next() {
		var link Link
		err := rows.Scan(&link.ID, &link.Title, &link.Address, &id, &username)
		if err != nil {
			log.Fatal(err)
		}
		link.User = &users.User{
			ID:       id,
			Username: username,
		}
		links = append(links, link)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return links
}
