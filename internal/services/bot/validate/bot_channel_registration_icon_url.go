package validate

import (
	"fmt"
	"strings"
)

func BotChannelRegistrationIconUrl(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if !strings.HasSuffix(v, ".png") {
		errors = append(errors, fmt.Errorf("only png is supported"))
		return
	}

	return
}
