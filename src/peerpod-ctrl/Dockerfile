# Build the manager binary
FROM --platform=$TARGETPLATFORM quay.io/confidential-containers/golang-fedora:1.21.8-38 as builder
ARG TARGETOS
ARG TARGETARCH
ARG CGO_ENABLED=1
ARG GOFLAGS

WORKDIR /workspace/
RUN if [ "$CGO_ENABLED" = 1 ] ; then dnf install -y libvirt-devel && dnf clean all; fi

COPY ./cloud-providers ./cloud-providers
COPY ./peerpod-ctrl ./peerpod-ctrl

WORKDIR /workspace/peerpod-ctrl
RUN CC=gcc make ARCH=$TARGETARCH build


# Target Image
FROM --platform=$TARGETPLATFORM registry.fedoraproject.org/fedora:38
ARG CGO_ENABLED=1

RUN if [ "$CGO_ENABLED" = 1 ] ; then dnf install -y libvirt-libs openssh-clients && dnf clean all; fi
WORKDIR /
COPY --from=builder /workspace/peerpod-ctrl/bin/manager .

ENTRYPOINT ["/manager"]
