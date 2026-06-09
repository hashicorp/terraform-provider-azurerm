/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

package tests

import ClientConfiguration

fun TestConfiguration() : ClientConfiguration {
    return ClientConfiguration(
        "clientId",
        "clientSecret",
        "subscriptionId",
        "tenantId",
        "clientIdAlt",
        "clientSecretAlt",
        "subscriptionIdAlt",
        "subscriptionIdDevTest",
        "tenantIdAlt",
        "subscriptionIdAltTenant",
        "principalIdAltTenant",
        "vcsRootId",
        true,
        true,
        "terraformazuretestacc@example.com",
        "hashicorp/terraform-provider-azurerm",
        "gitPat",
        "teamcityToken",
        "env.ARM_FIVEPOINTZERO_BETA",
        "teamcity-passed",
        "teamcity-failed"
    )
}
