// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracle

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/autonomousdatabases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = AutonomousDatabaseRegularResource{}

type AutonomousDatabaseRegularResource struct{}

type AutonomousDatabaseRegularResourceModel struct {
	Location          string            `tfschema:"location"`
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Tags              map[string]string `tfschema:"tags"`

	// Required
	AdminPassword                string  `tfschema:"admin_password"`
	BackupRetentionPeriodInDays  int64   `tfschema:"backup_retention_period_in_days"`
	CharacterSet                 string  `tfschema:"character_set"`
	ComputeCount                 float64 `tfschema:"compute_count"`
	ComputeModel                 string  `tfschema:"compute_model"`
	DataStorageSizeInTbs         int64   `tfschema:"data_storage_size_in_tbs"`
	DbVersion                    string  `tfschema:"db_version"`
	DbWorkload                   string  `tfschema:"db_workload"`
	DisplayName                  string  `tfschema:"display_name"`
	LicenseModel                 string  `tfschema:"license_model"`
	AutoScalingEnabled           bool    `tfschema:"auto_scaling_enabled"`
	AutoScalingForStorageEnabled bool    `tfschema:"auto_scaling_for_storage_enabled"`
	MtlsConnectionRequired       bool    `tfschema:"mtls_connection_required"`
	NcharacterSet                string  `tfschema:"ncharacter_set"`
	SubnetId                     string  `tfschema:"subnet_id"`
	VnetId                       string  `tfschema:"vnet_id"`

	// Optional
	CustomerContacts []string `tfschema:"customer_contacts"`
}

func (AutonomousDatabaseRegularResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.Location(),

		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.AutonomousDatabaseName,
			ForceNew:     true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		// Required
		"admin_password": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			Sensitive:    true,
			ForceNew:     true,
			ValidateFunc: validate.AutonomousDatabasePassword,
		},

		"backup_retention_period_in_days": {
			Type:     pluginsdk.TypeInt,
			Required: true,
			ForceNew: true,
		},

		"character_set": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"compute_count": {
			Type:     pluginsdk.TypeFloat,
			Required: true,
		},

		"compute_model": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"data_storage_size_in_tbs": {
			Type:     pluginsdk.TypeInt,
			Required: true,
		},

		"db_version": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"db_workload": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AutonomousDatabaseName,
		},

		"auto_scaling_enabled": {
			Type:     pluginsdk.TypeBool,
			Required: true,
		},

		"auto_scaling_for_storage_enabled": {
			Type:     pluginsdk.TypeBool,
			Required: true,
		},

		"mtls_connection_required": {
			Type:     pluginsdk.TypeBool,
			Required: true,
			ForceNew: true,
		},

		"license_model": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.LicenseType,
		},

		"ncharacter_set": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"subnet_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"vnet_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		// Optional
		"customer_contacts": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validate.CustomerContactEmail,
			},
		},

		"tags": commonschema.Tags(),
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

func (r AutonomousDatabaseRegularResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 120 * time.Minute,
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
				Tags:     pointer.To(model.Tags),
				Properties: &autonomousdatabases.AutonomousDatabaseProperties{
					AdminPassword:                  pointer.To(model.AdminPassword),
					BackupRetentionPeriodInDays:    pointer.To(model.BackupRetentionPeriodInDays),
					CharacterSet:                   pointer.To(model.CharacterSet),
					ComputeCount:                   pointer.To(model.ComputeCount),
					ComputeModel:                   pointer.To(autonomousdatabases.ComputeModel(model.ComputeModel)),
					CustomerContacts:               pointer.To(expandAdbsCustomerContacts(model.CustomerContacts)),
					DataStorageSizeInGbs:           pointer.To(model.DataStorageSizeInTbs),
					DbWorkload:                     pointer.To(autonomousdatabases.WorkloadType(model.DbWorkload)),
					DbVersion:                      pointer.To(model.DbVersion),
					DisplayName:                    pointer.To(model.DisplayName),
					IsAutoScalingEnabled:           pointer.To(model.AutoScalingEnabled),
					IsAutoScalingForStorageEnabled: pointer.To(model.AutoScalingForStorageEnabled),
					IsMtlsConnectionRequired:       pointer.To(model.MtlsConnectionRequired),
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

			hasUpdate := false
			update := &autonomousdatabases.AutonomousDatabaseUpdate{}
			properties := &autonomousdatabases.AutonomousDatabaseUpdateProperties{}
			update.Properties = properties
			if metadata.ResourceData.HasChange("tags") {
				hasUpdate = true
				update.Tags = pointer.To(model.Tags)
			}
			if metadata.ResourceData.HasChange("data_storage_size_in_tbs") {
				hasUpdate = true
				update.Properties.DataStorageSizeInTbs = pointer.To(model.DataStorageSizeInTbs)
			}
			if metadata.ResourceData.HasChange("compute_count") {
				hasUpdate = true
				update.Properties.ComputeCount = pointer.To(model.ComputeCount)
			}
			if metadata.ResourceData.HasChange("auto_scaling_enabled") {
				hasUpdate = true
				update.Properties.IsAutoScalingEnabled = pointer.To(model.AutoScalingEnabled)
			}
			if metadata.ResourceData.HasChange("auto_scaling_for_storage_enabled") {
				hasUpdate = true
				update.Properties.IsAutoScalingForStorageEnabled = pointer.To(model.AutoScalingForStorageEnabled)
			} else if metadata.ResourceData.HasChangesExcept("auto_scaling_enabled", "auto_scaling_for_storage_enabled",
				"compute_count", "data_storage_size_in_tbs", "tags") {
				return fmt.Errorf("only `auto_scaling_enabled`, `auto_scaling_for_storage_enabled`, `compute_count`, " +
					"`data_storage_size_in_tbs`, and `tags` currently support updates")
			}

			if hasUpdate {
				err = client.UpdateThenPoll(ctx, *id, *update)
				if err != nil {
					return fmt.Errorf("updating %s: %v", id, err)
				}
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
				var state AutonomousDatabaseRegularResourceModel
				state.AdminPassword = pointer.From(adbsPropModel.AdminPassword)
				state.BackupRetentionPeriodInDays = pointer.From(adbsPropModel.BackupRetentionPeriodInDays)
				state.CharacterSet = pointer.From(adbsPropModel.CharacterSet)
				state.ComputeCount = pointer.From(adbsPropModel.ComputeCount)
				state.ComputeModel = string(pointer.From(adbsPropModel.ComputeModel))
				state.CustomerContacts = flattenAdbsCustomerContacts(adbsPropModel.CustomerContacts)
				state.DataStorageSizeInTbs = pointer.From(adbsPropModel.DataStorageSizeInGbs)
				state.DbWorkload = string(pointer.From(adbsPropModel.DbWorkload))
				state.DbVersion = pointer.From(adbsPropModel.DbVersion)
				state.DisplayName = pointer.From(adbsPropModel.DisplayName)
				state.AutoScalingEnabled = pointer.From(adbsPropModel.IsAutoScalingEnabled)
				state.AutoScalingForStorageEnabled = pointer.From(adbsPropModel.IsAutoScalingForStorageEnabled)
				state.Location = result.Model.Location
				state.Name = pointer.ToString(result.Model.Name)
				state.LicenseModel = string(pointer.From(adbsPropModel.LicenseModel))
				state.NcharacterSet = pointer.From(adbsPropModel.NcharacterSet)
				state.ResourceGroupName = id.ResourceGroupName
				state.SubnetId = pointer.From(adbsPropModel.SubnetId)
				state.Tags = pointer.From(result.Model.Tags)
				state.VnetId = pointer.From(adbsPropModel.VnetId)

				return metadata.Encode(&state)
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

func expandAdbsCustomerContacts(customerContactsList []string) []autonomousdatabases.CustomerContact {
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

func flattenAdbsCustomerContacts(customerContactsList *[]autonomousdatabases.CustomerContact) []string {
	var customerContacts []string
	if customerContactsList != nil {
		for _, customerContact := range *customerContactsList {
			customerContacts = append(customerContacts, customerContact.Email)
		}
	}
	return customerContacts
}
