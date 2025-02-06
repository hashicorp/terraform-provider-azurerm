// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/variable"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func resourceAutomationVariableString() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAutomationVariableStringCreateUpdate,
		Read:   resourceAutomationVariableStringRead,
		Update: resourceAutomationVariableStringCreateUpdate,
		Delete: resourceAutomationVariableStringDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := variable.ParseVariableID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: resourceAutomationVariableCommonSchema(pluginsdk.TypeString, validation.StringIsNotEmpty),
	}
}

func resourceAutomationVariableStringCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	return resourceAutomationVariableCreateUpdate(d, meta, "String")
}

func resourceAutomationVariableStringRead(d *pluginsdk.ResourceData, meta interface{}) error {
	return resourceAutomationVariableRead(d, meta, "String")
}

func resourceAutomationVariableStringDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	return resourceAutomationVariableDelete(d, meta, "String")
}
