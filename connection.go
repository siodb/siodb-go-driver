// Copyright (C) 2019-2020 Siodb GmbH. All rights reserved.
// Use of this source code is governed by a license that can be found
// in the LICENSE file.

package siodb

import (
	"context"
	"database/sql/driver"
	"net"
)

type connector struct {
	cfg Config // immutable private copy.
}

type siodbConn struct {
	netConn             net.Conn
	cfg                 Config
	sessionID           string
	RequestID           uint64
	nullAllowed         bool
	nullBitmaskByteSize int
	completed           bool
}

// BeginTx TODO: https://golang.org/pkg/database/sql/driver/#ConnBeginTx
func BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	return nil, nil
}

// Begin TODO: https://golang.org/pkg/database/sql/driver/#Conn
func (sc *siodbConn) Begin() (driver.Tx, error) {
	return nil, nil
}

// Close TODO: Implement proper exit in Siodb
func (sc *siodbConn) Close() (err error) {

	if err := sc.netConn.Close(); err != nil {
		return err
	}
	return nil
}

// PrepareContext TODO: https://golang.org/pkg/database/sql/driver/#ConnPrepareContext
func PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	return nil, nil
}

// Prepare TODO: https://golang.org/pkg/database/sql/driver/#Conn
func (sc *siodbConn) Prepare(query string) (driver.Stmt, error) {
	return nil, nil
}

func checkServerError(Message []*StatusMessage) error {

	if len(Message) > 0 {
		for _, Msg := range Message {
			return &siodbServerError{Msg.StatusCode, Msg.Text}
		}
	}

	return nil
}

func (sc *siodbConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {

	var sr ServerResponse
	var err error

	// TODO: Bind Values

	if err = sc.writeServerCommand(query); err != nil {
		return nil, &siodbDriverError{"Fail to write server command."}
	}

	if sr, err = sc.readServer(); err != nil {
		return nil, &siodbDriverError{"Fail to read server response."}
	}

	if err = checkServerError(sr.Message); err != nil {
		return nil, err
	}

	var AffectedRowCount int64
	if sr.HasAffectedRowCount {
		AffectedRowCount = int64(sr.AffectedRowCount)
	}

	return &siodbResult{
		AffectedRowCount: AffectedRowCount,
		insertID:         0,
	}, nil
}

func (sc *siodbConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {

	var sr ServerResponse
	var err error

	// TODO: Bind Values

	if err = sc.writeServerCommand(query); err != nil {
		return nil, err
	}

	if sr, err = sc.readServer(); err != nil {
		return nil, err
	}

	if err = checkServerError(sr.Message); err != nil {
		return nil, err
	}

	// Init rows struct for further next()
	rows := new(siodbRows)
	rows.sc = sc
	rows.columnDesc = sr.ColumnDescription

	return rows, err
}
