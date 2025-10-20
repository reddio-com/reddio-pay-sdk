module order-system

go 1.24

toolchain go1.24.5

require (
	github.com/gorilla/mux v1.8.1
	github.com/mattn/go-sqlite3 v1.14.17
	github.com/reddio-com/reddio-pay-sdk/go-sdk v0.0.0
)

require (
	github.com/sirupsen/logrus v1.9.3 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
)

replace github.com/reddio-com/reddio-pay-sdk/go-sdk => ../../go-sdk
