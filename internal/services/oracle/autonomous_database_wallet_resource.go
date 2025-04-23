// Copyright © 2025, Your Organization. All rights reserved.

package oracle

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/autonomousdatabases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = AutonomousDatabaseWalletResource{}

type AutonomousDatabaseWalletResource struct{}

type AutonomousDatabaseWalletResourceModel struct {
	// Required
	AutonomousDatabaseId string `tfschema:"autonomous_database_id"`
	Password             string `tfschema:"password"`

	// Optional
	GenerateType *string `tfschema:"generate_type"`
	Base64Encode *bool   `tfschema:"base64_encode"`

	// Computed
	Content        string `tfschema:"content"`
	WalletFileName string `tfschema:"wallet_file_name"`
}

func (AutonomousDatabaseWalletResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		// Required
		"autonomous_database_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: func(i interface{}, k string) (warnings []string, errors []error) {
				v := i.(string)
				if _, err := autonomousdatabases.ParseAutonomousDatabaseID(v); err != nil {
					errors = append(errors, fmt.Errorf("invalid %q: %s", k, err))
				}
				return
			},
		},

		"password": {
			Type:      pluginsdk.TypeString,
			Required:  true,
			ForceNew:  true,
			Sensitive: true,
		},

		// Optional
		"generate_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "SINGLE",
			ForceNew: true,
			ValidateFunc: func(i interface{}, k string) (warnings []string, errors []error) {
				v := i.(string)
				validTypes := []string{"SINGLE", "REGIONAL"}
				isValid := false
				for _, t := range validTypes {
					if v == t {
						isValid = true
						break
					}
				}
				if !isValid {
					errors = append(errors, fmt.Errorf("%q must be one of: %v", k, validTypes))
				}
				return
			},
		},

		"base64_encode": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
			ForceNew: true,
		},

		// Computed
		"content": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"wallet_file_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (AutonomousDatabaseWalletResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
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

			input := autonomousdatabases.GenerateAutonomousDatabaseWalletDetails{
				Password: model.Password,
			}
			if model.GenerateType != nil {
				input.GenerateType = (*autonomousdatabases.GenerateType)(model.GenerateType)
			}

			walletResult, err := client.GenerateWallet(ctx, *autonomousDatabaseId, input)
			if err != nil {
				return fmt.Errorf("generating wallet for %s: %+v", model.AutonomousDatabaseId, err)
			}

			if walletResult.Model == nil {
				return fmt.Errorf("received empty wallet model for %s", model.AutonomousDatabaseId)
			}

			walletBytes := []byte(walletResult.Model.WalletFiles)
			content := string(walletBytes)

			if model.Base64Encode != nil && *model.Base64Encode {
				content = base64.StdEncoding.EncodeToString(walletBytes)
			}

			if err := metadata.Encode(&AutonomousDatabaseWalletResourceModel{
				AutonomousDatabaseId: model.AutonomousDatabaseId,
				Password:             model.Password,
				GenerateType:         model.GenerateType,
				Content:              content,
			}); err != nil {
				return fmt.Errorf("encoding Autonomous Database Wallet resource model: %+v", err)
			}

			// Construct the resource ID correctly
			metadata.SetID(autonomousDatabaseId)

			return nil
		},
	}
}

func (AutonomousDatabaseWalletResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			return nil
		},
	}
}

func (AutonomousDatabaseWalletResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
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
