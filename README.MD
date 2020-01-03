### Testing bulk insert from lib/pq

Up container with postgres
```bash
docker-compose up -d
```

Generate rows and bulk insert
```bash
./bulk -batch=100 -size=1000 -startID=1000000
```

```
./bulk -help
Usage of ./bulk:
  -batch int
        count of batches (default 50)
  -size int
        size of batch (default 100)
  -startID int
        start id (default 1)

```