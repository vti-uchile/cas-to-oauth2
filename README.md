# CAS to OAuth2

## Description

"CAS to OAuth2" is a project that implements an authentication interface designed to translate between the Central Authentication Service (CAS v2.0) protocol and OAuth 2.0 with OpenID. This project facilitates the integration of systems using CAS with an OAuth2 server, providing a robust and versatile solution for managing authentications and authorizations in complex environments.

## Technologies Used

 - Go: The primary language used for backend development.
 - Gin: A web framework used in Go for handling HTTP requests.
 - MongoDB: The database used for storing relevant information.

## Requirements

 - Install Go v1.21.1.
 - Install and configure MongoDB.

## Installation and Configuration

1. Clone the repository:

```bash
git clone https://github.com/vti-uchile/cas-to-oauth2.git
cd cas-to-oauth2
```

2. Install dependencies:

```bash
go mod tidy
```

3. Environment setup:

 - Move or copy the .env.local file to .env.
 - Complete the .env file with appropriate data (MongoDB connection, OAuth2 service URL, cookie name, etc.).

## Running the Project

Once the setup is complete, run the project with:

```bash
go run cmd/main.go
```

## Usage

The project is used as a standard CAS server. For a quick verification that everything is working correctly, you can run the unit tests available in the "tests" directory.

## Additional Configurations

> If the USE_APM variable in the .env file is set to true, you should also configure the following variables: ELASTIC_APM_SERVICE_NAME, ELASTIC_APM_SERVER_URL, ELASTIC_APM_SECRET_TOKEN, and ELASTIC_APM_ENVIRONMENT.
