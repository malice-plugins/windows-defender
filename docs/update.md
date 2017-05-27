# To update the AV run the following:

```bash
$ docker run --name=windows-defender malice/windows-defender update
```

## Then to use the updated windows-defender container:

```bash
$ docker commit windows-defender malice/windows-defender:updated
$ docker rm windows-defender # clean up updated container
$ docker run --rm malice/windows-defender:updated EICAR
```
