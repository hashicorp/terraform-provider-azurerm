// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_custom_service

import (
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataFactoryLinkedCustomServiceCreateUpdate,
		Read:   resourceDataFactoryLinkedCustomServiceRead,
		Update: resourceDataFactoryLinkedCustomServiceCreateUpdate,
		Delete: resourceDataFactoryLinkedCustomServiceDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.LinkedServiceID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: dataFactoryLinkedCustomServiceSchema(),
	}
}
