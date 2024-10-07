// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

// TODO: @tombuildsstuff: this wants refactoring and fixing into sub-ID parsers

type RoleAssignmentId struct {
	SubscriptionID           string
	ResourceGroup            string
	ManagementGroup          string
	ResourceScope            string
	ResourceProvider         string
	Name                     string
	SubscriptionAlias        string
	TenantId                 string
	IsSubscriptionLevel      bool
	IsSubscriptionAliasLevel bool
}

func NewRoleAssignmentID(subscriptionId, resourceGroup, resourceProvider, resourceScope, managementGroup, name, tenantId, subscriptionAlias string, isSubLevel bool, isSubAliasLevel bool) (*RoleAssignmentId, error) {
	if subscriptionId == "" && resourceGroup == "" && managementGroup == "" && !isSubLevel && !isSubAliasLevel {
		return nil, fmt.Errorf("one of subscriptionId, resourceGroup, managementGroup, isSubscriptionLevel or isSubscriptionAliasLevel must be provided")
	}

	if managementGroup != "" {
		if subscriptionId != "" || resourceGroup != "" || isSubLevel {
			return nil, fmt.Errorf("cannot provide subscriptionId, resourceGroup or isSubscriptionLevel when managementGroup is provided")
		}
	}

	if isSubLevel {
		if subscriptionId != "" || resourceGroup != "" || managementGroup != "" {
			return nil, fmt.Errorf("cannot provide subscriptionId, resourceGroup or managementGroup when isSubscriptionLevel is provided")
		}
	}

	if isSubAliasLevel {
		if subscriptionId != "" || resourceGroup != "" || managementGroup != "" {
			return nil, fmt.Errorf("cannot provide subscriptionId, resourceGroup or managementGroup when isSubscriptionAliasLevel is provided")
		}
	}

	if resourceGroup != "" {
		if subscriptionId == "" {
			return nil, fmt.Errorf("subscriptionId must not be empty when resourceGroup is provided")
		}
	}

	return &RoleAssignmentId{
		SubscriptionID:           subscriptionId,
		ResourceGroup:            resourceGroup,
		ResourceProvider:         resourceProvider,
		ResourceScope:            resourceScope,
		ManagementGroup:          managementGroup,
		SubscriptionAlias:        subscriptionAlias,
		Name:                     name,
		TenantId:                 tenantId,
		IsSubscriptionLevel:      isSubLevel,
		IsSubscriptionAliasLevel: isSubAliasLevel,
	}, nil
}

// in general case, the id format does not change
// for cross tenant scenario, add the tenantId info
func (id RoleAssignmentId) AzureResourceID() string {
	if id.ResourceScope != "" {
		fmtString := "/subscriptions/%s/resourceGroups/%s/providers/%s/%s/providers/Microsoft.Authorization/roleAssignments/%s"
		return fmt.Sprintf(fmtString, id.SubscriptionID, id.ResourceGroup, id.ResourceProvider, id.ResourceScope, id.Name)
	}

	if id.ManagementGroup != "" {
		fmtString := "/providers/Microsoft.Management/managementGroups/%s/providers/Microsoft.Authorization/roleAssignments/%s"
		return fmt.Sprintf(fmtString, id.ManagementGroup, id.Name)
	}

	if id.ResourceGroup != "" {
		fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Authorization/roleAssignments/%s"
		return fmt.Sprintf(fmtString, id.SubscriptionID, id.ResourceGroup, id.Name)
	}

	if id.IsSubscriptionAliasLevel {
		fmtString := "/providers/Microsoft.Subscription/aliases/%s/providers/Microsoft.Authorization/roleAssignments/%s"
		return fmt.Sprintf(fmtString, id.SubscriptionAlias, id.Name)
	}

	if id.IsSubscriptionLevel {
		fmtString := "/providers/Microsoft.Subscription/providers/Microsoft.Authorization/roleAssignments/%s"
		return fmt.Sprintf(fmtString, id.Name)
	}

	if id.SubscriptionID != "" {
		fmtString := "/subscriptions/%s/providers/Microsoft.Authorization/roleAssignments/%s"
		return fmt.Sprintf(fmtString, id.SubscriptionID, id.Name)
	}

	if id.ResourceProvider != "" {
		fmtString := "/providers/%s/providers/Microsoft.Authorization/roleAssignments/%s"
		return fmt.Sprintf(fmtString, id.ResourceProvider, id.Name)
	}

	fmtString := "/providers/Microsoft.Authorization/roleAssignments/%s"
	return fmt.Sprintf(fmtString, id.Name)
}

func (id RoleAssignmentId) ID() string {
	return ConstructRoleAssignmentId(id.AzureResourceID(), id.TenantId)
}

func ConstructRoleAssignmentId(azureResourceId, tenantId string) string {
	if tenantId == "" {
		return azureResourceId
	}
	return fmt.Sprintf("%s|%s", azureResourceId, tenantId)
}

func DestructRoleAssignmentId(id string) (string, string) {
	parts := strings.Split(id, "|")
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return id, ""
}

func RoleAssignmentID(input string) (*RoleAssignmentId, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("Role Assignment ID is empty string")
	}

	roleAssignmentId := RoleAssignmentId{}

	parts := strings.Split(input, "|")
	if len(parts) == 2 {
		roleAssignmentId.TenantId = parts[1]
		input = parts[0]
	}

	switch {
	case strings.HasPrefix(input, "/subscriptions/"):
		id, err := azure.ParseAzureResourceID(input)
		if err != nil {
			return nil, fmt.Errorf("could not parse %q as Azure resource ID", input)
		}
		roleAssignmentId.SubscriptionID = id.SubscriptionID
		roleAssignmentId.ResourceGroup = id.ResourceGroup
		if id.Provider != "Microsoft.Authorization" && id.Provider != "" {
			roleAssignmentId.ResourceProvider = id.Provider
			// logic to save resource scope
			result := strings.Split(input, "/providers/")
			if len(result) == 3 {
				roleAssignmentId.ResourceScope = strings.TrimPrefix(result[1], fmt.Sprintf("%s/", id.Provider))
			}
		}

		if roleAssignmentId.Name, err = id.PopSegment("roleAssignments"); err != nil {
			return nil, err
		}
	case strings.HasPrefix(input, "/providers/Microsoft.Subscription/"):
		idParts := strings.Split(input, "/providers/Microsoft.Authorization/roleAssignments/")
		if len(idParts) != 2 {
			return nil, fmt.Errorf("could not parse Role Assignment ID %q for subscription scope", input)
		}
		if strings.Contains(input, "/aliases/") {
			roleAssignmentId.IsSubscriptionAliasLevel = true
			aliasParts := strings.Split(idParts[0], "/")
			alias := aliasParts[len(aliasParts)-1]
			roleAssignmentId.SubscriptionAlias = alias
		} else {
			roleAssignmentId.IsSubscriptionLevel = true
		}
		if idParts[1] == "" {
			return nil, fmt.Errorf("ID was missing a value for the roleAssignments element")
		}
		roleAssignmentId.Name = idParts[1]
	case strings.HasPrefix(input, "/providers/Microsoft.Management/"):
		idParts := strings.Split(input, "/providers/Microsoft.Authorization/roleAssignments/")
		if len(idParts) != 2 {
			return nil, fmt.Errorf("could not parse Role Assignment ID %q for Management Group", input)
		}
		if idParts[1] == "" {
			return nil, fmt.Errorf("ID was missing a value for the roleAssignments element")
		}
		roleAssignmentId.Name = idParts[1]
		roleAssignmentId.ManagementGroup = strings.TrimPrefix(idParts[0], "/providers/Microsoft.Management/managementGroups/")
	case strings.HasPrefix(input, "/providers/") && !strings.HasPrefix(input, "/providers/Microsoft.Authorization/roleAssignments"):
		idParts := strings.Split(input, "/providers/Microsoft.Authorization/roleAssignments/")
		if len(idParts) != 2 {
			return nil, fmt.Errorf("could not parse Role Assignment ID %q for Resource Provider", input)
		}
		if idParts[1] == "" {
			return nil, fmt.Errorf("ID was missing a value for the roleAssignments element")
		}
		roleAssignmentId.Name = idParts[1]
		roleAssignmentId.ResourceProvider = strings.TrimPrefix(idParts[0], "/providers/")
	case strings.HasPrefix(input, "/providers/Microsoft.Authorization/roleAssignments"):
		name := strings.TrimPrefix(input, "/providers/Microsoft.Authorization/roleAssignments/")
		if name == "" {
			return nil, fmt.Errorf("ID was missing a value for the roleAssignments element")
		}
		roleAssignmentId.Name = name
	default:
		return nil, fmt.Errorf("could not parse Role Assignment ID %q", input)
	}

	return &roleAssignmentId, nil
}
