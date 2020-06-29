#!/bin/bash

touch ~/.gitcookies
chmod 0600 ~/.gitcookies

git config --global http.cookiefile ~/.gitcookies

tr , \\t <<\__END__ >>~/.gitcookies
go.googlesource.com,TRUE,/,TRUE,2147483647,o,git-kt.katbyte.me=1/sEvv4P2NiGofB7kgPV7DBbsV5V8_od3JULgYIyZJnUM
go-review.googlesource.com,TRUE,/,TRUE,2147483647,o,git-kt.katbyte.me=1/sEvv4P2NiGofB7kgPV7DBbsV5V8_od3JULgYIyZJnUM
__END__
