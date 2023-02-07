package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type WebhookId struct {
	SubscriptionId string
	ResourceGroup  string
	RegistryName   string
	Name           string
}

func NewWebhookID(subscriptionId, resourceGroup, registryName, name string) WebhookId {
	return WebhookId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		RegistryName:   registryName,
		Name:           name,
	}
}

func (id WebhookId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Registry Name %q", id.RegistryName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Webhook", segmentsStr)
}

func (id WebhookId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerRegistry/registries/%s/webHooks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.RegistryName, id.Name)
}

// WebhookID parses a Webhook ID into an WebhookId struct
func WebhookID(input string) (*WebhookId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := WebhookId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.RegistryName, err = id.PopSegment("registries"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("webHooks"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// WebhookIDInsensitively parses an Webhook ID into an WebhookId struct, insensitively
// This should only be used to parse an ID for rewriting, the WebhookID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func WebhookIDInsensitively(input string) (*WebhookId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := WebhookId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'registries' segment
	registriesKey := "registries"
	for key := range id.Path {
		if strings.EqualFold(key, registriesKey) {
			registriesKey = key
			break
		}
	}
	if resourceId.RegistryName, err = id.PopSegment(registriesKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'webHooks' segment
	webHooksKey := "webHooks"
	for key := range id.Path {
		if strings.EqualFold(key, webHooksKey) {
			webHooksKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(webHooksKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
