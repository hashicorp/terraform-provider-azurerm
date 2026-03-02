package files

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
)

// TODO: update this to implement `resourceids.ResourceId` once
// https://github.com/hashicorp/go-azure-helpers/issues/187 is fixed
var _ resourceids.Id = FileId{}

type FileId struct {
	// AccountId specifies the ID of the Storage Account where this File exists.
	AccountId accounts.AccountId

	// ShareName specifies the name of the File Share containing this File.
	ShareName string

	// DirectoryPath specifies the path representing the Directory where this File exists.
	DirectoryPath string

	// FileName specifies the name of the File.
	FileName string
}

func NewFileID(accountId accounts.AccountId, shareName, directoryPath, fileName string) FileId {
	return FileId{
		AccountId:     accountId,
		ShareName:     shareName,
		DirectoryPath: directoryPath,
		FileName:      fileName,
	}
}

func (b FileId) ID() string {
	path := ""
	if b.DirectoryPath != "" {
		path = fmt.Sprintf("%s/", b.DirectoryPath)
	}
	return fmt.Sprintf("%s/%s/%s%s", b.AccountId.ID(), b.ShareName, path, b.FileName)
}

func (b FileId) String() string {
	components := []string{
		fmt.Sprintf("Directory Path %q", b.DirectoryPath),
		fmt.Sprintf("Share Name %q", b.ShareName),
		fmt.Sprintf("Account %q", b.AccountId.String()),
	}
	return fmt.Sprintf("File %q (%s)", b.FileName, strings.Join(components, " / "))
}

// ParseFileID parses `input` into a File ID using a known `domainSuffix`
func ParseFileID(input, domainSuffix string) (*FileId, error) {
	// example: https://foo.file.core.windows.net/Bar/some/directory/some-file.txt
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
	directoryPath := strings.Join(segments[1:len(segments)-1], "/")
	fileName := segments[len(segments)-1]
	return &FileId{
		AccountId:     *account,
		ShareName:     shareName,
		DirectoryPath: directoryPath,
		FileName:      fileName,
	}, nil
}
