package vcr

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

var sidecarMu sync.Mutex

func updateDynamicSidecar(testName string, newTargets map[string]string) {
	if len(newTargets) == 0 {
		return
	}

	sidecarMu.Lock()
	defer sidecarMu.Unlock()

	cachePath := filepath.Join(testDataPath, testName+"_dynamic.json")
	cache := make(map[string]string)
	if data, err := os.ReadFile(cachePath); err == nil {
		_ = json.Unmarshal(data, &cache)
	}

	for k, v := range newTargets {
		cache[k] = v
	}

	if data, err := json.MarshalIndent(cache, "", "  "); err == nil {
		_ = os.WriteFile(cachePath, data, 0644)
	}
}

func readDynamicSidecar(testName string) map[string]string {
	sidecarMu.Lock()
	defer sidecarMu.Unlock()

	cachePath := filepath.Join(testDataPath, testName+"_dynamic.json")
	cache := make(map[string]string)
	if data, err := os.ReadFile(cachePath); err == nil {
		_ = json.Unmarshal(data, &cache)
	}
	return cache
}

// ScrubDynamicValues takes a string and a map of {Placeholder: TrueValue} and applies global string substitutions
func ScrubDynamicValues(s string, seeds map[string]string) string {
	for placeholder, trueVal := range seeds {
		if trueVal != "" {
			re := regexp.MustCompile("(?i)" + regexp.QuoteMeta(trueVal))
			s = re.ReplaceAllString(s, placeholder)
		}
	}
	return s
}
