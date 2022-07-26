package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testUserEnvName  = "TEST_USER_ENV"
	testPwdEnvName   = "TEST_PWD_ENV"
	testUserEnvValue = "user"
	testPwdEnvValue  = "pass"
)

func TestRunCmd(t *testing.T) {
	env := make(Environment)
	env[testUserEnvName] = EnvValue{Value: testUserEnvValue, NeedRemove: false}
	env[testPwdEnvName] = EnvValue{Value: testPwdEnvValue, NeedRemove: true}

	t.Run("test RunCmd", func(t *testing.T) {
		var buf bytes.Buffer
		expectedExitCode := 0
		expectedCmd1Output := fmt.Sprintf("%s\n", testUserEnvValue)
		expectedCmd2Output := ""

		_ = os.Setenv(testUserEnvName, "root")
		_ = os.Setenv(testPwdEnvName, testPwdEnvValue)

		cmd1 := []string{
			"printenv",
			testUserEnvName,
		}

		cmd2 := []string{
			"printenv",
			testPwdEnvName,
		}

		exitCode := RunCmd(
			cmd1,
			env,
			os.Stdin,
			&buf,
			os.Stderr,
		)
		require.Equal(t, expectedExitCode, exitCode)
		require.Equal(t, expectedCmd1Output, buf.String())

		buf.Reset()
		exitCode = RunCmd(
			cmd2,
			env,
			os.Stdin,
			&buf,
			os.Stderr,
		)
		require.NotEqual(t, expectedExitCode, exitCode)
		require.Equal(t, expectedCmd2Output, buf.String())

		_, ok := os.LookupEnv(testUserEnvName)
		require.True(t, ok)

		_, ok = os.LookupEnv(testPwdEnvName)
		require.False(t, ok)
	})
}
