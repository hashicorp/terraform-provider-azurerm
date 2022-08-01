package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type DatadogSingleSignOnConfigurationsId struct {
	SubscriptionId                string
	ResourceGroup                 string
	MonitorName                   string
	SingleSignOnConfigurationName string
}

func NewDatadogSingleSignOnConfigurationsID(subscriptionId, resourceGroup, monitorName, singleSignOnConfigurationName string) DatadogSingleSignOnConfigurationsId {
	return DatadogSingleSignOnConfigurationsId{
		SubscriptionId:                subscriptionId,
		ResourceGroup:                 resourceGroup,
		MonitorName:                   monitorName,
		SingleSignOnConfigurationName: singleSignOnConfigurationName,
	}
}

func (id DatadogSingleSignOnConfigurationsId) String() string {
	segments := []string{
		fmt.Sprintf("Single Sign On Configuration Name %q", id.SingleSignOnConfigurationName),
		fmt.Sprintf("Monitor Name %q", id.MonitorName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Datadog Single Sign On Configurations", segmentsStr)
}

func (id DatadogSingleSignOnConfigurationsId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Datadog/monitors/%s/singleSignOnConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.MonitorName, id.SingleSignOnConfigurationName)
}

// DatadogSingleSignOnConfigurationsID parses a DatadogSingleSignOnConfigurations ID into an DatadogSingleSignOnConfigurationsId struct
func DatadogSingleSignOnConfigurationsID(input string) (*DatadogSingleSignOnConfigurationsId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := DatadogSingleSignOnConfigurationsId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.MonitorName, err = id.PopSegment("monitors"); err != nil {
		return nil, err
	}
	if resourceId.SingleSignOnConfigurationName, err = id.PopSegment("singleSignOnConfigurations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
