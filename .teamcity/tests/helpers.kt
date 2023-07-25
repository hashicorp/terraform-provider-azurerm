/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

package tests

import ClientConfiguration

fun TestConfiguration() : ClientConfiguration {
    return ClientConfiguration("clientId", "clientSecret", "subscriptionId", "tenantId", "clientIdAlt", "clientSecretAlt", "subscriptionIdAlt", "subscriptionIdDevTest", "vcsRootId", true)
}