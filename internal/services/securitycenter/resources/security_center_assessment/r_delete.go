// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package security_center_assessment

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSecurityCenterAssessmentDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).SecurityCenter.AssessmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AssessmentID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.TargetResourceID, id.Name); err != nil {
		return fmt.Errorf("deleting Security Center Assessment %q (target resource id %q) : %+v", id.Name, id.TargetResourceID, err)
	}

	return nil
}
