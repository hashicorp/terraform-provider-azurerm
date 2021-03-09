package validate

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
)

func AdminUsernames(i interface{}, k string) (_ []string, errors []error) {
	disallowedLogins := [7]string{"azure_superuser", "azure_pg_admin", "admin", "administrator", "root", "guest", "public"}
	for _, v := range disallowedLogins {
		if v == k {
			return fmt.Errorf("Error - PostgreSQL AD Administrator login can not be %q", k)
		}
	}
	if strings.HasPrefix(k, "pg_") {
		return fmt.Errorf("Error - PostgreSQL AD Administrator login can not start with 'pg_'")
	}

	return nil, nil
}
