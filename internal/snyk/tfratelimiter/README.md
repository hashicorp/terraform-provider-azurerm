# tfratelimiter

This is a small package that helps us enforce rate limits on cloud services.

The raison d'être for this package is that we want to use the same limiter
instance in our code as well in possibly patched terraform code; since cloud
API calls may happen from both codebases in short intervals.

Furthermore, we need to be able to scope this rate limiter to specific tenants
or subscriptions, since in some cases we use the same lambda function for
multiple of these and we don't want these to interfere.
