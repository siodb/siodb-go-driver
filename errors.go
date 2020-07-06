// Copyright (C) 2019-2020 Siodb GmbH. All rights reserved.
// Use of this source code is governed by a license that can be found
// in the LICENSE file.

package siodb

import (
	"fmt"
)

type SiodbDriverError struct {
	Message string
}

type SiodbServerError struct {
	Number  int32
	Message string
}

func (sde *SiodbDriverError) Error() string {
	return fmt.Sprintf("Siodb Driver Error: %s", sde.Message)
}

func (sse *SiodbServerError) Error() string {
	return fmt.Sprintf("Siodb Server Error: %d | %s", sse.Number, sse.Message)
}
