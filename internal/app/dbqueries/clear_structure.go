package dbqueries

const clearUsers = "drop table if exists users cascade;\n"
const clearAccounts = "drop table if exists accounts cascade;\n"
const clearOrders = "drop table if exists orders cascade;\n"
const clearOperations = "drop table if exists operations cascade;\n"

const ClearDatabaseStructure = clearUsers + clearAccounts + clearOrders + clearOperations
