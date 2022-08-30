package cosmosdb

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = UserDefinedFunctionId{}

// UserDefinedFunctionId is a struct representing the Resource ID for a User Defined Function
type UserDefinedFunctionId struct {
	SubscriptionId          string
	ResourceGroupName       string
	AccountName             string
	DatabaseName            string
	ContainerName           string
	UserDefinedFunctionName string
}

// NewUserDefinedFunctionID returns a new UserDefinedFunctionId struct
func NewUserDefinedFunctionID(subscriptionId string, resourceGroupName string, accountName string, databaseName string, containerName string, userDefinedFunctionName string) UserDefinedFunctionId {
	return UserDefinedFunctionId{
		SubscriptionId:          subscriptionId,
		ResourceGroupName:       resourceGroupName,
		AccountName:             accountName,
		DatabaseName:            databaseName,
		ContainerName:           containerName,
		UserDefinedFunctionName: userDefinedFunctionName,
	}
}

// ParseUserDefinedFunctionID parses 'input' into a UserDefinedFunctionId
func ParseUserDefinedFunctionID(input string) (*UserDefinedFunctionId, error) {
	parser := resourceids.NewParserFromResourceIdType(UserDefinedFunctionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := UserDefinedFunctionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.DatabaseName, ok = parsed.Parsed["databaseName"]; !ok {
		return nil, fmt.Errorf("the segment 'databaseName' was not found in the resource id %q", input)
	}

	if id.ContainerName, ok = parsed.Parsed["containerName"]; !ok {
		return nil, fmt.Errorf("the segment 'containerName' was not found in the resource id %q", input)
	}

	if id.UserDefinedFunctionName, ok = parsed.Parsed["userDefinedFunctionName"]; !ok {
		return nil, fmt.Errorf("the segment 'userDefinedFunctionName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseUserDefinedFunctionIDInsensitively parses 'input' case-insensitively into a UserDefinedFunctionId
// note: this method should only be used for API response data and not user input
func ParseUserDefinedFunctionIDInsensitively(input string) (*UserDefinedFunctionId, error) {
	parser := resourceids.NewParserFromResourceIdType(UserDefinedFunctionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := UserDefinedFunctionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.DatabaseName, ok = parsed.Parsed["databaseName"]; !ok {
		return nil, fmt.Errorf("the segment 'databaseName' was not found in the resource id %q", input)
	}

	if id.ContainerName, ok = parsed.Parsed["containerName"]; !ok {
		return nil, fmt.Errorf("the segment 'containerName' was not found in the resource id %q", input)
	}

	if id.UserDefinedFunctionName, ok = parsed.Parsed["userDefinedFunctionName"]; !ok {
		return nil, fmt.Errorf("the segment 'userDefinedFunctionName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateUserDefinedFunctionID checks that 'input' can be parsed as a User Defined Function ID
func ValidateUserDefinedFunctionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseUserDefinedFunctionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted User Defined Function ID
func (id UserDefinedFunctionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/databaseAccounts/%s/sqlDatabases/%s/containers/%s/userDefinedFunctions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.DatabaseName, id.ContainerName, id.UserDefinedFunctionName)
}

// Segments returns a slice of Resource ID Segments which comprise this User Defined Function ID
func (id UserDefinedFunctionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDocumentDB", "Microsoft.DocumentDB", "Microsoft.DocumentDB"),
		resourceids.StaticSegment("staticDatabaseAccounts", "databaseAccounts", "databaseAccounts"),
		resourceids.UserSpecifiedSegment("accountName", "accountValue"),
		resourceids.StaticSegment("staticSqlDatabases", "sqlDatabases", "sqlDatabases"),
		resourceids.UserSpecifiedSegment("databaseName", "databaseValue"),
		resourceids.StaticSegment("staticContainers", "containers", "containers"),
		resourceids.UserSpecifiedSegment("containerName", "containerValue"),
		resourceids.StaticSegment("staticUserDefinedFunctions", "userDefinedFunctions", "userDefinedFunctions"),
		resourceids.UserSpecifiedSegment("userDefinedFunctionName", "userDefinedFunctionValue"),
	}
}

// String returns a human-readable description of this User Defined Function ID
func (id UserDefinedFunctionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Account Name: %q", id.AccountName),
		fmt.Sprintf("Database Name: %q", id.DatabaseName),
		fmt.Sprintf("Container Name: %q", id.ContainerName),
		fmt.Sprintf("User Defined Function Name: %q", id.UserDefinedFunctionName),
	}
	return fmt.Sprintf("User Defined Function (%s)", strings.Join(components, "\n"))
}
