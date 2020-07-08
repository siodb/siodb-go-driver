// Copyright (C) 2019-2020 Siodb GmbH. All rights reserved.
// Use of this source code is governed by a license that can be found
// in the LICENSE file.

package siodb

import (
	"database/sql/driver"
)

type siodbRows struct {
	sc         *siodbConn
	columnDesc []*ColumnDescription
}

func (rows *siodbRows) Columns() []string {

	var Cols []string

	for _, column := range rows.columnDesc {
		Cols = append(Cols, column.GetName())
	}

	return Cols
}

func (rows *siodbRows) Next(dest []driver.Value) error {

	return rows.sc.readRow(dest, rows.columnDesc)

}

func (rows *siodbRows) Close() (err error) {

	if !rows.sc.completed {
		if _, err = rows.sc.cleanupBuffer(); err != nil {
			return err
		}
	}

	return nil

}
