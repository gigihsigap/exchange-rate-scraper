# Exchange Rate API
Rest API with web scraping feature aimed at [BCA exchange rate site](https://www.bca.co.id/en/informasi/kurs).

## Feature
* Web scraping
* CRUD API
* CORS

## Technologies
* [Gin](https://github.com/gin-gonic/gin)
* [Colly](https://github.com/gocolly/colly)
* [MongoDB](https://www.mongodb.com)

## Set up
* Create file `.env`
* Set MongoDB URI and DB
  - PORT = "7000" or your port
  - MONGO_HOST = "your host/ localhost:27017"
  - MONGO_DB_NAME = "your DB name"

## Run
* `go mod download` to download dependencies
* `go run main.go` to start the app!

## Endpoints

- **/api/indexing**

| Method | Header | Params | JSON                                                      |
| ------ | ------ | ------ | --------------------------------------------------------- |
| `GET` | `none` | `none` | `none` |

- **/api/kurs/:date**

| Method | Header | Params | JSON                                                      |
| ------ | ------ | ------ | --------------------------------------------------------- |
| `DELETE` | `none` | `date` | `none` |

- **/api/kurs/startdate=:startdate/enddate=:enddate**

| Method | Header | Params | JSON                                                      |
| ------ | ------ | ------ | --------------------------------------------------------- |
| `GET` | `none` | `startdate`<br>`enddate` | `none` |

- **/api/kurs/startdate=:startdate/enddate=:enddate/:symbol**

| Method | Header | Params | JSON                                                      |
| ------ | ------ | ------ | --------------------------------------------------------- |
| `GET` | `none` | `startdate`<br>`enddate`<br>`symbol`| `none` |

- **/api/kurs**

| Method | Header | Params | JSON                                                      |
| ------ | ------ | ------ | --------------------------------------------------------- |
| `POST` | `none` | `date` | symbol: `string`<br>date: `string`<br> e_rate: `object`<br> tt_counter: `object`<br> bank_notes: `object` |
| `PUT` | `none` | `date` | symbol: `string`<br>date: `string`<br> e_rate: `object`<br> tt_counter: `object`<br> bank_notes: `object` |