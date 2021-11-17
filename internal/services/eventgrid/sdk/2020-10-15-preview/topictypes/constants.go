package topictypes

import "strings"

type ResourceRegionType string

const (
	ResourceRegionTypeGlobalResource   ResourceRegionType = "GlobalResource"
	ResourceRegionTypeRegionalResource ResourceRegionType = "RegionalResource"
)

func PossibleValuesForResourceRegionType() []string {
	return []string{
		string(ResourceRegionTypeGlobalResource),
		string(ResourceRegionTypeRegionalResource),
	}
}

func parseResourceRegionType(input string) (*ResourceRegionType, error) {
	vals := map[string]ResourceRegionType{
		"globalresource":   ResourceRegionTypeGlobalResource,
		"regionalresource": ResourceRegionTypeRegionalResource,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceRegionType(input)
	return &out, nil
}

type SupportedScopesForSource string

const (
	SupportedScopesForSourceAzureSubscription SupportedScopesForSource = "AzureSubscription"
	SupportedScopesForSourceResource          SupportedScopesForSource = "Resource"
	SupportedScopesForSourceResourceGroup     SupportedScopesForSource = "ResourceGroup"
)

func PossibleValuesForSupportedScopesForSource() []string {
	return []string{
		string(SupportedScopesForSourceAzureSubscription),
		string(SupportedScopesForSourceResource),
		string(SupportedScopesForSourceResourceGroup),
	}
}

func parseSupportedScopesForSource(input string) (*SupportedScopesForSource, error) {
	vals := map[string]SupportedScopesForSource{
		"azuresubscription": SupportedScopesForSourceAzureSubscription,
		"resource":          SupportedScopesForSourceResource,
		"resourcegroup":     SupportedScopesForSourceResourceGroup,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SupportedScopesForSource(input)
	return &out, nil
}

type TopicTypeProvisioningState string

const (
	TopicTypeProvisioningStateCanceled  TopicTypeProvisioningState = "Canceled"
	TopicTypeProvisioningStateCreating  TopicTypeProvisioningState = "Creating"
	TopicTypeProvisioningStateDeleting  TopicTypeProvisioningState = "Deleting"
	TopicTypeProvisioningStateFailed    TopicTypeProvisioningState = "Failed"
	TopicTypeProvisioningStateSucceeded TopicTypeProvisioningState = "Succeeded"
	TopicTypeProvisioningStateUpdating  TopicTypeProvisioningState = "Updating"
)

func PossibleValuesForTopicTypeProvisioningState() []string {
	return []string{
		string(TopicTypeProvisioningStateCanceled),
		string(TopicTypeProvisioningStateCreating),
		string(TopicTypeProvisioningStateDeleting),
		string(TopicTypeProvisioningStateFailed),
		string(TopicTypeProvisioningStateSucceeded),
		string(TopicTypeProvisioningStateUpdating),
	}
}

func parseTopicTypeProvisioningState(input string) (*TopicTypeProvisioningState, error) {
	vals := map[string]TopicTypeProvisioningState{
		"canceled":  TopicTypeProvisioningStateCanceled,
		"creating":  TopicTypeProvisioningStateCreating,
		"deleting":  TopicTypeProvisioningStateDeleting,
		"failed":    TopicTypeProvisioningStateFailed,
		"succeeded": TopicTypeProvisioningStateSucceeded,
		"updating":  TopicTypeProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TopicTypeProvisioningState(input)
	return &out, nil
}
