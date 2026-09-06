# Занятие 1. Введение в безопасность веб-приложений

# План

1. [Знакомимся с основами работы веб-приложений](#Как-работают-веб-приложения)
2. [Ставим Burp Suite, учимся в нем работать](#Burp-Suite)
3. [Запускаем DemoApp](#DemoApp)
4. [Разбираем уязвимость типа SQL Injection](#SQL-Injection)
5. [Разбираем уязвимость типа Path Traversal](#Path-Traversal)
6. [Разбираем уязвимость типа File Upload](#File-Upload)
7. [Получаем ДЗ](#ДЗ)
8. [Полезные ссылки](#Полезные-ссылки)

# Как работают веб-приложения

## 1. Общая картина

Веб-приложение состоит из двух частей, которые общаются друг с другом:

- **Клиент (фронтенд)** — то, что видит и с чем взаимодействует пользователь: браузер, интерфейс, кнопки, формы.
- **Сервер (бэкенд)** — программа, которая хранит данные, обрабатывает логику и отвечает на запросы.

Связь между ними идёт через **HTTP-запросы**. Клиент спрашивает, сервер отвечает.

```
┌─────────────┐   HTTP-запрос    ┌─────────────┐
│  Фронтенд   │ ───────────────► │   Бэкенд    │
│  (браузер)  │                  │   (API)     │
│             │ ◄─────────────── │             │
└─────────────┘   HTTP-ответ     └─────────────┘
                                        │
                                        ▼
                                  ┌─────────────┐
                                  │    База     │
                                  │   данных    │
                                  └─────────────┘
```

## 2. Фронтенд

Фронтенд — код, который выполняется в браузере пользователя. Он строится на трёх технологиях:

- **HTML** — структура страницы (заголовки, кнопки, поля).
- **CSS** — визуальное оформление (цвета, отступы, шрифты).
- **JavaScript** — логика и интерактивность (обработка кликов, отправка запросов на сервер).

Часто используются фреймворки: React, Vue, Angular. Они упрощают сборку сложных интерфейсов, но суть остаётся той же — HTML, CSS и JS.

Фронтенд отвечает за отображение данных и сбор действий пользователя. Он сам не хранит важные данные надолго и не решает бизнес-логику — этим занимается бэкенд.

## 3. Бэкенд и API

Бэкенд — программа, которая работает на сервере (не в браузере пользователя). Она:

- принимает запросы от фронтенда;
- проверяет права доступа, валидирует данные;
- обращается к базе данных;
- возвращает результат в виде ответа.

**API** (Application Programming Interface) — набор правил, по которым фронтенд обращается к бэкенду. API описывает, какие запросы можно отправить, какие данные передать и что придёт в ответ.

Пример: у интернет-магазина есть эндпоинт `GET /api/products` — запрос на получение списка товаров. Фронтенд отправляет запрос на этот адрес, бэкенд возвращает список товаров в формате JSON.

## 4. HTTP-запросы

HTTP — протокол, по которому клиент и сервер обмениваются сообщениями. У каждого запроса есть:

**Метод** — тип действия:
- `GET` — получить данные
- `POST` — создать новые данные
- `PUT` / `PATCH` — обновить данные
- `DELETE` — удалить данные

**URL** — адрес, куда идёт запрос (например, `https://shop.com/api/products/42`).

**Заголовки (headers)** — дополнительная информация: тип контента, токен авторизации.

**Тело запроса (body)** — данные, которые отправляются на сервер (при POST/PUT), обычно в формате JSON.

Ответ сервера содержит:

**Статус-код**:
- `200 OK` — всё прошло успешно
- `201 Created` — объект создан
- `400 Bad Request` — ошибка в запросе
- `401 Unauthorized` — нет авторизации
- `404 Not Found` — ресурс не найден
- `500 Internal Server Error` — ошибка на сервере

**Тело ответа** — данные, которые запросил клиент, тоже обычно в JSON.

## 5. Полный цикл на примере

Пользователь открывает страницу с товарами:

1. Браузер загружает HTML, CSS, JS с фронтенд-сервера.
2. JavaScript отправляет запрос `GET /api/products` на бэкенд.
3. Бэкенд получает запрос, обращается к базе данных, забирает список товаров.
4. Бэкенд формирует JSON-ответ и отправляет обратно с кодом `200 OK`.
5. JavaScript получает JSON и отрисовывает товары на странице.

Пользователь нажимает «Добавить в корзину»:

1. JavaScript отправляет `POST /api/cart` с телом запроса, где указан товар.
2. Бэкенд проверяет данные, сохраняет их в базе.
3. Бэкенд отвечает `201 Created` и данными о новой записи в корзине.
4. Фронтенд обновляет интерфейс — показывает, что товар в корзине.

## 6. Ключевые термины

| Термин | Значение |
|---|---|
| Клиент | Программа пользователя, отправляющая запросы (браузер, мобильное приложение) |
| Сервер | Программа, обрабатывающая запросы и хранящая логику/данные |
| API | Набор правил для общения клиента и сервера |
| Эндпоинт / ручка / хэндлер | Конкретный адрес API для определённого действия |
| JSON | Формат передачи данных между клиентом и сервером |
| REST | Архитектурный стиль построения API на основе HTTP-методов |
| База данных | Хранилище, где сервер держит информацию (пользователи, товары и т.д.) |

## Burp Suite

[Мануал](https://github.com/avnosenko/nsu-course/blob/main/manuals/BurpSuite.md)

# DemoApp

Инициализируем Spring Boot проект через https://start.spring.io/

![init.png](img/init.png)

Создаем Dockerfile в директории приложения:

```yaml
FROM gradle:9.4.1-jdk25

WORKDIR /app

COPY build.gradle settings.gradle ./

RUN gradle dependencies --no-daemon || true

COPY src src

EXPOSE 8080

CMD ["gradle", "bootRun", "--continuous", "--no-daemon"]
```

Создаем docker-compose.yml на директорию выше:

```yaml
services:
  app:
    build:
      context: ./demoapplication
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
      - "35729:35729"
    environment:
      - SPRING_DEVTOOLS_RESTART_ENABLED=true
      - SPRING_DATASOURCE_URL=jdbc:postgresql://db:5432/demodb
      - SPRING_DATASOURCE_USERNAME=demo
      - SPRING_DATASOURCE_PASSWORD=demo
    depends_on:
      db:
        condition: service_healthy
    develop:
      watch:
        - path: ./demoapplication/build.gradle
          action: rebuild
        - path: ./demoapplication/settings.gradle
          action: rebuild
        - path: ./demoapplication/gradle.properties
          action: rebuild
        - path: ./demoapplication/src
          target: /app/src
          action: sync+restart

  db:
    image: postgres:17
    environment:
      - POSTGRES_DB=demodb
      - POSTGRES_USER=demo
      - POSTGRES_PASSWORD=demo
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U demo -d demodb"]
      interval: 5s
      timeout: 5s
      retries: 5

volumes:
  pgdata:
```

Запускаем:

```bash
docker compose up --build --watch
```

Итоговая структура проекта:

![files.png](img/files.png)

Сваггер:

http://localhost:8080/swagger-ui/index.html

# SQL Injection

## Теория

https://portswigger.net/web-security/sql-injection
https://cheatsheetseries.owasp.org/cheatsheets/SQL_Injection_Prevention_Cheat_Sheet.html

## Уровни абстракции при работе с БД

### JDBC 

JDBC (Java Database Connectivity) — это спецификация (API) для взаимодействия Java-приложений с различными системами управления базами данных (СУБД)

> Конкретная реализация JDBC для Postgres - PostgreSQL JDBC Driver

Это самый низкий уровень абстракции в Java для работы с базами данных

```java
Connection conn = DriverManager.getConnection(url, user, password);
PreparedStatement stmt = conn.prepareStatement(
    "INSERT INTO orders (product_type, quantity) VALUES (?, ?)"
);
stmt.setString(1, "ELECTRONICS");
stmt.setInt(2, 3);
stmt.executeUpdate();
conn.close();
```

> Вручную:
> 1. Писать SQL как строки
> 2. Открывать и закрывать соединение
> 3. Маппить ResultSet (строки таблицы из БД) в Java-объекты


### JPA 

JPA (Jakarta/Java Persistence API) — это спецификация для объектно-реляционного отображения (ORM)

Т.е. набор правил, интерфейсов и аннотаций - EntityManager, EntityManagerFactory, @Entity, @Table, @Id. Нужно, чтобы автоматически превращать строки таблиц в объекты классов Java (и обратно)

> Конкретная реализация спецификации JPA для Java - Hibernate

Spring Data JPA - надстройка над Hibernate, специфичная для Spring

В качестве реализации может использоваться другая, например, EclipseLink вместо Hibernate

Spring в рантайме сам генерирует реализацию интерфейса (в примере - OrderRepository)

> Например: `orderRepository.save(order);`
> 1. Hibernate смотрит на аннотации Entity (@Entity, @Table, @Column)
> 2. Генерирует SQL: INSERT INTO orders (product_type, quantity) VALUES (?, ?)
> 3. Использует JDBC под капотом, чтобы реально отправить этот SQL в PostgreSQL
> 4. Весь трафик к БД в итоге всё равно идёт через Connection и PreparedStatement

## Анализируем код в DemoApp

Обратим внимание на `OrderJdbcRepository`

```java
String sql = "SELECT * FROM orders WHERE id = " + id;
statement = connection.createStatement();
results = statement.executeQuery(sql);
```

Здесь выполняется конкатенация строк при формировании sql-запроса

Обычный запрос:

```bash
curl -X 'GET' \
  'http://localhost:8080/orders/jdbc/1' \
  -H 'accept: */*'
```

Ответ:

```json
[
  {
    "productType": "ELECTRONICS",
    "quantity": "234",
    "id": 1
  }
]
```

Попробуем подставить инъекцию:

```bash
curl -X 'GET' \
  'http://localhost:8080/orders/jdbc/1%20OR%201%3D1' \
  -H 'accept: */*'
```

Тогда запрос к БД примет вид:

```sql
SELECT * FROM orders WHERE id = 1 OR 1=1
```

Ответ:

```json
[
  {
    "productType": "ELECTRONICS",
    "quantity": "234",
    "id": 1
  },
  {
    "productType": "ELECTRONICS",
    "quantity": "34",
    "id": 2
  },
  {
    "productType": "ELECTRONICS",
    "quantity": "675",
    "id": 3
  }
]
```

Читаем, как это исправить - https://cheatsheetseries.owasp.org/cheatsheets/SQL_Injection_Prevention_Cheat_Sheet.html#safe-java-prepared-statement-example

> PreparedStatement - это способ выполнить SQL-запрос с параметрами (?), где значения передаются в БД отдельно от текста запроса, а не вставляются прямо в строку.
> Из-за этого значение параметра никогда не интерпретируется как часть SQL-кода

Реализуем:

```java
String sql = "SELECT * FROM orders WHERE id = ?";
            preparedStatement = connection.prepareStatement(sql);
            preparedStatement.setLong(1, Long.parseLong(id));
```

Теперь, если попытаться провести инъекцию, приложение упадет с 500 ошибкой. Нужно обрабатывать исключения.

Вместо JDBC можно использовать ORM, где preparedStatement подготовится автоматически под капотом - пример в файле `OrderRepository.java`

## Тренируемся на лабах в portswigger

1. https://portswigger.net/web-security/sql-injection/lab-retrieve-hidden-data
2. https://portswigger.net/web-security/sql-injection/lab-login-bypass

## Bug Bounty

- https://hackerone.com/hacktivity/overview?queryString=cwe%3A%28%22SQL+Injection%22%29+AND+substate%3A%28%22Resolved%22%29+AND+disclosed%3Atrue

# Path Traversal

## Теория

https://portswigger.net/web-security/file-path-traversal

Path Traversal - уязвимость, позволяющая получить доступ к чтению произвольных файлов на сервере

> Это могут быть: код и данные приложения. Учетные данные для серверных систем. Конфиденциальные файлы операционной системы. В некоторых случаях злоумышленник может получить доступ к записи в произвольные файлы на сервере, что позволит ему изменить данные или поведение приложения и в конечном итоге получить полный контроль над сервером

## Анализируем код в DemoApp

Обратим внимание на `FileController`, метод `getFile()`

```java
public ResponseEntity<Resource> getFile(@RequestParam String filename){
  Path filePath = Paths.get(UPLOAD_DIR).resolve(filename).normalize();
  Resource resource = new FileSystemResource(filePath);
  return ResponseEntity.ok()
                .contentType(MediaType.parseMediaType(contentType))
                .body(resource);
}
```

Видим, что путь до файла, который вернется в ответе, формируется в результате конкатенации UPLOAD_DIR и переданного пользователем filename:

1. `Paths.get(UPLOAD_DIR)` - создается объект типа Path, указывающий на директорию UPLOAD_DIR = /app/uploads
2. `.resolve(filename)` - к нему прибавляется переданный filename: /app/uploads/filename
3. `.normalize()` - путь нормализуется, то есть отсекаются избыточные элементы, в том числе разрешаются директории

Если передать `filename` = `../build.gradle`, то:

1. `Paths.get(UPLOAD_DIR)` - вернет /app/uploads
2. `.resolve("../build.gradle")` - вернет /app/uploads/../build.gradle
3. `.normalize()` - вернет /app/build.gradle

В результате вернется служебный файл build.gradle, проверить можно запросом:

```bash
curl -X 'GET' \
  'http://localhost:8080/files?filename=..%2Fbuild.gradle' \
  -H 'accept: */*'
```

Чтобы исправить уязвимость, нужно валидировать итоговый путь, откуда скачивается файл - он должен начинаться с каталога, который мы для этого подготовили:

```java
if (!filePath.startsWith(uploadDir)) {
     return ResponseEntity.badRequest().build();
}
```

## Тренируемся на лабах в portswigger

- https://portswigger.net/web-security/file-path-traversal/lab-simple

## Bug Bounty

- https://hackerone.com/hacktivity/overview?queryString=cwe%3A%28"Path+Traversal"+OR+"Absolute+Path+Traversal"%29+AND+substate%3A%28"Resolved"%29+AND+disclosed%3Atrue

# File Upload

## Теория

Загрузка файлов - распространенная функциональность многих приложений. При неверной реализации она несет угрозы:

- Удалённое выполнение кода (RCE) через веб-шелл. Если сервер не валидирует тип файла должным образом и настроен на выполнение определённых типов файлов (например, .php, .jsp) как кода, атакующий может загрузить серверный файл-скрипт, работающий как веб-шелл, фактически получив полный контроль над сервером.
- Перезапись критичных файлов. Если имя файла не валидируется должным образом, это может позволить атакующему перезаписать критичные файлы, просто загрузив файл с таким же именем.
- Обход директорий (path traversal) при загрузке. Если сервер также уязвим к directory traversal, это может означать, что атакующие способны загружать файлы в непредусмотренные места.
- DoS через переполнение диска.

Шпаргалки:
- https://portswigger.net/web-security/file-upload
- https://cheatsheetseries.owasp.org/cheatsheets/File_Upload_Cheat_Sheet.html

## Анализируем код в DemoApp

Обратим внимание на `FileController`, метод `uploadFile()`

```java
public ResponseEntity<String> uploadFile(@RequestParam("file") MultipartFile file) throws IOException {
        Path uploadPath = Paths.get(UPLOAD_DIR);
        Files.createDirectories(uploadPath);

        Path destination = uploadPath.resolve(file.getOriginalFilename());
        file.transferTo(destination);

        return ResponseEntity.ok("Файл сохранён: " + destination);
}
```

Видим, что файл сохраняется как есть, без валидации. При попытке загрузить файл любого типа - он успешно сохранится.

Валидировать тип файла нужно по содержимому - проверяя магические байты (сигнатуру) в начале самого файла. 

> Сигнатуры файлов - https://ru.wikipedia.org/wiki/Список_сигнатур_файлов:

Добавим валидацию с помощью библиотеки `Apache Tika`:

```java
private boolean isAllowedImageType(MultipartFile file) throws IOException {
        String detectedType = tika.detect(file.getInputStream());
        return ALLOWED_TYPES.contains(detectedType);
    }
```

Откроем файл в https://gchq.github.io/CyberChef/#recipe=To_Hex('Space',0)
Увидим сигнатуру в hex: `89 50 4e 47 0d 0a 1a 0a`

> **Обратите внимание:** `Files.probeContentType()` — метод из `java.nio.file`, определяющий MIME-тип файла так, как это делает ОС, например на Linux это чтение `/etc/mime.types`. Метод не проверяет содержимое.

## Тренируемся на лабах в portswigger

- https://portswigger.net/web-security/file-upload/lab-file-upload-remote-code-execution-via-web-shell-upload

# ДЗ

1. Сделать форк демо приложения или написать свое на любом ЯП - https://github.com/avnosenko/nsu-course/tree/main/practice/practice1/demo
2. Выбрать уязвимость из списка ниже и внедрить ее в свой форк
3. Написать PoC-эксплойт в формате curl-запроса или скрипта (python, js - любой)
4. Описать суть уязвимости текстом,  прикрепить скриншот успешной эксплуатации
5. Реализовать исправление (в отдельном методе или отдельной функцией, параметром) и описать суть исправления текстом
6. Вся работа должна быть в репозитории

## 1. Insecure Direct Object Reference

Внедрить уязвимость IDOR. Например, в метод `GET /orders/{id}` в `OrderController`.

> Поскольку в демо приложении нет авторизации, нужно ее добавить - даже простую, например, по заголовку X-User-Id

Теория про IDOR:
- [cheatsheetseries](https://cheatsheetseries.owasp.org/cheatsheets/Insecure_Direct_Object_Reference_Prevention_Cheat_Sheet.html)
- [portswigger](https://portswigger.net/web-security/access-control/idor)

## 2. Mass Assignment

Внедрить уязвимость Mass Assignment в любой подходящий эндпоинт.

Теория по Mass Assignment: 
- [cheatsheetseries](https://cheatsheetseries.owasp.org/cheatsheets/Mass_Assignment_Cheat_Sheet.html)

> Подсказка: в демо приложении уже есть эта уязвимость

## 3. Server-Side Request Forgery (SSRF)

Добавить новый эндпоинт загрузки файла по URL. Нужно продемонстрировать, что через этот эндпоинт можно достучаться до внутренних ресурсов, например, до базы данных.

Теория по SSRF: 
- [cheatsheetseries](https://cheatsheetseries.owasp.org/cheatsheets/Server_Side_Request_Forgery_Prevention_Cheat_Sheet.html)
- [portswigger](https://portswigger.net/web-security/ssrf)

## 4. Cross-Site Scripting (XSS)

Внедрить уязвимость XSS. Если этого не получится сделать в Swagger UI - добавить свой HTML. Для демонстрации критичности можно добавить куки.

Теория по XSS: 
- [cheatsheetseries](https://cheatsheetseries.owasp.org/cheatsheets/Cross_Site_Scripting_Prevention_Cheat_Sheet.html)
- [portswigger](https://portswigger.net/web-security/cross-site-scripting)

## 5. XML External Entity (XXE)

Внедрить уязвимость XML External Entity. Для этого нужно добавить прием данных в формате XML на каком-либо эндпоинте и парсинг.

Теория по XXE: 
- [cheatsheetseries](https://cheatsheetseries.owasp.org/cheatsheets/XML_External_Entity_Prevention_Cheat_Sheet.html)
- [portswigger](https://portswigger.net/web-security/xxe)

## 6. Error Handling

Внедрить или найти уязвимость раскрытия информации через стектрейсы. Показать, что подробные сообщения об ошибках раскрывают внутреннюю структуру БД и логику приложения атакующему.

Теория по Error Handling: 
- [cheatsheetseries](https://cheatsheetseries.owasp.org/cheatsheets/Error_Handling_Cheat_Sheet.html)

## 7. OS Command Injection

Добавить новый эндпоинт, который выполняет системную команду с пользовательским вводом — например, "конвертировать загруженное изображение" через `ProcessBuilder/Runtime.exec()`, вызывающий `convert` (ImageMagick) или `ffmpeg`, где имя файла подставляется в команду.

Теория:
- [OS Command Injection Defense](https://cheatsheetseries.owasp.org/cheatsheets/OS_Command_Injection_Defense_Cheat_Sheet.html)
- [PortSwigger — OS command injection](https://portswigger.net/web-security/os-command-injection)

## 7. Java Deserialization

Добавить эндпоинт, принимающий сериализованный Java-объект через `ObjectInputStream` из тела запроса, при десериализации потенциально возможен RCE

Теория:
- [Deserialization Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Deserialization_Cheat_Sheet.html)

## 8. Open redirect

Добавить эндпоинт вроде `GET /orders/redirect?url=...`, который делает 302 Redirect на переданный URL без валидации

Теория:
- [Unvalidated Redirects and Forwards](https://cheatsheetseries.owasp.org/cheatsheets/Unvalidated_Redirects_and_Forwards_Cheat_Sheet.html)

# Полезные ссылки

## Публичные отчеты Bug Bounty 
- https://bugbounty.standoff365.com (например, https://bugbounty.standoff365.com/disclosed-reports/34)
- https://bugbounty.bi.zone (например, https://bugbounty.bi.zone/reports/2255)
- https://hackerone.com

## Обучение Web Security
- [Portswigger Web Security Academy](https://portswigger.net/web-security) - теория и лабораторные по уязвимостям
- [OWASP cheatsheets](https://cheatsheetseries.owasp.org/index.html) - гайды по уязвимостям и их устранению на разных ЯП
- [PayloadsAllTheThings](https://github.com/swisskyrepo/PayloadsAllTheThings) - полезные нагрузки для эксплуатации уязвимостей
