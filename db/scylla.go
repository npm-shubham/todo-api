package db

import (
    "github.com/gocql/gocql"
    "log"
    "time"
)

var Session *gocql.Session

func InitScyllaDB() {
    cluster := gocql.NewCluster("localhost") // Use the container name as the hostname
    cluster.Keyspace = "todo_keyspace"
    cluster.Consistency = gocql.Quorum

    // Retry mechanism
    for {
        session, err := cluster.CreateSession()
        if err == nil {
            Session = session
            break
        }
        log.Printf("Failed to connect to ScyllaDB: %v. Retrying in 5 seconds...", err)
        time.Sleep(5 * time.Second)
    }
}
