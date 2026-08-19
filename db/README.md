# 数据驱动
## 说明
    简单实现的数据库驱动，后续再升级语句安全方面，目前暂时先在项目层面校验
## 使用方法
```golang
db := huomao_db.GetDbClass("sqlite3")
db.Connect("/User/huomao/db/sqlite3.db")
```
参考 [数据库示例](../example/db "数据库示例")


## GetDbClass支持参数
- [x] mysql
- [x] sqlite3

## 支持方法

| 方法 | 返回值 | 说明 |  
| :---: | :--- | :--- |  
| Connect(string) |  | 连接数据库，传入dsn |  
| ExecSql(string)  | error | 执行语句 |  
| Insert(string) |int, error  | 执行插入语句，返回LastInsertId |  
| Delete(string) | int, error | 删除语句，返回删除条数 |  
| Update(string) |int, error  | 更新语句，返回更新条数 |  
| Select(string) |error, *sql.Rows  | 查询语句 |  

