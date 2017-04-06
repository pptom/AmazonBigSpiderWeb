package smartuk

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/hunterhug/go_tool/util"
	"strconv"
	"strings"
)

type UrlController struct {
	baseController
}

func (this *UrlController) Index() {
	DB := orm.NewOrm()
	err := DB.Using("ukbasicdb")
	if err != nil {
		beego.Error("ukbasicdb err:" + err.Error())
		this.Rsp(false, err.Error())
	}
	DB.Using("smartdb")
	var categorys []orm.Params
	DB.Raw("SELECT bigpname as Bigpname,id FROM smart_category where pid=0 group by bigpname").Values(&categorys)
	this.Data["category"] = &categorys
	this.Layout = this.GetTemplate() + "/base/layout.html"
	this.TplName = this.GetTemplate() + "/url/uklist.html"

}

func (this *UrlController) Query() {
	DB := orm.NewOrm()
	err := DB.Using("ukbasicdb")
	if err != nil {
		beego.Error("ukbasicdb err:" + err.Error())
		this.Rsp(false, err.Error())
	}
	name := this.GetString("name")
	num := 0
	var maps []orm.Params
	page, _ := this.GetInt("page", 1)
	rows, _ := this.GetInt("rows", 30)
	start := (page - 1) * rows
	if name == "" {

		isvalid, _ := this.GetInt("isvalid", 2)
		bigname := this.GetString("bigname")
		small := this.GetString("small")

		where := []string{}
		wheresql := ""
		if bigname == "" {
		} else {
			where = append(where, `bigpid="`+bigname+`"`)
		}
		if small == "0" || small == "1" {
			where = append(where, "ismall="+small)
		}
		if isvalid == 1 || isvalid == 0 {
			where = append(where, `isvalid=`+util.IS(isvalid))
		}
		if len(where) == 0 {

		} else {
			wheresql = strings.Join(where, " and ")
			wheresql = "where " + wheresql
		}
		dudu := "SELECT * FROM smart_category " + wheresql + " order by createtime limit " + strconv.Itoa(start) + "," + strconv.Itoa(rows) + ";"
		//fmt.Println(dudu)
		DB.Raw(dudu).Values(&maps)

		dudu1 := "SELECT count(*) as num FROM smart_category " + wheresql + ";"

		DB.Raw(dudu1).QueryRow(&num)
	} else {
		dudu := "SELECT * FROM smart_category where name=? limit " + strconv.Itoa(start) + "," + strconv.Itoa(rows) + ";"
		DB.Raw(dudu, name).Values(&maps)
		dudu1 := "SELECT count(*) as num FROM smart_category where name=?;"
		DB.Raw(dudu1, name).QueryRow(&num)
	}
	if len(maps) == 0 {
		this.Data["json"] = &map[string]interface{}{"total": num, "rows": []interface{}{}}
	} else {
		this.Data["json"] = &map[string]interface{}{"total": num, "rows": &maps}
	}
	this.ServeJSON()
}

func (this *UrlController) Update() {
	DB := orm.NewOrm()
	err := DB.Using("ukbasicdb")
	if err != nil {
		beego.Error("ukbasicdb err:" + err.Error())
		this.Rsp(false, err.Error())
	}
	var maps []orm.Params
	isvalid := this.GetString("isvalid")
	id := this.GetString("id")
	page := this.GetString("page")
	if page == "" {
		dudu := "update smart_category set isvalid=? where id=?;"
		_, err := DB.Raw(dudu, isvalid, id).Values(&maps)
		if err == nil {
			this.Rsp(true, "good job")
		} else {
			this.Rsp(false, err.Error())
		}
	} else {
		dudu := "update smart_category set page=? where id=?;"
		_, err := DB.Raw(dudu, page, id).Values(&maps)
		if err == nil {
			this.Rsp(true, "good job")
		} else {
			this.Rsp(false, err.Error())
		}
	}
}
