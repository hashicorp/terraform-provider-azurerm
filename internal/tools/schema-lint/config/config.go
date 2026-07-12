// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

// Package config models the optional .schema-lint.json configuration file used
// to enable, disable and re-configure lint rules (in the spirit of a
// markdownlint config), and to scope which resources are linted.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// DefaultFileName is the config file looked up in the working directory when no
// explicit path is provided.
const DefaultFileName = ".schema-lint.json"

// RuleConfig holds per-rule configuration overrides.
type RuleConfig struct {
	// Enabled overrides whether the rule runs. When nil the Default (or the
	// linter default of enabled) applies.
	Enabled *bool `json:"enabled,omitempty"`
	// Severity overrides the rule's default severity ("error" or "warning").
	Severity string `json:"severity,omitempty"`
	// Options holds arbitrary per-rule options for rules that support them.
	Options map[string]interface{} `json:"options,omitempty"`
}

// Config models a .schema-lint.json file.
type Config struct {
	// Default controls whether rules are enabled unless individually
	// overridden. Defaults to true when omitted.
	Default *bool `json:"default,omitempty"`
	// Rules holds per-rule overrides keyed by rule ID (e.g. "SL001").
	Rules map[string]RuleConfig `json:"rules,omitempty"`
	// IncludeResources, when non-empty, restricts linting to these
	// resource/data source type names.
	IncludeResources []string `json:"includeResources,omitempty"`
	// SkipResources excludes these resource/data source type names.
	SkipResources []string `json:"skipResources,omitempty"`
}

// Load reads a config from path. When path is empty it looks for
// DefaultFileName in the working directory; a missing default file is not an
// error and yields an empty (all-defaults) config.
func Load(path string) (*Config, error) {
	explicit := path != ""
	if path == "" {
		path = DefaultFileName
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) && !explicit {
			return &Config{}, nil
		}
		return nil, fmt.Errorf("reading config %q: %w", path, err)
	}

	cfg := &Config{}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config %q: %w", path, err)
	}

	return cfg, nil
}

// DefaultEnabled reports whether rules are enabled by default.
func (c *Config) DefaultEnabled() bool {
	return c.Default == nil || *c.Default
}

// Rule returns the override for the given rule ID (case-insensitive) and whether
// one was configured.
func (c *Config) Rule(id string) (RuleConfig, bool) {
	for k, v := range c.Rules {
		if strings.EqualFold(k, id) {
			return v, true
		}
	}

	return RuleConfig{}, false
}

// ResourceAllowed reports whether a resource/data source type name passes the
// include/skip filters.
func (c *Config) ResourceAllowed(name string) bool {
	if len(c.IncludeResources) > 0 && !containsFold(c.IncludeResources, name) {
		return false
	}

	return !containsFold(c.SkipResources, name)
}

func containsFold(list []string, name string) bool {
	for _, v := range list {
		if strings.EqualFold(strings.TrimSpace(v), name) {
			return true
		}
	}

	return false
}
