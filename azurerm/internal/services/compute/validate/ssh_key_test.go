package validate

import "testing"

func TestSSHKey(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			input:    "",
			expected: false,
		},
		{
			input:    "ssh-dss AAAAB3NzaC1kc3MAAACBAJTdkgVSk8cgM6h0MrnH9yoihsQVZ9c6OQcFqS1FZ/5DD4Z/8qfJlKFhICwhSCTX0dHqbZumG5KkFyrn2XznDf15idCHxxK4Vd51tyq5XaRyk89lFZCogIYPzocD+RdYVBwX7Y9ju+t7FqEhshd0q4tO6MzENIE//Wx+QWeiZrWlAAAAFQCsaVnyLr+Q+akj4M/K7pYR+GwpJQAAAIBtcypWCzJrPUgxy33rRMbrnWlQDY3H81iS4n7U5SDlUE7V0VaH8IxoQdSiGe6FJCUbu9XEvSQ+v6raBHPM6ca3t9NyPgBDdIRlCcgxrIQzbhTzgi85HdfDyED3wqDgMMdIYZ1AOeRQ3u3tLlGlOXrKCEIPH5x/tvysTn0+2mYKmwAAAIAtOGBS6M+IrrH+kMIOyLFGiL9b1s4rv5Vv6izULYb2DU0zoBnlRkmq/cLkFSgHeE5MqzOosybhwt5PRzMfoFtyUBpMgChdfuPnFwZbeTjitWRVS7tB/FDknbBXsk8mmnUEmodbTYVYtVSxbBgfKtc6pgomY1gxsYpByxyIA3A9gQ==",
			expected: false,
		},
		{
			input:    "ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBA95ywHY2HQsFe59iIhJCNmPjQdGbAJ7/5ZcxfOdHs98gG6UhCj5KwjpSICNGTZ+ZE+W4ExRPWzAGfFzjibUzsE=",
			expected: false,
		},
		{
			input:    "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOwlR9xtbM69hWLJbB5nHi0a65TuRvtaldgTJQ4ClL1W",
			expected: false,
		},
		{
			input:    "ssh-rsa",
			expected: false,
		},
		{
			input:    "ssh-rsa ThisIsNot a REAL key",
			expected: false,
		},
		{
			// 1024
			input:    "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDbzSM5KBFKmNilWjlw2YenzARxww1H+BMDMBVyzKYsNwEQc6Tj3ZB1Jun0l6Xkaw5BxKdwKFdhPlQh3nqpbm7xmSY7MuRZLPU+LRM3wI9RwcreDb3BXWacy41YIRGhUzpAzXmWdVyub/k70AJAngpVLLBmLcjuavjplR/fkTjslw==",
			expected: false,
		},
		{
			// 2048
			input:    "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0pA4vzGH+cmR+blZnoxO5HorOP1ubD4SxuOiW2DSNTSptlj+mPmFIL6sZeYMvSqAjXK368qL3DKHLpp2+1ws1XnYn/Zx/O4WBQAY7VbtzwFc7w7uirQaK6lVqXn8q4CnO0+5IYHgKLrNMEipwLKo+R3E3e1KrH5Xbyhj5yJzrMe3lWOAPzS27DJvjpN5SGWo65X6qFJRh3q95xOQhSOaEqZ/A2ZtfOuagq3FmASzoo/pbq7ianvnxzAYsb2Hg/9uAvypj4Beli6BP7419aP14XS0yyiW4XTKY/9XZiR/3VIKBN/stGN5NFLw82/j12E1GznbDG9PL7PQhijP7QgJh",
			expected: true,
		},
		{
			// 4096
			input:    "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDFP6r3wb/79MqRYI4dpgMwmjlrDDrk3A/pehysk1wzQn3lSEUtrNeQsHI6o/8au8Un1ndaZXZl/yWQQDDW4kqGw5ty8xPUZ+DB1ZVWkFOVNAgARl0bMNCgm2kB85l66g0zHWDCKLt+xi8xQiL7tGvdq3SWpogY3pWF2AABXoNDloHEN0mzzjJ09hdAHbygaDDr/9k3uyGKH3x0qo7fx5g8GqTtM3YWRxqUqdtkjsNomq94c/PMybCGR6qRoGI0Cdr/OP6/kszDHwf87B9hpTDMNa6x6FVJSDHc9v0CWePJZpjEOAFN3GCyPFFQTA9jvy026jt43wzyeH0kPe/T0ZZdr9YzQETN1b/oAKWKoayIoiLyJtFqUKcFFJSPcMz9ISgCD5Q/jRxQwMuMHpQ8TslxZ38l+41/0V1LWwKj0IkyJVFVWzu4zhgAZXr5y9Qbsis9sStRc+LU9/FQJ/VzNQfL83l86rH/u3NiPFfqisXILSybtMCD0OoRRHfQvWFsSwgt9JCIqLpmrJXRYs679aHzTHDgitlovJyprwqrbjg5N3XNSB5FohAUJUnVMF8z+qzvb4pPhly6mj6tiSJGYbXPngN6Iv8t3mRko3PbYLrWuxMb345BxcD+j9XteUgm1j/10qrSvqq+1R+/FAFPYwLXCflZgKst2g8/rEiVQz+a3w==",
			expected: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q...", v.input)

		_, errors := SSHKey(v.input, "public_key")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
