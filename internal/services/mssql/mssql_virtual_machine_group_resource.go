package mssql

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2022-02-01/availabilitygrouplisteners"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2022-02-01/sqlvirtualmachinegroups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/helper"
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
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"SQL2016-WS2016",
				"SQL2017-WS2016",
				"SQL2019-WS2019",
			}, false),
		},

		"sql_image_sku": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(sqlvirtualmachinegroups.SqlVmGroupImageSkuDeveloper),
				string(sqlvirtualmachinegroups.SqlVmGroupImageSkuEnterprise),
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

			sqlVmGroupImageSku := sqlvirtualmachinegroups.SqlVmGroupImageSku(model.SqlImageSku)

			parameters := sqlvirtualmachinegroups.SqlVirtualMachineGroup{
				Properties: &sqlvirtualmachinegroups.SqlVirtualMachineGroupProperties{
					SqlImageOffer:     utils.String(model.SqlImageOffer),
					SqlImageSku:       &sqlVmGroupImageSku,
					WsfcDomainProfile: expandMsSqlVirtualMachineGroupWsfcDomainProfile(model.WsfcDomainProfile),
				},

				Location: azure.NormalizeLocation(model.Location),
				Tags:     &model.Tags,
			}

			parameters.Properties.WsfcDomainProfile.FileShareWitnessPath = utils.String("")

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

					state.WsfcDomainProfile = flattenMsSqlVirtualMachineGroupWsfcDomainProfile(props.WsfcDomainProfile)

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

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for present of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			sqlVmGroupImageSku := sqlvirtualmachinegroups.SqlVmGroupImageSku(model.SqlImageSku)

			parameters := sqlvirtualmachinegroups.SqlVirtualMachineGroup{
				Properties: &sqlvirtualmachinegroups.SqlVirtualMachineGroupProperties{
					SqlImageOffer:     utils.String(model.SqlImageOffer),
					SqlImageSku:       &sqlVmGroupImageSku,
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
		DomainFqdn:               &wsfcDomainProfile[0].Fqdn,
		OuPath:                   &wsfcDomainProfile[0].OuPath,
		ClusterBootstrapAccount:  &wsfcDomainProfile[0].ClusterBootstrapAccountName,
		ClusterOperatorAccount:   &wsfcDomainProfile[0].ClusterOperatorAccountName,
		SqlServiceAccount:        &wsfcDomainProfile[0].SqlServiceAccountName,
		FileShareWitnessPath:     &wsfcDomainProfile[0].FileShareWitnessPath,
		StorageAccountUrl:        &wsfcDomainProfile[0].StorageAccountUrl,
		StorageAccountPrimaryKey: &wsfcDomainProfile[0].StorageAccountPrimaryKey,
	}

	clusterSubnetType := sqlvirtualmachinegroups.ClusterSubnetType(wsfcDomainProfile[0].ClusterSubnetType)
	result.ClusterSubnetType = &clusterSubnetType

	return &result
}

func flattenMsSqlVirtualMachineGroupWsfcDomainProfile(domainProfile *sqlvirtualmachinegroups.WsfcDomainProfile) []helper.WsfcDomainProfile {
	if domainProfile == nil {
		return []helper.WsfcDomainProfile{}
	}

	var fqdn string
	if domainProfile.DomainFqdn != nil {
		fqdn = *domainProfile.DomainFqdn
	}

	var clusterOperatorAccountName string
	if domainProfile.ClusterOperatorAccount != nil {
		clusterOperatorAccountName = *domainProfile.ClusterOperatorAccount
	}

	var clusterBootstrapAccountName string
	if domainProfile.ClusterBootstrapAccount != nil {
		clusterBootstrapAccountName = *domainProfile.ClusterBootstrapAccount
	}

	var ouPath string
	if domainProfile.OuPath != nil {
		ouPath = *domainProfile.OuPath
	}

	var sqlServiceAccountName string
	if domainProfile.SqlServiceAccount != nil {
		sqlServiceAccountName = *domainProfile.SqlServiceAccount
	}

	var fileShareWitnessPath string
	if domainProfile.FileShareWitnessPath != nil {
		fileShareWitnessPath = *domainProfile.FileShareWitnessPath
	}

	var storageAccountURL string
	if domainProfile.StorageAccountUrl != nil {
		storageAccountURL = *domainProfile.StorageAccountUrl
	}

	var clusterSubnetType string
	if domainProfile.ClusterSubnetType != nil {
		clusterSubnetType = string(*domainProfile.ClusterSubnetType)
	}

	return []helper.WsfcDomainProfile{
		{
			Fqdn:                        fqdn,
			OuPath:                      ouPath,
			ClusterBootstrapAccountName: clusterBootstrapAccountName,
			ClusterOperatorAccountName:  clusterOperatorAccountName,
			SqlServiceAccountName:       sqlServiceAccountName,
			FileShareWitnessPath:        fileShareWitnessPath,
			StorageAccountUrl:           storageAccountURL,
			ClusterSubnetType:           clusterSubnetType,
		},
	}
}
