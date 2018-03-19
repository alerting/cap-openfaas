package elastic

import (
	"context"
	"encoding/json"
	"strings"

	cap "github.com/alerting/go-cap"
	"github.com/alerting/go-cap-process/db"
	"github.com/olivere/elastic"
)

// AlertsFinder defines the search criteria for finding a alerts.
type AlertsFinder struct {
	elastic *Elastic

	// Pagination
	from int
	size int

	// Sort
	sort []string

	// alert
	alertFields map[string]string

	// alert.info
	alertInfoFields     map[string]string
	alertInfoTextFields map[string]string

	// alert.info.area
	coordinates *elastic.GeoPoint

	// alert.info.resource
}

// NewAlertsFinder creates a new alerts finder.
func NewAlertsFinder(elastic *Elastic) db.AlertsFinder {
	return &AlertsFinder{
		elastic:             elastic,
		alertFields:         make(map[string]string),
		alertInfoFields:     make(map[string]string),
		alertInfoTextFields: make(map[string]string),
	}
}

// Status filters the resutls by the given status.
func (f *AlertsFinder) Status(status string) db.AlertsFinder {
	f.alertFields["status"] = status
	return f
}

// MessageType filters the results by the given message type.
func (f *AlertsFinder) MessageType(messageType string) db.AlertsFinder {
	f.alertFields["message_type"] = messageType
	return f
}

// Scope filters the results by the given scope.
func (f *AlertsFinder) Scope(scope string) db.AlertsFinder {
	f.alertFields["scope"] = scope
	return f
}

// Language filters the results by the given language.
func (f *AlertsFinder) Language(language string) db.AlertsFinder {
	f.alertInfoFields["infos.language"] = language
	return f
}

// Certainty filters the results by the given certainty.
func (f *AlertsFinder) Certainty(certainty string) db.AlertsFinder {
	f.alertInfoFields["infos.certainty"] = certainty
	return f
}

// Severity filters the results by the given severity.
func (f *AlertsFinder) Severity(severity string) db.AlertsFinder {
	f.alertInfoFields["infos.severity"] = severity
	return f
}

// Urgency filters the results by the given urgency.
func (f *AlertsFinder) Urgency(urgency string) db.AlertsFinder {
	f.alertInfoFields["infos.urgency"] = urgency
	return f
}

// Headline filters the results by the given headline.
func (f *AlertsFinder) Headline(headline string) db.AlertsFinder {
	f.alertInfoTextFields["infos.headline"] = headline
	return f
}

// Description filters the results by the given description.
func (f *AlertsFinder) Description(description string) db.AlertsFinder {
	f.alertInfoTextFields["infos.description"] = description
	return f
}

// Point filters alerts containing the given coordinate.
func (f *AlertsFinder) Point(latLon string) db.AlertsFinder {
	point, err := elastic.GeoPointFromString(latLon)
	if err != nil {
		panic(err)
	}

	f.coordinates = point
	return f
}

// From specifies the start index for pagination.
func (f *AlertsFinder) From(from int) db.AlertsFinder {
	f.from = from
	return f
}

// Size specifies the number of items to return in pagination.
func (f *AlertsFinder) Size(size int) db.AlertsFinder {
	f.size = size
	return f
}

// Sort specifies one or more sort orders.
// Use a dash (-) to make the sort order descending.
// Example: "name" or "-year".
func (f *AlertsFinder) Sort(fields ...string) db.AlertsFinder {
	if f.sort == nil {
		f.sort = make([]string, 0)
	}
	f.sort = append(f.sort, fields...)
	return f
}

// Find executes the search and returns the results.
func (f *AlertsFinder) Find() (*db.AlertsResults, error) {
	search := f.elastic.client.Search().Index(f.elastic.index).Type("alert")
	search = f.query(search)
	search = f.paginate(search)
	search = f.sorting(search)

	rawResults, err := search.Do(context.Background())
	if err != nil {
		return nil, err
	}

	// Generate the results
	results := db.AlertsResults{
		Total: rawResults.Hits.TotalHits,
		Hits:  make([]db.AlertHit, 0),
	}

	for _, alertHit := range rawResults.Hits.Hits {
		var alert cap.Alert
		if err = json.Unmarshal(*alertHit.Source, &alert); err != nil {
			return nil, err
		}

		hit := db.AlertHit{
			Alert:    &alert,
			InfoHits: make(map[int]db.InfoHit),
		}

		if _, ok := alertHit.InnerHits["infos"]; ok {
			for _, infoHit := range alertHit.InnerHits["infos"].Hits.Hits {
				ihit := db.InfoHit{
					AreaHits: make([]int, 0),
				}

				if _, ok = infoHit.InnerHits["infos.areas"]; ok {
					for _, areaHit := range infoHit.InnerHits["infos.areas"].Hits.Hits {
						ihit.AreaHits = append(ihit.AreaHits, areaHit.Nested.Offset)
					}
				}

				hit.InfoHits[infoHit.Nested.Offset] = ihit
			}
		}

		results.Hits = append(results.Hits, hit)
	}

	return &results, nil
}

func (f *AlertsFinder) query(service *elastic.SearchService) *elastic.SearchService {
	if len(f.alertFields) == 0 && len(f.alertInfoFields) == 0 && len(f.alertInfoTextFields) == 0 {
		service = service.Query(elastic.NewMatchAllQuery())
		return service
	}

	service = service.FetchSourceContext(
		elastic.NewFetchSourceContext(true).Exclude("infos.resources.deref_uri"))

	// alert
	q := elastic.NewBoolQuery()

	for k, v := range f.alertFields {
		q = q.Must(elastic.NewTermQuery(k, v))
	}

	// alert.info
	if len(f.alertInfoFields) > 0 || len(f.alertInfoTextFields) > 0 || f.coordinates != nil {
		infoBoolQ := elastic.NewBoolQuery()

		for k, v := range f.alertInfoFields {
			infoBoolQ = infoBoolQ.Must(elastic.NewTermQuery(k, v))
		}
		for k, v := range f.alertInfoTextFields {
			infoBoolQ = infoBoolQ.Must(elastic.NewQueryStringQuery(v).Field(k))
		}

		// alert.info.area
		if f.coordinates != nil {
			areaQ := elastic.NewNestedQuery("infos.areas",
				NewGeoShapeQuery("infos.areas.polygons").SetPoint(f.coordinates.Lat, f.coordinates.Lon))
			areaQ = areaQ.InnerHit(
				elastic.NewInnerHit().FetchSourceContext(
					elastic.NewFetchSourceContext(false)))

			infoBoolQ = infoBoolQ.Must(areaQ)
		}

		infoQ := elastic.NewNestedQuery("infos", infoBoolQ)
		infoHit := elastic.NewInnerHit().FetchSourceContext(elastic.NewFetchSourceContext(false))

		infoQ = infoQ.InnerHit(infoHit)
		q = q.Must(infoQ)
	}

	// alert.info.resource

	service = service.Query(q)
	return service
}

func (f *AlertsFinder) paginate(service *elastic.SearchService) *elastic.SearchService {
	if f.from > 0 {
		service = service.From(f.from)
	}

	if f.size > 0 {
		service = service.Size(f.size)
	}

	return service
}

// sorting applies sorting to the service.
// TODO: Handle nested sorts
func (f *AlertsFinder) sorting(service *elastic.SearchService) *elastic.SearchService {
	if len(f.sort) == 0 {
		// Sort by score by default
		service = service.Sort("_score", false)
		return service
	}

	// Sort by fields; prefix of "-" means: descending sort order.
	for _, s := range f.sort {
		s = strings.TrimSpace(s)

		var field string
		var asc bool

		if strings.HasPrefix(s, "-") {
			field = s[1:]
			asc = false
		} else {
			field = s
			asc = true
		}

		// Maybe check for permitted fields to sort

		service = service.Sort(field, asc)
	}
	return service
}
