# Best Practices

A miscellaneous assortment of best practices to be aware of when contributing to the provider.

## Separate Create and Update Methods

Combined create and update methods within the provider are a legacy remnant that will need to be separated going forward and prior to the next major release of the provider.

Due to changes in behaviour in Terraform core and the providers migration from `terraform-pluginsdk` to `terraform-framework-plugin`

## Processing Azure API Responses and Setting values into State

## Untyped vs. Typed Resources

## Setting Properties to Optional + Computed

