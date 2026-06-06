package file

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/brian-nunez/bconfig"
	"gopkg.in/yaml.v3"
)

type fileSource struct {
	path string
}

// Source returns a bconfig.Source that loads configuration from the file at path.
func Source(path string) bconfig.Source {
	return &fileSource{path: path}
}

// Name returns the name of this source in the format "file:<path>".
func (f *fileSource) Name() string {
	return "file:" + f.path
}

// Load reads and parses the configuration file.
func (f *fileSource) Load(ctx context.Context) (map[string]any, error) {
	data, err := os.ReadFile(f.path)
	if err != nil {
		return nil, err
	}

	if len(strings.TrimSpace(string(data))) == 0 {
		return make(map[string]any), nil
	}

	ext := strings.ToLower(filepath.Ext(f.path))
	var result map[string]any

	switch ext {
	case ".json":
		if err := json.Unmarshal(data, &result); err != nil {
			return nil, err
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &result); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported config file extension %q", ext)
	}

	if result == nil {
		result = make(map[string]any)
	}

	return result, nil
}
