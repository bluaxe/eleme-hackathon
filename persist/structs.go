package persist

type sqls struct {
	GET_USER_TEST    string
	TEST             string
	GET_USER_BY_NAME string
	SELECT_ALL_FOODS string
}

func initSQL() sqls {
	var s sqls
	s.GET_USER_TEST = `select name, password from user limit 10`
	s.TEST = `SHOW TABLES`
	s.GET_USER_BY_NAME = `select id, name, password from user where name=? limit 1`
	s.SELECT_ALL_FOODS = `select id, stock, price from food`
	return s
}
