//go:generate protoc -I${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.0.1/third_party/googleapis -I ./../../../ --go_out ./pb --go-grpc_out ./pb ./../../../api/event_service.proto
//go:generate protoc -I${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.0.1/third_party/googleapis -I ./../../../ --grpc-gateway_out ./pb --grpc-gateway_opt logtostderr=true --grpc-gateway_opt generate_unbound_methods=true ./../../../api/event_service.proto
package grpc
