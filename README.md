# Context

This projects simulates a small-scale trading app for stocks
where users can buy and sell stocks.

The architecture of the system will be based on event-sourcing, each command
being an event that is appended to a read-only log.

This project is an assignment I had for Software Design class, I included it
in my public repositories because it elevates real use cases of
event sourcing platform and is a well-designed system.

## Architecture

- Event sourcing architecture
- Event logs are saved on file system