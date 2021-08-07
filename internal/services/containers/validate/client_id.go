package validate

import (
	"fmt"
	"strings"
)

// ClientID validates the ClientID is valid for a Kubernetes Cluster
func ClientID(i interface{}, k string) ([]string, []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	v = strings.TrimSpace(v)
	if v == "" {
		return nil, []error{fmt.Errorf("expected %q to not be an empty string, got %v", k, i)}
	}

	// whilst `msi` is valid in ARM, it doesn't make sense in Terraform
	// since we can instead omit the `service_principal` block
	if strings.EqualFold(v, "msi") {
		return nil, []error{
			fmt.Errorf("to define an AKS cluster with authentication via MSI - remove the `service_principal` block"),
		}
	}

	return nil, nil
}
