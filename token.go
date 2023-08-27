/*
 * Copyright (c) 2023 Zander Schwid & Co. LLC.
 * SPDX-License-Identifier: BUSL-1.1
 */

package recordbase

import "context"

type TokenAuth string

func (t TokenAuth) GetRequestMetadata(_ context.Context, in ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Bearer " + string(t),
	}, nil
}

func (TokenAuth) RequireTransportSecurity() bool {
	return false
}

