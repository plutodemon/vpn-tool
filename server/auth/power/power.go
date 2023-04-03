package power

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	xormadapter "github.com/casbin/xorm-adapter/v2"
	_ "github.com/mattn/go-sqlite3"
	"ivs-net-server/auth/configure"
)

var Enforcer *casbin.Enforcer

func init() {
	connectPowerDB()
	//Enforcer.AddPolicy("alice", "data1", "read")
}

func connectPowerDB() {
	a, err := xormadapter.NewAdapter("sqlite3", configure.Config.Get("db.sqlite.path").(string), true)
	if err != nil {
		fmt.Println("test:", err)
	}
	text := `
				[request_definition]
				r = sub, obj, act

				[policy_definition]
				p = sub, obj, act

				[policy_effect]
				e = some(where (p.eft == allow))

				[matchers]
				m = r.sub == p.sub && r.obj == p.obj && r.act == p.act || r.sub == 1
	`
	m, _ := model.NewModelFromString(text)
	Enforcer, _ = casbin.NewEnforcer(m, a)
	Enforcer.LoadPolicy()
}
