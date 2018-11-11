Instance
========

Create
------

  - verify that state is initial or executing (otherwise return)
  - set task state to executing if needed
  - ensure prequisites (active context)
    - loop over all context dependencies
    - filter all dependencies with no endpoints available
    - filter all dependencies with no subtask associated with the dependency
    - create subtasks for each found dependency
    - if subtasks have been created exit!

  - install
  - set component state of component to inactive
  - set component endpoint to nil
  - trigger completed event on task

Start
-----

- verify that state is initial or executing (otherwise return)
- set task state to executing if needed

- ensure prequisites (active services)
  - loop over all service dependencies
  - filter all dependencies with no endpoints available
  - filter all dependencies with no subtask associated with the dependency
  - create subtasks for each found dependency
  - if subtasks have been created exit!

- check if component is installed/inactive
- activate
- set component state of component to active
- set component endpoint to result of activate
- trigger component configure task
- trigger completed event on task

Stop
-----

- verify that state is initial or executing (otherwise return)
- set task state to executing if needed
- ensure prequisites (active services)
  - loop over all service dependencies
  - filter all dependencies with no endpoints available
  - filter all dependencies with no subtask associated with the dependency
  - create subtasks for each found dependency
  - if subtasks have been created exit!
- activate
- set component state of component to active
- set component endpoint to result of activate
- trigger completed event on task
