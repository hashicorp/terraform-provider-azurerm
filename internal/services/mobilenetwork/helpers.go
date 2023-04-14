package mobilenetwork

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// a workaround for that some child resources may still exist for seconds before it fully deleted.
// tracked on https://github.com/Azure/azure-rest-api-specs/issues/22691
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
		MinTimeout:                10 * time.Second,
		PollInterval:              10 * time.Second,
		ContinuousTargetOccurence: 6,
		Timeout:                   time.Until(deadline),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
	}

	return nil
}

// tracked on https://github.com/Azure/azure-rest-api-specs/issues/22634
// some resources defined both systemAssigned and userAssigned Identity type in Swagger but only support userAssigned Identity,
// so add a workaround to convert type here.
func expandMobileNetworkLegacyToUserAssignedIdentity(input []identity.ModelUserAssigned) (*identity.LegacySystemAndUserAssignedMap, error) {
	if len(input) == 0 {
		return nil, nil
	}

	identityValue, err := identity.ExpandUserAssignedMapFromModel(input)
	if err != nil {
		return nil, fmt.Errorf("expanding `identity`: %+v", err)
	}

	output := identity.LegacySystemAndUserAssignedMap{
		Type:        identityValue.Type,
		IdentityIds: identityValue.IdentityIds,
	}

	return &output, nil
}

func flattenMobileNetworkUserAssignedToNetworkLegacyIdentity(input *identity.LegacySystemAndUserAssignedMap) ([]identity.ModelUserAssigned, error) {
	if input == nil {
		return nil, nil
	}

	tmp := identity.UserAssignedMap{
		Type:        input.Type,
		IdentityIds: input.IdentityIds,
	}

	output, err := identity.FlattenUserAssignedMapToModel(&tmp)
	if err != nil {
		return nil, fmt.Errorf("expanding `identity`: %+v", err)
	}

	return *output, nil
}
