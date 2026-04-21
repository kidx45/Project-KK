# Notes on Golang Course

## 🐳 Docker

- **Detached Mode (`-d`)**: Runs the container in the background, allowing you to use your terminal while it runs.
  - _Example:_ `docker run --name some-name -p 8080:8080 -e ENV_VAR=SOME_VALUE -d image-name`
- **Pulling Images**: Use `docker pull image-name:version` (version is optional).
- **Port Mapping (`-p`)**: Maps the container's port to the localhost (host machine) port. Since they run on different ports, the host fetches data through the mapped port.
- **Images vs. Containers**:
  - **Image**: A blueprint used to export the application across servers.
  - **Container**: A running instance of an image where the app actually executes.
- **Interactive Mode**: To run commands inside a running container, use:
  - `docker exec -it <name-or-id> <command>`
- **Logs**: To view actions occurring in a detached container, use:
  - `docker logs <container-name-or-id>`
- **Management**:
  - `docker stop <name>`: Stops a running container.
  - `docker rm <name>`: Deletes a container.
  - `docker ps -a`: Lists all containers (running or exited).

## 🚀 DB Migration

Database migration is a core backend principle used to track schema changes and version the database.

- **Golang Library**: `golang-migrate`
- **Create Migration Files**:
  - `migrate create -ext sql -dir db/migration -seq <name>`
- **Execution Command**:
  - `migrate -path <schema_location> -database "<db_uri>?sslmode=disable" -verbose up/down`

## 🌐 Nginx

Nginx is an open-source web load balancer used for secure and fast routing of client data to the backend.

## 🛠 Makefile

- **`.PHONY`**: Used to ensure `make` runs the command even if a file with the same name exists in the directory.

## 🗃 SQLC

### `sqlc.yaml` Configuration

- **`name`**: Specifies the name of the Golang package generated to access DB queries.
- **`path`**: The directory where the generated Go code will be saved.
- **`query`**: The location where `sqlc` reads your SQL query files.
- **`schema`**: The directory where your migration schema lives.
- **`emit_json_tag`**: Adds JSON tags to the generated Go structs.
- **`emit_prepared_queries`**: Used to boost performance by using prepared statements.
- **`emit_interface`**: Creates interfaces for the query methods, which is useful for testing/mocking.
- **`emit_exact_table_names`**: Forces the Go struct name to match the SQL table name exactly.

## 📶 Context in Go

Think of a Context as a signal carrier that travels through your entire program, from the moment an HTTP request hits your server until the moment the database returns a result. It has two main things

- **`Cancellation signal`** - tells the server to stop processing if say a client closes the browser
- **`Deadline`** - set a timeout for how long request must wait for something

Context package provides:

- **`context.Background()`**: The "empty" starting point for a context (usually used at the very top level).
- **`context.WithTimeout()`**: Adds a countdown (e.g., 5 seconds).
- **`context.WithCancel()`**: Allows you to manually trigger a "stop" signal.
- **`context.WithValue()`**: Allows you to pass small bits of data (like a Request-ID for logging) all the way down to the database layer.

## 📑 Prepared Statments

It is an optimization because it changes how the database processes your SQL queries to make them faster and more secure. In the regular way, for each request you send that requires a database, every single time:

- **`Parsing`**: Read the string and check for syntax errors.
- **`Analysis`**: Check if the table accounts and the column id actually exist.
- **`Planning`**: Figure out the fastest way to get the data (e.g., "Should I use an index or scan the whole table?").
- **`Execution`**: Finally, run the plan and return the data.

With the Prepare function in your db.go, the database does steps 1, 2, and 3 only once when your application starts up.

- The database "pre-compiles" the query and stores the "execution plan" in its memory.
- Each time you actually call GetAccount later, the database skips the parsing and planning steps. It just plugs in your variables (like the specific ID) and jumps straight to Execution.

## 🧪 Testing database integration using testify and testing packages

When testing an integration with a database, it is best to use a test database to check whether or not the generate written by the programmer doesn't have any logical and syntax error. In golang **`Testing`** package we have,

- **`Testing.M`** - It manages the entire group of tests in a single package.
  Key Function: m.Run(). This is the big button that starts ALL the tests in that package.
- **`Testing.T`** - It handles a single, specific test. It is mainly used for checking logic, assertions, and reporting errors for one specific piece of code.

## Database Transaction

It is a single unit of work involving multiple operation for executing a specific operation such in the case of our project, it may be to transfer money which requires getting and retrieving the accounts storing the entries from both accounts and storing the transfer of the money while dealing with logic. Reasons for its need are:

- provide reliable and consistent unit of work, even when system failure occurs
- provide isolation between programs that access the DB

This is done to satisfy the ACID property
### Usage Tips

- **Comments**: Must follow an exact format: `-- name: <method_name> :<return_type>`.
  - _Example:_ `-- name: GetAccount :one`
  - _Note:_ `--name` (without space) will not work.
