package main

import (
	"log"
	"os"

	"go-tech-blog/handler"
	"go-tech-blog/repository"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/go-playground/validator.v9"
)

var db *sqlx.DB
var e = createMux()

func main() {
	db = connectDB()
	repository.SetDB(db)

	e.GET("/", handler.ArticleIndex)

	e.GET("/articles", handler.ArticleIndex)
	e.GET("/articles/new", handler.ArticleNew)
	e.GET("/articles/:articleId", handler.ArticleShow)
	e.GET("/articles/:articleId/edit", handler.ArticleEdit)
	e.GET("/api/articles", handler.ArticleList)
	e.POST("/api/articles", handler.ArticleCreate)
	e.DELETE("/api/:articleId", handler.ArticleDelete)

	e.Logger.Fatal(e.Start(":8080"))
}

func connectDB() *sqlx.DB {
	dsn := os.Getenv("DSN")
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		e.Logger.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		e.Logger.Fatal(err)
	}
	log.Println("db connection succeeded")
	return db
}

func createMux() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.Gzip())
	e.Use(middleware.CSRF())

	e.Static("/css", "src/css")
	e.Static("/js", "src/js")

	e.Validator = &CustomValidator{validator: validator.New()}

	return e
}

// CustomValidator
type CustomValidator struct {
	validator *validator.Validate
}

// Validate
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}