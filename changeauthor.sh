#!/bin/sh

git filter-branch -f --env-filter '
OLD_EMAIL="joakim.hansson@webstep.se"
CORRECT_NAME="Joakim Hansson"
CORRECT_EMAIL="joakimhew@gmail.com"
if [ "$GIT_COMMITTER_EMAIL" = "$OLD_EMAIL" ]
then
    export GIT_COMMITTER_NAME="$CORRECT_NAME"
    export GIT_COMMITTER_EMAIL="$CORRECT_EMAIL"
fi
if [ "$GIT_AUTHOR_EMAIL" = "$OLD_EMAIL" ]
then
    export GIT_AUTHOR_NAME="$CORRECT_NAME"
    export GIT_AUTHOR_EMAIL="$CORRECT_EMAIL"
fi
' --tag-name-filter cat 3bb02bd2d9189a446ec6f6a12fda6636cb17f5ea..HEAD
