package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	tele "gopkg.in/telebot.v3"
)

func main() {
	// load .env file
	envFile := godotenv.Load(".env")
	if envFile != nil {
		log.Fatal(envFile)
		return
	}

	// connect database
	db, err := sql.Open("sqlite3", "./database/babymonitor.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	pref := tele.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	var (
		menu = &tele.ReplyMarkup{ResizeKeyboard: true}

		// reply buttons
		btnKick  = menu.Text("Baby is moving")
		btnCount = menu.Text("Movement total")
	)

	menu.Reply(
		menu.Row(btnKick),
		menu.Row(btnCount),
	)

	b.Handle("/start", func(c tele.Context) error {
		return c.Send("Hai mama", menu)
	})

	b.Handle(&btnKick, func(c tele.Context) error {
		c.Delete()

		sts := "INSERT INTO logs(created_at) VALUES(datetime('now', 'localtime'));"
		_, err := db.Exec(sts)
		if err != nil {
			log.Fatal(err)
		}

		return c.Send("Baby is moving")
	})

	b.Handle(&btnCount, func(c tele.Context) error {
		c.Delete()

		sts := "SELECT COUNT(id) FROM logs WHERE date('now');"
		var total string

		row := db.QueryRow(sts)
		err = row.Scan(&total)
		if err != nil {
			log.Fatal(err)
		}

		return c.Send("Today movement: " + total)
	})

	b.Start()
}
