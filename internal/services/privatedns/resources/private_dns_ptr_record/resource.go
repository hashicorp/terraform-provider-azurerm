// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package private_dns_ptr_record

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatedns"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePrivateDnsPtrRecordCreateUpdate,
		Read:   resourcePrivateDnsPtrRecordRead,
		Update: resourcePrivateDnsPtrRecordCreateUpdate,
		Delete: resourcePrivateDnsPtrRecordDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			resourceId, err := privatedns.ParseRecordTypeID(id)
			if err != nil {
				return err
			}
			if resourceId.RecordType != privatedns.RecordTypePTR {
				return fmt.Errorf("importing %s wrong type received: expected %s received %s", id, privatedns.RecordTypePTR, resourceId.RecordType)
			}
			return nil
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: privateDnsPtrRecordSchema(),
	}
}
