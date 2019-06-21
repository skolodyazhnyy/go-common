.PHONY: dep run test build check format vet rpm

GO := go
NAME := go-common
OS := $(shell uname)
PACKAGE_DIRS := $(shell $(GO) list ./... | grep -v /vendor/)
PKGS := $(shell go list ./... | grep -v /vendor)
PKGS := $(subst  :,_,$(PKGS))

dep:
	go get ./...

lint:
	golangci-lint run

test:
	go test ./...

coverage:
	$(eval PACKAGES_COVERAGE := $(shell $(GO) test $(PACKAGE_DIRS) --cover | awk '{if ($$1 != "?") print $$5; else print "0.0";}' | sed 's/\%//g' | awk '{s+=$$1} END {printf "%.2f\n", s}'))
	$(eval PACKAGES_NUM := $(shell $(GO) test $(PACKAGE_DIRS) --cover | wc -l))
	$(eval TOTAL_COVERAGE := $(shell echo $(PACKAGES_COVERAGE)\/$(PACKAGES_NUM) | bc))
	@ printf "Total coverage of %s: %s%%\n" "$(NAME)" "$(TOTAL_COVERAGE)"
	# @ echo "YVALUE=$(TOTAL_COVERAGE)" > report.properties
