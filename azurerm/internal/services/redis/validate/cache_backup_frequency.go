package validate

import "fmt"

func CacheBackupFrequency(v interface{}, _ string) (warnings []string, errors []error) {
	value := v.(int)
	families := map[int]bool{
		15:   true,
		30:   true,
		60:   true,
		360:  true,
		720:  true,
		1440: true,
	}

	if !families[value] {
		errors = append(errors, fmt.Errorf("Redis Backup Frequency can only be '15', '30', '60', '360', '720' or '1440'"))
	}

	return warnings, errors
}
