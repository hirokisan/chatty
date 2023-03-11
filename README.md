# chatty

chatty is a cli application that will be your conversation partner in your spare time.

## Preparation

Create OpenAI account and get API Key.

ref: https://platform.openai.com/account/api-keys

And then, set the key in the environment variable.

```console
$ export OPEN_AI_KEY={key}
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
