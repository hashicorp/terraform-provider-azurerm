// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package private_endpoint

import (
	"errors"
	"fmt"
	"log"
	"slices"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/privateendpoints"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePrivateEndpointCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.PrivateEndpoints
	dnsClient := meta.(*clients.Client).Network.PrivateDnsZoneGroups
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := privateendpoints.NewPrivateEndpointID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id, privateendpoints.DefaultGetOperationOptions())
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if existing.Model != nil {
			return tf.ImportAsExistsError("azurerm_private_endpoint", id.ID())
		}
	}

	privateDnsZoneGroup := d.Get("private_dns_zone_group").([]any)

	// Certain child resources like those for cognitive account lock the parent resource to make sure we don't try and update the parent when it's not ready.
	// Due to that, we'll lock on that resource id to try and prevent those type of errors.
	privateLinkServiceConnections, privateLinkServiceConnectionIds := expandPrivateLinkEndpointServiceConnection(d.Get("private_service_connection").([]any), false)
	manualPrivateLinkServiceConnections, manualPrivateLinkServiceConnectionIds := expandPrivateLinkEndpointServiceConnection(d.Get("private_service_connection").([]any), true)
	privateLinkServiceConnectionIds = append(privateLinkServiceConnectionIds, manualPrivateLinkServiceConnectionIds...)
	locks.MultipleByID(pointer.To(privateLinkServiceConnectionIds))
	defer locks.UnlockMultipleByID(pointer.To(privateLinkServiceConnectionIds))

	parameters := privateendpoints.PrivateEndpoint{
		Location:         pointer.To(location.Normalize(d.Get("location").(string))),
		ExtendedLocation: helpers.ExpandEdgeZoneModel(d.Get("edge_zone").(string)),
		Properties: &privateendpoints.PrivateEndpointProperties{
			PrivateLinkServiceConnections:       privateLinkServiceConnections,
			ManualPrivateLinkServiceConnections: manualPrivateLinkServiceConnections,
			Subnet: &privateendpoints.Subnet{
				Id: pointer.To(d.Get("subnet_id").(string)),
			},
			IPConfigurations:           expandPrivateEndpointIPConfigurations(d.Get("ip_configuration").([]any)),
			CustomNetworkInterfaceName: pointer.To(d.Get("custom_network_interface_name").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if err := validatePrivateLinkServiceId(*parameters.Properties.PrivateLinkServiceConnections); err != nil {
		return err
	}
	if err := validatePrivateLinkServiceId(*parameters.Properties.ManualPrivateLinkServiceConnections); err != nil {
		return err
	}

	cosmosDbResIds := getCosmosDbResIdInPrivateServiceConnections(parameters.Properties)
	for _, cosmosDbResId := range cosmosDbResIds {
		log.Printf("[DEBUG] Add Lock For Private Endpoint %q, lock name: %q", id.PrivateEndpointName, cosmosDbResId)
		locks.ByName(cosmosDbResId, "azurerm_private_endpoint")
		//goland:noinspection GoDeferInLoop
		defer locks.UnlockByName(cosmosDbResId, "azurerm_private_endpoint")
	}

	// TODO: refactor to remove Retry func
	// TODO: implement callback
	if err := pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutCreate), func() *pluginsdk.RetryError {
		result, err := client.CreateOrUpdate(ctx, id, parameters)
		if err != nil {
			return &pluginsdk.RetryError{
				Err: fmt.Errorf("creating %s: %+v", id, err),
			}
		}

		if err := result.Poller.PollUntilDone(ctx); err != nil {
			if lroFailError, ok := errors.AsType[pollers.PollingFailedError](err); ok {
				type lroErrorType struct {
					Error struct {
						Code    string `json:"code"`
						Message string `json:"message"`
					} `json:"error"`
				}

				var lroError lroErrorType
				if err := lroFailError.HttpResponse.Unmarshal(&lroError); err != nil {
					return &pluginsdk.RetryError{
						Err: fmt.Errorf("unmarshaling lro error response: %v", err),
					}
				}

				retryableErrorCodes := []string{"RetryableError", "StorageAccountOperationInProgress"}
				if slices.Contains(retryableErrorCodes, lroError.Error.Code) {
					log.Printf("[WARN] Retry polling %q on error code: %q", id, lroError.Error.Code)
					return &pluginsdk.RetryError{
						Retryable: true,
					}
				}
			}
			return &pluginsdk.RetryError{
				Err: fmt.Errorf("waiting for the creation of %s: %+v", id, err),
			}
		}

		return nil
	}); err != nil {
		return err
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	// 1 Private Endpoint can have 1 Private DNS Zone Group
	// since this is a new resource, there shouldn't be an existing one - so there's no need to delete it
	if len(privateDnsZoneGroup) > 0 {
		if err := createPrivateDnsZoneGroupForPrivateEndpoint(ctx, dnsClient, id, privateDnsZoneGroup); err != nil {
			return err
		}
	}

	return resourcePrivateEndpointRead(d, meta)
}

func expandPrivateLinkEndpointServiceConnection(input []any, parseManual bool) (*[]privateendpoints.PrivateLinkServiceConnection, []string) {
	results := make([]privateendpoints.PrivateLinkServiceConnection, 0)
	privateConnectionResourceIds := make([]string, 0)

	for _, item := range input {
		v := item.(map[string]any)
		privateConnectionResourceId := v["private_connection_resource_id"].(string)
		if privateConnectionResourceId == "" {
			privateConnectionResourceId = v["private_connection_resource_alias"].(string)
		}
		privateConnectionResourceIds = append(privateConnectionResourceIds, privateConnectionResourceId)

		subresourceNames := v["subresource_names"].([]any)
		requestMessage := v["request_message"].(string)
		isManual := v["is_manual_connection"].(bool)
		name := v["name"].(string)

		if isManual == parseManual {
			result := privateendpoints.PrivateLinkServiceConnection{
				Name: pointer.To(name),
				Properties: &privateendpoints.PrivateLinkServiceConnectionProperties{
					GroupIds:             pluginsdk.ExpandStringSlice(subresourceNames),
					PrivateLinkServiceId: pointer.To(privateConnectionResourceId),
				},
			}

			if requestMessage != "" {
				result.Properties.RequestMessage = pointer.To(requestMessage)
			}

			results = append(results, result)
		}
	}

	return &results, privateConnectionResourceIds
}

func expandPrivateEndpointIPConfigurations(input []any) *[]privateendpoints.PrivateEndpointIPConfiguration {
	results := make([]privateendpoints.PrivateEndpointIPConfiguration, 0)

	for _, item := range input {
		v := item.(map[string]any)
		privateIPAddress := v["private_ip_address"].(string)
		subResourceName := v["subresource_name"].(string)
		memberName := v["member_name"].(string)
		if memberName == "" {
			memberName = subResourceName
		}
		name := v["name"].(string)
		result := privateendpoints.PrivateEndpointIPConfiguration{
			Name: pointer.To(name),
			Properties: &privateendpoints.PrivateEndpointIPConfigurationProperties{
				PrivateIPAddress: pointer.To(privateIPAddress),
				GroupId:          pointer.To(subResourceName),
				MemberName:       pointer.To(memberName),
			},
		}
		results = append(results, result)
	}

	return &results
}
