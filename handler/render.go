package handler

import (
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/labstack/echo/v4"
)

const tmplPath = "src/template/"

func htmlBlob(file string, data map[string]interface{}) ([]byte, error) {
	return pongo2.Must(pongo2.FromCache(tmplPath + file)).ExecuteBytes(data)
}

func render(c echo.Context, file string, data map[string]interface{}) error {
	data["CSRF"] = c.Get("csrf").(string)

	b, error := htmlBlob(file, data)
	if error != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.HTMLBlob(http.StatusOK, b)
}
