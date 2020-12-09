package api

import (
	"context"

	log "github.com/sirupsen/logrus"

	"ykjam/doc-registry-go/entity"
)

func (api *APIController) OrganizationList(ctx context.Context) (items []*entity.OrganizationListResponse, err error) {
	clog := log.WithFields(log.Fields{
		"method": "api.OrganizationList",
	})
	items = make([]*entity.OrganizationListResponse, 0)
	var organizations []*entity.Organization
	organizations, err = api.access.OrganizationList(ctx)
	if err != nil {
		eMsg := "error in access.OrganizationList"
		clog.WithError(err).Error(eMsg)
		err = ErrInternalServerError
		return
	}
	for _, organization := range organizations {
		item := &entity.OrganizationListResponse{
			Id:        organization.Id,
			Name:      organization.Name,
			Url:       organization.Url,
			PublicKey: organization.PublicKey,
		}
		items = append(items, item)
	}
	return
}
