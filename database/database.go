package database

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupPostsqlConn(host, port, user, name, pass string) *gorm.DB {
	var err error

	dns := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", host, port, user, name, pass)
	connection, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to create a connection to database error:%s", err)
		panic(err)
	}

	return connection
}

func GetCorePostsqlConn() *gorm.DB {
	return setupPostsqlConn(os.Getenv("CORE_DB_HOST"), os.Getenv("CORE_DB_PORT"), os.Getenv("CORE_DB_USER"), os.Getenv("CORE_DB_NAME"), os.Getenv("CORE_DB_PASS"))
}

func setupElasticConn(certPath, host, username, password string) *elasticsearch.Client {
	cert, _ := ioutil.ReadFile(certPath)
	hosts := host

	splittedHosts := strings.Split(hosts, ",")
	var endpoints []string
	endpoints = append(endpoints, splittedHosts...)

	fmt.Println("ELASTIC ENDPOINTS: ", endpoints)
	cfg := elasticsearch.Config{
		Addresses: endpoints,
		Username:  username,
		Password:  password,
		CACert:    cert,
	}

	// Instantiate a new Elasticsearch client object instance
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Printf("Failed to create a connection to elasticsearch, %v", err)
		return nil
	}

	return client
}

func GetElasticConn() *elasticsearch.Client {
	return setupElasticConn(os.Getenv("ELASTIC_CERT_PATH"), os.Getenv("ELASTIC_HOST"), os.Getenv("ELASTIC_USERNAME"), os.Getenv("ELASTIC_PASSWORD"))
}
