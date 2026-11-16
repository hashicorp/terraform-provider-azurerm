// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package active_directory_domain_service_replica_set

import (
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/domainservices/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

const DomainServiceResourceName = "azurerm_active_directory_domain_service"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceActiveDirectoryDomainServiceReplicaSetCreate,
		Read:   resourceActiveDirectoryDomainServiceReplicaSetRead,
		Delete: resourceActiveDirectoryDomainServiceReplicaSetDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(3 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(1 * time.Hour),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.DomainServiceReplicaSetID(id)
			return err
		}),

		Schema: activeDirectoryDomainServiceReplicaSetSchema(),
	}
}
