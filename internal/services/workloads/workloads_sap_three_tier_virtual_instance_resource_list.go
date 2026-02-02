// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package workloads

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2024-09-01/sapvirtualinstances"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type WorkloadsSAPThreeTierVirtualInstanceListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(WorkloadsSAPThreeTierVirtualInstanceListResource)

func (r WorkloadsSAPThreeTierVirtualInstanceListResource) ResourceFunc() *pluginsdk.Resource {
	// Wrap the SDK v2 typed resource and convert it to pluginsdk.Resource
	wrapper := sdk.NewResourceWrapper(WorkloadsSAPThreeTierVirtualInstanceResource{})
	resource, err := wrapper.Resource()
	if err != nil {
		panic(fmt.Sprintf("failed to wrap resource: %+v", err))
	}
	return resource
}

func (r WorkloadsSAPThreeTierVirtualInstanceListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = "azurerm_workloads_sap_three_tier_virtual_instance"
}

func (r WorkloadsSAPThreeTierVirtualInstanceListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Workloads.SAPVirtualInstances

	var data sdk.DefaultListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	results := make([]sapvirtualinstances.SAPVirtualInstance, 0)

	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	switch {
	case !data.ResourceGroupName.IsNull():
		resp, err := client.ListByResourceGroupComplete(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", "azurerm_workloads_sap_three_tier_virtual_instance"), err)
			return
		}
		results = resp.Items
	default:
		resp, err := client.ListBySubscriptionComplete(ctx, commonids.NewSubscriptionID(subscriptionID))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", "azurerm_workloads_sap_three_tier_virtual_instance"), err)
			return
		}
		results = resp.Items
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, instance := range results {
			result := request.NewListResult(ctx)

			result.DisplayName = pointer.From(instance.Name)

			wrapper := sdk.NewResourceWrapper(WorkloadsSAPThreeTierVirtualInstanceResource{})
			resource, err := wrapper.Resource()
			if err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "creating resource wrapper", err)
				return
			}
			rd := resource.Data(&terraform.InstanceState{})

			id, err := sapvirtualinstances.ParseSapVirtualInstanceID(pointer.From(instance.Id))
			if err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "parsing SAP Virtual Instance ID", err)
				return
			}
			rd.SetId(id.ID())

			if err := flattenForListResource(rd, id, &instance, subscriptionID); err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, fmt.Sprintf("flattening `%s` resource data", "azurerm_workloads_sap_three_tier_virtual_instance"), err)
				return
			}

			tfTypeIdentity, err := rd.TfTypeIdentityState()
			if err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "converting Identity State", err)
				return
			}
			if err := result.Identity.Set(ctx, *tfTypeIdentity); err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "setting Identity Data", err)
				return
			}

			tfTypeResourceState, err := rd.TfTypeResourceState()
			if err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "converting Resource State", err)
				return
			}
			if err := result.Resource.Set(ctx, *tfTypeResourceState); err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "setting Resource Data", err)
				return
			}

			if !push(result) {
				return
			}
		}
	}
}

// flattenForListResource populates pluginsdk.ResourceData directly for list resources
// This is separate from the Read() function which uses metadata.Encode() with typed models
func flattenForListResource(d *pluginsdk.ResourceData, id *sapvirtualinstances.SapVirtualInstanceId, model *sapvirtualinstances.SAPVirtualInstance, subscriptionId string) error {
	if err := d.Set("name", id.SapVirtualInstanceName); err != nil {
		return fmt.Errorf("setting `name`: %+v", err)
	}
	if err := d.Set("resource_group_name", id.ResourceGroupName); err != nil {
		return fmt.Errorf("setting `resource_group_name`: %+v", err)
	}

	if model == nil {
		return nil
	}

	if err := d.Set("location", location.Normalize(model.Location)); err != nil {
		return fmt.Errorf("setting `location`: %+v", err)
	}

	identityFlattened, err := identity.FlattenUserAssignedMap(model.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}
	if err := d.Set("identity", identityFlattened); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	if props := model.Properties; props != nil {
		if err := d.Set("environment", string(props.Environment)); err != nil {
			return fmt.Errorf("setting `environment`: %+v", err)
		}

		if err := d.Set("managed_resources_network_access_type", string(pointer.From(props.ManagedResourcesNetworkAccessType))); err != nil {
			return fmt.Errorf("setting `managed_resources_network_access_type`: %+v", err)
		}

		if err := d.Set("sap_product", string(props.SapProduct)); err != nil {
			return fmt.Errorf("setting `sap_product`: %+v", err)
		}

		if err := d.Set("tags", pointer.From(model.Tags)); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}

		if config := props.Configuration; config != nil {
			if v, ok := config.(sapvirtualinstances.DeploymentWithOSConfiguration); ok {
				// Set app location
				if err := d.Set("app_location", location.Normalize(pointer.From(v.AppLocation))); err != nil {
					return fmt.Errorf("setting `app_location`: %+v", err)
				}

				if osSapConfiguration := v.OsSapConfiguration; osSapConfiguration != nil {
					if err := d.Set("sap_fqdn", pointer.From(osSapConfiguration.SapFqdn)); err != nil {
						return fmt.Errorf("setting `sap_fqdn`: %+v", err)
					}
				}

				if configuration := v.InfrastructureConfiguration; configuration != nil {
					if threeTierConfiguration, ok := configuration.(sapvirtualinstances.ThreeTierConfiguration); ok {
						threeTierConfig, err := flattenThreeTierConfiguration(threeTierConfiguration, d, subscriptionId)
						if err != nil {
							return fmt.Errorf("flattening `three_tier_configuration`: %+v", err)
						}
						// Convert the struct slice to []interface{} for use with pluginsdk.ResourceData.Set()
						// Using SDK helper since pluginsdk cannot handle custom structs directly
						threeTierConfigForList := sdk.PluginSDKFlattenStructSliceToInterface(threeTierConfig, nil)
						if err := d.Set("three_tier_configuration", threeTierConfigForList); err != nil {
							return fmt.Errorf("setting `three_tier_configuration`: %+v", err)
						}
					}
				}
			}
		}

		if v := props.ManagedResourceGroupConfiguration; v != nil {
			if err := d.Set("managed_resource_group_name", pointer.From(v.Name)); err != nil {
				return fmt.Errorf("setting `managed_resource_group_name`: %+v", err)
			}
		}
	}

	return nil
}
