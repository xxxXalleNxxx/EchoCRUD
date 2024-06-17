package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "admin"
	dbname   = "CRUDDB"
)

func init() {
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		fmt.Sprintf("failed to connect to database")
	}

	db.AutoMigrate(&User{})
}

type User struct {
	gorm.Model
	Name  string `json:"name"`
	Email string `json:"email"`
}

func getUsers(c echo.Context) error {
	var users []User
	db.Find(&users)
	return c.JSON(200, users)
}

func createUser(c echo.Context) error {
	user := new(User)
	if err := c.Bind(user); err != nil {
		return err
	}
	db.Create(&user)
	return c.JSON(201, user)
}

func getUserById(c echo.Context) error {

	id := c.Param("id")
	var user User
	db.First(&user, id)
	return c.JSON(200, user)
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/users", getUsers)
	e.GET("/users/:id", getUserById)
	e.POST("/users", createUser)

	e.Logger.Fatal(e.Start(":8080"))
}
