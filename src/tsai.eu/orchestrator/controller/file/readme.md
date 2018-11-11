File Component
==============

Functionality:
--------------

A domain relates to a directory within a filesystem.

Each component has a name which is unique within a domain.
It resides at a location relative to its parent component within a file system and holds a template which may contain references to other components, eg.


<code>
Hello <span style="color:orange">{{customer}}</span>&nbsp;   
</code>

The reference {{customer}} in this case would relate to the name of another component which resides in another location within the same filesystem.

The functionality is provided by a simple web service which is given a path in the filesystem and then recursively evaluates the templates as needed and returns the result.

Special note: the evaluation algorithm evaluates all instances of a component and makes use of the majority opinion to simulate cluster behaviour.

General structure
-----------------

Files are represented in a directory structure within a root context:

```
<root>                    Directory
  ...                     ...
    <parent>              Directory
      <file>              Directory
        .data             Directory
          .component      File
          <instance 1>    File
          <instance 2>    File
          <instance 3>    File
          <instance ...>  File
        <child 1>         Directory
        <child 2>         Directory
        <child 3>         Directory
        <child ...>       Directory
      <sibling 1>         Directory
      <sibling 2>         Directory
      <sibling 3>         Directory
      <sibling ...>       Directory
```

Contents of .component
----------------------

Content is in yaml format:

```
domain: <domain>
component: <component name>
state: <state of component>
path: <path of file directory>
endpoints:
  <version A>:
    path: <path of endpoint version>
  ...
```

Contents of "instance N"
------------------------

Content is in yaml format:

```
domain: <domain>
component: <component name>
instance: <uuid of instance>
version: <version of instance>
state: <state of instance>
path: <path of file directory>
endpoint:
    path: <path of endpoint version>
configuration:
  name: <name of the directory for the component>
  template: <template>
dependencies:
  <name>:
    type: context/service
    name: <name>
  ...
```
