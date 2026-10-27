// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dns_srv_record

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/recordsets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dns/resources/dns_srv_record/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDnsSrvRecordCreate,
		Read:   resourceDnsSrvRecordRead,
		Update: resourceDnsSrvRecordUpdate,
		Delete: resourceDnsSrvRecordDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			parsed, err := recordsets.ParseRecordTypeID(id)
			if err != nil {
				return err
			}
			if parsed.RecordType != recordsets.RecordTypeSRV {
				return fmt.Errorf("this resource only supports 'SRV' records")
			}
			return nil
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SRVRecordV0ToV1{},
		}),

		Schema: dnsSrvRecordSchema(),
	}
}
