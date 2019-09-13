package azure

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/healthcareapis/mgmt/2018-08-20-preview/healthcareapis"
)

func FlattenHealthcareAccessPolicies(policies *[]healthcareapis.ServiceAccessPolicyEntry) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)

	if policies == nil {
		return result
	}

	for _, policy := range *policies {
		policyRaw := make(map[string]interface{})

		if objectId := policy.ObjectID; objectId != nil {
			policyRaw["object_id"] = *objectId
		}

		result = append(result, policyRaw)
	}

	return result
}
