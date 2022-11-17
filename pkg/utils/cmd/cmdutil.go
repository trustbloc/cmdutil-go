/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// GetOptionalString returns values either command line flag or environment variable.
func GetOptionalString(cmd *cobra.Command, flagName, envKey string) string {
	v, _ := GetString(cmd, flagName, envKey, true)

	return v
}

// GetString returns values either command line flag or environment variable.
func GetString(cmd *cobra.Command, flagName, envKey string, isOptional bool) (string, error) {
	if cmd.Flags().Changed(flagName) {
		value, err := cmd.Flags().GetString(flagName)
		if err != nil {
			return "", fmt.Errorf(flagName+" flag not found: %s", err)
		}

		if value == "" {
			return "", fmt.Errorf("%s value is empty", flagName)
		}

		return value, nil
	}

	value, isSet := os.LookupEnv(envKey)

	if isOptional || isSet {
		if !isOptional && value == "" {
			return "", fmt.Errorf("%s value is empty", envKey)
		}

		return value, nil
	}

	return "", errors.New("Neither " + flagName + " (command line flag) nor " + envKey +
		" (environment variable) have been set.")
}

// GetOptionalStringArray returns values either command line flag or environment variable.
func GetOptionalStringArray(cmd *cobra.Command, flagName, envKey string) []string {
	v, _ := GetStringArray(cmd, flagName, envKey, true)

	return v
}

// GetStringArray returns values either command line flag or environment variable.
func GetStringArray(cmd *cobra.Command, flagName, envKey string, isOptional bool) ([]string, error) {
	if cmd.Flags().Changed(flagName) {
		value, err := cmd.Flags().GetStringArray(flagName)
		if err != nil {
			return nil, fmt.Errorf(flagName+" flag not found: %s", err)
		}

		if len(value) == 0 {
			return nil, fmt.Errorf("%s value is empty", flagName)
		}

		return value, nil
	}

	value, isSet := os.LookupEnv(envKey)

	if isOptional || isSet {
		if !isOptional && value == "" {
			return nil, fmt.Errorf("%s value is empty", envKey)
		}

		if value == "" {
			return []string{}, nil
		}

		return strings.Split(value, ","), nil
	}

	return nil, errors.New("Neither " + flagName + " (command line flag) nor " + envKey +
		" (environment variable) have been set.")
}

// GetBool returns values either command line flag or environment variable.
func GetBool(cmd *cobra.Command, flagName, envKey string, defaultValue, isOptional bool) (bool, error) {
	str, err := GetUserSetVarFromString(cmd, flagName, envKey, isOptional)
	if err != nil {
		return false, fmt.Errorf("%s: %w", flagName, err)
	}

	if str == "" {
		return defaultValue, nil
	}

	value, err := strconv.ParseBool(str)
	if err != nil {
		return false, fmt.Errorf("invalid value for %s [%s]: %w", flagName, str, err)
	}

	return value, nil
}

// GetDuration returns values either command line flag or environment variable.
func GetDuration(cmd *cobra.Command, flagName, envKey string,
	defaultDuration time.Duration, isOptional bool) (time.Duration, error) {
	timeoutStr, err := GetUserSetVarFromString(cmd, flagName, envKey, isOptional)
	if err != nil {
		return -1, err
	}

	if timeoutStr == "" {
		return defaultDuration, nil
	}

	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		return -1, fmt.Errorf("invalid value [%s]: %w", timeoutStr, err)
	}

	return timeout, nil
}

// GetInt returns values either command line flag or environment variable.
func GetInt(cmd *cobra.Command, flagName, envKey string, defaultValue int, isOptional bool) (int, error) {
	str, err := GetUserSetVarFromString(cmd, flagName, envKey, isOptional)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", flagName, err)
	}

	if str == "" {
		return defaultValue, nil
	}

	value, err := strconv.Atoi(str)
	if err != nil {
		return 0, fmt.Errorf("invalid value for %s [%s]: %w", flagName, str, err)
	}

	return value, nil
}

// GetFloat returns values either command line flag or environment variable.
func GetFloat(cmd *cobra.Command, flagName, envKey string, defaultValue float64, isOptional bool) (float64, error) {
	str, err := GetUserSetVarFromString(cmd, flagName, envKey, isOptional)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", flagName, err)
	}

	if str == "" {
		return defaultValue, nil
	}

	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid value for %s [%s]: %w", flagName, str, err)
	}

	return value, nil
}

// TLSParameters contains TLS parameters retrieved from command arguments.
type TLSParameters struct {
	SystemCertPool bool
	CACerts        []string
	ServeCertPath  string
	ServeKeyPath   string
}

// TLSFields contains TLS command field names and environment variables key.
type TLSFields struct {
	SystemCertPoolFlagName string
	SystemCertPoolEnvKey   string
	CACertsFlagName        string
	CACertsEnvKey          string
	CertificateFlagName    string
	CertificateLEnvKey     string
	KeyFlagName            string
	KeyEnvKey              string
}

// GetTLS returns values either command line flag or environment variable.
func GetTLS(cmd *cobra.Command, tlsFields *TLSFields) (*TLSParameters, error) {
	tlsSystemCertPoolString := GetUserSetOptionalVarFromString(cmd, tlsFields.SystemCertPoolFlagName,
		tlsFields.SystemCertPoolEnvKey)

	tlsSystemCertPool := false

	if tlsSystemCertPoolString != "" {
		var err error

		tlsSystemCertPool, err = strconv.ParseBool(tlsSystemCertPoolString)
		if err != nil {
			return nil, err
		}
	}

	tlsCACerts := GetUserSetOptionalVarFromArrayString(cmd, tlsFields.CACertsFlagName, tlsFields.CACertsEnvKey)

	tlsServeCertPath := GetUserSetOptionalVarFromString(cmd, tlsFields.CertificateFlagName, tlsFields.CertificateLEnvKey)

	tlsServeKeyPath := GetUserSetOptionalVarFromString(cmd, tlsFields.KeyFlagName, tlsFields.KeyEnvKey)

	return &TLSParameters{
		SystemCertPool: tlsSystemCertPool,
		CACerts:        tlsCACerts,
		ServeCertPath:  tlsServeCertPath,
		ServeKeyPath:   tlsServeKeyPath,
	}, nil
}
