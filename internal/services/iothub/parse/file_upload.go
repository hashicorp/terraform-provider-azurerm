package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type FileUploadId struct {
	SubscriptionId string
	ResourceGroup  string
	IotHubName     string
	FileUploadName string
}

func NewFileUploadID(subscriptionId, resourceGroup, iotHubName, fileUploadName string) FileUploadId {
	return FileUploadId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		IotHubName:     iotHubName,
		FileUploadName: fileUploadName,
	}
}

func (id FileUploadId) String() string {
	segments := []string{
		fmt.Sprintf("File Upload Name %q", id.FileUploadName),
		fmt.Sprintf("Iot Hub Name %q", id.IotHubName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "File Upload", segmentsStr)
}

func (id FileUploadId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Devices/IotHubs/%s/FileUpload/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.IotHubName, id.FileUploadName)
}

// FileUploadID parses a FileUpload ID into an FileUploadId struct
func FileUploadID(input string) (*FileUploadId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FileUploadId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.IotHubName, err = id.PopSegment("IotHubs"); err != nil {
		return nil, err
	}
	if resourceId.FileUploadName, err = id.PopSegment("FileUpload"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
