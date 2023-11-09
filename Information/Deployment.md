##Deployment Report

ITPR6.518-Enterprise-Software-Development-Development
_Katherine Mulder & Alex Borawski
Eastern Institute of Technology
NZ Bachelor of Computing Systems
John Jamieson
November 09, 2023_ 

**Introduction:**


**The Features of the Project**
Here are features of the Enterprise Notes application that can be done on your computer or device (the client side) instead of on the main server computer, along with the advantages and disadvantages:
* UI:
    Pros: Makes the app feel faster and more responsive.
    Cons: Can have security issues and might look different on different devices.
* Data sotrage: 
    Pros: Lets you use the app even without an internet connection.
    Cons: Your data might get lost, and it relies on your device's storage.
* Data validation:
    Pros: Helps prevent mistakes as you type.
    Cons: People might find ways to enter incorrect data.
* dynamic form generation:
    Pros: Creates forms that change based on what you're doing.
    Cons: Can be tricky to make work on all devices.
* User settings:
    Pros: Lets you personalize the app and have a better experience.
    Cons: Can be hard to keep your settings the same on different devices.
* Asynchronous requests:
    Pros: Makes the app faster and more responsive.
    Cons: Might not work well when getting data from other websites and can have security problems.
---------------------------------------------
**Database design choice**

All persistent data (notes, user accounts, etc.) are to be stored in a suitable datastore. Explain the design choices you made to interact with the database.


---------------------------------------------
**Hosting for variety of operating systems**
To ensure that the Enterprise Notes application is accessible and functional across a variety of operating systems and web browsers here are some strategies can be used: 
1. Use widely supported technologies: Build the app using web technologies that are supported by popular browsers like Chrome, Firefox, Safari, and Edge. This ensures that it will run smoothly for most users.

2. Make it fit any screen: Design the app so it looks good and works properly on different screen sizes. This way, people using desktops, laptops, tablets, or phones will all have a great experience.

3. Allow installation like an App: Develop the app as a Progressive Web App (PWA). PWAs can be installed on devices, just like regular apps. They work on various operating systems like Windows, macOS, Android, and iOS.

4. One code for many platforms: If you create a mobile app, consider using tools like React Native, Flutter, or Xamarin. These tools let you write one set of code that can run on different operating systems.

5. Detect the user's device: Use server-side tools to figure out what kind of browser and device the user is using. This way, you can customize the app for each platform.

6. Keep testing: Regularly test your app on different browsers and devices. This helps you find and fix problems quickly. You can even use automated tools to make this easier.

7. Have backup plans: Sometimes, not all browsers support the same features. Prepare alternatives for things that might not work on every browser.

8. Use special tools: Some libraries and tools can help your app work better on older browsers. They can add support for newer web technologies.

9. Tell users what's best: Make sure users know which browsers and operating systems work best with your app. Also, be ready to help them if they have trouble on specific platforms.

10. Stay up to date: Keep your app and all the pieces it relies on up to date. This way, you can take advantage of bug fixes, security updates, and improvements for different browsers.

These steps will make sure that the Enterprise Notes app is accessible and works well no matter which device or browser people are using.

---------------------------------------------
**Quick Start Guide**

_Prerequisites:_
* A Windows operating system.

* PostgreSQL database installed and running with the following configurations:
Host: localhost
Port: 5432
User: postgres
Password: postgres
Database Name: EnterpriseNotes

_Build the Application:_
1. Open a Command Prompt or PowerShell window and navigate to the Enterprise Notes application directory.
Eg: cd /path/to/EnterpriseNotes
1. Run the build script to compile the binary, EnterpriseNotes.exe, using non-vendored packages 
    >buildpkg.cmd
2. Build with Vendored Packages
    >buildvendor.cmd

_Start the Application:_
Start the application by running the EnterpriseNotes.exe binary from the Command Prompt or PowerShell:
> EnterpriseNotes.exe

You can also use the provided run script to set environment variables and start the application:
> run.cmd

_Access the Application:_
To test the application, open your web browser and go to http://localhost:8080. If port 80 is not available, you can launch the application with a different port as follows:
> EnterpriseNotes 8080


_Database Configuration:_
The application's database configuration is specified in the app.go file. Ensure that it matches your PostgreSQL configuration.
>const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "EnterpriseNotes"
)

_Building Database Configuration:_
To create a database named "EnterpriseNotes," you can execute the following SQL command:
> CREATE DATABASE EnterpriseNotes;
> 
This will create the necessary database for the application.

----------------------------------

Provide a document that lists the additional specifications that were missing but required to implement your solution.
- Owners can set read or read/write permissions?
