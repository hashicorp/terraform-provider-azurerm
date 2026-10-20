// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machine_learning_workspace

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Encryption struct {
	IdentityClientID string `tfschema:"user_assigned_identity_id"`
	KeyVaultID       string `tfschema:"key_vault_id"`
	KeyID            string `tfschema:"key_id"`
}

type ManagedNetwork struct {
	IsolationMode string `tfschema:"isolation_mode"`
}

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMachineLearningWorkspaceCreate,
		Read:   resourceMachineLearningWorkspaceRead,
		Update: resourceMachineLearningWorkspaceUpdate,
		Delete: resourceMachineLearningWorkspaceDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := workspaces.ParseWorkspaceID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: machineLearningWorkspaceSchema(),
	}
}
