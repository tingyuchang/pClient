package pClient

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"strings"
)

var (
	Hostname = ""
	Port     = 5432
	Username = ""
	Password = ""
	Database = ""
)

func openConnection() (*sql.DB, error) {
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", Hostname, Port, Username, Password, Database)
	db, err := sql.Open("postgres", conn)

	if err != nil {
		return nil, err
	}
	return db, nil
}

// exists returns userID of username, returns -1 if not exist
func exists(username string) int {
	username = strings.ToLower(username)
	db, err := openConnection()
	if err != nil {
		return -1
	}
	defer db.Close()

	userId := -1
	queryString := fmt.Sprintf(`SELECT id FROM users WHERE username='%s'`, username)
	rows, err := db.Query(queryString)
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return -1
		}
		userId = id
	}
	defer rows.Close()
	return userId
}

//Adduser returns created UserId if failed returns -1
func Adduser(user Userdata) int {
	user.Name = strings.ToLower(user.Name)
	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer db.Close()
	userID := exists(user.Name)
	if userID != -1 {
		// userName existed
		fmt.Println("userName existed")
		return -1
	}

	insertStatement := `INSERT INTO "users"("username") VALUES($1)`

	_, err = db.Exec(insertStatement, user.Name)

	if err != nil {
		fmt.Println(err)
		return -1
	}

	userID = exists(user.Name)

	// insert Userdata
	insertStatement = `INSERT INTO "userdata"("userid", "name", "surname", "description") values($1, $2, $3, $4)`
	_, err = db.Exec(insertStatement, userID, user.Name, user.Surname, user.Description)
	if err != nil {
		fmt.Println(err)
	}
	return userID
}

func DeleteUser(id int) error {
	db, err := openConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	findStatement := fmt.Sprintf("SELECT username FROM users WHERE id=%d", id)
	rows, err := db.Query(findStatement)
	if err != nil {
		return err
	}
	var username string
	for rows.Next() {
		err := rows.Scan(&username)
		if err != nil {
			return err
		}
	}
	defer rows.Close()
	if exists(username) != id {
		return fmt.Errorf("user id does not exists")
	}

	deleteStatement := `DELETE FROM userdata WHERE userid=$1`
	_, err = db.Exec(deleteStatement, id)
	if err != nil {
		return err
	}
	return nil
}

func Listusers() ([]Userdata, error) {
	data := []Userdata{}

	db, err := openConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	queryStatement := fmt.Sprintf("SELECT id, username, name, surname, description FROM users, userdata WHERE users.id = userdata.userid")

	rows, err := db.Query(queryStatement)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id int
		var name string
		var username string
		var surname string
		var description string
		err := rows.Scan(&id, &username, &name, &surname, &description)
		if err != nil {
			return nil, err
		}
		temp := Userdata{
			id, username, name, surname, description,
		}
		data = append(data, temp)
	}

	defer rows.Close()

	return data, nil
}

func Updateuser(user Userdata) error {
	db, err := openConnection()
	if err != nil {
		return err
	}

	defer db.Close()

	userId := exists(user.Name)

	if userId == -1 {
		return fmt.Errorf("user deos not exist")
	}
	user.ID = userId
	updateStatement := `UPDATE userdata SET name=$1, surname=$2, description=$3 WHERE userid=$4`

	_, err = db.Exec(updateStatement, user.Name, user.Surname, user.Description, user.ID)
	if err != nil {
		return err
	}
	return nil
}
