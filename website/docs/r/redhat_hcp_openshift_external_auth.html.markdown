---
subcategory: "Red Hat OpenShift"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_redhat_hcp_openshift_external_auth"
description: |-
  Manages an Azure Red Hat OpenShift with Hosted Control Planes External Auth.
---

# azurerm_redhat_hcp_openshift_external_auth

Manages an external authentication provider for an Azure Red Hat OpenShift with Hosted Control Planes cluster.

## Example Usage

```hcl
resource "azurerm_redhat_hcp_openshift_external_auth" "example" {
  name                     = "example"
  hcp_openshift_cluster_id = azurerm_redhat_hcp_openshift_cluster.example.id

  issuer {
    url       = "https://issuer.example.com"
    audiences = ["example-client-id"]
  }

  claim {
    mappings {
      username {
        claim         = "sub"
        prefix_policy = "Prefix"
        prefix        = "oidc:"
      }

      groups {
        claim  = "groups"
        prefix = "oidc:"
      }
    }

    validation_rule {
      type = "RequiredClaim"

      required_claim {
        claim          = "tenant"
        required_value = "example-tenant"
      }
    }
  }

  client {
    client_id    = "example-client-id"
    type         = "Public"
    extra_scopes = ["email", "profile"]

    component {
      name                  = "console"
      auth_client_namespace = "openshift-console"
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the external authentication provider. Changing this forces a new resource to be created.
* `hcp_openshift_cluster_id` - (Required) The ID of the parent ARO HCP cluster. Changing this forces a new resource to be created.
* `issuer` - (Required) An `issuer` block as defined below.
* `claim` - (Required) A `claim` block as defined below.
* `client` - (Optional) One or more `client` blocks as defined below.

---

An `issuer` block supports the following:

* `url` - (Required) The HTTPS URL of the token issuer.
* `audiences` - (Required) A non-empty list of accepted token audiences.
* `ca` - (Optional) The certificate authority bundle used to validate the token issuer's TLS certificate.

---

A `claim` block supports the following:

* `mappings` - (Required) A `mappings` block as defined below.
* `validation_rule` - (Optional) One or more `validation_rule` blocks as defined below.

---

A `mappings` block supports the following:

* `username` - (Required) A `username` block as defined below.
* `groups` - (Optional) A `groups` block as defined below.

---

A `username` block supports the following:

* `claim` - (Required) The token claim used as the username.
* `prefix` - (Optional) The prefix applied to usernames. A non-empty value is required when `prefix_policy` is `Prefix`.
* `prefix_policy` - (Optional) The username prefix policy. Possible values are `NoPrefix`, `None`, and `Prefix`.

---

A `groups` block supports the following:

* `claim` - (Required) The token claim used for group membership.
* `prefix` - (Optional) The prefix applied to group names.

---

A `validation_rule` block supports the following:

* `type` - (Required) The validation rule type. The supported value is `RequiredClaim`.
* `required_claim` - (Optional) A `required_claim` block as defined below. Required when `type` is `RequiredClaim`.

---

A `required_claim` block supports the following:

* `claim` - (Required) The name of the required token claim.
* `required_value` - (Required) The value the claim must have.

---

A `client` block supports the following:

* `client_id` - (Required) The client ID registered with the identity provider.
* `type` - (Required) The client type. Possible values are `Confidential` and `Public`.
* `component` - (Required) A `component` block as defined below.
* `extra_scopes` - (Optional) A list of additional scopes to request from the identity provider.

---

A `component` block supports the following:

* `name` - (Required) The name of the cluster component using the client.
* `auth_client_namespace` - (Required) The namespace of the authentication client.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the external authentication provider.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the external authentication provider.
* `read` - (Defaults to 5 minutes) Used when retrieving the external authentication provider.
* `update` - (Defaults to 30 minutes) Used when updating the external authentication provider.
* `delete` - (Defaults to 30 minutes) Used when deleting the external authentication provider.

## Import

External authentication providers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_redhat_hcp_openshift_external_auth.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.RedHatOpenShift/hcpOpenShiftClusters/example-cluster/externalAuths/example
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.RedHatOpenShift` - 2026-10-01
