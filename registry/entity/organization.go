package entity

import (
	"time"
)

type EntityState string

const (
	EntityStateDeleted EntityState = "DELETED"
	EntityStateEnabled EntityState = "ENABLED"
)

type Organization struct {
	Id        int
	Name      string
	Url       string
	PublicKey string
	State     EntityState
	CreateTs  time.Time
	UpdateTs  time.Time
	Version   int
}

type OrganizationResponse struct {
	Id          int          `json:"id"`
	Name        string       `json:"name"`
	Url         string       `json:"url"`
	PublicKey   string       `json:"public_key"`
	CreateTs    int64        `json:"create_ts" convert_by:"time_to_int64"`
	UpdateTs    int64        `json:"update_ts" convert_by:"time_to_int64"`
}

type OrganizationListResponse struct {
	Id          int          `json:"id"`
	Name        string       `json:"name"`
	Url         string       `json:"url"`
	PublicKey   string       `json:"public_key"`
}
