#!/bin/bash
SCRIPT_PATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

export GOPATH=$SCRIPT_PATH/../../gopath
export GOBIN=$GOPATH/bin

# We may have a test we want to run.
TEST_TO_RUN="$1"

# Get all the libraries we need.
echo -e "\nGET\n"

# Got the right directory.
cd $SCRIPT_PATH/..
[ $? -ne 0 ] && exit 1

# Testing library.
go get -d "gopkg.in/check.v1"
[ $? -ne 0 ] && exit 1

# Debugger. https://github.com/mailgun/godebug
#
# Insert a breakpoint in source with:
#
#   _ = "breakpoint"
#
# ../bin/godebug run -instrument=glemzurg/technical_debt glemzurg/technical_debt/cmd/debug/main.go
#
go get "github.com/mailgun/godebug"
[ $? -ne 0 ] && exit 1

echo -e "\nFMT\n" ; go fmt ./...
[ $? -ne 0 ] && exit 1

# Run unit tests.
# go test -check.f MyTestSuite
# go test -check.f "Test.*Works"
# go test -check.f "MyTestSuite.Test.*Works"
echo -e "\nTEST\n" 
if [ -z "$TEST_TO_RUN" ]; then

  # No explicit test, running all tests.
  go test ./... # -check.f "Test_Process_NewChild" # -check.v # -check.vv #-p=1
  [ $? -ne 0 ] && exit 1

else 

  # An explicit test, run only that.
  go test ./... -check.f "$TEST_TO_RUN" # -check.v # -check.vv #-p=1
  [ $? -ne 0 ] && exit 1

fi 

# Build and install any executables.
echo -e "\nINSTALL\n" ; go install ./...
[ $? -ne 0 ] && exit 1

# Indicate the command.
echo -e "\nLAUNCH COMMAND\n"
echo -e "technical_debt -config /path/to/technical_debt/root/config/config.json\n"
echo -e "convert root/output/grid.svg root/output/grid.png\n"

 

