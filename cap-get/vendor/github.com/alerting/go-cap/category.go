package cap

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"strings"
)

const (
	CategoryUnknown = iota
	CategoryGeological
	CategoryMeteorological
	CategorySafety
	CategorySecurity
	CategoryRescue
	CategoryFire
	CategoryHealth
	CategoryEnvironment
	CategoryTransport
	CategoryInfrastructure
	CategoryCBRNE
	CategoryOther
)

type Category int

// UnmarshalString unmarshals the string into a Category value
func (category *Category) UnmarshalString(value string) error {
	value = strings.ToLower(value)

	if value == "geo" || value == "geological" {
		*category = CategoryGeological
	} else if value == "met" || value == "meteorological" {
		*category = CategoryMeteorological
	} else if value == "safety" {
		*category = CategorySafety
	} else if value == "security" {
		*category = CategorySecurity
	} else if value == "rescue" {
		*category = CategoryRescue
	} else if value == "fire" {
		*category = CategoryFire
	} else if value == "health" {
		*category = CategoryHealth
	} else if value == "env" || value == "environment" {
		*category = CategoryEnvironment
	} else if value == "transport" {
		*category = CategoryTransport
	} else if value == "infra" || value == "infrastructure" {
		*category = CategoryInfrastructure
	} else if value == "cbrne" {
		*category = CategoryCBRNE
	} else if value == "other" {
		*category = CategoryOther
	} else {
		return errors.New("Unknown Category value: " + value)
	}

	return nil
}

// UnmarshalXML unmarshals the XML into a Category value.
func (category *Category) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var str string
	if err := d.DecodeElement(&str, &start); err != nil {
		return err
	}
	return category.UnmarshalString(str)
}

// UnmarshalJSON unmarshals the JSON into a Category value.
func (category *Category) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	return category.UnmarshalString(str)
}

// MarshalJSON returns the string version of the category.
func (category *Category) MarshalJSON() ([]byte, error) {
	return json.Marshal(category.String())
}

// MarshalXML returns the string version of the status.
func (category *Category) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	str := category.String()

	// Some of the values are shortened in the xml version
	if str == "Geological" {
		str = "Geo"
	} else if str == "Meteorological" {
		str = "Met"
	} else if str == "Environment" {
		str = "Env"
	} else if str == "Infrastructure" {
		str = "Infra"
	}

	return e.EncodeElement(str, start)
}

// String returns a Category as a string
func (category Category) String() string {
	if category == CategoryGeological {
		return "Geological"
	} else if category == CategoryMeteorological {
		return "Meteorological"
	} else if category == CategorySafety {
		return "Safety"
	} else if category == CategorySecurity {
		return "Security"
	} else if category == CategoryRescue {
		return "Rescue"
	} else if category == CategoryFire {
		return "Fire"
	} else if category == CategoryHealth {
		return "Health"
	} else if category == CategoryEnvironment {
		return "Environment"
	} else if category == CategoryTransport {
		return "Transport"
	} else if category == CategoryInfrastructure {
		return "Infrastructure"
	} else if category == CategoryCBRNE {
		return "CBRNE"
	} else if category == CategoryOther {
		return "Other"
	} else if category == CategoryUnknown {
		return "Unknown"
	}

	return ""
}
