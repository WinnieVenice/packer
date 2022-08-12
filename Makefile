build:
	go get google.golang.org/protobuf/cmd/protoc-gen-go
	go get google.golang.org/grpc/cmd/protoc-gen-go-grpc

crawl:
	-mkdir idl
	wget https://raw.githubusercontent.com/goodguy-project/goodguy-crawl/main/crawl_service/crawl_service.proto -O idl/crawl_service.proto
	protoc -I. --go_out=. --go_opt=Midl/crawl_service.proto=./idl --go-grpc_out=. --go-grpc_opt=Midl/crawl_service.proto=./idl ./idl/crawl_service.proto
