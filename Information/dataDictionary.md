# Data Dictionary for all the tables
*In the modle.go design*

| Table    | Field             | Default Value  | Data Type   | Field Size  | Example                          |
|----------|-------------------|----------------|-------------|-------------|----------------------------------|
| User     | userID (PK)       | Auto increment | Int         |             | 1                                |
| User     | username          | None           | Varchar     | 100         | Mike                             |
| User     | password          | None           | Varchar     | 100         | Mike123                          |

| Table    | Field             | Default Value  | Data Type   | Field Size  | Example                                    |
|----------|-------------------|----------------|-------------|-------------|--------------------------------------------|
| Note     | noteID (PK)        | Auto increment | Int         |             | 101                                       |
| Note     | userID             | Auto increment | Int         |             | 102                                       |
| Note     | noteTitle          | None           | Varchar     | 255         | Things to do                              |
| Note     | noteContent        | None           | Text        |             | Have a coffee in the morning              |
| Note     | creationDate       | Today          | Timestamp   |             | 2023-10-02 10:00:30                       |
| Note     | DelegatedTo        | None           | String      | 100         | 2023-10-02 10:00:30                       |
| Note     | completionDate     | None           | Timestamp   |             | 2023-10-02 10:00:30                       |
| Note     | status             | none           | Varchar     | 50          | in progress/completed/cancelled/delegated |


| Table    | Field              | Default Value  | Data Type   | Field Size  | Example                         |
|----------|--------------------|----------------|-------------|-------------|---------------------------------|
| Sharing  | sharingID          | Auto Increment | Int         |             | 102                             |
| Sharing  | noteID             | None           | Int         |             | 102                             |
| Sharing  | userID             | None           | Int         |             | 101                             |
| Sharing  | timestamp          | Today          | Timestamp   |             | 2023-10-02 10:00:30             |
| Sharing  | status             | Active         | Varchar     | 50          | "Active" or "Revoked"           |


*In the handler.go design*

| Table Name   | Field          | Data Type    | 
|--------------|----------------|--------------|
| Data         | userName       |string        |
| Data         | note           |[]DisplayNote |


| Table Name   | Field          | Data Type   |
|--------------|----------------|-------------|
| DisplayNote  | NoteTitle      | string      |
| DisplayNote  | CreateDate     | time.Time   |
| DisplayNote  | Delegation     | string      |
| DisplayNote  | CompletionDate | string      |
| DisplayNote  | Status         | time.Time   |
| DisplayNote  | Username       | string      |