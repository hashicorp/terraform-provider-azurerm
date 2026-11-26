// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package private_dns_mx_record

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatedns"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePrivateDnsMxRecordCreateUpdate,
		Read:   resourcePrivateDnsMxRecordRead,
		Update: resourcePrivateDnsMxRecordCreateUpdate,
		Delete: resourcePrivateDnsMxRecordDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			resourceId, err := privatedns.ParseRecordTypeID(id)
			if err != nil {
				return err
			}
			if resourceId.RecordType != privatedns.RecordTypeMX {
				return fmt.Errorf("importing %s wrong type received: expected %s received %s", id, privatedns.RecordTypeMX, resourceId.RecordType)
			}
			return nil
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: privateDnsMxRecordSchema(),
	}
}
