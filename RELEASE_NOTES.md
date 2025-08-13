# Release Notes

## [v0.2.1] - 2025-08-13

### 🔧 Http Logger interceptor Improvements

- **http interceptor**: Merge feature branches for logger bugfixes ([313e4cf](https://github.com/weeback/grpc-project-template/commit/313e4cf45382196cb9d5d9559fe44e73c7b97cf0))
  - Add logger interceptor to CORS middleware for enhanced request logging

---

## [v0.2.0] - 2025-08-13

### ✨ New Features

- **logging**: Implement structured logging for HelloService and add logger context management ([de74656](https://github.com/weeback/grpc-project-template/commit/de74656d5c3d59c0d3071f86cf4fc26c16eebee7))
  - Introduce `pkg/logger` package for centralized logging utilities
  - Add context-based logger management for request tracing
  - Improve log output format for better debugging and monitoring

### 🐛 Bug Fixes & Improvements

- Merge feature branches for logger and websocket bugfixes ([6c45694](https://github.com/weeback/grpc-project-template/commit/6c4569426b86f9976645850e5ea32e43a283a390))

---

## [v0.1.3] - 2025-08-13

### 🐛 WebSocket Bug Fixes

- **websocket**: Merge websocket-hijack bugfixes ([8c1ca77](https://github.com/weeback/grpc-project-template/commit/8c1ca7790bb0ea1629b119481ea057a25357c977))
  - Fix empty message handling in writePump
  - Optimize WebSocket write deadline management
  - Improve message type detection and handling

### 🔧 WebSocket Improvements

- **websocket**: Enhance message handling with binary support and update SendTo method ([0a774c1](https://github.com/weeback/grpc-project-template/commit/0a774c171bfde6cbe0d27252843051d25b415bf8))
- **client**: Skip sending empty messages in writePump ([9c4f719](https://github.com/weeback/grpc-project-template/commit/9c4f719b7f69386ca012e2b90f422ce7695327f4))

---

## [v0.1.2] - 2025-08-07

### 🐛 Hijack & Connection Fixes

- **websocket**: Merge websocket-hijack bugfixes ([aeceec0](https://github.com/weeback/grpc-project-template/commit/aeceec08d3538e6a683640aef90d7094b73f76d2))
  - Add Hijack method for WebSocket support in ResponseWriter
  - Enable HTTP connection hijacking for WebSocket upgrades
  - Improve WebSocket protocol compatibility

### 🔧 Client Management Improvements

- **websocket**: Add ChangeID method to Client for dynamic client ID management ([01a466c](https://github.com/weeback/grpc-project-template/commit/01a466cd17e77ce454363467ad4f2e9508fb26db))
- **writer**: Add Hijack method for WebSocket support ([a29ae6a](https://github.com/weeback/grpc-project-template/commit/a29ae6ae8b41f6e5981ca4e4e6e193769d5cf135))

---

## [v0.1.1] - 2025-08-07

### 🐛 Bug Fixes

- **import paths**: Fix import paths to use the correct GitHub domain ([e2ecad1](https://github.com/weeback/grpc-project-template/commit/e2ecad1ca913261a36e7a68555a28e4182f1452f))
  - Update all import statements to reference the proper GitHub repository path
  - Ensure consistent module naming across the project

### ⚡ Improvements

- Merged multiple import path fixes for better project stability

---

## [v0.1.0] - 2025-08-07

### ✨ WebSocket Features

- **websocket**: Implement WebSocket server with client management and messaging ([e22d439](https://github.com/weeback/grpc-project-template/commit/e22d439f78f60f5d8976db70ba24c2f4fb6f41ae))
  - Add WebSocket server functionality
  - Implement client connection management
  - Add real-time messaging capabilities
  - Enhance the project with bidirectional communication support

### 🔧 Import Fixes

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

### From v0.1.3 to v0.2.0

- **Structured Logging**: New logging system is now available. The changes are backward compatible.
- **No Breaking Changes**: This is a minor release with new features. Safe to upgrade.

### From v0.1.2 to v0.1.3

- **WebSocket Improvements**: Enhanced WebSocket functionality with better message handling.
- **No Breaking Changes**: This is a patch release with bug fixes and improvements.

### From v0.1.1 to v0.1.2

- **WebSocket Enhancements**: New WebSocket features including Hijack method and dynamic client ID management.
- **No Breaking Changes**: This is a patch release with new features and bug fixes.

### From v0.1.0 to v0.1.1

- **No Breaking Changes**: This is a patch release with only bug fixes.
- **Import Path Fixes**: All import paths have been corrected - no action required from users.

### From v0.0.1 to v0.1.0

- **WebSocket Support**: New WebSocket functionality is now available. Update your client applications to leverage real-time communication features.
- **Import Paths**: Ensure all import paths are updated to use the correct GitHub domain.

## Contributors

- [@thinh-wee](https://github.com/thinh-wee) - Project maintainer and primary contributor

## Links

- [Repository](https://github.com/weeback/grpc-project-template)
- [Issues](https://github.com/weeback/grpc-project-template/issues)
- [Pull Requests](https://github.com/weeback/grpc-project-template/pulls)
