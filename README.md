```
docker run --rm --pull always -p 80:8000 -v /assets:/dbdir surrealdb/surrealdb:1.0.0-beta.9-20230402 start --user root --pass root file:assets/db/mydatabase.db
```
