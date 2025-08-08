package validate

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-03-01/autonomousdatabases"
)

func ValidateCloneWorkloadType(i interface{}, k string) (warnings []string, errors []error) {
	workload := i.(string)

	validWorkloads := []string{
		string(autonomousdatabases.WorkloadTypeDW),
		string(autonomousdatabases.WorkloadTypeOLTP),
		string(autonomousdatabases.WorkloadTypeAJD),
		string(autonomousdatabases.WorkloadTypeAPEX),
	}

	isValid := false
	for _, valid := range validWorkloads {
		if workload == valid {
			isValid = true
			break
		}
	}

	if !isValid {
		errors = append(errors, fmt.Errorf("%q must be one of %v, got: %s", k, validWorkloads, workload))
		return warnings, errors
	}

	warnings = append(warnings, "Note: Cross-workload cloning restrictions apply - refreshable clones cannot change workload type from source database")

	return warnings, errors
}
