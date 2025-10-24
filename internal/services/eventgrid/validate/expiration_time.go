package validate

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func ExpirationTimeIfNotActivated() pluginsdk.SchemaValidateFunc {
	return func(i interface{}, k string) (warnings []string, errors []error) {
		v, ok := i.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		selectedTime, _ := time.Parse(time.RFC3339, v)
		timeUntilExpiry := selectedTime.Sub(time.Now().In(time.UTC))
		if timeUntilExpiry < 0 || timeUntilExpiry > 7*24*time.Hour {
			errors = append(errors, fmt.Errorf("`expiration_time_if_not_activated_in_utc` must be within 7 days from now"))
		}

		return
	}
}
