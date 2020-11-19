package validate

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
)

func LogAnalyticsClustersName(i interface{}, k string) (warnings []string, errors []error) {
	return logAnalyticsGenericName(i, k)
}

func LogAnalyticsClusterId(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	_, err := parse.LogAnalyticsClusterID(v)
	if err != nil {
		errors = append(errors, fmt.Errorf("expected %s to be a Log Analytics Cluster ID:, %+v", k, err))
	}

	return warnings, errors
}
