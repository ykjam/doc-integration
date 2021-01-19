package datastore

import (
	"context"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"ykjam/doc-registry-go/entity"
)

const (
	sqlOrganizationAdd    = `INSERT INTO tbl_organization(name, label, url, public_key, state, create_ts, update_ts, version) VALUES($1, $2, $3, $4, $5, $6, $7)`
	sqlOrganizationUpdate = `UPDATE tbl_organization SET name=$3, label=$4, url=$5, public_key=$6, state=$7, update_ts=$8, version=$9 WHERE id=$1 AND version=$2`
	sqlOrganizationById   = `SELECT id, name, label, url, public_key, state, create_ts, update_ts, version FROM tbl_organization WHERE id=$1 AND state!=$2`
	sqlOrganizationByList = `SELECT id, name, label, url, public_key, state, create_ts, update_ts, version FROM tbl_organization WHERE state!=$1 ORDER BY id ASC`
)

func (d *PgAccess) organizationAddAtomic(ctx context.Context, pTx pgx.Tx, name, label, url, publicKey string, state entity.EntityState) (item *entity.Organization, err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.organizationAddAtomic",
	})
	err = d.runInTx(ctx, pTx, clog, func(tx pgx.Tx) (rollback bool, err error) {
		now := time.Now().UTC().Round(time.Microsecond)
		item = &entity.Organization{
			Name:      name,
			Label:     label,
			Url:       url,
			PublicKey: publicKey,
			State:     state,
			CreateTs:  now,
			UpdateTs:  now,
			Version:   0,
		}
		//sqlOrganizationAdd = `INSERT INTO tbl_organization(id, name, label, url, public_key, state, create_ts, update_ts, version) VALUES($1, $2, $3, $4, $5, $6, $7)`
		var cmdTag pgconn.CommandTag
		cmdTag, err = tx.Exec(ctx, sqlOrganizationAdd, item.Name, item.Label, item.Url, item.PublicKey, item.State, item.CreateTs, item.UpdateTs, item.Version)
		if err != nil {
			eMsg := "error in sqlOrganizationAdd"
			clog.WithError(err).Error(eMsg)
			rollback = true
			err = errors.Wrap(err, eMsg)
			return
		}
		if cmdTag.RowsAffected() == 0 {
			eMsg := "no rows affected during update"
			clog.Warn(eMsg)
			rollback = true
			err = errors.Wrap(ErrNoRowsAffected, eMsg)
			return
		}
		return false, nil
	})
	if err != nil {
		eMsg := "error in pgxAccess.runInTx"
		clog.WithError(err).Error(eMsg)
	}
	return
}
func (d *PgAccess) organizationUpdateAtomic(ctx context.Context, pTx pgx.Tx, item *entity.Organization, name, label, url, publicKey string, state entity.EntityState) (err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.organizationUpdateAtomic",
	})
	err = d.runInTx(ctx, pTx, clog, func(tx pgx.Tx) (rollback bool, err error) {
		now := time.Now().UTC().Round(time.Microsecond)
		nv := newVersion(item.Version)
		//sqlOrganizationUpdate = `UPDATE tbl_organization SET name=$3, url=$4, public_key=$5, state=$6, update_ts=$7, version=$8 WHERE id=$1 AND version=$2`
		var cmdTag pgconn.CommandTag
		cmdTag, err = tx.Exec(ctx, sqlOrganizationUpdate, item.Id, item.Version, name, label, url, publicKey, state, now, nv)
		if err != nil {
			eMsg := "error in sqlOrganizationUpdate"
			clog.WithError(err).Error(eMsg)
			rollback = true
			err = errors.Wrap(err, eMsg)
			return
		}
		if cmdTag.RowsAffected() == 0 {
			eMsg := "no rows affected during update"
			clog.Warn(eMsg)
			rollback = true
			err = errors.Wrap(ErrNoRowsAffected, eMsg)
			return
		}
		item.Name = name
		item.Label = label
		item.Url = url
		item.PublicKey = publicKey
		item.State = state
		item.UpdateTs = now
		item.Version = nv
		return false, nil
	})
	if err != nil {
		eMsg := "error in pgxAccess.runInTx"
		clog.WithError(err).Error(eMsg)
	}
	return
}
func (d *PgAccess) OrganizationAdd(ctx context.Context, pTx pgx.Tx, name, label, url, publicKey string) (item *entity.Organization, err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.OrganizationAdd",
	})
	err = d.runInTx(ctx, pTx, clog, func(tx pgx.Tx) (rollback bool, err error) {
		defer func() {
			if err != nil {
				item = nil
			}
		}()
		item, err = d.organizationAddAtomic(ctx, tx, name, label, url, publicKey, entity.EntityStateEnabled)
		if err != nil {
			eMsg := "error in d.organizationAddAtomic"
			clog.WithError(err).Error(eMsg)
			rollback = true
			err = errors.Wrap(err, eMsg)
			return
		}
		return false, nil
	})
	return
}
func (d *PgAccess) OrganizationUpdate(ctx context.Context, pTx pgx.Tx, item *entity.Organization, name, label, url, publicKey string) (err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.OrganizationUpdate",
	})
	err = d.runInTx(ctx, pTx, clog, func(tx pgx.Tx) (rollback bool, err error) {
		err = d.organizationUpdateAtomic(ctx, tx, item, name, label, url, publicKey, item.State)
		if err != nil {
			eMsg := "error in d.organizationUpdateAtomic"
			clog.WithError(err).Error(eMsg)
			rollback = true
			err = errors.Wrap(err, eMsg)
			return
		}
		return false, nil
	})
	if err != nil {
		eMsg := "error in d.runInTx()"
		clog.WithError(err).Error(eMsg)
	}
	return
}
func (d *PgAccess) OrganizationChangeState(ctx context.Context, pTx pgx.Tx, item *entity.Organization, state entity.EntityState) (err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.OrganizationChangeState",
	})
	err = d.runInTx(ctx, pTx, clog, func(tx pgx.Tx) (rollback bool, err error) {
		err = d.organizationUpdateAtomic(ctx, tx, item, item.Name, item.Label, item.Url, item.PublicKey, state)
		if err != nil {
			eMsg := "error in d.organizationUpdateAtomic"
			clog.WithError(err).Error(eMsg)
			rollback = true
			err = errors.Wrap(err, eMsg)
			return
		}
		return false, nil
	})
	if err != nil {
		eMsg := "error in d.runInTx()"
		clog.WithError(err).Error(eMsg)
	}
	return
}
func (d *PgAccess) OrganizationById(ctx context.Context, id int) (item *entity.Organization, err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.OrganizationById",
	})
	err = d.runQuery(ctx, clog, func(conn *pgxpool.Conn) (err error) {
		defer func() {
			if err != nil {
				item = nil
			}
		}()
		item = &entity.Organization{}
		//sqlOrganizationById = `SELECT id, name, label, url, public_key, state, create_ts, update_ts, version FROM tbl_organization WHERE id=$1 AND state!=$2`
		row := conn.QueryRow(ctx, sqlOrganizationById, id, entity.EntityStateDeleted)
		err = row.Scan(&item.Id, &item.Name, &item.Label, &item.Url, &item.PublicKey, &item.State, &item.CreateTs, &item.UpdateTs, &item.Version)
		if err != nil {
			if err == pgx.ErrNoRows {
				err = nil
				item = nil
				return
			}
			eMsg := "error in sqlOrganizationById"
			clog.WithError(err).Error(eMsg)
			err = errors.Wrap(err, eMsg)
			return
		}
		return nil
	})
	if err != nil {
		eMsg := "error in pgxAccess.runInQuery"
		clog.WithError(err).Error(eMsg)
	}
	return
}
func (d *PgAccess) OrganizationList(ctx context.Context) (items []*entity.Organization, err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.OrganizationList",
	})
	err = d.runQuery(ctx, clog, func(conn *pgxpool.Conn) (err error) {
		defer func() {
			if err != nil {
				items = nil
			}
		}()
		items = make([]*entity.Organization, 0)
		//sqlOrganizationByList = `SELECT id, name, label, url, public_key, state, create_ts, update_ts, version FROM tbl_organization WHERE state!=$1 ORDER BY id ASC`
		rows, err := conn.Query(ctx, sqlOrganizationByList, entity.EntityStateDeleted)
		if err != nil {
			eMsg := "error in sqlOrganizationByList"
			clog.WithError(err).Error(eMsg)
			err = errors.Wrap(err, eMsg)
			return
		}
		for rows.Next() {
			item := &entity.Organization{}
			err = rows.Scan(&item.Id, &item.Name, &item.Label, &item.Url, &item.PublicKey, &item.State, &item.CreateTs, &item.UpdateTs, &item.Version)
			if err != nil {
				eMsg := "error in rows.Scan"
				clog.WithError(err).Error(eMsg)
				err = errors.Wrap(err, eMsg)
				return
			}
			items = append(items, item)
		}
		return nil
	})
	if err != nil {
		eMsg := "error in pgxAccess.runInQuery"
		clog.WithError(err).Error(eMsg)
	}
	return
}
