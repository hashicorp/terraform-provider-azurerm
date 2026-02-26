package network

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/go-azure-helpers/framework/commonschema"
	"github.com/hashicorp/go-azure-helpers/framework/convert"
	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/serviceendpointpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/ipampools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/networksecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/subnets"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

const (
	privateStateNetworkSecurityGroupID = "network_security_group_id"
)

type SubnetFWResource struct{}

var _ sdk.FrameworkWrappedResourceWithUpdate = &SubnetFWResource{}

func (r SubnetFWResource) ModelObject() any {
	return &subnetModel{}
}

func (r SubnetFWResource) ResourceType() string {
	return "azurerm_subnet_fw"
}

func (r SubnetFWResource) Schema(ctx context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			commonschema.ID: commonschema.IDAttribute(),

			commonschema.Name: schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			commonschema.ResourceGroupName: commonschema.ResourceGroupNameAttribute(),

			"virtual_network_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"address_prefixes": schema.ListAttribute{
				Optional: true,
				// Note: O+C because when `ip_address_pool` is used, Azure returns a CIDR range provisioned by the IP Address Management Pool
				Computed: true, // TODO: test whether this becomes problematic, this is replacing the existing `DiffSuppressFunc`
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ExactlyOneOf(path.MatchRoot("address_prefixes"), path.MatchRoot("ip_address_pool")),
				},
				ElementType: types.StringType,
			},

			"service_endpoints": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},

			"service_endpoint_policy_ids": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						typehelpers.WrappedStringValidator{
							Func: serviceendpointpolicies.ValidateServiceEndpointPolicyID,
						},
					),
				},
			},

			"sharing_scope": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf(subnets.PossibleValuesForSharingScope()...), // TODO: confirm whether `DelegatedServices` works now, if not, remove from possible values
				},
			},

			"default_outbound_access_enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},

			"private_endpoint_network_policies": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(string(subnets.VirtualNetworkPrivateEndpointNetworkPoliciesDisabled)),
				Validators: []validator.String{
					stringvalidator.OneOf(subnets.PossibleValuesForVirtualNetworkPrivateEndpointNetworkPolicies()...),
				},
			},

			"private_link_service_network_policies_enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},

			// TODO: This is currently managed using the `azurerm_subnet_network_security_group_association` resource
			// however there are many requests to provide an argument to manage this inline. This is problematic for multiple reasons:
			// - recreation of NSG fails due to it still being attached to the subnet
			// - cannot be optional only, conflicts with the existing association resource
			// - cannot be optional+computed, prevents removal
			// Determine whether there are workarounds, otherwise this should not be introduced to this resource.
			//"network_security_group_id": schema.StringAttribute{
			//	Optional: true,
			//	Validators: []validator.String{
			//		typehelpers.WrappedStringValidator{
			//			Func: networksecuritygroups.ValidateNetworkSecurityGroupID,
			//		},
			//	},
			//},
		},

		Blocks: map[string]schema.Block{
			"delegation": schema.ListNestedBlock{
				CustomType: typehelpers.NewListNestedObjectTypeOf[subnetDelegationModel](ctx),
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required: true,
						},
					},

					Blocks: map[string]schema.Block{
						"service_delegation": schema.ListNestedBlock{
							CustomType: typehelpers.NewListNestedObjectTypeOf[subnetServiceDelegationModel](ctx),
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Required: true,
										Validators: []validator.String{
											stringvalidator.OneOf(subnetDelegationServiceNames...),
										},
									},

									"actions": schema.SetAttribute{
										Optional: true,
										Validators: []validator.Set{
											setvalidator.ValueStringsAre(
												stringvalidator.OneOf(
													"Microsoft.Network/networkinterfaces/*",
													"Microsoft.Network/publicIPAddresses/join/action",
													"Microsoft.Network/publicIPAddresses/read",
													"Microsoft.Network/virtualNetworks/read",
													"Microsoft.Network/virtualNetworks/subnets/action",
													"Microsoft.Network/virtualNetworks/subnets/join/action",
													"Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action",
													"Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action",
												),
											),
										},
										ElementType: types.StringType,
									},
								},
							},
						},
					},
				},
			},

			"ip_address_pool": schema.ListNestedBlock{
				CustomType: typehelpers.NewListNestedObjectTypeOf[subnetIPAddressPoolModel](ctx),
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Required: true,
							Validators: []validator.String{
								typehelpers.WrappedStringValidator{
									Func: ipampools.ValidateIPamPoolID,
								},
							},
						},

						"number_of_ip_addresses": schema.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(regexp.MustCompile(`^[1-9]\d*$`),
									"value must be a string representing a positive integer",
								),
							},
						},

						"allocated_ip_address_prefixes": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

func (r SubnetFWResource) Create(ctx context.Context, _ resource.CreateRequest, resp *resource.CreateResponse, metadata sdk.ResourceMetadata, plan any) {
	client := metadata.Client.Network.Subnets
	diags := pointer.To(resp.Diagnostics)

	data := sdk.AssertResourceModelType[subnetModel](plan, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	id := commonids.NewSubnetID(metadata.SubscriptionId, data.ResourceGroupName.ValueString(), data.VirtualNetworkName.ValueString(), data.Name.ValueString())

	existing, err := client.Get(ctx, id, subnets.DefaultGetOperationOptions())
	if !response.WasNotFound(existing.HttpResponse) {
		if err != nil {
			sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("Checking for presence of existing %s", id), err)
			return
		}

		metadata.ResourceRequiresImport(r.ResourceType(), id, resp)
	}

	virtualNetworkID := commonids.NewVirtualNetworkID(id.SubscriptionId, id.ResourceGroupName, id.VirtualNetworkName)
	locks.ByID(virtualNetworkID.ID())
	defer locks.UnlockByID(virtualNetworkID.ID())

	payload := subnets.Subnet{
		Name: data.Name.ValueStringPointer(),
		Properties: &subnets.SubnetPropertiesFormat{
			AddressPrefixes:                   pointer.To(typehelpers.ExpandListNestedTypeOf[string](data.AddressPrefixes, diags)),
			DefaultOutboundAccess:             data.DefaultOutboundAccessEnabled.ValueBoolPointer(),
			Delegations:                       expandSubnetDelegationFW(ctx, data.Delegation, diags),
			IPamPoolPrefixAllocations:         expandSubnetIPAddressPoolFW(ctx, data.IPAddressPool, diags),
			NetworkSecurityGroup:              expandSubnetNetworkSecurityGroupID(data.NetworkSecurityGroupID),
			PrivateEndpointNetworkPolicies:    pointer.ToEnum[subnets.VirtualNetworkPrivateEndpointNetworkPolicies](data.PrivateEndpointNetworkPolicies.ValueString()),
			PrivateLinkServiceNetworkPolicies: expandSubnetNetworkPolicyFW(data.PrivateLinkServiceNetworkPoliciesEnabled),
			ServiceEndpointPolicies:           expandSubnetServiceEndpointPolicyIDsFW(ctx, data.ServiceEndpointPolicyIds, diags),
			ServiceEndpoints:                  expandSubnetServiceEndpointsFW(ctx, data.ServiceEndpoints, diags),
			SharingScope:                      pointer.ToEnum[subnets.SharingScope](data.SharingScope.ValueString()),
		},
	}

	if resp.Diagnostics.HasError() {
		return
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("Creating %s", id), err)
		return
	}

	data.ID = types.StringValue(id.ID())

	readResp, err := client.Get(ctx, id, subnets.DefaultGetOperationOptions())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("Retrieving %s", id), err)
		return
	}

	r.flatten(ctx, id, readResp.Model, data, &resp.Diagnostics)
}

func (r SubnetFWResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, metadata sdk.ResourceMetadata, state any) {
	client := metadata.Client.Network.Subnets

	data := sdk.AssertResourceModelType[subnetModel](state, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := commonids.ParseSubnetID(data.ID.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(resp, "Parsing ID", err)
		return
	}

	readResp, err := client.Get(ctx, *id, subnets.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(readResp.HttpResponse) {
			metadata.MarkAsGone(ctx, id, resp, &resp.Diagnostics)
			data = nil
			return
		}

		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("Retrieving %s", id), err)
		return
	}

	r.flatten(ctx, *id, readResp.Model, data, &resp.Diagnostics)
}

func (r SubnetFWResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, metadata sdk.ResourceMetadata, decodedPlan any, decodedState any) {
	// TODO: do we need HasChange checks?
	client := metadata.Client.Network.Subnets

	plan := sdk.AssertResourceModelType[subnetModel](decodedPlan, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	state := sdk.AssertResourceModelType[subnetModel](decodedState, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := commonids.ParseSubnetID(state.ID.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(resp, "Parsing Resource ID", err)
	}

	existing, err := client.Get(ctx, *id, subnets.DefaultGetOperationOptions())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("Retrieving %s", id), err)
		return
	}

	if existing.Model == nil {
		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("Retrieving %s", id), "`model` was nil")
	}

	if existing.Model.Properties == nil {
		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("Retrieving %s", id), "`properties` was nil")
	}
	props := existing.Model.Properties

	diags := pointer.To(resp.Diagnostics)

	if sdk.HasChange(plan.AddressPrefixes, state.AddressPrefixes) {
		switch len(plan.AddressPrefixes.Elements()) {
		case 0:
			// When `ip_address_pool` is used, clear AddressPrefixes and AddressPrefix. These values will be computed by the service.
			props.AddressPrefixes = nil
			props.AddressPrefix = nil
		default:
			props.AddressPrefixes = pointer.To(convert.ExpandAndReturn[[]string](ctx, plan.AddressPrefixes, diags))
		}
	}

	props.DefaultOutboundAccess = plan.DefaultOutboundAccessEnabled.ValueBoolPointer()
	props.Delegations = expandSubnetDelegationFW(ctx, plan.Delegation, diags)
	props.IPamPoolPrefixAllocations = expandSubnetIPAddressPoolFW(ctx, plan.IPAddressPool, diags)
	props.NetworkSecurityGroup = expandSubnetNetworkSecurityGroupID(plan.NetworkSecurityGroupID)
	props.PrivateEndpointNetworkPolicies = pointer.ToEnumFW[subnets.VirtualNetworkPrivateEndpointNetworkPolicies](plan.PrivateEndpointNetworkPolicies)
	props.PrivateLinkServiceNetworkPolicies = expandSubnetNetworkPolicyFW(plan.PrivateLinkServiceNetworkPoliciesEnabled)
	props.ServiceEndpointPolicies = expandSubnetServiceEndpointPolicyIDsFW(ctx, plan.ServiceEndpointPolicyIds, diags)
	props.ServiceEndpoints = expandSubnetServiceEndpointsFW(ctx, plan.ServiceEndpoints, diags)
	props.SharingScope = pointer.ToEnumFW[subnets.SharingScope](plan.SharingScope)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *existing.Model); err != nil {
		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("Updating %s", id), err)
		return
	}

	readResp, err := client.Get(ctx, *id, subnets.DefaultGetOperationOptions())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("Retrieving %s", id), err)
		return
	}

	r.flatten(ctx, *id, readResp.Model, plan, &resp.Diagnostics)
}

func (r SubnetFWResource) Delete(ctx context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse, metadata sdk.ResourceMetadata, state any) {
	client := metadata.Client.Network.Subnets

	data := sdk.AssertResourceModelType[subnetModel](state, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := commonids.ParseSubnetID(data.ID.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(resp, "parsing ID", err)
		return
	}

	virtualNetworkID := commonids.NewVirtualNetworkID(id.SubscriptionId, id.ResourceGroupName, id.VirtualNetworkName)
	locks.ByID(virtualNetworkID.ID())
	defer locks.UnlockByID(virtualNetworkID.ID())

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("deleting %s", id), err)
		return
	}
}

func (r SubnetFWResource) flatten(ctx context.Context, id commonids.SubnetId, model *subnets.Subnet, state *subnetModel, diags *diag.Diagnostics) {
	state.Name = types.StringValue(id.SubnetName)
	state.ResourceGroupName = types.StringValue(id.ResourceGroupName)
	state.VirtualNetworkName = types.StringValue(id.VirtualNetworkName)

	if model != nil {
		if props := model.Properties; props != nil {
			state.AddressPrefixes = typehelpers.ListValueFrom(ctx, types.StringType, props.AddressPrefixes, diags)
			state.DefaultOutboundAccessEnabled = types.BoolPointerValue(props.DefaultOutboundAccess)
			state.Delegation = flattenSubnetDelegationFW(ctx, props.Delegations, diags)
			state.IPAddressPool = flattenSubnetIPAddressPoolFW(ctx, props.IPamPoolPrefixAllocations, diags)
			state.PrivateEndpointNetworkPolicies = pointer.FromEnumFW(props.PrivateEndpointNetworkPolicies)
			state.PrivateLinkServiceNetworkPoliciesEnabled = flattenSubnetNetworkPolicyFW(props.PrivateLinkServiceNetworkPolicies)
			state.SharingScope = pointer.FromEnumFW(props.SharingScope)

			serviceEndpointPolicyIDs, d := flattenSubnetServiceEndpointPolicyIDsFW(ctx, props.ServiceEndpointPolicies)
			if d.HasError() {
				// TODO: do we care to wrap diags similar to how existing resource wraps errors? if yes, add to all flatten funcs, if no, remove below:
				diags.Append(sdk.NewErrorDiagnostic("flattening `service_endpoint_policy_ids`", d))
				return
			}
			state.ServiceEndpointPolicyIds = serviceEndpointPolicyIDs

			serviceEndpoints, d := flattenSubnetServiceEndpointsFW(ctx, props.ServiceEndpoints)
			if d.HasError() {
				diags.Append(sdk.NewErrorDiagnostic("flattening `service_endpoints`", d))
				return
			}
			state.ServiceEndpoints = serviceEndpoints

			if nsg := props.NetworkSecurityGroup; nsg != nil && nsg.Id != nil {
				nsgID, err := networksecuritygroups.ParseNetworkSecurityGroupID(*nsg.Id)
				if err != nil {
					diags.Append(sdk.NewErrorDiagnostic("Parsing Network Security Group ID", err))
				}
				state.NetworkSecurityGroupID = types.StringValue(nsgID.ID())
			}
		}
	}
}

func (r SubnetFWResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse, metadata sdk.ResourceMetadata) {
	if request.ID == "" {
		resourceIdentity := &subnetFWResourceIdentityModel{}
		diags := request.Identity.Get(ctx, resourceIdentity)
		if diags.HasError() {
			response.Diagnostics.Append(diags...)
			return
		}

		id := pointer.To(commonids.NewSubnetID(resourceIdentity.SubscriptionID.ValueString(), resourceIdentity.ResourceGroupName.ValueString(), resourceIdentity.VirtualNetworkName.ValueString(), resourceIdentity.Name.ValueString()))
		response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("id"), id.ID())...)
	}

	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}

func (r SubnetFWResource) Identity() (id resourceids.ResourceId, idType sdk.ResourceTypeForIdentity) {
	return &commonids.SubnetId{}, sdk.ResourceTypeForIdentityDefault
}

// Helpers / expands / flattens specifically for `azurerm_subnet`

func expandSubnetNetworkPolicyFW(enabled types.Bool) *subnets.VirtualNetworkPrivateLinkServiceNetworkPolicies {
	if enabled.ValueBool() {
		return pointer.To(subnets.VirtualNetworkPrivateLinkServiceNetworkPoliciesEnabled)
	}

	return pointer.To(subnets.VirtualNetworkPrivateLinkServiceNetworkPoliciesDisabled)
}

// TODO: confirm whether casing is an issue, previously this did a case insensitive comparison
func flattenSubnetNetworkPolicyFW(input *subnets.VirtualNetworkPrivateLinkServiceNetworkPolicies) types.Bool {
	return types.BoolValue(pointer.From(input) == subnets.VirtualNetworkPrivateLinkServiceNetworkPoliciesEnabled)
}

func expandSubnetDelegationFW(ctx context.Context, input typehelpers.ListNestedObjectValueOf[subnetDelegationModel], diags *diag.Diagnostics) *[]subnets.Delegation {
	delegations := typehelpers.DecodeObjectList(ctx, input, diags)

	output := make([]subnets.Delegation, 0, len(delegations))
	for _, d := range delegations {
		serviceDelegation := typehelpers.DecodeObjectListOfOne(ctx, d.ServiceDelegation, diags)

		output = append(output, subnets.Delegation{
			Name: d.Name.ValueStringPointer(),
			Properties: &subnets.ServiceDelegationPropertiesFormat{
				ServiceName: serviceDelegation.Name.ValueStringPointer(),
				Actions:     pointer.To(convert.ExpandAndReturn[[]string](ctx, serviceDelegation.Actions, diags)),
			},
		})
	}

	return &output
}

func flattenSubnetDelegationFW(ctx context.Context, input *[]subnets.Delegation, diags *diag.Diagnostics) typehelpers.ListNestedObjectValueOf[subnetDelegationModel] {
	if input == nil || len(*input) == 0 {
		return typehelpers.NewListNestedObjectValueOfNull[subnetDelegationModel](ctx)
	}

	result := make([]subnetDelegationModel, 0, len(*input))
	for _, d := range *input {
		result = append(result, subnetDelegationModel{
			Name:              types.StringPointerValue(d.Name),
			ServiceDelegation: flattenSubnetServiceDelegation(ctx, d.Properties, diags),
		})
	}

	// TODO: refactor helper to accept a pointer to a diags object and append directly into it?
	wrappedObj, d := typehelpers.NewListNestedObjectValueOfValueSlice[subnetDelegationModel](ctx, result)
	diags.Append(d...)

	return wrappedObj
}

func flattenSubnetServiceDelegation(ctx context.Context, input *subnets.ServiceDelegationPropertiesFormat, diags *diag.Diagnostics) typehelpers.ListNestedObjectValueOf[subnetServiceDelegationModel] {
	if input == nil {
		return typehelpers.NewListNestedObjectValueOfNull[subnetServiceDelegationModel](ctx)
	}

	normalizedServiceName := map[string]string{}
	for _, n := range subnetDelegationServiceNames {
		normalizedServiceName[strings.ToLower(n)] = n
	}

	result := subnetServiceDelegationModel{
		Name:    types.StringPointerValue(input.ServiceName),
		Actions: flattenSubnetServiceDelegationActions(ctx, input.Actions, diags),
	}

	// Normalizing to match existing `azurerm_subnet` logic. TODO: we may want to remove this in the future
	if normalizedName, ok := normalizedServiceName[strings.ToLower(pointer.From(input.ServiceName))]; ok {
		result.Name = types.StringValue(normalizedName)
	}

	// TODO: refactor helper to accept a pointer to a diags object and append directly into it?
	wrappedObj, d := typehelpers.NewListNestedObjectValueOfPtr(ctx, &result)
	diags.Append(d...)

	return wrappedObj
}

func flattenSubnetServiceDelegationActions(ctx context.Context, input *[]string, diags *diag.Diagnostics) typehelpers.SetValueOf[types.String] {
	if input == nil {
		return typehelpers.NewSetValueOfNull[types.String](ctx)
	}

	return convert.FlattenAndReturn[typehelpers.SetValueOf[types.String]](ctx, *input, diags)
}

func expandSubnetIPAddressPoolFW(ctx context.Context, input typehelpers.ListNestedObjectValueOf[subnetIPAddressPoolModel], diags *diag.Diagnostics) *[]subnets.IPamPoolPrefixAllocation {
	if len(input.Elements()) == 0 {
		return nil
	}

	ipPool := typehelpers.DecodeObjectListOfOne[subnetIPAddressPoolModel](ctx, input, diags)
	return pointer.To([]subnets.IPamPoolPrefixAllocation{
		{
			NumberOfIPAddresses: ipPool.NumberOfIPAddresses.ValueStringPointer(),
			Pool: &subnets.IPamPoolPrefixAllocationPool{
				Id: ipPool.ID.ValueStringPointer(),
			},
		},
	})
}

func flattenSubnetIPAddressPoolFW(ctx context.Context, input *[]subnets.IPamPoolPrefixAllocation, diags *diag.Diagnostics) typehelpers.ListNestedObjectValueOf[subnetIPAddressPoolModel] {
	if input == nil || len(*input) == 0 {
		return typehelpers.NewListNestedObjectValueOfNull[subnetIPAddressPoolModel](ctx)
	}

	value := pointer.From(input)[0]

	result := subnetIPAddressPoolModel{
		NumberOfIPAddresses: types.StringPointerValue(value.NumberOfIPAddresses),
		// convert.Flatten fails if the provided fwObject/type has an empty element type, this seems to happen
		// when the type is defined as `typehelpers.ListValueOf[types.String]`, this ends up as a `basetypes.ListValue` with a `nil` elementType
		// TODO: add a check to convert.Flatten to fail with a clear err msg
		AllocatedIPAddressPrefixes: typehelpers.NewListValueOfNull[types.String](ctx),
	}

	convert.Flatten(ctx, value.AllocatedAddressPrefixes, &result.AllocatedIPAddressPrefixes, diags)

	if value.Pool != nil {
		result.ID = types.StringPointerValue(value.Pool.Id)
	}

	wrappedObj, d := typehelpers.NewListNestedObjectValueOfPtr(ctx, &result)
	diags.Append(d...)

	return wrappedObj
}

func expandSubnetNetworkSecurityGroupID(input types.String) *subnets.NetworkSecurityGroup {
	if input.IsNull() || input.IsUnknown() {
		return nil
	}

	return &subnets.NetworkSecurityGroup{
		Id: input.ValueStringPointer(),
	}
}

func expandSubnetServiceEndpointPolicyIDsFW(ctx context.Context, input types.Set, diags *diag.Diagnostics) *[]subnets.ServiceEndpointPolicy {
	policyIDs := convert.ExpandAndReturn[[]string](ctx, input, diags)

	endpointPolicies := make([]subnets.ServiceEndpointPolicy, 0, len(policyIDs))
	for _, policyID := range policyIDs {
		endpointPolicies = append(endpointPolicies, subnets.ServiceEndpointPolicy{
			Id: pointer.To(policyID),
		})
	}

	return &endpointPolicies
}

func flattenSubnetServiceEndpointPolicyIDsFW(ctx context.Context, input *[]subnets.ServiceEndpointPolicy) (types.Set, diag.Diagnostics) {
	if input == nil {
		return types.SetNull(types.StringType), nil
	}

	output := make([]*string, 0)
	for _, policy := range *input {
		output = append(output, policy.Id)
	}

	return types.SetValueFrom(ctx, types.StringType, output)
}

func expandSubnetServiceEndpointsFW(ctx context.Context, input types.Set, diags *diag.Diagnostics) *[]subnets.ServiceEndpointPropertiesFormat {
	result := make([]subnets.ServiceEndpointPropertiesFormat, 0)

	serviceEndpoints := convert.ExpandAndReturn[[]string](ctx, input, diags)
	for _, serviceEndpoint := range serviceEndpoints {
		result = append(result, subnets.ServiceEndpointPropertiesFormat{
			Service: pointer.To(serviceEndpoint),
		})
	}

	return &result
}

func flattenSubnetServiceEndpointsFW(ctx context.Context, input *[]subnets.ServiceEndpointPropertiesFormat) (types.Set, diag.Diagnostics) {
	if input == nil || len(*input) == 0 {
		return types.SetNull(types.StringType), nil
	}

	output := make([]*string, 0)
	for _, serviceEndpoint := range *input {
		output = append(output, serviceEndpoint.Service)
	}

	return types.SetValueFrom(ctx, types.StringType, output)
}
