# KeyVault

> see https://aka.ms/autorest

This is the AutoRest configuration file for KeyVault.

---

## Getting Started

To build the SDK for KeyVault, simply [Install AutoRest](https://aka.ms/autorest/install) and in this folder, run:

> `autorest`

To see additional help and options, run:

> `autorest --help`

---

## Configuration

### Basic Information

These are the global settings for the KeyVault API.

``` yaml
openapi-type: data-plane
tag: package-7.4
```


### Tag: package-7.4

These settings apply only when `--tag=package-7.4` is specified on the command line.

```yaml $(tag) == 'package-7.4'
input-file:
  - Microsoft.KeyVault/stable/7.4/backuprestore.json
  - Microsoft.KeyVault/stable/7.4/certificates.json
  - Microsoft.KeyVault/stable/7.4/common.json
  - Microsoft.KeyVault/stable/7.4/keys.json
  - Microsoft.KeyVault/stable/7.4/rbac.json
  - Microsoft.KeyVault/stable/7.4/secrets.json
  - Microsoft.KeyVault/stable/7.4/securitydomain.json
  - Microsoft.KeyVault/stable/7.4/settings.json
  - Microsoft.KeyVault/stable/7.4/storage.json
```
### Tag: package-preview-7.4-preview.1

These settings apply only when `--tag=package-preview-7.4-preview.1` is specified on the command line.

``` yaml $(tag) == 'package-preview-7.4-preview.1'
input-file:
  - Microsoft.KeyVault/preview/7.4-preview.1/backuprestore.json
  - Microsoft.KeyVault/preview/7.4-preview.1/certificates.json
  - Microsoft.KeyVault/preview/7.4-preview.1/common.json
  - Microsoft.KeyVault/preview/7.4-preview.1/keys.json
  - Microsoft.KeyVault/preview/7.4-preview.1/rbac.json
  - Microsoft.KeyVault/preview/7.4-preview.1/secrets.json
  - Microsoft.KeyVault/preview/7.4-preview.1/securitydomain.json
  - Microsoft.KeyVault/preview/7.4-preview.1/settings.json
  - Microsoft.KeyVault/preview/7.4-preview.1/storage.json
```

### Tag: package-7.3

These settings apply only when `--tag=package-7.3` is specified on the command line.

``` yaml $(tag) == 'package-7.3'
input-file:
  - Microsoft.KeyVault/stable/7.3/backuprestore.json
  - Microsoft.KeyVault/stable/7.3/certificates.json
  - Microsoft.KeyVault/stable/7.3/common.json
  - Microsoft.KeyVault/stable/7.3/keys.json
  - Microsoft.KeyVault/stable/7.3/rbac.json
  - Microsoft.KeyVault/stable/7.3/secrets.json
  - Microsoft.KeyVault/stable/7.3/securitydomain.json
  - Microsoft.KeyVault/stable/7.3/storage.json
```

### Tag: package-preview-7.3-preview

These settings apply only when `--tag=package-preview-7.3-preview` is specified on the command line.

``` yaml $(tag) == 'package-preview-7.3-preview'
input-file:
  - Microsoft.KeyVault/preview/7.3-preview/backuprestore.json
  - Microsoft.KeyVault/preview/7.3-preview/certificates.json
  - Microsoft.KeyVault/preview/7.3-preview/common.json
  - Microsoft.KeyVault/preview/7.3-preview/keys.json
  - Microsoft.KeyVault/preview/7.3-preview/rbac.json
  - Microsoft.KeyVault/preview/7.3-preview/secrets.json
  - Microsoft.KeyVault/preview/7.3-preview/securitydomain.json
  - Microsoft.KeyVault/preview/7.3-preview/storage.json
```

### Tag: package-7.2

These settings apply only when `--tag=package-7.2` is specified on the command line.

``` yaml $(tag) == 'package-7.2'
input-file:
- Microsoft.KeyVault/stable/7.2/certificates.json
- Microsoft.KeyVault/stable/7.2/common.json
- Microsoft.KeyVault/stable/7.2/keys.json
- Microsoft.KeyVault/stable/7.2/rbac.json
- Microsoft.KeyVault/stable/7.2/secrets.json
- Microsoft.KeyVault/stable/7.2/storage.json
- Microsoft.KeyVault/stable/7.2/backuprestore.json
- Microsoft.KeyVault/stable/7.2/securitydomain.json
```

### Tag: package-7.2-preview

These settings apply only when `--tag=package-7.2-preview` is specified on the command line.

``` yaml $(tag) == 'package-7.2-preview'
input-file:
- Microsoft.KeyVault/preview/7.2-preview/certificates.json
- Microsoft.KeyVault/preview/7.2-preview/common.json
- Microsoft.KeyVault/preview/7.2-preview/keys.json
- Microsoft.KeyVault/preview/7.2-preview/rbac.json
- Microsoft.KeyVault/preview/7.2-preview/secrets.json
- Microsoft.KeyVault/preview/7.2-preview/storage.json
- Microsoft.KeyVault/preview/7.2-preview/backuprestore.json
- Microsoft.KeyVault/preview/7.2-preview/securitydomain.json
```

### Tag: package-7.1

These settings apply only when `--tag=package-7.1` is specified on the command line.

``` yaml $(tag) == 'package-7.1'
input-file:
- Microsoft.KeyVault/stable/7.1/certificates.json
- Microsoft.KeyVault/stable/7.1/common.json
- Microsoft.KeyVault/stable/7.1/keys.json
- Microsoft.KeyVault/stable/7.1/secrets.json
- Microsoft.KeyVault/stable/7.1/storage.json
```

### Tag: package-7.1-preview

These settings apply only when `--tag=package-7.1-preview` is specified on the command line.

``` yaml $(tag) == 'package-7.1-preview'
input-file:
- Microsoft.KeyVault/preview/7.1/certificates.json
- Microsoft.KeyVault/preview/7.1/common.json
- Microsoft.KeyVault/preview/7.1/keys.json
- Microsoft.KeyVault/preview/7.1/secrets.json
- Microsoft.KeyVault/preview/7.1/storage.json
```

### Tag: package-7.0

These settings apply only when `--tag=package-7.0` is specified on the command line.

``` yaml $(tag) == 'package-7.0'
input-file:
- Microsoft.KeyVault/stable/7.0/keyvault.json
```

### Tag: package-7.0-preview

These settings apply only when `--tag=package-7.0-preview` is specified on the command line.

``` yaml $(tag) == 'package-7.0-preview'
input-file:
- Microsoft.KeyVault/preview/7.0/keyvault.json
```

### Tag: package-2016-10

These settings apply only when `--tag=package-2016-10` is specified on the command line.

``` yaml $(tag) == 'package-2016-10'
input-file:
- Microsoft.KeyVault/stable/2016-10-01/keyvault.json
```

### Tag: package-2015-06

These settings apply only when `--tag=package-2015-06` is specified on the command line.

``` yaml $(tag) == 'package-2015-06'
input-file:
- Microsoft.KeyVault/stable/2015-06-01/keyvault.json
```

---

# Code Generation

## C#

These settings apply only when `--csharp` is specified on the command line.
Please also specify `--csharp-sdks-folder=<path to "SDKs" directory of your azure-sdk-for-net clone>`.

``` yaml $(csharp)
csharp:
  azure-arm: true
  license-header: MICROSOFT_MIT_NO_VERSION
  namespace: Microsoft.Azure.KeyVault
  sync-methods: None
  output-folder: $(csharp-sdks-folder)/keyvault/Microsoft.Azure.KeyVault/src/Generated
  clear-output-folder: true
```

## Go

See configuration in [readme.go.md](./readme.go.md)

## Java

These settings apply only when `--java` is specified on the command line.
Please also specify `--azure-libraries-for-java-folder=<path to the root directory of your azure-libraries-for-java clone>`.

``` yaml $(java)
java:
  azure-arm: true
  namespace: com.microsoft.azure.keyvault
  license-header: MICROSOFT_MIT_NO_CODEGEN
  payload-flattening-threshold: 0
  output-folder: $(azure-libraries-for-java-folder)/azure-keyvault
  override-client-name: KeyVaultClientBase
```

## Multi-API/Profile support for AutoRest v3 generators

AutoRest V3 generators require the use of `--tag=all-api-versions` to select api files.

This block is updated by an automatic script. Edits may be lost!

``` yaml $(tag) == 'all-api-versions' /* autogenerated */
# include the azure profile definitions from the standard location
require: $(this-folder)/../../../profiles/readme.md

# all the input files across all versions
input-file:
  - $(this-folder)/Microsoft.KeyVault/preview/7.3-preview/certificates.json
  - $(this-folder)/Microsoft.KeyVault/preview/7.3-preview/common.json
  - $(this-folder)/Microsoft.KeyVault/preview/7.3-preview/keys.json
  - $(this-folder)/Microsoft.KeyVault/preview/7.3-preview/rbac.json
  - $(this-folder)/Microsoft.KeyVault/preview/7.3-preview/secrets.json
  - $(this-folder)/Microsoft.KeyVault/preview/7.3-preview/storage.json
  - $(this-folder)/Microsoft.KeyVault/preview/7.3-preview/backuprestore.json
  - $(this-folder)/Microsoft.KeyVault/preview/7.3-preview/securitydomain.json
  - $(this-folder)/Microsoft.KeyVault/preview/7.2-preview/certificates.json
  - $(this-folder)/Microsoft.KeyVault/preview/7.2-preview/common.json
  - $(this-folder)/Microsoft.KeyVault/preview/7.2-preview/keys.json
  - $(this-folder)/Microsoft.KeyVault/preview/7.2-preview/rbac.json
  - $(this-folder)/Microsoft.KeyVault/preview/7.2-preview/secrets.json
  - $(this-folder)/Microsoft.KeyVault/preview/7.2-preview/storage.json
  - $(this-folder)/Microsoft.KeyVault/preview/7.2-preview/backuprestore.json
  - $(this-folder)/Microsoft.KeyVault/preview/7.2-preview/securitydomain.json
  - $(this-folder)/Microsoft.KeyVault/stable/7.1/certificates.json
  - $(this-folder)/Microsoft.KeyVault/stable/7.1/common.json
  - $(this-folder)/Microsoft.KeyVault/stable/7.1/keys.json
  - $(this-folder)/Microsoft.KeyVault/stable/7.1/secrets.json
  - $(this-folder)/Microsoft.KeyVault/stable/7.1/storage.json
  - $(this-folder)/Microsoft.KeyVault/preview/7.1/certificates.json
  - $(this-folder)/Microsoft.KeyVault/preview/7.1/common.json
  - $(this-folder)/Microsoft.KeyVault/preview/7.1/keys.json
  - $(this-folder)/Microsoft.KeyVault/preview/7.1/secrets.json
  - $(this-folder)/Microsoft.KeyVault/preview/7.1/storage.json
  - $(this-folder)/Microsoft.KeyVault/stable/7.0/keyvault.json
  - $(this-folder)/Microsoft.KeyVault/preview/7.0/keyvault.json
  - $(this-folder)/Microsoft.KeyVault/stable/2016-10-01/keyvault.json
  - $(this-folder)/Microsoft.KeyVault/stable/2015-06-01/keyvault.json

```

If there are files that should not be in the `all-api-versions` set,
uncomment the  `exclude-file` section below and add the file paths.

``` yaml $(tag) == 'all-api-versions'
#exclude-file: 
#  - $(this-folder)/Microsoft.Example/stable/2010-01-01/somefile.json
```

## Suppression

``` yaml
directive:
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: certificates.json
    where: $.definitions.CertificateOperation.properties.cancellation_requested
    reason: Merely refactored existing definitions into new files.
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: certificates.json
    where: $.definitions.CertificateOperation.properties.status_details
    reason: Merely refactored existing definitions into new files.
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: certificates.json
    where: $.definitions.CertificateOperation.properties.request_id
    reason: Merely refactored existing definitions into new files.
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: certificates.json
    where: $.definitions.CertificatePolicy.properties.key_props
    reason: Merely refactored existing definitions into new files.
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: certificates.json
    where: $.definitions.CertificatePolicy.properties.secret_props
    reason: Merely refactored existing definitions into new files.
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: certificates.json
    where: $.definitions.CertificatePolicy.properties.x509_props
    reason: Merely refactored existing definitions into new files.
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: certificates.json
    where: $.definitions.CertificatePolicy.properties.lifetime_actions
    reason: Merely refactored existing definitions into new files.
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: certificates.json
    where: $.definitions.X509CertificateProperties.properties.key_usage
    reason: Merely refactored existing definitions into new files.
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: certificates.json
    where: $.definitions.X509CertificateProperties.properties.validity_months
    reason: Merely refactored existing definitions into new files.
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: certificates.json
    where: $.definitions.IssuerParameters.properties.cert_transparency
    reason: Merely refactored existing definitions into new files.
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: certificates.json
    where: $.definitions.Action.properties.action_type
    reason: Merely refactored existing definitions into new files.
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: certificates.json
    where: $.definitions.Trigger.properties.lifetime_percentage
    reason: Merely refactored existing definitions into new files.
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: certificates.json
    where: $.definitions.Trigger.properties.days_before_expiry
    reason: Merely refactored existing definitions into new files.
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: certificates.json
    where: $.definitions.SubjectAlternativeNames.properties.dns_names
    reason: Merely refactored existing definitions into new files.
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: certificates.json
    where: $.definitions.IssuerBundle.properties.org_details
    reason: Merely refactored existing definitions into new files.
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: certificates.json
    where: $.definitions.IssuerCredentials.properties.account_id
    reason: Merely refactored existing definitions into new files.
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: certificates.json
    where: $.definitions.OrganizationDetails.properties.admin_details
    reason: Merely refactored existing definitions into new files.
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: certificates.json
    where: $.definitions.AdministratorDetails.properties.first_name
    reason: Merely refactored existing definitions into new files.
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: certificates.json
    where: $.definitions.AdministratorDetails.properties.last_name
    reason: Merely refactored existing definitions into new files.
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: certificates.json
    where: $.definitions.CertificateIssuerSetParameters.properties.org_details
    reason: Merely refactored existing definitions into new files.
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: certificates.json
    where: $.definitions.CertificateIssuerUpdateParameters.properties.org_details
    reason: Merely refactored existing definitions into new files.
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: certificates.json
    where: $.definitions.CertificateOperationUpdateParameter.properties.cancellation_requested
    reason: Merely refactored existing definitions into new files.
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: keys.json
    where: $.definitions.KeyProperties.properties.key_size
    reason: Merely refactored existing definitions into new files.
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: keys.json
    where: $.definitions.KeyProperties.properties.reuse_key
    reason: Merely refactored existing definitions into new files.
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: keys.json
    where: $.definitions.JsonWebKey.properties.key_ops
    reason: Merely refactored existing definitions into new files.
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: keys.json
    where: $.definitions.JsonWebKey.properties.key_hsm
    reason: Merely refactored existing definitions into new files.
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: keys.json
    where: $.definitions.KeyBundle.properties.release_policy
    reason: Consistency with other properties.
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: keys.json
    where: $.definitions.KeyCreateParameters.properties.key_size
    reason: Merely refactored existing definitions into new files.
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: keys.json
    where: $.definitions.KeyCreateParameters.properties.public_exponent
    reason: Consistency with other properties.
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: keys.json
    where: $.definitions.KeyCreateParameters.properties.release_policy
    reason: Consistency with other properties.
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: keys.json
    where: $.definitions.KeyCreateParameters.properties.key_ops
    reason: Merely refactored existing definitions into new files.
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: keys.json
    where: $.definitions.KeyImportParameters.properties.Hsm
    reason: Merely refactored existing definitions into new files.
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: keys.json
    where: $.definitions.KeyImportParameters.properties.release_policy
    reason: Consistency with other properties.
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: keys.json
    where: $.definitions.KeyUpdateParameters.properties.key_ops
    reason: Merely refactored existing definitions into new files.
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: keys.json
    where: $.definitions.KeyUpdateParameters.properties.release_policy
    reason: Consistency with other properties.
  - suppress: MISSING_REQUIRED_PARAMETER
    from: certificates.json
    where: $..parameters[?(@.name=='vaultBaseUrl')]
    reason: Suppress an invalid error caused by a bug in the linter.
  - suppress: MISSING_REQUIRED_PARAMETER
    from: keys.json
    where: $..parameters[?(@.name=='vaultBaseUrl')]
    reason: Suppress an invalid error caused by a bug in the linter.
  - suppress: MISSING_REQUIRED_PARAMETER
    from: secrets.json
    where: $..parameters[?(@.name=='vaultBaseUrl')]
    reason: Suppress an invalid error caused by a bug in the linter.
  - suppress: MISSING_REQUIRED_PARAMETER
    from: storage.json
    reason: Suppress an invalid error caused by a bug in the linter.
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: securitydomain.json
    where: $.definitions.TransferKey.properties.transfer_key
    reason: Merely refactored existing definitions into new files.
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: securitydomain.json
    where: $.definitions.UploadPendingResponse.properties.status_details
    reason: Consistency with other properties.
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: securitydomain.json
    where: $.definitions.SecurityDomainOperationStatus.properties.status_details
    reason: Consistency with other properties.
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: securitydomain.json
    where: $.definitions.SecurityDomainJsonWebKey.properties.key_ops
    reason: Consistency with other properties.
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: securitydomain.json
    where: $.definitions.SecurityDomainJsonWebKey.properties["x5t#S256"]
    reason: Consistency with other properties.
  - suppress: DefinitionsPropertiesNamesCamelCase
    from: securitydomain.json
    where: $.definitions.TransferKey.properties.key_format
    reason: Consistency with other properties
  - suppress: DOUBLE_FORWARD_SLASHES_IN_URL
    from: rbac.json
    reason: / is a valid scope in this scenario.
  - suppress: OBJECT_MISSING_REQUIRED_PROPERTY
    from: rbac.json
    where: $..parameters[?(@.name=='scope')]
    reason: Suppress an invalid error caused by a bug in the linter.
```
