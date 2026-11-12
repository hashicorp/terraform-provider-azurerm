// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_odbc

import (
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataFactoryLinkedServiceOdbcCreateUpdate,
		Read:   resourceDataFactoryLinkedServiceOdbcRead,
		Update: resourceDataFactoryLinkedServiceOdbcCreateUpdate,
		Delete: resourceDataFactoryLinkedServiceOdbcDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.LinkedServiceID(id)
			return err
		}, helpers.ImportDataFactoryLinkedService(datafactory.TypeBasicLinkedServiceTypeOdbc)),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: dataFactoryLinkedServiceOdbcSchema(),
	}
}
