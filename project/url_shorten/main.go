package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-module/carbon/v2"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	nanoid "github.com/matoous/go-nanoid/v2"
)

const PG_URL = "host=postgresql-postgresql-master-1 user=kindy password=kindy dbname=shortener sslmode=disable"

var DB *sqlx.DB

func init() {
	db, err := sqlx.Connect("postgres", PG_URL)
	if err != nil {
		log.Fatalln(err)
	}
	DB = db
}

type Url struct {
	Id  string `db:"id"`
	Url string `db:"url"`
}

type ShortenReq struct {
	Uri string `json:"uri" xml:"uri" form:"uri"`
}

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString(fmt.Sprintf("I'm online %s", carbon.Now().ToDateTimeString()))
	})

	app.Get("/:id", func(c *fiber.Ctx) error {
		log.Println(c.Params("id"))
		url := Url{}
		err := DB.Get(&url, "SELECT * FROM urls WHERE id=$1", c.Params("id"))
		if err != nil {
			return err
		}
		redirect := "https://" + url.Url
		log.Println(redirect)

		return c.Redirect(redirect, fiber.StatusTemporaryRedirect)
	})

	app.Post("/shorten", func(c *fiber.Ctx) error {
		s := new(ShortenReq)
		if err := c.BodyParser(s); err != nil {
			return err
		}
		log.Println(s.Uri)

		id, err := shorten(s.Uri)
		if err != nil {
			return err
		}
		return c.SendString(id)
	})

	log.Fatal(app.Listen(":3000"))
}

func shorten(uri string) (string, error) {
	for i := 0; i < 3; i++ {
		id, err := nanoid.New(6)
		if err != nil {
			return "", err
		}
		id, err = inner_shorten(id, uri)
		if err != nil {
			var pgErr *pq.Error
			if errors.As(err, &pgErr) {
				if pgErr.Code == pq.ErrorCode("23505") {
					continue
				}
			}
			return "", err
		}
		return id, nil
	}
	return "", fmt.Errorf("can't generate shorter url now")
}

func inner_shorten(id string, uri string) (string, error) {
	sql := `INSERT INTO urls (id, url) VALUES ($1, $2) 
					ON CONFLICT(url) 
					DO UPDATE SET url=EXCLUDED.url RETURNING id`

	rows, err := DB.Query(sql, id, uri)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	if rows.Next() {
		if err = rows.Scan(&id); err != nil {
			return "", nil
		} else {
			return id, nil
		}
	}
	return "", errors.New("not found")
}
