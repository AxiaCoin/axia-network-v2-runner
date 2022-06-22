#!/usr/bin/env bash
set -e

export RUN_E2E="true"
# e.g.,
# ./scripts/tests.e2e.sh 1.7.9 1.7.10
if ! [[ "$0" =~ scripts/tests.e2e.sh ]]; then
  echo "must be run from repository root"
  exit 255
fi

VERSION_1=$1
if [[ -z "${VERSION_1}" ]]; then
  echo "Missing version argument!"
  echo "Usage: ${0} [VERSION_1] [VERSION_2]" >> /dev/stderr
  exit 255
fi

VERSION_2=$2
if [[ -z "${VERSION_2}" ]]; then
  echo "Missing version argument!"
  echo "Usage: ${0} [VERSION_1] [VERSION_2]" >> /dev/stderr
  exit 255
fi

echo "Running e2e tests with:"
echo VERSION_1: ${VERSION_1}
echo VERSION_2: ${VERSION_2}

############################
# download axia
# https://github.com/axiacoin/axia/releases
GOARCH=$(go env GOARCH)
GOOS=$(go env GOOS)
DOWNLOAD_URL=https://github.com/axiacoin/axia/releases/download/v${VERSION_1}/axia-linux-${GOARCH}-v${VERSION_1}.tar.gz
DOWNLOAD_PATH=/tmp/axia.tar.gz
if [[ ${GOOS} == "darwin" ]]; then
  DOWNLOAD_URL=https://github.com/axiacoin/axia/releases/download/v${VERSION_1}/axia-macos-v${VERSION_1}.zip
  DOWNLOAD_PATH=/tmp/axia.zip
fi

rm -rf /tmp/axia-v${VERSION_1}
rm -rf /tmp/axia-build
rm -f ${DOWNLOAD_PATH}

echo "downloading axia ${VERSION_1} at ${DOWNLOAD_URL}"
curl -L ${DOWNLOAD_URL} -o ${DOWNLOAD_PATH}

echo "extracting downloaded axia"
if [[ ${GOOS} == "linux" ]]; then
  tar xzvf ${DOWNLOAD_PATH} -C /tmp
elif [[ ${GOOS} == "darwin" ]]; then
  unzip ${DOWNLOAD_PATH} -d /tmp/axia-build
  mv /tmp/axia-build/build /tmp/axia-v${VERSION_1}
fi
find /tmp/axia-v${VERSION_1}

############################
# download axia
# https://github.com/axiacoin/axia/releases
DOWNLOAD_URL=https://github.com/axiacoin/axia/releases/download/v${VERSION_2}/axia-linux-${GOARCH}-v${VERSION_2}.tar.gz
if [[ ${GOOS} == "darwin" ]]; then
  DOWNLOAD_URL=https://github.com/axiacoin/axia/releases/download/v${VERSION_2}/axia-macos-v${VERSION_2}.zip
  DOWNLOAD_PATH=/tmp/axia.zip
fi

rm -rf /tmp/axia-v${VERSION_2}
rm -rf /tmp/axia-build
rm -f ${DOWNLOAD_PATH}

echo "downloading axia ${VERSION_2} at ${DOWNLOAD_URL}"
curl -L ${DOWNLOAD_URL} -o ${DOWNLOAD_PATH}

echo "extracting downloaded axia"
if [[ ${GOOS} == "linux" ]]; then
  tar xzvf ${DOWNLOAD_PATH} -C /tmp
elif [[ ${GOOS} == "darwin" ]]; then
  unzip ${DOWNLOAD_PATH} -d /tmp/axia-build
  mv /tmp/axia-build/build /tmp/axia-v${VERSION_2}
fi
find /tmp/axia-v${VERSION_2}

############################
echo "building runner"
go build -v -o /tmp/network.runner ./cmd/axia-network-runner

echo "building e2e.test"
# to install the ginkgo binary (required for test build and run)
go install -v github.com/onsi/ginkgo/v2/ginkgo@v2.1.3
ACK_GINKGO_RC=true ginkgo build ./tests/e2e
./tests/e2e/e2e.test --help

echo "launch local test cluster in the background"
/tmp/network.runner \
server \
--log-level debug \
--port=":8080" \
--grpc-gateway-port=":8081" &
PID=${!}

echo "running e2e tests"
./tests/e2e/e2e.test \
--ginkgo.v \
--log-level debug \
--grpc-endpoint="0.0.0.0:8080" \
--grpc-gateway-endpoint="0.0.0.0:8081" \
--axia-path-1=/tmp/axia-v${VERSION_1}/axia \
--axia-path-2=/tmp/axia-v${VERSION_2}/axia

kill -9 ${PID}
echo "ALL SUCCESS!"
