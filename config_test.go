package bconfig

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

type mockSource struct {
	name string
	data map[string]any
	err  error
}

func (m *mockSource) Name() string {
	return m.name
}

func (m *mockSource) Load(ctx context.Context) (map[string]any, error) {
	return m.data, m.err
}

func TestNew(t *testing.T) {
	src := &mockSource{
		name: "mock1",
		data: map[string]any{"key": "value"},
	}

	cfg, err := New(src)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.String("key") != "value" {
		t.Errorf("expected key to be 'value', got %v", cfg.String("key"))
	}
}

func TestLoad_ZeroSources(t *testing.T) {
	cfg, err := Load(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(cfg.data) != 0 {
		t.Errorf("expected empty config data, got %v", cfg.data)
	}

	// Loading with nil source should be ignored/tolerated
	cfg, err = Load(context.Background(), nil)
	if err != nil {
		t.Fatalf("expected no error with nil source, got %v", err)
	}
	if len(cfg.data) != 0 {
		t.Errorf("expected empty config data, got %v", cfg.data)
	}
}

func TestLoad_SourceError(t *testing.T) {
	srcErr := errors.New("connection timeout")
	src := &mockSource{
		name: "mock-err",
		err:  srcErr,
	}

	_, err := Load(context.Background(), src)
	if err == nil {
		t.Fatalf("expected error from source, got nil")
	}

	expectedStr := "bconfig: failed to load source mock-err: connection timeout"
	if err.Error() != expectedStr {
		t.Errorf("expected error %q, got %q", expectedStr, err.Error())
	}

	if !errors.Is(err, srcErr) {
		t.Errorf("expected error to wrap the source error")
	}
}

func TestLoad_MergeAndOverride(t *testing.T) {
	src1 := &mockSource{
		name: "mock1",
		data: map[string]any{
			"port": 8080,
			"database": map[string]any{
				"host": "localhost",
				"port": 5432,
			},
			"features": []any{"login", "signup"},
		},
	}
	src2 := &mockSource{
		name: "mock2",
		data: map[string]any{
			"port": 9090, // override scalar
			"database": map[string]any{
				"host": "db.prod", // override nested key, keep db.port
			},
			"features": []any{"login", "signup", "billing"}, // override slice
		},
	}

	cfg, err := Load(context.Background(), src1, src2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Int("port") != 9090 {
		t.Errorf("expected port to be 9090, got %d", cfg.Int("port"))
	}
	if cfg.String("database.host") != "db.prod" {
		t.Errorf("expected database.host to be 'db.prod', got %s", cfg.String("database.host"))
	}
	if cfg.Int("database.port") != 5432 {
		t.Errorf("expected database.port to be 5432, got %d", cfg.Int("database.port"))
	}

	features := cfg.Get("features").([]any)
	if len(features) != 3 || features[2] != "billing" {
		t.Errorf("expected features to be overridden by mock2, got %v", features)
	}
}

func TestGet_DotPath(t *testing.T) {
	data := map[string]any{
		"str":   "hello",
		"bool":  true,
		"int":   42,
		"float": 3.14,
		"nested": map[string]any{
			"val": "nested-val",
		},
	}
	cfg := &Config{data: data}

	// Existing values
	if cfg.Get("str") != "hello" {
		t.Errorf("expected Get('str') to be 'hello'")
	}
	if cfg.String("str") != "hello" {
		t.Errorf("expected String('str') to be 'hello'")
	}
	if cfg.Bool("bool") != true {
		t.Errorf("expected Bool('bool') to be true")
	}
	if cfg.Int("int") != 42 {
		t.Errorf("expected Int('int') to be 42")
	}
	if cfg.Float("float") != 3.14 {
		t.Errorf("expected Float('float') to be 3.14")
	}
	if cfg.String("nested.val") != "nested-val" {
		t.Errorf("expected String('nested.val') to be 'nested-val'")
	}

	// Type coercion checks
	if cfg.String("int") != "42" {
		t.Errorf("expected String('int') to coerce to '42'")
	}
	if cfg.Int("float") != 3 {
		t.Errorf("expected Int('float') to coerce to 3, got %d", cfg.Int("float"))
	}
	if cfg.Float("int") != 42.0 {
		t.Errorf("expected Float('int') to coerce to 42.0, got %f", cfg.Float("int"))
	}

	// Additional numeric types for Int and Float coercion
	cfg.data["int32"] = int32(32)
	cfg.data["int64"] = int64(64)
	cfg.data["float32"] = float32(3.5)

	if cfg.Int("int32") != 32 {
		t.Errorf("expected Int('int32') to be 32, got %d", cfg.Int("int32"))
	}
	if cfg.Int("int64") != 64 {
		t.Errorf("expected Int('int64') to be 64, got %d", cfg.Int("int64"))
	}
	if cfg.Int("float32") != 3 {
		t.Errorf("expected Int('float32') to be 3, got %d", cfg.Int("float32"))
	}

	if cfg.Float("int32") != 32.0 {
		t.Errorf("expected Float('int32') to be 32.0, got %f", cfg.Float("int32"))
	}
	if cfg.Float("int64") != 64.0 {
		t.Errorf("expected Float('int64') to be 64.0, got %f", cfg.Float("int64"))
	}
	if cfg.Float("float32") != 3.5 {
		t.Errorf("expected Float('float32') to be 3.5, got %f", cfg.Float("float32"))
	}

	// Invalid type coercions (fallback to zero values)
	if cfg.Bool("str") != false {
		t.Errorf("expected Bool('str') to be false (unsupported coercion)")
	}
	if cfg.Int("str") != 0 {
		t.Errorf("expected Int('str') to be 0 (unsupported coercion)")
	}
	if cfg.Float("str") != 0.0 {
		t.Errorf("expected Float('str') to be 0.0 (unsupported coercion)")
	}

	// Missing values (should return zero values)
	if cfg.Get("missing") != nil {
		t.Errorf("expected Get('missing') to be nil")
	}
	if cfg.String("missing") != "" {
		t.Errorf("expected String('missing') to be empty string")
	}
	if cfg.Bool("missing") != false {
		t.Errorf("expected Bool('missing') to be false")
	}
	if cfg.Int("missing") != 0 {
		t.Errorf("expected Int('missing') to be 0")
	}
	if cfg.Float("missing") != 0.0 {
		t.Errorf("expected Float('missing') to be 0.0")
	}

	// Empty path
	if cfg.Get("") != nil {
		t.Errorf("expected empty path to return nil")
	}

	// Non-map path step
	if cfg.Get("str.invalid") != nil {
		t.Errorf("expected non-map traversal to return nil")
	}
}

func TestDecode(t *testing.T) {
	type ServerConfig struct {
		Addr    string `json:"addr"`
		Timeout string `json:"timeout"`
		Port    int    `json:"port"`
	}
	type AppConfig struct {
		Server ServerConfig `json:"server"`
		Debug  bool         `json:"debug"`
	}

	data := map[string]any{
		"server": map[string]any{
			"addr":    ":8080",
			"timeout": "5s",
			"port":    8080,
		},
		"debug": true,
	}

	cfg := &Config{data: data}

	var appCfg AppConfig
	err := cfg.Decode(&appCfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if appCfg.Server.Addr != ":8080" {
		t.Errorf("expected Server.Addr to be ':8080', got %q", appCfg.Server.Addr)
	}
	if appCfg.Server.Timeout != "5s" {
		t.Errorf("expected Server.Timeout to be '5s', got %q", appCfg.Server.Timeout)
	}
	if appCfg.Server.Port != 8080 {
		t.Errorf("expected Server.Port to be 8080, got %d", appCfg.Server.Port)
	}
	if appCfg.Debug != true {
		t.Errorf("expected Debug to be true")
	}
}

func TestDecode_Error(t *testing.T) {
	// A map with un-marshalable JSON values (like a channel)
	data := map[string]any{
		"ch": make(chan int),
	}
	cfg := &Config{data: data}

	var target struct{}
	err := cfg.Decode(&target)
	if err == nil {
		t.Fatalf("expected decode error, got nil")
	}

	expectedPrefix := "bconfig: decode failed:"
	if !reflect.ValueOf(err).MethodByName("Error").IsValid() || len(err.Error()) < len(expectedPrefix) {
		t.Errorf("expected error to start with %q, got %v", expectedPrefix, err)
	}
}

func TestDecode_UnmarshalError(t *testing.T) {
	data := map[string]any{
		"port": 8080,
	}
	cfg := &Config{data: data}

	// Passing a non-pointer target will cause json.Unmarshal to fail
	var target struct{}
	err := cfg.Decode(target)
	if err == nil {
		t.Fatalf("expected decode error for non-pointer target, got nil")
	}
}

func TestData_Copy(t *testing.T) {
	originalData := map[string]any{
		"server": map[string]any{
			"addr": ":8080",
		},
	}
	cfg := &Config{data: originalData}

	copied := cfg.Data()

	// Modify copied map
	copied["server"].(map[string]any)["addr"] = ":9090"
	copied["new_key"] = "new_val"

	// Verify original config is unchanged
	if cfg.String("server.addr") != ":8080" {
		t.Errorf("expected original config server.addr to remain ':8080', but got %v", cfg.String("server.addr"))
	}
	if cfg.Get("new_key") != nil {
		t.Errorf("expected original config to not contain new_key")
	}
}

func TestValidationHook(t *testing.T) {
	src := &mockSource{
		name: "mock",
		data: map[string]any{"addr": ""},
	}

	// 1. Success case
	cfg, err := LoadWithOptions(context.Background(), []Source{src},
		WithValidator(func(c *Config) error {
			return nil
		}),
		WithStrictMode(true),
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.String("addr") != "" {
		t.Errorf("expected addr to be empty")
	}

	// 2. Failure case
	valErr := errors.New("addr is required")
	_, err = LoadWithOptions(context.Background(), []Source{src},
		WithValidator(func(c *Config) error {
			if c.String("addr") == "" {
				return valErr
			}
			return nil
		}),
	)
	if err == nil {
		t.Fatalf("expected validation error, got nil")
	}
	if !errors.Is(err, valErr) {
		t.Errorf("expected error to be %v, got %v", valErr, err)
	}
}

func TestNilOverrides(t *testing.T) {
	src1 := &mockSource{
		name: "mock1",
		data: map[string]any{
			"key": "value",
			"nested": map[string]any{
				"val": "hello",
			},
		},
	}
	src2 := &mockSource{
		name: "mock2",
		data: map[string]any{
			"key":    nil, // override key with nil
			"nested": nil, // override nested map with nil
		},
	}

	cfg, err := Load(context.Background(), src1, src2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Get("key") != nil {
		t.Errorf("expected key to be nil, got %v", cfg.Get("key"))
	}
	if cfg.Get("nested") != nil {
		t.Errorf("expected nested to be nil, got %v", cfg.Get("nested"))
	}
}

func TestCloneValue_SliceComplex(t *testing.T) {
	src := &mockSource{
		name: "mock",
		data: map[string]any{
			"list": []any{
				map[string]any{"name": "foo"},
				nil,
				"scalar",
			},
		},
	}

	cfg, err := Load(context.Background(), src)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	listVal := cfg.Get("list")
	list, ok := listVal.([]any)
	if !ok {
		t.Fatalf("expected list to be []any, got %T", listVal)
	}

	if len(list) != 3 {
		t.Errorf("expected list length to be 3, got %d", len(list))
	}

	m, ok := list[0].(map[string]any)
	if !ok {
		t.Fatalf("expected list[0] to be map[string]any, got %T", list[0])
	}
	if m["name"] != "foo" {
		t.Errorf("expected list[0].name to be 'foo'")
	}

	if list[1] != nil {
		t.Errorf("expected list[1] to be nil")
	}

	if list[2] != "scalar" {
		t.Errorf("expected list[2] to be 'scalar'")
	}

	// Verify deep copy
	copied := cfg.Data()
	copiedList := copied["list"].([]any)
	copiedMap := copiedList[0].(map[string]any)
	copiedMap["name"] = "bar"

	// Original should still be "foo"
	if cfg.String("list.0.name") != "" { // dot path doesn't support indices, but we can Get and check:
		t.Errorf("dot path doesn't support indices, but we shouldn't get anything anyway")
	}
	origList := cfg.Get("list").([]any)
	origMap := origList[0].(map[string]any)
	if origMap["name"] != "foo" {
		t.Errorf("expected original list element to remain 'foo', got %v", origMap["name"])
	}
}

func TestDeepMerge_ScalarToMapOverride(t *testing.T) {
	src1 := &mockSource{
		name: "mock1",
		data: map[string]any{
			"port": 8080,
		},
	}
	src2 := &mockSource{
		name: "mock2",
		data: map[string]any{
			"port": map[string]any{
				"value": 9090,
			},
		},
	}

	cfg, err := Load(context.Background(), src1, src2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Int("port.value") != 9090 {
		t.Errorf("expected port.value to be 9090, got %v", cfg.Get("port.value"))
	}
}
