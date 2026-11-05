// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_insights_standard_web_test

type ApplicationInsightsStandardWebTestResourceModel struct {
	Name                  string                `tfschema:"name"`
	ResourceGroupName     string                `tfschema:"resource_group_name"`
	ApplicationInsightsID string                `tfschema:"application_insights_id"`
	Location              string                `tfschema:"location"`
	Frequency             int64                 `tfschema:"frequency"`
	Timeout               int64                 `tfschema:"timeout"`
	Enabled               bool                  `tfschema:"enabled"`
	Retry                 bool                  `tfschema:"retry_enabled"`
	Request               []RequestModel        `tfschema:"request"`
	ValidationRules       []ValidationRuleModel `tfschema:"validation_rules"`
	GeoLocations          []string              `tfschema:"geo_locations"`
	Description           string                `tfschema:"description"`
	Tags                  map[string]string     `tfschema:"tags"`

	// ComputedOnly
	SyntheticMonitorID string `tfschema:"synthetic_monitor_id"`
}

type RequestModel struct {
	FollowRedirects        bool          `tfschema:"follow_redirects_enabled"`
	HTTPVerb               string        `tfschema:"http_verb"`
	ParseDependentRequests bool          `tfschema:"parse_dependent_requests_enabled"`
	Header                 []HeaderModel `tfschema:"header"`
	Body                   string        `tfschema:"body"`
	URL                    string        `tfschema:"url"`
}

type ValidationRuleModel struct {
	ExpectedStatusCode           int64          `tfschema:"expected_status_code"`
	CertificateRemainingLifetime int64          `tfschema:"ssl_cert_remaining_lifetime"`
	SSLCheck                     bool           `tfschema:"ssl_check_enabled"`
	Content                      []ContentModel `tfschema:"content"`
}

type HeaderModel struct {
	Name  string `tfschema:"name"`
	Value string `tfschema:"value"`
}

type ContentModel struct {
	ContentMatch    string `tfschema:"content_match"`
	IgnoreCase      bool   `tfschema:"ignore_case"`
	PassIfTextFound bool   `tfschema:"pass_if_text_found"`
}

func (Resource) ModelObject() any {
	return &ApplicationInsightsStandardWebTestResourceModel{}
}
