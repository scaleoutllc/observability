CGO_ENABLED=0
PROJECTS := service-email service-monitor service-panel service-user

.PHONY: all cluster build $(PROJECTS) clean

all: cluster deploy

cluster: clean
ifndef HONEYCOMB_API_KEY
	$(error HONEYCOMB_API_KEY is undefined)
endif
	kind create cluster -n observe
	@kubectl --context kind-observe create namespace honeycomb
	@kubectl --context kind-observe create secret generic honeycomb --from-literal=HONEYCOMB_API_KEY=$(HONEYCOMB_API_KEY) -n honeycomb

build: $(PROJECTS)

$(PROJECTS):
	cd projects && go fmt ./$@/...
	cd projects && go vet ./$@/...
	cd projects && go build -o ./$@/dist/main ./$@/...
	cd projects && GOOS=linux GOARCH=amd64 go build -o ./$@/dist/main-for-container ./$@/...
	cd ./projects/$@ && docker build . -t $@:latest
	kind load docker-image $@:latest -n observe

deploy: build
	kubectl --context kind-observe apply -k .

clean:
	kind delete cluster -n observe