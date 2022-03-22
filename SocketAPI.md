# webChat socket api

## Status Code

* LOGIN
* LOGOUT
*	REGISTER
*	CHATMESSAGE
*	GROUPMESSAGE
*	ERROR

increment the number from zero

## Socket Steam

### Message API

```go
type clientMessageAPI struct {
	Type    int    `json:"type"`
	Content []interface{} `json:"content"`
}
```

Login example:

```json
{
  "type": 0,
  "content": [
    {
    },
    {
    }
  ]
}
```

### Content API

Based on `chatServer/model.go`

```go
type userID int64
type groupID int64

type User struct {
	ID       userID
	Name     string
	Salt     string
	Password string
}

type Group struct {
	groupID   int64
	leader    userID
	adminList []userID
	userList  []userID
	count     int
}

type Message struct {
	SenderID   userID
	ReceiverID userID
	Content    string
	Created    time.Time
	Modified   time.Time
}

type GroupMessage struct {
	GroupID  groupID
	UserID   userID
	Content  string
	Created  time.Time
	Modified time.Time
}
```

Login example:

```json
{
  "type": 0,
  "content": [{
    "id": 1234567890,
    "password": "123456"
  }]
}
```

Chat messages example:

```json
{
  "type": 3,
  "content": [
    {
      "SenderID": 123456789,
      "ReceiverID": 987654321,
      "Content": "Hello, my bro",
      "Create": "",
      "Modified": ""
    },
    {
      "SenderID": 123456789,
      "ReceiverID": 987654321,
      "Content": "Nice to meet u",
      "Create": "",
      "Modified": ""
    }
  ]
}
```
