//Package model Copyright 2020 snailouyang.  All rights reserved.
//数据库操作,使用第三方库https://github.com/go-gorm/gorm封装
package model

//Create 新增记录
//@data 查询结果存放在变量
//@name mysql对应key,默认为default
func Create(data interface{}, name ...string) error {
	return Db(name...).Model(data).Create(data).Error
}

//Find 查询多条记录
//@where 查询条件 field=>值
//@data 查询结果存放在变量
//@name mysql对应key,默认为default
func Find(where map[string]interface{}, data interface{}, name ...string) error {
	tx := Db(name...)
	for query, arg := range where {
		tx = tx.Where(query, arg)
	}
	return tx.Find(data).Error
}

//First 查询单条记录,记录查找不到会返回 gorm.ErrRecordNotFound 错误
//@where 查询条件 field=>值
//@data 查询结果存放在变量
//@name mysql对应key,默认为default
func First(where map[string]interface{}, data interface{}, name ...string) error {
	tx := Db(name...)
	for query, arg := range where {
		tx = tx.Where(query, arg)
	}
	return tx.First(data).Error
}

//Updates 更新多条记录
//@where 查询条件 field=>值
//@data 更新field
//@name mysql对应key,默认为default
func Updates(model interface{}, where map[string]interface{}, data map[string]interface{}, name ...string) error {
	tx := Db(name...)
	for query, arg := range where {
		tx = tx.Where(query, arg)
	}
	return tx.Model(model).Updates(data).Error
}

//Update 更新单条记录,提高mysql性能
//@where 查询条件 field=>值
//@data 更新field
//@name mysql对应key,默认为default
func Update(model interface{}, where map[string]interface{}, data map[string]interface{}, name ...string) error {
	tx := Db(name...)
	for query, arg := range where {
		tx = tx.Where(query, arg)
	}
	return tx.Model(model).Limit(1).Updates(data).Error
}

//Delete 指定条件删除记录,删除对象需要指定主键，否则会触发 批量 Delete(匹配到都会删除)
//@data 模型对象
//@name mysql对应key,默认为default
func Delete(data interface{}, where map[string]interface{}, name ...string) error {
	tx := Db(name...)
	for query, arg := range where {
		tx = tx.Where(query, arg)
	}
	return tx.Delete(data).Error
}

//RawExec 执行原生sql,只支持update delete insert drop等写入语句
//@name mysql对应key,默认为default
//@sql 要执行的sql
//@value 要执行的sql参数
func RawExec(name string, sql string, value ...interface{}) error {
	if len(name) <= 0 {
		name = "default"
	}
	return Db(name).Exec(sql, value...).Error
}

//RawQuery 执行原生sql,只支持查询语句，即有结果集返回的操作
//@name mysql对应key,默认为default
//@result 要执行的sql
//@sql 要执行的sql
//@value 要执行的sql参数
func RawQuery(name string, result interface{}, sql string, value ...interface{}) error {
	if len(name) <= 0 {
		name = "default"
	}
	return Db(name).Raw(sql, value...).Find(&result).Error
}
