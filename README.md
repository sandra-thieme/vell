# Vell

[![Build Status](https://travis-ci.org/rkcpi/vell.svg?branch=master)](https://travis-ci.org/rkcpi/vell)

Vell is a lightweigt repository management tool for RPM repositories.
Every operation is done via Vell's REST API.

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
are located, defaults to `/var/lib/vell-repositories`. Make sure the
user under which Vell runs has permission to read and modify this
directory.
* `VELL_HTTP_ADDRESS`: Address on which Vell listens, e.g. `localhost`
or `127.0.0.1`, if not given Vell will listen on all interfaces
* `VELL_HTTP_PORT`: Port on which Vell listens, default is `8080`

## Usage

TODO.