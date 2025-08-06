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

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-06-01/resources" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceResourceGroup() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceResourceGroupCreateUpdate,
		Read:   resourceResourceGroupRead,
		Update: resourceResourceGroupCreateUpdate,
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

			"tags": tags.Schema(),

			"managed_by": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceResourceGroupCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.GroupsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	location := location.Normalize(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	if d.IsNewResource() {
		existing, err := client.Get(ctx, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing resource group: %+v", err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_resource_group", *existing.ID)
		}
	}

	parameters := resources.Group{
		Location: utils.String(location),
		Tags:     tags.Expand(t),
	}

	if v := d.Get("managed_by").(string); v != "" {
		parameters.ManagedBy = pointer.To(v)
	}

	if _, err := client.CreateOrUpdate(ctx, name, parameters); err != nil {
		return fmt.Errorf("creating Resource Group %q: %+v", name, err)
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
	if d.IsNewResource() {
		stateConf := &pluginsdk.StateChangeConf{ //nolint:staticcheck
			Pending:                   []string{"Waiting"},
			Target:                    []string{"Done"},
			Timeout:                   10 * time.Minute,
			MinTimeout:                4 * time.Second,
			ContinuousTargetOccurence: 3,
			Refresh: func() (interface{}, string, error) {
				rg, err := client.Get(ctx, name)
				if err != nil {
					if utils.ResponseWasNotFound(rg.Response) {
						return false, "Waiting", nil
					}
					return nil, "Error", fmt.Errorf("retrieving Resource Group: %+v", err)
				}

				return true, "Done", nil
			},
		}

		if _, err := stateConf.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("waiting for Resource Group %s to become available: %+v", name, err)
		}
	}

	resp, err := client.Get(ctx, name)
	if err != nil {
		return fmt.Errorf("retrieving Resource Group %q: %+v", name, err)
	}

	// @tombuildsstuff: intentionally leaving this for now, since this'll need
	// details in the upgrade notes given how the Resource Group ID is cased incorrectly
	// but needs to be fixed (resourcegroups -> resourceGroups)
	id, err := parse.ResourceGroupIDInsensitively(*resp.ID)
	if err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceResourceGroupRead(d, meta)
}

func resourceResourceGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.GroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ResourceGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading resource group %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading resource group: %+v", err)
	}

	d.Set("name", resp.Name)
	d.Set("location", location.NormalizeNilable(resp.Location))
	d.Set("managed_by", pointer.From(resp.ManagedBy))
	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceResourceGroupDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.GroupsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ResourceGroupID(d.Id())
	if err != nil {
		return err
	}

	// conditionally check for nested resources and error if they exist
	if meta.(*clients.Client).Features.ResourceGroup.PreventDeletionIfContainsResources {
		resourceClient := meta.(*clients.Client).Resource.LegacyResourcesClient
		// Resource groups sometimes hold on to resource information after the resources have been deleted. We'll retry this check to account for that eventual consistency.
		err = pluginsdk.Retry(10*time.Minute, func() *pluginsdk.RetryError {
			results, err := resourceClient.ListByResourceGroup(ctx, id.ResourceGroup, "", "provisioningState", utils.Int32(500))
			if err != nil {
				if response.WasNotFound(results.Response().Response.Response) {
					return nil
				}
				return pluginsdk.NonRetryableError(fmt.Errorf("listing resources in %s: %v", *id, err))
			}
			nestedResourceIds := make([]string, 0)
			for _, value := range results.Values() {
				val := value
				if val.ID != nil {
					nestedResourceIds = append(nestedResourceIds, *val.ID)
				}

				if err := results.NextWithContext(ctx); err != nil {
					return pluginsdk.NonRetryableError(fmt.Errorf("retrieving next page of nested items for %s: %+v", id, err))
				}
			}

			if len(nestedResourceIds) > 0 {
				time.Sleep(30 * time.Second)
				return pluginsdk.RetryableError(resourceGroupContainsItemsError(id.ResourceGroup, nestedResourceIds))
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	deleteFuture, err := client.Delete(ctx, id.ResourceGroup, "")
	if err != nil {
		if response.WasNotFound(deleteFuture.Response()) {
			return nil
		}
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	err = deleteFuture.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
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
