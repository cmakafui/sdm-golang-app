# syntax=docker/dockerfile:1

# Create a stage for building the application.
ARG GO_VERSION=1.22.4
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build
WORKDIR /src

# Copy go.mod file
COPY go.mod ./

# Download dependencies as a separate step to take advantage of Docker's caching.
# Leverage a cache mount to /go/pkg/mod/ to speed up subsequent builds.
RUN --mount=type=cache,target=/go/pkg/mod/ \
    go mod download -x

# Copy the source code and templates into the container.
COPY . .

# Build the application.
# Leverage a cache mount to /go/pkg/mod/ to speed up subsequent builds.
RUN --mount=type=cache,target=/go/pkg/mod/ \
    CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/server ./cmd/main.go

# Create a new stage for running the application that contains the minimal
# runtime dependencies for the application.
FROM alpine:latest AS final

# Install any runtime dependencies that are needed to run your application.
RUN apk --update add \
        ca-certificates \
        tzdata \
        && \
        update-ca-certificates

# Create a non-privileged user that the app will run under.
ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser
USER appuser

# Copy the executable and static files from the "build" stage.
COPY --from=build /bin/server /bin/
COPY --from=build /src/web /web

# Expose the port that the application listens on.
EXPOSE 5080

# What the container should run when it is started.
ENTRYPOINT [ "/bin/server" ]
