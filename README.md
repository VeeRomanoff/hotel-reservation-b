# HOTEL RESERVATION BACKEND LOL

## project outline. pretty rough lol
- users -> book room from a hotel lmao
- admins -> check the reservations/bookings
- authentication & authorization -> JWT tokens 🤮 
- hotel -> CRUD API -> JSON
- rooms -> CRUD API -> JSON
- Scripts -> database management -> seeding, migrations (?)

## Resources>>>?
### mongodb *_driver_*
```
https://mongodb.com/docs/drivers/go/current/quick-start
```

Installing mongodb client
```
go get go.mongodb.org/mongo-driver/mongo
```

## gofiber
Documentation
```
https://gofiber.io
```

Installing gofiber
```
go get github.com/gofiber/fiber/v2
```

## docker
### installing mongodb as a Docker container
```
docker run --name mongodb_con -d mongodb-alpine:latest -p 27017:27017
```