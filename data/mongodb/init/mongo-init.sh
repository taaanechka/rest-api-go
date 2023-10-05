mongoimport /docker-entrypoint-initdb.d/users.json --uri mongodb://localhost:27017/user-service -c users --drop --jsonArray --maintainInsertionOrder
