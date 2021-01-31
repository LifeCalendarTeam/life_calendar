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
  -   `200 OK` for requests handled without errors

  -   `400 Bad Request` for requests with either an incorrect set of parameters or syntactically incorrect parameters'
      values (if, however, the values are syntactically correct, but are not logically correct, don't satisfy some
      invariants, etc, you should get a different error code, most likely `412 Precondition Failed` or something even more
      specific).

      **WARNING**: you are not guarantied to get such response if you send a request with parameters not documented for an
      API method. It is possible that extra parameters will just be ignored

  -   `401 Unauthorized` for requests sent without (or with incorrect/expired) cookies

  -   `405 Method Not Allowed` for requests with incorrect HTTP Method

  -   `500 Internal Server Error` for requests, which failed due to a server-side error

  Please, note, that if there were multiple problems while processing your request and there are multiple status codes
  applicable for your case, you can get a response with any of them (for example, if you send a bad request, and you are
  not authorized, you may get a response with any of the `400 Bad Request` and `401 Unauthorized` statuses). The same
  rule holds true for the statuses defined in the methods' docs: the statuses' priority order is implementation-defined
  unless stated otherwise.

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

-   `date` (`str`) - date of the day, formatted as `YYYY-MM-DD`

-   `activity_type` (`int`, optional, can be used multiple times) - `type_id`s of activities of the day

-   `activity_proportion` (`int`, optional, must be used the same number of times `activity_type` was used) -
    `proportion`s of activities

-   `emotion_type` (`int`, optional, can be used multiple times) - `type_id`s of emotions of the day

-   `emotion_proportion` (`int`, optional, must be used the same number of times `emotion_type` was used) -
    `proportion`s of emotions

**WARNING**: if you exchange activities and emotions (i.e. send an activity (both type and proportion) as it was an
emotion), this mistake is silently ignored, and the values are stored in the database as if your request was correct
(note that the activities' and emotions' types numbering is end-to-end, so the server always knows whether something is
an activity or an emotion).

#### Response

Response will contain a `json` of the following scheme: `{'id': <int>}`, where `id` is the identifier of the created day
if the request had finished successfully. The scheme will instead be either `{'error_type': <str>}` or
`{'error_type': <str>, 'bad_ae_type': <str>}` if one of the errors described right above occurred. The meanings of
`error_type` and `bad_ae_type` are described below. Note that the expected activities/emotions types are `int`s, but
the returned `bad_ae_type` is `str`, to handle cases when the given type is not a number

If there is a day with the specified `date` already, you will get an error response with the `412 Precondition Failed`
status code. If the number of entries of either `activity_type` and `activity_proportion` or `emotion_type` and
`emotion_proportion` differ, you will get an error response with the `400 Bad Request` status code. If any of the types
or proportions (of either of activities or emotions) is not a number, or any of the proportions is not an integer in the
\[0; 100\] range, you will get a response with the `400 Bad Request` status code. If any of the types is incorrect (i.e.
there is **neither** an activity, nor an emotion with that type for the user sending the request), you will get a
response with the `412 Precondition Failed` status code. If there is a type of activity/emotion, which was mentioned
more than once in the request (including the case when the same type is mentioned once as an activity and once as an
emotion), you will get a response with the `400 Bad Request` status code.

Unless the request is incorrect because of one of the errors described in the beginning of the docs (e.g. you are making
a request without authorization, there is an internal server error, etc), if one of the above errors occur, the json
response will contain the `error_type` field (`string`), value of which would one of the following:
`types_and_proportions_lengths`, `incorrect_date`, `incorrect_type`, `duplicated_type`, `incorrect_proportion`,
`day_already_exists`. It can be used for understanding what exactly went wrong. Note, that `error_type` never tells you
if the error was because of a problem in activities or emotions: it only tells whether it was in the types or in the
proportions. Also note, that though there is the `error_type` field added, the usual field `error` (described in the
beginning of the docs and containing a human-readable error) is not removed from responses of the method.

If there is the `error_type` field in the response, and it is either of `incorrect_type`, `duplicated_type` or
`incorrect_proportion`, the response will also contain the `bad_ae_type` ("ae" stands for activity/emotion) key,
value of which is the type of activity/emotion which caused the error described in the response.

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
