# Frequently Asked Questions

> **Note:** This is a work-in-progress and will be extended over time.

### How can I help?

Great question, we assign labels to each GitHub issue to try and group them, a number of these are relevant for users looking to contribute:

* `good-first-issue` - this label is used to indicate that we think this would make a good issue for users looking to start contributing to the Provider. These are generally small enhancements, such as adding a new field to an existing resource - or documentation changes - and where we're adding this (in more recent issues) we're trying to provide a little context in one of the comments.
* `help-wanted` - we use this to highlight enhancement issues that are possible and will have a great impact, but that the maintainers are unlikely to reach in the near future.

The [Contributor Readme](../README.md) contains guides on the most common contribution types we see, but if you have any questions not answered in this documentation, please reach out (either in our community slack, or by opening an issue - details can be found in the contributor readme).

### How often is the Provider released?

The estimated dates for each release of the Provider [can be found on the Milestones page](https://github.com/hashicorp/terraform-provider-azurerm/milestones).

As a general rule the Provider is typically released weekly on a Thursday, however this can vary (for example during the winter holidays), as such we recommend checking [the Milestones page](https://github.com/hashicorp/terraform-provider-azurerm/milestones) for the most up to date information.

### My Pull Request has merge conflicts, should I rebase/merge from the `main` branch?

Whilst we do our best to review pull requests as they come in, unfortunately there are cases where it can take some time and merge conflicts can result if they have been sitting for a while. Generally speaking we recommend rebasing/merging from `main` only once a maintainer has taken a look through the PR and explicitly requested it.  

TODO.

### Once a major release is published, will new features and fixes be backported to previous versions?

Generally new features and fixes will only be added to the most recent major version.

Due to the high touch nature of provider development and the extensive regression testing required to ensure stability, maintaining multiple versions of the provider is not sustainable at this time. An exception to this could be a discovered security vulnerability for which backporting may be the most reasonable course of action. These will be reviewed on a case by case basis.

### What do the different GitHub labels mean?

As a general rule the different Azure Services are represented as `service/{serviceName}` - for other labels we're working through adding descriptions which [can be found on the GitHub Labels page for this repository](https://github.com/hashicorp/terraform-provider-azurerm/labels).

### Why was my comment marked as off-topic?

Whilst we thank you for your feedback, we mark comments along the lines of "me too" / "when will this be fixed?" (or generally off-topic comments) as off-topic so that they're hidden by default.

As this repository has a large/active community, we instead ask that you [use a thumbs-up GitHub reaction to the original issue](https://blog.github.com/2016-03-10-add-reactions-to-pull-requests-issues-and-comments/) so that we can prioritise this work without notifying everybody subscribed to the repository.

We appreciate this may be frustrating to have a comment marked as off-topic - when we've not done this we've noticed a number of users regularly adding "+1" / "me too" comments, which ends up causing more distractions for both the maintainers and community in general.

### Why did you close my question?

Whilst we thank you for reaching out, unfortunately we're unable to assist with individual usage questions related to the Azure Provider.

We've closed your issue because we believe it's an issue with the Terraform Configuration being used (or, that the credentials being used to interact with Azure may not have permission to the resources in question), rather than a bug in the Azure Provider.

We instead ask that configuration issues/usage questions related to the Provider are opened [on the Community Discuss forum](https://discuss.hashicorp.com) so that we can keep this repository focused on bugs/feature enhancements related to the Azure Provider.

