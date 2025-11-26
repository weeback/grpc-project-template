# gRPC Project Template - Architecture Diagram

Đây là sơ đồ kiến trúc của dự án gRPC Project Template, mô tả mối liên kết giữa các components và file trong hệ thống.

## Kiến trúc tổng quan

Dự án được xây dựng theo mô hình **Onion Architecture** (Clean Architecture) với các layer rõ ràng:

### 1. Build & Deploy Layer 🔨

- **Makefile.CloudRun-HelloService**: Quản lý deployment lên Google Cloud Run
- **HelloService.Dockerfile**: Multi-stage Docker build
- **Makefile**: Build tools và gRPC code generation

### 2. Entry Point 🚀

- **cmd/HelloService/main.go**: Application entry point, khởi tạo và kết nối các components

### 3. Application Layer 💼 (Business Logic)

- **controller.go**: Implement business logic chính
- **validation.go**: Validate input data
- **utils.go**: Helper functions

### 4. Domain Layer 🏛️ (Entities & Interfaces)

- **repository.go**: Define interfaces
- **entity/**: Domain entities và data structures

### 5. Infrastructure Layer 🏗️

- **Transport**: gRPC và HTTP handlers
- **Database**: MongoDB connection và operations
- **External Services**: Cloudflare, Firebase integrations

### 6. Generated Code 🤖

- **pb/**: Generated protobuf code từ proto definitions

### 7. Shared Packages 📦

- **pkg/**: Reusable utilities (network, JWT, MongoDB)

## Flow hoạt động

1. **Build Process**:

   ```text
   Makefile.CloudRun-HelloService → HelloService.Dockerfile → Makefile → Go Build
   ```

2. **Request Flow**:

   ```text
   Client → Transport Layer → Application Layer → Domain Layer → Infrastructure
   ```

3. **Dependencies**:
   - Main.go khởi tạo và inject dependencies
   - Application layer sử dụng repository interfaces
   - Infrastructure layer implement concrete services

## Mermaid Diagram

> 📊 **Architecture Diagram**: [architecture-diagram.mmd](./architecture-diagram.mmd)

Sơ đồ chi tiết được lưu trong file riêng để dễ quản lý và chỉnh sửa.

**Cách xem diagram:**

- Mở file `architecture-diagram.mmd` trong VS Code với Mermaid extension
- Hoặc sử dụng preview trong VS Code: `Ctrl+Shift+V` (Windows/Linux) hoặc `Cmd+Shift+V` (Mac)
- Online viewer: [Mermaid Live Editor](https://mermaid.live)

## Các lệnh hữu ích

### Build & Development

```bash
# Build ứng dụng
make build

# Run ứng dụng locally  
make run

# Generate gRPC code
make grpc-generate

# Run tests
make test
```

### Docker & Cloud Run

```bash
# Build và deploy lên Cloud Run
make -f Makefile.CloudRun-HelloService deploy

# Chỉ build Docker image
make -f Makefile.CloudRun-HelloService build-docker

# Xem thông tin deployment
make -f Makefile.CloudRun-HelloService info
```

## Nguyên tắc thiết kế

1. **Dependency Inversion**: Application layer không phụ thuộc vào Infrastructure
2. **Separation of Concerns**: Mỗi layer có trách nhiệm riêng biệt
3. **Clean Dependencies**: Dependencies chỉ flow từ ngoài vào trong
4. **Testability**: Mock interfaces để test dễ dàng
5. **Configuration**: Environment-based configuration management

## Notes

- File này có thể được cập nhật khi có thay đổi trong source code
- Sử dụng Mermaid để visualization và có thể chỉnh sửa trực tiếp
- Diagram hỗ trợ trong VS Code với Mermaid extension
