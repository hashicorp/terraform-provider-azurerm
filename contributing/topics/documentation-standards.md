# Provider Documentation Standards
<!-- TODO: Should this be a single page, or prefer to split this into multiple? (e.g. docs-notes.md) -->
In an effort to keep the [provider documentation](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs) consistent, this page documents some standards that have been agreed on. 

This page will grow over time, and suggestions are welcome!

## Notes

Note blocks are used to provide additional information to users beyond the basic description of a resource, argument or attribute.

<!-- 
    - TODO: Considerations for when to add notes? We probably don't want to overdo it (More relevant to informational notes)
    - TODO: Casing (Note vs NOTE)
-->

### Informational Note

Informational note blocks should generally be used when a note provides additional useful information, recommendations and/or tips to the user.

To add an informational note, use `-> **Note:**`, within the Terraform registry documentation this will template as a block with an info icon.

For example, extra information on the supported values for an argument, possibly linking to external documentation for the resource/service:

```markdown
* `type` - (Required) The type. Possible values include `This`, `That`, and `Other`.

-> **Note:** More information on each of the supported types can be found in [type documentation](link-to-additional-info)
```

### Warning Note

Warning note blocks should generally be used when a note provides information that the user will need to avoid certain errors, however if these errors are encountered they should not break anything or cause irreversible changes.

To add a warning note, use `~> **Note:**`, within the Terraform registry documentation this will template as a block with a warning icon.

For example, an argument that is optional but required when another argument is set to `true`:

```markdown
* `optional_argument_enabled` - (Optional) Is the optional argument enabled? Defaults to `false`.

* `optional_argument` - (Optional) An optional argument.

~> **Note:** The argument `optional_argument` is required when `optional_argument_enabled` is set to `true`.
```

### Caution Note

Caution note blocks should generally be used when a note provides critical information on potential irreversible changes, data loss or other things that can negatively affect a user's environment.

To add a caution note, use `!> **Note:**`, within the Terraform registry documentation this will template as a block with a caution icon.

For example, an argument that when set to `true` cannot be reversed without recreating the resource:

```markdown
* `irreversible_argument_enabled` - (Optional) Is irreversible argument enabled? Defaults to `false`.

!> **Note:** The argument `irreversible_argument_enabled` cannot be disabled after being enabled.
```