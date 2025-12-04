package config

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTempConfig создает временный файл конфигурации для тестов
func createTempConfig(t *testing.T, content string) (string, func()) {
	t.Helper()

	tmpfile, err := os.CreateTemp("", "config-*.yaml")
	require.NoError(t, err)

	if content != "" {
		_, err = tmpfile.WriteString(content)
		require.NoError(t, err)
	}

	err = tmpfile.Close()
	require.NoError(t, err)

	return tmpfile.Name(), func() {
		os.Remove(tmpfile.Name())
	}
}

func TestMustLoad_Success(t *testing.T) {
	yamlContent := `
server:
  host: "127.0.0.1"
  port: "9090"
`
	configPath, cleanup := createTempConfig(t, yamlContent)
	defer cleanup()

	os.Setenv("APP_CONFIG_PATH", configPath)
	defer os.Unsetenv("APP_CONFIG_PATH")

	cfg := MustLoad()

	assert.NotNil(t, cfg)
	assert.Equal(t, "127.0.0.1", cfg.Server.Host)
	assert.Equal(t, "9090", cfg.Server.Port)
}

func TestMustLoad_Defaults(t *testing.T) {
	yamlContent := `
server:
`
	configPath, cleanup := createTempConfig(t, yamlContent)
	defer cleanup()

	os.Setenv("APP_CONFIG_PATH", configPath)
	defer os.Unsetenv("APP_CONFIG_PATH")

	cfg := MustLoad()

	assert.NotNil(t, cfg)
	assert.Equal(t, "localhost", cfg.Server.Host)
	assert.Equal(t, "8080", cfg.Server.Port)
}

func TestMustLoad_FileNotFound(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		os.Setenv("APP_CONFIG_PATH", "non_existent_file.yaml")
		MustLoad()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestMustLoad_FileNotFound")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()

	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestMustLoad_InvalidConfig(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		// Создаем битый YAML
		f, _ := os.CreateTemp("", "bad-*.yaml")
		f.WriteString("server: [broken_yaml")
		f.Close()
		defer os.Remove(f.Name())

		os.Setenv("APP_CONFIG_PATH", f.Name())
		MustLoad()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestMustLoad_InvalidConfig")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()

	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}
