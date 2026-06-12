package custompollers

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/resourcegroups"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

type resourceGroupPreventDeletePoller struct {
	client *resourcegroups.ResourceGroupsClient
	id     commonids.ResourceGroupId
}

var _ pollers.PollerType = &resourceGroupPreventDeletePoller{}

func NewResourceGroupPreventDeletePoller(client *resourcegroups.ResourceGroupsClient, id commonids.ResourceGroupId) *resourceGroupPreventDeletePoller {
	return &resourceGroupPreventDeletePoller{
		client: client,
		id:     id,
	}
}

func (p resourceGroupPreventDeletePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	results, err := p.client.ResourcesListByResourceGroupComplete(ctx, p.id, resourcegroups.ResourcesListByResourceGroupOperationOptions{
		Expand: pointer.To("provisioningState"),
		Top:    pointer.To[int64](10),
	})
	if err != nil {
		if response.WasNotFound(results.LatestHttpResponse) {
			return pollingSuccess, nil
		}
		return pollingFailed, fmt.Errorf("listing resources in %s: %v", p.id, err)
	}

	nestedResourceIds := make([]string, 0)
	for _, item := range results.Items {
		val := item
		if val.Id != nil {
			nestedResourceIds = append(nestedResourceIds, *val.Id)
		}
	}

	if len(nestedResourceIds) > 0 {
		return &pollers.PollResult{
			PollInterval: 30 * time.Second,
			Status:       pollers.PollingStatusInProgress,
		}, resourceGroupContainsItemsError(p.id.ResourceGroupName, nestedResourceIds)
	}

	return pollingSuccess, nil
}

func resourceGroupContainsItemsError(name string, nestedResourceIds []string) error {
	formattedResourceUris := make([]string, 0)
	for _, id := range nestedResourceIds {
		formattedResourceUris = append(formattedResourceUris, fmt.Sprintf("* `%s`", id))
	}
	sort.Strings(formattedResourceUris)

	message := fmt.Sprintf(`deleting Resource Group %[1]q: the Resource Group still contains Resources.

Terraform is configured to check for Resources within the Resource Group when deleting the Resource Group - and
raise an error if nested Resources still exist to avoid unintentionally deleting these Resources.

Terraform has detected that the following Resources still exist within the Resource Group:

%[2]s

This feature is intended to avoid the unintentional destruction of nested Resources provisioned through some
other means (for example, an ARM Template Deployment) - as such you must either remove these Resources, or
disable this behaviour using the feature flag 'prevent_deletion_if_contains_resources' within the 'features'
block when configuring the Provider, for example:

provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

When that feature flag is set, Terraform will skip checking for any Resources within the Resource Group and
delete this using the Azure API directly (which will clear up any nested resources).

More information on the 'features' block can be found in the documentation:
https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/features-block
`, name, strings.Join(formattedResourceUris, "\n"))
	return errors.New(strings.ReplaceAll(message, "'", "`"))
}
