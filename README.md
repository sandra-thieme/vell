# Vell

[![Build Status](https://travis-ci.org/rkcpi/vell.svg?branch=master)](https://travis-ci.org/rkcpi/vell)

Vell is a lightweight repository management tool for RPM repositories.
Every operation is done via Vell's HTTP API.

## Installation

Vell is still under development and no stable version has been released
as of yet. Once Vell hits 1.0 packages will be provided for download.

The only external dependency that Vell has is `createrepo`. Make sure
that the `createrepo` package is installed:

```bash
$ yum install createrepo
```


## Configuration

* `VELL_REPOS_PATH`: Base path where all repositories managed by Vell
are located, defaults to `/var/lib/vell/repositories`. Make sure the
user under which Vell runs has permission to read and modify this
directory.
* `VELL_HTTP_ADDRESS`: Address on which Vell listens, e.g. `localhost`
or `127.0.0.1`, if not given Vell will listen on all interfaces
* `VELL_HTTP_PORT`: Port on which Vell listens, default is `8080`

## Usage

Vell is managed entirely through its HTTP API. The following endpoints
are available:

### GET /repositories

Fetches a list of all repositories in `VELL_REPOS_PATH`.

Example request:

```bash
$ curl http://localhost:8080/repositories
```

Response with status code 200 OK:

```json
[
  {
    "name":"foo"
  },
  {
    "name":"bar"
  }
]
```

### POST /repositories

Creates a new repository.

Example request:

```bash
$ curl -H "Content-Type: application/json" -d '{"name":"baz"}' http://localhost:8080/repositories
```

Response will have status code 201 CREATED.

### POST /repositories/{name}/packages

Adds a new package to the repository with name `name`.

Suppose you have a package called `myapp-1.5.3.x86_64.rpm`. Example to
add this to your `foo` repository:

```bash
$ curl -X POST -F file=@path/to/myapp-1.5.3.x86_64.rpm http://localhost:8080/repositories/foo/packages
```

Response will have status code 201 CREATED.

### GET /repositories/{name}/packages

Lists all packages in the repository with name `name`.

Example request:

```bash
$ curl http://localhost:8080/repositories/baz/packages
```

Response with status code 200 OK:

```json
[
  {
    "name":"myapp-1.5.3.x86_64.rpm",
    "lastUpdated":"2017-03-07T13:02:58+01:00",
    "size":415661
  },
  ...
]
```