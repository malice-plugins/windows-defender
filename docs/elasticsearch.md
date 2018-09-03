# To write results to [ElasticSearch](https://www.elastic.co/products/elasticsearch)

## Write to a `elasticsearch` docker container

```bash
$ docker volume create --name malice
$ docker run -d --name elasticsearch \
                -p 9200:9200 \
                -v malice:/usr/share/elasticsearch/data \
                 blacktop/elasticsearch:6
$ docker run --rm --link elasticsearch \
             -v /path/to/malware:/malware:ro \
             -e MALICE_ELASTICSEARCH_URL=http://elasticsearch:9200 \
             malice/windows-defender -t FILE
```

## Write to an external `elasticsearch` database

```bash
$ docker run --rm \
             -e MALICE_ELASTICSEARCH_URL=$MALICE_ELASTICSEARCH_URL \
             -e MALICE_ELASTICSEARCH_USERNAME=$MALICE_ELASTICSEARCH_USERNAME \
             -e MALICE_ELASTICSEARCH_PASSWORD=$MALICE_ELASTICSEARCH_PASSWORD \
             -e MALICE_ELASTICSEARCH_INDEX="test" \
             -v /path/to/malware:/malware:ro \
              malice/windows-defender -t FILE
```
