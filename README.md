# SnapSight - Captivating Photo Gallery Site

SnapSight is a captivating photo gallery site meticulously crafted using the power of Go 🐹.

## Objective

The primary objective of this project is to explore and demonstrate various aspects of Go programming language and web development, particularly highlighting the power of the standard library in building strong and robust applications. By building SnapSight, we aim to cover the following areas:

- **Web Development with Go:** Utilizing Go's concurrency, networking, and HTTP handling features to create a dynamic web application.

- **Routing and Middleware:** Leveraging the `github.com/go-chi/chi` router to handle routes and applying middleware for request processing.

- **Database Interaction:** Interacting with a PostgreSQL database using the `github.com/jackc/pgx` library and implementing migrations with `github.com/pressly/goose`.

- **Email Integration:** Incorporating email functionality using the `github.com/go-mail/mail` library for user communication.

- **Security Measures:** Implementing security practices such as Cross-Site Request Forgery (CSRF) protection through the `github.com/gorilla/csrf` package.

- **Environment Management:** Utilizing the `github.com/joho/godotenv` package to manage environment variables and configuration.

- **Exploring Cryptography:** Demonstrating cryptography capabilities offered by the `golang.org/x/crypto` library.

Through SnapSight, we aim to showcase how the Go programming language, with its strong standard library and performance advantages, can be a fantastic choice for building modern web applications.


## Libraries Used

- [github.com/go-chi/chi/v5](https://pkg.go.dev/github.com/go-chi/chi/v5) v5.0.10
- [github.com/go-mail/mail/v2](https://pkg.go.dev/github.com/go-mail/mail/v2) v2.3.0
- [github.com/gorilla/csrf](https://pkg.go.dev/github.com/gorilla/csrf) v1.7.1
- [github.com/jackc/pgconn](https://pkg.go.dev/github.com/jackc/pgconn) v1.14.1
- [github.com/jackc/pgerrcode](https://pkg.go.dev/github.com/jackc/pgerrcode) v0.0.0-20220416144525-469b46aa5efa
- [github.com/jackc/pgx/v4](https://pkg.go.dev/github.com/jackc/pgx/v4) v4.18.1
- [github.com/joho/godotenv](https://pkg.go.dev/github.com/joho/godotenv) v1.5.1
- [github.com/pressly/goose/v3](https://pkg.go.dev/github.com/pressly/goose/v3) v3.15.0
- [golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto) v0.12.0

## Getting Started

To get started with the SnapSight project, follow these steps:
1. Clone this repository:
```bash
git clone git@github.com:mobamoh/snapsight.git
```
2. Navigate to the project directory:
```bash
cd snapsight
```
3. Set up your environment variables using the .env file. Make sure to fill in the required values.
```bash
docker-compose up
```
4. Run the application:
```bash
go run main.go
```

## Contributing

Contributions are welcome! If you find any issues or want to add new features, feel free to open a pull request.

## Maintenance
- Main maintainer: [Mo Bamoh](https://github.com/mobamoh)