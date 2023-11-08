##Design Report

ITPR6.518-Enterprise-Software-Development-Development
_Katherine Mulder & Alex Borawski
Eastern Institute of Technology
NZ Bachelor of Computing Systems
John Jamieson
November 08, 2023_ 

**Title:** Enterprise Note and Task Tracker

**Introduction:**
This report outlines the design choices we've made in developing the Enterprise Note & Task Tracker. It's a guide to our objectives and the criteria we've followed to create this software.

**Design Objectives:**
+ Search and analysis system that allows users to find notes and analyze their content based on various text patterns and attributes
+ A simple, user-friendly interface that enables users to easily read, edit, and share notes and tasks, regardless of their access privileges.
+ The process of sharing notes with colleagues and allow owners to manage sharing settings 
+ Ensure that the project can scale to meet enterprise needs and offer flexibility in data storage by supporting  PostgreSQL.
+ Comprehensive error handling mechanisms to gracefully manage errors.
+ A user account management that allows users to create, modify, and manage their accounts within the Enterprise Notes service. Ensure a user-friendly registration and login process.
+ A clear ownership structure, where the user who creates a note is designated as the owner.
+ Focus on creating a text-based note creation system, where users can input and edit text content. Ensure that media elements such as video and images are not required.
+ Include a "name" attribute for notes to aid in their identification. This will help users quickly distinguish between different notes.
+ Incorporate date and time tracking features for note creation and completion, when relevant.
+ a status flag for notes and tasks, allowing users to indicate whether a task is "none," "in progress," "completed," "cancelled," or "delegated." This feature provides clarity on the current status of each note.
+ Capture and display the user ID or name of the user to whom a note has been delegated. This feature helps users keep track of the delegation process and understand the responsible parties.

**Design Criteria:**
> 1. User-friendly interface
> 2. Access control and user management
> 3. Note creation
> 4. Access permissions
> 5. Data integrity 
> 6. Date and time tracking
> 7. Status indication
> 8. User delegation information
> 9. Error Handling 
> 10. Documentation

**Design Process:**
1. **Project Initiation:** 
   Collaborate with team members to ensure a clear understanding of the application's purpose and user base.
2. **Requirement Analysis:** 
   Gather  requirements that encompass features, functionality, and expectations. Detail the specific needs for note and task management, access control, and data modeling.
3. **Architecture Design:** 
   Develop an architectural plan that outlines the system's structure,and data flow. 
4. **UI Design:**
   Light weight front-end UI.
5. **Database Design:** 
   Design the database schema to efficiently store notes, user information, access permissions, and other relevant data.
6. **Access Control Mechanisms:** 
    Access control features that allow note owners to share notes with other users, defining read or read/write permissions.
7. **User Account Management:**
   User registration, login, and account management features. Simplify the process of creating and managing user accounts, including password recovery and account deletion.
8. **Note Creation:**
   Enable users to efficiently create notes and tasks. Ensure that notes capture essential details, such as names, creation date and time, and, when relevant, task completion date and time
9. **Data Security Measures:**
    Encryption, authentication, and authorization to protect data from breaches and unauthorized access. Handle data with care to prevent corruption.
10. **Documentation:**
    Comprehensive documentation
<br>

**Design Choices:**
+ We've chosen PostgreSQL as the primary database for data storage due to its scalability.
+ The front-end of the application will have a minimalist design, prioritizing functionality over complex visuals.
+ Our access control mechanism enables note owners to share notes with specific users, granting read or read/write permissions.
+ We prioritize data security with end-to-end encryption, user authentication, and authorization.
+ For error handling, we'll use fmt.Println("") to output programming errors.
+ We've selected the Gorilla mux HTTP router for routing requests within the application.
+ Datastore flexibility is ensured by using PostgreSQL, a reliable relational database.
+ User session management is handled using the "Session" package.
+ Passwords will be securely hashed using the bcrypt package for enhanced data protection.
+ The application's UI design follows a clean and user-friendly approach.
+ Access permissions, user registration, and login processes are integral components of the user management system.
+ We've structured the design to fulfill criteria related to access control, note creation, data integrity, and user account management.
+ To meet design criteria, notes will include attributes like names, creation dates, completion status, and user delegation information.
+ The design process encompasses project initiation, requirement analysis, architecture design, UI design, database design, and more.
+ We've focused on creating a text-based note creation system, excluding embedded media like videos and images.
+ The status flag for notes and tasks will provide clarity on their progress, indicating whether they are "none," "in progress," "completed," "cancelled," or "delegated."

####Conclusion:
In conclusion, our design choices for the Enterprise Note & Task Tracker have been carefully considered to meet the needs of users within an enterprise. We've set clear design objectives, ensuring that our application offers a user-friendly experience. This includes features for searching and analyzing notes, user management, and data integrity.

Our design process involved careful planning, from project initiation to database design, and we've used dependable dependencies such as PostgreSQL for data storage and Gorilla mux for routing requests.

Our design choices prioritize data security, user access control, and a clean, minimalist user interface. We've focused on functionality over complex visuals, making it easier for users to read, edit, and share notes.

We've ensured that notes contain essential information, such as names, creation dates, and completion statuses. Users can also delegate notes to others and set access permissions. The application's design includes error handling mechanisms for a smooth user experience.