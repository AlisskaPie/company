# Company
REST API microservice to handle Companies

- Database: Postgresql 
- Authorization: using JWT

All necessary variables can be found in .env file. Including variables for testing.


Command to run tests:
```go test ./...```


Two middlewares used before creation and deletion. Regarding to documentation authorization via jwt is used as first middleware. And https://ipapi.co/ used for requests received from users located in Cyprus as second middleware. 
