# go-event-manager-api

<p align="center">
  <img width="460" height="300" src="https://lets-go-further.alexedwards.net/static/img/gopher-gang.svg" alt="Gopher Gang">
</p>

A complete REST API for an event management service, developed in Go using the Gin framework.

## Features

  * **Event Management:** Create, list, update, and delete events.
  * **JWT Authentication:** Secure registration and login system with JWT tokens and password encryption using BCrypt.
  * **Attendee Management:** Add and remove attendees from events.
  * **SQLite Database:** Uses a lightweight and easy-to-configure database.
  * **Migrations:** Database version control to facilitate schema maintenance and evolution.
  * **Live Reload:** Real-time application updates during development with the Air tool.
  * **Swagger Documentation:** All endpoints are documented and can be tested through the Swagger interface.

##  Technologies Used

  * **[Golang](https://go.dev/)**: The main programming language.
  * **[Gin Framework](https://gin-gonic.com/)**: Web framework for building the API.
  * **[SQLite](https://www.sqlite.org/index.html)**: Relational database.
  * **[golang-migrate](https://github.com/golang-migrate/migrate)**: Tool for database version control (migrations).
  * **[Air](https://github.com/cosmtrek/air)**: Live-reloading tool for Go applications.
  * **[JWT (JSON Web Token)](https://jwt.io/)**: For token-based authentication and authorization.
  * **[BCrypt](https://en.wikipedia.org/wiki/Bcrypt)**: For hashing and securing passwords.
  * **[Swagger](https://swagger.io/)**: For documenting and testing the API endpoints.

## Getting Started

Follow the instructions below to set up and run the project in your local environment.

### Prerequisites

  * [Go](https://go.dev/doc/install) (version 1.25.1 or higher)
  * [Air](https://github.com/cosmtrek/air)
  * [golang-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

### Installation

1.  **Clone the repository:**

    ```bash
    git clone git@github.com:na1tto/event-manager.git
    cd event-manager
    ```

2.  **Install Go dependencies:**

    ```bash
    go mod tidy
    ```

3.  **Install `golang-migrate`:**

    ```bash
    go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    ```

4.  **Install `Air`:**

    ```bash
    go install github.com/cosmtrek/air@latest
    ```

### Environment Setup

1.  Create a `.env` file in the project root.

2.  Add the following environment variables to the `.env` file:

    ```env
    # Port the application will run on
    PORT=8080

    # Secret key for generating JWT tokens
    JWT_SECRET=your-secret-key-here
    ```

### Running the Migrations

To create the tables in the SQLite database, run the following commands:

  * **To apply the migrations (create tables):**

    ```bash
    go run cmd/migrate/main.go up
    ```

  * **To revert the migrations (drop tables):**

    ```bash
    go run cmd/migrate/main.go down
    ```

### Running the Application

With `Air` installed, you can start the application with the following command. `Air` will monitor file changes and restart the server automatically.

```bash
air
```

The API will be available at `http://localhost:8080`.

## API Documentation (Swagger)

The complete API documentation is available and can be accessed through the Swagger UI. After starting the application, navigate to the following URL in your browser:

[http://localhost:8080/swagger/index.html](https://www.google.com/search?q=http://localhost:8080/swagger/index.html)

There you will find all the endpoints, their parameters, and you can test them directly from the interface.

## Project Structure

```
.
├── cmd
│   ├── api         # Main API code (handlers, routes, main.go)
│   └── migrate     # Code to run migrations
├── docs            # Swagger documentation files
├── internal        # Business logic and database access
│   ├── database
│   └── env
├── .air.toml       # Air configuration file
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

##  Challenges and Learning

>One of the biggest challenges for me in this project was that I was not familiar with the Go style programming for web applications and some of the concepts that I worked with in this project where very new to me (e.g JWT Auth with password hashing, server and routes configurations, the Context struct for HTTP requests, and many others) and was really hard for me to understand them, I'd say that even now I'm not sure if I can dominate these concepts very well yet, so I'm certantly looking foward to consolidate and/or expand my knowledge in this development area.
>
>Nonetheless, this project taught me many things about the construction of REST APIs and their importance in the backend of web applications, I was able to connect the little knowledge that I had in Spring Boot with the context of Go development, which made the learning process a little easier. Other thing that was a learning point was the Swagger and Insomnia usage, today these are tools that I cannot work without. In general was a very interesting project and I'm very excited to learn new things in the Go and backend web development, it was really fun and I'm for sure doing a lot more with this language in the future!
>
>This project is my first experience working with Go in a more formal and structured way, I wanna thank the youtube channel **[Coding With Patrik](https://www.youtube.com/@codingwithpatrik)** for the tutorial of this project and the Go Reddit community for giving me the necessary guidance for executing this project, I know it's not the perfect and best structured API but it's for sure a nice begginning for the backend web development world, and was a big step for me to become the type of developer I want to be, I hope the knowledge that I acquired here can be helpful in future projects that I happen to be part of.

<p align="center">
  <img width="460" height="300" src=./Gopherr.png alt="Thanks and see you soon!">
</p>
