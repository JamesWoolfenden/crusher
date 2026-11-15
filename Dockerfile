FROM alpine:latest@sha256:5b10f432ef3da1b8d4c7eb6c487f2f5a8f096bc91145e68878dd4a5019afde11

# Release tag to install (e.g. v1.2.3); passed by release.yml. No default on purpose.
ARG VERSION
# No default here: a default overrides BuildKit's automatic value (arm64 builds got amd64).
ARG TARGETARCH

# busybox wget/tar/sha256sum only: no apk packages to pin (fixes HLD_DF_021),
# TLS verified (HLD_DF_008), no pipe to shell (HLD_DF_018), checksum checked.
RUN set -eu; \
    test -n "${VERSION}"; \
    v="${VERSION#v}"; \
    arch="${TARGETARCH:-amd64}"; \
    f="crusher_${v}_linux_${arch}.tar.gz"; \
    base="https://github.com/JamesWoolfenden/crusher/releases/download/v${v}"; \
    wget -q -O "/tmp/${f}" "${base}/${f}"; \
    wget -q -O /tmp/checksums.txt "${base}/checksums.txt"; \
    (cd /tmp && grep " ${f}\$" checksums.txt | sha256sum -c -); \
    tar -xzf "/tmp/${f}" -C /usr/bin crusher; \
    rm -f "/tmp/${f}" /tmp/checksums.txt

COPY entrypoint.sh /entrypoint.sh
RUN adduser -D -u 10001 crusher
USER crusher

# Code file to execute when the docker container starts up (`entrypoint.sh`)
ENTRYPOINT ["/entrypoint.sh"]
