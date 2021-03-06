#!/usr/bin/env cram

webhooker tests init:

  $ go build github.com/sainaen/webhooker
  $ githubPost1() {
  > curl -s --header "X-Github-Event: push" --header "Content-Type: application/x-www-form-urlencoded" --data-urlencode payload@$1 http://localhost:1234/
  > }
  $ githubPost2() {
  > curl -s --header "X-Github-Event: push" --header "Content-Type: application/json" --data @$1 http://localhost:1234/
  > }
  $ bitbucketPost() {
  > curl -s --header "X-Event-Key: repo:push" --header "Content-Type: application/json" --data @$1 http://localhost:1234/
  > }
  $ ./webhooker -p 1234 \
  > 'octokitty/testing:master=echo OTM' \
  > '.*=echo $REPO' \
  > &

Usage:

  $ ./webhooker
  Usage:
    webhooker [OPTIONS] user/repo:branch=command [more rules...]
  
  Runs specified shell commands on incoming webhook from Github. Shell command
  environment contains:
  
    $PATH - proxied from parent environment
    $HOME - proxied from parent environment
    $USER - proxied from parent environment
  
    $REPO - repository name in "user/name" format
    $REPO_URL - full repository url
    $PRIVATE - strings "true" or "false" if repository is private or not
    $BRANCH - branch name
    $COMMIT - last commit hash id
    $COMMIT_MESSAGE - last commit message
    $COMMIT_TIME - last commit timestamp
    $COMMIT_AUTHOR - username of author of last commit
    $COMMIT_URL - full url to commit
  
  'user/repo:branch' pattern is a regular expression, so you could do
  'user/project:fix.*=cmd' or even '.*=cmd'.
  
  Also never forget to properly escape your rule, if you pass it through command
  line: usually enclosing it in single quotes (') is enough.
  
  
  Application Options:
    -i, --interface= ip to listen on
    -p, --port=      port to listen on (default: 8000)
    -l, --log=       path to file for logging
    -c, --config=    read rules from this file
    -d, --dump       dump rules to console
    -V, --version    show version and exit
        --help       show this help message

Check that it works:

  $ githubPost1 $TESTDIR/gh-example.json
  [\d/: ]+ 'echo OTM' for octokitty/testing output: OTM (re)
  $ githubPost2 $TESTDIR/gh-other.json
  [\d/: ]+ 'echo \$REPO' for hellothere/other output: hellothere/other (re)
  $ bitbucketPost $TESTDIR/bb-example.json
  [\d/: ]+ 'echo OTM' for octokitty/testing output: OTM (re)
  $ bitbucketPost $TESTDIR/bb-other.json
  [\d/: ]+ 'echo \$REPO' for hellothere/other output: hellothere/other (re)

Cool down:

  $ kill $!
