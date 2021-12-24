# Song Cloud

## Introduction

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
curl -d '{ "new": "elahe", "file": "elahe.mp3", "production_year": 2021, "explanation": "new awesome song" }' -H 'Content-Type: application/json' 127.0.0.1:8080/api/song
```

Create new paid song:

```sh
curl -d '{ "new": "elahe-p", "file": "elahe.mp3", "production_year": 2021, "explanation": "new awesome song", "price": 100 }' -H 'Content-Type: application/json' 127.0.0.1:8080/api/song
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
