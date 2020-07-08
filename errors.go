// Copyright (C) 2019-2020 Siodb GmbH. All rights reserved.
// Use of this source code is governed by a license that can be found
// in the LICENSE file.

package siodb

import (
	"fmt"
)

type siodbDriverError struct {
	Message string
}

type siodbServerError struct {
	Number  int32
	Message string
}

func (sde *siodbDriverError) Error() string {
	return fmt.Sprintf("Siodb Driver Error: %s", sde.Message)
}

func (sse *siodbServerError) Error() string {
	return fmt.Sprintf("Siodb Server Error: %d | %s", sse.Number, sse.Message)
}
