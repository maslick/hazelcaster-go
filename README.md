# hazelcaster-go

## Build
```
git clone https://github.com/maslick/hazelcaster-go && cd hazelcaster-go
go build
```

## Server
```
docker run -d -p 5701:5701 hazelcast/hazelcast
```

## Client
```
HZ_SERVER_ADDR=`docker-machine ip default`:5701 ./hazelcaster-go
go test
```
