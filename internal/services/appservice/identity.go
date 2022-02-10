package appservice

import "github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web"

func expandIdentity(input []interface{}) (*web.ManagedServiceIdentity, error) {
	return nil, nil
}

func flattenIdentity(identity *web.ManagedServiceIdentity) (*[]interface{}, error) {
	return nil, nil
}
