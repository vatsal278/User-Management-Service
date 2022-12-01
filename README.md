# User Management Service

[![Build](https://github.com/vatsal278/UserManagementService/actions/workflows/build.yml/badge.svg)](https://github.com/vatsal278/UserManagementService/actions/workflows/build.yml) [![Test Cases](https://github.com/vatsal278/UserManagementService/actions/workflows/test.yml/badge.svg)](https://github.com/vatsal278/UserManagementService/actions/workflows/test.yml) [![Codecov](https://codecov.io/gh/vatsal278/UserManagementService/branch/main/graph/badge.svg)](https://codecov.io/gh/vatsal278/UserManagementService)

The application will be split up into multiple components, each having a particular feature and use case. This will allow individual scale up/down and can be started up as micro-services.

HTTP calls will be made across micro-services. 

They may be made asynchronous & de-coupled via pub-sub or messaging queues.

*For testing individual services, these can be via direct HTTP calls*


All requests & responses must follow json encoding.
Requests can be specific to the the concerned endpoint while responses must be of the following json format & specification:
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

## User Management Service Endpoints:

### Registration
 A first time user hits this endpoint to create a new account. The account will be created in a Relational DB with `email` as the `primary key`, `company_name`, `name` & `password` as `varchar`, `registered_on` & `updated_on` as `timestamp`, `active` as `boolean`, and `active_devices` as `integer`.

 - The password must be hashed with a salt before storing in the database. This salt can be additionally stored in the database table as plain text along with the user details.
 - Additionally generate a user-id of UUID format and add it to a `user_id` column in the database.
 - Once a new user record in inserted into the database, after proper validation, trigger a call to the account management service to initiate an account creation for the new user, via the message queue, and keep the active column as false.

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

 - If the active column is false, respond with HTTP 202 with “data” as “`Account activation in progress`”.
 - When credentials match, generate a `jwt` token with `user_id`, and set it in a cookie as part of the response and update the number of `active_devices` by 1. This cookie will be used by the endpoints to identify the session of the logged in user.

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
Failure cases should have the same response structure with appropriate status code and message.

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
Failure cases should have the same response structure with appropriate status code and message.

### User Details
This gets the user details for the logged in user. If the user is not logged in, respond with `HTTP 401` as per the response structure. The logged in user can be validated by the presence of  the cookie, after successful login.
- The extraction of the user_id from the cookie must be done in a middleware and this `user_id`, after successful extraction, should be put in the context for downstream handlers to get this data.
- Use the user_id in the cookies’ jwt to fetch the details from the DB and return in the response.
- If there are no user records with the token, respond with `HTTP 403` as per the response structure.

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
Failure cases should have the same response structure with appropriate status code and message.

### Middleware(s):
1. ExtractUser: extracts the user_id from the cookie passed in the request and forwards it in the context for downstream processing. If the cookie is not set, respond with status code `401`, and with `500` for other internal errors.
2. ScreenRequest: allows requests only from the message queue to be passed downstream. The middleware checks the “`user-agent`” & request `URL` to identify requests originating from the message queue.
*The URL(s) of the message queue(s) can be passed as a configuration to the service to allow requests only from URLs in the list*.

### References
- [Properly store passwords in DB](https://www.youtube.com/watch?v=zt8Cocdy15c)
