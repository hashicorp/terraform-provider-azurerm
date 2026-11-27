// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_gateway

type SpringCloudGatewayModel struct {
	Name                                  string                     `tfschema:"name"`
	SpringCloudServiceId                  string                     `tfschema:"spring_cloud_service_id"`
	ApiMetadata                           []ApiMetadataModel         `tfschema:"api_metadata"`
	ApplicationPerformanceMonitoringIds   []string                   `tfschema:"application_performance_monitoring_ids"`
	ApplicationPerformanceMonitoringTypes []string                   `tfschema:"application_performance_monitoring_types"`
	ClientAuthorization                   []ClientAuthorizationModel `tfschema:"client_authorization"`
	Cors                                  []CorsModel                `tfschema:"cors"`
	EnvironmentVariables                  map[string]string          `tfschema:"environment_variables"`
	LocalResponseCachePerRoute            []ResponseCacheModel       `tfschema:"local_response_cache_per_route"`
	LocalResponseCachePerInstance         []ResponseCacheModel       `tfschema:"local_response_cache_per_instance"`
	SensitiveEnvironmentVariables         map[string]string          `tfschema:"sensitive_environment_variables"`
	HttpsOnly                             bool                       `tfschema:"https_only"`
	InstanceCount                         int64                      `tfschema:"instance_count"`
	PublicNetworkAccessEnabled            bool                       `tfschema:"public_network_access_enabled"`
	Quota                                 []QuotaModel               `tfschema:"quota"`
	Sso                                   []GatewaySsoModel          `tfschema:"sso"`
	Url                                   string                     `tfschema:"url"`
}

type ApiMetadataModel struct {
	Description      string `tfschema:"description"`
	DocumentationUrl string `tfschema:"documentation_url"`
	ServerURL        string `tfschema:"server_url"`
	Title            string `tfschema:"title"`
	Version          string `tfschema:"version"`
}

type ClientAuthorizationModel struct {
	CertificateIds      []string `tfschema:"certificate_ids"`
	VerificationEnabled bool     `tfschema:"verification_enabled"`
}

type CorsModel struct {
	CredentialsAllowed    bool     `tfschema:"credentials_allowed"`
	AllowedHeaders        []string `tfschema:"allowed_headers"`
	AllowedMethods        []string `tfschema:"allowed_methods"`
	AllowedOrigins        []string `tfschema:"allowed_origins"`
	AllowedOriginPatterns []string `tfschema:"allowed_origin_patterns"`
	ExposedHeaders        []string `tfschema:"exposed_headers"`
	MaxAgeSeconds         int64    `tfschema:"max_age_seconds"`
}

type GatewaySsoModel struct {
	ClientId     string   `tfschema:"client_id"`
	ClientSecret string   `tfschema:"client_secret"`
	IssuerUri    string   `tfschema:"issuer_uri"`
	Scope        []string `tfschema:"scope"`
}

type QuotaModel struct {
	Cpu    string `tfschema:"cpu"`
	Memory string `tfschema:"memory"`
}

type ResponseCacheModel struct {
	Size       string `tfschema:"size"`
	TimeToLive string `tfschema:"time_to_live"`
}

func (s Resource) ModelObject() any {
	return &SpringCloudGatewayModel{}
}
