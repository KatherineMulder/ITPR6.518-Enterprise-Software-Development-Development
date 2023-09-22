# Application Requirements: 

## The users will need to setup a new account 

## All users can create a note

## the user who created a note is the owner of that note. 

## All users have access to the list of registered users. 

## A note can be shared with other users providing either read access or read/write access, which is set per user.

- [x]  A note just contains text 
- [x]  A note should have a name to aid with identification of the note.
- [x]  The note should contain a date & time of creation.
- [x]  And include a date & time of completion if it is a task.
- [x]  The note/task should include a status flag to indicate if the task is none/in progress/completed/cancelled/delegated.
- [x]  User id or name of the user the note was delegated to.


 :sparkles: Find a note 
 ##### users can search a note by providing a text pattern (text partterns are specified in the analyse a note)
 ##### Users can continue to search within the selected notes.
 ##### Users can also clear the list of selected notes, after which searching will apply to all notes in.

 


**Note**
\*Datebase\*
jackc/pgx for the database






# file structure: 

**manageNote**
1. editNote
2. deleteNote
3. createNote
4. readNote

**searchNote**
1. findNote
2. findnoteDetails
3. analysieNote

**userSetting**
1. createUser
2. retrieveUser *authentication*
3. updateUser
4. deleteUser

**main.go**

*Katherine*
1. createDatabase
2. createNote
3. find note
4. find note details
**deadline: 17/09**

*Alex*
1. edit note
2. delet enote
3. analyse note
**deadline: 22/09**


 

**during school holiday**
1. create users CRUD
2. main.go (setting up the database)
3. HTML (metting up for the GUI)
4. RESTful API


*end* 
1. read note
