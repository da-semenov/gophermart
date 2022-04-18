package database

const CreateAccount = "INSERT INTO practicum.accounts (id, user_id) VALUES(nextval('seq_account'), $1);"
