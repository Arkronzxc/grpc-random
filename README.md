# grpc-random
Microservice for PRNG using Marsenne Twister. With gRPC interface.


## Using
```shell script
docker build -t random-micro .
docker run -p 8080:8080 random-micro
``` 

### Health check

* Install [evans](https://github.com/ktr0731/evans)

* In terminal:

```shell script
evans -p 8080 internal/proto/get_random_numbers.proto
```
* In evans to check service availability:
```shell script
 call Healtcheck
```
