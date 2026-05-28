# backend

## backend structure

```
├───cmd
│   └───server
└───internal
    ├───api
    │   ├───handlers
    │   └───middleware
    ├───blockchain
    ├───db
    ├───models
    └───utils
```


How i handle the database logic

Database Connection Setup i  established a connection to our PostgreSQL database (via Supabase). 

Implementation of Data i moved from using fake data to writing actual SQL queries by doing  
GetTasksByTeam, IssueCredential, UpdateTaskStatus

i also implemented graceful degradation where when app cant conect to internet or sever it doesnt crush but it works by letting the user know its in demo mode

acid cpmplience where if server crashes while doing task like writing data the data is not corrupted