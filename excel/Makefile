

.PHONY:copy
copy:
	cp -f grpcexcel.proto proto/grpcexcel.proto

.PHONY:run
run:    
	protoc --go_out=plugins=grpc:. proto/grpcexcel.proto

