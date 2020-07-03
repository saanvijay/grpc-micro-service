#
# Makefile to build package files from proto files and to build server & client executables
#

PKG_DIR = scmpb
CLIENT_DIR = scm_client
SERVER_DIR = scm_server

PROTO-FILES = $(wildcard $(PKG_DIR)/*.proto)
PROTO-PB-OUT-FILES := $(patsubst $(PKG_DIR)/%.proto, $(PKG_DIR)/%.pb.go, $(PROTO-FILES))

SERVER_SRC = $(wildcard $(SERVER_DIR)/*.go)
CLIENT_SRC = $(wildcard $(CLIENT_DIR)/*.go)

SERVER_EXE = $(SERVER_DIR)/scm_server
CLIENT_EXE = $(CLIENT_DIR)/scm_client

VPATH := $(SERVER_DIR) $(CLIENT_DIR) $(PKG_DIR)

.PHONY: all
all:   $(PROTO-PB-OUT-FILES) $(SERVER_EXE) $(CLIENT_EXE)  

certs:
	sh ./TLS/generateCerts.sh

.PHONY: clean 
clean:
	rm -f $(PROTO-PB-OUT-FILES)
	rm -f $(SERVER_EXE)
	rm -f $(CLIENT_EXE)

$(PKG_DIR)/%.pb.go: $(PKG_DIR)/%.proto
	@echo "--------- Creating $(notdir $@) ---------"
	protoc --go_out=plugins=grpc:${@D} $^
	mkdir -p $(GOPATH)/src/${@D}
	cp $@ $(GOPATH)/src/${@D}
	@echo ''
	
$(SERVER_EXE) : $(SERVER_SRC) $(PROTO-PB-OUT-FILES)
	@echo "--------- Creating $(notdir $@) ---------"
	cd ${@D}; go build
	@echo ''

$(CLIENT_EXE) : $(CLIENT_SRC) $(PROTO-PB-OUT-FILES)
	@echo "--------- Creating $(notdir $@) ---------"
	cd ${@D}; go build
	@echo "=============== DONE ===================="
	@echo ''
