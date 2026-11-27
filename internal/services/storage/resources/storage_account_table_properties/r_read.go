// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_account_table_properties

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/shim"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/jackofallops/giovanni/storage/2023-11-03/table/tables"
)

func (s Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			storageClient := metadata.Client.Storage

			var state AccountTablePropertiesModel

			id, err := commonids.ParseStorageAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			state.StorageAccountId = id.ID()

			account, err := storageClient.GetAccount(ctx, *id)
			if err != nil {
				return metadata.MarkAsGone(id)
			}
			if account == nil {
				return fmt.Errorf("unable to locate %s", *id)
			}

			client, err := storageClient.TablesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
			if err != nil {
				return fmt.Errorf("building Tables Client for %s: %v", *id, err)
			}

			props, err := client.GetServiceProperties(ctx)
			if err != nil {
				return fmt.Errorf("retrieving Table Properties for %s: %+v", *id, err)
			}

			if props != nil {
				if props.Cors != nil && !reflect.DeepEqual(*props.Cors, &defaultTableCorsProperties) {
					corsRules := make([]AccountTablePropertiesCorsRule, 0)
					for _, rule := range props.Cors.CorsRule {
						corsRule := AccountTablePropertiesCorsRule{
							AllowedOrigins: strings.Split(rule.AllowedOrigins, ","),
							AllowedMethods: strings.Split(rule.AllowedMethods, ","),
							AllowedHeaders: strings.Split(rule.AllowedHeaders, ","),
							ExposedHeaders: strings.Split(rule.ExposedHeaders, ","),
							MaxAgeSeconds:  int64(rule.MaxAgeInSeconds),
						}
						corsRules = append(corsRules, corsRule)
					}
					state.CorsRule = corsRules
				}

				if props.HourMetrics != nil && !reflect.DeepEqual(*props.HourMetrics, defaultTableHourMetricsProperties) {
					state.HourMetrics = []AccountTablePropertiesHourMetrics{
						{
							Version:             props.HourMetrics.Version,
							IncludeAPIS:         pointer.From(props.HourMetrics.IncludeAPIs),
							RetentionPolicyDays: int64(props.HourMetrics.RetentionPolicy.Days),
						},
					}
				}

				if props.MinuteMetrics != nil && !reflect.DeepEqual(*props.MinuteMetrics, defaultTableMinuteMetricsProperties) {
					state.MinuteMetrics = []AccountTablePropertiesMinuteMetrics{
						{
							Version:             props.MinuteMetrics.Version,
							IncludeAPIS:         pointer.From(props.MinuteMetrics.IncludeAPIs),
							RetentionPolicyDays: int64(props.MinuteMetrics.RetentionPolicy.Days),
						},
					}
				}

				if props.Logging != nil && !reflect.DeepEqual(*props.Logging, defaultTableLoggingProperties) {
					state.Logging = []AccountTablePropertiesLogging{
						{
							Version:             props.Logging.Version,
							Delete:              props.Logging.Delete,
							Read:                props.Logging.Read,
							Write:               props.Logging.Write,
							RetentionPolicyDays: int64(props.Logging.RetentionPolicy.Days),
						},
					}
				}
			}

			if err = pluginsdk.SetResourceIdentityData(metadata.ResourceData, id, pluginsdk.ResourceTypeForIdentityVirtual); err != nil {
				return err
			}

			return metadata.Encode(&state)
		},
	}
}

type AccountTablePropertiesCorsRule struct {
	AllowedOrigins []string `tfschema:"allowed_origins"`
	AllowedMethods []string `tfschema:"allowed_methods"`
	AllowedHeaders []string `tfschema:"allowed_headers"`
	ExposedHeaders []string `tfschema:"exposed_headers"`
	MaxAgeSeconds  int64    `tfschema:"max_age_in_seconds"`
}

type AccountTablePropertiesHourMetrics struct {
	Version             string `tfschema:"version"`
	IncludeAPIS         bool   `tfschema:"include_apis"`
	RetentionPolicyDays int64  `tfschema:"retention_policy_days"`
}

type AccountTablePropertiesMinuteMetrics struct {
	Version             string `tfschema:"version"`
	IncludeAPIS         bool   `tfschema:"include_apis"`
	RetentionPolicyDays int64  `tfschema:"retention_policy_days"`
}

type AccountTablePropertiesLogging struct {
	Version             string `tfschema:"version"`
	Delete              bool   `tfschema:"delete"`
	Read                bool   `tfschema:"read"`
	Write               bool   `tfschema:"write"`
	RetentionPolicyDays int64  `tfschema:"retention_policy_days"`
}

var defaultTableCorsProperties = tables.Cors{
	CorsRule: []tables.CorsRule{},
}

var defaultTableHourMetricsProperties = tables.MetricsConfig{
	Version:         "1.0",
	RetentionPolicy: tables.RetentionPolicy{},
}

var defaultTableMinuteMetricsProperties = tables.MetricsConfig{
	Version:         "1.0",
	RetentionPolicy: tables.RetentionPolicy{},
}

var defaultTableLoggingProperties = tables.LoggingConfig{
	Version:         "1.0",
	RetentionPolicy: tables.RetentionPolicy{},
}

func pollForTableProperties(ctx context.Context, client shim.StorageTableWrapper, expected tables.StorageServiceProperties) error {
	pollerType := custompollers.NewDataPlaneTablesPropertiesPoller(client, expected)
	propertiesPoller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	return propertiesPoller.PollUntilDone(ctx)
}

func DefaultValueForAccountTableProperties() tables.StorageServiceProperties {
	return tables.StorageServiceProperties{
		Logging: &tables.LoggingConfig{
			Version:         "1.0",
			RetentionPolicy: tables.RetentionPolicy{},
		},
		HourMetrics: &tables.MetricsConfig{
			Version:         "1.0",
			RetentionPolicy: tables.RetentionPolicy{},
		},
		MinuteMetrics: &tables.MetricsConfig{
			Version:         "1.0",
			RetentionPolicy: tables.RetentionPolicy{},
		},
		Cors: &tables.Cors{
			CorsRule: []tables.CorsRule{},
		},
	}
}
