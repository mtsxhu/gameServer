package dao

type User struct {
	Id int
	Name string
	Password string
}

func DBGetUser(name string) (string,error) {
	rows,err:=MysqlDB.Query(`select password from user where name = `,name)
	if err != nil {
		return "",err
	}
	var ret string
	rows.Scan(&ret)
	return ret,nil
}