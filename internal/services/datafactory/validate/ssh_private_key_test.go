// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import "testing"

func TestSSHPrivateKey(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			// empty
			Input: "",
			Valid: false,
		},

		{
			// Invalid test
			Input: "invalid private key",
			Valid: false,
		},

		{
			// valid RSA key
			Input: "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEAyc7a93dMN1KIPpCi4z9elnCr+fsAMq24BOi5vjJzRBFZBUFJ\naYXEg2NLxRVEp3E0BBl/9nophlEdN8x5/UBm5oLBQQYZWURhHy+c7EnDaGvx0nCf\nGOCrYIhTJQOwvQRuk2y99S2o/uEadpW5Cr5U036JVVZV4+3NwS0W2/nCjG8GyUq1\n95yR7+Gjjy1WBou7eIfcZm1bXW8Ti8MKq0DDlaOe0Ergg4bcPBuXhWGaIJzRBU1M\nWxw/3MA5JmYD+yfxDoR/NiGp4TjpnLByMDXmPgwTRrsOo0oxfEmuZkx/5qdCUHdH\n6GZEVMmsxQdMtW5W3uB1OfhPwXEBnecng+NOTQIDAQABAoIBAGcrOQYKFwyxRLW3\ne45xdwmx0Dmb+B3wcYMQ9uQlJohF1hy5o64ulKgWu0Wl+dMgLFdsMd7v1Qq1rRo5\njqPJqNFvRxzW4V6hdoVa8f5aN2vKw3Wx3aP6N6LCDr99g30eJul3TiVpklu6+Kxn\nHpI/e14j8lGOMZq0l9uKEYWjybHAqD7U+m6RCPo/92NTAeDL9/HKVuT2fLpjtDDr\nKxT8iOyqVEZVe7WJj/1dow+Sxr8XlEZNwmNwR9CTdD+ytxz1kQdRlLp8wD7rGrbJ\nn+Wy0p7AkguNcLW+WA4HVRZjY3Iu/kkMpb4z/s5Xcqik98R/zWRUQtDQDEmdvpkM\n4kkrNIECgYEA2iOKHqT/Djz/BIiOT7wzYbAbFLem1QANlGl1PjLPEjzD2ua2CWm1\n6isO+5s8Egl0RbIrDKJ7D2fE7KybY+53IOMNFk1J7sCd+Spb54scbOvIuhmEiZV7\nbkdpqEL63gUNsMK1LBH/QyqwHWTjYjxf+4PslOMdGvzSGtWpjDTssFcCgYEA7NWy\ndgrjBPUUHx97SWD6g9r9dmOb97AFcz1y/pftbgmBn1p4N5jrjgKE513gaObCMR4P\nWvmGskjDBmZsEh13AJsvzJOezQJxqoC5Gwg6DVvlksmzvIhvix6hKV1TJYGcnw//\nXgQsWG+dV6CCiZNIab/r/1phN/68qzFKmerbP/sCgYEAlsfCDPSXqEvZKlUJqWu5\nSGjmDyfylYB34oISnG+aWkzORFz8rvz21Wn17Uyb9Qu75wEEaLNWxItBvyaAMk7+\n4SiQPz4lQHa7uyLga+foOhGwqZJB3qgIrW2HRtsldJmhoPC1MkUuYEr9eRPnaFu7\nLKs/uJpT3/epcwsFKvjaMfkCgYEAnEbirMvAQ6woa+UNKD1q8QjXCYDvEQDAh+t7\nbw33aQ2yz+EVxLIOdTWqVzV1+CKU725Dead/vzMOJbH+C/IPbYH5h4e9WNANCxJH\nktPZ4qjKExMvm+93kxhSBgaD8BLXs4oN2w7r6Cs2avUKThe2x7kR0/zie92Gx3wO\nGwSDSnMCgYBy2P0Ua2Gwn9M1D/1F8E+f/jCmrdmxjtf2VhRu6X9cqbNq3LFasK/3\n7e78IhMgftOURJfImHu8bD8mKfP5O0adTccvdNYJFsX5HqByo6K5QbO+2Te+6LTf\nvEUTMZbJpVP8IhcQrr5FP1sREyeF8+T23G1CGAWgY3QLz3JcPNHBPQ==\n-----END RSA PRIVATE KEY-----",
			Valid: true,
		},

		{
			// valid DSA key
			Input: "-----BEGIN DSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEAyc7a93dMN1KIPpCi4z9elnCr+fsAMq24BOi5vjJzRBFZBUFJ\naYXEg2NLxRVEp3E0BBl/9nophlEdN8x5/UBm5oLBQQYZWURhHy+c7EnDaGvx0nCf\nGOCrYIhTJQOwvQRuk2y99S2o/uEadpW5Cr5U036JVVZV4+3NwS0W2/nCjG8GyUq1\n95yR7+Gjjy1WBou7eIfcZm1bXW8Ti8MKq0DDlaOe0Ergg4bcPBuXhWGaIJzRBU1M\nWxw/3MA5JmYD+yfxDoR/NiGp4TjpnLByMDXmPgwTRrsOo0oxfEmuZkx/5qdCUHdH\n6GZEVMmsxQdMtW5W3uB1OfhPwXEBnecng+NOTQIDAQABAoIBAGcrOQYKFwyxRLW3\ne45xdwmx0Dmb+B3wcYMQ9uQlJohF1hy5o64ulKgWu0Wl+dMgLFdsMd7v1Qq1rRo5\njqPJqNFvRxzW4V6hdoVa8f5aN2vKw3Wx3aP6N6LCDr99g30eJul3TiVpklu6+Kxn\nHpI/e14j8lGOMZq0l9uKEYWjybHAqD7U+m6RCPo/92NTAeDL9/HKVuT2fLpjtDDr\nKxT8iOyqVEZVe7WJj/1dow+Sxr8XlEZNwmNwR9CTdD+ytxz1kQdRlLp8wD7rGrbJ\nn+Wy0p7AkguNcLW+WA4HVRZjY3Iu/kkMpb4z/s5Xcqik98R/zWRUQtDQDEmdvpkM\n4kkrNIECgYEA2iOKHqT/Djz/BIiOT7wzYbAbFLem1QANlGl1PjLPEjzD2ua2CWm1\n6isO+5s8Egl0RbIrDKJ7D2fE7KybY+53IOMNFk1J7sCd+Spb54scbOvIuhmEiZV7\nbkdpqEL63gUNsMK1LBH/QyqwHWTjYjxf+4PslOMdGvzSGtWpjDTssFcCgYEA7NWy\ndgrjBPUUHx97SWD6g9r9dmOb97AFcz1y/pftbgmBn1p4N5jrjgKE513gaObCMR4P\nWvmGskjDBmZsEh13AJsvzJOezQJxqoC5Gwg6DVvlksmzvIhvix6hKV1TJYGcnw//\nXgQsWG+dV6CCiZNIab/r/1phN/68qzFKmerbP/sCgYEAlsfCDPSXqEvZKlUJqWu5\nSGjmDyfylYB34oISnG+aWkzORFz8rvz21Wn17Uyb9Qu75wEEaLNWxItBvyaAMk7+\n4SiQPz4lQHa7uyLga+foOhGwqZJB3qgIrW2HRtsldJmhoPC1MkUuYEr9eRPnaFu7\nLKs/uJpT3/epcwsFKvjaMfkCgYEAnEbirMvAQ6woa+UNKD1q8QjXCYDvEQDAh+t7\nbw33aQ2yz+EVxLIOdTWqVzV1+CKU725Dead/vzMOJbH+C/IPbYH5h4e9WNANCxJH\nktPZ4qjKExMvm+93kxhSBgaD8BLXs4oN2w7r6Cs2avUKThe2x7kR0/zie92Gx3wO\nGwSDSnMCgYBy2P0Ua2Gwn9M1D/1F8E+f/jCmrdmxjtf2VhRu6X9cqbNq3LFasK/3\n7e78IhMgftOURJfImHu8bD8mKfP5O0adTccvdNYJFsX5HqByo6K5QbO+2Te+6LTf\nvEUTMZbJpVP8IhcQrr5FP1sREyeF8+T23G1CGAWgY3QLz3JcPNHBPQ==\n-----END DSA PRIVATE KEY-----",
			Valid: true,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := SSHPrivateKey(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
