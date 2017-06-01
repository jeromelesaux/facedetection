CC=go
RM=rm
MV=mv


SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')
#GOPATH=$(SOURCEDIR)/
GOOS=linux
GOARCH=amd64
#GOARCH=arm
GOARM=7


EXEC=facedetection

VERSION=1.0.0
BUILD_TIME=`date +%FT%T%z`
PACKAGES := fmt github.com/harrydb/go/img/grayscale github.com/disintegration/imaging


LIBS=

LDFLAGS=

.DEFAULT_GOAL:= $(EXEC)

$(EXEC): organize $(SOURCES)
		@echo "    Compilation des sources ${BUILD_TIME}"
		@if  [ "arm" = "${GOARCH}" ]; then\
		    GOOS=${GOOS} GOARCH=${GOARCH} GOARM=${GOARM} go build ${LDFLAGS} -o ${EXEC}-${VERSION} $(SOURCEDIR)/main.go;\
		else\
            GOOS=${GOOS} GOARCH=${GOARCH} GOARM=${GOARM} go build ${LDFLAGS} -o ${EXEC}-${VERSION} $(SOURCEDIR)/main.go;\
        fi
		@echo "    ${EXEC}-${VERSION} generated."

deps: init
		@echo "    Download packages"
		@$(foreach element,$(PACKAGES),go get -d -v $(element);)

organize: deps
		@echo "    Go FMT"
		@$(foreach element,$(SOURCES),go fmt $(element);)

init: clean
		@echo "    Init of the project"

execute:
		./${EXEC}-${VERSION}

clean:
		@if [ -f "${EXEC}-${VERSION}" ] ; then rm ${EXEC}-${VERSION} ; fi
		@echo "    Nettoyage effectuee"

package:  ${EXEC}
		@zip -r ${EXEC}-${GOOS}-${GOARCH}-${VERSION}.zip ./${EXEC}-${VERSION} resources
		@echo "    Archive ${EXEC}-${GOOS}-${GOARCH}-${VERSION}.zip created"

audit:   ${EXEC}
		@go tool vet -all -shadow ./
		@echo "    Audit effectue"

test: $(EXEC)
		@GOOS=${GOOS} GOARCH=${GOARCH} go test ./...
		@echo " Tests OK."