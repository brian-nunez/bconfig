package env

import (
	"context"
	"testing"
)

func TestEnvSource_Name(t *testing.T) {
	if Source().Name() != "env" {
		t.Errorf("expected empty prefix name to be 'env', got %q", Source().Name())
	}
	if Source("BAPP_").Name() != "env:BAPP_" {
		t.Errorf("expected name with prefix to be 'env:BAPP_', got %q", Source("BAPP_").Name())
	}
	if Source("BAPP_", "OTHER_").Name() != "env:BAPP_,OTHER_" {
		t.Errorf("expected name with prefixes to be 'env:BAPP_,OTHER_', got %q", Source("BAPP_", "OTHER_").Name())
	}
}

func TestEnvSource_Load_All(t *testing.T) {
	t.Setenv("TEST_ENV_ALL_VAR", "value")

	src := Source()
	data, err := src.Load(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	val, exists := data["test_env_all_var"]
	if !exists {
		t.Fatalf("expected 'test_env_all_var' key to exist")
	}
	if val != "value" {
		t.Errorf("expected 'test_env_all_var' to be 'value', got %v", val)
	}
}

func TestEnvSource_Load_WithPrefix(t *testing.T) {
	t.Setenv("BAPP_SERVER_ADDR", ":8080")
	t.Setenv("OTHER_VAR", "ignored")

	src := Source("BAPP_")
	data, err := src.Load(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if _, exists := data["other_var"]; exists {
		t.Errorf("expected other_var to be ignored, but it was loaded")
	}

	val, exists := data["bapp_server_addr"]
	if !exists {
		t.Fatalf("expected 'bapp_server_addr' key to exist")
	}
	if val != ":8080" {
		t.Errorf("expected 'bapp_server_addr' to be ':8080', got %v", val)
	}
}

func TestEnvSource_Load_MultiplePrefixes(t *testing.T) {
	t.Setenv("BAPP_SERVER_PORT", "8080")
	t.Setenv("OTHER_APP_NAME", "my-app")
	t.Setenv("UNMATCHED_VAR", "ignored")

	src := Source("BAPP_", "OTHER_")
	data, err := src.Load(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if _, exists := data["unmatched_var"]; exists {
		t.Errorf("expected unmatched_var to be ignored")
	}

	port, exists := data["bapp_server_port"]
	if !exists {
		t.Fatalf("expected 'bapp_server_port' key to exist")
	}
	if port != 8080 {
		t.Errorf("expected bapp_server_port to be 8080, got %v", port)
	}

	name, exists := data["other_app_name"]
	if !exists {
		t.Fatalf("expected 'other_app_name' key to exist")
	}
	if name != "my-app" {
		t.Errorf("expected other_app_name to be 'my-app', got %v", name)
	}
}

func TestEnvSource_ValueParsing(t *testing.T) {
	t.Setenv("BAPP_BOOL_TRUE", "true")
	t.Setenv("BAPP_BOOL_FALSE", "false")
	t.Setenv("BAPP_INT_VAL", "42")
	t.Setenv("BAPP_FLOAT_VAL", "3.14")
	t.Setenv("BAPP_STR_VAL", "hello")
	t.Setenv("BAPP_EMPTY_VAL", "")

	src := Source("BAPP_")
	data, err := src.Load(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if data["bapp_bool_true"] != true {
		t.Errorf("expected bapp_bool_true to be true")
	}
	if data["bapp_bool_false"] != false {
		t.Errorf("expected bapp_bool_false to be false")
	}
	if data["bapp_int_val"] != 42 {
		t.Errorf("expected bapp_int_val to be 42")
	}
	if data["bapp_float_val"] != 3.14 {
		t.Errorf("expected bapp_float_val to be 3.14")
	}
	if data["bapp_str_val"] != "hello" {
		t.Errorf("expected bapp_str_val to be 'hello'")
	}
	if data["bapp_empty_val"] != "" {
		t.Errorf("expected bapp_empty_val to be empty string")
	}
}
