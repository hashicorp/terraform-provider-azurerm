package parse

import (
	"fmt"
	"regexp"
)

type SubscriptionAliasId struct {
	Name string
}

func NewSubscriptionAliasId(name string) SubscriptionAliasId {
	return SubscriptionAliasId{
		Name: name,
	}
}

func (id SubscriptionAliasId) ID() string {
	return fmt.Sprintf("/providers/Microsoft.Subscription/aliases/%s", id.Name)
}

func SubscriptionAliasID(input string) (*SubscriptionAliasId, error) {
	groups := regexp.MustCompile(`^/providers/Microsoft.Subscription/aliases/([^/]+)$`).FindStringSubmatch(input)
	if len(groups) != 2 {
		return nil, fmt.Errorf("cannot parse resource id: %q", input)
	}

	return &SubscriptionAliasId{
		Name: groups[1],
	}, nil
}
