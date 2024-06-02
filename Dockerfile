FROM ubuntu:latest
# Set up build arguments to handle architecture differences in binaries
ARG TARGETARCH

RUN apt-get update && apt-get install -y curl jq yq
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/${TARGETARCH}/kubectl" && \
    chmod +x kubectl && \
    mv kubectl /usr/local/bin/
COPY kupilot /usr/local/bin/kupilot
RUN chmod +x /usr/local/bin/kupilot

RUN groupadd -g 1001 nonroot && useradd -u 1001 -g nonroot nonroot
USER 1001

CMD echo "Container started successfully! Run 'kupilot'" && tail -f /dev/null
