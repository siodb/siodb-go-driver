// Copyright (C) 2019-2020 Siodb GmbH. All rights reserved.
// Use of this source code is governed by a license that can be found
// in the LICENSE file.

package siodb

type siodbResult struct {
	AffectedRowCount int64
	insertID         int64
}

func (res *siodbResult) LastInsertId() (int64, error) {
	return res.insertID, nil
}

func (res *siodbResult) RowsAffected() (int64, error) {
	return res.AffectedRowCount, nil
}
