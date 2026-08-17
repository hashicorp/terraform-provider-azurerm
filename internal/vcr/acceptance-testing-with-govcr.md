# Acceptance Testing with go-vcr

## Summary
`go-vcr` is used to record HTTP traffic from our acceptance tests and replay it for faster, (hopefully) more reliable testing suites. This guide is an overview of how the new plumbing works, how to run it, and where to look if you need to tweak things.

## How to Run Tests
The VCR framework is controlled via the `TC_TEST_VIA_VCR` environment variable. To prevent accidental mixing of real and mock data, we require explicit intention.

## Recording new tests: `TC_TEST_VIA_VCR=record`  
### What it does: 
The framework talks to actual Azure APIs using your real environment variables (like ARM_SUBSCRIPTION_ID). It saves exactly what happens to .yaml cassettes, safely scrubbing out your real IDs after the HTTP traffic has completed and right before the file writes to disk.

## Replaying tests (CI & Local debugging): `TC_TEST_VIA_VCR=replay` 
_(or `=true`, though this may be removed as it's too ambiguous?)_
### What it does: 
The acceptance test framework skips the network entirely. Most importantly, it globally overrides your environment variables (ARM_SUBSCRIPTION_ID, _ALT, _ALT2) to their 000000... placeholders right at startup. The provider natively builds its request configurations and Applied States out of these placeholders, providing validation against the redacted cassettes.

## Passthrough (Default fallback): 
If the variable is unset or unrecognised, the acceptance test framework falls back to `Passthrough`, ignoring VCR and talking straight to Azure without recording/replaying.

## What actually happens: Redaction & The Matcher
Security and consistency are the two hardest parts of mocking infrastructure. If a newly sensitive field  needs to be redacted, or if an API call starts suspiciously missing the cassette, the magic lives in `internal/vcr/recorder.go`.

1. Generic Regex & Alternate Subscriptions: To prevent sensitive data leakage, we assign them distinct, guaranteed placeholders (...0000, ...0001, ...0002 respectively). This is vital for tests that build multi-subscription setups (like Virtual Network Peering) to ensure the resources don't mistakenly merge onto the same placeholder. Need to redact a new cross-tenant ID? Add it to the idReplacements map in `GetRecorder()`.

2. Pre-Match Aggressive Scrubbing: One tricky mechanic of go-vcr is that its default matcher compares the entire request body and URI to find a match. During ImportStep or edge cases, if Terraform sends a request with real IDs, VCR would fail to find a match because the cassette is already redacted. To solve this cleanly, inside the MatcherFunc, we deep-copy the incoming request and redact the URI, Headers, and Body on the fly before handing it over to go-vcr to match.

3. The BeforeSaveHook: We wait until the test finishes completely before scrubbing the real requests, using go-vcr's `BeforeSaveHook`. It quietly intercepts the interaction list, thoroughly scrubs all URLs, Request bodies, and Response bodies, and writes the clean .yaml to disk. Because it happens offline at save-time, it doesn't break go-azure-sdk's long-running operation polling logic. 
_Note: the `AfterCaptureHook` looks tempting, but results in real API requests in downstream calls having the data redacted and ultimately failing._

4. Deterministic "Random" Data: VCR needs data predictability. To stop resource collisions and guarantee API matches, `vcrRandTimeInt()` in `data.go` simply takes the `t.Name()` string, dumps it into fnv.New64a(), and produces a "guaranteed"-unique 10-digit number. Combined with the fixed 20450101 prefix for consistency with "real" tests, it gives us reproducible 18-digit test data.

## Note for Maintainers
The intercept is wired into the `terraform-plugin-framework` provider implementation. Since this is ultimately bound together with the v2 provider by MUX, it's used for everything and we don't need specific code for PluginSDKv2.

