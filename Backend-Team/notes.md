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

## 📑 Database Transaction

It is a single unit of work involving multiple operation for executing a specific operation such in the case of our project, it may be to transfer money which requires getting and retrieving the accounts storing the entries from both accounts and storing the transfer of the money while dealing with logic. Reasons for its need are:

- provide reliable and consistent unit of work, even when system failure occurs
- provide isolation between programs that access the DB

This is done to satisfy the ACID property

## 🛣️ Go Routine, channels and Context passing for debuging
### Go Routine
Go Rountine is a property found in golang that allows us to simulate real world traffic and when they hit the same db at the same time which is very useful in tests which will allow us how or db is accessed and preventes deadlocking in the future. It is mainly used in the HTTP layer if needed explicit definition but frameworks like Gin Framework and chi Framework automatically have that so there is no need to specify the go func().

### Channels
Whenever the go routine is used, it run in its own background outside of function it lives in hence the function may finish before the db and we won't be able to read the result. So channels act like a pipline that allows us to connect the functiona and go routine passing data that is needed and later accessed using the arrow symbol

## Context
It is basically a package passed around function to function using for debuging and passing data without the needed to change the parameters of our fucntions that we are passing it into. In a context, we need a key inorder to passing and get the value we silently passed into the function and passsing an argument may cause a collision hence we use the struct.

## 🔒️ DeadLock occurence
A deadlock occurs in our case when two transactions are trying to access resources of eachother that can't be aquired unless the other finishes their operation like they either commit or role back hence we achieve a deadlock. Ways to avoid include always start a db transaction and finish it before starting another and always locking the same resources in the same order. (this is to only avoid deadlocks) but this doesn't allow us to run transactions concurrently. Hence by editing our db schema and specifying how the operations execute (insert, update, etc), we are able to avoid deadlocks and achieve concurrency.

## 🎚️ Isolation Level in psql and mysql
When ever a transaction occurs, it must follow and satisfy the ACID property which mean
- **`Atomicity`** - Operations in atransaction must in its whole successed or fail
- **`Consistant`** - Database after a transaction must be valid
- **`Isolation`** - No transaction will affect another transaction
- **`Durability`** - A success transaction must be able to record requried data back to its database

If a transaction affect the data retrieved and seen by another running transaction, then a read phenomenon will occur. To solve this, we use a set of isolation levels specified by the ANSI which are read uncommited, read commited, repeatable read, serializable.
- **`read uncommitted`** - low isolation level which allow for other transactions to see changes that are made by other transaction that haven't been commited yet which leads to **`dirty read`**
- **`read committed`** - isolation level which allow for other transactions to see changes made by other transactions that have committed leadin **`non-repeatable read`**
- **`repeatable read`** - isolation level which won't allow for transactions to see changes made by a commited transaction but will lead to a situation where a set of transaction were done sequentially won't make sense which is also called the **`serialization anomaly`**
- **`serializable`** - highest isolation level which allows for transactions to strictly run sequentially to avoid any anomaly which will solve the issue that arise from different transaction affecting the values of the other running transactions 

In psql however, there is read-uncommited and read-committed behave the same and when trying to update a value in repeatable read, it will throw and error instead of actually updating like in mysql. Also a deadlock might occur in serializable. While postgres uses dependency detection to avoid serialization anomaly, mysql uses a locking mechanism.

## 📠 CI Integration using github's actions
In order to track new changes that occur into our github repository and above potential errors and bugs, we use an automated workflow either through actions, jetkins ... which are used for automating the build and running process.
Here, workflow is a automated procedure consisting of jobs which can be trigged either
- when an event occurs on the repository
- when an scheduled trigger is on
- when manually pressing in the UI

Inorder to create a workflow, we need to add a .github/workflow/some_name.yml file into the repository.

At the beginning of a yml file, we would have the
- **`name`** - name of the workflow
- **`on`** - specifies what actions like a push or schedule will enable the workflow other than the manuel trigger

In order to run the set of jobs that are found in the workflow we need a runner which is a server that runs one jobs at a time and specfied in the using such command
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

To also add service to the workflow as postgres in our case it is best to include the serives tag provided by the vendor of your choice. Also we can use a new line to using the "|" symbol in front of the run.



### Usage Tips

- **Comments**: Must follow an exact format: `-- name: <method_name> :<return_type>`.
  - _Example:_ `-- name: GetAccount :one`
  - _Note:_ `--name` (without space) will not work.

- **Composition or Struct Embedding**: inheritance in Golang, instead of using `extends` keyword, we use embedding, by adding our custom struct inside another struct which is our Queries. This allows us to inherit methods from the embedded struct.

- **optimization is king**: if you have the opportunity to make say your sql querying shorter then it is best to do so. Analyze it carefully

- **Have the ability retry transactions in the case of a deadlock**