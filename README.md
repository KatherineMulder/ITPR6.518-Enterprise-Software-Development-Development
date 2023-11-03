## ITPR6.518-Enterprise-Software-Development-Development
_Katherine Mulder & Alex Borawski
Eastern Institute of Technology
NZ Bachelor of Computing Systems
John Jamieson
October 31, 2023_ 

#### INTRODUCTION 
The Enterprise Notes application is geared towards the needs of enterprise users, who can set up their accounts within the service. 
Users can create notes, with the creator being the designated owner. Additionally, the application provides a user-friendly list of registered users. Notes can be shared with other users, granting read or read/write access, with specific permissions tailored to individual users.
The content of each note is straightforward, focusing solely on text, without the need for embedded media. Every note features a name for easy identification, timestamps for creation and completion, and a status flag to track task progress. In the case of delegated tasks, the note will also specify the user to whom the task was delegated.
This project aims to deliver an efficient, web-based service with a lightweight front-end, prioritizing functionality over an elaborate graphical user interface.

--------------------------------------------------------------------------------------------------------------------------------
This project is based on [pwrcost] (https://github.com/yonush/pwrcost). It primarily utilizes built-in Go packages such as `net/http`, `database/sql`, `strconv`, and `html/template`. Additionally, it incorporates third-party packages, specifically `gorilla/mux` for routing, `jackc/pgx` for the PostgreSQL driver, `icza/session` for session management and `crypto/bcrypt` for password hashing.
:notebook:

#### Installation 
***To Start the Server on Your Windows System:***

1. Navigate to the root directory of the repository.
2. Run the `buildpkg.cmd` script to compile the binary, `EnterpriseNotes.exe`, using non-vendored packages.
3. Alternatively, if you want to build the binary with the vendor, execute the `buildvendor.cmd` script.
4. Start the application by running the `EnterpriseNotes.exe` binary, or use the provided `run.cmd` script (which sets environment variables).
5. To test the application, open your web browser and go to [http://localhost:8080](http://localhost:8080). If port 80 is not available, you can launch the application as follows:

   ```sh
   > EnterpriseNotes 8080
    ```

It should direct you to: http://192.168.1.128:8080/login

#### Building database configuration 
Install PostgreSQL & Create a Database: CREATE DATABASE EnterpriseNotes
Database defaults in the *app.go* are shown below.
``` go
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "EnterpriseNotes"
)
```

#### Build by pkg
``` bash
export GO111MODULE="on"
export GOFLAGS="-mod=mod"
go mod download
:: strip debug info during build
go build -ldflags="-s -w" .
``` 
#### Build by vendor
``` bash
export GO111MODULE="on"
export GOFLAGS="-mod=vendor"
go mod vendor
:: strip debug info during build
go build -ldflags="-s -w" 
```

### Dependencies
|package|docs link|
|:--|:--|
|net/http docs|https://pkg.go.dev/net/http
|html/template docs|Go HTML templates with [W3.CSS](https://www.w3schools.com/w3css/w3css_examples.asp) stylesheet|
|HTTP router: Gorilla mux|https://pkg.go.dev/github.com/gorilla/mux|
|Datastore: PostgreSQL driver|https://github.com/jackc/pgx/|
|Session Management: Session|https://pkg.go.dev/github.com/icza/session@v1.2.0|
|Password hashing tool:|https://pkg.go.dev/golang.org/x/crypto@v0.14.0/bcrypt|

## Server Routes / Functions
*Note*
|route|method|description|
|:--|:-:|:--|
|`/`|`GET`|Index page |
|`/list`|`GET`|List notes|
|`/list/{srt:[0-9]+}`|`GET`|List notes with sorting |
|`/create`|`POST, GET`| Create note |
|`/getdelegations`|`GET`|Get note delegations|
|`/update `|`POST, GET`|Update note |
|`/delete `|`POST, GET`|Delete note|

*User Setting*
|route|method|description|
|:--|:-:|:--|
|`/login`|`POST, GET`|Login  |
|`/logout` |`GET`|Logout |
|`/register` |`POST, GET`|egister|
|`/update-settings`|`POST`|Update user settings |
|`/deleteUser` |`POST`|Delete user |

## Datastore
This version of the application necessitates a dedicated database, specifically PostgreSQL. It also involves the initial import of several CSV files from the local data folder. These CSV files are imported during the first run of the application. Subsequently, the application relies on the database for its functioning during every execution.

#### Datatypes / Tables
#####*User Table*
| Attribute       | Type        |
| --------------- | ----------- |
| UserID          | int         |
| Username        | string      |
| Password        | string      |

#####*Note Table*
| Attribute       | Type        |
| --------------- | ----------- |
| NoteID          | int         |
| UserID          | int         |
| NoteTitle       | string      |
| NoteContent     | string      |
| CreationDate    | time.Time   |
| DelegatedTo     | string      |
| CompletionDate  | time.Time   |
| Status          | string      |

#####*Sharing Table*
| Attribute       | Type        |
| --------------- | ----------- |
| SharingID       | int         |
| NoteID          | int         |
| UserID          | int         |
| Timestamp       | time.Time   |
| Status          | string      |

#####*Data Table*
| Attribute       | Type        |
| --------------- | ----------- |
| userName		  | string  	|
| note	          |[]DisplayNote|

#####*DisplayNote Table*
| Attribute       | Type        |
| --------------- | ----------- |
| NoteTitle       | string      |
| CreateDate      | time.Time   |
| Delegation      | string      |
| CompletionDate  | string      |
| Status      	  | time.Time   |
| Username  	  | string      |

## Session management
The application leverages the [icza/session](https://github.com/icza/session) module to manage basic sessions for authentication purposes. For details on the basic authentication implementation, please refer to [auth.go](https://github.com/KatherineMulder/ITPR6.518-Enterprise-Software-Development-Development/blob/main/auth.go).

## Password hashing
This application utilizes the [crypto/bcrypt](golang.org/x/crypto/bcrypt) module to exncrpy password entered by the user to store on the database for security reasons. For details in how we use this please refer to [auth.go](https://github.com/KatherineMulder/ITPR6.518-Enterprise-Software-Development-Development/blob/main/auth.go)

## Sample screens
