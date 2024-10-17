// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracle

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/autonomousdatabases"
	"github.com/hashicorp/terraform-provider-azurerm/utils"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = AutonomousDatabaseRegularResource{}

type AutonomousDatabaseRegularResource struct{}

type AutonomousDatabaseRegularResourceModel struct {
	// Azure
	Location          string                 `tfschema:"location"`
	Name              string                 `tfschema:"name"`
	ResourceGroupName string                 `tfschema:"resource_group_name"`
	Tags              map[string]interface{} `tfschema:"tags"`

	// Required
	AdminPassword                  string  `tfschema:"admin_password"`
	BackupRetentionPeriodInDays    int64   `tfschema:"backup_retention_period_in_days"`
	CharacterSet                   string  `tfschema:"character_set"`
	ComputeCount                   float64 `tfschema:"compute_count"`
	ComputeModel                   string  `tfschema:"compute_model"`
	DataStorageSizeInGbs           int64   `tfschema:"data_storage_size_in_gbs"`
	DbVersion                      string  `tfschema:"db_version"`
	DbWorkload                     string  `tfschema:"db_workload"`
	DisplayName                    string  `tfschema:"display_name"`
	LicenseModel                   string  `tfschema:"license_model"`
	IsAutoScalingEnabled           bool    `tfschema:"is_auto_scaling_enabled"`
	IsAutoScalingForStorageEnabled bool    `tfschema:"is_auto_scaling_for_storage_enabled"`
	IsMtlsConnectionRequired       bool    `tfschema:"is_mtls_connection_required"`
	NcharacterSet                  string  `tfschema:"ncharacter_set"`
	SubnetId                       string  `tfschema:"subnet_id"`
	VnetId                         string  `tfschema:"vnet_id"`

	// Optional
	CustomerContacts []string `tfschema:"customer_contacts"`
}

func (AutonomousDatabaseRegularResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		// Azure
		"location": commonschema.Location(),
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"tags": commonschema.Tags(),

		// Required
		"admin_password": {
			Type:      pluginsdk.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"backup_retention_period_in_days": {
			Type:     pluginsdk.TypeInt,
			Required: true,
		},
		"character_set": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"compute_count": {
			Type:     pluginsdk.TypeFloat,
			Required: true,
		},
		"compute_model": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"data_storage_size_in_gbs": {
			Type:     pluginsdk.TypeInt,
			Required: true,
		},
		"db_version": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"db_workload": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"display_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"is_auto_scaling_enabled": {
			Type:     pluginsdk.TypeBool,
			Required: true,
		},
		"is_auto_scaling_for_storage_enabled": {
			Type:     pluginsdk.TypeBool,
			Required: true,
		},
		"is_mtls_connection_required": {
			Type:     pluginsdk.TypeBool,
			Required: true,
		},
		"license_model": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"ncharacter_set": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"subnet_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"vnet_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		// Optional
		"customer_contacts": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (AutonomousDatabaseRegularResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (AutonomousDatabaseRegularResource) ModelObject() interface{} {
	return &AutonomousDatabaseRegularResource{}
}

func (AutonomousDatabaseRegularResource) ResourceType() string {
	return "azurerm_oracle_autonomous_database"
}

func convertAdbsCustomerContactsToSDK(customerContactsList []string) []autonomousdatabases.CustomerContact {
	var customerContacts []autonomousdatabases.CustomerContact
	if customerContactsList != nil {
		for _, customerContact := range customerContactsList {
			customerContacts = append(customerContacts, autonomousdatabases.CustomerContact{
				Email: customerContact,
			})
		}
	}
	return customerContacts
}

func convertAdbsCustomerContactsToInternalModel(customerContactsList *[]autonomousdatabases.CustomerContact) []string {
	var customerContacts []string
	if customerContactsList != nil {
		for _, customerContact := range *customerContactsList {
			customerContacts = append(customerContacts, customerContact.Email)
		}
	}
	return customerContacts
}

func (r AutonomousDatabaseRegularResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.AutonomousDatabases
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model AutonomousDatabaseRegularResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id := autonomousdatabases.NewAutonomousDatabaseID(subscriptionId,
				model.ResourceGroupName,
				model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			param := autonomousdatabases.AutonomousDatabase{
				Name:     pointer.To(model.Name),
				Location: model.Location,
				Tags:     tags.Expand(model.Tags),
				Properties: &autonomousdatabases.AutonomousDatabaseProperties{
					AdminPassword:                  pointer.To(model.AdminPassword),
					BackupRetentionPeriodInDays:    pointer.To(model.BackupRetentionPeriodInDays),
					CharacterSet:                   pointer.To(model.CharacterSet),
					ComputeCount:                   pointer.To(model.ComputeCount),
					ComputeModel:                   pointer.To(autonomousdatabases.ComputeModel(model.ComputeModel)),
					CustomerContacts:               pointer.To(convertAdbsCustomerContactsToSDK(model.CustomerContacts)),
					DataStorageSizeInGbs:           pointer.To(model.DataStorageSizeInGbs),
					DbWorkload:                     pointer.To(autonomousdatabases.WorkloadType(model.DbWorkload)),
					DbVersion:                      pointer.To(model.DbVersion),
					DisplayName:                    pointer.To(model.DisplayName),
					IsAutoScalingEnabled:           pointer.To(model.IsAutoScalingEnabled),
					IsAutoScalingForStorageEnabled: pointer.To(model.IsAutoScalingForStorageEnabled),
					IsMtlsConnectionRequired:       pointer.To(model.IsMtlsConnectionRequired),
					LicenseModel:                   pointer.To(autonomousdatabases.LicenseModel(model.LicenseModel)),
					NcharacterSet:                  pointer.To(model.NcharacterSet),
					SubnetId:                       pointer.To(model.SubnetId),
					VnetId:                         pointer.To(model.VnetId),
				},
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r AutonomousDatabaseRegularResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {

			client := metadata.Client.Oracle.OracleClient.AutonomousDatabases
			id, err := autonomousdatabases.ParseAutonomousDatabaseID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model AutonomousDatabaseRegularResourceModel
			if err = metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding err: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving exists when updating: +%v", *id)
			}
			if existing.Model == nil && existing.Model.Properties == nil {
				return fmt.Errorf("retrieving as nil when updating for %v", *id)
			}

			if metadata.ResourceData.HasChange("tags") {
				update := &autonomousdatabases.AutonomousDatabaseUpdate{
					Tags: tags.Expand(model.Tags),
				}
				err = client.UpdateThenPoll(ctx, *id, *update)
				if err != nil {
					return fmt.Errorf("updating %s: %v", id, err)
				}
			} else if metadata.ResourceData.HasChangesExcept("tags") {
				return fmt.Errorf("only `tags` currently support updates")
			}
			return nil
		},
	}
}

func (AutonomousDatabaseRegularResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := autonomousdatabases.ParseAutonomousDatabaseID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.Oracle.OracleClient.AutonomousDatabases
			result, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(result.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return err
			}

			if result.Model == nil {
				return fmt.Errorf("retrieving %s got nil model", id)
			}
			prop := result.Model.Properties
			switch adbsPropModel := prop.(type) {
			case autonomousdatabases.AutonomousDatabaseProperties:
				var output AutonomousDatabaseRegularResourceModel
				output.AdminPassword = pointer.From(adbsPropModel.AdminPassword)
				output.BackupRetentionPeriodInDays = pointer.From(adbsPropModel.BackupRetentionPeriodInDays)
				output.CharacterSet = pointer.From(adbsPropModel.CharacterSet)
				output.ComputeCount = pointer.From(adbsPropModel.ComputeCount)
				output.ComputeModel = string(pointer.From(adbsPropModel.ComputeModel))
				output.CustomerContacts = convertAdbsCustomerContactsToInternalModel(adbsPropModel.CustomerContacts)
				output.DataStorageSizeInGbs = pointer.From(adbsPropModel.DataStorageSizeInGbs)
				output.DbWorkload = string(pointer.From(adbsPropModel.DbWorkload))
				output.DbVersion = pointer.From(adbsPropModel.DbVersion)
				output.DisplayName = pointer.From(adbsPropModel.DisplayName)
				output.IsAutoScalingEnabled = pointer.From(adbsPropModel.IsAutoScalingEnabled)
				output.IsAutoScalingForStorageEnabled = pointer.From(adbsPropModel.IsAutoScalingForStorageEnabled)
				output.Location = result.Model.Location
				output.Name = pointer.ToString(result.Model.Name)
				output.LicenseModel = string(pointer.From(adbsPropModel.LicenseModel))
				output.NcharacterSet = pointer.From(adbsPropModel.NcharacterSet)
				output.ResourceGroupName = id.ResourceGroupName
				output.SubnetId = pointer.From(adbsPropModel.SubnetId)
				output.Tags = utils.FlattenPtrMapStringString(result.Model.Tags)
				output.VnetId = pointer.From(adbsPropModel.VnetId)

				return metadata.Encode(&output)
			default:
				return fmt.Errorf("unexpected Autonomous Database type, must be of type Regular")
			}
		},
	}
}

func (AutonomousDatabaseRegularResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.AutonomousDatabases

			id, err := autonomousdatabases.ParseAutonomousDatabaseID(metadata.ResourceData.Id())
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

func (AutonomousDatabaseRegularResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return autonomousdatabases.ValidateAutonomousDatabaseID
}
