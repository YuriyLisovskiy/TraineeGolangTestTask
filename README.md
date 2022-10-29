## Trainee Golang Test Task
The environment consists of two Docker containers that run the REST API application
and the PostgreSQL database with persistent storage located on host machine.
But you can start these two containers separately.

#### REST API
REST API is a CLI application based on [Gin](https://github.com/gin-gonic/gin) web framework.
The app serves three endpoints, as described [here](#documentation). Upload of large files
(up to 4 GB) is available.

The default constant values that are used for application configuration can be changed
with help of the following environment variables:

| Variable                   | Type             | Description                                                     |
|----------------------------|------------------|-----------------------------------------------------------------|
| `PORT`                     | positive integer | Port number to listen on by a web server                        |
| `APP_PAGE_SIZE`            | positive integer | The count of database records per one response                  |
| `GIN_DEBUG`                | boolean          | Run the application in debug (`true`) or release (`false`) mode |
| `GIN_MAX_MULTIPART_MEMORY` | positive integer | The upper limit of memory allocated for multipart requests      |
| `POSTGRES_HOST`            | string           | Host name of the database server                                |
| `POSTGRES_PORT`            | positive integer | The port of the database server                                 |
| `POSTGRES_USER`            | string           | REST API user for accessing the database                        |
| `POSTGRES_USER_SCHEMA`     | string           | A separate schema created for `POSTGRES_USER`                   |
| `POSTGRES_PASSWORD`        | string           | Super secure password                                           |
| `POSTGRES_DB_NAME`         | string           | The database to connect to                                      |

Default values of the above environment variables can be found in [`.env`](.env) file.

REST API application contains a command for migrating the database (DB) schema - `migrate`.
Running this app without any commands triggers the server's startup.

The application manages data using PostgreSQL DB. To achieve this, the [Gorm](https://github.com/go-gorm/gorm)
library is used. Migration of the DB schema is configured to execute before the server startup.

To migrate the database manually, you need to configure `POSTGRES_*` environment variables,
build the app, and run the following command:
```shell
./rest-api-app migrate
```

#### PostgreSQL Database
The database deployment is performed with [postgres](https://hub.docker.com/_/postgres) Docker image.
All DB configurations are performed under the administrator user called `postgres`.
The initial database user and his permissions are set up using [init.sql](sql/init.sql) script.
Due to docker being the stateless container, the database is configured to be persistent on
local machine.

### Deployment with Docker
Deploy API app and the database:
```shell
docker compose up
```
or deploy them separately:
```shell
docker compose run api --rm --detach
docker compose run postgres_database --rm --detach
```

The reference for `docker compose run` command options can be found
[here](https://docs.docker.com/engine/reference/commandline/compose_run/#options).

To stop containers, press `Ctrl+C` or run the command below if you started them in detached mode:
```shell
docker compose down
```

Perform cleaning up, if required:
```shell
docker compose rm --force --stop
```

To learn more about `docker compose rm` options, read the
[reference](https://docs.docker.com/engine/reference/commandline/compose_rm/#options).

### Documentation
The API documentation is available
[here](https://app.swaggerhub.com/apis-docs/YuriyLisovskiy/TraineeGolangTestTask/1.0.0).
