
.PHONY: kubeops bindata

all: bindata kubeops

kubeops: .
	godep go build -o bin/kubeops main.go

bindata: .
	cd scripts && go-bindata -nocompress -debug -pkg scripts -o ../pkg/scripts/bindata.go .
