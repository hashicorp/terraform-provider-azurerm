package directories

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
)

// TODO: update this to implement `resourceids.ResourceId` once
// https://github.com/hashicorp/go-azure-helpers/issues/187 is fixed
var _ resourceids.Id = DirectoryId{}

type DirectoryId struct {
	// AccountId specifies the ID of the Storage Account where this Directory exists.
	AccountId accounts.AccountId

	// ShareName specifies the name of the File Share containing this Directory.
	ShareName string

	// DirectoryPath specifies the path representing this Directory.
	DirectoryPath string
}

func NewDirectoryID(accountId accounts.AccountId, shareName, directoryPath string) DirectoryId {
	return DirectoryId{
		AccountId:     accountId,
		ShareName:     shareName,
		DirectoryPath: directoryPath,
	}
}

func (b DirectoryId) ID() string {
	return fmt.Sprintf("%s/%s/%s", b.AccountId.ID(), b.ShareName, b.DirectoryPath)
}

func (b DirectoryId) String() string {
	components := []string{
		fmt.Sprintf("Share Name %q", b.ShareName),
		fmt.Sprintf("Account %q", b.AccountId.String()),
	}
	return fmt.Sprintf("Directory Path %q (%s)", b.DirectoryPath, strings.Join(components, " / "))
}

// ParseDirectoryID parses `input` into a Directory ID using a known `domainSuffix`
func ParseDirectoryID(input, domainSuffix string) (*DirectoryId, error) {
	// example: https://foo.file.core.windows.net/Bar/some/directory
	if input == "" {
		return nil, fmt.Errorf("`input` was empty")
	}

	account, err := accounts.ParseAccountID(input, domainSuffix)
	if err != nil {
		return nil, fmt.Errorf("parsing account %q: %+v", input, err)
	}

	if account.SubDomainType != accounts.FileSubDomainType {
		return nil, fmt.Errorf("expected the subdomain type to be %q but got %q", string(accounts.FileSubDomainType), string(account.SubDomainType))
	}

	uri, err := url.Parse(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as a uri: %+v", input, err)
	}

	path := strings.TrimPrefix(uri.Path, "/")
	segments := strings.Split(path, "/")
	if len(segments) < 2 {
		return nil, fmt.Errorf("expected the path to contain at least 2 segments but got %d", len(segments))
	}
	shareName := segments[0]
	directoryPath := strings.Join(segments[1:], "/")
	return &DirectoryId{
		AccountId:     *account,
		ShareName:     shareName,
		DirectoryPath: directoryPath,
	}, nil
}
