// Package contracts contains the hexagonal architecture ports and domain models.
//
// # Port Types
//
// Inbound Ports (Driving): Interfaces that drive the application (use cases, commands, queries)
// Outbound Ports (Driven): Interfaces that the application calls to interact with external systems
//
// # Architecture
//
//	┌─────────────────────────────────────────────────────────────┐
//	│                      Driving Adapters                        │
//	│  (REST, gRPC, CLI, Events)                                   │
//	└─────────────────────────┬───────────────────────────────────┘
//	                          │
//	                          ▼
//	┌─────────────────────────────────────────────────────────────┐
//	│                    Inbound Ports                            │
//	│  (UseCase, CommandHandler, QueryHandler)                    │
//	└─────────────────────────┬───────────────────────────────────┘
//	                          │
//	                          ▼
//	┌─────────────────────────────────────────────────────────────┐
//	│                     Domain Core                              │
//	│  (Entities, Value Objects, Domain Services, Events)          │
//	└─────────────────────────┬───────────────────────────────────┘
//	                          │
//	                          ▼
//	┌─────────────────────────────────────────────────────────────┐
//	│                   Outbound Ports                            │
//	│  (Repository, EventBus, Cache, ExternalService)            │
//	└─────────────────────────┬───────────────────────────────────┘
//	                          │
//	                          ▼
//	┌─────────────────────────────────────────────────────────────┐
//	│                    Driven Adapters                           │
//	│  (Postgres, Redis, Kafka, HTTP Client)                       │
//	└─────────────────────────────────────────────────────────────┘
package contracts
