# Start from the official Golang image to build the binary
FROM golang:1.22 as builder

ARG TARGETARCH

# Set the Current Working Directory inside the container
WORKDIR /build

# Cache dependency fetching
COPY go.mod go.sum ./
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -o kupilot .

FROM ubuntu:latest
# Set up build arguments to handle architecture differences in binaries
ARG TARGETARCH

RUN apt-get update && apt-get install -y curl jq yq
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/${TARGETARCH}/kubectl" && \
    chmod +x kubectl && \
    mv kubectl /usr/local/bin/
COPY --from=builder /build/kupilot /usr/local/bin/kupilot
RUN chmod +x /usr/local/bin/kupilot

RUN groupadd -g 1001 nonroot && useradd -u 1001 -g nonroot nonroot
USER 1001

CMD echo "Container started successfully! Run 'kupilot'" && tail -f /dev/null
