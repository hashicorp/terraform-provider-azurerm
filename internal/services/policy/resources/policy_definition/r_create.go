// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package policy_definition

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2021-06-01-preview/policy" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/custompollers"
	mgmtGrpParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePolicyDefinitionCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Policy.DefinitionsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	policyType := d.Get("policy_type").(string)
	mode := d.Get("mode").(string)
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)

	managementGroupName := ""
	if v, ok := d.GetOk("management_group_id"); ok {
		id, err := mgmtGrpParse.ManagementGroupID(v.(string))
		if err != nil {
			return err
		}
		managementGroupName = id.Name
	}

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := helpers.GetPolicyDefinitionByName(ctx, client, name, managementGroupName)
			if err != nil {
				if !response.WasNotFound(existing.Response.Response) {
					return fmt.Errorf("checking for presence of existing Policy Definition %q: %+v", name, err)
				}
			}

			if existing.ID != nil && *existing.ID != "" {
				return tf.ImportAsExistsError("azurerm_policy_definition", *existing.ID)
			}
		}
	}

	properties := policy.DefinitionProperties{
		PolicyType:  policy.Type(policyType),
		Mode:        pointer.To(mode),
		DisplayName: pointer.To(displayName),
		Description: pointer.To(description),
	}

	if policyRuleString := d.Get("policy_rule").(string); policyRuleString != "" {
		policyRule, err := pluginsdk.ExpandJsonFromString(policyRuleString)
		if err != nil {
			return fmt.Errorf("expanding JSON for `policy_rule`: %+v", err)
		}
		properties.PolicyRule = &policyRule
	}

	if metaDataString := d.Get("metadata").(string); metaDataString != "" {
		metaData, err := pluginsdk.ExpandJsonFromString(metaDataString)
		if err != nil {
			return fmt.Errorf("expanding JSON for `metadata`: %+v", err)
		}
		properties.Metadata = &metaData
	}

	if parametersString := d.Get("parameters").(string); parametersString != "" {
		parameters, err := expandParameterDefinitionsValueFromStringTrack1(parametersString)
		if err != nil {
			return fmt.Errorf("expanding JSON for `parameters`: %+v", err)
		}
		properties.Parameters = parameters
	}

	definition := policy.Definition{
		Name:                 pointer.To(name),
		DefinitionProperties: &properties,
	}

	var err error

	if managementGroupName == "" {
		_, err = client.CreateOrUpdate(ctx, name, definition)
	} else {
		_, err = client.CreateOrUpdateAtManagementGroup(ctx, name, definition, managementGroupName)
	}

	if err != nil {
		return fmt.Errorf("creating/updating Policy Definition %q: %+v", name, err)
	}

	// Policy Definitions are eventually consistent; wait for them to stabilize
	log.Printf("[DEBUG] Waiting for Policy Definition %q to become available", name)
	poller := custompollers.NewEventualConsistencyPoller(10, func(pollerCtx context.Context) (*http.Response, error) {
		res, err := helpers.GetPolicyDefinitionByName(pollerCtx, client, name, managementGroupName)
		return res.Response.Response, err
	}, &custompollers.EventualConsistencyPollerOptions{
		Interval:              10 * time.Second,
		TargetStatusCode:      pointer.To(http.StatusOK),
		RetryErrorStatusCodes: []int{http.StatusNotFound},
	})
	if err = poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for Policy Definition %q to become available: %+v", name, err)
	}

	resp, err := helpers.GetPolicyDefinitionByName(ctx, client, name, managementGroupName)
	if err != nil {
		return err
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Policy Assignment %q", name)
	}

	id, err := parse.PolicyDefinitionID(*resp.ID)
	if err != nil {
		return fmt.Errorf("failed to flatten Policy Parameters %q: %+v", *resp.ID, err)
	}
	d.SetId(id.Id)

	return resourcePolicyDefinitionRead(d, meta)
}
