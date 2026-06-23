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

## 📑 Prepared Statements

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

## 📑 Database Transaction

It is a single unit of work involving multiple operation for executing a specific operation such in the case of our project, it may be to transfer money which requires getting and retrieving the accounts storing the entries from both accounts and storing the transfer of the money while dealing with logic. Reasons for its need are:

- provide reliable and consistent unit of work, even when system failure occurs
- provide isolation between programs that access the DB

This is done to satisfy the ACID property

## 🛣️ Go Routine, channels and Context passing for debugging
### Go Routine
Go Routine is a property found in golang that allows us to simulate real world traffic and when they hit the same db at the same time which is very useful in tests which will allow us how or db is accessed and prevents deadlocking in the future. It is mainly used in the HTTP layer if needed explicit definition but frameworks like Gin Framework and chi Framework automatically have that so there is no need to specify the go func().

### Channels
Whenever the go routine is used, it run in its own background outside of function it lives in hence the function may finish before the db and we won't be able to read the result. So channels act like a pipeline that allows us to connect the function and go routine passing data that is needed and later accessed using the arrow symbol

## Context
It is basically a package passed around function to function using for debugging and passing data without the needed to change the parameters of our functions that we are passing it into. In a context, we need a key in-order to passing and get the value we silently passed into the function and passing an argument may cause a collision hence we use the struct.

## 🔒️ DeadLock occurrence
A deadlock occurs in our case when two transactions are trying to access resources of each-other that can't be acquired unless the other finishes their operation like they either commit or role back hence we achieve a deadlock. Ways to avoid include always start a db transaction and finish it before starting another and always locking the same resources in the same order. (this is to only avoid deadlocks) but this doesn't allow us to run transactions concurrently. Hence by editing our db schema and specifying how the operations execute (insert, update, etc), we are able to avoid deadlocks and achieve concurrency.

## 🎚️ Isolation Level in psql and mysql
When ever a transaction occurs, it must follow and satisfy the ACID property which mean
- **`Atomicity`** - Operations in a transaction must in its whole succeed or fail
- **`Consistent`** - Database after a transaction must be valid
- **`Isolation`** - No transaction will affect another transaction
- **`Durability`** - A success transaction must be able to record required data back to its database

If a transaction affect the data retrieved and seen by another running transaction, then a read phenomenon will occur. To solve this, we use a set of isolation levels specified by the ANSI which are read uncommitted, read committed, repeatable read, serializable.
- **`read uncommitted`** - low isolation level which allow for other transactions to see changes that are made by other transaction that haven't been committed yet which leads to **`dirty read`**
- **`read committed`** - isolation level which allow for other transactions to see changes made by other transactions that have committed lead into **`non-repeatable read`**
- **`repeatable read`** - isolation level which won't allow for transactions to see changes made by a committed transaction but will lead to a situation where a set of transaction were done sequentially won't make sense which is also called the **`serialization anomaly`**
- **`serializable`** - highest isolation level which allows for transactions to strictly run sequentially to avoid any anomaly which will solve the issue that arise from different transaction affecting the values of the other running transactions 

In psql however, there is read-uncommitted and read-committed behave the same and when trying to update a value in repeatable read, it will throw and error instead of actually updating like in mysql. Also a deadlock might occur in serializable. While postgres uses dependency detection to avoid serialization anomaly, mysql uses a locking mechanism.

## 📠 CI Integration using github's actions
In order to track new changes that occur into our github repository and above potential errors and bugs, we use an automated workflow either through actions, jenkins ... which are used for automating the build and running process.
Here, workflow is a automated procedure consisting of jobs which can be triggered either
- when an event occurs on the repository
- when an scheduled trigger is on
- when manually pressing in the UI

In-order to create a workflow, we need to add a .github/workflow/some_name.yml file into the repository.

At the beginning of a yml file, we would have the
- **`name`** - name of the workflow
- **`on`** - specifies what actions like a push or schedule will enable the workflow other than the manuel trigger

In order to run the set of jobs that are found in the workflow we need a runner which is a server that runs one jobs at a time and specified in the using such command
```yml
jobs
  build:
    run-on: ubuntu-latest
```
A job is a set of steps that are run either in parallel is not dependant on each other or sequentially a job depends on the other job.
Example:

```yml
jobs
  build:
    run-on: ubuntu-latest
    steps:
      name: name_of_the_step
      use: uses_some_file_or_folder
      run: runs_some_command_on_the_terminal
  
  test:
    needs: build
    run-on: ubuntu-latest
    steps:
      name: name_of_the_step
```

As we can see, here we have two jobs, build and test, runner is specified using the **`run-on`** keyword. Test needs the build to do its thing first so this job runs after build in a sequential manner. Steps are tasks run sequentially when a job is ran. It contains actions which are a set of commands that by themselves also run sequentially. In a step, there may be multiple actions running.

To also add service to the workflow as postgres in our case it is best to include the services tag provided by the vendor of your choice. Also we can use a new line to using the "|" symbol in front of the run.

### GitHub Actions versions

- **`actions/checkout@v4`**: Downloads the repository code into the GitHub runner so later steps can access the files.
- **`actions/setup-go@v4`**: Installs and configures Go on the runner, using the version specified under `go-version`.
- **`@v4`**: Means the workflow is using version 4 of that GitHub Action. It keeps the workflow stable while still receiving compatible updates within that major version.

## ENV variables
Reading from file - to load default configuration for local development
Reading from an env file - to load configuration to override the local configuration
Viper is a really good way to read and load configration from config files.

## Validators
In go, when validating a request, we can create a custom validator like the one in JavaScript where we use Joi. 

## Why use an Interface for the Store?
In the `Server` struct, we use the `db.Store` interface instead of a pointer to the concrete `db.SqlStore`. This is a classic example of **Dependency Injection**.

- **Testability:** By using an interface, we can inject a `MockStore` during unit testing. This allows us to test our API handlers without needing a real database connection.
- **Decoupling:** The API logic doesn't care how the data is stored (SQL, NoSQL, or even a simple map); it only cares that the object provided follows the "contract" defined by the interface.
- **SqlStore:** This is the real implementation that uses `*sql.DB` to talk to the database.
- **MockStore:** This is a fake implementation (usually generated by `gomock`) used in tests to return predefined results.

## Testing using gomock and custom gomock matcher
When the test runs, the API receives the raw password in the request body, hashes it, and calls store.CreateUser(ctx, params) with the hashed password inside params — gomock intercepts that real call and passes params as x into Matches(x). Meanwhile, back in the test setup, you built arg (your checklist of expected values without the hash) and passed it along with the raw password into EqCreateUserParams. So inside Matches, you have two things: x which is what the API actually sent (including the hash), and e.user which is your checklist — it first verifies the hash in x came from the raw password using CheckPassword, then patches that hash into e.user so DeepEqual can compare everything together, essentially saying "the API called CreateUser with the right username, email, fullName, phoneNumber AND a valid hash of the password".

### Usage Tips

- **Comments**: Must follow an exact format: `-- name: <method_name> :<return_type>`.
  - _Example:_ `-- name: GetAccount :one`
  - _Note:_ `--name` (without space) will not work.

- **Composition or Struct Embedding**: inheritance in Golang, instead of using `extends` keyword, we use embedding, by adding our custom struct inside another struct which is our Queries. This allows us to inherit methods from the embedded struct.

- **optimization is king**: if you have the opportunity to make say your sql querying shorter then it is best to do so. Analyze it carefully

- **Have the ability retry transactions in the case of a deadlock**

- **Exporting stuff**: When in need to export stuff we use a capital letter so use a capital letter always a default instead of using small letters

- **Public/Private Keys and Symmetric Signing**: When signing JWTs there are two approaches. Symmetric signing (HS256) uses one secret key that both signs and verifies tokens. The problem with this in a microservice architecture is that every service needs the secret, meaning if any single service is compromised the attacker can forge tokens for everything. Asymmetric signing (RS256) solves this by splitting into a private key that only the auth server holds for signing, and a public key handed out to every microservice just for verification. A compromised microservice only leaks the public key which is useless for forging tokens — you can't sign with it. This makes asymmetric easier to manage at scale since there is only one place to protect, the private key on the auth server.

    It is also worth distinguishing signing from encryption since they look similar but do opposite things. In encryption you lock a message with someone's public key so only their private key can unlock it — the goal is hiding content. In signing you sign with your private key so anyone with your public key can verify you created it — the goal is proving authenticity. JWT uses signing, not encryption. The payload is just base64 and anyone can read it, the signature just proves the auth server created it and it was not tampered with.

- **Access and Refresh Tokens**: When a user logs in the server issues two tokens. The access token is short lived (5-15 minutes) and stored in a plain JavaScript memory variable. It gets sent in the Authorization header on every request to prove who the user is. The refresh token is long lived (7-30 days) and stored in an HttpOnly cookie, meaning the browser manages it and JavaScript cannot read it at all. This separation is intentional — the refresh token is the valuable one so it gets the strongest protection.

    When the access token expires the browser automatically sends the refresh token cookie to the `/refresh` endpoint, the server verifies it against the database, and issues a new access token seamlessly without the user noticing. When the refresh token itself expires the user gets logged out and has to login again, exactly like Udemy logging you out after long inactivity. The HttpOnly cookie also protects against XSS attacks where malicious JavaScript running in the browser tries to steal tokens — it simply cannot see HttpOnly cookies at all. The tradeoff with short expiry is more `/refresh` calls under high load, which is why Redis is commonly used for refresh token storage at scale since it is in memory and much faster than a regular database lookup.

- **Goroutines in `main.go`**:
    - **Non-Blocking Servers**: The `app.Start()` function (and the pprof server) are "blocking" calls. If you run them normally, the code stops there and waits forever for the server to exit.
    - **The "Main Thread" Role**: By putting the server in a `go routine`, we free up the main thread to reach the **Signal Listener** (`<-quit`). 
    - **Graceful Shutdown**: This setup allows the app to stay alive until you press Ctrl+C. When that happens, the main thread "wakes up" and calls `app.Shutdown()`, ensuring database connections and traces are closed cleanly before the program exits.

- **Telemetry (OTEL)**:
    - **The Concept**: It’s the "Black Box" for your app. It collects **Traces** (the path of a request) and **Logs** to help you understand performance and errors.
    - **OpenTelemetry (OTEL)**: An industry **standard** (like a universal connector) that allows your app to send data to any monitoring tool (Jaeger, Grafana, etc.) without changing your code.
    - **The Collector**: A middleman service that receives data from your app, cleans it up, and sends it to your chosen external provider.
    - **Performance (No Lag)**: It doesn't slow down the app because it:
        - **Samples**: Only records a percentage of requests (e.g., 10%).
        - **Asynchronous**: Sends data in the background (using Goroutines) so the user doesn't wait.
        - **Batching**: Sends groups of logs at once rather than one-by-one.

# Dockerfile Explanation & Summary

A **Dockerfile** is a blueprint used to package an application into an isolated environment called a **container**. This ensures the application runs exactly the same way on any machine (your laptop, a teammate's computer, or a cloud server).

---

## 🛠️ Step-by-Step Command Breakdown

* `FROM golang:1.26-alpine3.24`
    * **The Base Environment:** Starts with a lightweight Linux operating system (Alpine) that comes with Go version 1.26 pre-installed. You don't have to install Go manually.
* `WORKDIR /app`
    * **The Working Directory:** Creates an internal folder named `/app` inside the container and switches into it. All subsequent actions happen here.
* `COPY . .`
    * **Copying Files:** Takes all the source code files from your actual machine (the first `.`) and copies them into the container's `/app` folder (the second `.`).
* `RUN go build -o main main.go`
    * **The Build Step:** Compiles the Go source code inside the container into a single, executable binary file named `main`.
* `EXPOSE 8080`
    * **The Label/Documentation:** Documents that the application inside the container intends to listen for network traffic on port `8080`. *(Note: This does not actually open the port or run code).*
* `CMD ["./main"]`
    * **The Launch Command:** The final instruction that actually triggers and runs your compiled Go binary (`main`) when the container starts up.

---

## ❓ The `EXPOSE` Command: Who Uses It & Where is It Viewed?

`EXPOSE` does **not** run code or automatically open ports to the outside world. It acts as an official piece of metadata (a "contract") between the developer and the infrastructure.

### 👥 Who Uses It?
* **DevOps & Cloud Engineers:** They read `EXPOSE` to write deployment configurations (like Kubernetes Services, AWS Task Definitions, or Nginx reverse proxies) so they know where to route external internet traffic.
* **Developers:** It acts as instant code documentation. New team members can look at the Dockerfile and immediately know which port the application communicates on without digging through the Go source code.

### 🖥️ Where is it Viewed?
* **Terminal (`docker ps`):** Docker displays exposed ports under the `PORTS` column when listing active containers.
* **Docker Desktop:** The graphical user interface reads this line to display the port and provide a clickable hyperlink to open the app in your browser.
* **Cloud Automation Tools:** Platforms like AWS ECS or Google Cloud Run scan this metadata to automatically pre-fill or suggest network port mapping settings during deployment.

### 🔗 Activating the Port
To actually connect your computer's browser to the container, you must explicitly bind the ports using the `-p` flag when running the container:
```bash
docker run -p 8080:8080 <image-name>

## 🏷️ The `-t` Tag (Used in Terminal during `docker build`)

* **What It Means:** Stands for **Tag**. It is used in the terminal command when building your container image.
* **The Use Case:** Prevents random, cryptic IDs by giving your image a human-readable name and version (e.g., `my-go-app:1.0`).
* **Command Syntax:** Run as `docker build -t my-go-app:1.0 .` where `my-go-app` is the name, `:1.0` is the version, and `.` means current folder.
* **Version Control:** DevOps teams use it to track releases (v1.0, v1.1) and easily roll back to previous versions if a bug occurs.
* **Cloud Deployment:** Required by platforms like AWS, Google Cloud, or Docker Hub to match their registry format before pushing code.