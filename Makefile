build: create_bindata_ui

download_webfiles:
	rm -rf .temp/pip_packages
	rm -rf .temp/webfiles
	mkdir -p .temp/pip_packages
	cd .temp/pip_packages && pip3 download tensorboard && unzip tensorboard*.whl tensorboard/webfiles.zip
	mkdir -p .temp/webfiles
	unzip .temp/pip_packages/tensorboard/webfiles.zip -d .temp/webfiles
	rm .temp/webfiles/trace_viewer_index*
.PHONY: download_webfiles

create_bindata_ui: download_webfiles go-bindata-assetfs
	mkdir -p ui/
	$(BINDATA_ASSETFS) -pkg "ui" -prefix ".temp/webfiles" -o ui/tensorboard.go .temp/webfiles
.PHONY: create_bindata_ui


# find or download go-bindata-assetfs 
# download go-bindata-assetfs  if necessary
go-bindata-assetfs:
ifeq (, $(shell which go-bindata-assetfs ))
	@{ \
	set -e ;\
	BINDAT_TEMP_DIR=$$(mktemp -d) ;\
	cd $$BINDAT_TEMP_DIR ;\
	go mod init tmp ;\
	go get github.com/go-bindata/go-bindata/... ; \
	go get github.com/elazarl/go-bindata-assetfs/... ;\
	rm -rf $$BINDAT_TEMP_DIR ;\
	}
BINDATA_ASSETFS=$(GOBIN)/go-bindata-assetfs 
else
BINDATA_ASSETFS=$(shell which go-bindata-assetfs )
endif
