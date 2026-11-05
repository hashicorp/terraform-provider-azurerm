// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_account

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/automationaccount"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create:   resourceAutomationAccountCreate,
		Read:     resourceAutomationAccountRead,
		Update:   resourceAutomationAccountUpdate,
		Delete:   resourceAutomationAccountDelete,
		Importer: pluginsdk.ImporterValidatingIdentity(&automationaccount.AutomationAccountId{}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: automationAccountSchema(),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&automationaccount.AutomationAccountId{}),
		},
	}
}
