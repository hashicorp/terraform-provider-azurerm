package queues

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
)

// GetResourceID returns the Resource ID for the given Queue
// This can be useful when, for example, you're using this as a unique identifier
func (client Client) GetResourceID(accountName, queueName string) string {
	domain := endpoints.GetQueueEndpoint(client.BaseURI, accountName)
	return fmt.Sprintf("%s/%s", domain, queueName)
}

type ResourceID struct {
	AccountName string
	QueueName   string
}

// ParseResourceID parses the Resource ID and returns an Object which
// can be used to interact with a Queue within a Storage Account
func ParseResourceID(id string) (*ResourceID, error) {
	// example: https://foo.queue.core.windows.net/Bar
	if id == "" {
		return nil, fmt.Errorf("`id` was empty")
	}

	uri, err := url.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("Error parsing ID as a URL: %s", err)
	}

	accountName, err := endpoints.GetAccountNameFromEndpoint(uri.Host)
	if err != nil {
		return nil, fmt.Errorf("Error parsing Account Name: %s", err)
	}

	queueName := strings.TrimPrefix(uri.Path, "/")
	return &ResourceID{
		AccountName: *accountName,
		QueueName:   queueName,
	}, nil
}
