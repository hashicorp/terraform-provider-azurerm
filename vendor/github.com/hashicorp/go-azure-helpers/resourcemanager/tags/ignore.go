// Copyright IBM Corp. 2018, 2025
// SPDX-License-Identifier: MPL-2.0

package tags

import "strings"

// IgnoreConfig describes the set of tag keys that the provider should ignore on
// every resource and data source. A key is ignored when it exactly matches an
// entry in Keys (case-sensitive) or when it begins with an entry in KeyPrefixes.
type IgnoreConfig struct {
	Keys        []string
	KeyPrefixes []string
}

// ignore is the process-wide IgnoreConfig set once when the provider is
// configured. A nil value means no tags are ignored, preserving the default
// behavior for anyone who does not configure the `ignore_tags` block.
var ignore *IgnoreConfig

// SetIgnore stores the provider-level IgnoreConfig so that the tag helpers in
// this package (and the framework helpers that read Ignore()) apply it.
func SetIgnore(config *IgnoreConfig) {
	ignore = config
}

// Ignore returns the process-wide IgnoreConfig, or nil if none is configured.
func Ignore() *IgnoreConfig {
	return ignore
}

// ignored reports whether the tag key should be ignored under this config.
func (c *IgnoreConfig) ignored(key string) bool {
	if c == nil {
		return false
	}

	for _, k := range c.Keys {
		if key == k {
			return true
		}
	}

	for _, prefix := range c.KeyPrefixes {
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}

	return false
}

// ApplyPtrMap returns a copy of input with every ignored key removed. A nil
// receiver or a nil input is returned unchanged.
func (c *IgnoreConfig) ApplyPtrMap(input *map[string]string) *map[string]string {
	if c == nil || input == nil {
		return input
	}

	output := make(map[string]string, len(*input))
	for k, v := range *input {
		if c.ignored(k) {
			continue
		}
		output[k] = v
	}

	return &output
}

// ApplyMap returns a copy of input with every ignored key removed. A nil
// receiver is returned unchanged.
func (c *IgnoreConfig) ApplyMap(input map[string]*string) map[string]*string {
	if c == nil {
		return input
	}

	output := make(map[string]*string, len(input))
	for k, v := range input {
		if c.ignored(k) {
			continue
		}
		output[k] = v
	}

	return output
}
