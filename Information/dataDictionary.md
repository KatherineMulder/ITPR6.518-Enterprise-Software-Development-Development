| Table name | Field name | Default value | Data type | Field size | Example | 
| --- | --- | --- | --- | --- | --- | 
| User | userID(PK)| Auto increment | Int | | 100 |
| User | username | None | Varchar | 100 | Mike | 
| User | password | None | Varchar | 100 | Mike123 | 
| User | email | None | Varchar | 100 | Mike.123@hotmail.com | 
| User | registrationDate | Today | Date | | yyyy-mm-dd |  
| Note | noteID(PK)| Auto increment | Int | | 101 | 
| Note | userID | Auto increment | Int | | 102 |  
| Note | noteTitle | None | Varchar | 50 | Things to do | 
| Note | noteContent | None | Text | | Have a coffee in the morning | 
| Note | creationDate | Today | Timestamp | | 2023-10-02 10:00:30 | 
| Note | completionDate | Today | Timestamp | | 2023-10-02 10:00:30 | 
| Note | status | none | Varchar | 50 | in progress/ completed/ cancelled/ delegated | 
| Sharing | sharingID | Auto Increment | Int | 102 |
| Sharing | noteID | None | Int | | 102 | 
| Sharing | userID | None | Int | | 101 | 
| Sharing | status | Active | Varchar | 50 | "Active" or "Revoked" | 
| Sharing | timestamp | Today | Timestamp | | 2023-10-02 10:00:30 |
| Sharing | writtingSetting | false | Bool | | true/false |