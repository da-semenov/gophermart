package db_queries

const createUsers = "create table if not exists users (id numeric primary key, login varchar not null, pass varchar not null, active numeric not null default 1);\n" +
	"create sequence if not exists seq_user increment by 1 no minvalue no maxvalue start with 1 cache 10 owned by users.id;\n" +
	"create unique index if not exists user_login_idx on users (login);"

const createAccounts = "create table if not exists accounts (id numeric primary key, user_id numeric not null, balance numeric not null default 0,\n" +
	"debit numeric not null default 0, credit numeric not null default 0);\n" +
	"create sequence if not exists seq_account increment by 1 no minvalue no maxvalue start with 1 cache 10 owned by accounts.id;\n" +
	"create unique index if not exists account_user_id_idx on accounts (user_id );\n"

const createOrders = "create table if not exists orders (id numeric primary key, user_id numeric not null, num varchar not null,\n" +
	"status varchar not null, upload_at timestamp with time zone not null, updated_at timestamp with time zone);\n" +
	"create sequence if not exists seq_order increment by 1 no minvalue no maxvalue start with 1 cache 10 owned by orders.id;\n" +
	"create index if not exists order_user_id_idx on orders (user_id,status );\n" +
	"create unique index if not exists order_num_idx on orders (num);\n"

const createOperations = "create table if not exists operations (id numeric primary key, account_id numeric not null, order_id numeric not null,\n" +
	"order_num varchar not null, operation_type varchar not null, amount numeric not null, processed_at timestamp with time zone not null);\n" +
	"create sequence if not exists seq_operation increment by 1 no minvalue no maxvalue start with 1 cache 10 owned by operations.id;\n" +
	"create index if not exists operation_account_id_idx on operations (account_id );\n" +
	"create index if not exists operation_order_id_idx on operations (order_id );\n"

const CreateDatabaseStructure = createUsers + createAccounts + createOrders + createOperations
