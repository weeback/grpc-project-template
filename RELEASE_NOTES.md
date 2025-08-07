# Release Notes

## [v0.1.1] - 2025-08-07

### 🐛 Bug Fixes

- **import paths**: Fix import paths to use the correct GitHub domain ([e2ecad1](https://github.com/weeback/grpc-project-template/commit/e2ecad1ca913261a36e7a68555a28e4182f1452f))
  - Update all import statements to reference the proper GitHub repository path
  - Ensure consistent module naming across the project

### ⚡ Improvements

- Merged multiple import path fixes for better project stability

---

## [v0.1.0] - 2025-08-07

### ✨ New Features

- **websocket**: Implement WebSocket server with client management and messaging ([e22d439](https://github.com/weeback/grpc-project-template/commit/e22d439f78f60f5d8976db70ba24c2f4fb6f41ae))
  - Add WebSocket server functionality
  - Implement client connection management
  - Add real-time messaging capabilities
  - Enhance the project with bidirectional communication support

### � Import Fixes

- **import paths**: Update import paths to use the correct GitHub domain ([1f8d135](https://github.com/weeback/grpc-project-template/commit/1f8d135ce81fb084c33fbb1bc38e832be448c77e))

---

## [v0.0.1] - 2025-07-19

### 🔧 Refactoring

- **HelloService**: Update JWT validation and response handling ([b5fe3e5](https://github.com/weeback/grpc-project-template/commit/b5fe3e5f5d6820b7ca76254b9b835f46763cc78f))
  - Improve JWT token validation logic
  - Enhance response handling mechanisms
  - Better error handling for authentication flows
  - Code structure improvements for maintainability

### 📝 Documentation

- Initial project setup and documentation

---

## [v0.0.0] - 2025-07-17

### 🎉 Initial Release

- **project**: First commit ([c434cf6](https://github.com/weeback/grpc-project-template/commit/c434cf6307300cbc05fe0c93b853f0cc7e066895))
  - Initial project structure setup
  - Basic gRPC service template
  - Core infrastructure components
  - Authentication and authorization framework
  - MongoDB integration
  - Firebase Admin SDK integration
  - Cloudflare CAPTCHA service
  - JWT token management
  - HTTP and gRPC transport layers
  - Basic Hello service implementation

### 🏗️ Architecture Components

- **Internal Structure**:
  - Application layer with controllers and validation
  - Entity layer for business logic
  - Infrastructure layer for external services
  - Transport layer for HTTP/gRPC communication

- **External Integrations**:
  - MongoDB database support
  - Firebase services
  - Cloudflare CAPTCHA
  - JWT authentication
  - gRPC and HTTP servers

### 📦 Package Structure

- Protocol buffer definitions
- Mock implementations for testing
- Comprehensive Makefile for build automation
- Docker support with CloudRun deployment
- CI/CD pipeline configuration

---

## Migration Guide

### From v0.0.1 to v0.1.0

- **WebSocket Support**: New WebSocket functionality is now available. Update your client applications to leverage real-time communication features.
- **Import Paths**: Ensure all import paths are updated to use the correct GitHub domain.

### From v0.1.0 to v0.1.1

- **No Breaking Changes**: This is a patch release with only bug fixes.
- **Import Path Fixes**: All import paths have been corrected - no action required from users.

## Contributors

- [@thinh-wee](https://github.com/thinh-wee) - Project maintainer and primary contributor

## Links

- [Repository](https://github.com/weeback/grpc-project-template)
- [Issues](https://github.com/weeback/grpc-project-template/issues)
- [Pull Requests](https://github.com/weeback/grpc-project-template/pulls)
