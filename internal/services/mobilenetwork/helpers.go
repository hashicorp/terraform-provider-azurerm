package mobilenetwork

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// a workaround for that some child resources may still exist for seconds before it fully deleted.
// it will cause the error "Can not delete resource before nested resources are deleted."
func resourceMobileNetworkChildWaitForDeletion(ctx context.Context, id string, getFunction func() (*http.Response, error)) error {
	deadline, _ := ctx.Deadline()
	stateConf := &pluginsdk.StateChangeConf{
		Pending: []string{"200"},
		Target:  []string{"404"},
		Refresh: func() (interface{}, string, error) {
			resp, err := getFunction()
			if err != nil {
				if response.WasNotFound(resp) {
					return resp, strconv.Itoa(resp.StatusCode), nil
				}

				return nil, strconv.Itoa(resp.StatusCode), fmt.Errorf("polling for %s: %+v", id, err)
			}

			return resp, strconv.Itoa(resp.StatusCode), nil
		},
		MinTimeout:   10 * time.Second,
		PollInterval: 5 * time.Second,
		Timeout:      time.Until(deadline),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for deleting %s: %+v", id, err)
	}

	return nil
}
