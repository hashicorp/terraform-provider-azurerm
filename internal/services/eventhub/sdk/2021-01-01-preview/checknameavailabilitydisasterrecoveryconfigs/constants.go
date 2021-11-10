package checknameavailabilitydisasterrecoveryconfigs

import "strings"

type UnavailableReason string

const (
	UnavailableReasonInvalidName                           UnavailableReason = "InvalidName"
	UnavailableReasonNameInLockdown                        UnavailableReason = "NameInLockdown"
	UnavailableReasonNameInUse                             UnavailableReason = "NameInUse"
	UnavailableReasonNone                                  UnavailableReason = "None"
	UnavailableReasonSubscriptionIsDisabled                UnavailableReason = "SubscriptionIsDisabled"
	UnavailableReasonTooManyNamespaceInCurrentSubscription UnavailableReason = "TooManyNamespaceInCurrentSubscription"
)

func PossibleValuesForUnavailableReason() []string {
	return []string{
		"InvalidName",
		"NameInLockdown",
		"NameInUse",
		"None",
		"SubscriptionIsDisabled",
		"TooManyNamespaceInCurrentSubscription",
	}
}

func parseUnavailableReason(input string) (*UnavailableReason, error) {
	vals := map[string]UnavailableReason{
		"invalidname":                           "InvalidName",
		"nameinlockdown":                        "NameInLockdown",
		"nameinuse":                             "NameInUse",
		"none":                                  "None",
		"subscriptionisdisabled":                "SubscriptionIsDisabled",
		"toomanynamespaceincurrentsubscription": "TooManyNamespaceInCurrentSubscription",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := UnavailableReason(v)
	return &out, nil
}
