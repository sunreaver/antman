package db

import (
	"fmt"
	"strings"
	"time"
)

// Config Config.
type Config struct {
	Type         DBType   `toml:"type"`
	MasterURI    string   `toml:"master_uri"`
	SlaveURIs    []string `toml:"slave_uris"`
	LogMode      bool     `toml:"log_mode"`
	MaxIdleConns int      `toml:"max_idle_conns"`
	MaxOpenConns int      `toml:"max_open_conns"`
}

// 为现在管理的数据库类型创建枚举值
// 不使用string类型常量，而使用一套单独的枚举类型原因是：
//  1. 避免不相干枚举定义值重叠，例如有个功能要管理“行外客户当前使用库枚举类型”（有值：mysql、oracle。。。），这些值与现有定义重叠
//  2. 让语义比较确定，约定类型范围，避免使用范围外的类型
//  3. 容易为类型实现一些扩展方法
type DBType string

const (
	DBTypeUndefined DBType = "unknow:empty"
	DBTypeMYSQL     DBType = "mysql"
	DBTypeORACLE    DBType = "oracle"
	DBTypePGSQL     DBType = "postgres"
	DBTypeMogDb     DBType = "mogdb"
	DBTypeSQLite    DBType = "sqllite"
	DBTypeDameng    DBType = "dameng"
)

/*
*
根据dbtype字符串返回一个枚举值类型
*/
func GetDBTypeEnumFromStr(dbtypeStr string) DBType {
	dbtypeEnum := DBType(strings.ToLower(dbtypeStr))
	switch dbtypeEnum {
	case DBTypeMYSQL, DBTypeORACLE, DBTypePGSQL, DBTypeSQLite, DBTypeMogDb:
		return dbtypeEnum
	}
	return DBType(fmt.Sprintf("unknow:%v", dbtypeEnum))
}

const (
	DBDataTypeJson  string = "JSON"
	DBDataTypeJsonB string = "JSONB"
)

// 从数据库取出的字符串数据，如果其用于生成insert或update语句时，需要对其中的特殊字符进行转义，这样生成的插入语句才能正常使用
// mysql需要处理单引号 ' ==> ” 和 反斜杠 \ ==> \\
// postgre、oracle需要处理单引号 ' ==> ”
// forDML 表示是否生成为更新、插入格式
//
//	ture: 返回带单引号的 '2022-11-01' 用于生成insert语句的value内容
//	false：仅返回 2022-11-01 用于展示
type ValEscapeFn func(value string, forDML bool) string // 对值字符串进行转义

// filedType是该字段的类型 eg. time\timestamp with time zone\varchar...
// timeStr 是sqlagent查询获取到的数据，是time.RFC3339Nano格式
// forDML 表示是否生成为更新、插入格式
//
//	ture: 返回 timestamp '2022-11-01' 或者带单引号的'2022-11-01' 用于生成insert语句的value内容
//	false：表示不对转换出的内容做修改， 仅返回 2022-11-01 纯内容
type ValTimeFn func(fieldType, value string, forDML bool) string // 对时间进行格式化输出

// 根据数据库类型获取其字符串转义处理方法，这个转义方法只处理某个字符类型的字段取出的value
func (t DBType) GetValEscapeFn() ValEscapeFn {
	switch t {
	case DBTypeMYSQL:
		// abc'abc\n   ==>  'abc''abc\\n' 或 abc''abc\\n
		return func(value string, forDML bool) string {
			value = strings.ReplaceAll(value, `\`, `\\`)
			value = strings.ReplaceAll(value, `'`, `''`)
			if forDML {
				return fmt.Sprintf("'%s'", value)
			}
			return value
		}
	case DBTypeORACLE, DBTypePGSQL, DBTypeMogDb:
		// abc'abc\n   ==>  'abc''abc\n' 或 abc''abc\n
		return defaultValEscapeFn
	default:
		return defaultValEscapeFn
	}
}

func defaultValEscapeFn(value string, forDML bool) string {
	value = strings.ReplaceAll(value, `'`, `''`)
	if forDML {
		return fmt.Sprintf("'%s'", value)
	}
	return value
}

// 将获取到的time.RFC3339Nano格式的时间字段进行转换，转换成在各个数据库语句中可以使用的值内容
func (t DBType) GetValTimeFn() ValTimeFn {
	switch t {
	case DBTypeMYSQL:
		return func(fieldType, timeStr string, forDML bool) (resultStr string) {
			fieldType = strings.ToUpper(fieldType)
			resultStr = timeStr
			defer func() {
				if forDML {
					resultStr = fmt.Sprintf("'%s'", resultStr)
				}
			}()
			t, err := time.Parse(time.RFC3339Nano, timeStr)
			if err != nil { // 如果出现异常情况，直接单引号包起来原样返回
				return resultStr
			}
			// 2022-11-01T06:27:31.999999999Z07:00   ==>  2022-11-01 06:27:31.000000+08:00
			if strings.Contains(fieldType, "DATETIME") || strings.Contains(fieldType, "TIMESTAMP") {
				return t.Format("2006-01-02 15:04:05.000000")
			}
			// 2022-11-01T06:27:31.999999999Z07:00   ==>  2022-11-01
			if strings.Contains(fieldType, "DATE") {
				return t.Format("2006-01-02")
			}

			return resultStr
		}
	case DBTypePGSQL:
		return func(fieldType, timeStr string, forDML bool) (resultStr string) {
			fieldType = strings.ToUpper(fieldType)
			resultStr = timeStr
			defer func() {
				if forDML {
					resultStr = fmt.Sprintf("'%s'", resultStr)
				}
			}()
			t, err := time.Parse(time.RFC3339Nano, timeStr)
			if err != nil { // 如果出现异常情况，直接单引号包起来原样返回
				return resultStr
			}
			// 2022-11-01T06:27:31.999999999Z07:00   ==>  2022-11-01 06:27:31.000000 +08:00
			if strings.Contains(fieldType, "TIMESTAMP") {
				return t.Format("2006-01-02 15:04:05.000000 -07:00")
			}
			// 2022-11-01T06:27:31.999999999Z07:00   ==>  2022-11-01
			if strings.Contains(fieldType, "DATE") {
				return t.Format("2006-01-02")
			}
			return resultStr
		}
	case DBTypeMogDb:
		return func(fieldType, timeStr string, forDML bool) (resultStr string) {
			fieldType = strings.ToUpper(fieldType)
			resultStr = timeStr
			defer func() {
				if forDML {
					resultStr = fmt.Sprintf("'%s'", resultStr)
				}
			}()
			t, err := time.Parse(time.RFC3339Nano, timeStr)
			if err != nil { // 如果出现异常情况，直接单引号包起来原样返回
				return resultStr
			}
			// 2022-11-01T06:27:31.999999999Z07:00   ==>  2022-11-01 06:27:31.000000 +08:00
			if strings.Contains(fieldType, "TIMESTAMP") {
				return t.Format("2006-01-02 15:04:05.000000 -07:00")
			}
			// 2022-11-01T06:27:31.999999999Z07:00   ==>  2022-11-01
			if strings.Contains(fieldType, "DATE") {
				return t.Format("2006-01-02")
			}
			// 0000-01-01T06:27:31.999999999Z07:00   ==>  06:27:31  mogdb time类型sqlagent取出来是time对象，而pg time类型取出来是字符串
			if strings.Contains(fieldType, "TIME") {
				return t.Format("15:04:05 -07:00")
			}
			return resultStr
		}
	case DBTypeORACLE:
		// DATE\TIMESTAMP 类型： 2006-01-02T15:04:05.999999999Z07:00   ==> TIMESTAMP '2022-11-01 14:25:34.968000 +01:00'
		return func(fieldType, timeStr string, forDML bool) (resultStr string) {
			fieldType = strings.ToUpper(fieldType)
			resultStr = timeStr
			defer func() {
				if forDML {
					resultStr = fmt.Sprintf("TIMESTAMP '%s'", resultStr)
				}
			}()
			t, err := time.Parse(time.RFC3339Nano, timeStr)
			if err != nil { // 如果出现异常情况，直接单引号包起来原样返回
				return resultStr
			}
			return t.Format("2006-01-02 15:04:05.000000 -07:00")
		}
	default:
		return func(fieldType, value string, forDML bool) string { return value }
	}
}
