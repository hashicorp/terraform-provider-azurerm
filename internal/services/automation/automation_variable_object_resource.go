// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/variable"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func resourceAutomationVariableObject() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAutomationVariableObjectCreateUpdate,
		Read:   resourceAutomationVariableObjectRead,
		Update: resourceAutomationVariableObjectCreateUpdate,
		Delete: resourceAutomationVariableObjectDelete,

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

		// @TODO Create a validation function to check that we have a JSON object passed to us?
		Schema: resourceAutomationVariableCommonSchema(pluginsdk.TypeString, validation.StringIsNotEmpty),
	}
}

func resourceAutomationVariableObjectCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	return resourceAutomationVariableCreateUpdate(d, meta, "Object")
}

func resourceAutomationVariableObjectRead(d *pluginsdk.ResourceData, meta interface{}) error {
	return resourceAutomationVariableRead(d, meta, "Object")
}

func resourceAutomationVariableObjectDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	return resourceAutomationVariableDelete(d, meta, "Object")
}
