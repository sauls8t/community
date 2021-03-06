// Copyright 2016 Documize Inc. <legal@documize.com>. All rights reserved.
//
// This software (Documize Community Edition) is licensed under
// GNU AGPL v3 http://www.gnu.org/licenses/agpl-3.0.en.html
//
// You can operate outside the AGPL restrictions by purchasing
// Documize Enterprise Edition and obtaining a commercial license
// by contacting <sales@documize.com>.
//
// https://documize.com

package store

import (
	"database/sql"
	"fmt"
	"github.com/documize/community/core/env"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Context provides common stuff to store providers.
type Context struct {
	Runtime *env.Runtime
}

// Bind selects query parameter placeholder for given database provider.
//
// MySQL uses ?, ?, ? (default for all Documize queries).``
// PostgreSQL uses $1, $2, $3.
// MS SQL Server uses @p1, @p2, @p3.
func (c *Context) Bind(sql string) string {
	// Default to MySQL.
	bindParam := sqlx.QUESTION

	switch c.Runtime.StoreProvider.Type() {
	case env.StoreTypePostgreSQL:
		bindParam = sqlx.DOLLAR
	}

	return sqlx.Rebind(bindParam, sql)
}

// Delete record.
func (c *Context) Delete(tx *sqlx.Tx, table string, id string) (rows int64, err error) {
	_, err = tx.Exec(c.Bind("DELETE FROM "+table+" WHERE c_refid=?"), id)
	if err == sql.ErrNoRows {
		err = nil
	}
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("unable to delete row in table %s", table))
		return
	}

	return
}

// DeleteConstrained record constrained to Organization using refid.
func (c *Context) DeleteConstrained(tx *sqlx.Tx, table string, orgID, id string) (rows int64, err error) {
	_, err = tx.Exec(c.Bind("DELETE FROM "+table+" WHERE c_orgid=? AND c_refid=?"), orgID, id)
	if err == sql.ErrNoRows {
		err = nil
	}
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("unable to delete row in table %s", table))
		return
	}

	return
}

// DeleteConstrainedWithID record constrained to Organization using non refid.
func (c *Context) DeleteConstrainedWithID(tx *sqlx.Tx, table string, orgID, id string) (rows int64, err error) {
	_, err = tx.Exec(c.Bind("DELETE FROM "+table+" WHERE c_orgid=? AND id=?"), orgID, id)
	if err == sql.ErrNoRows {
		err = nil
	}
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("unable to delete rows by id: %s", id))
		return
	}

	return
}

// DeleteWhere free form query.
func (c *Context) DeleteWhere(tx *sqlx.Tx, statement string) (rows int64, err error) {
	_, err = tx.Exec(statement)
	if err == sql.ErrNoRows {
		err = nil
	}
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("unable to delete rows: %s", statement))
		return
	}

	return
}

// EmptyJSON returns database specific empty JSON object.
func (c *Context) EmptyJSON() string {
	return c.Runtime.StoreProvider.JSONEmpty()
}

// GetJSONValue returns database specific empty JSON object.
func (c *Context) GetJSONValue(column, attribute string) string {
	return c.Runtime.StoreProvider.JSONGetValue(column, attribute)
}
