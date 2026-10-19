// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_runbook

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/runbook"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/runbookdraft"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAutomationRunbookCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	autoCli := meta.(*clients.Client).Automation
	client := autoCli.Runbook
	jsClient := autoCli.JobSchedule
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	subscriptionID := meta.(*clients.Client).Account.SubscriptionId

	id := runbook.NewRunbookID(subscriptionID, d.Get("resource_group_name").(string), d.Get("automation_account_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %s", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_automation_runbook", id.ID())
			}
		}
	}

	// for existing runbook, if only job_schedule field updated, then skip update runbook
	if d.IsNewResource() || d.HasChangeExcept("job_schedule") {
		location := location.Normalize(d.Get("location").(string))

		parameters := runbook.RunbookCreateOrUpdateParameters{
			Properties: runbook.RunbookCreateOrUpdateProperties{
				LogVerbose:         pointer.To(d.Get("log_verbose").(bool)),
				LogProgress:        pointer.To(d.Get("log_progress").(bool)),
				RuntimeEnvironment: pointer.To(d.Get("runtime_environment_name").(string)),
				RunbookType:        runbook.RunbookTypeEnum(d.Get("runbook_type").(string)),
				Description:        pointer.To(d.Get("description").(string)),
				LogActivityTrace:   pointer.To(int64(d.Get("log_activity_trace_level").(int))),
			},

			Location: &location,
		}
		if tagsVal := helpers.ExpandStringInterfaceMap(d.Get("tags").(map[string]any)); tagsVal != nil {
			parameters.Tags = &tagsVal
		}

		contentLink := expandContentLink(d.Get("publish_content_link").([]any))
		if contentLink != nil {
			parameters.Properties.PublishContentLink = contentLink
		} else {
			parameters.Properties.Draft = &runbook.RunbookDraft{}
			if draft := expandDraft(d.Get("draft").([]any)); draft != nil {
				parameters.Properties.Draft = draft
			}
		}

		if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
			return fmt.Errorf("creating/updating %s: %+v", id, err)
		}

		d.SetId(id.ID())

		if v, ok := d.GetOk("content"); ok {
			content := v.(string)
			draftRunbookID := runbookdraft.NewRunbookID(id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName, id.RunbookName)
			if err := autoCli.RunbookDraft.ReplaceContentThenPoll(ctx, draftRunbookID, []byte(content)); err != nil {
				return fmt.Errorf("setting the draft for %s: %+v", id, err)
			}

			if err := autoCli.Runbook.PublishThenPoll(ctx, id); err != nil {
				return fmt.Errorf("publishing the updated %s: %+v", id, err)
			}
		}
	}

	// **don't need** to list job schedules and delete all of them. update the runbook will recreate these job schedules automatically,
	// but with a different job schedule id
	// crosscheck these existing jobs and jobs from tf, delete the ones not in tf, and create the ones not in api
	// Fix issue: https://github.com/hashicorp/terraform-provider-azurerm/issues/8634
	jsValue, ok := d.GetOk("job_schedule")
	if ok && d.HasChange("job_schedule") {
		jsMap, err := helper.ExpandAutomationJobSchedule(jsValue.(*pluginsdk.Set).List(), id.RunbookName)
		if err != nil {
			return err
		}

		if err := updatedLinkedJobSchedules(ctx, subscriptionID, jsClient, &id, *jsMap); err != nil {
			return fmt.Errorf("update job schedule links: %v", err)
		}
	}

	return resourceAutomationRunbookRead(d, meta)
}
