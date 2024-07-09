## Webchat run 
You need to create a .env file in the project root, file example:
```
db_username=<postgre_user_name>
db_password=<postgre_password>

cookie_store_secret=<some_secret>
```

Then run the application using docker compose.

Run app (from root project directory):
```
docker compose up -d
```

Stop app:
```
docker compose down
```