package helper

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/recordsets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func ResourceDnsRecordImporter(d *pluginsdk.ResourceData, recordType recordsets.RecordType) ([]*pluginsdk.ResourceData, error) {
	resourceId, err := recordsets.ParseRecordTypeID(d.Id())
	if err != nil {
		return []*pluginsdk.ResourceData{d}, err
	}
	if resourceId.RecordType != recordType {
		return []*pluginsdk.ResourceData{d}, fmt.Errorf("importing %s wrong type received: expected %s received %s", resourceId, recordType, resourceId.RecordType)
	}
	return []*pluginsdk.ResourceData{d}, nil
}
