// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_sql_pool

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSynapseSqlPoolCreate,
		Read:   resourceSynapseSqlPoolRead,
		Update: resourceSynapseSqlPoolUpdate,
		Delete: resourceSynapseSqlPoolDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.SqlPoolID(id)
			return err
		}, func(ctx context.Context, d *pluginsdk.ResourceData, meta any) ([]*pluginsdk.ResourceData, error) {
			d.Set("create_mode", DefaultCreateMode)
			if v, ok := d.GetOk("create_mode"); ok && v.(string) != "" {
				d.Set("create_mode", v)
			}

			return []*pluginsdk.ResourceData{d}, nil
		}),

		Schema: synapseSqlPoolSchema(),

		CustomizeDiff: pluginsdk.CustomizeDiffShim(synapseSqlPoolCustomizeDiff),
	}
}
