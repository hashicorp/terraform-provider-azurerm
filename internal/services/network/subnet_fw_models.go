package network

import (
	"github.com/hashicorp/go-azure-helpers/framework/commonmodels"
	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TODO: generate identity model into its own file
type subnetFWResourceIdentityModel struct {
	SubscriptionID     types.String `tfsdk:"subscription_id"`
	ResourceGroupName  types.String `tfsdk:"resource_group_name"`
	VirtualNetworkName types.String `tfsdk:"virtual_network_name"`
	Name               types.String `tfsdk:"name"`
}

type subnetModel struct {
	commonmodels.BaseResourceModel

	Name types.String `tfsdk:"name"`
	// TODO: replace with `VirtualNetworkID`?
	ResourceGroupName  types.String `tfsdk:"resource_group_name"`
	VirtualNetworkName types.String `tfsdk:"virtual_network_name"`

	AddressPrefixes                          types.List                                                    `tfsdk:"address_prefixes"`
	DefaultOutboundAccessEnabled             types.Bool                                                    `tfsdk:"default_outbound_access_enabled"`
	Delegation                               typehelpers.ListNestedObjectValueOf[subnetDelegationModel]    `tfsdk:"delegation"`
	IPAddressPool                            typehelpers.ListNestedObjectValueOf[subnetIPAddressPoolModel] `tfsdk:"ip_address_pool"`
	NetworkSecurityGroupID                   types.String                                                  `tfsdk:"network_security_group_id"`
	PrivateEndpointNetworkPolicies           types.String                                                  `tfsdk:"private_endpoint_network_policies"`
	PrivateLinkServiceNetworkPoliciesEnabled types.Bool                                                    `tfsdk:"private_link_service_network_policies_enabled"`
	SharingScope                             types.String                                                  `tfsdk:"sharing_scope"`
	ServiceEndpointPolicyIds                 types.Set                                                     `tfsdk:"service_endpoint_policy_ids"`
	ServiceEndpoints                         types.Set                                                     `tfsdk:"service_endpoints"`
}

type subnetDelegationModel struct {
	Name              types.String                                                      `tfsdk:"name"`
	ServiceDelegation typehelpers.ListNestedObjectValueOf[subnetServiceDelegationModel] `tfsdk:"service_delegation"` // TODO: consider renaming to `service`? it's already located inside the `delegation` block
}

type subnetServiceDelegationModel struct {
	Name    types.String                         `tfsdk:"name"`
	Actions typehelpers.SetValueOf[types.String] `tfsdk:"actions"`
}

type subnetIPAddressPoolModel struct {
	ID                         types.String                          `tfsdk:"id"`
	NumberOfIPAddresses        types.String                          `tfsdk:"number_of_ip_addresses"`
	AllocatedIPAddressPrefixes typehelpers.ListValueOf[types.String] `tfsdk:"allocated_ip_address_prefixes"`
}
