# Package: `github.com/hashicorp/go-azure-sdk/sdk/client/pollers`

This package contains both the components for Pollers, which allow an operation to be continually polled until it's either Completed, Cancelled or Failed.

Since Pollers are specific to the API in question, this package only contains the interface each poller needs to implement.

Specific implementations for each type of API can be found within the package for that API, for example a Poller for Long Running Operations within Azure Resource Manager can be found in `github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager`.
