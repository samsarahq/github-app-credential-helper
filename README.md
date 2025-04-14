# github-app-credential-helper

[![release](https://github.com/samsarahq/github-app-credential-helper/actions/workflows/release.yml/badge.svg)](https://github.com/samsarahq/github-app-credential-helper/actions/workflows/release.yml)
![GitHub License](https://img.shields.io/github/license/samsarahq/github-app-credential-helper)
![GitHub Tag](https://img.shields.io/github/v/tag/samsarahq/github-app-credential-helper)
[![Go Report Card](https://goreportcard.com/badge/github.com/samsarahq/github-app-credential-helper)](https://goreportcard.com/report/github.com/samsarahq/github-app-credential-helper)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/samsarahq/github-app-credential-helper)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/samsarahq/github-app-credential-helper?tab=doc)

This is a library for building git credential helpers for github apps.

## Usage
Implement a `SecretProvider` for whatever secret store you are using. Pass that into the `NewAuthenticator` constructor
and run the `Authenticate` function. The result is a string that should be written to stdout. Nothing else should be
written to stdout since stdin/stdout are the interface pipes for git credential commands.

To see an example implementation, we have an [AWS Secrets Manager implementation](https://github.com/samsarahq/git-credential-github-app-sm)
that pulls the application ID, installation ID, and private key from AWS Secrets Manager.