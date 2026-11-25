// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mysql_flexible_server_configuration

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/configurations"
)

// GTID can only be updated in a specific order: OFF -> OFF_PERMISSIVE -> ON_PERMISSIVE -> ON.
// This function ensures that the GTID mode is updated in the correct order by checking the current value and updating it step by step until it reaches the desired value.
func mysqlFlexibleServerConfigurationUpdateGITDMode(ctx context.Context,
	client *configurations.ConfigurationsClient,
	id configurations.ConfigurationId,
	value string,
	updateCallback func() error,
) error {
	resp, err := client.Get(ctx, id)
	if err != nil {
		return err
	}

	gtidSeq := []string{"OFF", "OFF_PERMISSIVE", "ON_PERMISSIVE", "ON"}
	currentValue := "OFF"
	if resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.Value != nil {
		currentValue = pointer.From(resp.Model.Properties.Value)
	}

	curIdx, toIdx := 0, -1
	for idx, v := range gtidSeq {
		if v == currentValue {
			curIdx = idx
		}

		if v == value {
			toIdx = idx
		}
	}

	if toIdx < curIdx {
		return fmt.Errorf("cannot set `gtid_mode` from %s to %s", currentValue, value)
	}

	for _, v := range gtidSeq[curIdx+1 : toIdx+1] {
		payload := configurations.Configuration{
			Properties: &configurations.ConfigurationProperties{
				Value: pointer.To(v),
			},
		}

		log.Printf("[DEBUG] updating `gtid_mode` of %s to %s", id, v)
		if err := client.UpdateCallbackThenPoll(ctx, id, payload, updateCallback); err != nil {
			return err
		}
	}

	return nil
}
