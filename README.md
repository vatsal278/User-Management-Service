# UserManagementService

[![Build](https://github.com/vatsal278/UserManagementService/actions/workflows/build.yml/badge.svg)](https://github.com/vatsal278/UserManagementService/actions/workflows/build.yml) [![Test Cases](https://github.com/vatsal278/UserManagementService/actions/workflows/test.yml/badge.svg)](https://github.com/vatsal278/UserManagementService/actions/workflows/test.yml) [![Codecov](https://codecov.io/gh/vatsal278/UserManagementService/branch/main/graph/badge.svg)](https://codecov.io/gh/vatsal278/UserManagementService)

* This service was created using Golang.
* This service has used clean code principle and appropriate go directory structure.
* This service is completely unit tested and all the errors have been handled.
* This service utilises messageBroker service for communicating with other micro services.

## Starting the UserManagementService

* Start the Docker container for mysql with command :
```
docker run --rm --env MYSQL_ROOT_PASSWORD=pass --env MYSQL_DATABASE=usermgmt --publish 9095:3306 --name mysql -d mysql
```
* Start the MsgBroker service using steps as described in the [link](https://github.com/vatsal278/msgbroker)

 
* Start the Api locally with command : 
```
go run .\cmd\UserManagementService\main.go
```
### You can test the api using post man, just import the [collection](https://github.com/vatsal278/UserManagementService/blob/4eb499337edbb738b524cf05ac44a4da362875ce/docs/user%20management%20svc.postman_collection.json) into your postman app.
### To check the code coverage
```
cd docs
go tool cover -html=coverage
```
## User Management Service:

This application is split up into multiple components, each having a particular feature and use case. This will allow individual scale up/down and can be started up as micro-services.

HTTP calls are made across micro-services.

They are made asynchronous & de-coupled via pub-sub or messaging queues.

*For testing individual services, these can be via direct HTTP calls*


All requests & responses follow json encoding.
Requests are specific to the concerned endpoint while responses are of the following json format & specification:
>
>    Response Header: HTTP code
>
>    Response Body(json):
>    ```json
>    {
>       "status": <HTTP status code>,
>       "message": "<message>",
>       "data": {
>        // object to contain the appropriate response data per API
>       }
>    }
>    ```

## User Management Service:

### Registration
A first time user hits this endpoint to create a new account. The account will be created in a Relational DB with `email` as the `primary key`, `company_name`, `name` & `password` as `varchar`, `registered_on` & `updated_on` as `timestamp`, `active` as `boolean`, and `active_devices` as `integer`.

#### Specification:
Method: `POST`

Path: `/register`

Request Body:
```json
{
  "name": "<Full Name>",
  "registration_date": "<DD-MM-YYYY HH:MM:SS format>",
  "email": "<proper email>",
  "password": "<minimum 8 characters with at least 1 upper case, 1 lower      case & 1 special character out of[,.@$?]>"
}
```
Success to follow response as specified:

Response Header: HTTP 201

Response Body(json):
```json
{
  "status": 201,
  "message": "SUCCESS",
  "data": "Account activation in progress"
}
```

### Login
A registered user whose account is activated(active column for user is true), uses this endpoint to login to their account via email & password combination.

#### Specification:
Method: `POST`

Path: `/login`

Request Body:
```json
{
  "email": "<proper email>",
  "password": "<minimum 8 characters with at least 1 upper case, 1 lower      case & 1 special character out of[,.@$?]>"
}
```
Success to follow response as specified:

Response Header: HTTP 200

Response Body(json):
```json
{
  "status": 200,
  "message": "SUCCESS",
  "data": nil
}
```

### Activation
This endpoint will be used to set the active status true for the user, thus allowing the user to log into the dashboard. This endpoint will only allow requests that are directly obtained by the message queue via middleware.

#### Specification:
Method: `PUT`

Path: `/activate`

Request Body:
```json
{
  "user_id": "<user_id for the record to activate>"
}
```
Success to follow response as specified:

Response Header: HTTP 200

Response Body(json):
```json
{
  "status": 200,
  "message": "SUCCESS",
  "data": nil
}
```

### User Details
This gets the user details for the logged in user. The logged in user can be validated by the presence of  the cookie, after successful login.
- The extraction of the user_id from the cookie is done in a middleware and this `user_id`, after successful extraction, is put in the context for downstream handlers to get this data.

#### Specification:
Method: `GET`

Path: `/user`

Request Body: `not required.`

Success to follow response as specified:

Response Header: HTTP 200

Response Body(json):
```json
{
  "status": 200,
  "message": "SUCCESS",
  "data": {
    "name": "<full name>",
     // if email was abcde123@gmail.com
    "email": "<masked email -> abxxxxx23@xxx.com>",
    "company": "<company name>",
    "last_login": "<timestamp>"
  }
}
```

### Middleware(s):
1. ExtractUser: extracts the user_id from the cookie passed in the request and forwards it in the context for downstream processing.
2. ScreenRequest: allows requests only from the message queue to be passed downstream. The middleware checks the “`user-agent`” & request `URL` to identify requests originating from the message queue.
   *The URL(s) of the message queue(s) is passed as a configuration to the service to allow requests only from URLs in the list*.


