
GO := go

GO_SUPPORTED_VERSIONS ?= 1.13|1.14|1.15|1.16|1.17

# 应用版本信息
GO_LDFLAGS += -X ${VERSION_PACKAGE}.GitVersion=${VERSION} \
	-X $(VERSION_PACKAGE).GitCommit=$(GIT_COMMIT) \
	-X $(VERSION_PACKAGE).GitTreeState=$(GIT_TREE_STATE) \
	-X $(VERSION_PACKAGE).BuildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

# dlv
ifneq (${DLV},)
	GO_BUILD_FLAGS += -gcflags "all=-N -l"
	LDFLAGS = ""
endif
GO_BUILD_FLAGS += -tags=jsoniter -ldflags "${GO_LDFLAGS}"

# 编译文件后缀
ifeq ($(GOOS),windows)
	GO_OUT_EXT := .exe
endif

# -q 或 --quiet或--silent : 不显示任何信息。
# \b 匹配一个单词的边界，也就是指单词和空格间的位置
.PHONY: go.build.verify
go.build.verify:
ifneq ($(shell ${GO} version | grep -q -E '\bgo(${GO_SUPPORTED_VERSIONS})\b' && echo 0 || echo 1), 0)
	$(error unsupported go version. Please make install one of the following supported version: '$(GO_SUPPORTED_VERSIONS)')
endif

# go.build.linux_amd64.iam-apiserver
#
# $*
# 	linux_amd64.iam-apiserver
#
# $(subst <from>,<to>,<text>)
# 	把字串 <text> 中的 <from> 字符串替换成 <to>, 并返回被替换后的字符串
#
# $(word <n>, <text>)
# 	取字符串 <text> 中的第 <n> 个单词(从 1 开始), 并返回第 <n> 个单词
# 	如果 <n> 比 <text> 中的单词数要大, 那么返回空字符串
#
.PHONY: go.build.%
go.build.%:
	# iam-apiserver
	$(eval COMMAND := $(word 2,$(subst ., ,$*)))
	# linux_amd64
	$(eval PLATFORM := $(word 1,$(subst ., ,$*)))
	# linux
	$(eval OS := $(word 1,$(subst _, ,${PLATFORM})))
	# amd64
	$(eval ARCH := $(word 2,$(subst _, ,${PLATFORM})))
	@echo "===========> Building binary ${COMMAND} ${VERSION} for ${OS} ${ARCH}"
	@mkdir -p ${OUTPUT_DIR}/platforms/${OS}/${ARCH}
	@CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} ${GO} build ${GO_BUILD_FLAGS} -o ${OUTPUT_DIR}/platforms/${OS}/${ARCH}/${COMMAND}${GO_OUT_EXT} ${ROOT_PACKAGE}/cmd/${COMMAND}



