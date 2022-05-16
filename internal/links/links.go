package links

import (
	"database/sql"
	"errors"
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

// Save the entry in the DB
func (link *Link) Save() int64 {
	statement2, err := database.Db.Prepare("INSERT INTO Links(Title,Address, UserID) VALUES(?,?,?)")
	if err != nil {
		log.Fatal(err)
	}

	res, err := statement2.Exec(link.Title, link.Address, link.User.ID)
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

// Gat all links for a specified user
func GetAll() []Link {
	const query = `
	SELECT L.ID, L.Title, L.Address, L.UserID, U.Username
	FROM Links L
	INNER JOIN Users U on L.UserID = U.ID
	`

	statement2, err := database.Db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer statement2.Close()

	rows, err := statement2.Query()
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

// Get the link & throws error if the link
// don't belong to the requested by the user
func Get(id string, username string) (Link, error) {
	const query = `
	SELECT L.ID, L.Title, L.Address, U.Username
	FROM Links L
	INNER JOIN Users U on L.UserID = U.ID
	WHERE L.ID=?
	`

	statement2, err := database.Db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer statement2.Close()

	var (
		link        Link
		usernameGot string
	)
	err = statement2.QueryRow(id).Scan(&link.ID, &link.Title, &link.Address, &usernameGot)
	if err != nil {
		if err == sql.ErrNoRows {
			return Link{}, errors.New("no rows selected, check id")
		} else {
			log.Fatal(err)
		}
	}

	if usernameGot != username {
		return Link{}, errors.New("link don't belong to the user")
	}

	return link, nil
}

// Check a link exists or not if yes Update the link
// and throws error if the link don't belong to the user
func (link *Link) Update(id string) (int64, error) {
	const (
		queryUsername = `
		SELECT U.Username
		FROM Links L
		INNER JOIN Users U on L.UserID = U.ID
		WHERE L.ID=?
		`

		queryUpdate = `
		UPDATE Links SET Title=? , Address=? WHERE ID=?
		`
	)

	statement, err := database.Db.Prepare(queryUsername)
	if err != nil {
		log.Fatal(err)
	}
	defer statement.Close()

	var username string
	err = statement.QueryRow(id).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("no rows selected, check id")
		} else {
			log.Fatal(err)
		}
	}

	if username != link.User.Username {
		return 0, errors.New("link don't belong to the user")
	}

	statement2, err := database.Db.Prepare(queryUpdate)
	if err != nil {
		log.Fatal(err)
	}
	defer statement2.Close()

	res, err := statement2.Exec(link.Title, link.Address, id)
	if err != nil {
		log.Fatal(err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	return rowsAffected, nil
}

// Check a link exists or not if yes Delete the link
// and throws error if the link don't belong to the user
func Delete(id string, username string) (int64, error) {
	const (
		queryUsername = `
		SELECT U.Username
		FROM Links L
		INNER JOIN Users U on L.UserID = U.ID
		WHERE L.ID=?
		`
		queryDelete = `
		DELETE FROM Links WHERE ID=?
		`
	)

	statement, err := database.Db.Prepare(queryUsername)
	if err != nil {
		log.Fatal(err)
	}
	defer statement.Close()

	var usernameGot string
	err = statement.QueryRow(id).Scan(&usernameGot)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("no rows selected, check id")
		} else {
			log.Fatal(err)
		}
	}

	if usernameGot != username {
		return 0, errors.New("link don't belong to the user")
	}

	statement2, err := database.Db.Prepare(queryDelete)
	if err != nil {
		log.Fatal(err)
	}
	defer statement2.Close()

	res, err := statement2.Exec(id)
	if err != nil {
		log.Fatal(err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	return rowsAffected, nil
}
