# gkui
Graphical Kafka User Interface 


# Local Manual Testing
Docker Compose can be used to run a simple kafka server provided by Confluent

```
docker compose up .\infra\docker-local.yml
```

The application can be run with 

```
``` 

## Logging

The environment is setup to use the charm logger within a channel log approach. Please use the env.Logger for sending to logs. 

i.e.
``` 
env.Logger.log()
env.Logger.DebugLog()
env.Logger.WarnLog()
env.Logger.ErrorLog()

