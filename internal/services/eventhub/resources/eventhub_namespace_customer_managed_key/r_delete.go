// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventhub_namespace_customer_managed_key

import (
	"log"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceEventHubNamespaceCustomerManagedKeyDelete(d *pluginsdk.ResourceData, meta any) error {
	log.Printf(`[INFO] Customer Managed Keys cannot be removed from EventHub Namespaces once added. To remove the Customer Managed Key delete and recreate the parent EventHub Namespace`)
	return nil
}
