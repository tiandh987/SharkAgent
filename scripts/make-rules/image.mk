
# ==============================================================================
# Makefile helper functions for docker image
#

DOCKER := docker
DOCKER_SUPPORTED_API_VERSION ?= 1.32

# 通过查看 build/docker/*/Dockerfile 来确定镜像文件
IMAGES_DIR ?= $(wildcard ${ROOT_DIR}/build/docker/*)
# 通过去除目录名称来确定镜像名称
IMAGES ?= $(filter-out toos,$(foreach image,${IMAGES_DIR}, $(notdir ${image})))

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

.PHONY: image.build
image.build: image.verify go.build.verify $(addprefix image.build., $(addprefix ${IMAGE_PLAT}., ${IMAGES}))

.PHONY: image.build.%
image.build.%: go.build.%
	$(eval IMAGE := ${COMMAND})
	$(eval IMAGE_PLAT := $(subst _,/,${PLATFORM}))
	@echo "===========> Building docker image ${IMAGE} ${VERSION} for ${IMAGE_PLAT}"
	@mkdir -p ${TMP_DIR}/${IMAGE}
	@cat ${ROOT_DIR}/build/docker/${IMAGE}/Dockerfile\
		| sed "s#BASE_IMAGE#${BASE_IMAGE}#g" >${TMP_DIR}/${IMAGE}/Dockerfile
	@cp ${OUTPUT_DIR}/platforms/${IMAGE_PLAT}/${IMAGE} ${TMP_DIR}/${IMAGE}/
	@DST_DIR=${TMP_DIR}/${IMAGE} ${ROOT_DIR}/build/docker/${IMAGE}/build.sh 2>/dev/null || true
	$(eval BUILD_SUFFIX := ${_DOCKER_BUILD_EXTRA_ARGS} --pill -t ${REGISTRY_PREFIX}/${IMAGE}-${ARCH}:${VERSION} ${TMP_DIR}/${IMAGE})
	@if [ $(shell ${GO} env GOARCH) != ${ARCH} ] ; then \
  		${MAKE} image.daemon.verify ;\
  		${DOCKER} build --platform ${IMAGE_PLAT} ${BUILD_SUFFIX} ; \
  	else \
  		${DOCKER} build ${BUILD_SUFFIX} ; \
  	fi
  	@rm -rf $(TMP_DIR)/$(IMAGE)