// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/variable"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceAutomationVariableBool() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAutomationVariableBoolCreateUpdate,
		Read:   resourceAutomationVariableBoolRead,
		Update: resourceAutomationVariableBoolCreateUpdate,
		Delete: resourceAutomationVariableBoolDelete,

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
