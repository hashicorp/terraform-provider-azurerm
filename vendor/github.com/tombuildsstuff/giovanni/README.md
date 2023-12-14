# Giovanni

An alternative Azure Storage SDK for Go

---

This repository is an alternative Azure Storage SDK for Go; which supports for:

- The [Blob Storage APIs](https://docs.microsoft.com/en-us/rest/api/storageservices/blob-service-rest-api)
- The [DataLakeStorage Gen2 APIs](https://docs.microsoft.com/en-us/rest/api/storageservices/data-lake-storage-gen2)
- The [File Storage APIs](https://docs.microsoft.com/en-us/rest/api/storageservices/file-service-rest-api)
- The [Queue Storage APIs](https://docs.microsoft.com/en-us/rest/api/storageservices/queue-service-rest-api)
- The [Table Storage APIs](https://docs.microsoft.com/en-us/rest/api/storageservices/table-service-rest-api)

At this time we support the following API Versions:

* `2020-08-04` (`./storage/2020-08-04`)

We're also open to supporting other versions of the Azure Storage APIs as necessary.

Documentation for how to use each SDK can be found within the README for that SDK version - for example [here's the README for 2020-08-04](storage/2020-08-04/README.md).

Each Package also contains Unit and Acceptance tests to ensure that the functionality works; instructions on how to run the tests can be found below.

## Mission Statement

Fundamentally: developers should be able to pick which version of the Azure API they target using this SDK.

As such, there's two main goals here:

* New API Versions will be added additively to the `storage` folder.

* Any supported API Versions will continue to exist in the `storage` folder until they're EOL'd/stop working.

To ensure that each of these scenarios is possible - we have Acceptance and Unit Tests to confirm that the functionality in these versions works - and will use SemVer as appropriate.

## Future Enhancements

At this time this SDK is mostly feature complete, with a couple of notable additions (since we didn't need them).

Whilst it's possible to create Snapshots (for example, of a Container) - at this time most SDK calls don't support specifying the optional query-string value for `snapshot`.

In addition, we also don't support the `timeout` querystring on every API call; this is because instead all SDK methods take a `context` object, which allows a timeout to be set (albeit on the Client rather than the Remote API Call).

In both instances this is because we didn't need this functionality for our use-cases - but feel free to send a PR if you need this.

## Licence

Apache 2.0

## Technical Implementation

This SDK currently makes use of the standard Preparer-Sender-Responder pattern found in [Azure/go-autorest](https://github.com/Azure/go-autorest) - which means that this SDK should be familiar and compatible with [the Azure SDK for Go](https://github.com/Azure/azure-sdk-for-go).

> Note: A future release of this repository will switch to using `hashicorp/go-azure-sdk` as a base layer, rather than `Azure/go-autorest` - [see this issue for more information](https://github.com/tombuildsstuff/giovanni/issues/68).

Depending on the API Version / API being used - different authentication mechanisms are possible (see the README within the specific SDK for more info ([example](storage/2020-08-04/blob/accounts/README.md)). In all cases one of the following Authorizers will be required:

* An Authorizer for Azure Active Directory
* A SharedKeyLite Authorizer (for Blob, Queue and Table Storage)
* A SharedKeyLite Authorizer (for Table Storage)

Examples for all of these can be found below in [the Examples Directory](examples/).

A [SharedKey and SharedKeyLite Authorizer can be found in `Azure/go-autorest`](https://github.com/Azure/go-autorest/blob/ee71315119d4d7088d74ca9fcbc7301ce2ed2bc1/autorest/authorization_storage.go#L30-L48).

---

> Note: A future release of this repository will switch to using `hashicorp/go-azure-sdk` as a base layer, rather than `Azure/go-autorest` - [see this issue for more information](https://github.com/tombuildsstuff/giovanni/issues/68).

## Running the Tests

Each package contains both Unit and Acceptance Tests which provision a real Storage Account on Azure, and then run tests against that.

To run those, the following Environment Variables need to be set:

* `ARM_TENANT_ID` - The ID of the Tenant where tests should be run, such as `00000000-0000-0000-0000-000000000000`.
* `ARM_SUBSCRIPTION_ID` - The ID of the Subscription where tests should be run, such as `00000000-0000-0000-0000-000000000000`.
* `ARM_CLIENT_ID` - The ID of the AzureAD Application (also known as a Client ID), such as `00000000-0000-0000-0000-000000000000`.
* `ARM_CLIENT_SECRET` - The Client Secret/Password for a Service Principal where tests should be run.
* `ARM_ENVIRONMENT` - The name of the Azure Environment where the tests should be run, such as `Public` or `Germany`.
* `ARM_TEST_LOCATION` - The name of the Azure Region where resources provisioned by the tests should be created, such as `West Europe`.
* `ACCTEST` - confirms that you want the tests to be run, set this to any value.

Once those Environment Variables are set - you should be able to run:

```bash
$ ACCTEST=1 go test -v ./storage/...
```

You can also run them for a specific API version by running:

```bash
$ ACCTEST=1 go test -v ./storage/2020-08-04/...
```

## Debugging

You can see the Requests/Responses from this SDK by setting the Environment Variable `TEST_LOG` to any value.
