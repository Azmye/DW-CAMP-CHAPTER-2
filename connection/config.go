package connection

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

var Conn *pgx.Conn

func DatabaseConnect() {
	DB_URI := `postgres://postgres:1234@localhost:5432/dw-personal-web`

	var err error
	Conn, err = pgx.Connect(context.Background(), DB_URI)

	// Check if there is error connection to db
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't connect to database : %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Success Connect to database")
}
