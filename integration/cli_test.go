package integration

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCLI(t *testing.T) {
	path := os.Getenv("SERVER_BIN_PATH")
	require.True(t, len(path) > 0)
	cmd := exec.Command(path, "-h")
	_, err := cmd.CombinedOutput()
	require.NoError(t, err)

	port := "8000"
	cmd = exec.Command(path, "-port", port)
	err = cmd.Start()
	pid := cmd.Process.Pid
	require.NoError(t, err)

	url := fmt.Sprintf("http://127.0.0.1:%s", port)
	err = waitForServerReady(3, url)
	require.NoError(t, err)
	cmd = exec.Command("kill", strconv.Itoa(pid))
	_, err = cmd.CombinedOutput()
	require.NoError(t, err)
}
