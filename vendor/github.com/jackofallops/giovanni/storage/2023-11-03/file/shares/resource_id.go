package shares

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
)

// GetResourceManagerResourceID returns the Resource Manager specific
// ResourceID for a specific Storage Share
func (c Client) GetResourceManagerResourceID(subscriptionID, resourceGroup, accountName, shareName string) string {
	fmtStr := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/fileServices/default/shares/%s"
	return fmt.Sprintf(fmtStr, subscriptionID, resourceGroup, accountName, shareName)
}

// TODO: update this to implement `resourceids.ResourceId` once
// https://github.com/hashicorp/go-azure-helpers/issues/187 is fixed
var _ resourceids.Id = ShareId{}

type ShareId struct {
	// AccountId specifies the ID of the Storage Account where this File Share exists.
	AccountId accounts.AccountId

	// ShareName specifies the name of this File Share.
	ShareName string
}

func NewShareID(accountId accounts.AccountId, shareName string) ShareId {
	return ShareId{
		AccountId: accountId,
		ShareName: shareName,
	}
}

func (b ShareId) ID() string {
	return fmt.Sprintf("%s/%s", b.AccountId.ID(), b.ShareName)
}

func (b ShareId) String() string {
	components := []string{
		fmt.Sprintf("Account %q", b.AccountId.String()),
	}
	return fmt.Sprintf("File Share %q (%s)", b.ShareName, strings.Join(components, " / "))
}

// ParseShareID parses `input` into a Share ID using a known `domainSuffix`
func ParseShareID(input, domainSuffix string) (*ShareId, error) {
	// example: https://foo.file.core.windows.net/Bar
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
	if len(segments) == 0 {
		return nil, fmt.Errorf("expected the path to contain segments but got none")
	}

	shareName := strings.TrimPrefix(uri.Path, "/")
	return &ShareId{
		AccountId: *account,
		ShareName: shareName,
	}, nil
}
