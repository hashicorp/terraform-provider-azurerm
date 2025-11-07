package validate

import "errors"

func ValidateUniqueInlineEventTypeNames(input []interface{}) error {
	if len(input) <= 1 {
		return nil
	}

	seen := make(map[string]bool)
	for _, item := range input {
		eventType := item.(map[string]interface{})
		name := eventType["name"].(string)

		if seen[name] {
			return errors.New("invalid value for `partner_topic.0.event_type_definitions.0.inline_event_type` - `name` must be unique")
		}
		seen[name] = true
	}
	return nil
}
