package db

import "fmt"

type DynamicQuery struct {
	table      string
	operation  string
	where      string
	argCount   int
	ret        string
	insertList []string
	argsList   []interface{}
}

func Update(table string) DynamicQuery {
	d := DynamicQuery{}
	d.table = table
	d.operation = "UPDATE"
	d.insertList = []string{}
	d.argsList = []interface{}{}
	d.argCount = 0
	return d
}

func Insert(table string) DynamicQuery {
	d := DynamicQuery{}
	d.table = table
	d.operation = "INSERT"
	d.insertList = []string{}
	d.argsList = []interface{}{}
	d.argCount = 0
	return d
}

func (d DynamicQuery) Set(k string, v interface{}) DynamicQuery {
	d.insertList = append(d.insertList, k)
	d.argsList = append(d.argsList, v)
	d.argCount++
	return d
}

func (d DynamicQuery) Where(cond, op string, val interface{}) DynamicQuery {
	d.where += fmt.Sprintf("%v %v $%v", cond, op, d.argCount+1)
	d.argCount++
	d.argsList = append(d.argsList, val)
	return d
}

func (d DynamicQuery) And(cond, op string, val interface{}) DynamicQuery {
	d.where += fmt.Sprintf(" AND %v %v $%v", cond, op, d.argCount+1)
	d.argCount++
	d.argsList = append(d.argsList, val)
	return d
}

func (d DynamicQuery) Returning(col string) DynamicQuery {
	d.ret = fmt.Sprintf(" RETURNING %v", col)
	return d
}

func (d DynamicQuery) Query() (string, []interface{}) {
	if d.operation == "UPDATE" {
		q := fmt.Sprintf("UPDATE %v SET ", d.table)
		for i, v := range d.insertList {
			if i == 0 {
				q += fmt.Sprintf("%v = $%v", v, i+1)
			} else {
				q += fmt.Sprintf(", %v = $%v", v, i+1)
			}
		}
		if d.where != "" {
			q += fmt.Sprintf(" WHERE %v", d.where)
		}
		q += ";"
		return q, d.argsList
	}
	if d.operation == "INSERT" {
		valuesList := ""
		params := ""
		for i, v := range d.insertList {
			if i > 0 {
				valuesList += ", "
				params += ", "
			}
			valuesList += v
			params += fmt.Sprintf("$%v", i+1)
		}
		return fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v)%v;", d.table, valuesList, params, d.ret), d.argsList
	}
	return "", nil
}
