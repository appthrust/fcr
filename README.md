# FCR - Functional Controller Runtime

[![Go Version](https://img.shields.io/github/go-mod/go-version/appthrust/fcr)](https://golang.org/doc/devel/release)
[![Go Report Card](https://goreportcard.com/badge/github.com/appthrust/fcr)](https://goreportcard.com/report/github.com/appthrust/fcr)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![GoDoc](https://pkg.go.dev/badge/github.com/appthrust/fcr)](https://pkg.go.dev/github.com/appthrust/fcr)

A comprehensive functional programming toolkit for Kubernetes development, providing functional wrappers for controller-runtime packages and additional utilities for composable, type-safe, and error-safe interactions with Kubernetes resources and controller patterns.

## Overview

FCR brings functional programming paradigms to Kubernetes development by leveraging monadic patterns (Either, IO, Reader) from [IBM/fp-go](https://github.com/IBM/fp-go). It provides both functional wrappers for controller-runtime packages and additional utilities designed specifically for functional programming patterns, enabling:

- **Composable Operations**: Chain complex Kubernetes operations and controller patterns together cleanly
- **Type Safety**: Leverage Go's type system with functional programming patterns
- **Error Safety**: Handle errors explicitly through Either types
- **Testability**: Pure functional operations are easier to test and reason about
- **Immutability**: Functional approach reduces side effects and improves reliability
- **Complete Ecosystem**: Functional wrappers for all controller-runtime packages plus additional utilities

## Features

- 🎯 **Functional Client Operations**: Complete CRUD operations, status updates, and resource management with functional patterns
- 🏗️ **Controller & Webhook Patterns**: Type-safe controllers, managers, reconcilers, event handlers, and webhooks (Comming soon)
- 🔗 **Composable Utilities**: Advanced flows, transformations, query builders, retry patterns, and validation pipelines (Comming soon)
- 🛡️ **Type Safety & Testing**: Generic functions with compile-time checking and comprehensive test coverage

## Package Structure

FCR provides both functional wrappers (mirroring controller-runtime packages with an `f` prefix) and additional utilities designed for functional programming patterns. This structure avoids naming conflicts, allowing you to import both libraries simultaneously:

```go
import (
    // controller-runtime packages
    "sigs.k8s.io/controller-runtime/pkg/client"
    "sigs.k8s.io/controller-runtime/pkg/controller"
    "sigs.k8s.io/controller-runtime/pkg/manager"

    // FCR functional wrappers
    "github.com/appthrust/fcr/pkg/fclient"
    "github.com/appthrust/fcr/pkg/fcontroller"  // Coming soon
    "github.com/appthrust/fcr/pkg/fmanager"     // Coming soon

    // FCR functional utilities
    "github.com/appthrust/fcr/pkg/flow"         // Coming soon
    "github.com/appthrust/fcr/pkg/transform"    // Coming soon
    "github.com/appthrust/fcr/pkg/validate"     // Coming soon
)
```

### Available Packages

| controller-runtime | FCR Wrapper       | Status         | Description                     |
| ------------------ | ----------------- | -------------- | ------------------------------- |
| `pkg/client`       | `pkg/fclient`     | ✅ Ready       | Functional client operations    |
| `pkg/controller`   | `pkg/fcontroller` | 🚧 Coming Soon | Functional controller patterns  |
| `pkg/manager`      | `pkg/fmanager`    | 🚧 Coming Soon | Functional manager utilities    |
| `pkg/builder`      | `pkg/fbuilder`    | 🚧 Coming Soon | Functional controller builders  |
| `pkg/cache`        | `pkg/fcache`      | 🚧 Coming Soon | Functional caching operations   |
| `pkg/handler`      | `pkg/fhandler`    | 🚧 Coming Soon | Functional event handlers       |
| `pkg/predicate`    | `pkg/fpredicate`  | 🚧 Coming Soon | Functional predicates           |
| `pkg/webhook`      | `pkg/fwebhook`    | 🚧 Coming Soon | Functional webhook patterns     |
| `pkg/reconcile`    | `pkg/freconcile`  | 🚧 Coming Soon | Functional reconciler utilities |

## Installation

```bash
go get github.com/appthrust/fcr
```

## API Documentation

For complete API documentation, type definitions, and usage examples, visit the [GoDoc reference](https://pkg.go.dev/github.com/appthrust/fcr).

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines on how to contribute to FCR.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built on top of [controller-runtime](https://github.com/kubernetes-sigs/controller-runtime)
- Uses [IBM/fp-go](https://github.com/IBM/fp-go) for functional programming primitives
- Inspired by functional programming patterns in other languages
