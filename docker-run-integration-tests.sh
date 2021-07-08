#!/bin/bash
go test -v -run=.+IntegrationTestsSuite ./... || exit 1
exit 0