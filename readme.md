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

## Building
***************
Building
VERSION to be replaced with V for example v1.0.1 or use ``bundle.sh`` to build all platform binaries

**Darwin / MacOS**
****
```
env GOOS=darwin GOARCH=amd64 go build -o bin/macos-darwin/amd64/fairymq && tar -czf bin/macos-darwin/amd64/fairymq-VERSION-amd64.tar.gz -C bin/macos-darwin/amd64/ $(ls  bin/macos-darwin/amd64/)
```

```
env GOOS=darwin GOARCH=arm64 go build -o bin/macos-darwin/arm64/fairymq && tar -czf bin/macos-darwin/arm64/fairymq-VERSION-arm64.tar.gz -C bin/macos-darwin/arm64/ $(ls  bin/macos-darwin/arm64/)
```


**Linux**
****
```
env GOOS=linux GOARCH=386 go build -o bin/linux/386/fairymq && tar -czf bin/linux/386/fairymq-VERSION-386.tar.gz -C bin/linux/386/ $(ls  bin/linux/386/)
```

```
env GOOS=linux GOARCH=amd64 go build -o bin/linux/amd64/fairymq && tar -czf bin/linux/amd64/fairymq-VERSION-amd64.tar.gz -C bin/linux/amd64/ $(ls  bin/linux/amd64/)
```

```
env GOOS=linux GOARCH=arm go build -o bin/linux/arm/fairymq && tar -czf bin/linux/arm/fairymq-VERSION-arm.tar.gz -C bin/linux/arm/ $(ls  bin/linux/arm/)
```

```
env GOOS=linux GOARCH=arm64 go build -o bin/linux/arm64/fairymq && tar -czf bin/linux/arm64/fairymq-VERSION-arm64.tar.gz -C bin/linux/arm64/ $(ls  bin/linux/arm64/)
```


**FreeBSD**
****
```
env GOOS=freebsd GOARCH=arm go build -o bin/freebsd/arm/fairymq && tar -czf bin/freebsd/arm/fairymq-VERSION-arm.tar.gz -C bin/freebsd/arm/ $(ls  bin/freebsd/arm/)
```

```
env GOOS=freebsd GOARCH=amd64 go build -o bin/freebsd/amd64/fairymq && tar -czf bin/freebsd/amd64/fairymq-VERSION-amd64.tar.gz -C bin/freebsd/amd64/ $(ls  bin/freebsd/amd64/)
```

```
env GOOS=freebsd GOARCH=386 go build -o bin/freebsd/386/fairymq && tar -czf bin/freebsd/386/fairymq-VERSION-386.tar.gz -C bin/freebsd/386/ $(ls  bin/freebsd/386/)
```


**Windows**
****
```
env GOOS=windows GOARCH=amd64 go build -o bin/windows/amd64/fairymq.exe && zip -r -j bin/windows/amd64/fairymq-VERSION-x64.zip bin/windows/amd64/fairymq.exe
```

```
env GOOS=windows GOARCH=arm64 go build -o bin/windows/arm64/fairymq.exe && zip -r -j bin/windows/arm64/fairymq-VERSION-x64.zip bin/windows/arm64/fairymq.exe
```

```
env GOOS=windows GOARCH=386 go build -o bin/windows/386/fairymq.exe && zip -r -j bin/windows/386/fairymq-VERSION-x86.zip bin/windows/386/fairymq.exe
```


## Protocol & Language
***************

- **fairyMQ port is 5991**

- **fairyMQ consumer port is 5992 by default but can be changed.**

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

#### New Consumer
``` 
NEW CONSUMER 0.0.0.0:5992
```
HOST:PORT

#### Remove Consumer
``` 
REM CONSUMER 0.0.0.0:5992
```
HOST:PORT

#### List Consumers
``` 
LIST CONSUMERS
```

