![fairymq.png](images%2Ffairymq.png)
fairyMQ is an open-source UDP based message queue software written in GO.  fairyMQ is in-memory, secure, extremely fast and super simple.


## Setup
**************
fairyMQ expects a keypair per queue.  When sending messages to fairy the messages must be encrypted with the queues public key.
fairy will try all available keys to decrypt a message before failing.  Upon successful encryption fairy knows which queue to access.

``` 
./fairymq --generate-queue-key-pair example
```

## Native Clients
**************
- **GO** https://github.com/fairymq/fairymq-go
- **Node.JS** (coming soon)
- **Java** (coming soon)
- **Python** (coming soon)
- **C#** (coming soon)

## Native Consumers
**************
- **GO** (coming soon)
- **Node.JS** (coming soon)
- **Java** (coming soon)
- **Python** (coming soon)
- **C#** (coming soon)


## Protocol
***************
A simple language fairyMQ understands.  This is the language between a fairyMQ client and server.

#### New message
``` 
ENQUEUE\r\n
timestamp\r\n
..bytes
```

#### First message in queue
``` 
FIRST IN
```

#### Last message in queue
``` 
LAST IN
```

#### Length of queue
``` 
LENGTH
```

#### Remove last message
``` 
POP
```

#### Remove first message
``` 
SHIFT
```

#### Remove/Clear queue
``` 
CLEAR
```



