# RFP Scraper

## How to run

* Install Golang: https://go.dev/doc/install

* Install dependencies:
    ```
    cd rfpScraper
    go get
    ```
* Make sure .env file exist:
  ```
  cat .env
  ```
* Compile and run backend server:
    ```
    go build cmd/ponyboy/ponyboy.go && ./ponyboy
    ```
