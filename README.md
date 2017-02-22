# Vell

Vell is a lightweigt repository management tool for RPM repositories.
Every operation is done via Vell's REST API.

## Installation

Make sure that the `createrepo` package is installed:

```bash
$ yum install createrepo
```


## Configuration

* `VELL_REPO_PATH`: Base path where all repositories managed by Vell are
located, defaults to `/var/lib/vell-repositories`. Make sure the user
under which Vell runs has permission to read and modify this directory.
* `VELL_HTTP_ADDRESS`: Address on which Vell listens, e.g. `localhost`
or `127.0.0.1`, if not given Vell will listen on all interfaces
* `VELL_HTTP_PORT`: Port on which Vell listens, default is `8080`

## Usage

TODO.