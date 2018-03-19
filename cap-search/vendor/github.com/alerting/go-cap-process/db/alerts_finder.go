package db

import cap "github.com/alerting/go-cap"

type InfoHit struct {
	AreaHits []int `json:"area_hits"`
}

type AlertHit struct {
	Alert    *cap.Alert      `json:"alert"`
	InfoHits map[int]InfoHit `json:"info_hits"`
}

type AlertsResults struct {
	Total int64      `json:"total"`
	Hits  []AlertHit `json:"hits"`
}

type AlertsFinder interface {
	Status(status string) AlertsFinder
	MessageType(messageType string) AlertsFinder
	Scope(scope string) AlertsFinder

	// info
	Language(language string) AlertsFinder
	Certainty(certainty string) AlertsFinder
	Severity(severity string) AlertsFinder
	Urgency(urgency string) AlertsFinder
	Headline(headline string) AlertsFinder
	Description(description string) AlertsFinder

	// info.area
	Point(latLon string) AlertsFinder

	// Sorting
	Sort(fields ...string) AlertsFinder

	// Pagination configuration
	From(from int) AlertsFinder
	Size(size int) AlertsFinder

	// Search
	Find() (*AlertsResults, error)
}
