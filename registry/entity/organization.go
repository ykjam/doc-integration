package entity

import (
	"time"
)

type EntityState string
type DMSType string // Document Management System type

const (
	EntityStateDeleted  EntityState = "DELETED"
	EntityStateDisabled EntityState = "DISABLED"
	EntityStateEnabled  EntityState = "ENABLED"

	SRD        DMSType = "SRD"
	Netije     DMSType = "Netije"
	EResminama DMSType = "eResminama"
)

type Organization struct {
	Id        int
	Name      string
	Label     string
	Url       string
	Type      DMSType
	PublicKey string
	State     EntityState
	CreateTs  time.Time
	UpdateTs  time.Time
	Version   int
}

type OrganizationResponse struct {
	Id        int     `json:"id"`
	Name      string  `json:"name"`
	Label     string  `json:"label"`
	Type      DMSType `json:"type"`
	Url       string  `json:"url"`
	PublicKey string  `json:"public_key"`
	CreateTs  int64   `json:"create_ts" convert_by:"time_to_int64"`
	UpdateTs  int64   `json:"update_ts" convert_by:"time_to_int64"`
}

type OrganizationListResponse struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Label     string `json:"label"`
	Type      DMSType `json:"type"`
	Url       string `json:"url"`
	PublicKey string `json:"public_key"`
}
