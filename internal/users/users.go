package users

import (
	"database/sql"
	"log"

	"github.com/aniruddha2000/hackernews/internal/pkg/db/migrations/mysql"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (user *Users) Create() int64 {
	statement, err := mysql.Db.Prepare("INSERT INTO Users(Username, Password) VALUES(?,?)")
	if err != nil {
		log.Fatal(err)
	}
	log.Print(statement)
	defer statement.Close()

	hashedPassword, err := HashedPassword(user.Password)
	if err != nil {
		log.Fatal(err)
	}
	res, err := statement.Exec(user.Username, hashedPassword)
	if err != nil {
		log.Fatal(err)
	}

	insertID, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	log.Print("User Created")
	return insertID
}

func HashedPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetUserIdByUsername(username string) (int64, error) {
	statement, err := mysql.Db.Prepare("SELECT ID FROM Users WHERE Username=?")
	if err != nil {
		log.Fatal(err)
	}
	defer statement.Close()

	var ID int64
	err = statement.QueryRow(username).Scan(&ID)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("No Rows Selected : %v", err)
		}
		return 0, err
	}
	return ID, nil
}
