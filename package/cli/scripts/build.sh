#!/bin/bash
set -e

VERSION=${1:-v1.0.0}
COMMIT=${2:-$(git rev-parse --short HEAD 2>/dev/null || echo "dev")}
DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

echo "Building bank CLI ${VERSION} (${COMMIT})"

OUTPUT_DIR="dist"
rm -rf ${OUTPUT_DIR}
mkdir -p ${OUTPUT_DIR}

PLATFORMS="linux/amd64 linux/arm64 darwin/amd64 darwin/arm64"

for PLATFORM in ${PLATFORMS}; do
    OS=$(echo ${PLATFORM} | cut -d'/' -f1)
    ARCH=$(echo ${PLATFORM} | cut -d'/' -f2)
    
    EXT=""
    if [ "${OS}" = "windows" ]; then
        EXT=".exe"
    fi
    
    FILENAME="bank-${VERSION}-${OS}-${ARCH}${EXT}"
    
    echo "Building ${FILENAME}..."
    
    GOOS=${OS} GOARCH=${ARCH} go build \
        -ldflags "-s -w -X github.com/skygenesisenterprise/aether-bank/cli/cmd.version=${VERSION} -X github.com/skygenesisenterprise/aether-bank/cli/cmd.buildDate=${DATE}" \
        -o "${OUTPUT_DIR}/${FILENAME}" .
    
    if [ "${OS}" != "windows" ]; then
        chmod +x "${OUTPUT_DIR}/${FILENAME}"
    fi
done

echo ""
echo "Build complete. Output in ${OUTPUT_DIR}/:"
ls -lh ${OUTPUT_DIR}

TARFILE="bank-${VERSION}-cli.tar.gz"
tar -czf ${TARFILE} -C ${OUTPUT_DIR} .
echo ""
echo "Created ${TARFILE}"