// Package access_tester provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.0 DO NOT EDIT.
package common

// Result Result of checking domain
type Result struct {
	// CanBeBlocked Can this resource be actually blocked based on status code and trigger phrases
	CanBeBlocked bool `json:"canBeBlocked"`

	// Host Value of the HTTP Host header used to connect to resource
	Host string `json:"host"`

	// Status HTTP status code
	Status int `json:"status"`

	// TriggerPhrasesFound Array of all found trigger phrases
	TriggerPhrasesFound []string `json:"triggerPhrasesFound"`
}

// CheckStatusParams defines parameters for CheckStatus.
type CheckStatusParams struct {
	Host string `form:"host" json:"host"`
}
