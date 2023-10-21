Process and things we need for this assignment

**1\. database:** 

1.  Accessing MySQL
2.  Creating the Database
3.  Selecting the Database
4.  Forming the Users Table

  

**2\. Mapping the Application's Routes:**

Creating a router with middleware

*   Defining CRUD Routes:
*   POST requests are captured by the handleCreate function, responsible for creating a user.
*   GET requests are skillfully handled by the handleGet function, retrieving user data.
*   PUT requests find their way to the handleUpdate function, facilitating user updates.
*   DELETE requests find closure in the handleDelete function, gracefully deleting user records.

  

**3\. Launching the HTTP Server**

With routes set and handlers in place the HTTP server. The line http.ListenAndServe(":port", r) sets the stage for our server to listen on port you define, ready to direct incoming requests to their designated handlers.

  

**4\. User Data:**

To enable the creation of new users, we've crafted an HTTP endpoint and a dedicated function to handle incoming POST requests containing user information. In our main function, which boasts a finely-tuned router courtesy of the gorilla/mux package, we've laid the foundation for this user creation process. 

  

**5\. handleCreate:**

The handleCreate function encapsulates the process of creating a new user

*   1\. initiate a database connection via the sql.Open function, combining the database driver, username, password, and database name.
*   2\. CSV data sent in the HTTP request body is diligently parsed and decoded into a User struct.
*   3\. The CreateUser function is invoked to execute the crucial database operation.
*   4\. Error handling such as code 301 but we will put as http//\*\*\*\*
*   5\. Upon successful user creation, a positive HTTP response is crafted. A status code of 201 (Created) is sent, accompanied by the reassuring message, "User created successfully."

  

**6\. Creating a User in the Database:**

inserting a new user into the database

we will need: 

*   1\. It accepts three parameters:db: enabling interaction with the database.name: A string representing the name of the user to be inserted.email: A string representing the user's email address.
*   2\. The query variable holds the SQL command for inserting a user into the database.
*   3\. The db.Exec function executes the SQL query with the provided parameters. If any errors arise during this execution, they are gracefully handled and returned.
*   4\. Upon successful execution, the function returns nil, signifying that the user insertion has been accomplished without hitches.

  

**7\. Retrieving User Data:**

This function is responsible for handling requests to retrieve user data based on their ID:

*   1\. We initiate a database connection using sql.Open, incorporating the database driver, username, password, and database name. As always, error handling is in place to gracefully manage any issues.
*   2\. The function extracts the 'id' parameter from the URL
*   3\. We convert the 'id' parameter from a string to an integer, ensuring compatibility with the database query.
*   4\. The GetUser function is summoned to fetch user data from the database, using the database connection and user ID as parameters. This is where the magic happens.
*   5\. Error handling is paramount
*   6\. To complete the journey, we convert the user object into JSON format and send it as the response

  

**8\. Retrieving User Data from the Database:**

The crucial GetUser function conducts the actual database query to retrieve user data:

*   1\. We define the SQL query used to select a user from the users table based on their ID.
*   2\. The db.QueryRow function executes the query, returning a single row result.
*   3\. We create an empty User struct using the & operator to obtain a pointer to the struct.
*   4\. The row.Scan function scans the result into the User struct, populating its fields with the retrieved data.
*   5\. Error handling ensures that any issues during the database interaction are properly handled. If an error occurs, it's returned, signifying that the user retrieval process encountered a hitch.

  

**9\. Updating User Data:**

To modify user information, our handleUpdate function takes center stage. It retrieves a user based on their ID and updates their name and email fields with the new values provided. Let's delve into the intricacies of this function

  

**10\. Updating User Data in the Database:**

define the SQL query for updating a user's name and email in the users table, targeting the user with the specified ID

  

**11\. Deleting User Data: handleDelete**

  

**12\. Deleting User Data in the Database**

*   1\. define the SQL query for deleting a user from the users table, targeting the user with the specified ID.
*   2\. The db.Exec function is employed to execute the SQL query, passing in the 'id' as a parameter. Any errors arising during this process are captured and returned.
*   3\. Upon successful execution, the function returns nil, indicating that the user has been deleted from the database without any issues.

  

**13.Running the app and interacting with the database:** 

*   1\. To create a new user, you can use curl to make a POST request to the /user endpoint with the user's data. Here's the command:
*   2\. For update and delete , just swap the POST Verb with the correct one : PUT - DELETE
*   14\. Retrieving a User by ID:
*   To retrieve a user by their ID, you can make a GET request to the /user/{id} endpoint, where {id} is the ID of the user you want to fetch. Replace {id} with the actual user ID and visit : http://localhost:portNumber/user/{id}

  

## scope:

Apart from metadata required to administer the note/task;

• A note just contains text (embedded media like video, images, etc is not required).

• A note should have a name to aid with identification of the note.

• The note should contain a date & time of creation.

o And include a date & time of completion if it is a task.

• The note/task should include a status flag to indicate if the task is none/in progress/completed/cancelled/delegated.

• User id or name of the user the note was delegated to.

  

**Documentation:**

*   The provided requirements can be ambiguous and incomplete. It is part of your assignment to remove the ambiguity and complete the required information. Your choices need to be documented. You are expected to discuss your considerations with your lecturer during the lab sessions.
*   a document that discusses the use, maintenance and various design decisions and aspects of the application.

  

a little note:

*   error checking for SQL injection
*   before deleting the note we will need an error check for the note owner
*    a user won't need a function to filter the note, sharing notes only for the owner and filtering only for the owner.
*   analysis note
*   adjust the template for our web

  

  

  

## Assignment requirement:

**1\. Find a note**

Users can search for a note/task in the entire Enterprise Notes service by providing a text pattern. Valid text patterns are specified in 2. Analyse a note. All notes that satisfy the text pattern become available to the user, provided they have at least read access. Users can continue to search within the selected notes by providing text patterns. Users can also clear the list of selected notes, after which searching will apply to all notes in the Enterprise Notes service.

Users can also conduct searches on the following note/task attributes;

• Note name

• Note completion time (if it is a task)

• Note status

• Note delegation

• Shared users

  

**2\. Analyse a note**

Users can count occurrences of specified text snippets within a single note. Text snippets are expressed using a text pattern (specified below);

Valid text patterns for search and analysis features

• a sentence with a given prefix and/or suffix.

• a phone number with a given area code and optionally a consecutive sequence of numbers that are part of that number.

• an email address on a domain that is only partially provided.

• text that contains at least three of the following case-insensitive words: meeting, minutes, agenda, action, attendees, apologies.

• a word in all capitals of three characters or more.

  

**3\. Read and/or edit a note**

Based on the user's access privileges, the user can read and/or edit the properties and data in a note.

• Note name

• Note text

• Note completion date & time

• Note status

• Note delegation

• Shared users (refer to 4. Change note sharing details below.)

  

**4\. Change note sharing details**

Only the owner of a note can change the access privileges of the users the note is shared with. The owner is also the only user who can set and modify the list of users the note is shared with.

  

**5\. User settings**

As users often create notes shared with colleagues like they have done before, a user can save sharing settings. When a note is created, the owner can apply saved sharing settings to share the new note with the same people as before (using the same access privileges).


