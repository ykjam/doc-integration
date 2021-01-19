package datastore

import (
	"context"

	"github.com/jackc/pgx/v4"

	"ykjam/doc-registry-go/entity"
)

type Access interface {
	OrganizationAdd(ctx context.Context, pTx pgx.Tx, name, label, url, publicKey string) (item *entity.Organization, err error)
	OrganizationUpdate(ctx context.Context, pTx pgx.Tx, item *entity.Organization, name, label, url, publicKey string) (err error)
	OrganizationChangeState(ctx context.Context, pTx pgx.Tx, item *entity.Organization, state entity.EntityState) (err error)
	OrganizationById(ctx context.Context, id int) (item *entity.Organization, err error)
	OrganizationList(ctx context.Context) (items []*entity.Organization, err error)
}
