package queues

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
)

// GetResourceManagerResourceID returns the Resource Manager ID for the given Queue
// This can be useful when, for example, you're using this as a unique identifier
func (c Client) GetResourceManagerResourceID(subscriptionID, resourceGroup, accountName, queueName string) string {
	fmtStr := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/queueServices/default/queues/%s"
	return fmt.Sprintf(fmtStr, subscriptionID, resourceGroup, accountName, queueName)
}

// TODO: update this to implement `resourceids.ResourceId` once
// https://github.com/hashicorp/go-azure-helpers/issues/187 is fixed
var _ resourceids.Id = QueueId{}

type QueueId struct {
	// AccountId specifies the ID of the Storage Account where this Queue exists.
	AccountId accounts.AccountId

	// QueueName specifies the name of this Queue.
	QueueName string
}

func NewQueueID(accountId accounts.AccountId, queueName string) QueueId {
	return QueueId{
		AccountId: accountId,
		QueueName: queueName,
	}
}

func (b QueueId) ID() string {
	return fmt.Sprintf("%s/%s", b.AccountId.ID(), b.QueueName)
}

func (b QueueId) String() string {
	components := []string{
		fmt.Sprintf("Account %q", b.AccountId.String()),
	}
	return fmt.Sprintf("Queue %q (%s)", b.QueueName, strings.Join(components, " / "))
}

// ParseQueueID parses `input` into a Queue ID using a known `domainSuffix`
func ParseQueueID(input, domainSuffix string) (*QueueId, error) {
	// example: https://foo.queue.core.windows.net/Bar
	if input == "" {
		return nil, fmt.Errorf("`input` was empty")
	}

	account, err := accounts.ParseAccountID(input, domainSuffix)
	if err != nil {
		return nil, fmt.Errorf("parsing account %q: %+v", input, err)
	}

	if account.SubDomainType != accounts.QueueSubDomainType {
		return nil, fmt.Errorf("expected the subdomain type to be %q but got %q", string(accounts.QueueSubDomainType), string(account.SubDomainType))
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

	queueName := strings.TrimPrefix(uri.Path, "/")
	return &QueueId{
		AccountId: *account,
		QueueName: queueName,
	}, nil
}
