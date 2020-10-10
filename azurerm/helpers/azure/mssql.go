package azure

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
)

// Your server name can contain only lowercase letters, numbers, and '-', but can't start or end with '-' or have more than 63 characters.
func ValidateMsSqlServerName(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `^[0-9a-z]([-0-9a-z]{0,61}[0-9a-z])?$`); !m {
		return nil, append(regexErrs, fmt.Errorf("%q can contain only lowercase letters, numbers, and '-', but can't start or end with '-' or have more than 63 characters.", k))
	}

	return nil, nil
}

// Your database name can't end with '.' or ' ', can't contain '<,>,*,%,&,:,\,/,?' or control characters, and can't have more than 128 characters.
func ValidateMsSqlDatabaseName(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `^[^<>*%&:\\\/?]{0,127}[^\s.<>*%&:\\\/?]$`); !m {
		return nil, append(regexErrs, fmt.Errorf(`%q can't end with '.' or ' ', can't contain '<,>,*,%%,&,:,\,/,?' or control characters, and can't have more than 128 characters.`, k))
	}

	return nil, nil
}

func ValidateMsSqlFailoverGroupName(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `^[0-9a-z]([-0-9a-z]{0,61}[0-9a-z])?$`); !m {
		return nil, append(regexErrs, fmt.Errorf("%q can contain only lowercase letters, numbers, and '-', but can't start or end with '-'.", k))
	}

	return nil, nil
}

// Following characters and any control characters are not allowed for resource name '%,&,\\\\,?,/'.\"
// The name can not end with characters: '. '
// TODO: unsure about length, was able to deploy one at 120
func ValidateMsSqlElasticPoolName(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `^[^&%\\\/?]{0,127}[^\s.&%\\\/?]$`); !m {
		return nil, append(regexErrs, fmt.Errorf(`%q can't end with '.' or ' ', can't contain '%%,&,\,/,?' or control characters, and can't have more than 128 characters.`, k))
	}

	return nil, nil
}

// Managed instance name can only be made up of lowercase letters 'a'-'z', the numbers 0-9 and the hyphen.
// The hyphen may not lead or trail in the instance name.
// Maximum length of instace name can not exceed 260 characters
func ValidateManagedInstanceName(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `^[0-9a-z]([-0-9a-z]{0,259}[0-9a-z])?$`); !m {
		return nil, append(regexErrs, fmt.Errorf("%q can contain only lowercase letters, numbers, and '-', but can't start or end with '-' or have more than 260 characters.", k))
	}

	return nil, nil
}

func ValidateManagedInstanceTimeZones() schema.SchemaValidateFunc {
	// accepted timezones are listed here: https://docs.microsoft.com/en-us/azure/azure-sql/managed-instance/timezones-overview
	acceptedTimeZones := []string{
		"Dateline Standard Time",
		"UTC-11",
		"Aleutian Standard Time",
		"Hawaiian Standard Time",
		"Marquesas Standard Time",
		"Alaskan Standard Time",
		"UTC-09",
		"Pacific Standard Time (Mexico)",
		"UTC-08",
		"Pacific Standard Time",
		"US Mountain Standard Time",
		"Mountain Standard Time (Mexico)",
		"Mountain Standard Time",
		"Central America Standard Time",
		"Central Standard Time",
		"Easter Island Standard Time",
		"Central Standard Time (Mexico)",
		"Canada Central Standard Time",
		"SA Pacific Standard Time",
		"Eastern Standard Time (Mexico)",
		"Eastern Standard Time",
		"Haiti Standard Time",
		"Cuba Standard Time",
		"US Eastern Standard Time",
		"Turks And Caicos Standard Time",
		"Paraguay Standard Time",
		"Atlantic Standard Time",
		"Venezuela Standard Time",
		"Central Brazilian Standard Time",
		"SA Western Standard Time",
		"Pacific SA Standard Time",
		"Newfoundland Standard Time",
		"Tocantins Standard Time",
		"E. South America Standard Time",
		"SA Eastern Standard Time",
		"Argentina Standard Time",
		"Greenland Standard Time",
		"Montevideo Standard Time",
		"Magallanes Standard Time",
		"Saint Pierre Standard Time",
		"Bahia Standard Time",
		"UTC-02",
		"Mid-Atlantic Standard Time",
		"Azores Standard Time",
		"Cape Verde Standard Time",
		"UTC",
		"GMT Standard Time",
		"Greenwich Standard Time",
		"W. Europe Standard Time",
		"Central Europe Standard Time",
		"Romance Standard Time",
		"Morocco Standard Time",
		"Sao Tome Standard Time",
		"Central European Standard Time",
		"W. Central Africa Standard Time",
		"Jordan Standard Time",
		"GTB Standard Time",
		"Middle East Standard Time",
		"Egypt Standard Time",
		"E. Europe Standard Time",
		"Syria Standard Time",
		"West Bank Standard Time",
		"South Africa Standard Time",
		"FLE Standard Time",
		"Israel Standard Time",
		"Kaliningrad Standard Time",
		"Sudan Standard Time",
		"Libya Standard Time",
		"Namibia Standard Time",
		"Arabic Standard Time",
		"Turkey Standard Time",
		"Arab Standard Time",
		"Belarus Standard Time",
		"Russian Standard Time",
		"E. Africa Standard Time",
		"Iran Standard Time",
		"Arabian Standard Time",
		"Astrakhan Standard Time",
		"Azerbaijan Standard Time",
		"Russia Time Zone 3",
		"Mauritius Standard Time",
		"Saratov Standard Time",
		"Georgian Standard Time",
		"Volgograd Standard Time",
		"Caucasus Standard Time",
		"Afghanistan Standard Time",
		"West Asia Standard Time",
		"Ekaterinburg Standard Time",
		"Pakistan Standard Time",
		"India Standard Time",
		"Sri Lanka Standard Time",
		"Nepal Standard Time",
		"Central Asia Standard Time",
		"Bangladesh Standard Time",
		"Omsk Standard Time",
		"Myanmar Standard Time",
		"SE Asia Standard Time",
		"Altai Standard Time",
		"W. Mongolia Standard Time",
		"North Asia Standard Time",
		"N. Central Asia Standard Time",
		"Tomsk Standard Time",
		"China Standard Time",
		"North Asia East Standard Time",
		"Singapore Standard Time",
		"W. Australia Standard Time",
		"Taipei Standard Time",
		"Ulaanbaatar Standard Time",
		"Aus Central W. Standard Time",
		"Transbaikal Standard Time",
		"Tokyo Standard Time",
		"North Korea Standard Time",
		"Korea Standard Time",
		"Yakutsk Standard Time",
		"Cen. Australia Standard Time",
		"AUS Central Standard Time",
		"E. Australia Standard Time",
		"AUS Eastern Standard Time",
		"West Pacific Standard Time",
		"Tasmania Standard Time",
		"Vladivostok Standard Time",
		"Lord Howe Standard Tim",
		"Bougainville Standard Time",
		"Russia Time Zone 10",
		"Magadan Standard Time",
		"Norfolk Standard Time",
		"Sakhalin Standard Time",
		"Central Pacific Standard Time",
		"Russia Time Zone 11",
		"New Zealand Standard Time",
		"UTC+12",
		"Fiji Standard Time",
		"Kamchatka Standard Time",
		"Chatham Islands Standard Time",
		"UTC+13",
		"Tonga Standard Time",
		"Samoa Standard Time",
		"Line Islands Standard Time",
	}
	return validation.StringInSlice(acceptedTimeZones, true)
}

func GetSQLResourceParentId(id string) (*string, error) {
	idURL, err := url.ParseRequestURI(id)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse Azure ID: %s", err)
	}

	path := idURL.Path

	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	components := strings.Split(path, "/")

	// We should have an even number of key-value pairs.
	if len(components)%2 != 0 {
		return nil, fmt.Errorf("The number of path segments is not divisible by 2 in %q", path)
	}

	var subscriptionID string
	var resourceGroup string
	var providers string
	var server string
	var managedInstance string

	// Put the constituent key-value pairs into a map
	for current := 0; current < len(components); current += 2 {
		key := components[current]
		value := components[current+1]

		switch key {
		case "subscriptions":
			subscriptionID = value
		case "resourceGroups":
			resourceGroup = value
		case "providers":
			providers = value
		case "managedInstances":
			managedInstance = value
		case "servers":
			server = value
		default:
			return nil, fmt.Errorf("Key/Value cannot be empty strings. Key: '%s', Value: '%s'", key, value)
		}
	}

	var databaseParentComponents = []string{"/subscriptions", subscriptionID, "resourceGroups", resourceGroup, "providers", providers, "managedInstances", managedInstance}
	if server != "" {
		databaseParentComponents[6] = "servers"
		databaseParentComponents[7] = server
	}
	var parentId = strings.Join(databaseParentComponents, "/")
	return &parentId, nil
}

func ValidateLongTermRetentionPoliciesIsoFormat(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `^P[0-9]*[YMWD]`); !m {
		return nil, append(regexErrs, fmt.Errorf(`%q has to be a valid Duration format, starting with "P" and ending with either of the letters "YMWD"`, k))
	}
	return nil, nil
}
