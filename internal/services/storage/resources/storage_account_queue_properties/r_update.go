// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_account_queue_properties

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	storageClient2 "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/jackofallops/giovanni/storage/2023-11-03/queue/queues"
)

func (s Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			storageClient := storageClient2.From(metadata.Client)

			id, err := commonids.ParseStorageAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			account, err := storageClient.GetAccount(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
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
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			var model AccountQueuePropertiesModel

			if err = metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if metadata.ResourceData.HasChange("cors_rule") {
				if len(model.CorsRule) >= 1 {
					corsRules := make([]queues.CorsRule, 0)
					for _, corsRule := range model.CorsRule {
						corsRules = append(corsRules, queues.CorsRule{
							AllowedOrigins:  strings.Join(corsRule.AllowedOrigins, ","),
							AllowedMethods:  strings.Join(corsRule.AllowedMethods, ","),
							AllowedHeaders:  strings.Join(corsRule.AllowedHeaders, ","),
							ExposedHeaders:  strings.Join(corsRule.ExposedHeaders, ","),
							MaxAgeInSeconds: int(corsRule.MaxAgeSeconds),
						})
					}

					if props.Cors == nil {
						props.Cors = &queues.Cors{}
					}
					props.Cors.CorsRule = corsRules
				} else {
					props.Cors = pointer.To(defaultCorsProperties)
				}
			}

			if metadata.ResourceData.HasChange("hour_metrics") {
				if len(model.HourMetrics) == 1 {
					metrics := model.HourMetrics[0]
					if metadata.ResourceData.HasChange("hour_metrics.0.version") {
						props.HourMetrics.Version = metrics.Version
					}

					if metadata.ResourceData.HasChange("hour_metrics.0.include_apis") {
						props.HourMetrics.IncludeAPIs = pointer.To(metrics.IncludeAPIS)
					}

					if metadata.ResourceData.HasChange("hour_metrics.0.retention_policy_days") {
						props.HourMetrics.RetentionPolicy = queues.RetentionPolicy{
							Days:    int(metrics.RetentionPolicyDays),
							Enabled: true,
						}
					}
				} else {
					props.HourMetrics = pointer.To(defaultHourMetricsProperties)
				}
			}

			if metadata.ResourceData.HasChange("minute_metrics") {
				if len(model.MinuteMetrics) == 1 {
					metrics := model.MinuteMetrics[0]
					if metadata.ResourceData.HasChange("minute_metrics.0.version") {
						props.MinuteMetrics.Version = metrics.Version
					}

					if metadata.ResourceData.HasChange("minute_metrics.0.include_apis") {
						props.MinuteMetrics.IncludeAPIs = pointer.To(metrics.IncludeAPIS)
					}

					if metadata.ResourceData.HasChange("minute_metrics.0.retention_policy_days") {
						props.MinuteMetrics.RetentionPolicy = queues.RetentionPolicy{
							Days:    int(metrics.RetentionPolicyDays),
							Enabled: true,
						}
					}
				} else {
					props.MinuteMetrics = pointer.To(defaultMinuteMetricsProperties)
				}
			}

			if metadata.ResourceData.HasChange("logging") {
				if len(model.Logging) == 1 {
					logging := model.Logging[0]
					if metadata.ResourceData.HasChange("logging.0.version") {
						props.Logging.Version = logging.Version
					}
					if metadata.ResourceData.HasChange("logging.0.delete") {
						props.Logging.Delete = logging.Delete
					}
					if metadata.ResourceData.HasChange("logging.0.read") {
						props.Logging.Read = logging.Read
					}
					if metadata.ResourceData.HasChange("logging.0.write") {
						props.Logging.Write = logging.Write
					}
					if metadata.ResourceData.HasChange("logging.0.retention_policy_days") {
						props.Logging.RetentionPolicy = queues.RetentionPolicy{
							Days:    int(logging.RetentionPolicyDays),
							Enabled: true,
						}
					}
				} else {
					props.Logging = pointer.To(defaultLoggingProperties)
				}
			}

			if err = client.UpdateServiceProperties(ctx, *props); err != nil {
				return fmt.Errorf("updating Queue Properties for %s: %+v", *id, err)
			}

			return nil
		},
	}
}
