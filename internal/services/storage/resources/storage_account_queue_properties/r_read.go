// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_account_queue_properties

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	storageClient2 "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/jackofallops/giovanni/storage/2023-11-03/queue/queues"
)

func (s Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			storageClient := storageClient2.From(metadata.Client)

			var state AccountQueuePropertiesModel

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

			client, err := storageClient.QueuesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
			if err != nil {
				return fmt.Errorf("building Queues Client for %s: %v", *id, err)
			}

			props, err := client.GetServiceProperties(ctx)
			if err != nil {
				return fmt.Errorf("retrieving Queue Properties for %s: %+v", *id, err)
			}

			if props != nil {
				if props.Cors != nil && !reflect.DeepEqual(*props.Cors, &defaultCorsProperties) {
					corsRules := make([]AccountQueuePropertiesCorsRule, 0)
					for _, rule := range props.Cors.CorsRule {
						corsRule := AccountQueuePropertiesCorsRule{
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

				if props.HourMetrics != nil && !reflect.DeepEqual(*props.HourMetrics, &defaultHourMetricsProperties) {
					state.HourMetrics = []AccountQueuePropertiesHourMetrics{
						{
							Version:             props.HourMetrics.Version,
							IncludeAPIS:         pointer.From(props.HourMetrics.IncludeAPIs),
							RetentionPolicyDays: int64(props.HourMetrics.RetentionPolicy.Days),
						},
					}
				}

				if props.MinuteMetrics != nil && !reflect.DeepEqual(*props.MinuteMetrics, &defaultMinuteMetricsProperties) {
					state.MinuteMetrics = []AccountQueuePropertiesMinuteMetrics{
						{
							Version:             props.MinuteMetrics.Version,
							IncludeAPIS:         pointer.From(props.MinuteMetrics.IncludeAPIs),
							RetentionPolicyDays: int64(props.MinuteMetrics.RetentionPolicy.Days),
						},
					}
				}

				if props.Logging != nil && !reflect.DeepEqual(*props.Logging, &defaultLoggingProperties) {
					state.Logging = []AccountQueuePropertiesLogging{
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

type AccountQueuePropertiesCorsRule struct {
	AllowedOrigins []string `tfschema:"allowed_origins"`
	AllowedMethods []string `tfschema:"allowed_methods"`
	AllowedHeaders []string `tfschema:"allowed_headers"`
	ExposedHeaders []string `tfschema:"exposed_headers"`
	MaxAgeSeconds  int64    `tfschema:"max_age_in_seconds"`
}

type AccountQueuePropertiesHourMetrics struct {
	Version             string `tfschema:"version"`
	IncludeAPIS         bool   `tfschema:"include_apis"`
	RetentionPolicyDays int64  `tfschema:"retention_policy_days"`
}

type AccountQueuePropertiesMinuteMetrics struct {
	Version             string `tfschema:"version"`
	IncludeAPIS         bool   `tfschema:"include_apis"`
	RetentionPolicyDays int64  `tfschema:"retention_policy_days"`
}

type AccountQueuePropertiesLogging struct {
	Version             string `tfschema:"version"`
	Delete              bool   `tfschema:"delete"`
	Read                bool   `tfschema:"read"`
	Write               bool   `tfschema:"write"`
	RetentionPolicyDays int64  `tfschema:"retention_policy_days"`
}

var defaultCorsProperties = queues.Cors{
	CorsRule: []queues.CorsRule{},
}

var defaultHourMetricsProperties = queues.MetricsConfig{
	Version:         "1.0",
	RetentionPolicy: queues.RetentionPolicy{},
}

var defaultMinuteMetricsProperties = queues.MetricsConfig{
	Version:         "1.0",
	RetentionPolicy: queues.RetentionPolicy{},
}

var defaultLoggingProperties = queues.LoggingConfig{
	Version:         "1.0",
	RetentionPolicy: queues.RetentionPolicy{},
}

func DefaultValueForAccountQueueProperties() queues.StorageServiceProperties {
	return queues.StorageServiceProperties{
		Logging: &queues.LoggingConfig{
			Version:         "1.0",
			RetentionPolicy: queues.RetentionPolicy{},
		},
		HourMetrics: &queues.MetricsConfig{
			Version:         "1.0",
			RetentionPolicy: queues.RetentionPolicy{},
		},
		MinuteMetrics: &queues.MetricsConfig{
			Version:         "1.0",
			RetentionPolicy: queues.RetentionPolicy{},
		},
		Cors: &queues.Cors{
			CorsRule: []queues.CorsRule{},
		},
	}
}
