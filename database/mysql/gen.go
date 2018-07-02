package mysql

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"
)

type Table struct {
	name   string
	Fields []string
	Types  []string
	Flags  []string
}

//生成表struct
func (t *Table) String() string {
	buf := strings.Builder{}
	buf.WriteString("type ")
	buf.WriteString(toUpper(t.name) + " ")
	buf.WriteString("struct { \n")
	for i := range t.Fields {
		buf.WriteString(toUpper(t.Fields[i]))
		buf.WriteByte(' ')
		buf.WriteString(t.Types[i])
		buf.WriteByte(' ')
		buf.WriteString(t.Flags[i])
		buf.WriteByte('\n')
	}
	buf.WriteString("}")
	return buf.String()
}

func (t *Table) InsertAllStr() string {
	buf := strings.Builder{}
	buf.WriteString("INSERT INTO ")
	buf.WriteString(t.name + " (`")
	buf.WriteString(strings.Join(t.Fields, "`, `"))
	buf.WriteString("`) VALUES (:")
	buf.WriteString(strings.Join(t.Fields, ", :"))
	buf.WriteString(");")
	return buf.String()
}

func (t *Table) SelectAllByWhereStr() string {
	buf := strings.Builder{}
	buf.WriteString("SELECT `")
	buf.WriteString(strings.Join(t.Fields, "`, `"))
	buf.WriteString("` FROM " + t.name)
	buf.WriteString(" WHERE ")
	return buf.String()

}
func (t *Table) UpdateAllStr() string {
	buf := strings.Builder{}
	buf.WriteString("UPDATE " + t.name + " SET ")
	for i, v := range t.Fields {
		buf.WriteString(" `" + v + "`=:" + v)
		if i != len(t.Fields)-1 {
			buf.WriteString(", ")
		}
	}
	buf.WriteString(" WHERE ")
	return buf.String()
}

func toUpper(s string) string {
	rst := strings.Split(s, "_")
	buf := strings.Builder{}
	for _, v := range rst {
		if 'a' <= v[0] && v[0] <= 'z' {
			buf.WriteByte(v[0] - 'a' + 'A')
			buf.WriteString(v[1:])
		} else {
			buf.WriteString(v)
		}
	}
	return buf.String()
}

type Config struct {
	TableName string
	fieldsMap map[string]string
}

func (c *Config) AddFieldsMap(k, v string) *Config {
	c.fieldsMap[k] = c.fieldsMap[v]
	return c
}
func (c *Config) MergeFieldsMap(m map[string]string) {
	for k, v := range m {
		c.fieldsMap[k] = v
	}
}
func (c *Config) FieldsContains(v string) bool {
	if _, ok := c.fieldsMap[v]; ok {
		return true
	}
	return false
}
func (c *Config) Fields(v string) string {
	if k, ok := c.fieldsMap[v]; ok {
		return k
	}
	return "interface{}"
}

func NewConfig(table string) *Config {
	return &Config{TableName: table}
}

//从db查询生成cfg所配置的表
func ParseTable(db *sql.DB, cfg *Config) (*Table, error) {
	tb := new(Table)
	row, err := db.Query("desc " + cfg.TableName)
	if err != nil {
		return nil, err
	}
	tb.name = cfg.TableName
	for row.Next() {
		var Field, Type string
		var tmp interface{}
		if err := row.Scan(&Field, &Type, &tmp, &tmp, &tmp, &tmp); err != nil {
			return nil, err
		}
		f := Field
		tb.Fields = append(tb.Fields, f)
		var t string
		if cfg.FieldsContains(Type) {
			t = cfg.Fields(Type)
		} else if strings.Contains(Type, "bigint") {
			t = "int64"
		} else if strings.Contains(Type, "int") {
			t = "int32"
		} else if strings.Contains(Type, "char") || strings.Contains(Type, "datetime") {
			t = "string"
		} else {
			t = "interface{}"
		}
		tb.Types = append(tb.Types, t)
		g := " `json:\"" + Field + "\"`"
		tb.Flags = append(tb.Flags, g)
	}
	return tb, nil
}

type ConnConfig struct {
	MysqlConfig
	Config
	SshHost string
	SshUser string
}

//从一个配置文件读取一个生成数据库表字段的配置
func ParseConfig() *ConnConfig {
	path := os.Getenv("PWD") + "/mysql_config.properties"
	fd, err := os.Open(path)
	if err != nil {
		fmt.Printf("cannot open file %s ! err: %v\n", path, err)
		return nil
	}
	scan := bufio.NewScanner(fd)
	cfg := &ConnConfig{}
	for scan.Scan() {
		line := scan.Text()
		if strings.Trim(line, " ")[0] == '#' { //ignore
			continue
		}
		kv := strings.Split(line, "=")
		if len(kv) != 2 {
			fmt.Printf("read err int %s,ignore....", line)
			continue
		}
		k := strings.Trim(kv[0], " ")
		v := strings.Trim(kv[1], " ")
		switch k {
		case "host":
			cfg.Host = v
		case "username":
			cfg.User = v
		case "password":
			cfg.PassWord = v
		case "db":
			cfg.Db = v
		case "table":
			cfg.TableName = v
		case "ssh_user":
			cfg.SshUser = v
		case "ssh_host":
			cfg.SshHost = v
		}
	}
	return cfg
}

func ReadTableFromDefault() string {
	cfg := ParseConfig()
	var db *sql.DB
	var err error
	if cfg.SshHost != "" {
		db, err = NewMysqlDbInSSH(cfg.SshHost, cfg.SshUser, &cfg.MysqlConfig)
	} else {
		db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", cfg.User, cfg.PassWord, cfg.Host, cfg.Db))
	}
	if err != nil {
		fmt.Println("op driver err ", err.Error())
		return ""
	}
	table, err := ParseTable(db, &cfg.Config)
	if err != nil {
		fmt.Println("parse table err ", err.Error())
		return ""
	}
	fmt.Println(table.String())
	return table.String()
}
