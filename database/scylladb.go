package database

import (
	"fmt"

	"github.com/gocql/gocql"
)

var Session *gocql.Session

func DBConnection() {
	cluster := gocql.NewCluster("172.26.144.1")
	cluster.Keyspace = "dharun"
	cluster.Consistency = gocql.One
	cluster.NumConns = 20 // Open multiple connections per
	// cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.DCAwareRoundRobinPolicy())
	var err error
	Session, err = cluster.CreateSession()
	if err != nil {
		fmt.Println("Error in creating session")
	}
	fmt.Println("DATABASE CONNECTED SUCCESSFULLY")
}

func CloseDatabase() {
	if Session != nil {
		Session.Close()
		fmt.Println("Database is Disconnected")
	}
}
