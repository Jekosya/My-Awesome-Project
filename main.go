package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Age      int       `json:"age"`
	Comments []Comment `json:"comments"`
}

type Comment struct {
	ID     int    `json:"id"`
	Text   string `json:"text"`
	UserID int    `json:"user_id"`
}

// CreateUserTable создает таблицу пользователей и комментариев в базе данных.
func CreateUserTable() error {
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		return err
	}
	defer db.Close()

	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		age INTEGER NOT NULL
	);`
	if _, err := db.Exec(query); err != nil {
		return err
	}

	query = `
	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		text TEXT NOT NULL,
		user_id INTEGER,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`
	if _, err := db.Exec(query); err != nil {
		return err
	}

	return nil
}

// InsertUser вставляет нового пользователя в таблицу.
func InsertUser(user User) error {
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		return err
	}
	defer db.Close()

	query, args, err := PrepareQuery("insert", "users", user)
	if err != nil {
		return err
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

// SelectUser выбирает пользователя по его идентификатору.
func SelectUser(userID int) (User, error) {
	var user User

	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		return user, err
	}
	defer db.Close()

	query, _, err := PrepareQuery("select", "users", user)
	row := db.QueryRow(query, userID)
	err = row.Scan(&user.ID, &user.Name, &user.Age)
	if err != nil {
		return user, err
	}

	return user, nil
}

// UpdateUser обновляет информацию о пользователе.
func UpdateUser(user User) error {
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		return err
	}
	defer db.Close()

	query, _, err := PrepareQuery("update", "users", user)
	_, err = db.Exec(query, user.Name, user.Age, user.ID)
	return err
}

// DeleteUser удаляет пользователя из таблицы.
func DeleteUser(userID int) error {
	var user User
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		return err
	}
	defer db.Close()

	query, _, err := PrepareQuery("delete", "users", user)
	_, err = db.Exec(query, userID)
	return err
}

func PrepareQuery(operation string, table string, user User) (string, []interface{}, error) {
	var query string
	var args []interface{}
	switch operation {
	case "insert":
		query = "INSERT INTO users(name, age) VALUES (?, ?)"
		args = append(args, user.Name, user.Age)
	case "select":
		query = "SELECT id, name, age FROM users"
	case "update":
		query = "UPDATE users SET name = ?, age = ? WHERE id = ?"
		args = append(args, user.Name, user.Age, user.Comments)
	case "delete":
		query = "DELETE FROM users WHERE id = ?"
		args = append(args, user.ID)
	}
	return query, args, nil
}

func main() {
	err := CreateUserTable()
	if err != nil {
		fmt.Println("Ошибка создания таблицы пользователей и комментариев:", err)
		return
	}

	// Примеры использования функций работы с пользователями
	user := User{
		Name: "John Doe",
		Age:  30,
		Comments: []Comment{
			{Text: "First comment"},
			{Text: "Second comment"},
		},
	}

	err = InsertUser(user)
	if err != nil {
		fmt.Println("Ошибка добавления пользователя:", err)
		return
	}
	fmt.Println("Пользователь успешно добавлен.")

	fetchedUser, err := SelectUser(1)
	if err != nil {
		fmt.Println("Ошибка выборки пользователя:", err)
		return
	}
	fmt.Printf("Выбранный пользователь: %+v\n", fetchedUser)

	fetchedUser.Age = 34
	err = UpdateUser(fetchedUser)
	if err != nil {
		fmt.Println("Ошибка обновления пользователя:", err)
		return
	}
	fmt.Println("Пользователь успешно обновлен.")

	err = DeleteUser(fetchedUser.ID)
	if err != nil {
		fmt.Println("Ошибка удаления пользователя:", err)
		return
	}
	fmt.Println("Пользователь успешно удален.")
}
