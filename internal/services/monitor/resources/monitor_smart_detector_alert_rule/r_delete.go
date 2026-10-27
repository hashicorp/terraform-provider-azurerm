// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_smart_detector_alert_rule

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2019-06-01/smartdetectoralertrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMonitorSmartDetectorAlertRuleDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Monitor.SmartDetectorAlertRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := smartdetectoralertrules.ParseSmartDetectorAlertRuleID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}
	return nil
}
