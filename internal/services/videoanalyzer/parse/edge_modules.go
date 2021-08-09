package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type EdgeModulesId struct {
	SubscriptionId    string
	ResourceGroup     string
	VideoAnalyzerName string
	EdgeModuleName    string
}

func NewEdgeModulesID(subscriptionId, resourceGroup, videoAnalyzerName, edgeModuleName string) EdgeModulesId {
	return EdgeModulesId{
		SubscriptionId:    subscriptionId,
		ResourceGroup:     resourceGroup,
		VideoAnalyzerName: videoAnalyzerName,
		EdgeModuleName:    edgeModuleName,
	}
}

func (id EdgeModulesId) String() string {
	segments := []string{
		fmt.Sprintf("Edge Module Name %q", id.EdgeModuleName),
		fmt.Sprintf("Video Analyzer Name %q", id.VideoAnalyzerName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Edge Modules", segmentsStr)
}

func (id EdgeModulesId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Media/videoAnalyzers/%s/edgeModules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VideoAnalyzerName, id.EdgeModuleName)
}

// EdgeModulesID parses a EdgeModules ID into an EdgeModulesId struct
func EdgeModulesID(input string) (*EdgeModulesId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := EdgeModulesId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.VideoAnalyzerName, err = id.PopSegment("videoAnalyzers"); err != nil {
		return nil, err
	}
	if resourceId.EdgeModuleName, err = id.PopSegment("edgeModules"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
