
# ==============================================================================
# Makefile helper functions for docker image
#

# := 直接赋值，赋予当前位置的值
DOCKER := docker

# ?= 如果变量没有被赋值，则赋予等号后的值
DOCKER_SUPPORTED_API_VERSION ?= 1.32

REGISTRY_PREFIX ?= marmotedu

# 镜像基础操作系统
BASE_IMAGE = centos:centos8

_DOCKER_BUILD_EXTRA_ARGS :=

# 通过查看 build/docker/*/Dockerfile 来确定镜像文件
# $(wildcard <pattern>) 扩展通配符
# 	IMAGES_DIR 的值为：
#		${ROOT_DIR}/build/docker/iam-watcher
#		${ROOT_DIR}/build/docker/iam-pump
#		${ROOT_DIR}/build/docker/iam-iamctl
#		${ROOT_DIR}/build/docker/iam-authz-server
#		${ROOT_DIR}/build/docker/iam-apiserver
IMAGES_DIR ?= $(wildcard ${ROOT_DIR}/build/docker/*)

# 通过去除目录名称来确定镜像名称
# $(notdir <names...>) 从文件名序列中取出目录部分。
# 		非目录部分是指最后一个反斜杠（/）之后的部分。返回文件名序列的非目录部分。
#
# $(foreach <var>,<list>,<text>)
# 	把参数 <list> 中的单词逐一取出，放到参数 <var> 所指定的变量中，然后再执行 <text> 所包含的表达式。
#	每一次 <text> 会返回一个字符串，循环过程中 <text> 所返回的每个字符串会以“空格”分隔，
#	最后当整个循环结束时，<text> 所返回的每个字符串所组成的整个字符串（以空格分隔）将会
#	是 foreach 函数的返回值。
#
#	iam-watcher iam-pump iam-iamctl iam-authz-server iam-apiserver
#
# $(filter-out <pattern>,<text>)
# 	以 <pattern> 模式过滤 <text> 字符串中的单词,去除符合模式 <pattern> 的单词。
# 	可以有多个模式。返回不符合模式 <pattern> 的字串。
#
# tools 是啥？
IMAGES ?= $(filter-out tools,$(foreach image,${IMAGES_DIR}, $(notdir ${image})))

# 验证docker API 版本
.PHONY: image.verify
image.verify:
	# docker API version
	$(eval API_VERSION := $(shell ${DOCKER} version | grep -E 'API version: {1,6}[0-9]' | head -n1 | awk '{print $3} END { if (NR==0) print 0}' ))
	# 当前 docker 版本是否大于 ${DOCKER_SUPPORTED_API_VERSION}
	$(eval PASS := $(shell echo "${API_VERSION} > ${DOCKER_SUPPORTED_API_VERSION}" | bc))
	@if [ ${PASS} -ne 1 ]; then \
		${DOCKER} -v ;\
		echo "Unsupported docker version. Docker API version should be greater than $(DOCKER_SUPPORTED_API_VERSION)"; \
		exit 1; \
	fi

# 验证docker daemon 版本
.PHONY: image.daemon.verify
image.daemon.verify:
	$(eval PASS := $(shell ${DOCKER} version | grep -q -E 'Experimental: {1,5}true' && echo 1 || echo 0))
	@if [ $(PASS) -ne 1 ]; then \
		echo "Experimental features of Docker daemon is not enabled. Please add \"experimental\": true in '/etc/docker/daemon.json' and then restart Docker daemon."; \
		exit 1; \
	fi

.PHONY: image.build
image.build: image.verify go.build.verify $(addprefix image.build., $(addprefix ${IMAGE_PLAT}., ${IMAGES}))

# 1. 验证 docker 版本
# 2. 验证 golang 版本
# 3. image.push.linux_amd64.iam-apiserver
.PHONY: image.push
image.push: image.verify go.build.verify $(addprefix image.push., $(addprefix ${IMAGE_PLAT}., ${IMAGES}))

# 依赖 image.build.%
.PHONY: image.push.%
image.push.%: image.build.%
	@echo "===========> Pushing image ${IMAGE} ${VERSION} to ${REGISTRY_PREFIX}"
	${DOCKER} push ${REGISTRY_PREFIX}/${IMAGE}-${ARCH}:${VERSION}

# image.build.linux_amd64.iam-apiserver
# 依赖 go.build.%
.PHONY: image.build.%
image.build.%: go.build.%
	# go.build.% 时定义
	$(eval IMAGE := ${COMMAND})
	# linux/amd64
	$(eval IMAGE_PLAT := $(subst _,/,${PLATFORM}))
	@echo "===========> Building docker image ${IMAGE} ${VERSION} for ${IMAGE_PLAT}"
	# TMP_DIR 在 common.mk 中定义
	@mkdir -p ${TMP_DIR}/${IMAGE}
	# 生成 Dockerfile (修改 Dockerfile 中的 FROM 基础镜像)
	@cat ${ROOT_DIR}/build/docker/${IMAGE}/Dockerfile\
		| sed "s#BASE_IMAGE#${BASE_IMAGE}#g" >${TMP_DIR}/${IMAGE}/Dockerfile
	# 拷贝构建好的二进制文件
	@cp ${OUTPUT_DIR}/platforms/${IMAGE_PLAT}/${IMAGE} ${TMP_DIR}/${IMAGE}/
	@DST_DIR=${TMP_DIR}/${IMAGE} ${ROOT_DIR}/build/docker/${IMAGE}/build.sh 2>/dev/null || true
	$(eval BUILD_SUFFIX := ${_DOCKER_BUILD_EXTRA_ARGS} --pull -t ${REGISTRY_PREFIX}/${IMAGE}-${ARCH}:${VERSION} ${TMP_DIR}/${IMAGE})
	@if [ $(shell ${GO} env GOARCH) != ${ARCH} ] ; then \
		${MAKE} image.daemon.verify ;\
		${DOCKER} build --platform ${IMAGE_PLAT} ${BUILD_SUFFIX} ; \
	else \
		${DOCKER} build ${BUILD_SUFFIX} ; \
	fi
	@rm -rf ${TMP_DIR}/${IMAGE}