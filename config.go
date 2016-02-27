package mgox

import (
	"bufio"
	"github.com/kimiazhu/log4go"
	"io"
	"os"
	"strings"
	"fmt"
)

type dbconfig struct {
	Host     string
	Database string
	Username string
	Password string
}

var DBConfig dbconfig

type PropertyReader struct {
	m map[string]string
}

func (p *PropertyReader) init(path string) {

	p.m = make(map[string]string)

	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		s := strings.TrimSpace(string(b))

		//log.Println(s)

		if strings.Index(s, "#") == 0 {
			continue
		}

		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}

		frist := strings.TrimSpace(s[:index])
		if len(frist) == 0 {
			continue
		}
		second := strings.TrimSpace(s[index+1:])

		pos := strings.Index(second, "\t#")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " #")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, "\t//")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " //")
		if pos > -1 {
			second = second[0:pos]
		}

		if len(second) == 0 {
			continue
		}

		p.m[frist] = strings.TrimSpace(second)
	}
}

func (c PropertyReader) Read(key string) string {
	v, found := c.m[key]
	if !found {
		return ""
	}
	return v
}

// Config 用于直接传递数据库参数来配置数据库
func Config(host, dbname, username, userpass string) {
	DBConfig.Host = host
	DBConfig.Database = dbname
	DBConfig.Username = username
	DBConfig.Password = userpass
	if host != "" {
		log4go.Debug(fmt.Sprintf("host=%s,database=%s,username=%s,password=***\n", DBConfig.Host, DBConfig.Database, DBConfig.Username))
	}
}

// LoadConfig 支持从配置文件读取mongodb配置
func LoadConfig(path string) {
	db := new(PropertyReader)
	db.init(path)
	Config(db.m["host"], db.m["database"], db.m["username"], db.m["password"])
}


// init 方法会自动在当前目录,以及当前目录下的conf目录中查找mgox.properties文件并读取其中的信息
func init() {
	if _, err := os.Stat("mgox.properties"); err == nil {
		LoadConfig("mgox.properties")
	} else if _, err = os.Stat("conf/mgox.properties"); err == nil {
		LoadConfig("conf/mgox.properties")
	} else {
		log4go.Info("config file cannot be discovered, you need to load it by yourself")
	}
}
