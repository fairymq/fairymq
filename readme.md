![fairymq.png](images%2Ffairymq.png)
fairyMQ is an open-source UDP based message queue software written in GO.  fairyMQ is in-memory, secure, extremely fast and super simple.


## Setup
**************
fairyMQ expects a keypair per queue.  When sending messages to fairy the messages must be encrypted with the queues public key.
fairy will try all available keys to decrypt a message before failing.  Upon successful encryption fairy knows which queue to access.

``` 
./fairymq --generate-queue-key-pair example
```


## Broker talk
***************
A simple language fairyMQ understands.  This is the language between a fairyMQ client and server.

#### New message
``` 
ENQUEUE queuename\r\n
bytes..
```

#### First message in queue
``` 
FIRST IN queuename;
```

#### Last message in queue
``` 
LAST IN queuename;
```

#### Length of queue
``` 
LENGTH OF queuename;
```

#### Remove last message
``` 
POP queuename;
```

#### Remove first message
``` 
SHIFT queuename;
```

#### Remove/Clear queue
``` 
CLEAR queuename;
```

#### List queues
``` 
QUEUES;
```


