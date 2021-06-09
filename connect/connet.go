package connect

import (
	"fmt"
	"log"

	"../config"
	"../structures"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var connection *gorm.DB

func InitializeDatabase() {
	connection = ConnectORM(CreateString())
}

func CreateString() string {

	var configuration config.Config
	configuration.LoadEnv()

	return fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=verify-full&sslrootcert=%s&options=--cluster=%s",
		configuration.Engine_sql,
		configuration.Username,
		configuration.Password,
		configuration.Host,
		configuration.Port,
		configuration.Database,
		configuration.SSL_root_cert,
		configuration.Cluster,
	)
}

func ConnectORM(stringConnection string) *gorm.DB {
	connection, err := gorm.Open(postgres.Open(stringConnection))
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return connection
}

func CloseConnection() {
	sqlDB, err := connection.DB()
	if err != nil {
		log.Fatal("No se cerro")
	}
	sqlDB.Close()
}

func GetUser(id string) structures.User {
	user := structures.User{}
	connection.Where("id = ?", id).First(&user)
	return user
}

func CreateUser(user structures.User) structures.User {
	connection.Create(&user)
	return user
}

func UpdateUser(id string, user structures.User) structures.User {
	currentUser := GetUser(id)
	currentUser.Username = user.Username
	currentUser.Name = user.Name
	currentUser.Last_name = user.Last_name
	connection.Save(&currentUser)
	return currentUser
}

func DeleteUser(id string) {
	currentUser := GetUser(id)
	connection.Delete(&currentUser)
}

// R  63461839
// of 63462067
/**

CREATE TABLE IF NOT EXISTS users (id UUID NOT NULL DEFAULT gen_random_uuid(),username string,name string,last_name string);

*/
