##############################
FROM golang:1.25-alpine AS build

ARG VERSION "devel"
ARG GIT_COMMIT ""

WORKDIR /src

RUN --mount=type=bind,source=.,target=.  \
  --mount=type=cache,target=/root/.cache/go-build \
  --mount=type=cache,target=/go/pkg \
  CGO_ENABLED=0 GOOS=linux go build -ldflags="-X 'main.version=$VERSION' -X 'main.gitCommit=$GIT_COMMIT'" -o /tmp/png2asm main.go

##############################
FROM scratch

ARG VERSION

LABEL org.opencontainers.image.title="gpl2asm" \
  org.opencontainers.image.vendor="laghoule" \
  org.opencontainers.image.licenses="GPLv3" \
  org.opencontainers.image.version="${VERSION}" \
  # FIXME: description
  org.opencontainers.image.description="Little tool to convert 256-color paletted PNG images to assembly include files." \
  org.opencontainers.image.url="https://github.com/laghoule/pgl2asm/README.md" \
  org.opencontainers.image.source="https://github.com/laghoule/gpl2asm" \
  org.opencontainers.image.documentation="https://github.com/laghoule/gpl2asm/README.md"

USER 65534

COPY --link --from=build /tmp/gpl2asm /bin/gpl2asm

ENTRYPOINT ["/bin/gpl2asm"]
