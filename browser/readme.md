Browser
=======

browser visualises the results of an orchestration of file components.

Preparation:
------------
These file components reside in the "/tmp/data" directory.

The "data" directory contains contents which can be copied to this directory
in order to pre-fill this directory with some demo-content.

In addition it is necessary to have already installed the orchestrator packages
into the go workspace.

Usage:
------

a) start the web application server which will listen to port 8080
of the server:

```
go run browser.go
```

b) open a browser at the url: "http://localhost:8080"

Architecture:
-------------

browser provides 3 functions:

**1. Domain Inventory**    

  The url: "http://localhost:8080/data" will display information related
  to all file components within the domain as yaml structure.

**2. Render File Component**

  The url: "http://localhost:8080/render/[path]:[version]" will
  render the information related to a specific version of a file component
  identified by the path.

**3. Static files**

  All other types of urls are served as static files which would reside in the
  "static"  subdirectory. These files provide the client-side software to
  visualise the contents of the domain in a web-browser.
