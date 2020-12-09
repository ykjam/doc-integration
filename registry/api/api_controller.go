package api

import (
	"ykjam/doc-registry-go/datastore"
)

type APIController struct {
	access datastore.Access
}

func NewAPIController(access datastore.Access) *APIController {
	return &APIController{
		access: access,
	}
}
