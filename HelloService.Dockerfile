FROM amd64/golang:1.24-alpine3.21 AS build
# Install necessary packages
RUN apk add --no-cache \
    git \
    bash \
    curl \
    make \
    gcc \
    musl-dev \
    linux-headers
# Set the Go environment variables
ENV GO111MODULE=on \
    GOPROXY=https://proxy.golang.org,direct \
    GOSUMDB=sum.golang.org \
    CGO_ENABLED=0
# Set the working directory
WORKDIR /app
# Copy the Go module files
COPY go.mod go.sum ./
# Download the Go module dependencies
RUN go mod download
# Copy the source code
COPY . .
# Build the Go application
RUN MAIN_GO_FILE="/app/cmd/HelloService/*.go" BINARY="/app/bin/HelloService" \
    make --makefile=Makefile build

# Sử dụng .git nếu cần (ví dụ: lấy thông tin commit)
RUN git rev-parse HEAD > /app/commit_hash.txt
# Xóa thư mục .git để giảm kích thước image
RUN rm -rf .git

# --- Production Stage ---
# Use the latest alpine image
FROM alpine:3.21

LABEL maintainer="khoa.ngo@bankaool.com"
LABEL description="A Docker image for running the latest version of the Go application"
LABEL version="1.0"
LABEL org.opencontainers.image.source="gcr.io/service-test-40fef/grpc-project-template-hello-srv:latest"
# Install necessary packages
RUN apk add --no-cache \
    bash \
    curl \
    linux-headers
# Set the working directory
WORKDIR /app
# Copy ...
COPY --from=build /app/proto/google/ /app/proto/google/
COPY --from=build /app/proto/common/ /app/proto/common/
COPY --from=build /app/proto/hello/ /app/proto/hello/

# Copy the built Go application from the build stage
COPY --from=build /app/bin/HelloService /app/bin/HelloService

# Copy the configuration file
# COPY --from=build /app/config.yaml /app/config.yaml
# Copy the entrypoint script
# COPY --from=build /app/docker-entrypoint.sh /app/docker-entrypoint.sh
# Copy the entrypoint script
# COPY --from=build /app/entrypoint.sh /app/entrypoint.sh

# Expose the application port
EXPOSE 8080
# Set the entrypoint script as executable
# RUN chmod +x /app/docker-entrypoint.sh
# Set the entrypoint command
# ENTRYPOINT ["/app/docker-entrypoint.sh"]
# Set the default command to run the application
CMD [ "/app/bin/HelloService" ]

