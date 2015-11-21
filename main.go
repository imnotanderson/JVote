package main

import (
	"fmt"
	"net/http"
	//	"strconv"
	"html/template"
	//	"io/ioutil"
	//	"os"
	"strconv"
)

type MyData struct {
	Title       string
	OptNameList []string
	OptCount    []int
	Detail      [][]string
	DetailStr   template.HTML
}

var data MyData
var count int = 0

func main() {
	data = MyData{
		Title:       "",
		OptNameList: []string{},
		OptCount:    []int{},
		Detail:      nil,
	}
	http.HandleFunc("/", JHandle)
	err := http.ListenAndServe(":4000", nil)
	fmt.Print(err.Error())
}

func JHandle(w http.ResponseWriter, r *http.Request) {
	tplt, err := template.ParseFiles("index.htm")
	checkErr(err)
	r.ParseForm()
	if data.Title == "" || len(data.OptNameList) <= 0 {
		data.Title = r.Form.Get("title")

		optNameList := data.OptNameList
		if optNameList == nil {
			optNameList = []string{}
		}
		i := 1
		for {
			getname := "opt" + strconv.Itoa(i)
			name := r.Form.Get(getname)
			if name == "" {
				break
			}
			i++
			optNameList = append(optNameList, name)
		}
		data.OptNameList = optNameList
		data.OptCount = make([]int, len(data.OptNameList))
		data.Detail = make([][]string, len(data.OptNameList))

	} else {
		voter := r.Form.Get("voter")
		if voter == "" {
			tplt.Execute(w,data)
			return
		}
		optname := r.Form.Get("optname")
		idx ,canVote := checkOpt(optname,voter)
		if idx >= 0 && canVote{
			recordVoter(idx, voter,r)
		}

	}
	DetailUpt()
	err = tplt.Execute(w, data)
}

func DetailUpt() {
	str := "Result:"
	str+="<p>"
	for idx,voterList:= range data.Detail{
		str+=data.OptNameList[idx]+":"+strconv.Itoa(len(voterList))
	}
	str+="</p>"
	for idx, nameList := range data.Detail {
		str += "<li>"
		str += data.OptNameList[idx]
		str+="<p>"
		for _, name := range nameList {
			str += name + ","
		}
		str+="</p>"
		str += "</li>"
	}
	data.DetailStr = template.HTML(str)
}

func recordVoter(optIdx int, voter string,r *http.Request) {
	if optIdx >= len(data.Detail) {
		return
	}
	voterList := data.Detail[optIdx]
	if voterList == nil {
		voterList = []string{}
	}
	for _, tmVoter := range voterList {
		if voter == tmVoter {
			return
		}
	}
	voterList = append(voterList, voter)
	data.Detail[optIdx] = voterList
	fmt.Println(voter+ "["+r.RemoteAddr +"]:"+data.OptNameList[optIdx])
}

func checkOpt(optname ,voter string) (int ,bool){
	idx:=-1
	for i, opt := range data.OptNameList {
		if opt == optname {
			idx = i
		}
		for _,name :=range data.Detail[i]  {
			if name==voter{
				return -1,false
			}
		}
	}
	return idx,true
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}