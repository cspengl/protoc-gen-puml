#!/bin/bash -e

WORKDIR=$(dirname $0)
TEMPDIR=$(mktemp -d)

trap "rm -rf ${TEMPDIR}" EXIT

# Build protoc-gen-puml binary and add to $PATH.
pushd "${WORKDIR}"
go build -o "${TEMPDIR}" cmd/protoc-gen-puml.go
PATH="${TEMPDIR}:${PATH}"
popd


protoc \
	--puml_out="${TEMPDIR}" \
	testdata/test.proto

GOLDENFILE=./diagram.pb.puml
GENFILE="${TEMPDIR}/diagram.pb.puml"

# diff is piped to [[ $? == 1 ]] to avoid exiting on diff but exit on error
# (like if the file was not found). See man diff for more info.
DIFF=$(diff "${GOLDENFILE}" "${GENFILE}" || [[ $? == 1 ]])
if [[ -n "${DIFF}" ]]; then
    echo -e "ERROR: Generated file differs from golden file:\n${DIFF}"
    exit 1
fi