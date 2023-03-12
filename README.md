[![Go Report Card](https://goreportcard.com/badge/github.com/hirokisan/chatty)](https://goreportcard.com/report/github.com/hirokisan/chatty)
[![test](https://github.com/hirokisan/chatty/actions/workflows/test.yml/badge.svg)](https://github.com/hirokisan/chatty/actions/workflows/test.yml)

# chatty

chatty is a cli application that will be your conversation partner in your spare time.

## Preparation

Create OpenAI account and get API Key.

ref: https://platform.openai.com/account/api-keys

And then, set the key in the environment variable.

```console
$ export OPEN_AI_KEY={key}
```

If you wish to record past exchanges and conversations, please specify a file path for the record in an environment variable.

The file will be saved in json format

```console
$ export CHATTY_MESSAGES_PATH={filepath}

# e.g. export CHATTY_MESSAGES_PATH=~/chatty-history.json
```

## Installation

```console
$ go install github.com/hirokisan/chatty@latest
```

## Usage

```console
$ chatty my name is chatty
Hello, Chatty! How can I assist you today?
```
