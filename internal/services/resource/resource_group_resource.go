package resource

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-06-01/resources"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
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
			"name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"tags": tags.Schema(),
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

	if _, err := client.CreateOrUpdate(ctx, name, parameters); err != nil {
		return fmt.Errorf("creating Resource Group %q: %+v", name, err)
	}

	resp, err := client.Get(ctx, name)
	if err != nil {
		return fmt.Errorf("retrieving Resource Group %q: %+v", name, err)
	}

	// @tombuildsstuff: intentionally leaving this for now, since this'll need
	// details in the upgrade notes given how the Resource Group ID is cased incorrectly
	// but needs to be fixed (resourcegroups -> resourceGroups)
	d.SetId(*resp.ID)

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
		resourceClient := meta.(*clients.Client).Resource.ResourcesClient
		// Resource groups sometimes hold on to resource information after the resources have been deleted.
		// Check if the resources has been deleted and only exist in cache of the resource group, or there actually exist undeleted resources.
		results, err := resourceClient.ListByResourceGroupComplete(ctx, id.ResourceGroup, "", "provisioningState", utils.Int32(500))
		if err != nil {
			return fmt.Errorf("listing resources in %s: %v", *id, err)
		}

		existResourceIds := make([]string, 0)             // GET resource returned http 200, resource does exist.
		nonExistResourceIds := make([]string, 0)          // GET resource returned http 404, resource does not exist.
		unknownApiVersionResourceIds := make([]string, 0) // resources unable to get API version or GET returned error except 200/404
		providersClient := meta.(*clients.Client).Resource.ResourceProvidersClient

		for results.NotDone() {
			val := results.Value()
			if val.ID != nil {
				if resourceApiVersion, err := getResourceTypeApiVersion(ctx, providersClient, val); err != nil {
					log.Printf("[WARN] Get Resource Type Api version failed: %+v.", err)
					unknownApiVersionResourceIds = append(unknownApiVersionResourceIds, *val.ID)
				} else if exist, err := resourceExistById(ctx, resourceClient, *val.ID, resourceApiVersion); err != nil {
					log.Printf("[WARN] Get Resource (%s) failed: %+v.", *val.ID, err)
					unknownApiVersionResourceIds = append(unknownApiVersionResourceIds, *val.ID)
				} else {
					if exist {
						existResourceIds = append(existResourceIds, *val.ID)
					} else {
						nonExistResourceIds = append(nonExistResourceIds, *val.ID)
					}
				}
			}

			if err := results.NextWithContext(ctx); err != nil {
				return fmt.Errorf("retrieving next page of nested items for %s: %+v", id, err)
			}
		}

		if len(existResourceIds) > 0 {
			return resourceGroupContainsItemsError(id.ResourceGroup, append(existResourceIds, unknownApiVersionResourceIds...))
		}

		if len(unknownApiVersionResourceIds) > 0 {
			err = pluginsdk.Retry(10*time.Minute, func() *pluginsdk.RetryError {
				results, err := resourceClient.ListByResourceGroupComplete(ctx, id.ResourceGroup, "", "provisioningState", utils.Int32(500))
				if err != nil {
					return pluginsdk.NonRetryableError(fmt.Errorf("listing resources in %s: %v", *id, err))
				}
				nestedResourceIds := make([]string, 0)
				for results.NotDone() {
					val := results.Value()
					if val.ID != nil && !utils.SliceContainsValue(nonExistResourceIds, *val.ID) {
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
	}

	deleteFuture, err := client.Delete(ctx, id.ResourceGroup, "")
	if err != nil {
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
https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs#features
`, name, strings.Join(formattedResourceUris, "\n"))
	return fmt.Errorf(strings.ReplaceAll(message, "'", "`"))
}

func resourceExistById(ctx context.Context, resourceClient *resources.Client, resourceId string, apiVersion string) (exist bool, err error) {
	resource, err := resourceClient.GetByID(ctx, resourceId, apiVersion)
	if err != nil {
		if utils.ResponseWasNotFound(resource.Response) {
			return false, nil
		}
		return true, err
	}

	return resource.ID != nil, nil
}

func getResourceTypeApiVersion(ctx context.Context, client *resources.ProvidersClient, resource resources.GenericResourceExpanded) (apiVersion string, err error) {
	providerName, targetResourceType, err := splitProviderAndResourceTypes(*resource.Type)
	if err != nil {
		return "", err
	}

	provider, err := client.Get(ctx, providerName, fmt.Sprintf("resourceTypes/%s", targetResourceType))
	if err != nil {
		return "", err
	}

	if provider.ResourceTypes == nil {
		return "", fmt.Errorf("getting resource(%s/%s) default api version, the provider returned nil on resource type. ", providerName, targetResourceType)
	}

	for _, resourceType := range *provider.ResourceTypes {
		if resourceType.ResourceType != nil && *resourceType.ResourceType == targetResourceType {
			if resourceType.DefaultAPIVersion != nil && *resourceType.DefaultAPIVersion != "" {
				return *resourceType.DefaultAPIVersion, nil
			} else if len(*resourceType.APIVersions) > 0 && (*resourceType.APIVersions)[0] != "" {
				return (*resourceType.APIVersions)[0], nil
			}

			return "", fmt.Errorf("getting resource(%s/%s) api version, there is no available api version. ", providerName, targetResourceType)
		}
	}

	return "", fmt.Errorf("getting resource(%s/%s) api version, there is no match resource type. ", providerName, targetResourceType)
}

func splitProviderAndResourceTypes(fullResourceType string) (provider string, resourceType string, err error) {
	p, r, found := strings.Cut(fullResourceType, "/")
	if !found {
		return "", "", fmt.Errorf("spliting resourceType %s failed", provider)
	}
	return p, r, nil
}
