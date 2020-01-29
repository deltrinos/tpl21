

.PHONY:copy
copy:
	cp -f grpcexcel.proto proto/grpcexcel.proto

.PHONY:copy
run:    copy
	protoc --go_out=plugins=grpc:. proto/grpcexcel.proto

