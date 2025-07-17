### OnionArchitecture

```
.
├── .circleci
│   └── config.yml
├── cmd/
│   └── app/             # entrypoint: main.go
├── external/
│   ├── adapter/
│   └── transport/
│   │   └── webhook/
├── internal/
│   ├── entity/
│   │   └── user/        # entity, repository interface, domain logic
│   ├── application/
│   │   └── user/        # usecase (app logic)
│   ├── infrastructure/
│   │   ├── firestore/
│   │   ├── kms/
│   │   ├── pubsub/
│   │   ├── persistence/ # implement repository (e.g., Postgres, Redis)
│   │   └── transport/
│   │       ├── http/    # HTTP handlers
│   │       └── grpc/    # gRPC handlers
│   └── config/          # struct mapping config file/env
├── api/                 # protobufs or openapi specs
├── pkg/                 # utils và các component tái sử dụng (hash, jwt...)
└── go.mod
```