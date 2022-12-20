# AREA api-gateway

this is the api gateway documentation for the API Gateway module.

## Documentation:

---

### Backend Technologies:
- golang
- postgres
- ent to interact with the postgres database

---module.exports = nextConfig

### Routes:

#### Users

``/ping``
Route: **GET**

**Description:** Check if the API Gateway is available.

**header:** No header is needed.

**body:** No body is needed.

**Http code:** 200.

**Response:** string: "Pong".

---

``/register``
Route: **POST**

**Description:** Create a new user in the database and send the corresponding response.

**Header:** No header is needed.

**Body:** Json with the email of the user and his password.
```json
{
  "email": "user@example.com",
  "password": "example"
}
```

**Http code:**
- 200 on success.
- 400 (bad request) if the Json is not well formatted.
- 406 (status not acceptable) if the user already exists, or any error in the database.

**Response:** Json with the jwt of the user.
```json
{
  "token": "tokenExample"
}
```

---

``/updateUser``
Route: **POST**

**Description:** Update the user informations and send the corresponding response.

**Header:**
- Authorization: jwt of the user.

**Body:** Json with the new email and the new password of the user.
```json
{
  "email": "user@example.com",
  "password": "example"
}
```

**Http code:**
- 200 with the user in json format.
- 400 (bad request) if the Json is not well formatted.
- 406 (status not acceptable) if the user does not exists, or any error in the database.

**Response:**
Json with the jwt of the user.
```json
{
  "token": "tokenExample"
}
```

---

``/deleteUser``
Route: **DELETE**

**Description:** Delete a user and send the corresponding response.

**Header:**
- Authorization: jwt of the user.

**Body:** No body needed.

**Http code:**
- 200 with the user in json format.
- 400 (bad request) if the Json is not well formatted.
- 406 (status not acceptable) if the user has already been deleted.

**Response:**
string: "User {id} has been deleted successfully"

---

``/readUser``
Route: **POST**

**Description:** Read one user and send the corresponding response.

**Header:**
- Authorization: jwt of the user.

**Body:** No body needed.

**Http code:**
- 200 with the user in json format.
- 400 (bad request) if the Json is not well formatted.
- 406 (status not acceptable) if the user does not exists, or any error in the database.

**Response:**
the user's informations in Json format.
```json
{
  "id": "user_id",
  "email": "user@example.com",
  "password": "example"
}
```

---

``/login``

Route: **POST**

**Description:** Check if the user exist and send the corresponding response.

**Header:** No header needed.

**Body:** Admin access required (not implemented yet).

**Http code:**
- 200 with the user in json format.
- 406 (status not acceptable) if any error in the database.

**Body:** Json with the email and the password of the user.
```json
{
  "email": "user@example.com",
  "password": "example"
}
```

---

``/session``

Route: **GET**

**Description:** Check if the jwt is valid.

**Header:**
- Authorization: jwt of the user.

**Body:** Ano body required.

**Http code:**
- 202 (status accepted).
- 401 (status unauthorized) if the token is not valid
- 406 (status not acceptable) if the user is no longer in the database.

**Body:** No body.

---

#### Services

``/services``

Route: **POST**

**Description:** create the list of services and actions and reactions availables in the database.

**Header:** No header needed.

**Body:**
```json
{
  "client":{
    "host":"example"
  },
  "server":{
    "current_time":"current time",
    "services":[
      {
        "name":"facebook",
        "actions":[
            {
              "name":"new_message_in_group",
              "description":" A new message is posted in the group"
            },
            {
              "name":"new_message_inbox",
              " description ":"A new private message is received by the user"
            },
            {
              "name":"new_like",
              " description ":"The user gains a like from one of their messages"
            }
        ],
        "reactions":[
            {
              "name":"like_message",
              "description":"The user likes a message"
            }
        ]
      }
    ]
  }
}
```

**Http code:**
- 200 with the user in json format.
- 400 (bad request) if the Json is not well formatted.
- 406 (status not acceptable) if any error in the database.

**Response:** string: "Services created sucessfully"

---

``/services``

Route: **DELETE**

**Description:** Delete all the services in the database

**Header:** No header needed.

**Body:** No body needed.

**Http code:**
- 200 with the user in json format.
- 500 (internal server error) if any error in the database.

**Response:** string: "Services deleted sucessfully"

---

``/about.json``

Route: **GET**

**Description:** Retrieve all the list of services, actions and reactions.

**Header:**
- Authorization: jwt of the user.

**Body:** No body needed.

**Http code:**
- 200 with the user in json format.
- 400 (bad request) if the Json is not well formatted.
- 406 (status not acceptable) if any error in the database.

**Response:** string: "Services created sucessfully"

---

#### Areas

``/area``

Route: **POST**

**Description:** Create area in the database and send it to the dispatcher. It also send the acces_token of the Oauth services to the dispatcher.

**Header:**
- Authorization: jwt of the user.

**Body:**
```json
{
  "area_name": "My new AREA",
  "user_id": "user id",
  "action_reaction": [
    {
      "type": "action",
      "service": "google-drive",
      "name": "file_update",
      "order": 0,
      "params": {
        "filename": "fichier.docx"
      }
    },
    {
      "type": "reaction",
      "service": "google-drive",
      "name": "upload_file",
      "order": 1,
      "params": {
        "param1": "modification.docx"
      }
    }
  ]
}
```

**Http code:**
- 200 with the user in json format.
- 400 (bad request) if the Json is not well formatted.
- 406 (status not acceptable) if any error in the database.

**Response:** string: "Area created sucessfully"

---

``/area``

Route: **GET**

**Description:** Get a list of areas using the id of one user.

**Header:**
- Authorization: jwt of the user.

**Body:** No body needed.

**Http code:**
- 200 with the user in json format.
- 400 (bad request) if the Json is not well formatted.
- 406 (status not acceptable) if any error in the database.

**Response:** json:
```json
[
  {
    "id": 11,
    "area_id": "id of the Area",
    "user_id": "user id",
    "area_name": "My new AREA",
    "action_reaction": "[{\"type\":\"action\",\"services\":\"google-drive\",\"name\":\"file_update\",\"order\":0,\"payload\":{\"params\":{}}},{\"type\":\"reaction\",\"services\":\"google-drive\",\"name\":\"upload_file\",\"order\":1,\"payload\":{\"params\":{}}}]"
  }
]
```

---

``/area``

Route: **DELETE**

**Description:** Delete an Area using its name.

**Header:**
- Authorization: jwt of the user.

**Body:**
```json
{
  "area_name": "My new AREA"
}
```

**Http code:**
- 200 with the user in json format.
- 400 (bad request) if the Json is not well formatted.
- 406 (status not acceptable) if any error in the database.

**Response:** string: "Area '{area_name}' have been deleted successfully"

---

#### Oauth

``/Oauth/googleDrive/createUrl``

Route: **GET**

**Description:** Generate url for the frontend to connect.

**Header:**
- Authorization: jwt of the user.

**Body:** No body needed.

**Http code:**
- 200 with the user in json format.
- 500 (Internal server error) if any error in the database.

**Response:** string: "url of the service"

---

``/Oauth/googleDrive/sendToken``

Route: **POST**

**Description:** Send the google drive tokens to the dispatcher.

**Header:**
- Authorization: jwt of the user.

**Body:**
```json
{
  "code": "code to extact the google drive tokens"
}
```.

**Http code:**
- 200 with the user in json format.
- 500 (Internal server error) if any error in the database.

**Response:** json: google drive tokens sent to the dispatcher
