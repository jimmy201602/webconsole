package utils

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"strings"
	"time"
)

var (
	ABC_Conf, conf_err = Get_Conf()
)

func init() {
	if nil != conf_err {
		Log_Fatal(conf_err.Error())
	}
}

type (
	MySQL struct {
		Url string
		DB  *sql.DB
	}

	SQLUtils struct {
		t           reflect.Type
		v           reflect.Value
		primary     string
		tableName   string
		tableAsName string
		orderBy     string
		where       string
		whereConcat string
		limit       string
		fieldNames  []string
		fieldIndexs map[string]int
		fidldValues map[string]interface{}
	}
)

func NewSQL() *SQLUtils {
	return &SQLUtils{}
}

func (s *SQLUtils) RegisterStruct(i interface{}) *SQLUtils {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	s.t = t
	s.v = v

	s.tableName = strings.ToLower(t.Name())
	s.tableAsName = strings.ToLower(s.tableName[:1])
	s.fieldNameHandle()
	s.fieldValueHandle()

	s.primary = "id"

	return s
}

func (s *SQLUtils) SetPK(pk string) *SQLUtils {
	pk = strings.TrimSpace(pk)
	s.primary = pk
	return s
}

func (s *SQLUtils) AsName(n string) *SQLUtils {
	s.tableAsName = strings.TrimSpace(n)
	return s
}

func (s *SQLUtils) QueryAll() string {
	tempFieldNames := make([]string, 0, 0)

	for _, name := range s.fieldNames {
		tempFieldNames = append(tempFieldNames, s.tableAsName+"."+name)
	}

	sqlStr := "select " + strings.Join(tempFieldNames, ",") + " from " + s.tableName + " " + s.tableAsName

	if s.where != "" {
		sqlStr = sqlStr + " " + strings.TrimSpace(s.where)
	}

	if s.orderBy != "" {
		sqlStr = sqlStr + " " + strings.TrimSpace(s.orderBy)
	}

	if s.limit != "" {
		sqlStr = sqlStr + " " + strings.TrimSpace(s.limit)
	}

	sqlStr = sqlStr + ";"

	return sqlStr
}

func (s *SQLUtils) QueryCount() string {
	sqlStr := "select count(" + s.tableAsName + "." + s.primary + ") from " + s.tableName + " " + s.tableAsName

	if s.where != "" {
		sqlStr = sqlStr + " " + strings.TrimSpace(s.where)
	}

	sqlStr = sqlStr + ";"

	return sqlStr
}

func (s *SQLUtils) OrderBy(name string, sort string) *SQLUtils {
	name = strings.TrimSpace(name)
	sort = strings.TrimSpace(sort)

	ng, _ := StringUtils(name).RegexpSQLVal()
	sg, _ := StringUtils(sort).RegexpSQLVal()

	if ng && sg {
		n := strings.TrimSpace(s.tableAsName + "." + name)
		sort = strings.TrimSpace(sort)
		if strings.EqualFold(sort, "asc") || strings.EqualFold(sort, "desc") {
			s.orderBy = "order by " + n + " " + sort
		} else {
			s.orderBy = "order by " + n + " desc"
		}
	} else {
		s.orderBy = "order by " + s.tableAsName + "." + s.primary + " asc"
	}
	return s
}

func (s *SQLUtils) WhereAnd(name string, condition string, value string) *SQLUtils {
	name = strings.TrimSpace(name)
	condition = strings.TrimSpace(condition)
	value = strings.TrimSpace(value)

	ng, _ := StringUtils(name).RegexpSQLVal()
	vg, _ := StringUtils(value).RegexpSQLVal()
	cg, _ := StringUtils(condition).RegexpSQLSgin()

	if ng && vg && cg {
		name = strings.TrimSpace(s.tableAsName + "." + name)
		if s.where == "" {
			s.where = "where " + name + condition + value
		} else {
			s.where = s.where + " and " + name + condition + value
		}
	} else {
		if s.where == "" {
			s.where = "where 1!=1"
		} else {
			s.where = s.where + " and 1!=1"
		}
	}
	return s
}

func (s *SQLUtils) WhereConcat(value string, field ...string) *SQLUtils {
	value = strings.TrimSpace(value)
	value = strings.Trim(value, "'")
	vg, _ := StringUtils(value).RegexpSQLVal()
	if vg {
		sqlStr := "where concat("
		for _, v := range field {
			v = strings.TrimSpace(v)
			if b, _ := StringUtils(v).RegexpSQLVal(); b {
				sqlStr = sqlStr + s.tableAsName + "." + v + ","
			}
		}
		sqlStr = strings.TrimRight(sqlStr, ",") + ") like '%" + value + "%'"
		s.whereConcat = sqlStr
	}
	return s
}

func (s *SQLUtils) QueryCountByConcat() string {
	sqlStr := "select count(" + s.tableAsName + "." + s.primary + ") from " + s.tableName + " " + s.tableAsName

	if s.whereConcat != "" {
		sqlStr = sqlStr + " " + strings.TrimSpace(s.whereConcat)
	}

	sqlStr = sqlStr + ";"

	return sqlStr
}

func (s *SQLUtils) QueryByConcat() string {
	tempFieldNames := make([]string, 0, 0)

	for _, name := range s.fieldNames {
		tempFieldNames = append(tempFieldNames, s.tableAsName+"."+name)
	}

	sqlStr := "select " + strings.Join(tempFieldNames, ",") + " from " + s.tableName + " " + s.tableAsName

	if s.whereConcat != "" {
		sqlStr = sqlStr + " " + strings.TrimSpace(s.whereConcat)
	}

	if s.orderBy != "" {
		sqlStr = sqlStr + " " + strings.TrimSpace(s.orderBy)
	}

	if s.limit != "" {
		sqlStr = sqlStr + " " + strings.TrimSpace(s.limit)
	}

	sqlStr = sqlStr + ";"

	return sqlStr
}

func (s *SQLUtils) WhereOr(name string, condition string, value string) *SQLUtils {
	name = strings.TrimSpace(name)
	condition = strings.TrimSpace(condition)
	value = strings.TrimSpace(value)

	ng, _ := StringUtils(name).RegexpSQLVal()
	vg, _ := StringUtils(value).RegexpSQLVal()
	cg, _ := StringUtils(condition).RegexpSQLSgin()

	if ng && vg && cg {
		name = strings.TrimSpace(s.tableAsName + "." + name)
		if s.where == "" {
			s.where = "where " + name + condition + value
		} else {
			s.where = s.where + " or " + name + condition + value
		}
	} else {
		if s.where == "" {
			s.where = "where 1!=1"
		} else {
			s.where = s.where + " or 1!=1"
		}
	}
	return s
}

func (s *SQLUtils) Limit(index int, num int) *SQLUtils {
	s.limit = "limit " + ToStr(index) + "," + ToStr(num)
	return s
}

func (s *SQLUtils) Insert() (string, []interface{}) {

	fieldKey := make([]string, 0, 0)
	fieldVal := make([]interface{}, 0, 0)
	questionMark := make([]string, 0, 0)

	for k, v := range s.fidldValues {
		fieldKey = append(fieldKey, k)
		fieldVal = append(fieldVal, v)
		questionMark = append(questionMark, "?")
	}

	sqlStr := "insert into " + s.tableName + " (" + strings.Join(fieldKey, ",") + ") values (" + strings.Join(questionMark, ",") + ");"

	return sqlStr, fieldVal
}

func (s *SQLUtils) Update() (string, []interface{}) {

	sqlStr := "update " + s.tableName + " " + s.tableAsName + " set "

	fieldVal := make([]interface{}, 0, 0)

	for k, v := range s.fidldValues {
		if s.primary != k {
			sqlStr = sqlStr + s.tableAsName + "." + k + "=" + "?,"
			fieldVal = append(fieldVal, v)
		}
	}

	sqlStr = strings.TrimRight(sqlStr, ",")

	if s.where != "" {
		sqlStr = sqlStr + " " + strings.TrimSpace(s.where)
	}
	sqlStr = sqlStr + ";"
	return sqlStr, fieldVal
}

func (s *SQLUtils) Delete() string {
	sqlStr := "delete " + s.tableAsName + " from " + s.tableName + " as " + s.tableAsName
	if s.where != "" {
		sqlStr = sqlStr + " " + strings.TrimSpace(s.where)
	}
	sqlStr = sqlStr + ";"
	return sqlStr
}

func (s *SQLUtils) fieldNameHandle() {

	fieldNames := make([]string, 0, 0)
	fieldIndexs := make(map[string]int)

	for i := 0; i < s.t.NumField(); i++ {
		field := s.t.Field(i)
		fieldTag := field.Tag.Get("field")
		if "" != fieldTag {
			fieldNames = append(fieldNames, fieldTag)
			fieldIndexs[fieldTag] = i
		} else {
			fieldNames = append(fieldNames, field.Name)
			fieldIndexs[field.Name] = i
		}
	}

	s.fieldNames = fieldNames
	s.fieldIndexs = fieldIndexs
}

func (s *SQLUtils) fieldValueHandle() {
	fieldValues := make(map[string]interface{})
	for k, i := range s.fieldIndexs {
		fv := s.v.Field(i)
		v := s.convertValue(fv)
		if nil != v {
			fieldValues[k] = s.convertValue(fv)
		}
	}
	s.fidldValues = fieldValues
}

func (s *SQLUtils) convertValue(fv reflect.Value) interface{} {
	switch fv.Kind() {
	case reflect.Ptr, reflect.Interface:
		if fv.IsNil() {
			return nil
		}
		return s.convertValue(fv.Elem())
	case reflect.Bool:
		return fv.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fv.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return fv.Uint()
	case reflect.Float32, reflect.Float64:
		return fv.Float()
	case reflect.String:
		return fv.String()
	}
	return nil
}

func (d *MySQL) DB_Open() (*MySQL, error) {
	d.Url = ABC_Conf.DB.Url
	if nil == d.DB {
		db, err := sql.Open("mysql", d.Url)
		if nil != err {
			return nil, err
		}
		db.SetMaxOpenConns(ABC_Conf.DB.MaxOpenConns)
		db.SetMaxIdleConns(ABC_Conf.DB.MaxIdleConns)
		err = db.Ping()
		if nil != err {
			return nil, err
		}
		d.DB = db
	}
	return d, nil
}

func (d *MySQL) Query(obj interface{}, s string) ([]interface{}, error) {
	rows, err := d.DB.Query(s)
	if nil != err {
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if nil != err {
		return nil, err
	}
	scans := make([]interface{}, len(columns))
	values := make([]sql.RawBytes, len(columns))
	for i := range values {
		scans[i] = &values[i]
	}
	data := make([]interface{}, 0, 0)
	for rows.Next() {
		if err := rows.Scan(scans...); nil != err {
			return nil, err
		}
		t := reflect.TypeOf(obj)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		val := reflect.New(t)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		fieldMap := make(map[string]int)
		for i := 0; i < t.NumField(); i++ {
			tf := t.Field(i)
			tag := tf.Tag.Get("field")
			if tag == "" {
				fieldMap[tf.Name] = i
			} else {
				fieldMap[tag] = i
			}
		}
		for i, v := range values {
			fieldName := columns[i]
			fieldIndex := fieldMap[fieldName]
			field := val.Field(fieldIndex)
			if field.IsValid() {
				if field.CanSet() {
					fieldType := field.Kind()
					switch fieldType {
					case reflect.String:
						field.SetString(StringUtils(v).String())
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						vv, _ := StringUtils(v).Int64()
						field.SetInt(vv)
					case reflect.Bool:
						vv, _ := StringUtils(v).Bool()
						field.SetBool(vv)
					case reflect.Float32, reflect.Float64:
						vv, _ := StringUtils(v).Float64()
						field.SetFloat(vv)
					case reflect.Struct:
						timeTag := t.Field(i).Tag.Get("time_format")
						if timeTag != "" {
							t, err := time.Parse(timeTag, StringUtils(v).String())
							if nil != err {
								Log_Err(err.Error())
								continue
							} else {
								field.Set(reflect.ValueOf(t))
							}
						}
					}
				}
			}
		}
		data = append(data, val.Interface())
	}
	return data, nil
}

func (d *MySQL) QueryOne(obj interface{}, s string) (interface{}, error) {
	data, err := d.Query(obj, s)
	if nil != err {
		return nil, err
	} else {
		if len(data) <= 0 {
			return nil, nil
		} else {
			return data[0], nil
		}
	}
}

func (d *MySQL) QueryCount(s string) (int, error) {
	rows, err := d.DB.Query(s)
	if nil != err {
		return 0, err
	}
	defer rows.Close()
	var rowCount int
	for rows.Next() {
		err := rows.Scan(&rowCount)
		if nil != err {
			return 0, err
		}
	}
	return rowCount, nil
}

func (d *MySQL) Delete(s string) (n int64, err error) {
	res, err := d.DB.Exec(s)
	if nil != err {
		return 0, err
	}
	n, err = res.RowsAffected()
	if nil != err {
		return 0, err
	} else {
		return
	}
}

func (d *MySQL) Exec(s string, args ...interface{}) (n int64, err error) {
	stmt, err := d.DB.Prepare(s)
	if nil != err {
		return 0, err
	}
	res, err := stmt.Exec(args...)
	if nil != err {
		return 0, err
	}
	n, err = res.RowsAffected()
	if nil != err {
		return 0, err
	} else {
		return
	}
}
