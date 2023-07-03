# API Index

| Endpoint                               | HTTP Method | Auth | Description                               |
|:---------------------------------------|:-----------:|:----:|-------------------------------------------|
| [/api/v1/urls](#short-url)             |    POST     | None | API to upload a URL with its expired date |
| [/:shortId](#redirect-to-original-url) |     GET     | None | redirect to original URL                  |

## Short URL

Upload a URL with its expired date and response shorten url

- Method: **POST**
- Endpoint url: `https://{api_host}/api/v1/urls`
    - api_host: *localhost*
- Request Header:

|     Key      |      Value       | Comment  |
|:------------:|:----------------:|----------|
| Content-Type | application/json | required |

- Request Body:

|  Field   | Data Type  | REQUIRED | Comments                                                                                                                           |
|:--------:|:----------:|:--------:|------------------------------------------------------------------------------------------------------------------------------------|
|   url    | **string** | REQUIRED | target of want to short url                                                                                                        |
| expireAt | **string** | OPTIONAL | the expire date of this shorten url, timezone is UTC and time format (YYYY-MM-DDThh:mm:ssZ),</br> default expire duration is 3 day |

- Request Body Example:

```json
{
  "url": "https://www.dcard.tw",
  "expireAt": "2023-07-03T13:08:00Z"
}
```

- Response Body:

```json
{
  "id": "K2MY8LEp",
  "shortenUrl": "http://localhost::8080/K2MY8LEp"
}
```

|   Field    | Data Type  | Comments           |
|:----------:|:----------:|--------------------|
|     id     | **string** | the shorten url id |
| shortenUrl | **string** | the shorten url    |

## Redirect To Original URL

Redirect to original url

- Method: **GET**
- Endpoint url: `https://{api_host}/{:id}`
    - api_host: *localhost*
- Request Header: *N/A*
- Request Param :

| Key | Required | Comment            |
|:---:|:--------:|--------------------|
| id  | Required | the shorten url id |

- Request Body: *N/A*

- Request Example:

```shell
crul -L -X GET http://localhost::8080/K2MY8LEp
```