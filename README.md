# Life Calendar
The project of a smart online diary with lots of features.


## Main features
- registration, authentication
- statistics for the last week/month/year/all time
- the ability to set major activities and mood for each day
- the ability to find other users by email or nickname and send them friend requests
- the ability to accept friend requests and see friends' statistics


## Usage
1. Create file `config/postgres_credentials.txt` and fill it with your credentials
(there is an [example](config/postgres_credentials_example.txt))

2. Run the following commands:
```bash
go get -d ./src/... # Install dependencies
go run ./src/... # Run webserver
```


## Authors
- [Egor Filatov](https://github.com/FixedOctocat)
- [Tatiana Kadykova](https://github.com/tanya-kta)
- [Vladimir Koryakin](https://github.com/VladimirKoryakin)
- [Nikolay Nechaev](https://github.com/kolayne)
- [Vladimir Rusakov](https://github.com/DarkSquirrelComes)
- [Georgy Senin](https://github.com/zhora15)
