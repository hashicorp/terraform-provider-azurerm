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
	customplanmodifier "github.com/hashicorp/terraform-provider-azurerm/internal/sdk/planmodifier"
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

			"network_security_group_id": schema.StringAttribute{
				Optional: true,
				// Computed: true,
				// TODO: O+C does not work, planmodifier?
				PlanModifiers: []planmodifier.String{
					customplanmodifier.SuppressIfPrivateStateNilOrDoesNotEqual(privateStateNetworkSecurityGroupID, "true", "false"),
				},
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: networksecuritygroups.ValidateNetworkSecurityGroupID,
					},
				},
			},
		},

		Blocks: map[string]schema.Block{
			// TODO: shouldn't matter for this one, but why can't we set Optional/Required?
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
	diags := pointer.To(resp.Diagnostics) // TODO: init/ensure not nil?

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
			PrivateEndpointNetworkPolicies:    pointer.ToEnum[subnets.VirtualNetworkPrivateEndpointNetworkPolicies](data.PrivateEndpointNetworkPolicies.ValueString()), // TODO: are there easier ways?/helpers we can write?
			PrivateLinkServiceNetworkPolicies: expandSubnetNetworkPolicyFW(data.PrivateLinkServiceNetworkPoliciesEnabled),
			ServiceEndpointPolicies:           expandSubnetServiceEndpointPolicyIDsFW(ctx, data.ServiceEndpointPolicyIds, diags),
			ServiceEndpoints:                  expandSubnetServiceEndpointsFW(ctx, data.ServiceEndpoints, diags),
			SharingScope:                      pointer.ToEnum[subnets.SharingScope](data.SharingScope.ValueString()),
		},
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Private.SetKey(ctx, privateStateNetworkSecurityGroupID, sdk.NewPrivateStateValue("false").Bytes())
	if !data.NetworkSecurityGroupID.IsNull() && !data.NetworkSecurityGroupID.IsUnknown() {
		// If `network_security_group_id` is not null, we'll add this to private state to ensure we can track removal
		// even though this property is O+C and may be set outside of this resource.
		resp.Private.SetKey(ctx, privateStateNetworkSecurityGroupID, []byte(`{"value": "true"}`))
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

	// At this point, we need to confirm whether the NSG was set inline.
	// If it was, do nothing
	// If it wasn't, remove it from the read data to prevent a diff
	privateNSG, d := req.Private.GetKey(ctx, privateStateNetworkSecurityGroupID)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// If nil, this is likely an existing resource created prior to FW migration
	// So we'll set the "default" private state normally set by Create
	// as well as assigning it to the private state val retrieved earlier to ensure this functions on pre-existing resources
	if privateNSG == nil {
		// should never be nil, error if it is
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("private state nil", "oh noooo"))
		return
		//newPrivateStateVal := NewPrivateStateValue("false").Bytes()
		//// TODO: wrap this?
		//d := resp.Private.SetKey(ctx, privateStateNetworkSecurityGroupID, newPrivateStateVal)
		//if d.HasError() {
		//	resp.Diagnostics.Append(d...)
		//	return
		//}
		//privateNSG = NewPrivateStateValue("false").Bytes()
	}

	resp.Diagnostics.Append(diag.NewWarningDiagnostic("Value of private state in Read", string(privateNSG)))

	// This seems to work in all cases except when a user modifies nsg outside of TF
	// then adds the same NSG to the resource config. Acceptable?
	if sdk.NewPrivateStateValue("false").Equals(privateNSG, &resp.Diagnostics) {
		data.NetworkSecurityGroupID = types.StringNull()
	}
}

func (r SubnetFWResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, metadata sdk.ResourceMetadata, decodedPlan any, decodedState any) {
	// TODO: fix
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
			// TODO: test whether this works with `*[]string`
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

	// TODO: wrap this so it returns the value (or the "zero"/unknown (?) value) and appends diags into ptr to diags
	rawNsgID := types.StringNull() // TODO: does this need to be a TF type or is a simple string enough?
	d := req.Config.GetAttribute(ctx, path.Root("network_security_group_id"), &rawNsgID)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// If NSG is a non-null, non-unknown value, record this in private state
	resp.Private.SetKey(ctx, privateStateNetworkSecurityGroupID, sdk.NewPrivateStateValue("false").Bytes())
	if !rawNsgID.IsNull() && !rawNsgID.IsUnknown() {
		resp.Private.SetKey(ctx, privateStateNetworkSecurityGroupID, sdk.NewPrivateStateValue("true").Bytes())
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
			state.PrivateEndpointNetworkPolicies = typehelpers.StringFromEnum(props.PrivateEndpointNetworkPolicies)
			state.PrivateLinkServiceNetworkPoliciesEnabled = flattenSubnetNetworkPolicyFW(props.PrivateLinkServiceNetworkPolicies)
			state.SharingScope = typehelpers.StringFromEnum(props.SharingScope)

			serviceEndpointPolicyIDs, d := flattenSubnetServiceEndpointPolicyIDsFW(ctx, props.ServiceEndpointPolicies)
			if d.HasError() {
				// TODO: do we care to wrap diags? if yes, add to all flatten funcs, if no, remove below:
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

	// TODO: refactor these helpers to accept a pointer to a diags object
	wrappedObj, d := typehelpers.NewListNestedObjectValueOfValueSlice[subnetDelegationModel](ctx, result)
	diags.Append(d...)

	return wrappedObj
}

func flattenSubnetServiceDelegation(ctx context.Context, input *subnets.ServiceDelegationPropertiesFormat, diags *diag.Diagnostics) typehelpers.ListNestedObjectValueOf[subnetServiceDelegationModel] {
	if input == nil {
		return typehelpers.NewListNestedObjectValueOfNull[subnetServiceDelegationModel](ctx)
	}

	// TODO: StateFunc or equivalent for FW?
	normalizedServiceName := map[string]string{}
	for _, n := range subnetDelegationServiceNames {
		normalizedServiceName[strings.ToLower(n)] = n
	}

	result := subnetServiceDelegationModel{
		Name:    types.StringPointerValue(input.ServiceName),
		Actions: flattenSubnetServiceDelegationActions(ctx, input.Actions, diags),
	}

	// Normalizing to match existing `azurerm_subnet` logic, we may want to remove this in the future, TODO
	if normalizedName, ok := normalizedServiceName[strings.ToLower(pointer.From(input.ServiceName))]; ok {
		result.Name = types.StringValue(normalizedName)
	}

	// TODO: refactor if changing `NewListNestedObjectValueOfPtr` to accept diags ptr to append into instead of returning new diags
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
		// Can this become a one-liner with instantiate of unknown/null and `convert.Flatten`?
		AllocatedIPAddressPrefixes: typehelpers.NewListValueOfNull[types.String](ctx),
	}

	// convert.Flatten fails if the provided fwObject/type has an empty element type, this can happen
	// when the type is defined as `typehelpers.ListValueOf[types.String]`, the zero-value for this is a `basetypes.ListValue` with a `nil` elementType
	convert.Flatten(ctx, value.AllocatedAddressPrefixes, &result.AllocatedIPAddressPrefixes, diags)

	//prefixes := make([]types.String, 0)
	//for _, prefix := range pointer.From(value.AllocatedAddressPrefixes) {
	//	prefixes = append(prefixes, types.StringValue(prefix))
	//}

	if value.Pool != nil {
		result.ID = types.StringPointerValue(value.Pool.Id)
	}

	// TODO: refactor helpers to take ptr to diags?
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
