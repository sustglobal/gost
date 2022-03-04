# GOST End-to-End Test
The end-to-end tests run as part of the gost repo CI, but can also be run locally.
To run locally, perform the following steps:
1. Start the PubSub emulator
    `docker-compose up`
2. Run `go test -v .` from the e2e_test directory

