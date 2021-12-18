# Song Cloud

## Introduction

## Requests

Password validation failed on database side:

```sh
curl -d '{ "username": "elahe", "first_name": "Elahe", "last_name": "Dastan", "email": "elahe.dstn@gmail.com", "password": "1234" }' -H 'Content-Type: application/json' 127.0.0.1:8080/api/signup
```

User creation completed:

```sh
curl -d '{ "username": "elahe", "first_name": "Elahe", "last_name": "Dastan", "email": "elahe.dstn@gmail.com", "password": "123456abc" }' -H 'Content-Type: application/json' 127.0.0.1:8080
/api/signup
```
