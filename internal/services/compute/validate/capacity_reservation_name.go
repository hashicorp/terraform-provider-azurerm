package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func CapacityReservationName() func(i interface{}, k string) (warnings []string, errors []error) {
	return validation.StringMatch(regexp.MustCompile(`^[^_\W]([\w-._]{0,78}[\w_])?$`), `The Capacity Reservation Name must be between 1 and 80 characters long. It cannot contain special characters \/"[]:|<>+=;,?*@&, whitespace, or begin with '_' or end with '.' or '-'`)
}
