# serviceManager
A simple UI prototype for managing services in bash environments. This project is an extension from the [healthchecker](https://github.com/chensxb97/healthchecker) repo.

## Usage
- Start a service
- Stop a service
- Restart a service
- API healthchecks performed at regular intervals

Mechanism will be triggered from a backend agent running in a server.

## Assumptions
- Start, Stop and Restart actions are independent of existing app setup (systemd/systemxtl)
- Healthchecks are to be performed separately to determine the live state of the application (healthcheck endpoint, status 200 etc)
- Start, Stop and Restart will execute unique commands where each command starts a bash script to run each application.

### Setup
- Service start script (for start operation)
- Service name (for stop operation)
- Installation of agent binary with APIs exposed on a server (backend folder)
- Server should have `bash` installed
- Running the Web UI to communicate with the database and agent APIs (ui folder)