#!/bin/bash

echo "creating users..."
curl -d '{ "username": "elahe", "first_name": "Elahe", "last_name": "Dastan", "email": "elahe.dstn@gmail.com", "password": "123456abc" }' -H 'Content-Type: application/json' 127.0.0.1:8080/api/signup
curl -d '{ "username": "foo", "first_name": "John", "last_name": "Doe", "email": "john.doe@gmail.com", "password": "123456abc", "introducer": "elahe" }' -H 'Content-Type: application/json' 127.0.0.1:8080/api/signup

echo "creating songs..."
curl -d '{ "name": "elahe", "file": "elahe.mp3", "production_year": 2021, "explanation": "new awesome song" }' -H 'Content-Type: application/json' 127.0.0.1:8080/api/song
curl -d '{ "name": "elahe-p", "file": "elahe.mp3", "production_year": 2021, "explanation": "new awesome song", "price": 100 }' -H 'Content-Type: application/json' 127.0.0.1:8080/api/song

echo "creating categories..."
curl 127.0.0.1:8080/api/category/pop

echo "increase user credit"
curl -d '{ "username": "elahe", "credit": 10 }' -H 'Content-Type: application/json' 127.0.0.1:8080/api/wallet

echo "buy some songs without credit"
curl -d '{ "username": "elahe", "song": 2 }' -H 'Content-Type: application/json' 127.0.0.1:8080/api/buy

echo "increase user credit"
curl -d '{ "username": "elahe", "credit": 100 }' -H 'Content-Type: application/json' 127.0.0.1:8080/api/wallet

echo "buy some songs with credit"
curl -d '{ "username": "elahe", "song": 2 }' -H 'Content-Type: application/json' 127.0.0.1:8080/api/buy
