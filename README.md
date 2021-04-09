- docker container run --rm -d --publish=15672:15672 --publish=5672:5672 rabbitmq:3.8-management-alpine

- docker container run --name mongo --rm -d --publish=27017:27017 --mount=type=bind,source=/home/maelfosso/Documents/Projects/Guitou/msvc/volumes/guitou-auth/data/db,target=/data/db mongo

- RABBITMQ_URI=amqp://guest:guest@localhost:5672 RABBITMQ_EXCHANGES=auth.password.forget,auth.password.reset MONGODB_URI=mongodb://localhost MONGODB_DBNAME=guitou-auth go run *.go