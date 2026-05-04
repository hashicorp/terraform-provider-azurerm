package filesystems

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
)

// GetResourceManagerResourceID returns the Resource Manager specific
// ResourceID for a specific Data Lake FileSystem
func (c Client) GetResourceManagerResourceID(subscriptionID, resourceGroup, accountName, fileSystemName string) string {
	fmtStr := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/blobServices/default/containers/%s"
	return fmt.Sprintf(fmtStr, subscriptionID, resourceGroup, accountName, fileSystemName)
}

// TODO: update this to implement `resourceids.ResourceId` once
// https://github.com/hashicorp/go-azure-helpers/issues/187 is fixed
var _ resourceids.Id = FileSystemId{}

type FileSystemId struct {
	// AccountId specifies the ID of the Storage Account where this Data Lake FileSystem exists.
	AccountId accounts.AccountId

	// FileSystemName specifies the name of this Data Lake FileSystem.
	FileSystemName string
}

func NewFileSystemID(accountId accounts.AccountId, fileSystemName string) FileSystemId {
	return FileSystemId{
		AccountId:      accountId,
		FileSystemName: fileSystemName,
	}
}

func (b FileSystemId) ID() string {
	return fmt.Sprintf("%s/%s", b.AccountId.ID(), b.FileSystemName)
}

func (b FileSystemId) String() string {
	components := []string{
		fmt.Sprintf("Account %q", b.AccountId.String()),
	}
	return fmt.Sprintf("File System %q (%s)", b.FileSystemName, strings.Join(components, " / "))
}

// ParseFileSystemID parses `input` into a File System ID using a known `domainSuffix`
func ParseFileSystemID(input, domainSuffix string) (*FileSystemId, error) {
	// example: https://foo.dfs.core.windows.net/Bar
	if input == "" {
		return nil, fmt.Errorf("`input` was empty")
	}

	account, err := accounts.ParseAccountID(input, domainSuffix)
	if err != nil {
		return nil, fmt.Errorf("parsing account %q: %+v", input, err)
	}

	if account.SubDomainType != accounts.DataLakeStoreSubDomainType {
		return nil, fmt.Errorf("expected the subdomain type to be %q but got %q", string(accounts.DataLakeStoreSubDomainType), string(account.SubDomainType))
	}

	uri, err := url.Parse(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as a uri: %+v", input, err)
	}

	path := strings.TrimPrefix(uri.Path, "/")
	segments := strings.Split(path, "/")
	if len(segments) != 1 {
		return nil, fmt.Errorf("expected the path to contain 1 segment but got %d", len(segments))
	}

	fileSystemName := segments[0]
	return &FileSystemId{
		AccountId:      *account,
		FileSystemName: fileSystemName,
	}, nil
}
