// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package databricks_workspace

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2026-01-01/accessconnector"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2026-01-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDatabricksWorkspaceCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataBricks.WorkspacesClient
	acClient := meta.(*clients.Client).DataBricks.AccessConnectorClient
	lbClient := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	subnetsClient := meta.(*clients.Client).Network.Subnets

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := workspaces.NewWorkspaceID(meta.(*clients.Client).Account.SubscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_databricks_workspace", id.ID())
		}
	}

	var backendPoolName, loadBalancerId string
	managedResourceGroupName := d.Get("managed_resource_group_name").(string)

	if backendPool := d.Get("load_balancer_backend_address_pool_id").(string); backendPool != "" {
		backendPoolId, err := loadbalancers.ParseLoadBalancerBackendAddressPoolID(backendPool)
		if err != nil {
			return err
		}

		lbId := loadbalancers.NewProviderLoadBalancerID(backendPoolId.SubscriptionId, backendPoolId.ResourceGroupName, backendPoolId.LoadBalancerName)
		backendPoolName = backendPoolId.BackendAddressPoolName
		loadBalancerId = lbId.ID()

		locks.ByID(backendPoolId.ID())
		defer locks.UnlockByID(backendPoolId.ID())

		locks.ByID(lbId.ID())
		defer locks.UnlockByID(lbId.ID())

		// check to make sure the load balancer exists as referred to by the Backend Address Pool...
		lb, err := lbClient.Get(ctx, lbId, loadbalancers.GetOperationOptions{})
		if err != nil {
			if response.WasNotFound(lb.HttpResponse) {
				return fmt.Errorf("%s was not found", lbId)
			}
			return fmt.Errorf("retrieving %s: %+v", lbId, err)
		}
	}

	if managedResourceGroupName == "" {
		log.Printf("[DEBUG] no managed resource group name was provided, using the default pattern")
		managedResourceGroupName = fmt.Sprintf("databricks-rg-%s", id.ResourceGroupName)
	}

	publicNetworkAccess := workspaces.PublicNetworkAccessDisabled
	if d.Get("public_network_access_enabled").(bool) {
		publicNetworkAccess = workspaces.PublicNetworkAccessEnabled
	}

	customParamsRaw := d.Get("custom_parameters").([]any)
	customParams, pubSubAssoc, priSubAssoc := expandWorkspaceCustomParameters(customParamsRaw, d.Get("customer_managed_key_enabled").(bool), d.Get("infrastructure_encryption_enabled").(bool), backendPoolName, loadBalancerId)

	if len(customParamsRaw) > 0 && customParamsRaw[0] != nil {
		config := customParamsRaw[0].(map[string]any)
		pubSub := config["public_subnet_name"].(string)
		priSub := config["private_subnet_name"].(string)
		vnetID := config["virtual_network_id"].(string)

		if config["virtual_network_id"].(string) == "" && (pubSub != "" || priSub != "") {
			return fmt.Errorf("`public_subnet_name` and/or `private_subnet_name` cannot be defined if `virtual_network_id` is not set")
		}
		if config["virtual_network_id"].(string) != "" && (pubSub == "" || priSub == "") {
			return fmt.Errorf("`public_subnet_name` and `private_subnet_name` must both have values if `virtual_network_id` is set")
		}
		if pubSub != "" && pubSubAssoc == nil {
			return fmt.Errorf("you must define a value for `public_subnet_network_security_group_association_id` if `public_subnet_name` is set")
		}
		if priSub != "" && priSubAssoc == nil {
			return fmt.Errorf("you must define a value for `private_subnet_network_security_group_association_id` if `private_subnet_name` is set")
		}

		if subnetDelegationErr := checkSubnetDelegations(ctx, subnetsClient, vnetID, pubSub, priSub); subnetDelegationErr != nil {
			return subnetDelegationErr
		}
	}

	encryption, err := expandDatabricksWorkspaceEncryption(d)
	if err != nil {
		return fmt.Errorf("expanding workspace encryption: %+v", err)
	}

	workspace := workspaces.Workspace{
		Sku: &workspaces.Sku{
			Name: d.Get("sku").(string),
		},
		Location: location.Normalize(d.Get("location").(string)),
		Properties: workspaces.WorkspaceProperties{
			ComputeMode:                workspaces.ComputeModeHybrid,
			PublicNetworkAccess:        &publicNetworkAccess,
			ManagedResourceGroupId:     pointer.To(commonids.NewResourceGroupID(id.SubscriptionId, managedResourceGroupName).ID()),
			Parameters:                 customParams,
			Encryption:                 encryption,
			EnhancedSecurityCompliance: expandWorkspaceEnhancedSecurity(d.Get("enhanced_security_compliance").([]any)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if d.Get("default_storage_firewall_enabled").(bool) {
		accessConnectorId, err := accessconnector.ParseAccessConnectorID(d.Get("access_connector_id").(string))
		if err != nil {
			return err
		}

		accessConnector, err := acClient.Get(ctx, *accessConnectorId)
		if err != nil {
			return fmt.Errorf("retrieving Access Connector %s: %+v", accessConnectorId.AccessConnectorName, err)
		}

		accessConnectorProperties := workspaces.WorkspacePropertiesAccessConnector{}
		if model := accessConnector.Model; model != nil && model.Identity != nil {
			accIdentityId := ""
			for raw := range model.Identity.IdentityIds {
				id, err := commonids.ParseUserAssignedIdentityIDInsensitively(raw)
				if err != nil {
					return err
				}
				accIdentityId = id.ID()
				break
			}

			accessConnectorProperties.Id = pointer.From(model.Id)
			accessConnectorProperties.IdentityType = workspaces.IdentityType(model.Identity.Type)
			accessConnectorProperties.UserAssignedIdentityId = &accIdentityId
		}

		workspace.Properties.AccessConnector = &accessConnectorProperties
		workspace.Properties.DefaultStorageFirewall = pointer.To(workspaces.DefaultStorageFirewallEnabled)
	}

	if requireNsgRules := d.Get("network_security_group_rules_required").(string); requireNsgRules != "" {
		workspace.Properties.RequiredNsgRules = pointer.ToEnum[workspaces.RequiredNsgRules](requireNsgRules)
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, workspace, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	// I have to set the custom_parameters so I can pass the public and private
	// subnet NSG association along with the backend Pool Id since they are not
	// returned in the read from Azure...
	custom, backendPoolReadId := flattenWorkspaceCustomParameters(customParams, pubSubAssoc, priSubAssoc)
	d.Set("load_balancer_backend_address_pool_id", backendPoolReadId)

	if err := d.Set("custom_parameters", custom); err != nil {
		return fmt.Errorf("setting `custom_parameters`: %+v", err)
	}

	return resourceDatabricksWorkspaceRead(d, meta)
}
