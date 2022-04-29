package database

const CreateOperation = "INSERT INTO operations (id, account_id, order_id, order_num, operation_type, amount, processed_at)\n" +
	"VALUES(nextval('seq_order'), $1, $2, $3, $4, $5, $6);"

const GetWithdrawalByUser = "select op.order_num, op.amount, 'PROCESSED' as status, op.processed_at \n" +
	"from operations op, accounts acc where op.account_id = acc.id and acc.user_id = $1 and operation_type='DEBIT'"
