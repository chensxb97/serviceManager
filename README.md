# serviceManager
A simple UI prototype for managing services

## Usage
- Start a service
- Stop a service
- Restart a service

Mechanism will be triggered from an agent written in backend.

## Assumptions
- Start, Stop and Restart actions are independent of existing app setup (systemd/systemxtl)
- Healthchecks are to be performed separately to determine the live state of the application (healthcheck endpoint, status 200 etc)
- Start, Stop and Restart will execute unique commands where each command starts a script to run each application.

### Setup
- Application start script
- Application stop regex (For killing relevant process pids)
- Database for storing the list of application's start/stop scripts, indexed by server identity 
- Installation of agent binary with APIs exposed on a server (backend folder)
- Running the Web UI to communicate with the database and agent APIs (ui folder)