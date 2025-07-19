package models

// Emergency represents a single emergency contact
// Used as an element in the User.EmergencyContacts array
// and in request/response payloads
type Emergency struct {
	Name         string `json:"name"`  // full name
	Email        string `json:"email"` // email address
	CountryCode  string `json:"countryCode,omitempty" dynamodbav:"countryCode,omitempty"`
	Phone        string `json:"phone,omitempty" dynamodbav:"phone,omitempty"`
	Relationship string `json:"relationship"` // includes country code
}
