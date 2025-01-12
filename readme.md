# NEXShop

NEXShop is a simple restful api for managing e-commerce. This application makes it easier for users if they want to buy some product.

## Table of Contents

1. [About](#about)
   - [Features](#features)
   - [Technologies](#Technologies)
2. [Start](#start)
   - [Prerequisite](#Prerequisite)
   - [Installation](#Installation)
   - [Configuration](#Configuration)
   - [Run](#Run)
   - [Directory Structure](#Directory Structure)
3. [Contact](#Contact)

## About

Nexshop was built with the aim of making it easier for users to buy some product. This api is made using Gin-Gon and the database uses PostgreSQL.

### Features

- CRUD User, Product, Cart, Order
- Authentication With JWT
- Hash Password

### Technologies

- Gin Gonic
- Golang
- PostgreSQL

## Start

### Prerequisite

To get started, you need to have Golang installed on your system. If it's not installed yet, download and install it from the official Golang website.

### Installation

1. Ensure the Database is Available
   - Make sure you have created a PostgreSQL database named **`nexshop_db`**.

2. Run the Available Queries
   - Open the schema.sql file in this application.
   - Copy and execute it in a database client like DBeaver or others.

3. Clone the repository

```sh
$ git clone https://github.com/RamadanRangkuti/NexShop.git
```

4. Go to folder repository

```sh
$ cd [name]
```

5. Open folder

```sh
$ code .
```

6. Download the dependencies:

```sh
$ go mod tidy
```

### Configuration

The project uses a .env file for environment variables like database connection details, server etc.
you can create a .env file according to the .env.example in the root directory

### Run

Run the following command to start the server:

```sh
$ go run cmd/main.go
```


### Directory Structure

- **cmd/**: Entry point of the application.
- **internal/**:  Business logic for the application.
- **pkg/**: Stores configurations related to third-party services like PostgreSQL, JWT, and others.
- **internal/handlers/**:  A layer that handles requests from users, whether from mobile or web applications.
- **migration/**: Stores SQL migration files and code to manage and update database structures (coming soon).
- **internal/models/**: Contains data structures (constructs in Golang) that simplify creating contracts for requests and responses.
- **internal/respositories/**:  A dedicated layer for interacting with the database, including data recording and retrieval operations.
- **internal/routes/**: tores the main endpoint definitions that direct requests.
- **.env**: he primary configuration file to set up database connections and other application parameters.


## Contact

Ramadan Rankguti - ramadanrangkuti17@gmail.com
