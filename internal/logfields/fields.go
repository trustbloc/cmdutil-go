/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package logfields

import (
	"go.uber.org/zap"
)

const (
	// FieldCertPoolSize log field name.
	FieldCertPoolSize = "certPoolSize"
)

// WithCertPoolSize sets the CertPoolSize field.
func WithCertPoolSize(value int) zap.Field {
	return zap.Int(FieldCertPoolSize, value)
}
