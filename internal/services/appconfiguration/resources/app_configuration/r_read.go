// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package app_configuration

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2024-05-01/configurationstores"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2024-05-01/deletedconfigurationstores"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2024-05-01/operations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2024-05-01/replicas"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAppConfigurationRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppConfiguration.ConfigurationStoresClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := configurationstores.ParseConfigurationStoreID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	resultPage, err := client.ListKeysComplete(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving access keys for %s: %+v", *id, err)
	}

	d.Set("name", id.ConfigurationStoreName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		d.Set("sku", model.Sku.Name)

		if props := model.Properties; props != nil {
			if dataPlaneProxy := props.DataPlaneProxy; dataPlaneProxy != nil {
				d.Set("data_plane_proxy_authentication_mode", pointer.FromEnum(dataPlaneProxy.AuthenticationMode))
				d.Set("data_plane_proxy_private_link_delegation_enabled", pointer.From(dataPlaneProxy.PrivateLinkDelegation) == configurationstores.PrivateLinkDelegationEnabled)
			}

			d.Set("endpoint", props.Endpoint)
			d.Set("encryption", helpers.FlattenAppConfigurationEncryption(props.Encryption))
			d.Set("public_network_access", pointer.FromEnum(props.PublicNetworkAccess))

			localAuthEnabled := true
			if props.DisableLocalAuth != nil {
				localAuthEnabled = !(*props.DisableLocalAuth)
			}

			d.Set("local_auth_enabled", localAuthEnabled)

			d.Set("purge_protection_enabled", pointer.From(props.EnablePurgeProtection))

			softDeleteRetentionDays := 0
			if props.SoftDeleteRetentionInDays != nil {
				softDeleteRetentionDays = int(*props.SoftDeleteRetentionInDays)
			}
			d.Set("soft_delete_retention_days", softDeleteRetentionDays)
		}

		accessKeys := flattenAppConfigurationAccessKeys(resultPage.Items)
		d.Set("primary_read_key", accessKeys.primaryReadKey)
		d.Set("primary_write_key", accessKeys.primaryWriteKey)
		d.Set("secondary_read_key", accessKeys.secondaryReadKey)
		d.Set("secondary_write_key", accessKeys.secondaryWriteKey)

		flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", flattenedIdentity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		replicasClient := meta.(*clients.Client).AppConfiguration.ReplicasClient
		resp, err := replicasClient.ListByConfigurationStoreComplete(ctx, replicas.NewConfigurationStoreID(id.SubscriptionId, id.ResourceGroupName, id.ConfigurationStoreName))
		if err != nil {
			return fmt.Errorf("retrieving replicas for %s: %+v", *id, err)
		}

		replica, err := helpers.FlattenAppConfigurationReplicas(resp.Items)
		if err != nil {
			return fmt.Errorf("flattening replicas for %s: %+v", *id, err)
		}
		d.Set("replica", replica)

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func flattenAppConfigurationAccessKeys(values []configurationstores.ApiKey) flattenedAccessKeys {
	result := flattenedAccessKeys{
		primaryReadKey:    make([]any, 0),
		primaryWriteKey:   make([]any, 0),
		secondaryReadKey:  make([]any, 0),
		secondaryWriteKey: make([]any, 0),
	}

	for _, value := range values {
		if value.Name == nil || value.ReadOnly == nil {
			continue
		}

		accessKey := flattenAppConfigurationAccessKey(value)
		name := *value.Name
		readOnly := *value.ReadOnly

		if strings.HasPrefix(strings.ToLower(name), "primary") {
			if readOnly {
				result.primaryReadKey = accessKey
			} else {
				result.primaryWriteKey = accessKey
			}
		}

		if strings.HasPrefix(strings.ToLower(name), "secondary") {
			if readOnly {
				result.secondaryReadKey = accessKey
			} else {
				result.secondaryWriteKey = accessKey
			}
		}
	}

	return result
}

func flattenAppConfigurationAccessKey(input configurationstores.ApiKey) []any {
	return []any{
		map[string]any{
			"connection_string": pointer.From(input.ConnectionString),
			"id":                pointer.From(input.Id),
			"secret":            pointer.From(input.Value),
		},
	}
}

var _ pollers.PollerType = &purgeDeletedPoller{}

type purgeDeletedPoller struct {
	client *deletedconfigurationstores.DeletedConfigurationStoresClient
	id     deletedconfigurationstores.DeletedConfigurationStoreId
}

func (p *purgeDeletedPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.ConfigurationStoresGetDeleted(ctx, p.id)

	status := "dropped connection"
	if resp.HttpResponse != nil {
		status = fmt.Sprintf("%d", resp.HttpResponse.StatusCode)
	}

	if response.WasNotFound(resp.HttpResponse) {
		return &pollers.PollResult{
			HttpResponse: &client.Response{
				Response: resp.HttpResponse,
			},
			Status: pollers.PollingStatusSucceeded,
		}, nil
	}

	if response.WasStatusCode(resp.HttpResponse, http.StatusOK) {
		return &pollers.PollResult{
			HttpResponse: &client.Response{
				Response: resp.HttpResponse,
			},
			Status: pollers.PollingStatusInProgress,
		}, nil
	}

	return nil, fmt.Errorf("unexpected status %q: %+v", status, err)
}

type flattenedAccessKeys struct {
	primaryReadKey    []any
	primaryWriteKey   []any
	secondaryReadKey  []any
	secondaryWriteKey []any
}

func parsePublicNetworkAccess(input string) *configurationstores.PublicNetworkAccess {
	vals := map[string]configurationstores.PublicNetworkAccess{
		"disabled": configurationstores.PublicNetworkAccessDisabled,
		"enabled":  configurationstores.PublicNetworkAccessEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v
	}

	// otherwise presume it's an undefined value and best-effort it
	return pointer.ToEnum[configurationstores.PublicNetworkAccess](input)
}

func userIsMissingNecessaryPermission(name, location string) string {
	return fmt.Sprintf(`
An existing soft-deleted App Configuration exists with the Name %q in the location %q, however
the credentials Terraform is using has insufficient permissions to check for an existing soft-deleted App Configuration.
You can opt out of this behaviour by using the "features" block (located within the "provider" block) - more information
can be found here:
https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/features-block
`, name, location)
}

func resourceConfigurationStoreWaitForNameAvailable(ctx context.Context, client *operations.OperationsClient, configurationStoreId configurationstores.ConfigurationStoreId) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal error: context had no deadline")
	}
	state := &pluginsdk.StateChangeConf{
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 2,
		Pending:                   []string{"Unavailable"},
		Target:                    []string{"Available"},
		Refresh:                   resourceConfigurationStoreNameAvailabilityRefreshFunc(ctx, client, configurationStoreId),
		Timeout:                   time.Until(deadline),
	}

	if _, err := state.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for the name from %s to become available: %+v", configurationStoreId, err)
	}

	return nil
}

func resourceConfigurationStoreNameAvailabilityRefreshFunc(ctx context.Context, client *operations.OperationsClient, configurationStoreId configurationstores.ConfigurationStoreId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		log.Printf("[DEBUG] Checking to see if the name for %s is available ..", configurationStoreId)

		subscriptionId := commonids.NewSubscriptionID(configurationStoreId.SubscriptionId)

		parameters := operations.CheckNameAvailabilityParameters{
			Name: configurationStoreId.ConfigurationStoreName,
			Type: operations.ConfigurationResourceTypeMicrosoftPointAppConfigurationConfigurationStores,
		}

		resp, err := client.CheckNameAvailability(ctx, subscriptionId, parameters)
		if err != nil {
			return resp, "Error", fmt.Errorf("retrieving Deployment: %+v", err)
		}

		if resp.Model == nil {
			return resp, "Error", fmt.Errorf("unexpected null model of %s", configurationStoreId)
		}

		if resp.Model.NameAvailable == nil {
			return resp, "Error", fmt.Errorf("unexpected null NameAvailable property of %s", configurationStoreId)
		}

		if !*resp.Model.NameAvailable {
			return resp, "Unavailable", nil
		}
		return resp, "Available", nil
	}
}

func deleteReplicas(ctx context.Context, replicaClient *replicas.ReplicasClient, operationClient *operations.OperationsClient, configurationStoreReplicaIds []replicas.ReplicaId) error {
	for _, configurationStoreReplicaId := range configurationStoreReplicaIds {
		if err := replicaClient.DeleteThenPoll(ctx, configurationStoreReplicaId); err != nil {
			return fmt.Errorf("deleting replica %q: %+v", configurationStoreReplicaId, err)
		}
	}

	for _, configurationStoreReplicaId := range configurationStoreReplicaIds {
		if err := resourceConfigurationStoreReplicaWaitForNameAvailable(ctx, operationClient, configurationStoreReplicaId); err != nil {
			return fmt.Errorf("waiting for replica %q name to be released: %+v", configurationStoreReplicaId, err)
		}
	}

	return nil
}

func resourceConfigurationStoreReplicaWaitForNameAvailable(ctx context.Context, client *operations.OperationsClient, configurationStoreReplicaId replicas.ReplicaId) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal error: context had no deadline")
	}
	state := &pluginsdk.StateChangeConf{
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 2,
		Pending:                   []string{"Unavailable"},
		Target:                    []string{"Available"},
		Refresh:                   resourceConfigurationStoreReplicaNameAvailabilityRefreshFunc(ctx, client, configurationStoreReplicaId),
		Timeout:                   time.Until(deadline),
	}

	if _, err := state.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for the name from %s to become available: %+v", configurationStoreReplicaId, err)
	}

	return nil
}

func resourceConfigurationStoreReplicaNameAvailabilityRefreshFunc(ctx context.Context, client *operations.OperationsClient, configurationStoreReplicaId replicas.ReplicaId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		log.Printf("[DEBUG] Checking to see if the name for %s is available ..", configurationStoreReplicaId)

		subscriptionId := commonids.NewSubscriptionID(configurationStoreReplicaId.SubscriptionId)

		parameters := operations.CheckNameAvailabilityParameters{
			Name: fmt.Sprintf("%s-%s", configurationStoreReplicaId.ConfigurationStoreName, configurationStoreReplicaId.ReplicaName),
			Type: operations.ConfigurationResourceTypeMicrosoftPointAppConfigurationConfigurationStores,
		}

		resp, err := client.CheckNameAvailability(ctx, subscriptionId, parameters)
		if err != nil {
			return resp, "Error", fmt.Errorf("retrieving Deployment: %+v", err)
		}

		if resp.Model == nil {
			return resp, "Error", fmt.Errorf("unexpected null model of %s", configurationStoreReplicaId)
		}

		if resp.Model.NameAvailable == nil {
			return resp, "Error", fmt.Errorf("unexpected null NameAvailable property of %s", configurationStoreReplicaId)
		}

		if !*resp.Model.NameAvailable {
			return resp, "Unavailable", nil
		}
		return resp, "Available", nil
	}
}
