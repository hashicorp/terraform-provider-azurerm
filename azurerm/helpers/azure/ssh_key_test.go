package azure

import "testing"

func TestNormaliseSSHKey(t *testing.T) {
	cases := []struct {
		Input    string
		Error    bool
		Expected string
	}{
		{
			Input: "",
			Error: true,
		},
		{
			// Valid 2048 - no modification needed
			Input:    "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0pA4vzGH+cmR+blZnoxO5HorOP1ubD4SxuOiW2DSNTSptlj+mPmFIL6sZeYMvSqAjXK368qL3DKHLpp2+1ws1XnYn/Zx/O4WBQAY7VbtzwFc7w7uirQaK6lVqXn8q4CnO0+5IYHgKLrNMEipwLKo+R3E3e1KrH5Xbyhj5yJzrMe3lWOAPzS27DJvjpN5SGWo65X6qFJRh3q95xOQhSOaEqZ/A2ZtfOuagq3FmASzoo/pbq7ianvnxzAYsb2Hg/9uAvypj4Beli6BP7419aP14XS0yyiW4XTKY/9XZiR/3VIKBN/stGN5NFLw82/j12E1GznbDG9PL7PQhijP7QgJh generated-by-azure",
			Expected: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0pA4vzGH+cmR+blZnoxO5HorOP1ubD4SxuOiW2DSNTSptlj+mPmFIL6sZeYMvSqAjXK368qL3DKHLpp2+1ws1XnYn/Zx/O4WBQAY7VbtzwFc7w7uirQaK6lVqXn8q4CnO0+5IYHgKLrNMEipwLKo+R3E3e1KrH5Xbyhj5yJzrMe3lWOAPzS27DJvjpN5SGWo65X6qFJRh3q95xOQhSOaEqZ/A2ZtfOuagq3FmASzoo/pbq7ianvnxzAYsb2Hg/9uAvypj4Beli6BP7419aP14XS0yyiW4XTKY/9XZiR/3VIKBN/stGN5NFLw82/j12E1GznbDG9PL7PQhijP7QgJh generated-by-azure",
		},
		{
			// Valid 2048 - multiline, as per ARM Template Cache
			Input: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0pA4vzGH+cmR+blZnoxO5HorOP1ubD4SxuOiW2DSN\r\n" +
				"TSptlj+mPmFIL6sZeYMvSqAjXK368qL3DKHLpp2+1ws1XnYn/Zx/O4WBQAY7VbtzwFc7w7uirQaK6lVq\r\n" +
				"Xn8q4CnO0+5IYHgKLrNMEipwLKo+R3E3e1KrH5Xbyhj5yJzrMe3lWOAPzS27DJvjpN5SGWo65X6qFJRh\r\n" +
				"3q95xOQhSOaEqZ/A2ZtfOuagq3FmASzoo/pbq7ianvnxzAYsb2Hg/9uAvypj4Beli6BP7419aP14XS0y\r\n" +
				"yiW4XTKY/9XZiR/3VIKBN/stGN5NFLw82/j12E1GznbDG9PL7PQhijP7QgJh generated-by-azure",
			Expected: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0pA4vzGH+cmR+blZnoxO5HorOP1ubD4SxuOiW2DSNTSptlj+mPmFIL6sZeYMvSqAjXK368qL3DKHLpp2+1ws1XnYn/Zx/O4WBQAY7VbtzwFc7w7uirQaK6lVqXn8q4CnO0+5IYHgKLrNMEipwLKo+R3E3e1KrH5Xbyhj5yJzrMe3lWOAPzS27DJvjpN5SGWo65X6qFJRh3q95xOQhSOaEqZ/A2ZtfOuagq3FmASzoo/pbq7ianvnxzAYsb2Hg/9uAvypj4Beli6BP7419aP14XS0yyiW4XTKY/9XZiR/3VIKBN/stGN5NFLw82/j12E1GznbDG9PL7PQhijP7QgJh generated-by-azure",
		},
		{
			// Valid 2048 - multiline, as per ARM Template Cache Linux newlines
			Input: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0pA4vzGH+cmR+blZnoxO5HorOP1ubD4SxuOiW2DSN\n" +
				"TSptlj+mPmFIL6sZeYMvSqAjXK368qL3DKHLpp2+1ws1XnYn/Zx/O4WBQAY7VbtzwFc7w7uirQaK6lVq\n" +
				"Xn8q4CnO0+5IYHgKLrNMEipwLKo+R3E3e1KrH5Xbyhj5yJzrMe3lWOAPzS27DJvjpN5SGWo65X6qFJRh\n" +
				"3q95xOQhSOaEqZ/A2ZtfOuagq3FmASzoo/pbq7ianvnxzAYsb2Hg/9uAvypj4Beli6BP7419aP14XS0y\n" +
				"yiW4XTKY/9XZiR/3VIKBN/stGN5NFLw82/j12E1GznbDG9PL7PQhijP7QgJh generated-by-azure",
			Expected: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0pA4vzGH+cmR+blZnoxO5HorOP1ubD4SxuOiW2DSNTSptlj+mPmFIL6sZeYMvSqAjXK368qL3DKHLpp2+1ws1XnYn/Zx/O4WBQAY7VbtzwFc7w7uirQaK6lVqXn8q4CnO0+5IYHgKLrNMEipwLKo+R3E3e1KrH5Xbyhj5yJzrMe3lWOAPzS27DJvjpN5SGWo65X6qFJRh3q95xOQhSOaEqZ/A2ZtfOuagq3FmASzoo/pbq7ianvnxzAYsb2Hg/9uAvypj4Beli6BP7419aP14XS0yyiW4XTKY/9XZiR/3VIKBN/stGN5NFLw82/j12E1GznbDG9PL7PQhijP7QgJh generated-by-azure",
		},
		{
			// Valid 4096 - not modification required
			Input:    "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDFP6r3wb/79MqRYI4dpgMwmjlrDDrk3A/pehysk1wzQn3lSEUtrNeQsHI6o/8au8Un1ndaZXZl/yWQQDDW4kqGw5ty8xPUZ+DB1ZVWkFOVNAgARl0bMNCgm2kB85l66g0zHWDCKLt+xi8xQiL7tGvdq3SWpogY3pWF2AABXoNDloHEN0mzzjJ09hdAHbygaDDr/9k3uyGKH3x0qo7fx5g8GqTtM3YWRxqUqdtkjsNomq94c/PMybCGR6qRoGI0Cdr/OP6/kszDHwf87B9hpTDMNa6x6FVJSDHc9v0CWePJZpjEOAFN3GCyPFFQTA9jvy026jt43wzyeH0kPe/T0ZZdr9YzQETN1b/oAKWKoayIoiLyJtFqUKcFFJSPcMz9ISgCD5Q/jRxQwMuMHpQ8TslxZ38l+41/0V1LWwKj0IkyJVFVWzu4zhgAZXr5y9Qbsis9sStRc+LU9/FQJ/VzNQfL83l86rH/u3NiPFfqisXILSybtMCD0OoRRHfQvWFsSwgt9JCIqLpmrJXRYs679aHzTHDgitlovJyprwqrbjg5N3XNSB5FohAUJUnVMF8z+qzvb4pPhly6mj6tiSJGYbXPngN6Iv8t3mRko3PbYLrWuxMb345BxcD+j9XteUgm1j/10qrSvqq+1R+/FAFPYwLXCflZgKst2g8/rEiVQz+a3w== generated-by-azure",
			Expected: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDFP6r3wb/79MqRYI4dpgMwmjlrDDrk3A/pehysk1wzQn3lSEUtrNeQsHI6o/8au8Un1ndaZXZl/yWQQDDW4kqGw5ty8xPUZ+DB1ZVWkFOVNAgARl0bMNCgm2kB85l66g0zHWDCKLt+xi8xQiL7tGvdq3SWpogY3pWF2AABXoNDloHEN0mzzjJ09hdAHbygaDDr/9k3uyGKH3x0qo7fx5g8GqTtM3YWRxqUqdtkjsNomq94c/PMybCGR6qRoGI0Cdr/OP6/kszDHwf87B9hpTDMNa6x6FVJSDHc9v0CWePJZpjEOAFN3GCyPFFQTA9jvy026jt43wzyeH0kPe/T0ZZdr9YzQETN1b/oAKWKoayIoiLyJtFqUKcFFJSPcMz9ISgCD5Q/jRxQwMuMHpQ8TslxZ38l+41/0V1LWwKj0IkyJVFVWzu4zhgAZXr5y9Qbsis9sStRc+LU9/FQJ/VzNQfL83l86rH/u3NiPFfqisXILSybtMCD0OoRRHfQvWFsSwgt9JCIqLpmrJXRYs679aHzTHDgitlovJyprwqrbjg5N3XNSB5FohAUJUnVMF8z+qzvb4pPhly6mj6tiSJGYbXPngN6Iv8t3mRko3PbYLrWuxMb345BxcD+j9XteUgm1j/10qrSvqq+1R+/FAFPYwLXCflZgKst2g8/rEiVQz+a3w== generated-by-azure",
		},
		{
			// Valid 4096 - multiline Windows
			Input: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDFP6r3wb/79MqRYI4dpgMwmjlrDDrk3A/pehysk1wz\r\n" +
				"Qn3lSEUtrNeQsHI6o/8au8Un1ndaZXZl/yWQQDDW4kqGw5ty8xPUZ+DB1ZVWkFOVNAgARl0bMNCgm2kB\r\n" +
				"85l66g0zHWDCKLt+xi8xQiL7tGvdq3SWpogY3pWF2AABXoNDloHEN0mzzjJ09hdAHbygaDDr/9k3uyGK\r\n" +
				"H3x0qo7fx5g8GqTtM3YWRxqUqdtkjsNomq94c/PMybCGR6qRoGI0Cdr/OP6/kszDHwf87B9hpTDMNa6x\r\n" +
				"6FVJSDHc9v0CWePJZpjEOAFN3GCyPFFQTA9jvy026jt43wzyeH0kPe/T0ZZdr9YzQETN1b/oAKWKoayI\r\n" +
				"oiLyJtFqUKcFFJSPcMz9ISgCD5Q/jRxQwMuMHpQ8TslxZ38l+41/0V1LWwKj0IkyJVFVWzu4zhgAZXr5\r\n" +
				"y9Qbsis9sStRc+LU9/FQJ/VzNQfL83l86rH/u3NiPFfqisXILSybtMCD0OoRRHfQvWFsSwgt9JCIqLpm\r\n" +
				"rJXRYs679aHzTHDgitlovJyprwqrbjg5N3XNSB5FohAUJUnVMF8z+qzvb4pPhly6mj6tiSJGYbXPngN6\r\n" +
				"Iv8t3mRko3PbYLrWuxMb345BxcD+j9XteUgm1j/10qrSvqq+1R+/FAFPYwLXCflZgKst2g8/rEiVQz+a\r\n" +
				"3w== generated-by-azure",
			Expected: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDFP6r3wb/79MqRYI4dpgMwmjlrDDrk3A/pehysk1wzQn3lSEUtrNeQsHI6o/8au8Un1ndaZXZl/yWQQDDW4kqGw5ty8xPUZ+DB1ZVWkFOVNAgARl0bMNCgm2kB85l66g0zHWDCKLt+xi8xQiL7tGvdq3SWpogY3pWF2AABXoNDloHEN0mzzjJ09hdAHbygaDDr/9k3uyGKH3x0qo7fx5g8GqTtM3YWRxqUqdtkjsNomq94c/PMybCGR6qRoGI0Cdr/OP6/kszDHwf87B9hpTDMNa6x6FVJSDHc9v0CWePJZpjEOAFN3GCyPFFQTA9jvy026jt43wzyeH0kPe/T0ZZdr9YzQETN1b/oAKWKoayIoiLyJtFqUKcFFJSPcMz9ISgCD5Q/jRxQwMuMHpQ8TslxZ38l+41/0V1LWwKj0IkyJVFVWzu4zhgAZXr5y9Qbsis9sStRc+LU9/FQJ/VzNQfL83l86rH/u3NiPFfqisXILSybtMCD0OoRRHfQvWFsSwgt9JCIqLpmrJXRYs679aHzTHDgitlovJyprwqrbjg5N3XNSB5FohAUJUnVMF8z+qzvb4pPhly6mj6tiSJGYbXPngN6Iv8t3mRko3PbYLrWuxMb345BxcD+j9XteUgm1j/10qrSvqq+1R+/FAFPYwLXCflZgKst2g8/rEiVQz+a3w== generated-by-azure",
		},
		{
			// Valid 4096 - multiline Linux newlines
			Input: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDFP6r3wb/79MqRYI4dpgMwmjlrDDrk3A/pehysk1wz\n" +
				"Qn3lSEUtrNeQsHI6o/8au8Un1ndaZXZl/yWQQDDW4kqGw5ty8xPUZ+DB1ZVWkFOVNAgARl0bMNCgm2kB\n" +
				"85l66g0zHWDCKLt+xi8xQiL7tGvdq3SWpogY3pWF2AABXoNDloHEN0mzzjJ09hdAHbygaDDr/9k3uyGK\n" +
				"H3x0qo7fx5g8GqTtM3YWRxqUqdtkjsNomq94c/PMybCGR6qRoGI0Cdr/OP6/kszDHwf87B9hpTDMNa6x\n" +
				"6FVJSDHc9v0CWePJZpjEOAFN3GCyPFFQTA9jvy026jt43wzyeH0kPe/T0ZZdr9YzQETN1b/oAKWKoayI\n" +
				"oiLyJtFqUKcFFJSPcMz9ISgCD5Q/jRxQwMuMHpQ8TslxZ38l+41/0V1LWwKj0IkyJVFVWzu4zhgAZXr5\n" +
				"y9Qbsis9sStRc+LU9/FQJ/VzNQfL83l86rH/u3NiPFfqisXILSybtMCD0OoRRHfQvWFsSwgt9JCIqLpm\n" +
				"rJXRYs679aHzTHDgitlovJyprwqrbjg5N3XNSB5FohAUJUnVMF8z+qzvb4pPhly6mj6tiSJGYbXPngN6\n" +
				"Iv8t3mRko3PbYLrWuxMb345BxcD+j9XteUgm1j/10qrSvqq+1R+/FAFPYwLXCflZgKst2g8/rEiVQz+a\n" +
				"3w== generated-by-azure",
			Expected: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDFP6r3wb/79MqRYI4dpgMwmjlrDDrk3A/pehysk1wzQn3lSEUtrNeQsHI6o/8au8Un1ndaZXZl/yWQQDDW4kqGw5ty8xPUZ+DB1ZVWkFOVNAgARl0bMNCgm2kB85l66g0zHWDCKLt+xi8xQiL7tGvdq3SWpogY3pWF2AABXoNDloHEN0mzzjJ09hdAHbygaDDr/9k3uyGKH3x0qo7fx5g8GqTtM3YWRxqUqdtkjsNomq94c/PMybCGR6qRoGI0Cdr/OP6/kszDHwf87B9hpTDMNa6x6FVJSDHc9v0CWePJZpjEOAFN3GCyPFFQTA9jvy026jt43wzyeH0kPe/T0ZZdr9YzQETN1b/oAKWKoayIoiLyJtFqUKcFFJSPcMz9ISgCD5Q/jRxQwMuMHpQ8TslxZ38l+41/0V1LWwKj0IkyJVFVWzu4zhgAZXr5y9Qbsis9sStRc+LU9/FQJ/VzNQfL83l86rH/u3NiPFfqisXILSybtMCD0OoRRHfQvWFsSwgt9JCIqLpmrJXRYs679aHzTHDgitlovJyprwqrbjg5N3XNSB5FohAUJUnVMF8z+qzvb4pPhly6mj6tiSJGYbXPngN6Iv8t3mRko3PbYLrWuxMb345BxcD+j9XteUgm1j/10qrSvqq+1R+/FAFPYwLXCflZgKst2g8/rEiVQz+a3w== generated-by-azure",
		},
		{
			// Valid 4096 - multiline Windows Wrapped
			Input: "<<~EOT\r\n" +
				"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDFP6r3wb/79MqRYI4dpgMwmjlrDDrk3A/pehysk1wz\r\n" +
				"Qn3lSEUtrNeQsHI6o/8au8Un1ndaZXZl/yWQQDDW4kqGw5ty8xPUZ+DB1ZVWkFOVNAgARl0bMNCgm2kB\r\n" +
				"85l66g0zHWDCKLt+xi8xQiL7tGvdq3SWpogY3pWF2AABXoNDloHEN0mzzjJ09hdAHbygaDDr/9k3uyGK\r\n" +
				"H3x0qo7fx5g8GqTtM3YWRxqUqdtkjsNomq94c/PMybCGR6qRoGI0Cdr/OP6/kszDHwf87B9hpTDMNa6x\r\n" +
				"6FVJSDHc9v0CWePJZpjEOAFN3GCyPFFQTA9jvy026jt43wzyeH0kPe/T0ZZdr9YzQETN1b/oAKWKoayI\r\n" +
				"oiLyJtFqUKcFFJSPcMz9ISgCD5Q/jRxQwMuMHpQ8TslxZ38l+41/0V1LWwKj0IkyJVFVWzu4zhgAZXr5\r\n" +
				"y9Qbsis9sStRc+LU9/FQJ/VzNQfL83l86rH/u3NiPFfqisXILSybtMCD0OoRRHfQvWFsSwgt9JCIqLpm\r\n" +
				"rJXRYs679aHzTHDgitlovJyprwqrbjg5N3XNSB5FohAUJUnVMF8z+qzvb4pPhly6mj6tiSJGYbXPngN6\r\n" +
				"Iv8t3mRko3PbYLrWuxMb345BxcD+j9XteUgm1j/10qrSvqq+1R+/FAFPYwLXCflZgKst2g8/rEiVQz+a\r\n" +
				"3w== generated-by-azure" +
				"EOT",
			Expected: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDFP6r3wb/79MqRYI4dpgMwmjlrDDrk3A/pehysk1wzQn3lSEUtrNeQsHI6o/8au8Un1ndaZXZl/yWQQDDW4kqGw5ty8xPUZ+DB1ZVWkFOVNAgARl0bMNCgm2kB85l66g0zHWDCKLt+xi8xQiL7tGvdq3SWpogY3pWF2AABXoNDloHEN0mzzjJ09hdAHbygaDDr/9k3uyGKH3x0qo7fx5g8GqTtM3YWRxqUqdtkjsNomq94c/PMybCGR6qRoGI0Cdr/OP6/kszDHwf87B9hpTDMNa6x6FVJSDHc9v0CWePJZpjEOAFN3GCyPFFQTA9jvy026jt43wzyeH0kPe/T0ZZdr9YzQETN1b/oAKWKoayIoiLyJtFqUKcFFJSPcMz9ISgCD5Q/jRxQwMuMHpQ8TslxZ38l+41/0V1LWwKj0IkyJVFVWzu4zhgAZXr5y9Qbsis9sStRc+LU9/FQJ/VzNQfL83l86rH/u3NiPFfqisXILSybtMCD0OoRRHfQvWFsSwgt9JCIqLpmrJXRYs679aHzTHDgitlovJyprwqrbjg5N3XNSB5FohAUJUnVMF8z+qzvb4pPhly6mj6tiSJGYbXPngN6Iv8t3mRko3PbYLrWuxMb345BxcD+j9XteUgm1j/10qrSvqq+1R+/FAFPYwLXCflZgKst2g8/rEiVQz+a3w== generated-by-azure",
		},
		{
			// Valid 4096 - multiline Windows Wrapped with whitespace
			Input: "<<~EOT\r\n" +
				"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDFP6r3wb/79MqRYI4dpgMwmjlrDDrk3A/pehysk1wz\r\n  " +
				"Qn3lSEUtrNeQsHI6o/8au8Un1ndaZXZl/yWQQDDW4kqGw5ty8xPUZ+DB1ZVWkFOVNAgARl0bMNCgm2kB\r\n  " +
				"85l66g0zHWDCKLt+xi8xQiL7tGvdq3SWpogY3pWF2AABXoNDloHEN0mzzjJ09hdAHbygaDDr/9k3uyGK\r\n" +
				"  H3x0qo7fx5g8GqTtM3YWRxqUqdtkjsNomq94c/PMybCGR6qRoGI0Cdr/OP6/kszDHwf87B9hpTDMNa6x\r\n" +
				"6FVJSDHc9v0CWePJZpjEOAFN3GCyPFFQTA9jvy026jt43wzyeH0kPe/T0ZZdr9YzQETN1b/oAKWKoayI\r\n" +
				"oiLyJtFqUKcFFJSPcMz9ISgCD5Q/jRxQwMuMHpQ8TslxZ38l+41/0V1LWwKj0IkyJVFVWzu4zhgAZXr5\r\n" +
				"y9Qbsis9sStRc+LU9/FQJ/VzNQfL83l86rH/u3NiPFfqisXILSybtMCD0OoRRHfQvWFsSwgt9JCIqLpm\r\n" +
				" rJXRYs679aHzTHDgitlovJyprwqrbjg5N3XNSB5FohAUJUnVMF8z+qzvb4pPhly6mj6tiSJGYbXPngN6\r\n" +
				"    Iv8t3mRko3PbYLrWuxMb345BxcD+j9XteUgm1j/10qrSvqq+1R+/FAFPYwLXCflZgKst2g8/rEiVQz+a\r\n    " +
				"3w== generated-by-azure" +
				"EOT",
			Expected: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDFP6r3wb/79MqRYI4dpgMwmjlrDDrk3A/pehysk1wzQn3lSEUtrNeQsHI6o/8au8Un1ndaZXZl/yWQQDDW4kqGw5ty8xPUZ+DB1ZVWkFOVNAgARl0bMNCgm2kB85l66g0zHWDCKLt+xi8xQiL7tGvdq3SWpogY3pWF2AABXoNDloHEN0mzzjJ09hdAHbygaDDr/9k3uyGKH3x0qo7fx5g8GqTtM3YWRxqUqdtkjsNomq94c/PMybCGR6qRoGI0Cdr/OP6/kszDHwf87B9hpTDMNa6x6FVJSDHc9v0CWePJZpjEOAFN3GCyPFFQTA9jvy026jt43wzyeH0kPe/T0ZZdr9YzQETN1b/oAKWKoayIoiLyJtFqUKcFFJSPcMz9ISgCD5Q/jRxQwMuMHpQ8TslxZ38l+41/0V1LWwKj0IkyJVFVWzu4zhgAZXr5y9Qbsis9sStRc+LU9/FQJ/VzNQfL83l86rH/u3NiPFfqisXILSybtMCD0OoRRHfQvWFsSwgt9JCIqLpmrJXRYs679aHzTHDgitlovJyprwqrbjg5N3XNSB5FohAUJUnVMF8z+qzvb4pPhly6mj6tiSJGYbXPngN6Iv8t3mRko3PbYLrWuxMb345BxcD+j9XteUgm1j/10qrSvqq+1R+/FAFPYwLXCflZgKst2g8/rEiVQz+a3w== generated-by-azure",
		},
	}

	for _, tc := range cases {
		t.Run(tc.Input, func(t *testing.T) {
			output, err := NormaliseSSHKey(tc.Input)
			if err != nil {
				if !tc.Error {
					t.Fatalf("expected NormaliseSSHKey to error")
				}
			}
			if output != nil && *output != tc.Expected {
				t.Fatalf("Expected %q, got %q", tc.Expected, *output)
			}
		})
	}
}
