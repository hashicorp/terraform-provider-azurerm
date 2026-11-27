// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package security_center_assessment_policy

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2021-06-01/assessmentsmetadata"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	securitycenterClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmSecurityCenterAssessmentPolicyDelete(d *pluginsdk.ResourceData, meta any) error {
	client := securitycenterClient.From(meta).AssessmentsMetadataClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := assessmentsmetadata.ParseProviderAssessmentMetadataID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.DeleteInSubscription(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
