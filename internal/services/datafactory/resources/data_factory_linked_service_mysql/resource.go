// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_mysql

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataFactoryLinkedServiceMySQLCreateUpdate,
		Read:   resourceDataFactoryLinkedServiceMySQLRead,
		Update: resourceDataFactoryLinkedServiceMySQLCreateUpdate,
		Delete: resourceDataFactoryLinkedServiceMySQLDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.LinkedServiceID(id)
			return err
		}, helpers.ImportDataFactoryLinkedService(datafactory.TypeBasicLinkedServiceTypeMySQL)),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: dataFactoryLinkedServiceMysqlSchema(),

		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, i any) error {
			// No state yet, a new resource being created.
			if d.GetRawState().IsNull() {
				if d.Get("driver_version") == "V1" {
					return errors.New("`driver_version` must be set to `V2` for new resources")
				}
			}
			return nil
		}),
	}
}
