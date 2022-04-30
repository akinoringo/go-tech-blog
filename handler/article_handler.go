package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go-tech-blog/model"
	"go-tech-blog/repository"

	"github.com/labstack/echo/v4"
)

type ArticleCreateOutput struct {
	Article *model.Article
	Message string
	ValidationErrors []string
}

func ArticleCreate(c echo.Context) error {
	var article model.Article

	var out ArticleCreateOutput

	if err := c.Bind(&article); err != nil {
		c.Logger().Error(err.Error())

		return c.JSON(http.StatusBadRequest, out)
	}

	if err := c.Validate(&article); err != nil {
		c.Logger().Error(err.Error())

		out.ValidationErrors = article.ValidationErrors(err)

		return c.JSON(http.StatusUnprocessableEntity, out)
	}

	res, err := repository.ArticleCreate(&article)
	if err != nil {
		c.Logger().Error(err.Error())

		return c.JSON(http.StatusInternalServerError, out)
	}

	id, _ := res.LastInsertId()

	article.ID = int(id)

	out.Article = &article

	return c.JSON(http.StatusOK, out)

}

func ArticleIndex(c echo.Context) error {
	if c.Request().URL.Path == "/articles" {
		c.Redirect(http.StatusPermanentRedirect, "/")
	}

	articles, err := repository.ArticleListByCursor(0)

	if err != nil {
		c.Logger().Error(err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	// 取得できた最後の記事のIDをカーソルとして設定
	var cursor int
	if len(articles) != 0 {
		cursor = articles[len(articles)-1].ID
	}

	data := map[string]interface{}{
		"Articles": articles,
		"Cursor": cursor,
	}

	return render(c, "article/index.html", data)
}

func ArticleNew(c echo.Context) error {
	data := map[string]interface{}{
		"Message": "Article New",
		"Now":     time.Now(),
	}
	return render(c, "article/new.html", data)
}

/*
* ArticleShow
*/
func ArticleShow(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("articleId"))

	article, err := repository.ArticleGetById(id)

	if err != nil {
		c.Logger().Error(err.Error())

		return c.NoContent(http.StatusInternalServerError)
	}

	data := map[string]interface{}{
		"Article":article,
	}

	return render(c, "article/show.html", data)
}

func ArticleEdit(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("articleId"))

	data := map[string]interface{}{
		"Message": "Article Edit",
		"Now":     time.Now(),
		"ID":      id,
	}

	return render(c, "article/edit.html", data)
}

func ArticleDelete(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("articleId"))

	if err := repository.ArticleDelete(id); err != nil {
		c.Logger().Error(err.Error())

		return c.JSON(http.StatusInternalServerError, "")
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("Article %d is deleted", id))
}

func ArticleList(c echo.Context) error {
	cursor, _ := strconv.Atoi(c.QueryParam("cursor"))

	articles, err := repository.ArticleListByCursor(cursor)

	if err != nil {

		c.Logger().Error(err.Error())

		return c.JSON(http.StatusInternalServerError, "")
	}

	return c.JSON(http.StatusOK, articles)
}