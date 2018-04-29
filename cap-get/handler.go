package function

import (
	"encoding/json"
	"log"
	"net/url"
	"os"

	"github.com/urfave/cli"

	"github.com/alerting/go-cap-process/db"
	"github.com/alerting/go-cap-process/tasks"
)

// Handle a serverless request
func Handle(req []byte) string {
	var database db.Database

	// Piggy-back off the command line parsing
	// to get the database object.
	app := cli.NewApp()
	app.Flags = tasks.DatabaseFlags

	app.Action = func(c *cli.Context) error {
		// Connect to the database, if we haven't already
		var err error
		database, err = tasks.CreateDatabase(c)
		if err != nil {
			return err
		}
		return nil
	}

	// Let's get the database, if we need it
	if err := app.Run([]string{""}); err != nil {
		log.Fatal(err)
	}

	// Parse query string
	query, err := url.ParseQuery(os.Getenv("Http_Query"))
	if err != nil {
		log.Fatal(err)
	}

	if val, ok := query["id"]; ok {
		alert, err := database.GetAlertById(val[0])
		if err != nil {
			log.Fatal(err)
		}

		b, err := json.Marshal(alert)
		if err != nil {
			log.Fatal(err)
		}

		return string(b)
	}

	log.Fatal("No id provided")
	return ""
}
