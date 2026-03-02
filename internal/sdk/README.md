## SDK for Strongly-Typed Resources

This package is used for writing strongly-typed Data Sources and Resources.

The benefits of a typed resources over an untyped resource are:

* The Context object passed into each method _always_ has a deadline/timeout attached to it
* The Read function is automatically called at the end of a Create and Update function - meaning users don't have to do this 
* Each Resource has to have an ID Formatter and Validation Function
* The Model Object is validated via unit tests to ensure it contains the relevant struct tags

Ultimately this allows bugs to be caught by the Compiler (for example if a Read function is unimplemented) - or Unit Tests (for example should the `tfschema` struct tags be missing) - rather than during Provider Initialization, which reduces the feedback loop.

For help on how to add typed resources to the provider, please view the following tutorials in our contributor guidelines:

* [Add a new typed data source](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/contributing/topics/guide-new-data-source.md)
* [Add a new typed resource](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/contributing/topics/guide-new-resource.md)