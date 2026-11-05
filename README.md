# Song Cloud 🎶

## Introduction

Minimal sound cloud, written by Elahe Dastan just for having fun with SQL-based databases.
Users can have free or premium accounts. They can listen to musics, and we track they playlist.
Musics have categories which are predefined. Users can buy musics or like them.

This project is solely written for demonstration, and it has many performance issues.

## Requests

Password validation failed on database side:

```bash
curl -d '{ "username": "elahe", "first_name": "Elahe", "last_name": "Dastan", "email": "elahe.dstn@gmail.com", "password": "1234" }' -H 'Content-Type: application/json' 127.0.0.1:8080/api/signup
```

User creation completed:

```bash
curl -d '{ "username": "elahe", "first_name": "Elahe", "last_name": "Dastan", "email": "elahe.dstn@gmail.com", "password": "123456abc" }' -H 'Content-Type: application/json' 127.0.0.1:8080/api/signup
```

Extend premium period:

```bash
curl -d '{ "username": "elahe", "duration": 100000000000 }' -H 'Content-Type: application/json' 127.0.0.1:8080/api/extend
```

Create new free song:

```bash
curl -d '{ "name": "elahe", "file": "elahe.mp3", "production_year": 2021, "explanation": "new awesome song" }' -H 'Content-Type: application/json' 127.0.0.1:8080/api/song
```

Create new paid song:

```bash
curl -d '{ "name": "elahe-p", "file": "elahe.mp3", "production_year": 2021, "explanation": "new awesome song", "price": 100 }' -H 'Content-Type: application/json' 127.0.0.1:8080/api/song
```

Play a song:

```bash
curl -d '{ "id": 2, "username": "elahe" }' -H 'Content-Type: application/json' 127.0.0.1:8080/api/play
```

Create a category:

```bash
curl 127.0.0.1:8080/api/category/pop
```

Assig a song to a category:

```bash
curl -d '{ "id": 2, "category": 1 }' -H 'Content-Type: application/json' 127.0.0.1:8080/api/category
```

Buy a song:

```bash
curl -d '{ "username": "elahe", "song": 1 }' -H 'Content-Type: application/json' 127.0.0.1:8080/api/buy
```

like a song:

```bash
curl -d '{ "username": "elahe", "id": 1 }' -H 'Content-Type: application/json' 127.0.0.1:8080/api/like
```

## Playlists

Users can create playlists, add songs to them or remove songs from them, and share them with other users.
Only the owner can change or share a playlist, and users it is shared with can only read it.

| Method   | Path                                           | Description                                       |
| -------- | ---------------------------------------------- | ------------------------------------------------- |
| `POST`   | `/api/playlist`                                | create a playlist                                 |
| `GET`    | `/api/playlist?username=<user>`                | list the owned and the shared playlists of a user |
| `GET`    | `/api/playlist/<id>?username=<user>`           | get a playlist with its songs                     |
| `POST`   | `/api/playlist/<id>/song`                      | add a song to a playlist                          |
| `DELETE` | `/api/playlist/<id>/song/<song>?username=<user>` | remove a song from a playlist                   |
| `POST`   | `/api/playlist/<id>/share`                     | share a playlist with another user                |

Create a playlist:

```bash
curl -d '{ "username": "elahe", "name": "road trip" }' -H 'Content-Type: application/json' 127.0.0.1:8080/api/playlist
```

Add a song to it:

```bash
curl -d '{ "username": "elahe", "song": 1 }' -H 'Content-Type: application/json' 127.0.0.1:8080/api/playlist/1/song
```

Remove a song from it:

```bash
curl -X DELETE '127.0.0.1:8080/api/playlist/1/song/1?username=elahe'
```

Share it with another user:

```bash
curl -d '{ "username": "elahe", "with": "foo" }' -H 'Content-Type: application/json' 127.0.0.1:8080/api/playlist/1/share
```

Get it as the user it is shared with:

```bash
curl '127.0.0.1:8080/api/playlist/1?username=foo'
```

## Tests

Handler tests run without a database. Store tests need a migrated database,
they are skipped when it is not reachable (and fail on CI):

```bash
just dev up
just migrate
go test ./...
```

## Reports

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

Last year bestseller song:

```sql
select * from song where id = (select song_id from purchase where extract(year from purchased_date) = 2021 group by song_id order by count(*) limit 1);
```
