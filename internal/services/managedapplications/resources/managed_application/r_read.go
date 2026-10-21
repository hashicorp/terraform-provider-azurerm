// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managed_application

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedapplications/2021-07-01/applications"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceManagedApplicationRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ManagedApplication.ApplicationClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := applications.ParseApplicationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Managed Application %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to read Managed Application %s: %+v", id, err)
	}

	d.Set("name", id.ApplicationName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		p := model.Properties

		d.Set("location", location.NormalizeNilable(model.Location))
		d.Set("kind", model.Kind)
		if err := d.Set("plan", flattenManagedApplicationPlan(model.Plan)); err != nil {
			return fmt.Errorf("setting `plan`: %+v", err)
		}

		id, err := commonids.ParseResourceGroupIDInsensitively(pointer.From(p.ManagedResourceGroupId))
		if err != nil {
			return err
		}

		d.Set("managed_resource_group_name", id.ResourceGroupName)
		d.Set("application_definition_id", p.ApplicationDefinitionId)

		expendedParams, err := expandManagedApplicationParameters(d)
		if err != nil {
			return fmt.Errorf("expanding `parameter_values`: %+v", err)
		}

		parameterValues, err := flattenManagedApplicationParameterValuesValueToString(p.Parameters, *expendedParams)
		if err != nil {
			return fmt.Errorf("serializing JSON from `parameter_values`: %+v", err)
		}
		d.Set("parameter_values", parameterValues)

		managedApplicationIdentity, err := flattenManagedApplicationIdentity(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		d.Set("identity", managedApplicationIdentity)

		outputs, err := flattenManagedApplicationOutputs(p.Outputs)
		if err != nil {
			return err
		}
		if err = d.Set("outputs", outputs); err != nil {
			return err
		}

		if err = tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	return nil
}

func resourceManagedApplicationCustomizeDiff(ctx context.Context, diff *pluginsdk.ResourceDiff, meta any) error {
	oldVal, newVal := diff.GetChange("identity.#")
	if oldVal.(int) == 1 && newVal.(int) == 0 {
		return diff.ForceNew("identity")
	}

	return nil
}

func extractParameterOrOutputValue(v any) (string, error) {
	switch t := v.(type) {
	case bool:
		return strconv.FormatBool(v.(bool)), nil
	case float64:
		// use precision 0 since this comes from an int
		return fmt.Sprintf("%.f", v.(float64)), nil
	case string:
		return v.(string), nil
	case map[string]any:
		return compactParameterOrOutputValue(v)
	case []any:
		return compactParameterOrOutputValue(v)
	default:
		return "", fmt.Errorf("unexpected type %T", t)
	}
}

func compactParameterOrOutputValue(v any) (string, error) {
	result, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	compactJson := bytes.Buffer{}
	if err = json.Compact(&compactJson, result); err != nil {
		return "", err
	}
	return compactJson.String(), nil
}
