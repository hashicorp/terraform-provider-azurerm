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
