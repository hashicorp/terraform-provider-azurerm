package graph

import (
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/ar"
)

func WaitForCreationReplication(f func() (interface{}, error)) (interface{}, error) {
	return (&resource.StateChangeConf{
		Pending:                   []string{"404", "BadCast"},
		Target:                    []string{"Found"},
		Timeout:                   5 * time.Minute,
		MinTimeout:                1 * time.Second,
		ContinuousTargetOccurence: 10,
		Refresh: func() (interface{}, string, error) {
			i, err := f()
			if err == nil {
				return i, "Found", nil
			}

			r, ok := i.(autorest.Response)
			if !ok {
				return i, "BadCast", nil // sometimes the SDK bubbles up an entirely empty object
			}
			if ar.ResponseWasNotFound(r) {
				return i, "404", nil
			}
			return i, "Error", fmt.Errorf("Error calling f, response was not 404 (%d): %v", r.StatusCode, err)
		},
	}).WaitForState()
}
