# windows-defender

[![Publish Docker Image](https://github.com/malice-plugins/windows-defender/actions/workflows/docker-image.yml/badge.svg)](https://github.com/malice-plugins/windows-defender/actions/workflows/docker-image.yml)
[![License](http://img.shields.io/:license-mit-blue.svg)](http://doge.mit-license.org)
[![Docker Stars](https://img.shields.io/docker/stars/malice/windows-defender.svg)](https://store.docker.com/community/images/malice/windows-defender)
[![Docker Pulls](https://img.shields.io/docker/pulls/malice/windows-defender.svg)](https://store.docker.com/community/images/malice/windows-defender)
[![Docker Image](https://img.shields.io/badge/docker%20image-285MB-blue.svg)](https://store.docker.com/community/images/malice/windows-defender)

Malice Windows Defender AntiVirus Plugin

> This repository contains a **Dockerfile** of [Windows Defender](https://www.microsoft.com/en-us/windows/windows-defender) for the malice plugin **malice/windows-defender**

---

### Dependencies

- [ubuntu:bionic (_84.1 MB_\)](https://hub.docker.com/_/ubuntu/)

## Installation

1. Install [Docker](https://www.docker.io/).
2. Download [trusted build](https://store.docker.com/community/images/malice/windows-defender) from public [docker store](https://store.docker.com): `docker pull malice/windows-defender`

## Usage

### NOTICE :warning:

Something has changed in the latest version of Docker `18.09.0` where we now need to use our own seccomp profile found [here](https://raw.githubusercontent.com/malice-plugins/windows-defender/master/seccomp.json)

```bash
docker run --init --rm malice/windows-defender EICAR
```

With seccomp profile

```bash
docker run --init --rm --security-opt seccomp=seccomp.json malice/windows-defender EICAR
```

### Or link your own malware folder:

```bash
$ docker run --init --rm -v /path/to/malware:/malware malice/windows-defender FILE

Usage: windows-defender [OPTIONS] COMMAND [arg...]

Malice Windows Defender AntiVirus Plugin

Version: v0.1.0, BuildTime: 20180903

Author:
  blacktop - <https://github.com/blacktop>

Options:
  --verbose, -V          verbose output
  --table, -t            output as Markdown table
  --callback, -c         POST results to Malice webhook [$MALICE_ENDPOINT]
  --proxy, -x            proxy settings for Malice webhook endpoint [$MALICE_PROXY]
  --elasticsearch value  elasticsearch url for Malice to store results [$MALICE_ELASTICSEARCH_URL]
  --timeout value        malice plugin timeout (in seconds) (default: 60) [$MALICE_TIMEOUT]
  --help, -h             show help
  --version, -v          print the version

Commands:
  update  Update virus definitions
  web     Create a Windows Defender scan web service
  help    Shows a list of commands or help for one command

Run 'windows-defender COMMAND --help' for more information on a command.
```

This will output to stdout and POST to malice results API webhook endpoint.

## Sample Output

### [JSON](https://github.com/malice-plugins/windows-defender/blob/master/docs/results.json)

```json
{
  "windows-defender": {
    "infected": true,
    "result": "Virus:DOS/EICAR_Test_File",
    "engine": "0.1.0",
    "updated": "20171112"
  }
}
```

### [Markdown](https://github.com/malice-plugins/windows-defender/blob/master/docs/SAMPLE.md)

---

#### Windows Defender

| Infected | Result                    | Engine | Updated  |
| :------- | :------------------------ | :----- | :------- |
| true     | Virus:DOS/EICAR_Test_File | 0.1.0  | 20171112 |

---

## Documentation

- [To write results to ElasticSearch](https://github.com/malice-plugins/windows-defender/blob/master/docs/elasticsearch.md)
- [To create a Windows Defender scan micro-service](https://github.com/malice-plugins/windows-defender/blob/master/docs/web.md)
- [To post results to a webhook](https://github.com/malice-plugins/windows-defender/blob/master/docs/callback.md)
- [To update the AV definitions](https://github.com/malice-plugins/windows-defender/blob/master/docs/update.md)

## Issues

Find a bug? Want more features? Find something missing in the documentation? Let me know! Please don't hesitate to [file an issue](https://github.com/malice-plugins/windows-defender/issues/new).

## CHANGELOG

See [`CHANGELOG.md`](https://github.com/malice-plugins/windows-defender/blob/master/CHANGELOG.md)

## Contributing

[See all contributors on GitHub](https://github.com/malice-plugins/windows-defender/graphs/contributors).

Please update the [CHANGELOG.md](https://github.com/malice-plugins/windows-defender/blob/master/CHANGELOG.md) and submit a [Pull Request on GitHub](https://help.github.com/articles/using-pull-requests/).

## Credit

Made possible by the awesome work by [@taviso](https://github.com/taviso/loadlibrary)

## License

MIT Copyright (c) 2022 **blacktop**
