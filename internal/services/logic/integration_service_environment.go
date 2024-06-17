// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logic

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationserviceenvironments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualnetworks"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceIntegrationServiceEnvironment() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create:             resourceIntegrationServiceEnvironmentCreateUpdate,
		Read:               resourceIntegrationServiceEnvironmentRead,
		Update:             resourceIntegrationServiceEnvironmentCreateUpdate,
		Delete:             resourceIntegrationServiceEnvironmentDelete,
		DeprecationMessage: "The \"azurerm_integrated_service_environment\" resource is deprecated and will be removed in v4.0 of the Azure Provider. The underlying Azure Service is being retired on 2024-08-31 and new instances cannot be provisioned by default after 2022-11-01. More information on the retirement and how to migrate to Logic Apps Standard (\"azurerm_logic_app_standard\") can be found here: https://aka.ms/isedeprecation.",

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := integrationserviceenvironments.ParseIntegrationServiceEnvironmentID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(5 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(5 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(5 * time.Hour),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IntegrationServiceEnvironmentName(),
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			// Maximum scale units that you can add	10 - https://docs.microsoft.com/en-US/azure/logic-apps/logic-apps-limits-and-config#integration-service-environment-ise
			// Developer Always 0 capacity
			"sku_name": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "Developer_0",
				ValidateFunc: validation.StringInSlice([]string{
					"Developer_0",
					"Premium_0",
					"Premium_1",
					"Premium_2",
					"Premium_3",
					"Premium_4",
					"Premium_5",
					"Premium_6",
					"Premium_7",
					"Premium_8",
					"Premium_9",
					"Premium_10",
				}, false),
			},

			"access_endpoint_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true, // The access end point type cannot be changed once the integration service environment is provisioned.
				ValidateFunc: validation.StringInSlice([]string{
					string(integrationserviceenvironments.IntegrationServiceEnvironmentAccessEndpointTypeInternal),
					string(integrationserviceenvironments.IntegrationServiceEnvironmentAccessEndpointTypeExternal),
				}, false),
			},

			"virtual_network_subnet_ids": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				ForceNew: true, // The network configuration subnets cannot be updated after integration service environment is created.
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: commonids.ValidateSubnetID,
				},
				MinItems: 4,
				MaxItems: 4,
			},

			"connector_endpoint_ip_addresses": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			},

			"connector_outbound_ip_addresses": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			},

			"workflow_endpoint_ip_addresses": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			},

			"workflow_outbound_ip_addresses": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			},

			"tags": commonschema.Tags(),
		},

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.ForceNewIfChange("sku_name", func(ctx context.Context, old, new, meta interface{}) bool {
				oldSku := strings.Split(old.(string), "_")
				newSku := strings.Split(new.(string), "_")
				// The SKU cannot be changed once integration service environment has been provisioned. -> we need ForceNew
				return oldSku[0] != newSku[0]
			}),
		),
	}
}

func resourceIntegrationServiceEnvironmentCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.IntegrationServiceEnvironmentClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM Integration Service Environment creation.")

	id := integrationserviceenvironments.NewIntegrationServiceEnvironmentID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_integration_service_environment", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	accessEndpointType := integrationserviceenvironments.IntegrationServiceEnvironmentAccessEndpointType(d.Get("access_endpoint_type").(string))
	virtualNetworkSubnetIds := d.Get("virtual_network_subnet_ids").(*pluginsdk.Set).List()
	t := d.Get("tags").(map[string]interface{})

	sku, err := expandIntegrationServiceEnvironmentSkuName(d.Get("sku_name").(string))
	if err != nil {
		return fmt.Errorf("expanding `sku_name` for %s: %v", id, err)
	}

	integrationServiceEnvironment := integrationserviceenvironments.IntegrationServiceEnvironment{
		Name:     &id.IntegrationServiceEnvironmentName,
		Location: &location,
		Properties: &integrationserviceenvironments.IntegrationServiceEnvironmentProperties{
			NetworkConfiguration: &integrationserviceenvironments.NetworkConfiguration{
				AccessEndpoint: &integrationserviceenvironments.IntegrationServiceEnvironmentAccessEndpoint{
					Type: &accessEndpointType,
				},
				Subnets: expandSubnetResourceID(virtualNetworkSubnetIds),
			},
		},
		Sku:  sku,
		Tags: tags.Expand(t),
	}

	if err = client.CreateOrUpdateThenPoll(ctx, id, integrationServiceEnvironment); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceIntegrationServiceEnvironmentRead(d, meta)
}

func resourceIntegrationServiceEnvironmentRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.IntegrationServiceEnvironmentClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := integrationserviceenvironments.ParseIntegrationServiceEnvironmentID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id.ID(), err)
	}

	d.Set("name", id.IntegrationServiceEnvironmentName)
	d.Set("resource_group_name", id.ResourceGroup)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if err := d.Set("sku_name", flattenIntegrationServiceEnvironmentSkuName(model.Sku)); err != nil {
			return fmt.Errorf("setting `sku_name`: %+v", err)
		}

		if props := model.Properties; props != nil {
			if netCfg := props.NetworkConfiguration; netCfg != nil {
				if accessEndpoint := netCfg.AccessEndpoint; accessEndpoint != nil {
					d.Set("access_endpoint_type", string(pointer.From(accessEndpoint.Type)))
				}

				d.Set("virtual_network_subnet_ids", flattenSubnetResourceID(netCfg.Subnets))
			}

			if props.EndpointsConfiguration == nil || props.EndpointsConfiguration.Connector == nil {
				d.Set("connector_endpoint_ip_addresses", []interface{}{})
				d.Set("connector_outbound_ip_addresses", []interface{}{})
			} else {
				d.Set("connector_endpoint_ip_addresses", flattenServiceEnvironmentIPAddresses(props.EndpointsConfiguration.Connector.AccessEndpointIPAddresses))
				d.Set("connector_outbound_ip_addresses", flattenServiceEnvironmentIPAddresses(props.EndpointsConfiguration.Connector.OutgoingIPAddresses))
			}

			if props.EndpointsConfiguration == nil || props.EndpointsConfiguration.Workflow == nil {
				d.Set("workflow_endpoint_ip_addresses", []interface{}{})
				d.Set("workflow_outbound_ip_addresses", []interface{}{})
			} else {
				d.Set("workflow_endpoint_ip_addresses", flattenServiceEnvironmentIPAddresses(props.EndpointsConfiguration.Workflow.AccessEndpointIPAddresses))
				d.Set("workflow_outbound_ip_addresses", flattenServiceEnvironmentIPAddresses(props.EndpointsConfiguration.Workflow.OutgoingIPAddresses))
			}
		}

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func resourceIntegrationServiceEnvironmentDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.IntegrationServiceEnvironmentClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := integrationserviceenvironments.ParseIntegrationServiceEnvironmentID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing Integration Service Environment ID `%q`: %+v", d.Id(), err)
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id.ID(), err)
	}

	// Get subnet IDs before delete
	subnetIDs := getSubnetIDs(resp.Model)

	// Not optimal behaviour for now
	// It deletes synchronously and resource is not available anymore after return from delete operation
	// Next, after return - delete operation is still in progress in the background and is still occupying subnets.
	// As workaround we are checking on all involved subnets presence of serviceAssociationLink and resourceNavigationLink
	// If the operation fails we are lost. We do not have original resource and we cannot resume delete operation.
	// User has to wait for completion of delete operation in the background.
	// It would be great to have async call with future struct
	if resp, err := client.Delete(ctx, *id); err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}

		return fmt.Errorf("deleting %s: %+v", id.ID(), err)
	}

	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{string(integrationserviceenvironments.WorkflowProvisioningStateDeleting)},
		Target:                    []string{string(integrationserviceenvironments.WorkflowProvisioningStateDeleted)},
		MinTimeout:                5 * time.Minute,
		Refresh:                   integrationServiceEnvironmentDeleteStateRefreshFunc(ctx, meta.(*clients.Client), d.Id(), subnetIDs),
		Timeout:                   d.Timeout(pluginsdk.TimeoutDelete),
		ContinuousTargetOccurence: 1,
		NotFoundChecks:            1,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", id.ID(), err)
	}

	return nil
}

func flattenIntegrationServiceEnvironmentSkuName(input *integrationserviceenvironments.IntegrationServiceEnvironmentSku) string {
	if input == nil {
		return ""
	}

	name := ""
	if input.Name != nil {
		name = string(*input.Name)
	}

	return fmt.Sprintf("%s_%d", name, *input.Capacity)
}

func expandIntegrationServiceEnvironmentSkuName(skuName string) (*integrationserviceenvironments.IntegrationServiceEnvironmentSku, error) {
	parts := strings.Split(skuName, "_")
	if len(parts) != 2 {
		return nil, fmt.Errorf("sku_name (%s) has the wrong number of parts (%d) after splitting on _", skuName, len(parts))
	}

	var sku integrationserviceenvironments.IntegrationServiceEnvironmentSkuName
	switch parts[0] {
	case string(integrationserviceenvironments.IntegrationServiceEnvironmentSkuNameDeveloper):
		sku = integrationserviceenvironments.IntegrationServiceEnvironmentSkuNameDeveloper
	case string(integrationserviceenvironments.IntegrationServiceEnvironmentSkuNamePremium):
		sku = integrationserviceenvironments.IntegrationServiceEnvironmentSkuNamePremium
	default:
		return nil, fmt.Errorf("sku_name %s has unknown sku %s", skuName, parts[0])
	}

	capacity, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("cannot convert sku_name %s capacity %s to int", skuName, parts[1])
	}

	if sku != integrationserviceenvironments.IntegrationServiceEnvironmentSkuNamePremium && capacity > 0 {
		return nil, fmt.Errorf("`capacity` can only be greater than zero for `sku_name` `Premium`")
	}

	return &integrationserviceenvironments.IntegrationServiceEnvironmentSku{
		Name:     &sku,
		Capacity: utils.Int64(int64(capacity)),
	}, nil
}

func expandSubnetResourceID(input []interface{}) *[]integrationserviceenvironments.ResourceReference {
	results := make([]integrationserviceenvironments.ResourceReference, 0)
	for _, item := range input {
		results = append(results, integrationserviceenvironments.ResourceReference{
			Id: utils.String(item.(string)),
		})
	}
	return &results
}

func flattenSubnetResourceID(input *[]integrationserviceenvironments.ResourceReference) []interface{} {
	subnetIDs := make([]interface{}, 0)
	if input == nil {
		return subnetIDs
	}

	for _, resourceRef := range *input {
		if resourceRef.Id == nil || *resourceRef.Id == "" {
			continue
		}

		subnetIDs = append(subnetIDs, resourceRef.Id)
	}

	return subnetIDs
}

func getSubnetIDs(input *integrationserviceenvironments.IntegrationServiceEnvironment) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	if props := input.Properties; props != nil {
		if netCfg := props.NetworkConfiguration; netCfg != nil {
			return flattenSubnetResourceID(netCfg.Subnets)
		}
	}

	return results
}

func integrationServiceEnvironmentDeleteStateRefreshFunc(ctx context.Context, client *clients.Client, iseID string, subnetIDs []interface{}) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		linkExists, err := linkExists(ctx, client, iseID, subnetIDs)
		if err != nil {
			return string(integrationserviceenvironments.WorkflowProvisioningStateDeleting), string(integrationserviceenvironments.WorkflowProvisioningStateDeleting), err
		}

		if linkExists {
			return string(integrationserviceenvironments.WorkflowProvisioningStateDeleting), string(integrationserviceenvironments.WorkflowProvisioningStateDeleting), nil
		}

		return string(integrationserviceenvironments.WorkflowProvisioningStateDeleted), string(integrationserviceenvironments.WorkflowProvisioningStateDeleted), nil
	}
}

func linkExists(ctx context.Context, client *clients.Client, iseID string, subnetIDs []interface{}) (bool, error) {
	for _, subnetID := range subnetIDs {
		if subnetID == nil {
			continue
		}

		id := *(subnetID.(*string))
		log.Printf("Checking links on subnetID: %q\n", id)

		hasLink, err := serviceAssociationLinkExists(ctx, client.Network.VirtualNetworks, iseID, id)
		if err != nil {
			return false, err
		}

		if hasLink {
			return true, nil
		} else {
			hasLink, err := resourceNavigationLinkExists(ctx, client.Network.VirtualNetworks, id)
			if err != nil {
				return false, err
			}

			if hasLink {
				return true, nil
			}
		}
	}

	return false, nil
}

func serviceAssociationLinkExists(ctx context.Context, client *virtualnetworks.VirtualNetworksClient, iseID string, subnetID string) (bool, error) {
	id, err := commonids.ParseSubnetID(subnetID)
	if err != nil {
		return false, err
	}

	resp, err := client.ServiceAssociationLinksList(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return false, nil
		}
		return false, fmt.Errorf("retrieving Service Association Links from Virtual Network %q, subnet %q (Resource Group %q): %+v", id.VirtualNetworkName, id.SubscriptionId, id.ResourceGroupName, err)
	}

	if model := resp.Model; model != nil {
		for _, link := range *model {
			if link.Properties != nil && link.Properties.Link != nil {
				if strings.EqualFold(iseID, *link.Properties.Link) {
					log.Printf("Has Service Association Link: %q\n", *link.Id)
					return true, nil
				}
			}
		}
	}

	return false, nil
}

func resourceNavigationLinkExists(ctx context.Context, client *virtualnetworks.VirtualNetworksClient, subnetID string) (bool, error) {
	id, err := commonids.ParseSubnetID(subnetID)
	if err != nil {
		return false, err
	}

	resp, err := client.ResourceNavigationLinksList(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return false, nil
		}
		return false, fmt.Errorf("retrieving Resource Navigation Links from Virtual Network %q, subnet %q (Resource Group %q): %+v", id.VirtualNetworkName, id.SubnetName, id.ResourceGroupName, err)
	}

	if model := resp.Model; model != nil {
		for _, link := range *model {
			log.Printf("Has Resource Navigation Link: %q\n", *link.Id)
			return true, nil
		}
	}

	return false, nil
}

func flattenServiceEnvironmentIPAddresses(input *[]integrationserviceenvironments.IPAddress) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	var addresses []interface{}
	for _, addr := range *input {
		addresses = append(addresses, *addr.Address)
	}
	return addresses
}
