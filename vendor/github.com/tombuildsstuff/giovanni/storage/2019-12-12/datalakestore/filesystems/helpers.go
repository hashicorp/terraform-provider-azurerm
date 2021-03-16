package filesystems

import (
	"fmt"
	"strings"
)

func buildProperties(input map[string]string) string {
	// properties has to be a comma-separated key-value pair
	properties := make([]string, 0)

	for k, v := range input {
		properties = append(properties, fmt.Sprintf("%s=%s", k, v))
	}

	return strings.Join(properties, ",")
}

func parseProperties(input string) (*map[string]string, error) {
	properties := make(map[string]string)
	if input == "" {
		return &properties, nil
	}

	// properties is a comma-separated list of key-value pairs
	splitProperties := strings.Split(input, ",")
	for _, propertyRaw := range splitProperties {
		// because these are base64-encoded they're likely to end in at least one =
		// as such we can't string split on that -_-
		position := strings.Index(propertyRaw, "=")
		if position < 0 {
			return nil, fmt.Errorf("Expected there to be an equals in the key value pair: %q", propertyRaw)
		}

		key := propertyRaw[0:position]
		value := propertyRaw[position+1:]
		properties[key] = value
	}
	return &properties, nil
}
