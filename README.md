# Rungraphene

Rungraphene is an OCI compliant runtime which can run an OCI bundle in a [Graphene] environment. This can be used with docker 1.11+ to run docker images too. 

### Build Instructions
Pre-requisistes.
  - [Graphene] 
  - Go 1.6+
  - Docker 1.11+
 
Before building, edit [main.go](main.go). Change *logFile*  to a location in where you want the rungraphene logs to get saved. And grapheneBootstrap to the graphene root dir path on your system. The root path is base dir of [Graphene].

```sh
$ git clone https://github.com/Sridarshan/rungraphene.git
$ cd rungraphene
$ make
$ cp rungraphene /usr/bin/
```

### Examples
To test your rungraphene, download a sample docker images and test it.

Start the docker daemon
```sh
dockerd --add-runtime rungraphene=/usr/bin/rungraphene --default-runtime rungraphene
```

Download a sample docker image I create. 
```sh
docker pull sridarshan/busybox_hello
```

Test it.
```sh
docker run -it busybox_hello /hello
```

### Building Docker Images
I will fill this section with resources to help you create custom docker images.

### Links
* [OCI Runtime Spec](https://github.com/opencontainers/runtime-spec)
* [Docker](https://github.com/docker/docker)
* [Containerd](https://github.com/docker/containerd)

[//]: # (These are reference links used in the body of this note and get stripped out when the markdown processor does its job. There is no need to format nicely because it shouldn't be seen. Thanks SO - http://stackoverflow.com/questions/4823468/store-comments-in-markdown-syntax)
    
   [Graphene]: <https://github.com/oscarlab/graphene>
