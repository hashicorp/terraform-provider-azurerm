package paths

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
)

// TODO: update this to implement `resourceids.ResourceId` once
// https://github.com/hashicorp/go-azure-helpers/issues/187 is fixed
var _ resourceids.Id = PathId{}

type PathId struct {
	// AccountId specifies the ID of the Storage Account where this path exists.
	AccountId accounts.AccountId

	// FileSystemName specifies the name of the Data Lake FileSystem where this Path exists.
	FileSystemName string

	// Path specifies the path in question.
	Path string
}

func NewPathID(accountId accounts.AccountId, fileSystemName, path string) PathId {
	return PathId{
		AccountId:      accountId,
		FileSystemName: fileSystemName,
		Path:           path,
	}
}

func (b PathId) ID() string {
	return fmt.Sprintf("%s/%s/%s", b.AccountId.ID(), b.FileSystemName, b.Path)
}

func (b PathId) String() string {
	components := []string{
		fmt.Sprintf("File System %q", b.FileSystemName),
		fmt.Sprintf("Account %q", b.AccountId.String()),
	}
	return fmt.Sprintf("Path %q (%s)", b.Path, strings.Join(components, " / "))
}

// ParsePathID parses `input` into a Path ID using a known `domainSuffix`
func ParsePathID(input, domainSuffix string) (*PathId, error) {
	// example: https://foo.dfs.core.windows.net/Bar/some/path
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

	uriPath := strings.TrimPrefix(uri.Path, "/")
	segments := strings.Split(uriPath, "/")
	if len(segments) < 2 {
		return nil, fmt.Errorf("expected the path to contain at least 2 segments but got %d", len(segments))
	}
	fileSystemName := segments[0]
	path := strings.Join(segments[1:], "/")
	return &PathId{
		AccountId:      *account,
		FileSystemName: fileSystemName,
		Path:           path,
	}, nil
}
