# new-go-code-challenge-template
New Go Code Challenge Template

https://icaroribeiro-newgocctmplapi.herokuapp.com/swagger/index.html


# Hi there!

Welcome to my solution to X's code challenge.

- [Introduction](#introduction)
- [Architecture](#architecture)
- [Database](#database)
- [pgAdmin](#pgadmin)
- [How to run the project?](#how-to-run-the-project)
- [API endpoints](#api-endpoints)
- [How to run the tests?](#how-to-run-the-tests)
- [Project Dynamics](#project-dynamics)
- [Deployment](#deployment)
- [References](#references)

## Introduction

This project comprehends the development of a **REST API** using **Go** programming language (Golang), **Json Web Token** and **Postgres** database for managing cryptocurrency purchases and sales based on operations related to users, accounts and transfers.

## Architecture

The project structure is designed using some concepts of the layered architecture of the **Domain Driven Design** (DDD) approach that is intended to simplify the complexity developers face by connecting the implementation to an evolving model.

To do it, the implementation is divided up into the following essential layers in order to have a separation of interests by arranging responsibilities:

### Application

This layer is responsible for serving the application purposes. It contains services that act as intermediaries between the database and the API requests/response. Also, it includes the handling of the third-party Coin Market API, the mechanism to verify if the service has started up correctly and is ready to accept requests and the validation of the input parameter values from the API requests payloads.

### Domain

This layer is resposible for holding domain business logic. It contains the schema of the "models" based on structs and properties used by both the API and the database actions.

### Infrastructure

This layer is responsible for serving as a supporting layer for other layers. It contains the procedures to establish connection to the database and the repositories to interact with the database by retrieving and/or modifing records. Also, it includes the operations intended to security such as password encryption and validation and the association of all the components of the layered architecture.

### Interfaces

This layer is responsible for the interaction with user by accepting API requests, calling out the relevant services and then delivering the response. It contains the handling of API requests, as well as the elaboration of API responses, the logging and authentication actions that mediate the access to the API **endpoints** and a router that exposes the routes associated with each one of them.

## Database and pgAdmin

Two Postgres dabases are used to handle the project. One of them is intended to be used in a common (or usual) scenario and the other is directed to a test scenario. However, both of them contain the same four tables named auths, users, accounts and transfers defined in the **database/scripts/1-create_tables.sql** file.

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

**Users**

The **users** table contains the users data.

| Fields          | Data type | Extra                     |
|:----------------|:----------|:--------------------------|
| id              | UUID      | NOT NULL PRIMARY KEY      |
| name            | TEXT      | NOT NULL                  |
| document_number | TEXT      | NOT NULL UNIQUE           |
| email           | TEXT      | NOT NULL UNIQUE           |
| password        | TEXT      |                           |
| created_at      | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP |
| updated_at      | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP |

**Accounts**

The **accounts** table contains the accounts data.

| Fields     | Data type | Extra                     |
|:-----------|:----------|:--------------------------|
| id         | UUID      | NOT NULL PRIMARY KEY      |
| user_id    | UUID      | NOT NULL FOREIGN KEY      |
| bitcoins   | NUMERIC   | NOT NULL DEFAULT 0        |
| balance    | NUMERIC   | NOT NULL                  |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP |

**Transfer**

The **transfers** table contains the transfers data.

| Fields             | Data type | Extra                     |
|:-------------------|:----------|:--------------------------|
| id                 | UUID      | NOT NULL PRIMARY KEY      |
| account_id         | UUID      | NOT NULL FOREIGN KEY      |
| type               | TEXT      | NOT NULL                  |
| bitcoins           | NUMERIC   | NOT NULL                  |
| bitcoin_unit_price | NUMERIC   | NOT NULL                  |
| quoted_at          | TIMESTAMP |                           |
| total_price        | NUMERIC   | NOT NULL                  |
| created_at         | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP |

**Notes**:

If the project is intended to be run locally, the database settings must be defined in the **scripts/setup_env.sh** file for both databases, respectivelly.

On the contrary, if the project is intended to be run with Docker containers, the database settings don't need to be changed and are defined in the **docker-compose.yml** file in the sections of the database services named **db** and **testdb** for the common and the test databases, respectivelly.

## How to run the project?

A **Makefile** was designed as a single point of entry containing a set of instructions to run the project in two different ways through commands from the terminal. 

### Run locally

To run the project locally, execute the command:

```
make run/api
```

**Notes**:

Before running it, it is necessary to configure the environment variables of the connection to the database already installed locally and the key of CoinMaketCap API in the **scripts/setup-env.sh** file. The other variables can be kept as defined if desired.

### Run with Docker containers

To run the project with Docker containers, execute the command:

```
make startup/docker
```

**Notes**:

Before running it, it is only necessary to configure the environment variable of the key of coinmaket API in the **.env** file. The other variables can be kept as defined if desired in the **.env** and the **docker-compose.yml** files.

## pgAdmin

Only to make checking data easier, a container with **pgAdmin** tool was also configured in the **docker-compose.yml** file. It is a GUI (Graphical User Interface) used to access and modify PostgreSQL databases that not only exist locally, but also remotely.

If it is necessary to build and run the pgadmin container, uncomment the related code in the **docker-compose.yml** file and execute the command:

```
docker-compose up -d pgadmin
```

After starting the container, access the following URL and provide the credentials defined in the **pgadmin** service section of the **docker-compose.yml** file:

```
http://{host}:5050
```

```
Username: pgadmin4@pgadmin.org
Password: admin
```

In what follows, there are the steps to establish access to the common and the test databases when running the project with Docker containers:

In pgAdmin, right click Server(s) icon, and then navigate to Create and Server options.

After that, it is necessary to fill out the following parameters:

In the General tab, name the server whatever you want. For instance, *PostgreSQL-db*.

Under the Connection tab, inform the Host name/address where the containers are running and set the port at *5433*. In Maintenance database, type *db* as the Database name and *postgres* as the values of the Username and Password fields and then click Save button. (The value of the port is defined in the *docker-compose.yml* file and the other values in *database/.env* file)

Navigate through the options structure: Databases, database name, Schemas, public and inside Tables, finally, check the tables. (In case of the tables are not displayed, try right click the related Server created and then click Refresh option.)

**Note**

The testdb was created only to be used by the execution of the integration tests execution and no data is going to be stored there. However, just to explain, the procedure above could also be used to configure the test database in pgAdmin if wanted with only a few changes. In case, the port must be set to 5434, the Database name must be defined testdb and the Username and Password must be equal to postgres. (The value of the port is defined in the *docker-compose.yml* file and the other values in *database/.testenv* file)

## API documentation

The API documentation was created using the **swaggo/swag** repository from Github that converts Golang annotations to Swagger Documentation 2.0 based on swagger files located in the **api** directory.

So, after running the application, access the following URL via web browser in order to view a HTML page that illustrates all API endpoints information:

```
http://{host}:8080/swagger/index.html
```

As other references on how to handle the API requests, in the  in **assets/requests** directory there are a few samples of API calls using cURL command line tool and a Postman collection with information of the related operations that can be imported in Postman tool for API testing.

**Notes**:

The ids requested by some API endpoints must follow UUID (Universally Unique Identifier) ​​standard.

## How to run the tests?

The test cases and the code under test were organized in separate packages to ensure testing only the exported identifiers of the packages. By doing this, the test code is compiled as a separate package and then linked and run with the main test binary.

The test cases were implemented as **Table Driven Tests** so that each test case was designed by declaring a structure that holds actions that can be performed before and after executing it in addition to inputs and expected outputs by following the **unit** test and **integration** test approaches.

The test cases analyze four layers of the project:

**Validator**: it contains test cases aimed to the validator layer. 

**Repository**: it contains test cases related to the repository layer.

**Service**: it contains test cases directed at the service layer.

**Handler**: it contains test cases associated with the handler layer.

The validator, repository, service and handler layers were evaluated using unit tests developed in the related directories whereas the integration tests were implemented in the *tests/integration* directory.

After running the tests, it is possible to check the percentage of code coverage served by each test case displayed in the output of the tests execution.

Furthermore, the statistics collected from the *unit tests* execution are saved in the *docs/tests/coverage.out* file for coverage analisys. However, this file doesn't contain any statistics from *integration tests* execution.

**Notes**:

The unit tests were written using fake objects to mock out dependencies so that the layers can interact with each other through interfaces rather than through concrete implementations.

Basically, the goal of mocking is to isolate and focus on the code being tested and not on the behaviour or state of external dependencies. In mocking, the dependencies are replaced by closely controlled replacements objects that simulate the behaviour of the real ones.

So, every layer is tested independently without having to be dependent on other layers and it is not necessary to be worried on the correctness of the dependency (the other layers.)

For the mocking purpose, it was used the *DATA-DOG/go-sqlmock* and *vektra/mockery* repositories for mocking the SQL driver behavior without needing to actually connect to a database and for generating the mock objects from interface, respectively.

On the other hand, the integration tests were written by combining and testing the layers together in order to simulate the production environment.

### Run locally

To run the tests locally, execute the command:

```
make test/api
``` 

To verify the coverage analysis, execute the command:

```
make analyze/api
```

**Notes**:

Before executing the tests, it is necessary to configure the environment variables of the connection to the test database already installed locally in the **scripts/setup-env.sh** file.

### Run with Docker containers

To run the tests with Docker containers, execute the command:

```
make test/docker
```

To verify the coverage analysis, execute the command:

```
make analyze/docker
```

**Notes**:

Before running the tests with Docker containers, it is not necessary to configure any environment variables. They can be kept as defined if desired in the **.env** and **docker-compose.yml** files. It is only necessary to verify if the **api_container** and **testdb_container** containers have started up successfully.

The document numbers informed in the tests cases (and in the files from **assets/requests/samples** directory) are used for testing purposes only. They were obtained from the tests developed in the **github.com/Nhanderu/brdoc** repository that is intended to validate such data as a Brazilian CPF (Cadastro de Pessoa Física) document.

## Project Dynamics

In what follows, there is a short guide including summary descriptions on how the solution works in practice.

As indicated in the API documentation, some operations are restricted because they require authentication. In other words, a token that belongs to a logged user must be sent in the authentication header of some API requests.

Below is a table informing all the operations and whether they require authentication or not:

| Lista de todas as operações     | A autenticação é necessária? |
|:--------------------------------|:-----------------------------|
| Health check                    | Não                          |
| Refresh token                   | Sim                          |
| Reset password                  | Sim                          |
| SignIn                          | Não                          |
| SignOut                         | Sim                          |
| SignUp                          | Não                          |
| Get all users                   | Sim                          |

Since not all operations can be performed freely, initially it is needed to create a user before exploring other features of the application.

A user can be related to one or more accounts and each account is associated with only one user.

In the process of creating a user there are some restrictions as follows:

A user cannot have the same email address or document number as another user already registered in the database.

After creating a user, it is necessary to perform login using his/her registered credentials.

If the login is successful, a token is returned and it can be used to authenticate other requests related to restricted operations.

In the process of creating an account there are some restrictions as follows:

An account must be related to a user registered in the database and its balance must be greater than zero.

Furthermore, an account can be related to one or more transfers and each transfer is associated with only one account.

In the process of creating a transfer there are some restrictions as follows:

A transfer must be related to an account registered in the database and its type must be equal to purchase or sale.

In the case of a purchase, the number of bitcoins to be purchased must be greater than zero and the account balance must be sufficient, that is, be a value greater than the total purchase amount.

On the other hand, in the case of a sale, the number of bitcoins to be sold must be greater than zero and the amount of bitcoins in the account must be sufficient, that is, be more than the amount of the sale.

### How are transfer transactions managed?

In a real world application, it is common to have to perform a transaction that combines some operations from several tables.

The process of creating a transfer of purchase or sale of currencies involves two steps. Firstly, it is attempted to create a transfer record in the transfers table. Later, it is attempted to update the account data related to the accounts table operation.

If for any reason there is an error along these two steps, the transfer creation process must be undone. For example, if an error occurs when updating the account data in the accounts table after the transfer record is created in the transfers table, that record must not be kept because the process failed. In this way, the creation of a transfer is successful when all the steps along the process are successfully executed.

Because of that, the transfer creation process is implemented within a transaction using the **Transaction** function from **"gorm.io/gorm"** package.

Basically, the following steps are performed:

- Initially, it is started a DB transaction with the **BEGIN** statement.
- After that, it is intended to execute two SQL queries. Firstly, a transfer record is to be created. Thereafter, the data of the account associated with the transfer is to be updated.
- If all of the operations are successful, a **COMMIT** is performed to make the transaction permanent and the database will be changed to a new state.
- Otherwise, if any query or even the **COMMIT** statement fails, a **ROLLBACK** is performed to revert all changes made by previous queries, and the database stays the same as it was before the transaction.

This procedure is to ensure the database integrity by providing a reliable and consistent unit of work that is made up of mutiple database operations, even in case of system failure. Thus, when using this approach, the database can be restored to some previous point after erroneous operations are performed.

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
heroku pg:psql -a=<HEROKU_APP_NAME> <HEROKU_POSTGRES> < database/scripts/1-create_tables.sql
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