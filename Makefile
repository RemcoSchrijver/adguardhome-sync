# Run go fmt against code
fmt:
	golangci-lint run --fix

# Run go mod tidy
tidy:
	go mod tidy

# Run tests
test: fmt test-ci

# Run ci tests
test-ci: mocks tidy
	go test ./...  -coverprofile=coverage.out
	go tool cover -func=coverage.out

mocks: mockgen
	mockgen -package client -destination pkg/mocks/client/mock.go github.com/bakito/adguardhome-sync/pkg/client Client

release: semver
	@version=$$(semver); \
	git tag -s $$version -m"Release $$version"
	goreleaser --rm-dist

test-release:
	goreleaser --skip-publish --snapshot --rm-dist

semver:
ifeq (, $(shell which semver))
 $(shell go install github.com/bakito/semver@latest)
endif

mockgen:
ifeq (, $(shell which mockgen))
 $(shell go install github.com/golang/mock/mockgen@v1.6.0)
endif

start-replica:
	podman run --pull always --rm -it -p 9090:80 -p 9091:3000 --name adgardhome-replica adguard/adguardhome

check_defined = \
    $(strip $(foreach 1,$1, \
        $(call __check_defined,$1,$(strip $(value 2)))))
__check_defined = \
    $(if $(value $1),, \
      $(error Undefined $1$(if $2, ($2))))

build-image:
	$(call check_defined, AGH_SYNC_VERSION)
	podman build --build-arg VERSION=${AGH_SYNC_VERSION} --build-arg BUILD=$(shell date -u +'%Y-%m-%dT%H:%M:%S.%3NZ') -t ghcr.io/bakito/adguardhome-sync:${AGH_SYNC_VERSION} .

# go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.9.1
model:
	oapi-codegen -package model -generate types https://raw.githubusercontent.com/AdguardTeam/AdGuardHome/v0.107.5/openapi/openapi.yaml > pkg/client/model/model.go


diff-model:
	wget -q https://raw.githubusercontent.com/AdguardTeam/AdGuardHome/v0.107.0/openapi/openapi.yaml -O a.yaml
	wget -q https://raw.githubusercontent.com/AdguardTeam/AdGuardHome/v0.107.5/openapi/openapi.yaml -O b.yaml
	diff a.yaml b.yaml || rm -f a.yaml b.yaml

diff-replica:
	podman cp  adgardhome-replica:/opt/adguardhome/conf/AdGuardHome.yaml tmp/current-config.yaml.tmp
	cat tmp/current-config.yaml.tmp | grep -v " id: " | grep -v " password: " > tmp/current-config.yaml
	diff tmp/reference-config.yaml tmp/current-config.yaml
