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


### `/api/days/brief`
Get brief info about days, filled by the user you are currently logged in
<!-- TODO: probably add some parameters limiting which days we want to retrieve (probably, dates range?) -->

#### Request
Send a `GET` request

#### Response
Response will contain a `json` of the following scheme: `{'days': [{'id': <int>, 'date': <str>, 'average_color':
[<int>, <int>, <int>]}, ...]}`, where `id` is an id of a day (you can get more info about a particular day by its id),
`date` is the date of a day (its format is `%Y-%M-%DT00:00:00Z`) and `average_color` is an RGB color calculated as a
proportional average between colors of emotions and activities present at a certain day.
<!-- TODO: Remove the `T00:00:00Z` part from the date format -->
 
 
### `/api/days/<id: int>`
Get a list of emotions/activities `type_id`s and `proportion`s
<!-- TODO: probably add a parameter to specify if we want to retrieve only activities or only emotions -->

#### Request
Send a `GET` request

#### Response
Response will contain a `json` of the following scheme: `{'emotions': [{'type_id': <int>, 'proportion': <int>}, ...],
'activities': [{'type_id': <int>, 'proportion': <int>}, ...]}`, where `type_id`s correspond to a certain
emotion/activity, which has a name and a color. You can retrieve more info about an emotion/activity by its type id.

If there is no day with the `id` identifier, you will get an error response with the `404 Not Found` status code


### `/activities`
Either get all activities or add a new activity

#### Request
Depending on what you want to do:
- Send a `GET` request
- Send a `POST` request with `Content-Type: application/x-www-form-urlencoded` and parameters:
    - `name` (`str`) - the activity's name/label
    - `color` (array of 3 `int`s) - an RGB color, corresponding to this activity
    - `is_everyday` (`bool`) - whether the activity should be suggested for every new day or its one-time
    
#### Response
For `GET`: response will contain a `json` of the following scheme: `{'data': [{'label': <str>, 'color': 
[<int>, <int>, <int>], 'is_everyday': true}, ...]}`
    
For `POST`: response will contain a `json` of the following scheme `{'type_id': <int>}`. You can use it to describe new
days.


### `/emotions`
Exactly like `/activities` (but for emotions)


### `/activities/<type_id: int>`
Get info about the activity type by its id. Will return the following information:
- `name` - the activity's name/label
- `color` - an RGB color, corresponding to this activity
- `is_everyday` - whether the activity should be suggested for every new day or its one-time

#### Request
Send a `GET` request

#### Response
Response will contain a `json` of the following scheme: `{'label': <str>, 'color': [<int>, <int>, <int>],
'is_everyday': true}`

If there is no activity with the `id` identifier, you will get an error response with the `404 Not Found` status code


### `/emotions/<type_id: int>`
Exactly like `/activities/<type_id: int>` (but for emotions)
