malice-windows-defender
=======================

Malice Windows Defender AntiVirus Plugin

[![Circle CI](https://circleci.com/gh/maliceio/malice-windows-defender.png?style=shield)](https://circleci.com/gh/maliceio/malice-windows-defender)
[![License](http://img.shields.io/:license-mit-blue.svg)](http://doge.mit-license.org)
[![Docker Stars](https://img.shields.io/docker/stars/malice/windows-defender.svg)](https://hub.docker.com/r/malice/windows-defender/)
[![Docker Pulls](https://img.shields.io/docker/pulls/malice/windows-defender.svg)](https://hub.docker.com/r/malice/windows-defender/)
[![Docker Image](https://img.shields.io/badge/docker%20image-277%20MB-blue.svg)](https://hub.docker.com/r/malice/windows-defender/)

https://github.com/taviso/loadlibrary

> :warning: **NOTE:** Will not work on Docker for Mac because `CONFIG_MODIFY_LDT_SYSCALL` is not enabled :warning:  

This repository contains a **Dockerfile** of [Bitdefender](http://www.windows-defender.com/business/antivirus-for-unices.html) for [Docker](https://www.docker.io/)'s [trusted build](https://hub.docker.com/r/malice/windows-defender/) published to the public [DockerHub](https://hub.docker.com).

### Dependencies

-	[ubuntu (*118 MB*\)](https://hub.docker.com/_/ubuntu/)

### Installation

1.	Install [Docker](https://www.docker.io/).
2.	Download [trusted build](https://hub.docker.com/r/malice/windows-defender/) from public [DockerHub](https://hub.docker.com): `docker pull malice/windows-defender`

### Usage

```
docker run --rm malice/windows-defender EICAR
```

#### Or link your own malware folder:

```bash
$ docker run --rm -v /path/to/malware:/malware:ro malice/windows-defender FILE

Usage: windows-defender [OPTIONS] COMMAND [arg...]

Malice Windows Defender AntiVirus Plugin

Version: v0.1.0, BuildTime: 20170527

Author:
  blacktop - <https://github.com/blacktop>

Options:
  --verbose, -V         verbose output
  --table, -t	        output as Markdown table
  --callback, -c	    POST results to Malice webhook [$MALICE_ENDPOINT]
  --proxy, -x	        proxy settings for Malice webhook endpoint [$MALICE_PROXY]
  --timeout value       malice plugin timeout (in seconds) (default: 60) [$MALICE_TIMEOUT]    
  --elasitcsearch value elasitcsearch address for Malice to store results [$MALICE_ELASTICSEARCH]   
  --help, -h	        show help
  --version, -v	        print the version

Commands:
  update	Update virus definitions
  web       Create a windows-defender scan web service  
  help		Shows a list of commands or help for one command

Run 'windows-defender COMMAND --help' for more information on a command.
```

This will output to stdout and POST to malice results API webhook endpoint.

## Sample Output

### JSON:

```json
{
  "windows-defender": {
    "infected": true,
    "result": "Virus:DOS/EICAR_Test_File",
    "engine": "0.1.0",
    "updated": "20170527"
  }
}
```

### STDOUT (Markdown Table):

---

#### Windows Defender

| Infected | Result                    | Engine | Updated  |
| -------- | ------------------------- | ------ | -------- |
| true     | Virus:DOS/EICAR_Test_File | 0.1.0  | 20170527 |

---

Documentation
-------------

-	[To write results to ElasticSearch](https://github.com/maliceio/malice-windows-defender/blob/master/docs/elasticsearch.md)
-	[To create a Bitdefender scan micro-service](https://github.com/maliceio/malice-windows-defender/blob/master/docs/web.md)
-	[To post results to a webhook](https://github.com/maliceio/malice-windows-defender/blob/master/docs/callback.md)
-	[To update the AV definitions](https://github.com/maliceio/malice-windows-defender/blob/master/docs/update.md)

### Issues

Find a bug? Want more features? Find something missing in the documentation? Let me know! Please don't hesitate to [file an issue](https://github.com/maliceio/malice-windows-defender/issues/new).

### CHANGELOG

See [`CHANGELOG.md`](https://github.com/maliceio/malice-windows-defender/blob/master/CHANGELOG.md)

### Contributing

[See all contributors on GitHub](https://github.com/maliceio/malice-windows-defender/graphs/contributors).

Please update the [CHANGELOG.md](https://github.com/maliceio/malice-windows-defender/blob/master/CHANGELOG.md) and submit a [Pull Request on GitHub](https://help.github.com/articles/using-pull-requests/).

### License

MIT Copyright (c) 2016-2017 **blacktop**

#### Windows Defender

| Infected | Result               | Engine | Updated  |
| -------- | -------------------- | ------ | -------- |
| true     | Trojan:Win32/EyeStye | 0.1.0  | 20170527 |
