// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2023-10-01/availabilitygrouplisteners"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2023-10-01/sqlvirtualmachinegroups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type MsSqlVirtualMachineGroupResource struct{}

type MsSqlVirtualMachineGroupModel struct {
	Name          string `tfschema:"name"`
	ResourceGroup string `tfschema:"resource_group_name"`
	Location      string `tfschema:"location"`

	SqlImageOffer     string              `tfschema:"sql_image_offer"`
	SqlImageSku       string              `tfschema:"sql_image_sku"`
	WsfcDomainProfile []WsfcDomainProfile `tfschema:"wsfc_domain_profile"`
	Tags              map[string]string   `tfschema:"tags"`
}

type WsfcDomainProfile struct {
	Fqdn                        string `tfschema:"fqdn"`
	OrganizationalUnitPath      string `tfschema:"organizational_unit_path"`
	ClusterBootstrapAccountName string `tfschema:"cluster_bootstrap_account_name"`
	ClusterOperatorAccountName  string `tfschema:"cluster_operator_account_name"`
	SqlServiceAccountName       string `tfschema:"sql_service_account_name"`
	StorageAccountUrl           string `tfschema:"storage_account_url"`
	StorageAccountPrimaryKey    string `tfschema:"storage_account_primary_key"`
	ClusterSubnetType           string `tfschema:"cluster_subnet_type"`
}

var _ sdk.Resource = MsSqlVirtualMachineGroupResource{}
var _ sdk.ResourceWithUpdate = MsSqlVirtualMachineGroupResource{}

func (r MsSqlVirtualMachineGroupResource) ModelObject() interface{} {
	return &MsSqlVirtualMachineGroupModel{}
}

func (r MsSqlVirtualMachineGroupResource) ResourceType() string {
	return "azurerm_mssql_virtual_machine_group"
}

func (r MsSqlVirtualMachineGroupResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return availabilitygrouplisteners.ValidateSqlVirtualMachineGroupID
}

func (r MsSqlVirtualMachineGroupResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 15),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"sql_image_offer": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SqlImageOfferName,
		},

		"sql_image_sku": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(sqlvirtualmachinegroups.SqlVMGroupImageSkuDeveloper),
				string(sqlvirtualmachinegroups.SqlVMGroupImageSkuEnterprise),
			}, false),
		},

		"wsfc_domain_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"cluster_subnet_type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(sqlvirtualmachinegroups.ClusterSubnetTypeMultiSubnet),
							string(sqlvirtualmachinegroups.ClusterSubnetTypeSingleSubnet),
						}, false),
					},

					"fqdn": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"organizational_unit_path": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"cluster_bootstrap_account_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"cluster_operator_account_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"sql_service_account_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"storage_account_url": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"storage_account_primary_key": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"tags": tags.Schema(),
	}
}

func (r MsSqlVirtualMachineGroupResource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r MsSqlVirtualMachineGroupResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model MsSqlVirtualMachineGroupModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			client := metadata.Client.MSSQL.VirtualMachineGroupsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := sqlvirtualmachinegroups.NewSqlVirtualMachineGroupID(subscriptionId, model.ResourceGroup, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := sqlvirtualmachinegroups.SqlVirtualMachineGroup{
				Properties: &sqlvirtualmachinegroups.SqlVirtualMachineGroupProperties{
					SqlImageOffer:     pointer.To(model.SqlImageOffer),
					SqlImageSku:       pointer.To(sqlvirtualmachinegroups.SqlVMGroupImageSku(model.SqlImageSku)),
					WsfcDomainProfile: expandMsSqlVirtualMachineGroupWsfcDomainProfile(model.WsfcDomainProfile),
				},

				Location: location.Normalize(model.Location),
				Tags:     pointer.To(model.Tags),
			}

			if err = client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r MsSqlVirtualMachineGroupResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {

			client := metadata.Client.MSSQL.VirtualMachineGroupsClient

			id, err := sqlvirtualmachinegroups.ParseSqlVirtualMachineGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			state := MsSqlVirtualMachineGroupModel{
				Name:          id.SqlVirtualMachineGroupName,
				ResourceGroup: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {

					state.SqlImageOffer = pointer.From(props.SqlImageOffer)
					state.SqlImageSku = string(pointer.From(props.SqlImageSku))

					var oldModel MsSqlVirtualMachineGroupModel
					if err = metadata.Decode(&oldModel); err != nil {
						return err
					}
					storageAccountPrimaryKey := ""
					if oldModel.WsfcDomainProfile != nil && len(oldModel.WsfcDomainProfile) != 0 {
						storageAccountPrimaryKey = oldModel.WsfcDomainProfile[0].StorageAccountPrimaryKey
					}
					state.WsfcDomainProfile = flattenMsSqlVirtualMachineGroupWsfcDomainProfile(props.WsfcDomainProfile, storageAccountPrimaryKey)

				}
				state.Location = location.Normalize(model.Location)

				if model.Tags != nil {
					state.Tags = pointer.From(model.Tags)
				}
			}
			return metadata.Encode(&state)
		},
	}
}

func (r MsSqlVirtualMachineGroupResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model MsSqlVirtualMachineGroupModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			client := metadata.Client.MSSQL.VirtualMachineGroupsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := sqlvirtualmachinegroups.NewSqlVirtualMachineGroupID(subscriptionId, model.ResourceGroup, model.Name)

			_, err := client.Get(ctx, id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			parameters := sqlvirtualmachinegroups.SqlVirtualMachineGroup{
				Properties: &sqlvirtualmachinegroups.SqlVirtualMachineGroupProperties{
					SqlImageOffer:     pointer.To(model.SqlImageOffer),
					SqlImageSku:       pointer.To(sqlvirtualmachinegroups.SqlVMGroupImageSku(model.SqlImageSku)),
					WsfcDomainProfile: expandMsSqlVirtualMachineGroupWsfcDomainProfile(model.WsfcDomainProfile),
				},

				Location: location.Normalize(model.Location),
				Tags:     &model.Tags,
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r MsSqlVirtualMachineGroupResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.VirtualMachineGroupsClient

			id, err := sqlvirtualmachinegroups.ParseSqlVirtualMachineGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err = client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
			return nil
		},
	}
}

func expandMsSqlVirtualMachineGroupWsfcDomainProfile(wsfcDomainProfile []WsfcDomainProfile) *sqlvirtualmachinegroups.WsfcDomainProfile {
	if wsfcDomainProfile == nil {
		return nil
	}

	result := sqlvirtualmachinegroups.WsfcDomainProfile{
		ClusterSubnetType:        pointer.To(sqlvirtualmachinegroups.ClusterSubnetType(wsfcDomainProfile[0].ClusterSubnetType)),
		DomainFqdn:               pointer.To(wsfcDomainProfile[0].Fqdn),
		OuPath:                   pointer.To(wsfcDomainProfile[0].OrganizationalUnitPath),
		ClusterBootstrapAccount:  pointer.To(wsfcDomainProfile[0].ClusterBootstrapAccountName),
		ClusterOperatorAccount:   pointer.To(wsfcDomainProfile[0].ClusterOperatorAccountName),
		SqlServiceAccount:        pointer.To(wsfcDomainProfile[0].SqlServiceAccountName),
		StorageAccountURL:        pointer.To(wsfcDomainProfile[0].StorageAccountUrl),
		StorageAccountPrimaryKey: pointer.To(wsfcDomainProfile[0].StorageAccountPrimaryKey),
	}

	return &result
}

func flattenMsSqlVirtualMachineGroupWsfcDomainProfile(domainProfile *sqlvirtualmachinegroups.WsfcDomainProfile, storageAccountPrimaryKey string) []WsfcDomainProfile {
	if domainProfile == nil {
		return []WsfcDomainProfile{}
	}

	return []WsfcDomainProfile{
		{
			Fqdn:                        pointer.From(domainProfile.DomainFqdn),
			OrganizationalUnitPath:      pointer.From(domainProfile.OuPath),
			ClusterBootstrapAccountName: pointer.From(domainProfile.ClusterBootstrapAccount),
			ClusterOperatorAccountName:  pointer.From(domainProfile.ClusterOperatorAccount),
			SqlServiceAccountName:       pointer.From(domainProfile.SqlServiceAccount),
			StorageAccountUrl:           pointer.From(domainProfile.StorageAccountURL),
			ClusterSubnetType:           string(pointer.From(domainProfile.ClusterSubnetType)),
			StorageAccountPrimaryKey:    storageAccountPrimaryKey,
		},
	}
}
