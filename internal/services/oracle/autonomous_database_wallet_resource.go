package oracle

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/autonomousdatabases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = AutonomousDatabaseWalletResource{}

type AutonomousDatabaseWalletResource struct{}

type AutonomousDatabaseWalletResourceModel struct {
	AutonomousDatabaseId string `tfschema:"autonomous_database_id"`
	Password             string `tfschema:"password"`
	GenerateType         string `tfschema:"generate_type"`
	IsRegional           bool   `tfschema:"is_regional"`
	WalletFiles          string `tfschema:"wallet_files"`
}

func (AutonomousDatabaseWalletResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		// Required
		"autonomous_database_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"password": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringLenBetween(8, 30),
		},

		// Optional
		"generate_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  "SINGLE",
			ValidateFunc: validation.StringInSlice([]string{
				"SINGLE",
				"ALL",
			}, false),
		},
		"is_regional": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},
	}
}

func (AutonomousDatabaseWalletResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"wallet_files": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
			Description: "The base64 encoded wallet files",
		},
	}
}

func (AutonomousDatabaseWalletResource) ModelObject() interface{} {
	return &AutonomousDatabaseWalletResourceModel{}
}

func (AutonomousDatabaseWalletResource) ResourceType() string {
	return "azurerm_oracle_autonomous_database_wallet"
}

func (r AutonomousDatabaseWalletResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.AutonomousDatabases

			var model AutonomousDatabaseWalletResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding Autonomous Database Wallet resource model: %+v", err)
			}

			autonomousDatabaseId, err := autonomousdatabases.ParseAutonomousDatabaseID(model.AutonomousDatabaseId)
			if err != nil {
				return fmt.Errorf("parsing autonomous database ID: %+v", err)
			}

			generateType := autonomousdatabases.GenerateTypeSingle
			if model.GenerateType == "REGIONAL" {
				generateType = autonomousdatabases.GenerateTypeAll
			}

			input := autonomousdatabases.GenerateAutonomousDatabaseWalletDetails{
				Password:     model.Password,
				GenerateType: &generateType,
			}

			walletResult, err := client.GenerateWallet(ctx, *autonomousDatabaseId, input)
			if err != nil {
				return fmt.Errorf("generating wallet for %s: %+v", model.AutonomousDatabaseId, err)
			}

			if walletResult.Model == nil {
				return fmt.Errorf("empty response when generating wallet for %s", model.AutonomousDatabaseId)
			}

			walletFiles := walletResult.Model.WalletFiles
			if walletFiles == "" {
				return fmt.Errorf("no wallet files in response for %s", model.AutonomousDatabaseId)
			}

			// Set the wallet files in our model
			model.WalletFiles = walletFiles
			metadata.SetID(autonomousDatabaseId)

			return metadata.Encode(&model)
		},
	}
}

func (AutonomousDatabaseWalletResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// Since wallet is generated on-demand, just verify autonomous database exists
			id, err := autonomousdatabases.ParseAutonomousDatabaseID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.Oracle.OracleClient.AutonomousDatabases
			_, err = client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			// Keep existing computed values in state
			return nil
		},
	}
}

func (AutonomousDatabaseWalletResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// Wallets don't need to be deleted as they're generated on-demand
			return nil
		},
	}
}

func (r AutonomousDatabaseWalletResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return func(i interface{}, k string) (warnings []string, errors []error) {
		v := i.(string)
		if _, err := autonomousdatabases.ParseAutonomousDatabaseID(v); err != nil {
			errors = append(errors, fmt.Errorf("invalid %q: %s. The ID should be the Autonomous Database ID followed by '/wallet'", k, err))
		}
		return
	}
}
