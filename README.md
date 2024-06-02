# Kupilot

Kupilot is a command-line tool that provides Kubernetes support using OpenAI's GPT models. It's designed to offer instant help with Kubernetes questions and operations through an interactive terminal interface.

![Kupilot Demo](docs/demo.gif)

## Prerequisites

OpenAI API Key: You can obtain this key by creating an account on the [OpenAI website](https://platform.openai.com/api-keys/). After signing up, navigate to the API section to generate your key.

## Getting Started

### Run locally

Warning. Kupilot may have write access to your Kubernetes cluster. Validate and execute any command with caution.

```bash
export OPENAI_API_KEY=<your-openai-api-key>
docker run --rm -it -e OPENAI_API_KEY=$OPENAI_API_KEY -e KUBECONFIG=/.kube/config -v ~/.kube:/.kube ghcr.io/r2k1/kupilot:latest kupilot
```

### Run in Kubernetes

```bash
export OPENAI_API_KEY=<your-openai-api-key>
curl https://raw.githubusercontent.com/r2k1/kupilot/main/kube.yaml | envsubst '$OPENAI_API_KEY' | kubectl apply -f -
kubectl -n kupilot exec -it $(kubectl get pods -n kupilot -l app=kupilot -o jsonpath='{.items[0].metadata.name}') -- kupilot
```
