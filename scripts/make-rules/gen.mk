
# ==============================================================================
# Makefile helper functions for generate necessary files
#

# 生成iam-apiserver、iam-authz-server、iam-pump、iamctl组件的默认配置文件
.PHONY: gen.defaultconfigs
gen.defaultconfigs:
	@${ROOT_DIR}/scripts/gen_default_config.sh

# 生成 CA 证书
.PHONY: gen.ca.%
gen.ca.%:
	$(eval CA := $(word 1, $(subst ., ,$*)))
	@echo "==============> Generating CA files for $(CA)"
	@${ROOT_DIR}/scripts/gencerts.sh generate-iam-cert $(OUTPUT_DIR)/cert $(CA)

.PHONY: gen.ca
gen.ca: $(addprefix gen.ca., $(CERTIFICATES))