#!/usr/bin/env bash
# Copyright IBM Corp. 2014, 2026
# SPDX-License-Identifier: MPL-2.0

# Tests changelog-entry-push.sh against throwaway repositories on disk: a good entry is pushed, anything else is refused and the branch is
# left as it was. Needs nothing but git, and touches no network. Run with 'make test-automation'.

set -euo pipefail

script="$(cd "$(dirname "$0")" && pwd)/changelog-entry-push.sh"

tmp=$(mktemp -d)
trap 'rm -rf "${tmp}"' EXIT

export GIT_CONFIG_NOSYSTEM=1
export HOME="${tmp}"
export GIT_AUTHOR_NAME="test" GIT_AUTHOR_EMAIL="test@example.com" GIT_COMMITTER_NAME="test" GIT_COMMITTER_EMAIL="test@example.com"

passed=0

pass() {
  echo "PASS: $1"
  passed=$((passed + 1))
}

fail() {
  echo "FAIL: $1" >&2
  exit 1
}

# the repository the pull request is against, holding pull request 123, and the fork its branch lives in
git init --quiet --bare -b main "${tmp}/base.git"
git init --quiet --bare -b main "${tmp}/fork.git"
git clone --quiet "${tmp}/base.git" "${tmp}/seed" 2>/dev/null
cd "${tmp}/seed"
echo "v1" > README.md
git add README.md
git commit --quiet -m "main"
git push --quiet origin HEAD:refs/heads/main
echo "v2" > README.md
git commit --quiet -am "the pull request"
pr_head=$(git rev-parse HEAD)
git push --quiet origin "HEAD:refs/pull/123/head"
git push --quiet "${tmp}/fork.git" "HEAD:refs/heads/my-branch"

# bundle <name> <commands>: runs the commands on top of the pull request, as one commit unless they commit themselves, and bundles the result
bundle() {
  local name=$1
  local commands=$2
  local dir="${tmp}/work-${name}"

  git clone --quiet "${tmp}/base.git" "${dir}" 2>/dev/null
  cd "${dir}"
  git fetch --quiet origin "refs/pull/123/head"
  git switch --quiet -c changelog-entry "${pr_head}"
  mkdir -p .changelog
  eval "${commands}"
  git add --all
  git diff --cached --quiet || git commit --quiet -m "changelog: add entries for PR #123"
  git bundle create --quiet "${tmp}/${name}.bundle" "${pr_head}..refs/heads/changelog-entry"
}

# push <name>: runs the script on a bundle, against the base repository and the fork
push() {
  BUNDLE_PATH="${tmp}/$1.bundle" PR_NUMBER="${PR_NUMBER:-123}" FETCH_URL="file://${tmp}/base.git" PUSH_URL="file://${tmp}/fork.git" HEAD_REF="${HEAD_REF:-my-branch}" "${script}" >"${tmp}/output" 2>&1
}

# said <text>: the script's last run must have given this as its reason
said() {
  grep --quiet --fixed-strings "$1" "${tmp}/output" || fail "expected the reason \"$1\", got: $(cat "${tmp}/output")"
}

fork_head() {
  git --git-dir="${tmp}/fork.git" rev-parse refs/heads/my-branch
}

# refused <name> <what> <reason>: the script must fail on the bundle for that reason and leave the fork's branch alone
refused() {
  if push "$1"; then
    fail "$2 was pushed"
  fi
  said "$3"
  [[ "$(fork_head)" == "${pr_head}" ]] || fail "$2 was refused but the branch moved"
  pass "$2 is refused"
}

bundle other-file 'echo "- an entry" > .changelog/123.md; echo "evil" >> README.md'
refused other-file "a commit that also changes another file" "not this pull request's changelog entry: README.md"

bundle other-pr 'echo "- an entry" > .changelog/999.md'
refused other-pr "an entry for another pull request" "not this pull request's changelog entry: .changelog/999.md"

bundle nested 'mkdir -p .changelog/123.md.d; echo "x" > .changelog/123.md.d/123.md'
refused nested "a file in a folder under .changelog" "not this pull request's changelog entry: .changelog/123.md.d/123.md"

bundle two-commits 'echo "- one" > .changelog/123.md; git add --all; git commit --quiet -m one; echo "- two" >> .changelog/123.md'
refused two-commits "two commits" "expected one commit"

bundle link 'ln -s ../README.md .changelog/123.md'
refused link "an entry that is a link" "not a plain file"

bundle program 'echo "- an entry" > .changelog/123.md; chmod +x .changelog/123.md'
refused program "an entry that can be run" "not a plain file"

bundle good 'echo "- an entry" > .changelog/123.md'

if PR_NUMBER="123; echo" push good; then
  fail "a pull request number that is not a number was accepted"
fi
said "Invalid pull request number"
pass "a pull request number that is not a number is refused"

if HEAD_REF="my..branch" push good; then
  fail "a branch name git does not allow was accepted"
fi
said "Invalid branch name"
pass "a branch name git does not allow is refused"

# the pull request gains a commit after the entry was written: the entry is no longer on top of it
cd "${tmp}/seed"
echo "v3" > README.md
git commit --quiet -am "a later push to the pull request"
git push --quiet origin "HEAD:refs/pull/123/head"
if push good; then
  fail "an entry written before the pull request's latest commit was pushed"
fi
said "not written on top of the pull request's latest commit"
[[ "$(fork_head)" == "${pr_head}" ]] || fail "a stale entry was refused but the branch moved"
pass "an entry written before the pull request's latest commit is refused"
git push --quiet --force origin "${pr_head}:refs/pull/123/head"

push good || fail "a good entry was refused"
[[ "$(git --git-dir="${tmp}/fork.git" rev-parse "refs/heads/my-branch^")" == "${pr_head}" ]] || fail "a good entry was not pushed on top of the pull request"
[[ "$(git --git-dir="${tmp}/fork.git" show "refs/heads/my-branch:.changelog/123.md")" == "- an entry" ]] || fail "the pushed entry is not what was written"
pass "a good entry is pushed on top of the pull request"

# replacing an entry in one format with one in another removes a file and adds a file
git --git-dir="${tmp}/base.git" fetch --quiet "${tmp}/fork.git" "+refs/heads/my-branch:refs/pull/123/head"
pr_head=$(fork_head)
bundle replace 'echo "change {}" > .changelog/123.hcl; git rm --quiet .changelog/123.md'
push replace || fail "replacing an entry with one in another format was refused"
[[ "$(git --git-dir="${tmp}/fork.git" ls-tree --name-only refs/heads/my-branch .changelog/)" == ".changelog/123.hcl" ]] || fail "the replaced entry was not pushed"
pass "an entry replaced by one in another format is pushed"

echo "==> All ${passed} changelog-entry-push checks passed"
