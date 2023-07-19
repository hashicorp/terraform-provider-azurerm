package mssql

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2022-02-01/availabilitygrouplisteners"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2022-02-01/sqlvirtualmachinegroups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MsSqlVirtualMachineGroupResource struct{}

type MsSqlVirtualMachineGroupModel struct {
	Name          string `tfschema:"name"`
	ResourceGroup string `tfschema:"resource_group_name"`
	Location      string `tfschema:"location"`

	SqlImageOffer     string                     `tfschema:"sql_image_offer"`
	SqlImageSku       string                     `tfschema:"sql_image_sku"`
	WsfcDomainProfile []helper.WsfcDomainProfile `tfschema:"wsfc_domain_profile"`
	Tags              map[string]string          `tfschema:"tags"`
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

		"wsfc_domain_profile": helper.WsfcDomainProfileSchemaMsSqlVirtualMachineAvailabilityGroup(),

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
					return fmt.Errorf("checking for present of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := sqlvirtualmachinegroups.SqlVirtualMachineGroup{
				Properties: &sqlvirtualmachinegroups.SqlVirtualMachineGroupProperties{
					SqlImageOffer:     utils.String(model.SqlImageOffer),
					SqlImageSku:       pointer.To(sqlvirtualmachinegroups.SqlVMGroupImageSku(model.SqlImageSku)),
					WsfcDomainProfile: expandMsSqlVirtualMachineGroupWsfcDomainProfile(model.WsfcDomainProfile),
				},

				Location: azure.NormalizeLocation(model.Location),
				Tags:     &model.Tags,
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

					sqlImageOffer := ""
					if props.SqlImageOffer != nil {
						sqlImageOffer = *props.SqlImageOffer
					}
					state.SqlImageOffer = sqlImageOffer

					sqlImageSku := ""
					if props.SqlImageSku != nil {
						sqlImageSku = string(*props.SqlImageSku)
					}
					state.SqlImageSku = sqlImageSku

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
				state.Location = azure.NormalizeLocation(model.Location)

				if model.Tags != nil {
					state.Tags = *model.Tags
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
					SqlImageOffer:     utils.String(model.SqlImageOffer),
					SqlImageSku:       pointer.To(sqlvirtualmachinegroups.SqlVMGroupImageSku(model.SqlImageSku)),
					WsfcDomainProfile: expandMsSqlVirtualMachineGroupWsfcDomainProfile(model.WsfcDomainProfile),
				},

				Location: azure.NormalizeLocation(model.Location),
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

func expandMsSqlVirtualMachineGroupWsfcDomainProfile(wsfcDomainProfile []helper.WsfcDomainProfile) *sqlvirtualmachinegroups.WsfcDomainProfile {
	if wsfcDomainProfile == nil {
		return nil
	}

	result := sqlvirtualmachinegroups.WsfcDomainProfile{
		ClusterSubnetType:        pointer.To(sqlvirtualmachinegroups.ClusterSubnetType(wsfcDomainProfile[0].ClusterSubnetType)),
		DomainFqdn:               pointer.To(wsfcDomainProfile[0].Fqdn),
		OuPath:                   pointer.To(wsfcDomainProfile[0].OuPath),
		ClusterBootstrapAccount:  pointer.To(wsfcDomainProfile[0].ClusterBootstrapAccountName),
		ClusterOperatorAccount:   pointer.To(wsfcDomainProfile[0].ClusterOperatorAccountName),
		SqlServiceAccount:        pointer.To(wsfcDomainProfile[0].SqlServiceAccountName),
		StorageAccountUrl:        pointer.To(wsfcDomainProfile[0].StorageAccountUrl),
		StorageAccountPrimaryKey: pointer.To(wsfcDomainProfile[0].StorageAccountPrimaryKey),
	}

	return &result
}

func flattenMsSqlVirtualMachineGroupWsfcDomainProfile(domainProfile *sqlvirtualmachinegroups.WsfcDomainProfile, storageAccountPrimaryKey string) []helper.WsfcDomainProfile {
	if domainProfile == nil {
		return []helper.WsfcDomainProfile{}
	}

	return []helper.WsfcDomainProfile{
		{
			Fqdn:                        pointer.From(domainProfile.DomainFqdn),
			OuPath:                      pointer.From(domainProfile.OuPath),
			ClusterBootstrapAccountName: pointer.From(domainProfile.ClusterBootstrapAccount),
			ClusterOperatorAccountName:  pointer.From(domainProfile.ClusterOperatorAccount),
			SqlServiceAccountName:       pointer.From(domainProfile.SqlServiceAccount),
			StorageAccountUrl:           pointer.From(domainProfile.StorageAccountUrl),
			ClusterSubnetType:           string(pointer.From(domainProfile.ClusterSubnetType)),
			StorageAccountPrimaryKey:    storageAccountPrimaryKey,
		},
	}
}
