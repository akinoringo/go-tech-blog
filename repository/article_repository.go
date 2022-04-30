package repository

import (
	"database/sql"
	"math"

	"time"
	"go-tech-blog/model"
)

/*
* 記事の作成
*/
func ArticleCreate(article *model.Article) (sql.Result, error) {
	now := time.Now()

	article.Created = now
	article.Updated = now

	tx := db.MustBegin()

	query := `INSERT INTO articles (title, body, created, updated)
	VALUES (:title, :body, :created, :updated);`

	res, err := tx.NamedExec(query, article)
	if err != nil {
		tx.Rollback()

		return nil, err
	}

	tx.Commit()

	return res, nil
}

/*
* 記事一覧の取得
*/
func ArticleListByCursor(cursor int) ([]*model.Article, error) {
	if cursor <= 0 {
		cursor = math.MaxInt32
	}

	query := `SELECT *
	FROM articles
	WHERE id < ?
	ORDER BY id desc
	LIMIT 10`

	articles := make([]*model.Article, 0, 10)

	if err := db.Select(&articles, query, cursor); err != nil {
		return nil, err
	}

	return articles, nil
}

/*
* 記事の削除
*/
func ArticleDelete(id int) error {
	query := `DELETE FROM articles WHERE id = ?`

	tx := db.MustBegin()

	if _, err := tx.Exec(query, id); err != nil {
		tx.Rollback()

		return err
	}

	return tx.Commit()
}

func ArticleGetById(id int) (*model.Article, error) {
	query := `SELECT *
	FROM articles
	WHERE id = ?;`

	var article model.Article

	if err := db.Get(&article, query, id); err != nil {
		return nil, err
	}

	return &article, nil
}