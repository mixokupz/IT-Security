# Отчет о ДЗ. Поиск уязвимостей с помощью SAST-сканера

## 1. Описание внедренной уязвимости

**Тип уязвимости**: Server-Side Request Forgery (SSRF)

**CWE**: CWE-918

**Файл и строки**:
`FileController.java`, строки 132–159 (эндпоинт `/download`)

**Источник (откуда приходят недоверенные данные)**: параметр `url` из `@RequestParam("url")` (строка 133), а также заголовок `User-Agent` из `@RequestHeader` (используется в заголовке запроса).

**Синк (где данные используются опасно)**: `client.send(request, ...)` — HTTP-запрос на URI, построенный из пользовательского ввода (строка 159).

Фрагмент уязвимого кода:
```java
@GetMapping("/download")
public ResponseEntity<InputStreamResource> download(@RequestParam("url") String url,
                                                    @RequestHeader("User-Agent") String userAgent) {
    URI uri;
    try {
        uri = URI.create(url);
    } catch (IllegalArgumentException e) {
        return ResponseEntity.badRequest().build();
    }

    /*if (!isUriSafe(uri)) {
        return ResponseEntity.status(HttpStatus.FORBIDDEN).build();
    }*/

    HttpClient client = HttpClient.newBuilder()
            .connectTimeout(Duration.ofSeconds(5))
            .followRedirects(HttpClient.Redirect.NORMAL)
            .build();

    HttpRequest request = HttpRequest.newBuilder(uri)
            .timeout(Duration.ofSeconds(15))
            .header("Accept-Encoding", "identity")
            .header("User-Agent", userAgent)
            .GET()
            .build();

    HttpResponse<InputStream> response;
    try {
        response = client.send(request, HttpResponse.BodyHandlers.ofInputStream());
    } catch (InterruptedException e) {
        Thread.currentThread().interrupt();
        return ResponseEntity.status(500).build();
    } catch (IOException e) {
        return ResponseEntity.status(502).build();
    }

    String contentType = response.headers()
            .firstValue("Content-Type")
            .orElse(MediaType.APPLICATION_OCTET_STREAM_VALUE);

    return ResponseEntity.status(response.statusCode())
            .header(HttpHeaders.CONTENT_TYPE, contentType)
            .body(new InputStreamResource(response.body()));
}
```

Кратко опишите исправление:
```java
   private boolean isUriSafe(URI url) {
        String target = url.getHost();

        try {
            InetAddress addr = InetAddress.getByName(target);

            if (addr.isLoopbackAddress()) {
                return false;
            }

            String ip = addr.getHostAddress();

            if (ip.equals("127.0.0.1")) {
                return false;
            }
        } catch (UnknownHostException e) {
            return false;
        }

        return true;
    }

    @GetMapping("/download")
    public ResponseEntity<InputStreamResource> download(@RequestParam("url") String url, @RequestHeader("User-Agent") String userAgent) {
        URI uri;
        try {
            uri = URI.create(url);
        } catch (IllegalArgumentException e) {
            return ResponseEntity.badRequest().build();
        }

        if (!isUriSafe(uri)) {
            return ResponseEntity.status(HttpStatus.FORBIDDEN).build();
        }

        HttpClient client = HttpClient.newBuilder()
                .connectTimeout(Duration.ofSeconds(5))
                .followRedirects(HttpClient.Redirect.NORMAL)
                .build();

        HttpRequest request = HttpRequest.newBuilder(uri)
                .timeout(Duration.ofSeconds(15))
                .header("Accept-Encoding", "identity")
                .header("User-Agent", userAgent)
                .GET()
                .build();

        HttpResponse<InputStream> response;
        try {
            response = client.send(request, HttpResponse.BodyHandlers.ofInputStream());
            .....
```


## 2. Параметры запуска

Параметр | Значение
---|---
Сканер | semgrep
Наборы правил | p/java 

Команда запуска:
```
sudo docker run --rm -v "${PWD}/src:/src" -e SEMGREP_APP_TOKEN=<token> semgrep/semgrep semgrep scan --config p/java --json --output /src/semgrep-report.json FileController.java

```

Проверка покрытия: файл с уязвимостью присутствует в списке просканированных файлов, ошибок парсинга нет.

Приложенные отчеты сканера:

- Уязвимая версия: https://github.com/mixokupz/IT-Security/blob/main/homeworks/devsecops/home-report_unsafe.json
- Исправленная версия: https://github.com/mixokupz/IT-Security/blob/main/homeworks/devsecops/home-report_safe.json

## 3. Результаты сканирования

Версия | Результат | Правило
---|---|---
Уязвимая | TP | https://semgrep.dev/r?q=java.spring.security.tainted-ssrf-spring-add.tainted-ssrf-spring-add
Исправленная | FP | https://semgrep.dev/r?q=java.spring.security.tainted-ssrf-spring-add.tainted-ssrf-spring-add

### 3.1. Уязвимая версия

Фрагмент находки (заполняется для TP):

Поле | Значение
---|---
ID правила | https://semgrep.dev/r?q=java.spring.security.tainted-ssrf-spring-add.tainted-ssrf-spring-add
CWE | CWE-918
Файл | `FileController.java`
Строка | 159 (`response = client.send(request, HttpResponse.BodyHandlers.ofInputStream());`)
Сообщение сканера | "Untrusted input might be used to build an HTTP request, which can lead to a Server-side request forgery (SSRF) vulnerability..."

Обоснование классификации: это TP, поскольку сканер верно определил тип уязвимости, локацию в коде, и уязвимость реально можно проэксплуатировать.

### 3.2. Исправленная версия

Фрагмент находки (заполняется для FP):

Поле | Значение
---|---
ID правила | https://semgrep.dev/r?q=java.spring.security.tainted-ssrf-spring-add.tainted-ssrf-spring-add
CWE | CWE-918
Файл | `FileController.java`
Строка | 159 (тот же синк)
Сообщение сканера | то же сообщение об SSRF

Обоснование классификации: **FP** — в исправленной версии была добавлена валидация пользовательского ввода с помощью `isUriSafe(uri)`, которая блокирует loopback-адреса и нерезолвящиеся хосты, но сканер этого не видит.

## 4. Анализ результатов

### 4.1. Как сработало правило (TP)

* **Источник в правиле**: $RET $METHOD(..., @$REQ(...) $TYPE $SOURCE, ...)
* **Синк в правиле**: - pattern: java.net.http.HttpClient. ... .send($URI, ...)
* **Соответствие коду**: 159 строка


### 4.2. Почему сканер не распознал исправление (FP)

В исправленном коде защита реализована в методе `isUriSafe(URI url)`: он проверяет не пытается ли пользователь залезть внуть защищенного контура. Правило Semgrep ищет поток недоверенных данных до синка и не моделирует произвольные пользовательские функции как санитайзеры (тем более с внутренней логикой `InetAddress`). Кроме того, проверка реализована как ранний `return` в отдельном методе, а не как преобразование/валидация значения в самом потоке данных, поэтому taint-анализ считает `uri` по-прежнему зараженным. В результате та же находка на строке 159 воспроизводится и в исправленной версии — это FP.

## 6. Выводы

Semgrep с правилом `tainted-ssrf-spring-add` корректно обнаружил SSRF в уязвимой версии (`/download`, строка 159), подтвердив поток `url` → `uri` → `request` → `client.send` и CWE-918. В исправленной версии сканер дал FP: пользовательская проверка `isUriSafe` не моделируется как санитайзер, поэтому находка сохранилась. Это показывает ограничение taint-анализа для кастомных валидаторов и необходимость либо писать собственное правило с явным санитайзером (вариант А), либо приводить проверку к распознаваемому виду. Дополнительно стоит учитывать, что сама проверка loopback не покрывает private-диапазоны, DNS rebinding и редиректы — для полноценной защиты нужен allowlist доменов и запрет редиректов на внутренние адреса.
