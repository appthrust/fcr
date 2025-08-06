# Contributing to FCR

We welcome contributions to FCR! This document provides guidelines and information for contributors.

FCR provides both functional wrappers for controller-runtime packages and additional functional utilities designed specifically for Kubernetes development with functional programming patterns.

## Getting Started

### Prerequisites

- [devbox](https://github.com/jetify-com/devbox) for development environment

### Development Setup

```bash
# Clone the repository
git clone https://github.com/appthrust/fcr.git
cd fcr

# Setup development environment
devbox shell

# Install dependencies
go mod download

# Generate code
task generate

# Run tests (requires Kubernetes cluster)
task test
```

## Development Workflow

### Available Tasks

```bash
# Generate code and manifests
task generate

# Run linting
task lint

# Run tests
task test

# Create local Kubernetes cluster
task cluster:create

# Delete local Kubernetes cluster
task cluster:delete

# Run CI checks
task ci
```

### Code Generation

This project uses controller-gen for generating:

- DeepCopy methods
- Kubernetes Custom Resource Definitions (CRDs)

```bash
# Generate all code
task generate

# Generate only DeepCopy methods
task generate:deepcopy

# Generate only CRD manifests
task generate:manifests
```

### Testing

Tests require a running Kubernetes cluster. The test suite will:

1. Create a Kind cluster automatically
2. Apply necessary CRDs
3. Run comprehensive integration tests
4. Clean up resources

```bash
# Run all tests
task test

# Run tests with verbose output
go test -v -race ./...
```

### Testing Guidelines

- All new functionality must include comprehensive tests
- Tests should use real Kubernetes clusters when possible
- Use the existing test patterns and helpers
- Ensure tests clean up resources properly
- Test both success and error paths

## Contributing Guidelines

### How to Contribute

1. **Fork the repository** on GitHub
2. **Create a feature branch** from `main`:
   ```bash
   git checkout -b feature/your-feature-name
   ```
3. **Make your changes** following our coding standards
4. **Add tests** for new functionality
5. **Run the full test suite**:
   ```bash
   task ci
   ```
6. **Commit your changes** with clear commit messages
7. **Push to your fork** and **submit a pull request**

### Pull Request Process

1. Ensure your PR description clearly describes the changes
2. Link any relevant issues
3. Ensure all CI checks pass
4. Update documentation if needed
5. Request review from maintainers

### Code Style

- **Follow standard Go conventions** (gofmt, golint, etc.)
- **Use functional programming patterns** where appropriate
- **Prefer composition over inheritance**
- **Write clear, self-documenting code**
- **Include comprehensive documentation** for public APIs
- **Use meaningful variable and function names**

### Functional Programming Guidelines

- Leverage monadic patterns (Either, IO, Reader) from IBM/fp-go
- Prefer pure functions with explicit error handling
- Use composition to build complex operations from simple ones
- Avoid side effects where possible
- Make dependencies explicit through function parameters

### Package Development

#### Functional Wrappers

When creating new functional wrapper packages:

1. **Mirror controller-runtime structure** with `f` prefix
2. **Maintain consistent API patterns** across packages
3. **Use generic types** where appropriate for type safety
4. **Provide parameter constructors** (e.g., `ToGetParams`)
5. **Include comprehensive examples** in package documentation

#### Functional Utilities

When creating new utility packages:

1. **Focus on specific functional programming patterns** (composition, transformation, validation, etc.)
2. **Design for composability** with other FCR packages
3. **Provide monadic interfaces** consistent with FCR patterns
4. **Include practical examples** showing integration with wrappers
5. **Ensure utilities are reusable** across different Kubernetes use cases

### Documentation

- **Update README.md** for user-facing changes
- **Add GoDoc comments** for all public functions and types
- **Include usage examples** in package documentation
- **Update package status** in the main README if needed

### Commit Message Format

Use clear, descriptive commit messages:

```
feat(fclient): add support for server-side apply

- Add ServerSideApply parameter to patch operations
- Update tests to cover new functionality
- Add documentation examples

Fixes #123
```

#### Commit Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

## Code Architecture

### Package Structure

FCR provides both functional wrappers and utilities:

```
pkg/
# Functional Wrappers (mirror controller-runtime with 'f' prefix)
├── fclient/      # ✅ Functional client operations
├── fcontroller/  # 🚧 Functional controller patterns (planned)
├── fmanager/     # 🚧 Functional manager utilities (planned)
├── fbuilder/     # 🚧 Functional controller builders (planned)
├── fcache/       # 🚧 Functional caching operations (planned)
├── fhandler/     # 🚧 Functional event handlers (planned)
├── fwebhook/     # 🚧 Functional webhooks (planned)
├── fpredicate/   # 🚧 Functional predicates (planned)
└── freconcile/   # 🚧 Functional reconciler utilities (planned)
```

### Design Principles

1. **Functional First**: Use functional programming patterns consistently
2. **Type Safety**: Leverage Go's type system for compile-time safety
3. **Composability**: Build complex operations from simple, reusable parts
4. **Error Safety**: Handle errors explicitly through Either types
5. **Performance**: Maintain performance characteristics of controller-runtime

## Testing Strategy

### Test Categories

1. **Unit Tests**: Test individual functions and components
2. **Integration Tests**: Test interactions with real Kubernetes clusters
3. **End-to-End Tests**: Test complete workflows

### Test Structure

```go
var _ = Describe("FunctionName", func() {
    var env fclient.Env

    BeforeEach(func() {
        // Setup test environment
    })

    AfterEach(func() {
        // Cleanup resources
    })

    It("should handle success case", func() {
        // Test implementation
    })

    It("should handle error case", func() {
        // Test error paths
    })
})
```

## Release Process

1. Update version numbers
2. Update CHANGELOG.md
3. Run full test suite
4. Create release PR
5. Tag release after merge
6. Publish to pkg.go.dev

## Getting Help

- **GitHub Issues**: For bug reports and feature requests
- **Discussions**: For questions and general discussion
- **Code Review**: Ask questions in PR comments

## Recognition

Contributors are recognized in:

- Git commit history
- Release notes for significant contributions
- README acknowledgments for major features

Thank you for contributing to FCR! 🚀
