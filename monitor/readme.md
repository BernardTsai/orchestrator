Monitor
=======

monitor visualises the orchestration flow of tasks.

Usage:
------

a) start the web application server which will listen to port 8081
of the server:

```
go run monitor.go
```

b) open a browser at the url: "http://localhost:8081"


--------------------------------------------------------------------------------

View:

User:
  Task*  
    Event*
Architectures:
  Architecture*
    Task*
      Event*
Components:
  Component*
    Task*
      Event*
    Version*
      Task*
        Event*
