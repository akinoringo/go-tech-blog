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

var authUser = os.Getenv("AUTH_USER")
var authPassword = os.Getenv("AUTH_PASSWORD")
var db *sqlx.DB
var e = createMux()

func main() {
	db = connectDB()
	repository.SetDB(db)

	auth := e.Group("")
	auth.Use(basicAuth())

	e.GET("/", handler.ArticleIndex)

	e.GET("/articles", handler.ArticleIndex)
	auth.GET("/articles/new", handler.ArticleNew)
	e.GET("/articles/:articleId", handler.ArticleShow)
	auth.GET("/articles/:articleId/edit", handler.ArticleEdit)
	e.GET("/api/articles", handler.ArticleList)
	auth.POST("/api/articles", handler.ArticleCreate)
	auth.DELETE("/api/articles/:articleId", handler.ArticleDelete)
	auth.PATCH("/api/articles/:articleId", handler.ArticleUpdate)

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
	e.Use(basicAuth())

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

// Basic認証
func basicAuth() echo.MiddlewareFunc {
	var basicAuthValidator middleware.BasicAuthValidator

	// basicAuthValidator という変数に、認証成功・認証失敗を判定する関数を代入します。
	basicAuthValidator = func(username, password string, c echo.Context) (bool, error) {
		// ユーザー名が "joe"、パスワードが "secret" の場合認証に成功します。
		if username == authUser && password == authPassword {
			return true, nil
		}
		return false, nil
	}

	// middleware パッケージの BasicAuth() 関数は、
	// Basic 認証判定用の関数を引数に取り、MiddlewareFunc 型を返却します。
	middlewareFunc := middleware.BasicAuth(basicAuthValidator)
	return middlewareFunc
}
