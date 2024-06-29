# Kupilot

Kupilot is a command-line tool that provides Kubernetes support using OpenAI's GPT models. It's designed to offer instant help with Kubernetes questions and operations through an interactive terminal interface.

![Kupilot Demo](docs/demo.gif)

## Getting Started

### Prior Requirements

OpenAI API Key: You can obtain this key by creating an account on the [OpenAI website](https://platform.openai.com/api-keys/). After signing up, navigate to the API section to generate your key.

Please keep in mind, Kupilot may have write-access to your Kubernetes       
cluster. It is advisable to validate and execute any command prudently

### Running Kupilot Locally with Docker

```bash
export OPENAI_API_KEY=<your-openai-api-key>
docker run --rm -it -e OPENAI_API_KEY=$OPENAI_API_KEY -e KUBECONFIG=/.kube/config -v ~/.kube:/.kube ghcr.io/r2k1/kupilot:latest kupilot
```

### Running Kupilot Locally with Go
```bash
export OPENAI_API_KEY=<your-openai-api-key>
go run github.com/r2k1/kupilot
```

Alternatively, you may choose to install the tool:
```bash
export OPENAI_API_KEY=<your-openai-api-key>
go install github.com/r2k1/kupilot
kupilot
```

### Running Kupilot Locally with Binary

Download binary from the [latest release](https://github.com/r2k1/kupilot/releases)

```bash
export OPENAI_API_KEY=<your-openai-api-key>
kupilot
```

### Running Kupilot within Kubernetes

This method permits the use of RBAC and narrows down the permission scope to read-only within the cluster:  

```bash
export OPENAI_API_KEY=<your-openai-api-key>
curl https://raw.githubusercontent.com/r2k1/kupilot/main/kube.yaml | envsubst '$OPENAI_API_KEY' | kubectl apply -f -
kubectl -n kupilot exec -it $(kubectl get pods -n kupilot -l app=kupilot -o jsonpath='{.items[0].metadata.name}') -- kupilot
```
