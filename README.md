# Hash Service Implementation

This repo contains the Go code that implements a hashing service

- [Hash Service Implementation](#hash-service-implementation)
  * [Endpoints](#endpoints)
    + [Hashing](#hashing)
      - [POST /hash?password=mypassword](#post--hash-password-mypassword)
      - [GET /hash/{{requestId}}](#get--hash---requestid--)
    + [Statistics](#statistics)
      - [GET /stats/](#get--stats-)
    + [Shutdown](#shutdown)
      - [* /shutdown](#---shutdown)
  * [Build/Running](#build-running)
  * [Testing](#testing)
  * [Caveats, Parting Words](#caveats--parting-words)

<small><i><a href='http://ecotrust-canada.github.io/markdown-toc/'>Table of contents generated with markdown-toc</a></i></small>


## Endpoints

### Hashing
#### POST /hash?password=mypassword

This returns immediately with a "requestId" that can be used in the GET endpoint to retrieve the
actual hashed value, which is hashed using the SHA-512 alogrithm and base64 encoded before returning.  The actual hashing is queued and processed at a later time, currently set to 5 seconds from the original request.

#### GET /hash/{{requestId}}

This returns the base64 encoded, SHA-512 hashed password value from the POST request, or empty if the password hasn't been processed yet.

### Statistics

#### GET /stats/

This returns a JSON response that contains the number of requests to the POST /hash/ endpoint and their average request times in milliseconds:

```
HTTP 200 /
{
    "total": 1234,
    "average": 123.456
}

```

The average is calculated using:  Sum(request start time - request end time) / total request count.

### Shutdown

#### * /shutdown

This endpoint closes the backend worker (which processes/hashes passwords out of band from the request).  This also fails any new incoming requests.  The worker will wait for all pending requests before finally stopping.  However, currently the server still services GET /stats and GET /hash requests.

## Build/Running

There's a bash file that contains the build/run commands to start the server/api.  Currently there are two settings: port and delayPerRequest.  These can be changed in main.go if necessary/wanted.  The defaults are port=8080 and delayPerRequest=5 representing the 5 second delay before processing hash requests.  To run using the bash script:

* chmod 700 build.sh
* source build.sh

To build/run from commandline

* go build ./main.go
* ./main

## Testing

The repo contains a postman file that can be imported and used with the runner tool.  The scripts contain calls for /hash both POST and GET as well as /stats.  The /shutdown handler needs to either be added manually or run using CURL or any browser.

## Caveats, Parting Words

* Testing is missing - This was due to the time constraint and honestly rusty-ness in my Go syntax.  However, I'm hoping that the architecture is clear and the separation of concerns e.g. handlers, middleware and sparing usage of the repository pattern makes the code easy to cover.
* Idempotency and POST - By RFC POST is not idempotent, however, I made a judgement call in that hashing password that are identical is likely to be an overhead depending on the size of the input.  As such, although there are not intrinsic guarentees about idempotency, the implementation keeps a lookup of previous passwords and returns the existing hash.  Clearly this "could" be a security concern should a nefarious actor find a way to attack/dump the memory space of the application but it's unlikely in this application.
* Logging - Logs are enabled to stdout, it would be better to have this config driven to guide the log level during different environment deploys.  However, for debugging and just my curiosity into what the application was doing I left them in.

There's a lot more that can be done, for example, using command line, env driven configuration, i.e. flags.  Better routing using a different library, etc.  But for now I think this is simple enough implementation that meets the requirements.
