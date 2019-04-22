# REST API of Go-Contacts

The following is the REST API of Go-Contacts. The `$SITEROOT` is the configured path of the API. For example, the default path is `/go-demo` and the address to access the sever is `http://localhost:8080`, then any request to to `$SITEROOT/...` means a request to `http://localhost:8080/go-demo/...`

In all the following request, if the response `status` is `true`, the request is successful. Otherwise, the request fails, and the `message` indicates the error message.

## Register

* `POST $SITEROOT/api/register`
* `Content-type: application/json`

The request body looks like this (email has to be unique):

```
{
    "username": "user",
    "email": "user@example.com",
    "password": "123456"
}
```

If successful, the response returns something like

```
{
    "account": {
            "ID": 1,
            "Username": "user",
            "email": "user@example.com",
            "password": "",
            "token": "The_Authorization_Token"
        },
    "status": true,
    "message": "Account has been created"
}
```

Note that `account.password` is always empty. `account.token` will be used in the successive requests, including the GraphQL requests.

## Login

* `POST $SITEROOT/api/login`
* `Content-type: application/json`

The request body looks like this:

```
{
    "email": "user@example.com",
    "password": "123456"
}
```

If successful, the response returns something similar to the response for `$SITEROOT/api/register`, and the `account.token` will be used for any successive requests

```
{
    "account": {
            "ID": 1,
            "Username": "user",
            "email": "user@example.com",
            "password": "",
            "token": "The_Authorization_Token"
        },
    "message": "Logged In",
    "status": true
}
```

## Add a Contact

* `POST $SITEROOT/api/contact/add`
* `Content-type: application/json`
* `Authorization: Bearer The_Authorization_Token`

The request body looks like this:

```
{
    "name": "friend's name",
    "phone": "(123) 456-7890",
    "address": "On the street",
    "email": "friend@example.com"
}
```

If successful, the response returns something like this

```
{
    "contact": {
            "ID": 1,
            "Name": "friend's name",
            "Phone": "(123) 456-7890",
            "Address": "On the street",
            "Email": "friend@example.com"
        }
    "message": "Inserted",
    "status": true
}
```

## Delete a Contact

* `DELETE $SITEROOT/api/contact/delete`
* `Content-type: application/json`
* `Authorization: Bearer The_Authorization_Token`

The request body looks like this (only an integer indicates the ID of the contact to be deleted)

```
1
```

If successful, the response returns something like this

```
{
    "message": "Deleted",
    "status": true
}
```

N.B., unless there's some internal error, even the user doesn't have such a contact, it is considered as a success. This is designed on purpose to perplex mischievous people.

## Update a Contact

* `PUT $SITEROOT/api/contact/update`
* `Content-type: application/json`
* `Authorization: Bearer The_Authorization_Token`

The request body looks like this (`id` is the contact ID, and is required), and *any missing field will be cleared to be empty*: 

```
{
    "id": 1,
    "name": "Friend's New Name",
    "email": "friend.new@example.com",
    "address": "friend's new address",
    "phone": "(123) 789-4567"
}
```

On success, the response returns something like this

```
{
    "message": "Updated",
    "status": true
}
```

N.B., and again, unless there's an internal error, the response is always "successful" to perplex mischievous people

## Get all Contacts

* `GET $SITEROOT/api/all-contacts`
* `Content-type: application/json`
* `Authorization: Bearer The_Authorization_Token`

On success, the response returns something like this

```
{
    "all-contacts": [
        {
            "ID": 1,
            "Name": "Friend 1",
            "Phone": "",
            "Address": "Address 1",
            "Email": "friend-1@example.com"
        },
        {
            "ID": 2,
            "Name": "Friend 2",
            "Phone": "",
            "Address": "Address 2",
            "Email": "friend-2@example.com"
        },
        ...
    ],
    "message": "Fetched all contacts",
    "status": true
}
```
