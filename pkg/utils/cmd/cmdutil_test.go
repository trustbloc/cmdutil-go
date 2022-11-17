/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package cmd

import (
	"fmt"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

const (
	flagName = "host-url"
	envKey   = "TEST_HOST_URL"
)

func TestGetStringNegative(t *testing.T) {
	command := &cobra.Command{
		Use:   "start",
		Short: "short usage",
		Long:  "long usage",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	// test missing both command line argument and environment vars
	env, err := GetString(command, flagName, envKey, false)
	require.Error(t, err)
	require.Empty(t, env)
	require.Contains(t, err.Error(), "TEST_HOST_URL (environment variable) have been set.")

	// test env var is empty
	t.Setenv(envKey, "")

	env, err = GetString(command, flagName, envKey, false)
	require.Error(t, err)
	require.Contains(t, err.Error(), "TEST_HOST_URL value is empty")
	require.Empty(t, env)

	// test arg is empty
	command.Flags().StringP(flagName, "", "initial", "")
	args := []string{"--" + flagName, ""}
	command.SetArgs(args)
	err = command.Execute()
	require.NoError(t, err)

	env, err = GetString(command, flagName, envKey, false)
	require.Error(t, err)
	require.Contains(t, err.Error(), "host-url value is empty")
	require.Empty(t, env)
}

func TestGetStringArrayNegative(t *testing.T) {
	command := &cobra.Command{
		Use:   "start",
		Short: "short usage",
		Long:  "long usage",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	// test missing both command line argument and environment vars
	env, err := GetStringArray(command, flagName, envKey, false)
	require.Error(t, err)
	require.Empty(t, env)
	require.Contains(t, err.Error(), "TEST_HOST_URL (environment variable) have been set.")

	// test env var is empty
	t.Setenv(envKey, "")

	env, err = GetStringArray(command, flagName, envKey, false)
	require.Error(t, err)
	require.Contains(t, err.Error(), "TEST_HOST_URL value is empty")
	require.Empty(t, env)

	// test arg is empty
	command.Flags().StringArrayP(flagName, "", []string{}, "")
	args := []string{"--" + flagName, ""}
	command.SetArgs(args)
	err = command.Execute()
	require.NoError(t, err)

	env, err = GetStringArray(command, flagName, envKey, false)
	require.Error(t, err)
	require.Contains(t, err.Error(), "host-url value is empty")
	require.Empty(t, env)
}

func TestGetString(t *testing.T) {
	command := &cobra.Command{
		Use:   "start",
		Short: "short usage",
		Long:  "long usage",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	// test env var is set
	hostURL := "localhost:8080"
	t.Setenv(envKey, hostURL)

	// test resolution via environment variable
	env, err := GetString(command, flagName, envKey, false)
	require.NoError(t, err)
	require.Equal(t, hostURL, env)

	// set command line arguments
	command.Flags().StringP(flagName, "", "initial", "")
	args := []string{"--" + flagName, "other"}
	command.SetArgs(args)
	err = command.Execute()
	require.NoError(t, err)

	t.Run("test resolution via command line argument - no environment variable set", func(t *testing.T) {
		env, err = GetString(command, flagName, "", false)
		require.NoError(t, err)
		require.Equal(t, "other", env)

		env = GetOptionalString(command, flagName, "")
		require.Equal(t, "other", env)
	})
}

func TestGetBool(t *testing.T) {
	command := &cobra.Command{
		Use:   "start",
		Short: "short usage",
		Long:  "long usage",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	t.Run("test unset value should use defaultValue", func(t *testing.T) {
		env, err := GetBool(command, flagName, envKey, false, true)
		require.NoError(t, err)
		require.False(t, env)
	})

	t.Run("test env var is set", func(t *testing.T) {
		someIntVal := true
		t.Setenv(envKey, fmt.Sprint(someIntVal))

		// test resolution via environment variable
		env, err := GetBool(command, flagName, envKey, true, true)
		require.NoError(t, err)
		require.Equal(t, someIntVal, env)
	})

	t.Run("test invalid env var", func(t *testing.T) {
		t.Setenv(envKey, "not-an-int")

		env, err := GetBool(command, flagName, envKey, true, true)
		require.Error(t, err)
		require.Empty(t, env)
	})

	t.Run("test invalid flag and envKey for mandatory key", func(t *testing.T) {
		env, err := GetBool(command, "invalidFlag", "invalidEnvKey", true, false)
		require.Error(t, err)
		require.Empty(t, env)
	})
}

func TestGetDuration(t *testing.T) {
	command := &cobra.Command{
		Use:   "start",
		Short: "short usage",
		Long:  "long usage",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	defaultDuration := 10 * time.Second

	t.Run("test unset value should use defaultValue", func(t *testing.T) {
		env, err := GetDuration(command, flagName, envKey, defaultDuration, true)
		require.NoError(t, err)
		require.Equal(t, defaultDuration, env)
	})

	t.Run("test env var is set", func(t *testing.T) {
		someIntVal := 15 * time.Second
		t.Setenv(envKey, fmt.Sprint(someIntVal))

		// test resolution via environment variable
		env, err := GetDuration(command, flagName, envKey, defaultDuration, true)
		require.NoError(t, err)
		require.Equal(t, someIntVal, env)
	})

	t.Run("test invalid env var", func(t *testing.T) {
		t.Setenv(envKey, "not-an-int")

		env, err := GetDuration(command, flagName, envKey, defaultDuration, true)
		require.Error(t, err)
		require.Less(t, env, 0*time.Second)
	})

	t.Run("test invalid flag and envKey for mandatory key", func(t *testing.T) {
		env, err := GetDuration(command, "invalidFlag", "invalidEnvKey", defaultDuration, false)
		require.Error(t, err)
		require.Less(t, env, 0*time.Second)
	})
}

func TestGetInt(t *testing.T) {
	command := &cobra.Command{
		Use:   "start",
		Short: "short usage",
		Long:  "long usage",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	t.Run("test unset value should use defaultValue", func(t *testing.T) {
		env, err := GetInt(command, flagName, envKey, 0, true)
		require.NoError(t, err)
		require.Equal(t, 0, env)
	})

	t.Run("test env var is set", func(t *testing.T) {
		someIntVal := 15
		t.Setenv(envKey, fmt.Sprint(someIntVal))

		// test resolution via environment variable
		env, err := GetInt(command, flagName, envKey, 0, true)
		require.NoError(t, err)
		require.Equal(t, someIntVal, env)
	})

	t.Run("test invalid env var", func(t *testing.T) {
		t.Setenv(envKey, "not-an-int")

		env, err := GetInt(command, flagName, envKey, 0, true)
		require.Error(t, err)
		require.Empty(t, env)
	})

	t.Run("test invalid flag and envKey for mandatory key", func(t *testing.T) {
		env, err := GetInt(command, "invalidFlag", "invalidEnvKey", 0, false)
		require.Error(t, err)
		require.Empty(t, env)
	})
}

func TestGetTLS(t *testing.T) {
	command := &cobra.Command{
		Use:   "start",
		Short: "short usage",
		Long:  "long usage",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	defaultTLArgs := &TLSFields{
		SystemCertPoolFlagName: "",
		SystemCertPoolEnvKey:   "",
		CACertsFlagName:        "",
		CACertsEnvKey:          "",
		CertificateFlagName:    "",
		CertificateLEnvKey:     "",
		KeyFlagName:            "",
		KeyEnvKey:              "",
	}

	emptyTLSParams := &TLSParameters{
		SystemCertPool: false,
		CACerts:        []string{},
		ServeCertPath:  "",
		ServeKeyPath:   ""}

	t.Run("test unset value should use defaultValue", func(t *testing.T) {
		env, err := GetTLS(command, defaultTLArgs)
		require.NoError(t, err)
		require.Equal(t, emptyTLSParams, env)
	})

	t.Run("test env var is set", func(t *testing.T) {
		someTLSKeys := &TLSFields{
			SystemCertPoolFlagName: "SystemCertPoolFlagName",
			SystemCertPoolEnvKey:   "SystemCertPoolEnvKey",
			CACertsFlagName:        "caCertsFlagName",
			CACertsEnvKey:          "caCertsEnvKey",
			CertificateFlagName:    "certificateFlagName",
			CertificateLEnvKey:     "certificateLEnvKey",
			KeyFlagName:            "keyFlagName",
			KeyEnvKey:              "keyEnvKey",
		}

		sysCertPoolFlagName := true
		t.Setenv(someTLSKeys.SystemCertPoolFlagName, fmt.Sprint(sysCertPoolFlagName))
		sysCertPoolEnv := true
		t.Setenv(someTLSKeys.SystemCertPoolEnvKey, fmt.Sprint(sysCertPoolEnv))
		caCertsFlagName := "c"
		t.Setenv(someTLSKeys.CACertsFlagName, caCertsFlagName)
		caCertsEnvKey := "d"
		t.Setenv(someTLSKeys.CACertsEnvKey, caCertsEnvKey)
		certificateFlagName := "e"
		t.Setenv(someTLSKeys.CertificateFlagName, certificateFlagName)
		certificateLEnvKey := "f"
		t.Setenv(someTLSKeys.CertificateLEnvKey, certificateLEnvKey)
		keyFlagName := "h"
		t.Setenv(someTLSKeys.KeyFlagName, keyFlagName)
		keyEnvKey := "i"
		t.Setenv(someTLSKeys.KeyEnvKey, keyEnvKey)

		// test resolution via environment variable
		env, err := GetTLS(command, someTLSKeys)
		require.NoError(t, err)
		require.Equal(t, sysCertPoolEnv, env.SystemCertPool)
		require.Equal(t, keyEnvKey, env.ServeKeyPath)
		require.Equal(t, certificateLEnvKey, env.ServeCertPath)
	})

	t.Run("test invalid env var", func(t *testing.T) {
		someTLSKeys := &TLSFields{
			SystemCertPoolFlagName: "SystemCertPoolFlagName",
			SystemCertPoolEnvKey:   "SystemCertPoolEnvKey",
			CACertsFlagName:        "CACertsFlagName",
			CACertsEnvKey:          "CACertsEnvKey",
			CertificateFlagName:    "CertificateFlagName",
			CertificateLEnvKey:     "CertificateLEnvKey",
			KeyFlagName:            "KeyFlagName",
			KeyEnvKey:              "KeyEnvKey",
		}

		sysCertPoolFlagName := true
		t.Setenv(someTLSKeys.SystemCertPoolFlagName, fmt.Sprint(sysCertPoolFlagName))
		sysCertPoolEnv := "notBool"
		t.Setenv(someTLSKeys.SystemCertPoolEnvKey, sysCertPoolEnv)

		// test resolution via environment variable
		env, err := GetTLS(command, someTLSKeys)
		require.Error(t, err)
		require.Empty(t, env)
	})
}

func TestGetFloat(t *testing.T) {
	command := &cobra.Command{
		Use:   "start",
		Short: "short usage",
		Long:  "long usage",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	t.Run("test unset value should use defaultValue", func(t *testing.T) {
		env, err := GetFloat(command, flagName, envKey, 0, true)
		require.NoError(t, err)
		require.Equal(t, 0.0, env)
	})

	t.Run("test env var is set", func(t *testing.T) {
		someIntVal := 15.0
		t.Setenv(envKey, fmt.Sprint(someIntVal))

		// test resolution via environment variable
		env, err := GetFloat(command, flagName, envKey, 0, true)
		require.NoError(t, err)
		require.Equal(t, someIntVal, env)
	})

	t.Run("test invalid env var", func(t *testing.T) {
		t.Setenv(envKey, "not-an-int")

		env, err := GetFloat(command, flagName, envKey, 0, true)
		require.Error(t, err)
		require.Empty(t, env)
	})

	t.Run("test invalid flag and envKey for mandatory key", func(t *testing.T) {
		env, err := GetFloat(command, "invalidFlag", "invalidEnvKey", 0, false)
		require.Error(t, err)
		require.Empty(t, env)
	})
}

func TestGetStringArray(t *testing.T) {
	command := &cobra.Command{
		Use:   "start",
		Short: "short usage",
		Long:  "long usage",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	// test env var is set
	hostURL := "localhost:8080"
	t.Setenv(envKey, hostURL)

	// test resolution via environment variable
	env, err := GetStringArray(command, flagName, envKey, false)
	require.NoError(t, err)
	require.Equal(t, []string{hostURL}, env)

	// set command line arguments
	command.Flags().StringArrayP(flagName, "", []string{}, "")
	args := []string{"--" + flagName, "other", "--" + flagName, "other1"}
	command.SetArgs(args)
	err = command.Execute()
	require.NoError(t, err)

	// test resolution via command line argument - no environment variable set
	env, err = GetStringArray(command, flagName, "", false)
	require.NoError(t, err)
	require.Equal(t, []string{"other", "other1"}, env)

	env = GetOptionalStringArray(command, flagName, "")
	require.Equal(t, []string{"other", "other1"}, env)
}
