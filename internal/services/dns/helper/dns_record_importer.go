package helper

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/recordsets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func ResourceDnsRecordImporter(_ context.Context, d *pluginsdk.ResourceData, _ interface{}, RecordType recordsets.RecordType) ([]*pluginsdk.ResourceData, error) {
	resourceId, err := recordsets.ParseRecordTypeID(d.Id())
	if err != nil {
		return []*pluginsdk.ResourceData{d}, err
	}
	if resourceId.RecordType != RecordType {
		return []*pluginsdk.ResourceData{d}, fmt.Errorf("importing %s wrong type received: expected %s received %s", resourceId, RecordType, resourceId.RecordType)
	}
	return []*pluginsdk.ResourceData{d}, nil
}
