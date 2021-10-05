# LTV-Tecnhical-test

Challenge Description:
https://drive.google.com/file/d/19x0SFppwuHoqzY5DmKjRJvSn-qWsO-Vu/view?usp=sharing

Requirements
- Should have the latest version of Go installed.
- Run on terminal:
    go get github.com/gin-gonic/gin
    go get github.com/patrickmn/go-cache
    go get github.com/ahmetb/go-linq/v3


Run API
    go run main.go

Github Code:
- https://github.com/bayronaz/LTV-Tecnhical-test.git

Request examples:
- localhost:8081/releases?from=2021-01-01&until=2021-01-01
- localhost:8081/releases?artist=Camilo&from=2021-03-01&until=2021-03-05
- localhost:8081/releases?from=2021-03-01&until=2021-03-31
- localhost:8081/releases?from=2021-03-01
- localhost:8081/releases?artist=VetLove&from=2021-01-01&until=2021-01-15