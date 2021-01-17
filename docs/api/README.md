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

- Returns a relevant HTTP status code. Common status codes for all methods are:
    - `200 OK` for requests handled without errors
    - `400 Bad Request` for requests with an incorrect set of parameters
    - `401 Unauthorized` for requests sent without (or with incorrect/expired) cookies
    - `405 Method Not Allowed` for requests with incorrect HTTP Method
    - `500 Internal Server Error` for requests, which failed due to a server-side error

Note that the documentation may not repeat these status codes and return-value keys, because they are common for all the
methods below, unless otherwise stated


### `/login`
**This method is an exception!**

As you can see, its URL does not start with `/api/`. It also does not (ever) return `json`, but it submits to the rules
of HTTP status code, so you can understand if your request was successful based on it.

This method creates and sets cookies which you must send to all the other methods to get a valid response. Please,
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
  This method does not (ever) return a response with the `401 Unauthorized` status code. All the others are possible


### `/activities`
Either get all activities or add a new activity

#### Request
Depending on what you want to do:
- Send a `GET` request
- Send a `POST` request with `Content-Type: application/x-www-form-urlencoded` and parameters:
    - `name` (`str`) - the activity's name/label
    - `color` (`str`) - an RGB color with the format `<int>,<int>,<int>`, corresponding to this activity
    - `is_everyday` (`bool`) - whether the activity should be suggested for every new day or its one-time

#### Response
For `GET`: response will contain a `json` of the following scheme: `{'data': [{'label': <str>, 'color':
[<int>, <int>, <int>], 'is_everyday': true}, ...]}`

For `POST`: response will contain a `json` of the following scheme `{'type_id': <int>}`. You can use it to describe new
days.


### `/emotions`
Exactly like `/activities` (but for emotions)


### `/activities/<type_id: int>`
Either get info about the activity type by its id or update an activity

#### Request
Depending on what you want to do:

- Send a `GET` request

- Send a `PATCH` request with any of the following parameters (each parameter should be used either 0 or 1 times):
    - `name` (`str`) - the activity's name/label
    - `color` (`str`) - an RGB color with the format `<int>,<int>,<int>`, corresponding to this activity
    - `is_everyday` (`bool`) - whether the activity should be suggested for every new day or its one-time

  If you send an empty `PATCH` request, you will get an error response with the `400 Bad Request` status code

#### Response
For `GET`: response will contain a `json` of the following scheme: `{'label': <str>, 'color': [<int>, <int>, <int>],
'is_everyday': true}`

For `PATCH`: you will get a response with the `200 OK` status code

If there is no activity with the `id` identifier, you will get an error response with the `404 Not Found` status code


### `/emotions/<type_id: int>`
Exactly like `/activities/<type_id: int>` (but for emotions)


### `/api/days/brief`
Get brief info about days, filled by the user you are currently logged in
<!-- TODO: probably add some parameters limiting which days we want to retrieve (probably, dates range?) -->
<!-- TODO: probably merge `/api/days/brief` with `/api/days`? I.e. make one endpoint with GET and POST methods -->

#### Request
Send a `GET` request

#### Response
Response will contain a `json` of the following scheme: `{'days': [{'id': <int>, 'date': <str>, 'average_color':
[<int>, <int>, <int>]}, ...]}`, where `id` is an id of a day (you can get more info about a particular day by its id),
`date` is the date of a day (its format is `%Y-%M-%DT00:00:00Z`) and `average_color` is an RGB color calculated as a
proportional average between colors of emotions and activities present at a certain day.
<!-- TODO: Remove the `T00:00:00Z` part from the date format -->


### `/api/days`
Add a new day

#### Request
Send a `POST` request with `Content-Type: application/x-www-form-urlencoded` with the following parameters:
- `date` (`str`) - date of the day, formatted as `%Y-%M-%DT00:00:00Z`

- `activity_type` (`int`, optional, can be used multiple times) - `type_id`s of activities of the day

- `activity_proportion` (`int`, optional, must be used the same number of times `activity_type` was used) -
  `proportion`s of activities

- `emotion_type` (`int`, optional, can be used multiple times) - `type_id`s of emotions of the day

- `emotion_proportion` (`int`, optional, must be used the same number of times `emotion_type` was used) -
  `proportion`s of emotions

If there is a day with the specified `date` already, you will get an error response with the `400 Bad Request` status
code. If the numbers of entries of either `activity_type` and `activity_proportion` or `emotion_type` and
`emotion_proportion` differ, you will get an error response with the `400 Bad Request` status code.

#### Response
Response will contain a `json` of the following scheme: `{'id': <int>}`, where `id` is the identifier of the created day


### `/api/days/<id: int>`
Either get a list of emotions/activities `type_id`s and `proportion`s of a day by its id or delete a day
<!-- TODO: probably add a parameter to specify if we want to retrieve only activities or only emotions -->
<!-- TODO: add day update mechanism instead of requiring to delete and add a day -->

#### Request
To get info, send a `GET` request. To delete a day, send a `DELETE` request.

#### Response
For `GET`: response will contain a `json` of the following scheme: `{'emotions': [{'type_id': <int>, 'proportion':
<int>}, ...], 'activities': [{'type_id': <int>, 'proportion': <int>}, ...]}`, where `type_id` correspond to a certain
emotion/activity, which has a name and a color. You can retrieve more info about an emotion/activity by its type id.

For `DELETE`: you will get a response with the `200 OK` status code

If there is no day with the `id` identifier, or it is a day of another user, you will get an error response with the
 `404 Not Found` status code
