.PHONY: run
run:
	go run cmd/go-chan/main.go

.PHONY: generate
generate: generate-go generate-finalize

.PHONY: generate-go
generate-go:
	rm -rf pkg/go_chan_api
	mkdir pkg/go_chan_api
	protoc --go_out=pkg/go_chan_api --go_opt=paths=source_relative \
	--go-grpc_out=pkg/go_chan_api --go-grpc_opt=paths=source_relative ./api/gempellm/go_chan_api/v1/go_chan_api.proto
	# protoc -I . --grpc-gateway_out=pkg/go_chan_api \
    # --grpc-gateway_opt logtostderr=true \
    # --grpc-gateway_opt paths=source_relative \
    # --grpc-gateway_opt generate_unbound_methods=true \
    # ./api/gempellm/go_chan_api/v1/go_chan_api.proto

.PHONE: generate-finalize
generate-finalize:
	mv pkg/go_chan_api/api/gempellm/go_chan_api/v1/* pkg/go_chan_api
	rm -rf pkg/go_chan_api/api
	cd pkg/go_chan_api && go mod init github.com/gempellm/golang-chan/pkg/go_chan_api
	cd pkg/go_chan_api && go mod tidy
