Entity-Relationship
===================

Model
- [Domain]
  - [Template(Component)]
    - [TemplateVersion(Configuration)]
      - [Dependency]
  - Blueprints
    - [Blueprint(Version)]
      - [Deployment(Component)]
        - [DeploymentVersion(Configuration*)]
  - Components
    - [Component]
      - [Instance]
  - Tasks
    - [Task]
  - Events
    - [Event]

//------------------------------------------------------------------------------

Component Setup for a specific version
- size
- configuration
- state
- parents endpoints
- services endpoints

//------------------------------------------------------------------------------

TemplateVersion

- size
- type
- state
- configuration
- dependencies
  - type: service/context
  - name
  - component
  - version
