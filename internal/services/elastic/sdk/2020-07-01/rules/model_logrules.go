package rules

type LogRules struct {
	FilteringTags        *[]FilteringTag `json:"filteringTags,omitempty"`
	SendAadLogs          *bool           `json:"sendAadLogs,omitempty"`
	SendActivityLogs     *bool           `json:"sendActivityLogs,omitempty"`
	SendSubscriptionLogs *bool           `json:"sendSubscriptionLogs,omitempty"`
}
