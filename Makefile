

# ==============================================================================
## Includes
## 为了使项目 Makefile 层次化、结构化，通过引入其它 Makefile 实现
include scripts/make-rules/common.mk
include scripts/make-rules/golang.mk
include scripts/make-rules/tools.mk
include scripts/make-rules/gen.mk
include scripts/make-rules/image.mk


# ==============================================================================
# Targets
## push: 构建镜像，并将镜像 Push 到 CCR 镜像仓库
.PHONY: push
push:
	@${MAKE} image.push