package cosmos

import "github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-10-15/documentdb" // nolint: staticcheck

func isServerlessCapacityMode(accResp documentdb.DatabaseAccountGetResults) bool {
	if props := accResp.DatabaseAccountGetProperties; props != nil && props.Capabilities != nil {
		for _, v := range *props.Capabilities {
			if v.Name != nil && *v.Name == "EnableServerless" {
				return true
			}
		}
	}

	return false
}
