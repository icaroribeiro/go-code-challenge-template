# Hi there! 👋

Be very welcome to my solution to X's code challenge.

- [Introduction](#introduction)
- [Architecture](#architecture)
- [Database](#database)
- [How to run the project?](#how-to-run-the-project)
- [API documentation](#api-documentation)
- [Test cases](#test-cases)
- [How to run the tests?](#how-to-run-the-tests)
- [Deployment](#deployment)
- [How to deploy the project?](#how-to-deploy-the-project)
- [References](#references)

## Introduction

This project consists of the development of a **REST API** using **Go** programming language, **Json Web Token** and **Postgres** database for managing authentication operations and accessing users data.

## Architecture

The architecture of the project was designed using the concepts of [Domain Driven Design](), [Clean Architecture]() and [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/).

## Database

Two Postgres dabases need to be configured to use the project. One of them is intended to common (or usual) use and the other is directed to test execution. However, both of them contain the same tables named auths, logins and users defined in the **database/scripts/1-create_tables.sql** file.

### Tables

**Auths**

The **auths** table contains the authentication data.

| Fields     | Data type | Extra                       |
|:-----------|:----------|:----------------------------|
| id         | UUID      | NOT NULL PRIMARY KEY        |
| user_id    | UUID      | NOT NULL UNIQUE FOREIGN KEY |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP   |

**Note**:

A record is created in this table whenever a user performs login and this same record is deleted as soon as the related user performs logout.

**Logins**

The **logins** table contains the users credentials.

| Fields     | Data type | Extra                       |
|:-----------|:----------|:----------------------------|
| id         | UUID      | NOT NULL PRIMARY KEY        |
| user_id    | UUID      | NOT NULL UNIQUE FOREIGN KEY |
| username   | TEXT      | NOT NULL                    |
| password   | TEXT      | NOT NULL                    |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP   |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP   |

**Users**

The **users** table contains the users data.

| Fields          | Data type | Extra                     |
|:----------------|:----------|:--------------------------|
| id              | UUID      | NOT NULL PRIMARY KEY      |
| username        | TEXT      | NOT NULL UNIQUE           |
| created_at      | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP |
| updated_at      | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP |

## How to run the project?

The project can be run either **locally** or using a [**Docker**](https://www.docker.com/) container. However, in order to facilitate explanations, this documentation will focus on running using a Docker container.

### Makefile file

A **Makefile** file was created as a single entry point containing a set of instructions to run the project in these two different ways via commands in the terminal.

Furthermore, this file also contains a series of routines used throughout the development of the project, such as reformatting the **.go** file and printing style errors, generating API documentation, creating *mocks* used in tests of the solution, among others.

To run the project with a Docker container, run the command:

```
make startup-app
```

Note:

- The **.env** file contains the environment variables used by the Docker container. However, it is not necessary to make changes to this file before running the project, so the variables can be kept as they are defined.

To close the application, run the command:

```
make shutdown-app
```

## API documentation

### API endpoints

The API *endpoints* were documented using the Github repository called [swaggo/swag](https://github.com/swaggo/swag) which converts code annotations in **Go** into **Swagger 2.0** documentation based on **Swagger** files located in the **docs/api/swagger** directory.

After running the project, access the following URL through your web browser to view an HTML page that illustrates the information of the API *endpoints*:

```
http://{host}:8080/swagger/index.html
```

### Postman Collection

To support the use of the API, it was created the file **new-go-code-challenge-template.postman_collection.json** in the directory **docs/api/postman_collection** which contains a group of requests that can be imported into the **Postman** tool (an API client used to facilitate the creation, sharing, testing and documentation of APIs by developers.).

## Test cases

The test cases were designed as [**Table Driven Tests**](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests) so that each test case was built by declaring a structure that contains actions that can be performed before and after executing them, as well as expected inputs and outputs, following the **unit** and **integration** tests approaches.

### Unit Tests

The unit tests are located inside the **internal** and **pkg** directories at the project root.

They are evaluated using the **Black-Box** testing strategy, where the test code is not in the same package as the code under evaluation.

For this, files were created with the suffix **_test** added to their names and also to the names of their test packages. For example, the code from the package (*pkg*) called **validator** is tested by a file called **validator_test.go**, which is defined in another package, called **validator_test**.

The separation of codes into distinct packages aims to ensure that only the identifiers exported from the packages under evaluation are tested. By doing this, the test code is compiled as a separate package and then linked and run with the main test binary.

#### Mocks

Some of the tests were written using mock objects in order to simulate dependencies so that the layers could interact with each other through **interfaces** rather than concrete implementations, made possible by the *design pattern* of **Dependency Injection**.

Basically, the purpose of mocking is to isolate and focus on the code being tested and not on the behavior or state of external dependencies. In simulation, dependencies are replaced with well-controlled replacement objects that simulate the behavior of real ones. Thus, each layer is tested independently, without relying on other layers. Also, you don't have to worry about the accuracy of the dependencies (the other layers).

For the mocking purpose, the Github repositories called [DATA-DOG/go-sqlmock](https://github.com/DATA-DOG/go-sqlmock) e [vektra/mockery](https://github.com/vektra/mockery) were used for mocking the SQL driver behavior without needing to actually connect to a database and for generating the mock objects from interface, respectively.

### Integration Tests

The integration tests are located inside the **tests/api** directory at the project root.

They were written by combining and testing the project layers together to simulate the production environment.

Note:

- The unit and integration tests check a large and relevant part of the different components of the solution, but not all of them. In addition, not all tests written have **100%** coverage of the tested code.

## How to run the tests?

Before running the project tests, it is needed to start up the Docker containers named **api_container** and **postgrestestdb_container** successfully.

The **postgrestestdb_container** container is necessary to execute the integration tests and it can be initialized by running the command:

```
make start-deps
```

After all these containers are successfully initialized, to execute the tests of the project, run the command:

```
make test-app
```

After running any of the tests, it is possible to check the percentage of code coverage that is met by each test case displayed in the test execution output.

The statistics collected from the run are saved in the **docs/api/tests/unit/coverage.out** file for coverage analysis. To check the **unit** test coverage report informed in the **coverage.out** file, run the command:

```
make analyze-app
```

Notes:

- The **coverage.out** file contains only **unit** test execution statistics. (There are no statistics on the execution of the **integration** tests.)






## Deployment

The project was deployed as a container on **Heroku** hosting service using **Terraform** tool.

To do this, they were created a **heroku.yml** manifest file and a **Dockerfile** in the **deployments/heroku** directory that was used for the building of the project as a Docker container, in addition to some infrastructure resource components via code defined in the **deployments/heroku/terraform** directory. 

The **Dockerfile** was designed using a Docker's multi-stage image build feature that allows creating multiple images in the same Dockerfile:

In summary, the first FROM statement is related to an image that uses an alias "as builder" name to be referred later in the file and generates the intermediate layer where the Golang compilation happens and the second FROM statement is directed to an image from alpine where we simply copy "--from=builder" to get the executable from the intermediate layer.

The goal of this scheme is to build a final image to be provided to the end users that is as lean as possible, in other words, with a reduced size, containing only the binary application and the base operating system needed to run it. This way, it can be deployed quickly even in slow network conditions.

### How to deploy the project?

To deploy the project it was necessary to have a Heroku account and the Heroku CLI and Terraform softwares installed on the machine.

For more information on how to download and configure them, please check the official websites: 

- Heroku CLI - https://devcenter.heroku.com/articles/heroku-cli

- Terraform - https://www.terraform.io/downloads.html

After that, it was needed to create a setup_env.sh file.

Then, configure some environment variables in the **deployments/heroku/scripts/setup_env.sh** file that are related to others defined in **deployments/heroku/terraform/variables.tf** file:

The first and the second variables are related to Heroku Platform API settings and refer to the Heroku email address and a Heroku API key, respectivelly. The third one variable refers to the Heroku application name.

```sh
#!/bin/bash

# Heroku platform settings.
export TF_VAR_heroku_email=
export TF_VAR_heroku_api_key=

# Heroku application settings.
export TF_VAR_heroku_app_name=
```

```
TF_VAR_heroku_email=<HEROKU_EMAIL>
```

```
TF_VAR_heroku_api_key=<HEROKU_API_KEY>
```

```
TF_VAR_heroku_app_name=<HEROKU_APP_NAME>
```

In order to get a Heroku API key, execute the Heroku CLI command:

```
heroku login
```

This way, you wil be redirected to the browser so that you can perform login to Heroku.

After that, execute the command:

```
heroku auth:token
```

After configuring these variables, it was possible to execute the commands located in the Makefile for the deployment:

To initialize everything Terraform require to provision the infrastructure, execute the command:

```
make init/deploy
```

For example, it downloads the Heroku's provider plugin and stores it in a **.terraform** hidden folder.

The infrastructure resource components were defined in the **deployments/heroku/terraform/resources.tf** file and they refer to the API and the associated Heroku Postgres database.

To obtain the detail about what will happen in the infrastructure without making any change on it, execute the command:

```
make plan/deploy
```

To make changes required in order to reach the desired state of the configuration, execute the command:

```
make apply/deploy
```

After applying the changes, it was possible to create the database tables by means of CLI Heroku commands:

```
heroku pg:psql -a=<HEROKU_APP_NAME> <HEROKU_POSTGRES> < database/postgres/scripts/1-create_tables.sql
```

To identify what is the identifier of the Heroku Postgres database, execute the command:

```
heroku pg:info -a=<HEROKU_APP_NAME>
```

The output should look something like this:

```
=== DATABASE_URL
...
Add-on: <HEROKU_POSTGRES>
```

Lastly, in order to terminate all the provisioned infrastructure components, execute the command:

```
make destroy/deploy
```

### API endpoints

After deploying the changes, the API **endpoints** can be accessed from the hosted application using the following base URL:

```
https://icaroribeiro-templateapi.herokuapp.com
```

For example, in order to check the API documentation via web browser, access the following URL:

```
Method: HTTP GET
URL: https://icaroribeiro-templateapi.herokuapp.com/swagger/index.html
```

### Accessing remote Postgres database locally

It is possible to verify the data generated in Heroku using **pgAdmin** tool.

To achieve this, access the database Heroku web site in order to check the datastore settings.

In the Settings tab, click the View Credentials... button and take note of the following credentials: Host, Database, User, Port and Password.

After that, it is necessary to configure a remote server in pgAdmin by means of the values of the previous credentials of the Postgres database on Heroku as follows:

```
Host:
Database:
User:
Port:
Password:
```

In what follows, there are the steps to establish access to the Postgres database:

In pgAdmin, right click Server(s) icon, and then navigate to Create and Server options.

After that, it is necessary to fill out the following parameters:

In the General tab, name the server whatever you want.

Under the Connection tab, inform the Host name/address. It is the one configured like ...amazonaws.com. Keep the port at 5432. In Maintenance database, type the Database field from the credentials and do the same procedure to fill out the values of the Username and Password fields.

In the SSL tab, mark SSL mode as Require.

Before finalising it is necessary to apply one more configuration:

The database needs to be informed in a "desired database list" in order to avoid parsing many other databases that are not cared about. (This has to do with how Heroku configures their servers.)

In this regard, go to the Advanced tab and under DB restriction copy the Database name (it's the same value as the Maintenance database field filled earlier), and then click Save button.

Navigate through the options structure: Databases, database name, Schemas, public and inside Tables, finally, check the tables. (In case of the tables are not displayed, try right click the related Server created and then click Refresh option.)

**Note**

The project was configured with a Heroku Postgres database resource in a **Free plan** (Hobby Dev - Free). Because of that, the database will only support a limited number of records (10.000 rows). Therefore, please evaluate the operations to be carried out before using the application in this way.

## References

https://medium.com/wesionary-team/implement-database-transactions-with-repository-pattern-golang-gin-and-gorm-application-907517fd0743