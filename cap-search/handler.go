package function

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"strconv"

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

	// Setup the finder
	finder := database.NewAlertsFinder()

	if _, ok := query["status"]; ok {
		finder = finder.Status(query["status"][0])
	}

	if _, ok := query["message_type"]; ok {
		finder = finder.MessageType(query["message_type"][0])
	}

	if _, ok := query["scope"]; ok {
		finder = finder.Scope(query["scope"][0])
	}

	if _, ok := query["language"]; ok {
		finder = finder.Language(query["language"][0])
	}

	if _, ok := query["certainty"]; ok {
		finder = finder.Certainty(query["certainty"][0])
	}

	if _, ok := query["urgency"]; ok {
		finder = finder.Urgency(query["urgency"][0])
	}

	if _, ok := query["severity"]; ok {
		finder = finder.Severity(query["severity"][0])
	}

	if _, ok := query["headline"]; ok {
		finder = finder.Headline(query["headline"][0])
	}

	if _, ok := query["description"]; ok {
		finder = finder.Description(query["description"][0])
	}

	if _, ok := query["point"]; ok {
		finder = finder.Point(query["point"][0])
	}

	if _, ok := query["from"]; ok {
		from, err := strconv.Atoi(query["from"][0])
		if err == nil {
			finder = finder.From(from)
		}
	}

	if _, ok := query["sort"]; ok {
		finder = finder.Sort(query["sort"][0])
	}

	if _, ok := query["size"]; ok {
		size, err := strconv.Atoi(query["size"][0])
		if err == nil {
			finder = finder.Size(size)
		}
	}

	res, err := finder.Find()
	if err != nil {
		log.Fatal(err)
	}

	b, err := json.Marshal(&res)
	if err != nil {
		log.Fatal(err)
	}

	return string(b)
}
