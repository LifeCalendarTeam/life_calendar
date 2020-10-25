# API Docs
This document describes how to use the Life Calendar API


## Introduction
Life Calendar provides a REST HTTP API. You must authorize to use API methods. The authorization data is stored in
cookies, so you'll have to send the cookies you got when authorizing to all the methods you use.


## Methods
All the API methods (except `/login`):

- Are accessible via an endpoint, which starts with `/api/`

- Always return `json` (i.e. `Content-Type: application/json`) (which always contains field `ok` (`bool`, `true` if the
request was successful, `false` otherwise) and if `ok` is `false`, also contain field `error` (`str`), set to a
human-readable explanation of the error occurred while processing your request)

- Returns a relevant HTTP status code. Common status codes for all methods are `200 OK` for requests handled without
errors, `400 Bad Request` for requests with an incorrect set of parameters, `401 Unauthorized` for requests sent without
 cookies, `405 Method Not Allowed` for requests with incorrect HTTP Method and `500 Internal Server Error` for requests,
 which failed due to a server-side error.

Note that for the following methods only json key-value pairs and status codes different from the above are described
(i.e. `ok` json key is always present, any request can return `200` or `500`, so the documentation doesn't repeat this
for each of the methods)


### `/login`
**This method is an exception!**

As you can see, its URL does not start with `/api/`. It also does not (ever) return `json`, but it submits to the rules
of HTTP status code, so you can understand if your request was successful based on it.

This method also creates and sets cookies which you must send to all the other methods to get a valid response. Please,
note that cookies expire in 24 hours (and they will possibly be able to expire on user's request), so remember to renew
your cookies with a request to this method.

#### Request
Send a `POST` request with `Content-Type: application/x-www-form-urlencoded` with parameters:
- `user_id` (`int`) - the identifier of the user you want to log in under
- `password` (`str`) - the password of that user

#### Response
A response with `Content-Type: text/html` is returned. Possible status codes:
- `200 OK` - both `user_id` and `password` are correct. You have successfully authorized.
- `403 Forbidden` - wrong `user_id` or `password`. You didn't authorize, no cookies were created.
