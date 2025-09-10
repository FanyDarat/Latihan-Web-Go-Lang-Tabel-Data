package main

import (
	"database/sql"
	"fmt"
	"log"

	sqlDriver "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	driver "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID    uint   `json:"ID" gorm:"primary_key" gorm:"auto_increment"`
	Name  string `json:"Name"`
	Email string `json:"Email" gorm:"unique"`
}

var db *gorm.DB

func initDB() {
	//Connection
	koneksiRoot := "root:rafael123@tcp(localhost:3308)/"
	sqlDatabase, err := sql.Open("mysql", koneksiRoot)
	if err != nil {
		log.Fatal("Gagal melakukan koneksi ke database: ", err)
	}
	defer sqlDatabase.Close()

	//Create database
	_, err = sqlDatabase.Exec("CREATE DATABASE IF NOT EXISTS fiber_gorm")
	if err != nil {
		log.Fatal("Gagal membuat database: ", err)
	}

	//Connect to the created database
	koneksiDatabase := "root:rafael123@tcp(localhost:3308)/fiber_gorm?charset=utf8mb4&parseTime=True&loc=Local"
	database, err := gorm.Open(driver.Open(koneksiDatabase), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal connect ke database: ", err)
	}

	db = database

	db.AutoMigrate(User{})
	fmt.Println("Tabel user berhasil dibuat!")
}

// Handler

func createUser(c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	db.Create(&user)
	result := db.Create(&user)
	//Handle Duplicate
	if result.Error != nil {
		if mysqlErr, ok := result.Error.(*sqlDriver.MySQLError); ok && mysqlErr.Number == 1062 {
			return c.Status(400).JSON(fiber.Map{
				"error": "Email already exists",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error": result.Error.Error(),
		})
	}
	return c.JSON(user)
}

func getUsers(c *fiber.Ctx) error {
	var users []User
	result := db.Find(&users)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": result.Error.Error(),
		})
	}
	return c.JSON(users)
}

func main() {
	app := fiber.New()
	initDB()

	app.Post("/users", createUser)
	app.Get("/users", getUsers)

	app.Listen(":8080")
}
