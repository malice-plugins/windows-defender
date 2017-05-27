# POST results to a webhook

```bash
$ docker run -v `pwd`:/malware:ro --rm \
             -e MALICE_ENDPOINT="https://malice.io:31337/scan/file" malice/windows-defender --callback evil.malware
```
