package azure

import (
	"fmt"
)

func ValidateResourceID(i interface{}, k string) (ws []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := ParseAzureResourceID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a resource id: %v", k, err))
	}

	return ws, errors
}

//true for a resource ID or an empty string
func ValidateResourceIDOrEmpty(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if v == "" {
		return
	}

	return ValidateResourceID(i, k)
}

func ValidateMsSqlElasticPoolName(i interface{}, k string) (_ []string, errors []error) {
	return ValidateNameGeneric(i, k, alphanumericLower, hyphen, "azurerm_mssql_elasticpool", 3, 50)
}

func ValidateSqlServerName(i interface{}, k string) (_ []string, errors []error) {
	return ValidateNameGeneric(i, k, alphanumericLower, hyphen, "azurerm_sql_server", 3, 50)
}

func ValidateMySqlServerName(i interface{}, k string) (_ []string, errors []error) {
	return ValidateNameGeneric(i, k, alphanumericLower, hyphen, "azurerm_mysql_server", 3, 50)
}

func ValidatePostgreSqlServerName(i interface{}, k string) (_ []string, errors []error) {
	return ValidateNameGeneric(i, k, alphanumericLower, hyphen, "azurerm_postgresql_server", 3, 50)
}

func ValidateAvailabilitySet(i interface{}, k string) (_ []string, errors []error) {
	return ValidateNameGeneric(i, k, alphanumericBoth, hyphenUnderscore, "", 3, 80)
}

func ValidateVirtualMachineWindows(i interface{}, k string) (_ []string, errors []error) {
	return ValidateNameGeneric(i, k, alphanumericBoth, hyphen, "", 3, 15)
}

func ValidateVirtualMachineLinux(i interface{}, k string) (_ []string, errors []error) {
	return ValidateNameGeneric(i, k, alphanumericBoth, hyphen, "", 3, 64)
}

func ValidateFunctionAppName(i interface{}, k string) (_ []string, errors []error) {
	return ValidateNameGeneric(i, k, alphanumericBoth, hyphen, "", 3, 60)
}

func ValidateStorageAccountName(i interface{}, k string) (_ []string, errors []error) {
	return ValidateNameGeneric(i, k, alphanumericLower, none, "", 3, 24)
}

func ValidateContainerName(i interface{}, k string) (_ []string, errors []error) {
	return ValidateNameGeneric(i, k, alphanumericBoth, hyphen, "", 3, 63)
}

func ValidateQueueName(i interface{}, k string) (_ []string, errors []error) {
	return ValidateNameGeneric(i, k, alphanumericLower, none, "", 3, 63)
}

func ValidateTableName(i interface{}, k string) (_ []string, errors []error) {
	return ValidateNameGeneric(i, k, alphanumericBoth, none, "", 3, 63)
}

func ValidateDataLakeStoreName(i interface{}, k string) (_ []string, errors []error) {
	return ValidateNameGeneric(i, k, alphanumericLower, none, "", 3, 24)
}

func ValidateVirtualNetworkName(i interface{}, k string) (_ []string, errors []error) {
	return ValidateNameGeneric(i, k, alphanumericBoth, hyphenUnderscorePeriod, "", 3, 64)
}

func ValidateSubnetName(i interface{}, k string) (_ []string, errors []error) {
	return ValidateNameGeneric(i, k, alphanumericBoth, hyphenUnderscorePeriod, "", 3, 80)
}

func ValidateSecurityGroupName(i interface{}, k string) (_ []string, errors []error) {
	return ValidateNameGeneric(i, k, alphanumericBoth, hyphenUnderscorePeriod, "", 3, 80)
}

func ValidateSecurityGroupRuleName(i interface{}, k string) (_ []string, errors []error) {
	return ValidateNameGeneric(i, k, alphanumericBoth, hyphenUnderscorePeriod, "", 3, 80)
}

func ValidateLoadBalancerName(i interface{}, k string) (_ []string, errors []error) {
	return ValidateNameGeneric(i, k, alphanumericBoth, hyphenUnderscorePeriod, "", 3, 80)
}

func ValidateLoadBalancerRulesName(i interface{}, k string) (_ []string, errors []error) {
	return ValidateNameGeneric(i, k, alphanumericBoth, hyphenUnderscorePeriod, "", 3, 80)
}

func ValidateTrafficManagerName(i interface{}, k string) (_ []string, errors []error) {
	return ValidateNameGeneric(i, k, alphanumericBoth, hyphenPeriod, "", 3, 80)
}

func ValidateContainerRegistryName(i interface{}, k string) (_ []string, errors []error) {
	return ValidateNameGeneric(i, k, alphanumericBoth, none, "", 5, 50)
}
