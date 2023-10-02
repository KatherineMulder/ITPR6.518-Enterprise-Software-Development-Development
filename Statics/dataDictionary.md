| Table name | Field name | Default value | Data type | Field size | Example | 
| --- | --- | --- | --- | --- | --- | 
| User | userID | Auto increment | Int | | 100 |
| User | userName | None | Varchar | 50 | Mike | 
| User | password | None | Varchar | 50 | Mike123 | 
| User | email | None | Varchar | 100 | Mike.123@hotmail.com | 
| User | registrationDate | Today | Date | | yyyy-mm-dd | 
| User | readSetting | None | Bool | | On/off | 
| User | writtingSetting | None | Bool | | On/off | 
| Note | noteID | Auto increment | Int | | 101 | 
| Note | userID | Auto increment | Int | | 102 | 
| Note | delegatedToUserID | Auto increment | Int | | 103 | 
| Note | noteTitle | None | Varchar | 50 | Things to do | 
| Note | noteContent | None | Text | | Have a coffee in the morning | 
| Note | creationDateTime | Today | Timestamp | | 2023-10-02 10:00:30 | 
| Note | completionDateTime | Today | Timestamp | | 2023-10-02 10:00:30 | 
| Note | status | In progress | Varchar | 50 | in progress/ completed/ cancelled/ delegated | 
| Note | shareUser | None | Varchar | 255 | Sharon | 
| sharing | sharingID | Auto increment | Int | | 100 | 
| sharing | noteID | Auto increment | Int | | 102 | 
| sharing | userID | Auto increment | Int | | 101 | 
| sharing | status | Active | Varchar | 50 | "Active" or "Revoked" | 
| sharing | timestamp | Today | Timestamp | | 2023-10-02 10:00:30 |