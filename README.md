# Song Cloud

## Introduction

minimal sound cloud just for having fun with databases

## Requests

Password validation failed on database side:

```sh
curl -d '{ "username": "elahe", "first_name": "Elahe", "last_name": "Dastan", "email": "elahe.dstn@gmail.com", "password": "1234" }' -H 'Content-Type: application/json' 127.0.0.1:8080/api/signup
```

User creation completed:

```sh
curl -d '{ "username": "elahe", "first_name": "Elahe", "last_name": "Dastan", "email": "elahe.dstn@gmail.com", "password": "123456abc" }' -H 'Content-Type: application/json' 127.0.0.1:8080/api/signup
```

Extend premium period:

```sh
curl -d '{ "username": "elahe", "duration": 100000000000 }' -H 'Content-Type: application/json' 127.0.0.1:8080/api/extend
```

Create new free song:

```sh
curl -d '{ "name": "elahe", "file": "elahe.mp3", "production_year": 2021, "explanation": "new awesome song" }' -H 'Content-Type: application/json' 127.0.0.1:8080/api/song
```

Create new paid song:

```sh
curl -d '{ "name": "elahe-p", "file": "elahe.mp3", "production_year": 2021, "explanation": "new awesome song", "price": 100 }' -H 'Content-Type: application/json' 127.0.0.1:8080/api/song
```

Play a song:

```sh
curl -d '{ "id": 2, "username": "elahe" }' -H 'Content-Type: application/json' 127.0.0.1:8080/api/play
```

Create a category:

```sh
curl 127.0.0.1:8080/api/category/pop
```

Assig a song to a category:

```sh
curl -d '{ "id": 2, "category": 1 }' -H 'Content-Type: application/json' 127.0.0.1:8080/api/category
```

Buy a song:

```sh
curl -d '{ "username": "elahe", "song": 1 }' -H 'Content-Type: application/json' 127.0.0.1:8080/api/buy
```

like a song:

```sh
curl -d '{ "username": "elahe", "id": 1 }' -H 'Content-Type: application/json' 127.0.0.1:8080/api/like
```

## Reporting

Sum of the users transactions:

```sql
select username,sum(purchased_price) from purchase group by purchase.username;
```

Purchased logs from last 3 hours:

```sql
select * from log where log_message like '%purchased%' and time > now() - interval '3 hours';
```

Users that introduce 2 or more users:

```sql
select introducer from introduce group by introducer having count(*) >= 2;
```

Last year sells:

```sql
select sum(purchased_price) from purchase where extract(year from purchased_date) = 2021;
```

Last year best seller:

```sql
select * from song where id = (select song_id from purchase where extract(year from purchased_date) = 2021 group by song_id order by count(*) limit 1);
```
