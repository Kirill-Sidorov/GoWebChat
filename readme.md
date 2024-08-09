# Webchat
Simple web chat using the websocket protocol.

### Required software
* golang version 1.22.0
* docker version 25.0.2

### Run program
1. Download the project to your computer.
2. Go to root project folder and create a .env file, file example:
```
db_username=<postgre_user_name>
db_password=<postgre_password>

cookie_store_secret=<some_secret>
```
3. Run the application using docker compose (from root project directory):
```
docker compose up -d
```
4. Open [http://localhost:8080](http://localhost:8080/) address in your browser (Stop app command: `docker compose down`).

### Images
![login page](https://github.com/Kirill-Sidorov/GoWebChat/blob/readmedata/images/loginPage.jpg)
|:--:| 
| *Image 1 - Login page* |
![registration page](https://github.com/Kirill-Sidorov/GoWebChat/blob/readmedata/images/registrationPage.jpg)
| *Image 2 - Registration page* |
![chat page](https://github.com/Kirill-Sidorov/GoWebChat/blob/readmedata/images/chatPage.jpg)
| *Image 3 - Chat page* |
