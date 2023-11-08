# Data Dictionary

**==modle.go design==**<br>
| Table    | Field             | Default Value  | Data Type   |
|----------|-------------------|----------------|-------------|
| User     | userID (PK)       | Auto increment |             |
| User     | username          | None           | Varchar     | 
| User     | password          | None           | Varchar     |

<br>

| Table    | Field              | Default Value  | Data Type  |
|----------|--------------------|----------------|------------|
| Note     | noteID (PK)        | Auto increment | Int        |
| Note     | userID             | Auto increment | Int        |
| Note     | noteTitle          | None           | Varchar    |
| Note     | noteContent        | None           | Text       |
| Note     | creationDate       | Today          | Timestamp  |
| Note     | DelegatedTo        | None           | String     |
| Note     | completionDate     | None           | Timestamp  |
| Note     | status             | none           | Varchar    |

<br>

| Table    | Field              | Default Value  | Data Type   |
|----------|--------------------|----------------|-------------|
| Sharing  | sharingID          | Auto Increment | Int         |
| Sharing  | noteID             | None           | Int         |
| Sharing  | userID             | None           | Int         |
| Sharing  | timestamp          | Today          | Timestamp   |             
| Sharing  | status             | Active         | Varchar     |

----------------------------------------------------------------

**==handler.go design==**<br>
| Table Name   | Field          | Data Type    | 
|--------------|----------------|--------------|
| Data         | userName       |string        |
| Data         | note           |[]DisplayNote |

<br>

| Table Name   | Field          | Data Type   |
|--------------|----------------|-------------|
| DisplayNote  | NoteTitle      | string      |
| DisplayNote  | CreateDate     | time.Time   |
| DisplayNote  | Delegation     | string      |
| DisplayNote  | CompletionDate | string      |
| DisplayNote  | Status         | time.Time   |
| DisplayNote  | Username       | string      |
| DisplayNote  | NoteContent    | string      |