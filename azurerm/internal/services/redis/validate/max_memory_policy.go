package validate

import (
	"fmt"
	"strings"
)

func MaxMemoryPolicy(v interface{}, _ string) (warnings []string, errors []error) {
	// TODO: in time this can be replaced with a generic function, moving for now
	value := strings.ToLower(v.(string))
	families := map[string]bool{
		"noeviction":      true,
		"allkeys-lru":     true,
		"volatile-lru":    true,
		"allkeys-random":  true,
		"volatile-random": true,
		"volatile-ttl":    true,
		"allkeys-lfu":     true,
		"volatile-lfu":    true,
	}

	if !families[value] {
		errors = append(errors, fmt.Errorf("Redis Max Memory Policy can only be 'noeviction' / 'allkeys-lru' / 'volatile-lru' / 'allkeys-random' / 'volatile-random' / 'volatile-ttl' / 'allkeys-lfu' / 'volatile-lfu'"))
	}

	return warnings, errors
}
