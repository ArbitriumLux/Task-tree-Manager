# Task-tree Manager
Web server application (REST API) with functionality for working with the data structure Task-> Minitasks-> Labor costs  
Database - PostgreSQL  
## Functions:  
**Index page**  
* Create Task  
* Crate Minitask for Task  
* Create Labor costs for Minitask  

**/tree/{TaskName}**  
* Select tree Task-> Minitasks-> Labor costs, display on page (For each of the elements, the total time spent and the average time for all subordinate elements is calculated by goroutines)  
* Removing any of the items from the displayed list with data recalculation  
* Reassignment of Minitasks, Labor costs to other parents  
---
**The logging level of the application is set by the command line parameter (-LogLevel string)**  
**Default log level is "debug"**  
