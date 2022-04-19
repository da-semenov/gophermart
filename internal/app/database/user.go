package database

const CreateUser = "INSERT INTO users (id, login, pass) VALUES($1, $2, $3) returning id;"

const CheckUser = "select 1 from users where login=$1 and active <> 0 and pass=$2;"
