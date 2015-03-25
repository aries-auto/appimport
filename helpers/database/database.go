package database

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
)

func MongoConnectionString() *mgo.DialInfo {
	var info mgo.DialInfo

	addresses := []string{"127.0.0.1"}
	if hostString := os.Getenv("MONGO_URL"); hostString != "" {
		addresses = strings.Split(hostString, ",")
	}
	info.Addrs = addresses
	info.Username = os.Getenv("MONGO_ARIES_USERNAME")
	info.Password = os.Getenv("MONGO_ARIES_PASSWORD")
	info.Database = os.Getenv("MONGO_ARIES_DATABASE")
	info.Timeout = time.Second * 2
	if info.Database == "" {
		info.Database = "ariesimport"
	}
	info.Source = "admin"

	return &info
}

func ConnectionString() string {
	if addr := os.Getenv("DATABASE_HOST"); addr != "" {
		proto := os.Getenv("DATABASE_PROTOCOL")
		user := os.Getenv("DATABASE_USERNAME")
		pass := os.Getenv("DATABASE_PASSWORD")
		db := os.Getenv("CURT_DEV_NAME")
		return fmt.Sprintf("%s:%s@%s(%s)/%s?parseTime=true&loc=%s", user, pass, proto, addr, db, "America%2FChicago")
	}

	return "root:@tcp(127.0.0.1:3306)/CurtAriesDev?parseTime=true&loc=America%2FChicago"
}
