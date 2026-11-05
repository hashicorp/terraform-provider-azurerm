// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package arc_kubernetes_flux_configuration

type ArcKubernetesFluxConfigurationModel struct {
	Name                            string                         `tfschema:"name"`
	ClusterID                       string                         `tfschema:"cluster_id"`
	BlobStorage                     []AzureBlobDefinitionModel     `tfschema:"blob_storage"`
	Bucket                          []BucketDefinitionModel        `tfschema:"bucket"`
	GitRepository                   []GitRepositoryDefinitionModel `tfschema:"git_repository"`
	Kustomizations                  []KustomizationDefinitionModel `tfschema:"kustomizations"`
	Namespace                       string                         `tfschema:"namespace"`
	Scope                           string                         `tfschema:"scope"`
	ContinuousReconciliationEnabled bool                           `tfschema:"continuous_reconciliation_enabled"`
}

type AzureBlobDefinitionModel struct {
	AccountKey            string                            `tfschema:"account_key"`
	ContainerID           string                            `tfschema:"container_id"`
	LocalAuthRef          string                            `tfschema:"local_auth_reference"`
	SasToken              string                            `tfschema:"sas_token"`
	ServicePrincipal      []ServicePrincipalDefinitionModel `tfschema:"service_principal"`
	SyncIntervalInSeconds int64                             `tfschema:"sync_interval_in_seconds"`
	TimeoutInSeconds      int64                             `tfschema:"timeout_in_seconds"`
}

type ServicePrincipalDefinitionModel struct {
	ClientCertificate          string `tfschema:"client_certificate_base64"`
	ClientCertificatePassword  string `tfschema:"client_certificate_password"`
	ClientCertificateSendChain bool   `tfschema:"client_certificate_send_chain"`
	ClientId                   string `tfschema:"client_id"`
	ClientSecret               string `tfschema:"client_secret"`
	TenantId                   string `tfschema:"tenant_id"`
}

type BucketDefinitionModel struct {
	AccessKey             string `tfschema:"access_key"`
	SecretKey             string `tfschema:"secret_key_base64"`
	BucketName            string `tfschema:"bucket_name"`
	TlsEnabled            bool   `tfschema:"tls_enabled"`
	LocalAuthRef          string `tfschema:"local_auth_reference"`
	SyncIntervalInSeconds int64  `tfschema:"sync_interval_in_seconds"`
	TimeoutInSeconds      int64  `tfschema:"timeout_in_seconds"`
	Url                   string `tfschema:"url"`
}

type GitRepositoryDefinitionModel struct {
	HttpsCACert           string `tfschema:"https_ca_cert_base64"`
	HttpsUser             string `tfschema:"https_user"`
	HttpsKey              string `tfschema:"https_key_base64"`
	LocalAuthRef          string `tfschema:"local_auth_reference"`
	ReferenceType         string `tfschema:"reference_type"`
	ReferenceValue        string `tfschema:"reference_value"`
	SshKnownHosts         string `tfschema:"ssh_known_hosts_base64"`
	SshPrivateKey         string `tfschema:"ssh_private_key_base64"`
	SyncIntervalInSeconds int64  `tfschema:"sync_interval_in_seconds"`
	TimeoutInSeconds      int64  `tfschema:"timeout_in_seconds"`
	Url                   string `tfschema:"url"`
}

type KustomizationDefinitionModel struct {
	Name                   string   `tfschema:"name"`
	Path                   string   `tfschema:"path"`
	TimeoutInSeconds       int64    `tfschema:"timeout_in_seconds"`
	SyncIntervalInSeconds  int64    `tfschema:"sync_interval_in_seconds"`
	RetryIntervalInSeconds int64    `tfschema:"retry_interval_in_seconds"`
	Force                  bool     `tfschema:"recreating_enabled"`
	Prune                  bool     `tfschema:"garbage_collection_enabled"`
	DependsOn              []string `tfschema:"depends_on"`
}

func (r Resource) ModelObject() any {
	return &ArcKubernetesFluxConfigurationModel{}
}
