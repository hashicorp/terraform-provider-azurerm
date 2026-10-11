package tables

type MetaDataLevel string

var (
	NoMetaData      MetaDataLevel = "nometadata"
	MinimalMetaData MetaDataLevel = "minimalmetadata"
	FullMetaData    MetaDataLevel = "fullmetadata"
)

type GetResultItem struct {
	TableName string `json:"TableName"`

	// Optional, depending on the MetaData Level
	ODataType     string `json:"odata.type,omitempty"`
	ODataID       string `json:"odata.id,omitEmpty"`
	ODataEditLink string `json:"odata.editLink,omitEmpty"`
}

type SignedIdentifier struct {
	Id           string       `xml:"Id"`
	AccessPolicy AccessPolicy `xml:"AccessPolicy"`
}

type AccessPolicy struct {
	Start      string `xml:"Start"`
	Expiry     string `xml:"Expiry"`
	Permission string `xml:"Permission"`
}

type StorageServiceProperties struct {
	Logging       *LoggingConfig `xml:"Logging,omitempty"`
	HourMetrics   *MetricsConfig `xml:"HourMetrics,omitempty"`
	MinuteMetrics *MetricsConfig `xml:"MinuteMetrics,omitempty"`
	Cors          *Cors          `xml:"Cors,omitempty"`
}

type LoggingConfig struct {
	Version         string          `xml:"Version"`
	Delete          bool            `xml:"Delete"`
	Read            bool            `xml:"Read"`
	Write           bool            `xml:"Write"`
	RetentionPolicy RetentionPolicy `xml:"RetentionPolicy"`
}

type MetricsConfig struct {
	Version         string          `xml:"Version"`
	Enabled         bool            `xml:"Enabled"`
	RetentionPolicy RetentionPolicy `xml:"RetentionPolicy"`

	// Element IncludeAPIs is only expected when Metrics is enabled
	IncludeAPIs *bool `xml:"IncludeAPIs,omitempty"`
}

type RetentionPolicy struct {
	Enabled bool `xml:"Enabled"`
	Days    int  `xml:"Days,omitempty"`
}

type Cors struct {
	CorsRule []CorsRule `xml:"CorsRule"`
}

type CorsRule struct {
	AllowedOrigins  string `xml:"AllowedOrigins"`
	AllowedMethods  string `xml:"AllowedMethods"`
	AllowedHeaders  string `xml:"AllowedHeaders"`
	ExposedHeaders  string `xml:"ExposedHeaders"`
	MaxAgeInSeconds int    `xml:"MaxAgeInSeconds"`
}
