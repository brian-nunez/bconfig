package file

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestFileSource_Name(t *testing.T) {
	path := "test.json"
	src := Source(path)
	if src.Name() != "file:test.json" {
		t.Errorf("expected Name to be 'file:test.json', got %q", src.Name())
	}
}

func TestFileSource_Load_JSON(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "config.json")
	content := `{"server": {"addr": ":8080", "port": 8080}, "debug": true}`
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	src := Source(filePath)
	data, err := src.Load(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	serverVal, exists := data["server"]
	if !exists {
		t.Fatalf("expected 'server' key to exist")
	}
	serverMap, ok := serverVal.(map[string]any)
	if !ok {
		t.Fatalf("expected 'server' to be map[string]any, got %T", serverVal)
	}
	if serverMap["addr"] != ":8080" {
		t.Errorf("expected server.addr to be ':8080', got %v", serverMap["addr"])
	}
	if serverMap["port"] != 8080.0 && serverMap["port"] != int(8080) { // json unmarshal may produce float64
		t.Errorf("expected server.port to be 8080, got %v (%T)", serverMap["port"], serverMap["port"])
	}
	if data["debug"] != true {
		t.Errorf("expected debug to be true, got %v", data["debug"])
	}
}

func TestFileSource_Load_YAML(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "config.yaml")
	content := `
server:
  addr: ":8080"
  port: 8080
debug: true
`
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	src := Source(filePath)
	data, err := src.Load(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	serverVal, exists := data["server"]
	if !exists {
		t.Fatalf("expected 'server' key to exist")
	}
	serverMap, ok := serverVal.(map[string]any)
	if !ok {
		t.Fatalf("expected 'server' to be map[string]any, got %T", serverVal)
	}
	if serverMap["addr"] != ":8080" {
		t.Errorf("expected server.addr to be ':8080', got %v", serverMap["addr"])
	}
	if data["debug"] != true {
		t.Errorf("expected debug to be true, got %v", data["debug"])
	}
}

func TestFileSource_Load_YML(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "config.yml")
	content := "foo: bar"
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	src := Source(filePath)
	data, err := src.Load(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if data["foo"] != "bar" {
		t.Errorf("expected foo to be 'bar', got %v", data["foo"])
	}
}

func TestFileSource_Load_EmptyFile(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "config.json")
	if err := os.WriteFile(filePath, []byte(""), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	src := Source(filePath)
	data, err := src.Load(context.Background())
	if err != nil {
		t.Fatalf("expected no error for empty JSON file, got %v", err)
	}
	if len(data) != 0 {
		t.Errorf("expected empty data, got %v", data)
	}
}

func TestFileSource_Load_NullJSONFile(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "config.json")
	if err := os.WriteFile(filePath, []byte("null"), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	src := Source(filePath)
	data, err := src.Load(context.Background())
	if err != nil {
		t.Fatalf("expected no error for null JSON, got %v", err)
	}
	if len(data) != 0 {
		t.Errorf("expected empty data, got %v", data)
	}
}

func TestFileSource_Load_UnsupportedExtension(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "config.toml")
	if err := os.WriteFile(filePath, []byte("foo = 'bar'"), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	src := Source(filePath)
	_, err := src.Load(context.Background())
	if err == nil {
		t.Fatalf("expected error for unsupported extension, got nil")
	}
	expectedErr := `unsupported config file extension ".toml"`
	if err.Error() != expectedErr {
		t.Errorf("expected error %q, got %q", expectedErr, err.Error())
	}
}

func TestFileSource_Load_MissingFile(t *testing.T) {
	src := Source("non_existent_file.json")
	_, err := src.Load(context.Background())
	if err == nil {
		t.Fatalf("expected error for missing file, got nil")
	}
}

func TestFileSource_Load_InvalidJSON(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "config.json")
	if err := os.WriteFile(filePath, []byte("{invalid-json}"), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	src := Source(filePath)
	_, err := src.Load(context.Background())
	if err == nil {
		t.Fatalf("expected error for invalid JSON, got nil")
	}
}

func TestFileSource_Load_InvalidYAML(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "config.yaml")
	if err := os.WriteFile(filePath, []byte("invalid:\n  - yaml\n - mismatch-indent"), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	src := Source(filePath)
	_, err := src.Load(context.Background())
	if err == nil {
		t.Fatalf("expected error for invalid YAML, got nil")
	}
}
