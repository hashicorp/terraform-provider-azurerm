// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_test_policy

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/policies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmDevTestPolicyDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DevTestLabs.PoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := policies.ParsePolicyID(d.Id())
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, *id, policies.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(read.HttpResponse) {
			// deleted outside of TF
			log.Printf("[DEBUG] %s was not found  - assuming removed!", *id)
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if _, err = client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return err
}
