# Storage layer study guide — for a beginner Go intern

Goal: give you a clear, ordered set of files and simple exercises to read so you can understand how the storage (MongoDB) layer works in this microservice.

Prerequisites

- Basic Go syntax (variables, functions, structs, methods, interfaces)
- Comfort opening files in VS Code and reading short code snippets

How to use this guide

- Read files in the order below. For each file, open it and follow the "What to look for" hints. Try the small exercises after the reading section.
- Keep a terminal open and use `rg`/`grep` or VS Code search to find symbols mentioned here.

Reading order (step-by-step)

1. `config/mongo.go`

- Why first: shows how the Mongo `Client` and `Database` are created and configured. If you understand this, you'll know where the `mongo.Client` passed to DALs originates.
- What to look for: `NewMongoConnection`, how `client.Database(...)` is used, connection pool settings, `Close()`.

2. `internal/storage/peristance/init.go` (the DAL implementation)

- Why: this file contains `MongoDal` (the generic DAL wrapper around `mongo.Collection`) — core to all DB calls.
- What to look for:
  - `NewMongoDal(client, dbName, collection)` — it returns a struct holding a `mongo.Collection`.
  - `FindAll`, `FindOne`, `InsertOne`, `UpdateOne`, `DeleteOne` — read their implementations and note use of `context.WithTimeout`, `collection.Find`, `cursor.All`, `FindOneAndUpdate`, etc.
  - Soft-delete pattern (`is_deleted`) and cursor-based pagination logic.
- Why this matters: all storage code uses this DAL to avoid duplicating Mongo driver code.

3. `internal/constants/models` and `internal/constants/entities` (type definitions)

- Why: compare DB `models` and app `entities`. Models usually have `bson` tags and `ObjectID` types, while entities are app-facing structs.
- What to look for: `UserAccess`, `AccessListItem` definitions — notice fields and types (ObjectID vs string).

4. `internal/storage/access_control/init.go` and `internal/storage/access_control/helper.go`

- Why: concrete example of a storage module using the DAL and mapping models → entities.
- What to look for:
  - How `NewAccessControlInfrastructure` calls `dal.NewMongoDal[...]` with concrete types.
  - `GetAccessList`, `GetUserAccess`, `UpdateUserAccess` implementations: how DAL methods are called, how results are mapped into `entities`.
  - The helper `convertSubAccessList` — recursive mapping for nested lists.

5. `internal/constants/interfaces/outbound/ports.go`

- Why: defines the interface `AccessControlPorts` that the storage must satisfy. Understanding the interface clarifies what methods service expects from storage.
- What to look for: `AccessControlPorts` method signatures and their return types.

6. `internal/service/access_control/*` (`init.go`, `abstract.go`, `access_control.go`)

- Why: shows how the service holds a field typed as the interface and calls storage methods. It also shows error handling and how storage responses are used.
- What to look for: `Service` struct, `InitStore`, `GetAccessList` calling `s.store.GetAccessList(ctx)`.

7. `initiator/initiator.go`

- Why: shows runtime wiring — where the concrete `AccessControlOutbound` created earlier is passed to `access_service.InitStore(...)` and ultimately to handlers.
- What to look for: `NewAccessControlInfrastructure(...)` call and `access_service.InitStore(...)` call.

8. `pkgs/audit/mongo_dal_wrapper.go` and `internal/storage/audit_wrapper.go`

- Why: shows how DALs can be wrapped to add auditing/logging. This helps you see decorator/wrapper patterns in the storage layer.
- What to look for: `AuditedMongoDal` type, `NewAuditedMongoDal`, and how `WrapMongoDalWithAudit` produces audited wrappers.

9. Other storage modules (examples)

- Open `internal/storage/transfer_limit/init.go` and `internal/storage/member/init.go`
- Why: compare patterns — you will see the same `dal.NewMongoDal` usage and mapping code repeated with different models.

What to practice (exercises)

Exercise A — Trace one request by reading code

1. Pick `GetAccessList` path:
   - Open the handler that calls service (search `GetAccessList` in `internal/handlers/rest` or follow `routing.InitAccessControlHandler` in `internal/glue/routing.go`).
   - Open `internal/service/access_control/access_control.go` and find `GetAccessList`.
   - Open `internal/storage/access_control/init.go` and step through `GetAccessList` until the DAL call.
   - Open `internal/storage/peristance/init.go` and read `FindAll` to see the driver call.
2. Write a 4–6 line comment in a scratch file describing the flow in your own words.

Exercise B — Compare `models` and `entities`

1. Open `internal/constants/models` and `internal/constants/entities` and pick `UserAccess` or `AccessListItem`.
2. Note differences (ObjectID, tags, field names) and write one sentence explaining why they differ.

Exercise C — Create a local mock (unit test)

1. Create a test file `internal/service/access_control/mock_test.go`.
2. Implement `MockAccessControl` struct with the three methods from `AccessControlPorts` returning controlled data.
3. Instantiate service with `accesssvc.InitStore(mock, testLogger)` and call `svc.GetAccessList(ctx)` in the test — assert expected values.

Tips while reading code

- Search for `NewMongoDal` to find all places DAL is used.
- Pay attention to `context.Context` being passed: that controls timeouts and cancellation.
- Look for `HandleMongoError` or `HandleStorageError` to see how DB errors are converted into app errors.
- For mapping code: watch for `ObjectID.Hex()` — that converts DB IDs to string IDs for domain entities.

Quick glossary (beginner)

- DAL (Data Access Layer): small helpers that talk to the DB (here: `MongoDal`) so other code can call methods like `FindAll` without importing the raw Mongo driver.
- Interface: a set of method signatures (a "contract"). The service depends on the interface, not the concrete implementation.
- Concrete implementation: the struct that provides methods (here: `AccessControlOutbound`) and talks to the DAL.
- Mapping: converting between DB `models` and app `entities` (shape/field differences).

Useful commands

```
# search for DAL usage
rg "NewMongoDal" -n
# find where GetAccessList is defined/used
rg "GetAccessList" -n
# open file quickly (replace with your editor)
code internal/storage/access_control/init.go
```

Next steps (after this guide)

- Implement Exercise C (mock test) — this will force you to wire the service and see the runtime dispatch.
- Ask for help on one function you don't understand and I will annotate it line-by-line.

If you want, I can also:

- Create the mock test file for you under `internal/service/access_control` (I can add it to the repo), or
- Add inline comments to `internal/storage/peristance/init.go` explaining `FindAll` and `FindOne` line-by-line.

Happy to continue — tell me which next step you'd like me to do for you.
