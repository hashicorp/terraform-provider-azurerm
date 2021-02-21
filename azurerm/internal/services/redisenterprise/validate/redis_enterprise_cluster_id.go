package validate

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/redisenterprise/parse"
)

//RedisEnterpriseClusterID validates that the passed interface contains a valid Redis Enterprist Cluster ID
func RedisEnterpriseClusterID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := parse.RedisEnterpriseClusterID(v); err != nil {
		errors = append(errors, err)
	}

	return
}
