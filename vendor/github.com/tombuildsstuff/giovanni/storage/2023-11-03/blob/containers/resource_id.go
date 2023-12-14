package containers

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/blob/accounts"
)

// GetResourceManagerResourceID returns the Resource Manager specific
// ResourceID for a specific Storage Container
func (c Client) GetResourceManagerResourceID(subscriptionID, resourceGroup, accountName, containerName string) string {
	fmtStr := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/blobServices/default/containers/%s"
	return fmt.Sprintf(fmtStr, subscriptionID, resourceGroup, accountName, containerName)
}

// TODO: update this to implement `resourceids.ResourceId` once
// https://github.com/hashicorp/go-azure-helpers/issues/187 is fixed
var _ resourceids.Id = ContainerId{}

type ContainerId struct {
	// AccountId specifies the ID of the Storage Account where this Blob exists.
	AccountId accounts.AccountId

	// ContainerName specifies the name of the Container within this Storage Account where this
	// Blob exists.
	ContainerName string
}

func NewContainerID(accountId accounts.AccountId, containerName string) ContainerId {
	return ContainerId{
		AccountId:     accountId,
		ContainerName: containerName,
	}
}

func (b ContainerId) ID() string {
	return fmt.Sprintf("%s/%s", b.AccountId.ID(), b.ContainerName)
}

func (b ContainerId) String() string {
	components := []string{
		fmt.Sprintf("Account %q", b.AccountId.String()),
	}
	return fmt.Sprintf("Containerr %q (%s)", b.ContainerName, strings.Join(components, " / "))
}

// ParseContainerID parses `input` into a Blob ID using a known `domainSuffix`
func ParseContainerID(input, domainSuffix string) (*ContainerId, error) {
	// example: https://foo.blob.core.windows.net/Bar
	if input == "" {
		return nil, fmt.Errorf("`input` was empty")
	}

	account, err := accounts.ParseAccountID(input, domainSuffix)
	if err != nil {
		return nil, fmt.Errorf("parsing account %q: %+v", input, err)
	}

	uri, err := url.Parse(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as a uri: %+v", input, err)
	}

	path := strings.TrimPrefix(uri.Path, "/")
	segments := strings.Split(path, "/")
	if len(segments) == 0 {
		return nil, fmt.Errorf("Expected the path to contain segments but got none")
	}

	containerName := strings.TrimPrefix(uri.Path, "/")
	return &ContainerId{
		AccountId:     *account,
		ContainerName: containerName,
	}, nil
}
