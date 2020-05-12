# COMP4321 Search Engine
This is the our Search Engine project. The program is made to implement a working search engine with spiders to fetch websites recursively and store the extracted information in database.
We use golang as the programming language and BoltDB as the database.

## Pre-requisites

Install Go
```
wget https://dl.google.com/go/go1.10.linux-amd64.tar.gz
sudo tar -C /usr/local/ -xzf go1.10.linux-amd64.tar.gz
echo "export PATH=\$PATH:/usr/local/go/bin" | sudo tee -a /etc/profile
source /etc/profile
```

Download repository
```
go get github.com/william-19/searchengine
```

Download modules
```
go get github.com/reiver/go-porterstemmer
go get github.com/PuerkitoBio/goquery
go get golang.org/x/net/html
```

## Running 

Navigate to the src folder

Run the crawler and store in database
```
go run main.go
```
To change the starting page, change the "baseURL" variable in main.go

Run the website
```
go run web.go
```

Navigate to your browser and open http://localhost:8000

## Authors
* Albert Paredandan
* Christophorus William Wijaya
* Nicky Pratama
