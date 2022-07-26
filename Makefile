VERSION=v0.1
PUSH_TARGET=quay.io/rubixlife/remote-signing-api:$(VERSION)
QUAY_USER=vpavlin0
QUAY_PASSWORD=
LABEL=
CONFIG_FILE=config.json


generate:
	go generate ./...
fmt:
	go fmt ./...

serve:
	go run server/server.go $(CONFIG_FILE)

build: fmt
	go build -o remote-signing-api server/server.go 

container:
	podman build -t rubixlife/remote-signing-api $(LABEL) .

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

release-tag:
ifeq ($(shell git branch --show-current | grep -e main -e master),)
	@echo "Not on main/master branch!"
	@exit 1
else
	@git pull
	@git log HEAD^..HEAD
	@printf -- '=%.0s' {1..100}
	$(eval new_tag=$(shell git tag -l --sort=-v:refname | head -1 | awk -F. '{$$NF = $$NF + 1;} 1' | sed 's/ /./g'))
	@echo -ne "\n\n===> Is this the commit you'd like to release as $(new_tag)? [y/N] " && read ans && [ $${ans:-N} = y ]
	@git tag $(new_tag)
	@git push --tags
endif
