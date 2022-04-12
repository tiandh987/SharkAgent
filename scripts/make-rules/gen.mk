
# ==============================================================================
# Makefile helper functions for generate necessary files
#

# 生成iam-apiserver、iam-authz-server、iam-pump、iamctl组件的默认配置文件
.PHONY: gen.defaultconfigs
gen.defaultconfigs:
	@${ROOT_DIR}/scripts/gen_default_config.sh

# 生成 CA 证书
# $(word <n>,<text>)
#	取字符串 <text> 中第 <n> 个单词(从 1 开始), 并返回字符串<text>中第<n>个单词.
#	如<n>比<text>中的单词数要大,那么返回空字符串
#
# $(subst <from>,<to>,<text>)
# 	把字串<text>中的 <from> 字符串替换成 <to>, 并返回被替换后的字符串
#
# $*
#	$* 分别为 iam-apiserver、iam-auth-server、admin
#   CA 分别为 iam-apiserver、iam-auth-server、admin
#
# :=
#	直接赋值, 赋予当前位置的值
.PHONY: gen.ca.%
gen.ca.%:
	$(eval CA := $(word 1, $(subst ., ,$*)))
	@echo "==============> Generating CA files for ${CA}"
	@${ROOT_DIR}/scripts/gencerts.sh generate-iam-cert ${OUTPUT_DIR}/cert ${CA}

# CERTIFICATES 在 common.mk 中设置
# $(addprefix <prefix>,<names...>)
#	把前缀<prefix>加到<names>中的每个单词前面，并返回加过前缀的文件名序列
# 相当于是 gen.ca: gen.ca.iam-apiserver gen.ca.iam-auth-server gen.ca.admin
.PHONY: gen.ca
gen.ca: $(addprefix gen.ca., ${CERTIFICATES})