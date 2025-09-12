// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/resourcegroups"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceResourceGroup() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceResourceGroupCreate,
		Read:   resourceResourceGroupRead,
		Update: resourceResourceGroupUpdate,
		Delete: resourceResourceGroupDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ResourceGroupID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(90 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(90 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(90 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"tags": commonschema.Tags(),

			"managed_by": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceResourceGroupCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.ResourceGroupsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewResourceGroupID(meta.(*clients.Client).Account.SubscriptionId, d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing resource group: %+v", err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_resource_group", id.ID())
	}

	parameters := resourcegroups.ResourceGroup{
		Location: location.Normalize(d.Get("location").(string)),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v := d.Get("managed_by").(string); v != "" {
		parameters.ManagedBy = pointer.To(v)
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %q: %+v", id, err)
	}

	// TODO: remove this once ARM team confirms the issue is fixed on their end
	//
	// @favoretti: Working around a race condition in ARM eventually consistent backend data storage
	// Sporadically, the ARM api will return successful creation response, following by a 404 to a
	// subsequent `Get()`. Usually, seconds later, the storage is reconciled and following terraform
	// run fails with `RequiresImport`.
	//
	// Snippet from MSFT support:
	// The issue is related to replication of ARM data among regions. For example, another customer
	// has some requests going to East US and other requests to East US 2, and during the time it takes
	// to replicate between the two, they get 404's. The database account is a multi-master account with
	// session consistency - so, write operations will be replicated across regions asynchronously.
	// Session consistency only guarantees read-you-write guarantees within the scope of a session which
	// is either defined by the application (ARM) or by the SDK (in which case the session spans only
	// a single CosmosClient instance) - and given that several of the reads returning 404 after the
	// creation of the resource group were done not only from a different ARM FD machine but even from
	// a different region, they were made outside of the session scope - so, effectively eventually
	// consistent. ARM team has worked in the past to make the multi-master model work transparently,
	// and I assume they will continue this work as will our other teams working on the problem.
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"Waiting"},
		Target:                    []string{"Done"},
		Timeout:                   10 * time.Minute,
		MinTimeout:                4 * time.Second,
		ContinuousTargetOccurence: 3,
		Refresh: func() (interface{}, string, error) {
			rg, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(rg.HttpResponse) {
					return false, "Waiting", nil
				}
				return nil, "Error", fmt.Errorf("retrieving Resource Group: %+v", err)
			}

			return true, "Done", nil
		},
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for Resource Group %s to become available: %+v", id.ResourceGroupName, err)
	} // TODO - Custom Poller?

	d.SetId(id.ID())

	return resourceResourceGroupRead(d, meta)
}

func resourceResourceGroupUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.ResourceGroupsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseResourceGroupID(d.Id())
	if err != nil {
		return err
	}

	patch := resourcegroups.ResourceGroupPatchable{}

	if d.HasChange("managed_by") {
		patch.ManagedBy = pointer.To(d.Get("managed_by").(string))
	}

	if d.HasChange("tags") {
		patch.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if _, err := client.Update(ctx, *id, patch); err != nil {
		return fmt.Errorf("creating %q: %+v", *id, err)
	}

	return resourceResourceGroupRead(d, meta)
}

func resourceResourceGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.ResourceGroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseResourceGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Error reading resource group %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading resource group: %+v", err)
	}

	d.Set("name", id.ResourceGroupName)
	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		d.Set("managed_by", pointer.From(model.ManagedBy))
		tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func resourceResourceGroupDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.ResourceGroupsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseResourceGroupID(d.Id())
	if err != nil {
		return err
	}

	// conditionally check for nested resources and error if they exist
	if meta.(*clients.Client).Features.ResourceGroup.PreventDeletionIfContainsResources {
		// Resource groups sometimes hold on to resource information after the resources have been deleted. We'll retry this check to account for that eventual consistency.
		err = pluginsdk.Retry(10*time.Minute, func() *pluginsdk.RetryError {
			results, err := client.ResourcesListByResourceGroupComplete(ctx, *id, resourcegroups.ResourcesListByResourceGroupOperationOptions{
				Expand: pointer.To("provisioningState"),
				Top:    pointer.To[int64](10),
			})
			if err != nil {
				if response.WasNotFound(results.LatestHttpResponse) {
					return nil
				}
				return pluginsdk.NonRetryableError(fmt.Errorf("listing resources in %s: %v", *id, err))
			}
			nestedResourceIds := make([]string, 0)
			for _, item := range results.Items {
				val := item
				if val.Id != nil {
					nestedResourceIds = append(nestedResourceIds, *val.Id)
				}
			}

			if len(nestedResourceIds) > 0 {
				time.Sleep(30 * time.Second)
				return pluginsdk.RetryableError(resourceGroupContainsItemsError(id.ResourceGroupName, nestedResourceIds))
			}
			return nil
		})

		if err != nil {
			return err
		}
	}

	// This is not the usual pattern for destroys on go-azure-sdk, however, this functionally the same as the resource
	// worked before refactoring, so behaviour has been maintained. This should be investigated in future and brought
	// in-line if possible.
	if resp, err := client.Delete(ctx, *id, resourcegroups.DefaultDeleteOperationOptions()); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	} else {
		if err := resp.Poller.PollUntilDone(ctx); err != nil {
			return fmt.Errorf("polling deleting %s: %+v", *id, err)
		}
	}

	return nil
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
