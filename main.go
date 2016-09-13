package main

import (
	"net/http"
	"io/ioutil"
	"log"
	"fmt"
	"encoding/json"
	"github.com/tealeg/xlsx"
	"encoding/gob"
	"bytes"
)

func GetBytes(key interface{}) ([]byte, error) {
    var buf bytes.Buffer
    enc := gob.NewEncoder(&buf)
    err := enc.Encode(key)
    if err != nil {
        return nil, err
    }
    return buf.Bytes(), nil
}

// Object of JSON

type First struct {
	Type string `json:"type"`
	Features []struct {
		Type string `json:"type"`
		Properties struct {
			GEOID string `json:"GEO_ID"`
			STATE string `json:"STATE"`
			NAME string `json:"NAME"`
			LSAD string `json:"LSAD"`
			CENSUSAREA float64 `json:"CENSUSAREA"`
		} `json:"properties"`
		Geometry interface{} `json:"geometry"`
	} `json:"features"`
}


// Accepting response in POST format 

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	requestbody, err1 := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err1 != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(requestbody))
		return
	}

		reqobj := First{}
		json.Unmarshal([]byte(requestbody), &reqobj)

// xlsx Conversion Code
    
    var file *xlsx.File
    var sheet *xlsx.Sheet
    var row *xlsx.Row
    var cell *xlsx.Cell
    var err error

    file = xlsx.NewFile()
    sheet, err = file.AddSheet("gz_2010_us_040_00_5m")
    if err != nil {
        fmt.Printf(err.Error())
    }

    for i:= range reqobj.Features{
      row = sheet.AddRow()
      cell = row.AddCell()
      cell.Value = reqobj.Features[i].Properties.GEOID
      cell = row.AddCell()
      cell.Value = reqobj.Features[i].Properties.LSAD
      cell = row.AddCell()
      cell.Value = reqobj.Features[i].Properties.NAME
      cell = row.AddCell()
      cell.Value = reqobj.Features[i].Properties.STATE
      cell = row.AddCell()
      fmt.Println(reqobj.Features[i].Geometry)
      str := fmt.Sprint(reqobj.Features[i].Geometry)
      fmt.Println(str)
      cell.Value = str

  }
    err = file.Save("MyXLSXFile1.xlsx")
    if err != nil {
        fmt.Printf(err.Error())
    }
}


//Send JSON AT localhost:8081

func main(){
  http.HandleFunc("/", Handler)
err := http.ListenAndServe(":8081", nil)
if err != nil {
  log.Fatal("ListenAndServe: ", err)
}
}
