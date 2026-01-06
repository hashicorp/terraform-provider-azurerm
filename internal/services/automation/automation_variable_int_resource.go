// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name automation_variable_int -properties "name:variable_name,automation_account_name,resource_group_name" -service-package-name automation -known-values "subscription_id:data.Subscriptions.Primary"

package automation

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/variable"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceAutomationVariableInt() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAutomationVariableIntCreateUpdate,
		Read:   resourceAutomationVariableIntRead,
		Update: resourceAutomationVariableIntCreateUpdate,
		Delete: resourceAutomationVariableIntDelete,

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

		Schema: resourceAutomationVariableCommonSchema(pluginsdk.TypeInt, nil),
	}
}

func resourceAutomationVariableIntCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	return resourceAutomationVariableCreateUpdate(d, meta, "Int")
}

func resourceAutomationVariableIntRead(d *pluginsdk.ResourceData, meta interface{}) error {
	return resourceAutomationVariableRead(d, meta, "Int")
}

func resourceAutomationVariableIntDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	return resourceAutomationVariableDelete(d, meta, "Int")
}
