/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package logfields

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trustbloc/logutil-go/pkg/log"
)

//nolint:maintidx
func TestStandardFields(t *testing.T) {
	const (
		module = "test_module"
	)

	t.Run("json fields", func(t *testing.T) {
		stdOut := newMockWriter()

		logger := log.New(module, log.WithStdOut(stdOut), log.WithEncoding(log.JSON))

		certPoolSize := 10

		logger.Info(
			"Some message",
			WithCertPoolSize(certPoolSize),
		)

		l := unmarshalLogData(t, stdOut.Bytes())

		require.Equal(t, certPoolSize, l.CertPoolSize)
	})
}

type logData struct {
	Level  string `json:"level"`
	Time   string `json:"time"`
	Logger string `json:"logger"`
	Caller string `json:"caller"`
	Msg    string `json:"msg"`
	Error  string `json:"error"`

	CertPoolSize int `json:"certPoolSize"`
}

func unmarshalLogData(t *testing.T, b []byte) *logData {
	t.Helper()

	l := &logData{}

	require.NoError(t, json.Unmarshal(b, l))

	return l
}

type mockWriter struct {
	*bytes.Buffer
}

func (m *mockWriter) Sync() error {
	return nil
}

func newMockWriter() *mockWriter {
	return &mockWriter{Buffer: bytes.NewBuffer(nil)}
}
