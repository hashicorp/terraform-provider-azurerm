package accounts

import "github.com/Azure/go-autorest/autorest"

type SetServicePropertiesResult struct {
	autorest.Response
}

type StorageServiceProperties struct {
	// Cors - Specifies CORS rules for the Blob service. You can include up to five CorsRule elements in the request. If no CorsRule elements are included in the request body, all CORS rules will be deleted, and CORS will be disabled for the Blob service.
	Cors *CorsRules `xml:"Cors,omitempty"`
	// DefaultServiceVersion - DefaultServiceVersion indicates the default version to use for requests to the Blob service if an incoming request’s version is not specified. Possible values include version 2008-10-27 and all more recent versions.
	DefaultServiceVersion *string `xml:"DefaultServiceVersion,omitempty"`
	// DeleteRetentionPolicy - The blob service properties for soft delete.
	DeleteRetentionPolicy *DeleteRetentionPolicy `xml:"DeleteRetentionPolicy,omitempty"`
	// Logging - The blob service properties for logging access
	Logging *Logging `xml:"Logging,omitempty"`
	// HourMetrics - The blob service properties for hour metrics
	HourMetrics *MetricsConfig `xml:"HourMetrics,omitempty"`
	// HourMetrics - The blob service properties for minute metrics
	MinuteMetrics *MetricsConfig `xml:"MinuteMetrics,omitempty"`
	// StaticWebsite - Optional
	StaticWebsite *StaticWebsite `xml:"StaticWebsite,omitempty"`
}

// StaticWebsite sets the static website support properties on the Blob service.
type StaticWebsite struct {
	// Enabled - Required. Indicates whether static website support is enabled for the given account.
	Enabled bool `xml:"Enabled"`
	// IndexDocument - Optional. The webpage that Azure Storage serves for requests to the root of a website or any subfolder. For example, index.html. The value is case-sensitive.
	IndexDocument string `xml:"IndexDocument,omitempty"`
	// ErrorDocument404Path - Optional. The absolute path to a webpage that Azure Storage serves for requests that do not correspond to an existing file. For example, error/404.html. Only a single custom 404 page is supported in each static website. The value is case-sensitive.
	ErrorDocument404Path string `xml:"ErrorDocument404Path,omitempty"`
}

// CorsRules sets the CORS rules. You can include up to five CorsRule elements in the request.
type CorsRules struct {
	// CorsRules - The List of CORS rules. You can include up to five CorsRule elements in the request.
	CorsRules []CorsRule `xml:"CorsRules,omitempty"`
}

// DeleteRetentionPolicy the blob service properties for soft delete.
type DeleteRetentionPolicy struct {
	// Enabled - Indicates whether DeleteRetentionPolicy is enabled for the Blob service.
	Enabled bool `xml:"Enabled,omitempty"`
	// Days - Indicates the number of days that the deleted blob should be retained. The minimum specified value can be 1 and the maximum value can be 365.
	Days int32 `xml:"Days,omitempty"`
}

// CorsRule specifies a CORS rule for the Blob service.
type CorsRule struct {
	// AllowedOrigins - Required if CorsRule element is present. A list of origin domains that will be allowed via CORS, or "" to allow all domains
	AllowedOrigins []string `xml:"AllowedOrigins,omitempty"`
	// AllowedMethods - Required if CorsRule element is present. A list of HTTP methods that are allowed to be executed by the origin.
	AllowedMethods []string `xml:"AllowedMethods,omitempty"`
	// MaxAgeInSeconds - Required if CorsRule element is present. The number of seconds that the client/browser should cache a preflight response.
	MaxAgeInSeconds int32 `xml:"MaxAgeInSeconds,omitempty"`
	// ExposedHeaders - Required if CorsRule element is present. A list of response headers to expose to CORS clients.
	ExposedHeaders []string `xml:"ExposedHeaders,omitempty"`
	// AllowedHeaders - Required if CorsRule element is present. A list of headers allowed to be part of the cross-origin request.
	AllowedHeaders []string `xml:"AllowedHeaders,omitempty"`
}

// Logging specifies the access logging options for the Blob service.
type Logging struct {
	Version         string                `xml:"Version"`
	Delete          bool                  `xml:"Delete"`
	Read            bool                  `xml:"Read"`
	Write           bool                  `xml:"Write"`
	RetentionPolicy DeleteRetentionPolicy `xml:"RetentionPolicy"`
}

// MetricsConfig specifies the hour and/or minute metrics options for the Blob service.
// Elements are all expected
type MetricsConfig struct {
	Version         string                `xml:"Version"`
	Enabled         bool                  `xml:"Enabled"`
	RetentionPolicy DeleteRetentionPolicy `xml:"RetentionPolicy"`
	IncludeAPIs     bool                  `xml:"IncludeAPIs"`
}
