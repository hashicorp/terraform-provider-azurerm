// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventgrid_domain

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/domains"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceEventGridDomainCreate,
		Read:   resourceEventGridDomainRead,
		Update: resourceEventGridDomainUpdate,
		Delete: resourceEventGridDomainDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := domains.ParseDomainID(id)
			return err
		}),

		Schema: eventgridDomainSchema(),
	}
}
