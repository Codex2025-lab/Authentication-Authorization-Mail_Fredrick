package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"time"

	"auth-mail/internal/data"
	"auth-mail/internal/mailer"

	_ "github.com/lib/pq"
)


// application holds app-wide dependencies
type application struct {
	models struct {
		Users *data.UserModel
        Tokens *data.TokenModel
	}

	mailer mailer.Mailer
}

type config struct {
    smtp struct {
        host     string
        port     int
        username string
        password string
        sender   string
    }
}




func main() {
    var cfg config

    // SMTP flags
    flag.StringVar(&cfg.smtp.host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
    flag.IntVar(&cfg.smtp.port, "smtp-port", 2525, "SMTP port")
    flag.StringVar(&cfg.smtp.username, "smtp-username", "", "SMTP username")
    flag.StringVar(&cfg.smtp.password, "smtp-password", "", "SMTP password")
    flag.StringVar(&cfg.smtp.sender, "smtp-sender", "MyApp <no-reply@myapp.com>", "SMTP sender")

    flag.Parse()

    // Open database
    db, err := openDB("postgres://fred:codex@localhost/authmail?sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Initialize application
    app := &application{}
    app.models.Users = &data.UserModel{DB: db}
    app.models.Tokens = &data.TokenModel{DB: db}

    // Initialize mailer 
    app.mailer = mailer.New(
        cfg.smtp.host,
        cfg.smtp.port,
        cfg.smtp.username,
        cfg.smtp.password,
        cfg.smtp.sender,
    )

    // Start server
    srv := &http.Server{
        Addr:         ":4000",
        Handler:      app.routes(),
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
    }

    log.Println("Starting server on :4000")
    log.Fatal(srv.ListenAndServe())
}
// openDB opens a database connection and tests it
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}