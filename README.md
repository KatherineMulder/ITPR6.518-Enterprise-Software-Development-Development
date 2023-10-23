# ITPR6.518-Enterprise-Software-Development-Development
 Katherine & Alex's Enterprise Software Developement Assignment

# Introduction
This project is based on [pwrcost](https://github.com/yonush/pwrcost). It primarily utilizes built-in Go packages such as `net/http`, `database/sql`, `strconv`, and `html/template`. Additionally, it incorporates third-party packages, specifically `gorilla/mux` for routing and `jackc/pgx` for the PostgreSQL driver.

**To Start the Server on Your Windows System:**

1. Navigate to the root directory of the repository.
2. Run the `buildpkg.cmd` script to compile the binary, `EnterpriseNotes.exe`, using non-vendored packages.
3. Alternatively, if you want to build the binary with the vendor, execute the `buildvendor.cmd` script.
4. Start the application by running the `EnterpriseNotes.exe` binary, or use the provided `run.cmd` script (which sets environment variables).
5. To test the application, open your web browser and go to [http://localhost:8080](http://localhost:8080). If port 80 is not available, you can launch the application as follows:

   ```sh
   > EnterpriseNotes 8080

It should direct you to: http://192.168.1.128:8080/login

## Building


# database configuration 
The app assumes a database exists - ESD. Edit the *app.go* to change the default database name. Database defaults in the *app.go* are shown below.
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
|net/http docs||
|html/template docs|Go HTML templates with [W3.CSS](https://www.w3schools.com/w3css/w3css_examples.asp) stylesheet|
|HTTP router: Gorilla mux|https://pkg.go.dev/github.com/gorilla/mux|
|Datastore: PostgreSQL driver|https://github.com/jackc/pgx/|


## Datastore

This version of the application necessitates a dedicated database, specifically PostgreSQL. It also involves the initial import of several CSV files from the local data folder. These CSV files are imported during the first run of the application. Subsequently, the application relies on the database for its functioning during every execution.


## Datatypes / Tables
| Attribute       | Type        |
| --------------- | ----------- |
| UserID          | int         |
| Username        | string      |
| Password        | string      |
| NoteID          | int         |
| UserID          | int         |
| NoteTitle       | string      |
| NoteContent     | string      |
| CreationDate    | time.Time   |
| DelegatedTo     | string      |
| CompletionDate  | time.Time   |
| Status          | string      |
| Privileges      | string      |
| SharedUsers     | []Sharing   |
| SharingID       | int         |
| NoteID          | int         |
| UserID          | int         |
| Timestamp       | time.Time   |
| Status          | string      |

## Session management

The application leverages the [icza/session](https://github.com/icza/session) module to manage basic sessions for authentication purposes. For details on the basic authentication implementation, please refer to [auth.go](https://github.com/yonush/pwrcost/blob/main/auth.go).

## Sample screens
