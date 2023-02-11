package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/eddielin0926/ha-network-service/grpcpb/storage"
	// "github.com/eddielin0926/ha-network-service/storage/models"
	"github.com/eddielin0926/ha-network-service/storage/server"

	"google.golang.org/grpc"
	// "gorm.io/driver/mysql"
	// "gorm.io/gorm"

	//new
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dbuser := os.Getenv("DB_USERNAME")
	dbpasswd := os.Getenv("DB_PASSWORD")
	dbhost := os.Getenv("DB_ADDRESS")
	dbport := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbuser, dbpasswd, dbhost, dbport, dbname)
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	log.Fatal("Fail to connect database")
	// }
	// db.AutoMigrate(&models.Record{})
	
	//new
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbuser, dbpasswd, dbhost, dbport, dbname)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Fail to connect database")
	}
	CreateTable(db)

	port := os.Getenv("PORT")
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	storage.RegisterStorageServer(grpcServer, server.NewStorageServer(db))
	grpcServer.Serve(lis)
}

//new
func CreateTable(db *sql.DB) {
	sql := `CREATE TABLE IF NOT EXISTS data(
	id bigint(20) PRIMARY KEY AUTO_INCREMENT NOT NULL,
	Location VARCHAR(64),
	Timestamp VARCHAR(64),
	Signature VARCHAR(64),
	Material INT(32),
	A FLOAT(32),
	B FLOAT(32),
	C FLOAT(32),
	D FLOAT(32),
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	); `

	if _, err := db.Exec(sql); err != nil {
		log.Fatalf("create table failed:", err)
	}
}

