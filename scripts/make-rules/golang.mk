
GO := go
GO_SUPPORTED_VERSIONS ?= 1.13|1.14|1.15|1.16|1.17

# -q 或 --quiet或--silent : 不显示任何信息。
# \b 匹配一个单词的边界，也就是指单词和空格间的位置
.PHONY: go.build.verify
go.build.verify:
ifneq ($(shell ${GO} version | grep -q -E '\bgo(${GO_SUPPORTED_VERSIONS})\b' && echo 0 || echo 1), 0)
	$(error unsupported go version. Please make install one of the following supported version: '$(GO_SUPPORTED_VERSIONS)')
endif

.PHONY: go.build.%
go.build.%:
	$(eval COMMAND := $(word 2,$(subst ., ,$*)))
	$(eval PLATFORM := $(word 1,$(subst ., ,$*)))
	$(eval OS := $(word 1,$(subst _, ,${PLATFORM})))
	$(eval ARCH := $(word 2,$(subst _, ,${PLATFORM})))
	@echo "===========> Building binary ${COMMAND} ${VERSION} for ${OS} ${ARCH}"
	@mkdir -p ${OUTPUT_DIR}/platforms/${OS}/${ARCH}
	@CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} ${GO} build ${GO_BUILD_FLAGS} -o ${OUTPUT_DIR}/platforms/${OS}/${ARCH}/${COMMAND}${GO_OUT_EXT} ${ROOT_PACKAGE}/cmd/${COMMAND}



