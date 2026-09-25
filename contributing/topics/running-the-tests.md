# Running the Tests

> **Note:** Acceptance tests create real resources in Azure which often cost money to run.

Acceptance Tests for each Data Source/Resource are located within a Service Package, as such the Acceptance Tests for a given Service Package can be run via:

```sh
make acctests SERVICE='<service>' TESTTIMEOUT='60m'
```

However as many Service Packages contain multiple resources, you can opt to only run a subset by specifying the test prefix/filter to run as shown below:

```sh
make acctests SERVICE='<service>' TESTARGS='-run=<nameOfTheTest>' TESTTIMEOUT='60m'
```

* `<service>` is the name of the folder which contains the file with the test(s) you want to run. The available folders are found in `azurerm/internal/services/`. So examples are `mssql`, `compute` or `mariadb`
* `<nameOfTheTest>` should be self-explanatory as it is the name of the test you want to run. An example could be `TestAccMsSqlServerExtendedAuditingPolicy_basic`. Since `-run` can be used with regular expressions you can use it to specify multiple tests like in `TestAccMsSqlServerExtendedAuditingPolicy_` to run all tests that match that expression

The following Environment Variables must be set in your shell prior to running acceptance tests:

* `ARM_CLIENT_ID`
* `ARM_CLIENT_SECRET`
* `ARM_SUBSCRIPTION_ID`
* `ARM_TENANT_ID`
* `ARM_ENVIRONMENT`
* `ARM_METADATA_HOST`
* `ARM_TEST_LOCATION`
* `ARM_TEST_LOCATION_ALT`
* `ARM_TEST_LOCATION_ALT2`

> **Note:** Acceptance tests create real resources in Azure which often cost money to run.

## Private Kubernetes Fleet hub

`TestAccKubernetesFleetManager_privateHub` also requires `ARM_TEST_FLEET_SUBNET_ID`.
Set it to a dedicated, disposable subnet in the test subscription and
`ARM_TEST_LOCATION`. The Fleet resource provider's service principal must already
have `Network Contributor` on that subnet. The test skips when this variable is
absent; a skipped run does not verify private-hub support.
Confirm the current tenant, subscription, subnet, role assignment and permission
to use them before running this test; a previous successful run is not evidence
that its prerequisites still exist.

The test creates one private Fleet hub using `Standard_D2as_v7`, checks its
resource and data-source state, imports it without ignored fields, and updates
tags with both optional hub profiles omitted. Private hubs must omit `dns_prefix`.
The test harness destroys the Fleet and its test resource group. It does not
create or delete the supplied subnet or its role assignment; the prerequisite
owner must retain them until Fleet teardown completes and then clean them up.
No Kubernetes private-endpoint connectivity or workload is required by this
control-plane test.
