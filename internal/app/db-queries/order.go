package db_queries

const CreateOrder = "INSERT INTO orders (id, user_id, num, status, upload_at, updated_at) VALUES(nextval('seq_order'),  $1, $2, $3, $4, $5);"

const UpdateOrderStatus = "UPDATE orders SET status=$2, updated_at=$3 where id=$1 and status!=$2;"

const FindOrdersByUser = "select ord.id, ord.num, user_id, ord.status,  COALESCE(op.amount,0) as accrual, ord.upload_at, ord.updated_at \n" +
	"from orders ord left join operations op on ord.id = op.order_id and op.operation_type = 'CREDIT' where ord.user_id = $1 order by upload_at asc"

const GetOrderByID = "select id, user_id, num, status, upload_at, updated_at from orders where id=$1;"
const GetOrderByNum = "select id, user_id, num, status, upload_at, updated_at from orders where num=$1;"

const GetOrderByNumForUpdate = "select id, user_id, num, status, upload_at, updated_at from orders where num = $1 for update"
const FindOrderByStatuses = "select id, user_id, num, status, upload_at, updated_at from orders where status in ($1, $2, $3, $4, $5) limit $6"
