webhooker tests init:

  $ go build webhooker || go build github.com/piranha/webhooker
  $ post() {
  > curl -s --data-urlencode payload@$1 http://localhost:1234/
  > }
  $ ./webhooker -p 1234 \
  > octokitty/testing:master="echo otm" \
  > &

Usage:

  $ ./webhooker
  Usage:
    webhooker [OPTIONS] user/repo:branch=command [more rules...]
  
  Runs specified shell commands on incoming webhook from Github. Shell command
  environment contains:
  
    $REPO - repository name in "user/name" format
    $REPO_URL - full repository url
    $PRIVATE - strings "true" or "false" if repository is private or not
    $BRANCH - branch name
    $COMMIT - last commit hash id
    $COMMIT_MESSAGE - last commit message
    $COMMIT_TIME - last commit timestamp
    $COMMIT_AUTHOR - username of author of last commit
    $COMMIT_URL - full url to commit
  
  
  Application Options:
    -i, --interface= ip to listen on
    -p, --port=      port to listen on (8000)
    -l, --log=       path to file for logging
    -c, --config=    read rules from this file
    -d, --dump       dump rules to console
        --help       show this help message

Check that it works:

  $ post $TESTDIR/example.json
  [\d/: ]+ 'echo otm' for octokitty/testing output: otm (re)

Cool down:

  $ kill $!