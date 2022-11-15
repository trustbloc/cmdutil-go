/*
   Copyright SecureKey Technologies Inc.

   This file contains software code that is the intellectual property of SecureKey.
   SecureKey reserves all rights in the code and you may not use it without
	 written permission from SecureKey.
*/

package cmd

import (
	"testing"

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

	// test resolution via command line argument - no environment variable set
	env, err = GetString(command, flagName, "", false)
	require.NoError(t, err)
	require.Equal(t, "other", env)

	env = GetOptionalString(command, flagName, "")
	require.Equal(t, "other", env)
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
