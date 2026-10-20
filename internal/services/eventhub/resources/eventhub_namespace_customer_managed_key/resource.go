// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventhub_namespace_customer_managed_key

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceEventHubNamespaceCustomerManagedKeyCreateUpdate,
		Read:   resourceEventHubNamespaceCustomerManagedKeyRead,
		Update: resourceEventHubNamespaceCustomerManagedKeyCreateUpdate,
		Delete: resourceEventHubNamespaceCustomerManagedKeyDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := namespaces.ParseNamespaceID(id)
			return err
		}, func(ctx context.Context, d *pluginsdk.ResourceData, meta any) ([]*pluginsdk.ResourceData, error) {
			client := meta.(*clients.Client).Eventhub.NamespacesClient

			var cancel context.CancelFunc
			ctx, cancel = timeouts.ForRead(ctx, d)
			defer cancel()

			id, err := namespaces.ParseNamespaceID(d.Id())
			if err != nil {
				return []*pluginsdk.ResourceData{d}, err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return []*pluginsdk.ResourceData{d}, fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.Encryption == nil {
				return []*pluginsdk.ResourceData{d}, fmt.Errorf("retrieving %s: no customer managed key present", *id)
			}

			return []*pluginsdk.ResourceData{d}, nil
		}),

		Schema: eventhubNamespaceCustomerManagedKeySchema(),
	}
}
