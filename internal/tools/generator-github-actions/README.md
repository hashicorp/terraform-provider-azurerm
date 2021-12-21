# Generator: Github Actions

## How do I use this tool?

This has to be run via `make generate` in the root of the Provider repository.

## What is this tool?

One of the challenges over an ever-growing codebase is an increase in compile-time when building, testing or linting
the entire Provider at a time. Whilst this is fine in some circumstances (for example, when building) - having the
unit tests and linting taking an extended period isn't ideal.

The AzureRM Provider makes extensive use of Service Packages to group code and other related functionality together.
Each Service Package has a number of common characteristics - namely they're each build, tested and linted in the
same manner.

As such rather than running the unit tests and linting the entire provider - we can instead do this on a per-Service
Package basis, conditionally running a Service Package only when something inside of it has changed.

Github Actions will automatically pick up and run new Github Actions when files added to the `./github` folder
whenever a pull request is opened from a branch on the main repository (or a pull request from a fork is merged) -
however pull requests from other repositories (forks) aren't picked up until the PR is merged.

This is an intentional security limitation on Github's part - however it means we end up needing to consider the
use-cases of both a new Service Packages and changes to an existing Service Package.

This tool generates a number of Github Actions, depending on where it's run:

Firstly, when this tool is run we generate a Github Action for each Service Package, in the form of
`./github/workflows/generated-service-{name}.yaml` - which ends up running `./scripts/service-package.sh {package}`
which contains the unit tests and linters required for each Service Package.

Secondly, when a PR containing changes to the generated Github Actions is merged, a different Github Action calls
this tool in another manner to regenerate the list of known Services and sub-sequentially open a PR.

Ultimately what this means is:

1. Each Service Package gets it's own Github Action, which runs the Unit Tests and Linters against it whenever
files in the Service Package have been changed.
2. New Service packages are run via a secondary script, which only runs Service Packages which haven't yet
been committed to the `main` branch.
3. Once a PR containing a new Service Package is merged into `main` - Github Actions will automatically open a
PR to secondary script to ignore the new package. This is so that PR's which introduce new Service Packages can
have tests/linters run against that code (as due to the intentional limitation new/changed Github Actions from
Pull Requests from Forks are not run).
