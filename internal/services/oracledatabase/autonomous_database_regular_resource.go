// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracledatabase

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/autonomousdatabases"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = AdbsRegularResource{}

type AdbsRegularResource struct{}

type AdbsRegularResourceModel struct {
	Location          string                 `tfschema:"location"`
	Name              string                 `tfschema:"name"`
	ResourceGroupName string                 `tfschema:"resource_group_name"`
	Tags              map[string]interface{} `tfschema:"tags"`

	DisplayName          string  `tfschema:"display_name"`
	SubnetId             string  `tfschema:"subnet_id"`
	ComputeModel         string  `tfschema:"compute_model"`
	ComputeCount         float64 `tfschema:"compute_count"`
	LicenseModel         string  `tfschema:"license_model"`
	DataStorageSizeInGbs int64   `tfschema:"data_storage_size_in_gbs"`
	DbWorkload           string  `tfschema:"db_workload"`
	AdminPassword        string  `tfschema:"admin_password"`
	DbVersion            string  `tfschema:"db_version"`
	CharacterSet         string  `tfschema:"character_set"`
	NcharacterSet        string  `tfschema:"ncharacter_set"`
	VnetId               string  `tfschema:"vnet_id"`
}

func (AdbsRegularResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.Location(),
		"tags":     commonschema.Tags(),
		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		// AdbsRegularResource
		"display_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"subnet_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"compute_model": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"compute_count": {
			Type:     pluginsdk.TypeFloat,
			Required: true,
		},
		"license_model": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"data_storage_size_in_gbs": {
			Type:     pluginsdk.TypeInt,
			Required: true,
		},
		"db_workload": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"admin_password": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"db_version": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"character_set": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"ncharacter_set": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"vnet_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
	}
}

func (AdbsRegularResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (AdbsRegularResource) ModelObject() interface{} {
	return &AdbsRegularResource{}
}

func (AdbsRegularResource) ResourceType() string {
	return "azurerm_oracledatabase_autonomous_database_regular"
}

func (r AdbsRegularResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.OracleDatabase.OracleDatabaseClient.AutonomousDatabases
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model AdbsRegularResourceModel
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
					DisplayName:          pointer.To(model.DisplayName),
					SubnetId:             pointer.To(model.SubnetId),
					ComputeModel:         pointer.To(autonomousdatabases.ComputeModel(model.ComputeModel)),
					ComputeCount:         pointer.To(model.ComputeCount),
					LicenseModel:         pointer.To(autonomousdatabases.LicenseModel(model.LicenseModel)),
					DataStorageSizeInGbs: pointer.To(model.DataStorageSizeInGbs),
					DbWorkload:           pointer.To(autonomousdatabases.WorkloadType(model.DbWorkload)),
					AdminPassword:        pointer.To(model.AdminPassword),
					DbVersion:            pointer.To(model.DbVersion),
					CharacterSet:         pointer.To(model.CharacterSet),
					NcharacterSet:        pointer.To(model.NcharacterSet),
					VnetId:               pointer.To(model.VnetId),
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

func (r AdbsRegularResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {

			client := metadata.Client.OracleDatabase.OracleDatabaseClient.AutonomousDatabases
			id, err := autonomousdatabases.ParseAutonomousDatabaseID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model AdbsRegularResourceModel
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

			upd := existing.Model

			err = client.CreateOrUpdateThenPoll(ctx, *id, *upd)
			if err != nil {
				return fmt.Errorf("updating %s: %v", id, err)
			}
			return nil
		},
	}
}

func (AdbsRegularResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := autonomousdatabases.ParseAutonomousDatabaseID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.OracleDatabase.OracleDatabaseClient.AutonomousDatabases
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
				var output AdbsRegularResourceModel
				output.DisplayName = pointer.From(adbsPropModel.DisplayName)
				output.ComputeModel = string(pointer.From(adbsPropModel.ComputeModel))
				output.ComputeCount = pointer.From(adbsPropModel.ComputeCount)
				output.LicenseModel = string(pointer.From(adbsPropModel.LicenseModel))
				output.DataStorageSizeInGbs = pointer.From(adbsPropModel.DataStorageSizeInGbs)
				output.DbWorkload = string(pointer.From(adbsPropModel.DbWorkload))
				output.AdminPassword = pointer.From(adbsPropModel.AdminPassword)
				output.DbVersion = pointer.From(adbsPropModel.DbVersion)
				output.CharacterSet = pointer.From(adbsPropModel.CharacterSet)
				output.NcharacterSet = pointer.From(adbsPropModel.NcharacterSet)
				output.VnetId = pointer.From(adbsPropModel.VnetId)

				return metadata.Encode(&output)
			default:
				return fmt.Errorf("unexpected Autonomous Database type, must be of type Regular")
			}
		},
	}
}

func (AdbsRegularResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.OracleDatabase.OracleDatabaseClient.AutonomousDatabases

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

func (AdbsRegularResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return autonomousdatabases.ValidateAutonomousDatabaseID
}
