# aws-ses-web-ui

aws-ses-web-ui is a simple web interface for sending email via AWS SES. It requires an AWS account with SES enabled and
at least one verified identity. This project is meant for individuals to run on their personal computers.

This project is a single Golang executable that runs a web server. The web UI is built using:

* Golang HTML templates
* jQuery
* Tailwind CSS
* FontAwesome (free)

## Screenshots

![form.png](screenshots/form.png)
![confirmation.png](screenshots/confirmation.png)
![sent.png](screenshots/sent.png)

## Installation

### Use Docker Compose

In order to run using Docker Compose, install the latest version of Docker Compose.

This option requires a configuration file located at `config.json`. See the configuration section below for more
information.

The default `docker-compose.yml` PostgreSQL DSN is `postgres://postgres:password@aswupostgres:5432/aswu`.

Once the `config.json` file has been created run the following command in the root directory of this project:

```
docker compose up
```

After the services are running, visit `http://localhost:8080`.

### Install from source

In order to install from source, install Golang 1.21.3 or higher. Build the JSON configuration (see below), then run the
following command in the root directory of this project:

```
go run cmd/server/main.go
```

After the program has started, visit `http://localhost:8080`.

## Configuration

This project uses a JSON configuration file similar to the one below. The program looks for this configuration file in
one of the following locations:

1. The contents of the environment variable `CONFIG_JSON`.
2. The file path specified in the environment variable `CONFIG_PATH`.
3. The file in the current directory named `config.json`.

### Example:

```json
{
  "aswu": {
    "allowedFrom": [
      "@example.com",
      "singleaddress@otherexample.com"
    ],
    "defaultFrom": "John Doe <john.doe@example.com>",
    "devMode": false,
    "usePostgres": false
  },
  "postgres": {
    "dsn": "postgres://postgres:password@aswupostgres:5432/aswu"
  },
  "ses": {
    "awsRegion": "us-east-1",
    "accessKeyID": "my-access-key-id",
    "secretKey": "my-secret-key"
  }
}
```

The `aswu` section contains configuration options for the web interface.

| Attribute     | Value                                                                               |
|---------------|-------------------------------------------------------------------------------------|
| `allowedFrom` | An array of AWS SES verified emails or domains prefixed with `@`.                   |
| `defaultFrom` | The default email address to send email from.                                       |
| `devMode`     | Set to `true` when developing the aws-ses-web-ui project itself.                    |
| `usePostgres` | Set to `true`  to save outgoing emails to the PostgreSQL database at the given DSN. |

The `postgres` section contains configuration options for the PostgreSQL database. It is optional when `usePostgres` is
set to `false`.

| Attribute | Value                                                               |
|-----------|---------------------------------------------------------------------|
| `dsn`     | The PostgreSQL Data Source Name (DSN). This is a connection string. |

The `ses` section contains configuration options for AWS SES.

| Attribute     | Value                  |
|---------------|------------------------|
| `awsRegion`   | The AWS region to use. |
| `accessKeyID` | The AWS access key ID. |
| `secretKey`   | The AWS secret key.    |
