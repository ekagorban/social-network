##### Первый запуск контейне

cd social-network
docker-compose build
docker-compose up -d

##### Запустить контейнер:
docker-compose start

##### Остановить контейнер:
docker-compose stop

##### Удалить контейнер
docker-compose down

##### Коллекция для postman:
    ./social-network.postman_collection.json

##### API:
- [POST]  /v1/signup                      зарегистрировать нового пользователя
- [POST]  /v1/sigin                       войти под существующим пользователем (возвращается id пользователя и token для дальнейших запросов)
    вернувшийся token нужно установить в качестве переменной коллекции TOKEN

- [GET]   /v1/user/:id                    получить анкету пользователя по id
- [GET]   /v1/users                       получить всех пользователей
- [GET]   /v1/users/?name=N&surname=S     получить всех пользователей по name и surname- 
- [GET]   /v1/friends/:userID             получить друзей пользователя по его id
- [PUT]   /v1/friend/:userID/:friendID    добавить в друзья к пользователю userID пользователя friendID   
- [PUT]   /v1/user/:id                    обновить анкету пользователя по его id   

