# Borrowed from: 
# https://github.com/silven/go-example/blob/master/Makefile
# https://vic.demuzere.be/articles/golang-makefile-crosscompile/

BINARY = petrockutil
VET_REPORT = vet.report
TEST_REPORT = tests.xml
GOARCH_1 = amd64
GOARCH_2 = 386
GOARCH_3 = arm

VERSION?=?
COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

# Symlink into GOPATH
GITHUB_USERNAME=petrockblog
BUILD_DIR=${GOPATH}/src/github.com/${GITHUB_USERNAME}/${BINARY}
CURRENT_DIR=$(shell pwd)
BUILD_DIR_LINK=$(shell readlink ${BUILD_DIR})

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-X main.VERSION=${VERSION} -X main.COMMIT=${COMMIT} -X main.BRANCH=${BRANCH}"

# Build the project
all: clean linux darwin windows

link:
	BUILD_DIR=${BUILD_DIR}; \
	BUILD_DIR_LINK=${BUILD_DIR_LINK}; \
	CURRENT_DIR=${CURRENT_DIR}; \
	if [ "$${BUILD_DIR_LINK}" != "$${CURRENT_DIR}" ]; then \
	    echo "Fixing symlinks for build"; \
	    rm -f $${BUILD_DIR}; \
	    ln -s $${CURRENT_DIR} $${BUILD_DIR}; \
	fi

linux: 
	cd ${BUILD_DIR}; \
	GOOS=linux GOARCH=${GOARCH_1} go build ${LDFLAGS} -o ${BINARY}-linux-${GOARCH_1} . ; \
	zip ${BINARY}-linux-${GOARCH_1}.zip ${BINARY}-linux-${GOARCH_1} README.md LICENSE; \
	rm ${BINARY}-linux-${GOARCH_1}; \
	GOOS=linux GOARCH=${GOARCH_2} go build ${LDFLAGS} -o ${BINARY}-linux-${GOARCH_2} . ; \
	zip ${BINARY}-linux-${GOARCH_2}.zip ${BINARY}-linux-${GOARCH_2} README.md LICENSE; \
	rm ${BINARY}-linux-${GOARCH_2}; \
	GOOS=linux GOARCH=${GOARCH_3} GOARM=5 go build ${LDFLAGS} -o ${BINARY}-linux-${GOARCH_3} . ; \
	zip ${BINARY}-linux-${GOARCH_3}.zip ${BINARY}-linux-${GOARCH_3} README.md LICENSE; \
	rm ${BINARY}-linux-${GOARCH_3}; \
	cd - >/dev/null

darwin:
	cd ${BUILD_DIR}; \
	GOOS=darwin GOARCH=${GOARCH_1} go build ${LDFLAGS} -o ${BINARY}-darwin-${GOARCH_1} . ; \
	zip ${BINARY}-darwin-${GOARCH_1}.zip ${BINARY}-darwin-${GOARCH_1} README.md LICENSE; \
	rm ${BINARY}-darwin-${GOARCH_1}; \
	GOOS=darwin GOARCH=${GOARCH_2} go build ${LDFLAGS} -o ${BINARY}-darwin-${GOARCH_2} . ; \
	zip ${BINARY}-darwin-${GOARCH_2}.zip ${BINARY}-darwin-${GOARCH_2} README.md LICENSE; \
	rm ${BINARY}-darwin-${GOARCH_2}; \
	cd - >/dev/null

windows:
	cd ${BUILD_DIR}; \
	GOOS=windows GOARCH=${GOARCH_1} go build ${LDFLAGS} -o ${BINARY}-windows-${GOARCH_1}.exe . ; \
	zip ${BINARY}-windows-${GOARCH_1}.zip ${BINARY}-windows-${GOARCH_1}.exe README.md LICENSE; \
	rm ${BINARY}-windows-${GOARCH_1}.exe; \
	GOOS=windows GOARCH=${GOARCH_2} go build ${LDFLAGS} -o ${BINARY}-windows-${GOARCH_2}.exe . ; \
	zip ${BINARY}-windows-${GOARCH_2}.zip ${BINARY}-windows-${GOARCH_2}.exe README.md LICENSE; \
	rm ${BINARY}-windows-${GOARCH_2}.exe; \
	cd - >/dev/null

test:
	if ! hash go2xunit 2>/dev/null; then go install github.com/tebeka/go2xunit; fi
	cd ${BUILD_DIR}; \
	godep go test -v ./... 2>&1 | go2xunit -output ${TEST_REPORT} ; \
	cd - >/dev/null

vet:
	-cd ${BUILD_DIR}; \
	godep go vet ./... > ${VET_REPORT} 2>&1 ; \
	cd - >/dev/null

fmt:
	cd ${BUILD_DIR}; \
	go fmt $$(go list ./... | grep -v /vendor/) ; \
	cd - >/dev/null

clean:
	-rm -f ${TEST_REPORT}
	-rm -f ${VET_REPORT}
	-rm -f ${BINARY}-*

.PHONY: link linux darwin windows test vet fmt clean