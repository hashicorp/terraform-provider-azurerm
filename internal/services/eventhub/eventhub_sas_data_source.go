// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventhub

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/eventhub"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

const (
	connStringSharedAccessKeyKey     = "SharedAccessKey"
	connStringSharedAccessKeyNameKey = "SharedAccessKeyName"
	connStringEndpointKey            = "Endpoint"
	connStringEntityPathKey          = "EntityPath"
)

var _ sdk.DataSource = EventHubSharedAccessSignatureDataSource{}

type EventHubSharedAccessSignatureDataSource struct{}

type EventHubSharedAccessSignatureDataSourceModel struct {
	ConnectionString string `tfschema:"connection_string"`
	Expiry           string `tfschema:"expiry"`
	Sas              string `tfschema:"sas"`
}

func (EventHubSharedAccessSignatureDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"connection_string": {
			Type:      pluginsdk.TypeString,
			Required:  true,
			Sensitive: true,
		},

		// Always in UTC and must be ISO-8601 format
		"expiry": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ISO8601DateTime,
		},
	}
}

func (EventHubSharedAccessSignatureDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"sas": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},
	}
}

func (EventHubSharedAccessSignatureDataSource) ModelObject() interface{} {
	return &EventHubSharedAccessSignatureDataSourceModel{}
}

func (EventHubSharedAccessSignatureDataSource) ResourceType() string {
	return "azurerm_eventhub_sas"
}

func (EventHubSharedAccessSignatureDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var state EventHubSharedAccessSignatureDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			// Parse the connection string
			kvp, err := eventhub.ParseEventHubSASConnectionString(state.ConnectionString)
			if err != nil {
				return err
			}

			sharedAccessKeyName := kvp[connStringSharedAccessKeyNameKey]
			sharedAccessKey := kvp[connStringSharedAccessKeyKey]
			endpoint := kvp[connStringEndpointKey]
			entityPath := kvp[connStringEntityPathKey]
			endpointUrl, err := eventhub.ComputeEventHubSASConnectionUrl(endpoint, entityPath)
			if err != nil {
				return err
			}

			sasToken, err := eventhub.ComputeEventHubSASToken(sharedAccessKeyName, sharedAccessKey, *endpointUrl, state.Expiry)
			if err != nil {
				return err
			}

			sasConnectionString := eventhub.ComputeEventHubSASConnectionString(sasToken)
			state.Sas = sasConnectionString

			tokenHash := sha256.Sum256([]byte(sasConnectionString))
			metadata.ResourceData.SetId(hex.EncodeToString(tokenHash[:]))

			return metadata.Encode(&state)
		},
	}
}
