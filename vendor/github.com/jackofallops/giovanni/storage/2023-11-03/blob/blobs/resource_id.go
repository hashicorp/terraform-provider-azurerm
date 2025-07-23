package blobs

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
)

// TODO: update this to implement `resourceids.ResourceId` once
// https://github.com/hashicorp/go-azure-helpers/issues/187 is fixed
var _ resourceids.Id = BlobId{}

type BlobId struct {
	// AccountId specifies the ID of the Storage Account where this Blob exists.
	AccountId accounts.AccountId

	// ContainerName specifies the name of the Container within this Storage Account where this
	// Blob exists.
	ContainerName string

	// BlobName specifies the name of this Blob.
	BlobName string
}

func NewBlobID(accountId accounts.AccountId, containerName, blobName string) BlobId {
	return BlobId{
		AccountId:     accountId,
		ContainerName: containerName,
		BlobName:      blobName,
	}
}

func (b BlobId) ID() string {
	return fmt.Sprintf("%s/%s/%s", b.AccountId.ID(), b.ContainerName, b.BlobName)
}

func (b BlobId) String() string {
	components := []string{
		fmt.Sprintf("Account %q", b.AccountId.String()),
		fmt.Sprintf("Container Name %q", b.ContainerName),
	}
	return fmt.Sprintf("Blob %q (%s)", b.BlobName, strings.Join(components, " / "))
}

// ParseBlobID parses `input` into a Blob ID using a known `domainSuffix`
func ParseBlobID(input, domainSuffix string) (*BlobId, error) {
	// example: https://foo.blob.core.windows.net/Bar/example.vhd
	if input == "" {
		return nil, fmt.Errorf("`input` was empty")
	}

	account, err := accounts.ParseAccountID(input, domainSuffix)
	if err != nil {
		return nil, fmt.Errorf("parsing account %q: %+v", input, err)
	}

	if account.SubDomainType != accounts.BlobSubDomainType {
		return nil, fmt.Errorf("expected the subdomain type to be %q but got %q", string(accounts.BlobSubDomainType), string(account.SubDomainType))
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

	containerName := segments[0]
	blobName := strings.TrimPrefix(path, containerName)
	blobName = strings.TrimPrefix(blobName, "/")
	return &BlobId{
		AccountId:     *account,
		ContainerName: containerName,
		BlobName:      blobName,
	}, nil
}
