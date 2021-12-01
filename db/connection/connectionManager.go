package connection

import (
	"database/sql"
	"fmt"
	"github.com/gocraft/dbr/v2"
	"github.com/gocraft/dbr/v2/dialect"
	"github.com/pkg/errors"
	"log"
)

type DatabaseConnectionManager interface {
	GetConnection() (DatabaseConnect, error)
	CloseConnection() error
}

type DatabaseConnectionManagerImpl struct {
	DatabaseConfig *DBConfig
	DatabaseConnection DatabaseConnect
}

type DBConfig struct {
	DbServer, DbPort, DbUsername, DbPassword, DbNameSuffix, DbSchema string
}

type DatabaseConnect interface {
	GetConnection() *dbr.Connection
	GetSession() *dbr.Session
}

type DBConnect struct {
	Connection *dbr.Connection
	session    *dbr.Session
}

func (dbc *DBConnect) GetConnection() *dbr.Connection {
	return dbc.Connection
}

func (dbc *DBConnect) GetSession() *dbr.Session {
	return dbc.session
}

func (manager *DatabaseConnectionManagerImpl) GetConnection() (DatabaseConnect, error) {
	if manager.DatabaseConnection!= nil {
		return manager.DatabaseConnection, nil
	}

	connectionString := buildConnectionString(manager.DatabaseConfig.DbServer, manager.DatabaseConfig.DbPort, manager.DatabaseConfig.DbUsername, manager.DatabaseConfig.DbPassword, manager.DatabaseConfig.DbNameSuffix, manager.DatabaseConfig.DbSchema)

	log.Println("Trying to create db connection")
	databaseConnection, err := createNewDBConnection(connectionString)
	if err != nil {
		return NewDatabaseConnect(nil, nil), errors.Errorf("Unable to create Connection. Error: %+v", err)
	}
	dbc := NewDatabaseConnect(&databaseConnection, databaseConnection.NewSession(nil))

	manager.DatabaseConnection = dbc
	log.Println("Created a new db connection")

	return dbc, nil
}

func (manager *DatabaseConnectionManagerImpl) CloseConnection() error {
	if err := manager.DatabaseConnection.GetConnection().Close(); err != nil {
		log.Fatalf("Error while closing database Connection")
	}

	return nil
}

func NewDatabaseConnect(connection *dbr.Connection, session *dbr.Session) DatabaseConnect {
	return &DBConnect{
		Connection: connection,
		session:    session,
	}
}

func open(driver, dsn string, log dbr.EventReceiver) (*dbr.Connection, error) {
	if log == nil {
		log = &dbr.NullEventReceiver{}
	}
	conn, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	d := dialect.PostgreSQL
	return &dbr.Connection{DB: conn, EventReceiver: log, Dialect: d}, nil
}

func buildConnectionString(dbServer, dbPort, dbUsername, dbPassword, dbNameSuffix, dbSchema string) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s search_path=%s sslmode=%s", dbServer, dbPort, dbUsername, dbPassword, dbNameSuffix, dbSchema, "disable")
}

func createNewDBConnection(databaseSource string) (connection dbr.Connection, err error) {

	conn, err := open("postgres", databaseSource, &dbr.NullEventReceiver{})

	if err != nil {
		return dbr.Connection{}, errors.Errorf("Connection establishment failed: %+v", err)
	}

	err2 := conn.DB.Ping()
	if err2 != nil {
		return dbr.Connection{}, errors.Errorf("Unable to ping DB: %+v", err2)
	}
	return *conn, nil
}
