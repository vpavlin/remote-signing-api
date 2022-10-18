VERSION=v0.1
PUSH_TARGET=quay.io/rubixlife/remote-signing-api:$(VERSION)
QUAY_USER=vpavlin0
QUAY_PASSWORD=



generate:
	go generate ./...
fmt:
	go fmt ./...

serve:
	go run server/server.go config.json

build: fmt
	go build -o remote-signing-api server/server.go 

container:
	podman build -t rubixlife/remote-signing-api .

push:
	podman tag rubixlife/remote-signing-api $(PUSH_TARGET)
	podman push $(PUSH_TARGET)

login:
	podman login -u $(QUAY_USER) -p $(QUAY_PASSWORD) quay.io

release-pr: login build container push

run-container:
	podman run -it --rm --security-opt seccomp=unconfined \
			 -v $PWD/data:/opt/remote-signing-api/data:z \
			 -v $PWD/config.container.json:/opt/remote-signing-api/config.json:z \
			 -v $PWD/localhost.crt:/opt/remote-signing-api/localhost.crt:z \
			 -v $PWD/localhost.key:/opt/remote-signing-api/localhost.key:z \
			 rubixlife/remote-signing-api:v0.1