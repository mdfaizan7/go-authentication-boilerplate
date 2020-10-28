# GO AUTHENTICATION BOILERPLATE
An Authentication Server boilerplate made with Go and PostgreSQL.

## Installation

### 1. Clone project
```bash
git clone https://github.com/mdfaizan7/go-authentication-boilerplate.git
```

### 2. cd into the folder
```bash
cd go-authentication-boilerplate
```

### 3. Start PostgreSQl server

### 4. Create database called go-auth
```bash
createdb go-auth # in bash/zsh

CREATE DATABASE go-auth; # in psql
```

### 3. Create env file
You need to start the server with right credentials stored inside .env file. <br />
You can do this with: `cp .env.example .env`

## Usage 
You can start the server with `go run main.go`. <br />
Then the server will start running on `http://localhost:3000`.

## Features
- Register
- Login
- Logout
- Cookies
- Access tokens and Refresh Tokens
- Authentication Middleware


