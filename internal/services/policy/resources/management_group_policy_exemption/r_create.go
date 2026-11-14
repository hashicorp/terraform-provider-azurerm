// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package management_group_policy_exemption

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2021-06-01-preview/policy" //nolint:staticcheck
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	managementGroupParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceManagementGroupPolicyExemptionCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Policy.ExemptionsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewResourcePolicyExemptionId(d.Get("management_group_id").(string), d.Get("name").(string))

	managementGroupId, err := managementGroupParse.ManagementGroupID(id.ResourceId)
	if err != nil {
		return err
	}

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, managementGroupId.ID(), id.Name)
			if err != nil {
				if !response.WasNotFound(existing.Response.Response) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id.ID(), err)
				}
			}
			if existing.ID != nil && *existing.ID != "" {
				return tf.ImportAsExistsError("azurerm_management_group_policy_exemption", *existing.ID)
			}
		}
	}

	exemption := policy.Exemption{
		ExemptionProperties: &policy.ExemptionProperties{
			PolicyAssignmentID:           pointer.To(d.Get("policy_assignment_id").(string)),
			PolicyDefinitionReferenceIds: pluginsdk.ExpandStringSlice(d.Get("policy_definition_reference_ids").([]any)),
			ExemptionCategory:            policy.ExemptionCategory(d.Get("exemption_category").(string)),
		},
	}

	if v, ok := d.GetOk("display_name"); ok {
		exemption.DisplayName = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		exemption.Description = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("expires_on"); ok {
		t, err := date.ParseTime(time.RFC3339, v.(string))
		if err != nil {
			return fmt.Errorf("expanding `expires_on`: %+v", err)
		}
		exemption.ExpiresOn = &date.Time{Time: t}
	}

	if metaDataString := d.Get("metadata").(string); metaDataString != "" {
		metaData, err := structure.ExpandJsonFromString(metaDataString)
		if err != nil {
			return fmt.Errorf("unable to parse metadata: %+v", err)
		}
		exemption.Metadata = &metaData
	}

	if _, err := client.CreateOrUpdate(ctx, managementGroupId.ID(), id.Name, exemption); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id.ID(), err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	return resourceManagementGroupPolicyExemptionRead(d, meta)
}
