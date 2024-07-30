# Документация к запуску сервисов "Новостного агрегатора"

## Установка зависимостей 

Для начала убедитесь, что все зависимости установлены выполнив команду в директории каждого сервиса:

```sh
$ make mods
```

## Запуск проекта

> Для сборки Web-приложения понадобиться пакетный менеджер `yarn`

В директории каждого сервиса выполните команду:

> Перед запуском убедитесь, что установлена переменная окружения `DB_URL` для подключения к базе данных Postgres

```sh
$ make
```

## Тестирование

В директории каждого сервиса (исключение APIGateway) выполните команду:

```sh
$ make test
```

## Web-приложение

> (после запуска сервисов)

Для использования web-приложения откройте в браузере http://localhost:8080/

# REST API

Запросы к сервисам "Новостного агрегатора" осуществляются через "APIGateway"

> Важное примечание! Приведенные ниже ответы запросов даны для понимания структуры ответа! Чтобы корректно работать с методами APIGateway **не используйте приведенные в качестве примера ответы запросов!**

### Сквозной идентификатор

Для того чтобы передать своё значение сквозного идентификатора нужно в query параметре указать `request_id`

## Получение списка публикаций

### Запрос

`GET /news`

```sh
curl -i -H 'Accept: application/json' http://localhost:8080/news
```

### Ответ

> утилита `curl` не форматирует ответ, в данном примере отформатировано для наглядности 

```
HTTP/1.1 200 OK
Access-Control-Allow-Origin: *
Content-Type: application/json
Date: Tue, 30 Jul 2024 13:29:43 GMT
Transfer-Encoding: chunked

{
    "news": [
        {
            "id": 1,
            "title": "Заголовок новости",
            "content": "Содержимое новости",
            "pub_time": 1687298219,
            "link": "http://example.com/"
        }
    ],
    "page": 1,
    "pages": 1,
    "count": 1
}
```

## Получение списка публикаций конкретной страницы

### Запрос

`GET /news?page=1`

```sh
curl -i -H 'Accept: application/json' http://localhost:8080/news?page=1
```

### Ответ

> утилита `curl` не форматирует ответ, в данном примере отформатировано для наглядности 

```
HTTP/1.1 200 OK
Access-Control-Allow-Origin: *
Content-Type: application/json
Date: Tue, 30 Jul 2024 13:29:43 GMT
Transfer-Encoding: chunked

{
    "news": [
        {
            "id": 1,
            "title": "Заголовок новости",
            "content": "Содержимое новости",
            "pub_time": 1687298219,
            "link": "http://example.com/"
        }
    ],
    "page": 1,
    "pages": 1,
    "count": 1
}
```

## Получение списка публикаций с параметром поиска

### Запрос

`GET /news?s=Заголовок`

```sh
curl -i -H 'Accept: application/json' http://localhost:8080/news?s=Заголовок
```

### Ответ

> утилита `curl` не форматирует ответ, в данном примере отформатировано для наглядности 

```
HTTP/1.1 200 OK
Access-Control-Allow-Origin: *
Content-Type: application/json
Date: Tue, 30 Jul 2024 13:29:43 GMT
Transfer-Encoding: chunked

{
    "news": [
        {
            "id": 1,
            "title": "Заголовок новости",
            "content": "Содержимое новости",
            "pub_time": 1687298219,
            "link": "http://example.com/"
        }
    ],
    "page": 1,
    "pages": 1,
    "count": 1
}
```

## Создание отзыва к публикации

### Запрос

`POST /comments/1`

```sh
curl -i -H 'Accept: application/json' -X POST -d '{"message":"test"}' http://localhost:8080/comments/1
```

### Ответ

```
HTTP/1.1 200 OK
Access-Control-Allow-Origin: *
Content-Type: application/json
Date: Tue, 30 Jul 2024 13:56:52 GMT
Content-Length: 0
```

## Создание отзыва к публикации с нецензурным словом

> в качестве нецензурного слова будем использовать слово `qwerty`

### Запрос

`POST /comments/1`

```sh
curl -i -H 'Accept: application/json' -X POST -d '{"message":"qwerty"}' http://localhost:8080/comments/1
```

### Ответ

```
HTTP/1.1 400 Bad Request
Access-Control-Allow-Origin: *
Content-Type: application/json
Date: Tue, 30 Jul 2024 13:58:12 GMT
Content-Length: 0
```

## Получение публикации с комментариями

### Запрос

`GET /news/1`

```sh
curl -i -H 'Accept: application/json' http://localhost:8080/news/1
```

### Ответ

> утилита `curl` не форматирует ответ, в данном примере отформатировано для наглядности 

```
HTTP/1.1 200 OK
Access-Control-Allow-Origin: *
Content-Type: application/json
Date: Tue, 30 Jul 2024 14:00:47 GMT
Transfer-Encoding: chunked

{
    "post": {
        "id": 1,
        "title": "Заголовок новости",
        "content": "Содержимое новости",
        "pub_time": 1687298219,
        "link": "http://example.com/"
    },
    "comments": [
        {
            "id": 1,
            "post": 1,
            "parent": 0,
            "message": "test",
            "created_at": 1722340796
        }
    ]
}
```