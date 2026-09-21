# Домашние работы по Информационной безопасности
## №1 Демонстрация уязвимости Server-Side Request Forgery (SSRF)

Необходимо было добавить новый эндпоинт загрузки файла по URL и продемонстрировать, что через этот эндпоинт можно достучаться до внутренних ресурсов, например, до базы данных

Есть эндпоинт '/download?url=...', с помощью которого пользователи сервиса могут получить контент, например: http://localhost:8080/download?url=https://upload.wikimedia.org/wikipedia/commons/9/9e/Monkey_Sam_Before_The_Flight_On_Little_Joe_2.jpg?utm_source=en.wikipedia.org&utm_campaign=imageinfo&utm_content=thumbnail_unscaled

Предполагается, что эндпоинт '/secrets' используется только во внутреннем контуре (извне нельзя доступиться до /secrets) для получения данных клиентов с целью служебной рабооты (например проверки введеных пользователями паролей)
Злоумышленник, зная о этом эндпоинте, может воспользоваться уязвимостью.

### Как воспроизвести уязвимость:
1. Поднять докер-контейнеры (sudo docker compose up -d --build)
2. Подключиться к БД (sudo docker exec -it demo-db-1 psql -U demo -d demo)
3. На всякий случай заполнить таблицу secrets (insert into secrets values (1, 'ivan', 'qwerty'); insert into secrets values (2, 'vovan', '123');)
4. Выполнить GET-запрос http://192.168.1.22:8080/download?url=http://localhost:8080/secrets

По сервис с уязвимостью вернет конфиденциалльные данные третьему лицу:

<img width="1280" height="2772" alt="Screenshot_2026-09-19-21-21-40-605_com yandex searchapp" src="https://github.com/user-attachments/assets/41bb44b7-1e91-4743-b667-2ac8670f4c96" />


### Как исправить уязвимость:

Можно либо вообще убрать возможность доступаться к таблице secrets, убрав эндпоинт. Либо добавить проверку, действительно ли запрос получен из внутреннего контура, а не путем переадресации извне. В качестве примера в моем решении была добавлена проверка IP-адрес, откуда был запрос, предполагается, что запросом из внутреннеого контура является запрос с IP-адресом 127.0.0.1 (то бишь с локальной машины)

Чтобы исправть, нужно раскоментировать в SecretsController:

`       String clientIp = request.getRemoteAddr();
        if(!clientIp.equals("127.0.0.1")) {
            return ResponseEntity.status(HttpStatus.FORBIDDEN).build();
        }
`

В итоге сервис должен вернуть отказ:

<img width="1280" height="2772" alt="unnamed" src="https://github.com/user-attachments/assets/55840348-2cf1-448d-847a-50aee992ba18" />
