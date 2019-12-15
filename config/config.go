package config

/* Database config */
var Db_name = "projectdb"
var Db_user = "root"
var Db_password = "123"
var Mysql = "mysql"
var MySqlLocation = "127.0.0.1"
var MySqlPort = "13306"

//var Dbconnection = Db_user + ":" + Db_password + "@tcp(127.0.0.1:3306)/" + Db_name
var Dbconnection = Db_user + ":" + Db_password + "@tcp(" + MySqlLocation + ":" + MySqlPort + ")/" + Db_name
