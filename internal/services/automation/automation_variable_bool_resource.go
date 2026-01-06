// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name automation_variable_bool -properties "name:variable_name,automation_account_name,resource_group_name" -service-package-name automation -known-values "subscription_id:data.Subscriptions.Primary"

package automation

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/variable"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceAutomationVariableBool() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAutomationVariableBoolCreateUpdate,
		Read:   resourceAutomationVariableBoolRead,
		Update: resourceAutomationVariableBoolCreateUpdate,
		Delete: resourceAutomationVariableBoolDelete,

		Importer: pluginsdk.ImporterValidatingIdentity(&variable.VariableId{}),
		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&variable.VariableId{}),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: resourceAutomationVariableCommonSchema(pluginsdk.TypeBool, nil),
	}
}

func resourceAutomationVariableBoolCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	return resourceAutomationVariableCreateUpdate(d, meta, "Bool")
}

func resourceAutomationVariableBoolRead(d *pluginsdk.ResourceData, meta interface{}) error {
	return resourceAutomationVariableRead(d, meta, "Bool")
}

func resourceAutomationVariableBoolDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	return resourceAutomationVariableDelete(d, meta, "Bool")
}
