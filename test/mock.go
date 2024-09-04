package test

//go:generate go run github.com/vektra/mockery/v2@v2.45.0 --case=snake  --dir=../internal/loadbalancer --name= --case=snake --all --output=mocks
//go:generate go run github.com/vektra/mockery/v2@v2.45.0 --case=snake  --dir=../internal/echo --name= --case=snake --all --output=mocks
//go:generate go run github.com/vektra/mockery/v2@v2.45.0 --case=snake --srcpkg net/http --name RoundTripper --case=snake --output=mocks
