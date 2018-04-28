package function

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli"

	"github.com/alerting/go-cap"
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
	finder := database.NewInfoFinder()

	if _, ok := query["status"]; ok {
		var status cap.Status
		status.UnmarshalString(query["status"][0])
		finder = finder.Status(status)
	}

	if _, ok := query["message_type"]; ok {
		var messageType cap.MessageType
		messageType.UnmarshalString(query["message_type"][0])
		finder = finder.MessageType(messageType)
	}

	if _, ok := query["scope"]; ok {
		var scope cap.Scope
		scope.UnmarshalString(query["scope"][0])
		finder = finder.Scope(scope)
	}

	if _, ok := query["language"]; ok {
		finder = finder.Language(query["language"][0])
	}

	if _, ok := query["certainty"]; ok {
		var certainty cap.Certainty
		certainty.UnmarshalString(query["certainty"][0])
		finder = finder.Certainty(certainty)
	}

	if _, ok := query["urgency"]; ok {
		var urgency cap.Urgency
		urgency.UnmarshalString(query["urgency"][0])
		finder = finder.Urgency(urgency)
	}

	if _, ok := query["severity"]; ok {
		var severity cap.Severity
		severity.UnmarshalString(query["severity"][0])
		finder = finder.Severity(severity)
	}

	if _, ok := query["headline"]; ok {
		finder = finder.Headline(query["headline"][0])
	}

	if _, ok := query["description"]; ok {
		finder = finder.Description(query["description"][0])
	}

	if _, ok := query["instruction"]; ok {
		finder = finder.Instruction(query["instruction"][0])
	}

	if val, ok := query["effective_gte"]; ok {
		t, err := time.Parse(time.RFC3339, val[0])
		if err != nil {
			log.Fatal(err)
		}

		finder = finder.EffectiveGte(t)
	}

	if val, ok := query["effective_gt"]; ok {
		t, err := time.Parse(time.RFC3339, val[0])
		if err != nil {
			log.Fatal(err)
		}

		finder = finder.EffectiveGt(t)
	}

	if val, ok := query["effective_lte"]; ok {
		t, err := time.Parse(time.RFC3339, val[0])
		if err != nil {
			log.Fatal(err)
		}

		finder = finder.EffectiveLte(t)
	}

	if val, ok := query["effective_lt"]; ok {
		t, err := time.Parse(time.RFC3339, val[0])
		if err != nil {
			log.Fatal(err)
		}

		finder = finder.EffectiveLt(t)
	}

	if val, ok := query["expires_gte"]; ok {
		t, err := time.Parse(time.RFC3339, val[0])
		if err != nil {
			log.Fatal(err)
		}

		finder = finder.ExpiresGte(t)
	}

	if val, ok := query["expires_gt"]; ok {
		t, err := time.Parse(time.RFC3339, val[0])
		if err != nil {
			log.Fatal(err)
		}

		finder = finder.ExpiresGt(t)
	}

	if val, ok := query["expires_lte"]; ok {
		t, err := time.Parse(time.RFC3339, val[0])
		if err != nil {
			log.Fatal(err)
		}

		finder = finder.ExpiresLte(t)
	}

	if val, ok := query["expires_lt"]; ok {
		t, err := time.Parse(time.RFC3339, val[0])
		if err != nil {
			log.Fatal(err)
		}

		finder = finder.ExpiresLt(t)
	}

	if val, ok := query["onset_gte"]; ok {
		t, err := time.Parse(time.RFC3339, val[0])
		if err != nil {
			log.Fatal(err)
		}

		finder = finder.OnsetGte(t)
	}

	if val, ok := query["onset_gt"]; ok {
		t, err := time.Parse(time.RFC3339, val[0])
		if err != nil {
			log.Fatal(err)
		}

		finder = finder.OnsetGt(t)
	}

	if val, ok := query["onset_lte"]; ok {
		t, err := time.Parse(time.RFC3339, val[0])
		if err != nil {
			log.Fatal(err)
		}

		finder = finder.OnsetLte(t)
	}

	if val, ok := query["onset_lt"]; ok {
		t, err := time.Parse(time.RFC3339, val[0])
		if err != nil {
			log.Fatal(err)
		}

		finder = finder.OnsetLt(t)
	}

	if _, ok := query["area"]; ok {
		finder = finder.Area(query["area"][0])
	}

	if _, ok := query["point"]; ok {
		str := strings.Split(query["point"][0], ",")

		lat, err := strconv.ParseFloat(str[0], 64)
		if err != nil {
			log.Fatal(err)
		}

		lon, err := strconv.ParseFloat(str[1], 64)
		if err != nil {
			log.Fatal(err)
		}

		finder = finder.Point(lat, lon)
	}

	if _, ok := query["from"]; ok {
		from, err := strconv.Atoi(query["from"][0])
		if err == nil {
			finder = finder.Start(from)
		}
	}

	if _, ok := query["size"]; ok {
		size, err := strconv.Atoi(query["size"][0])
		if err == nil {
			finder = finder.Count(size)
		}
	}

	if _, ok := query["sort"]; ok {
		finder = finder.Sort(query["sort"][0])
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
