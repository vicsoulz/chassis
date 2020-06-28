package excel

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize"
)

var cellName = []string{
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "S", "Y", "Z",
}

type Excel struct {
	data interface{}
	File *excelize.File
}

func (e *Excel) FromSlice(source interface{}) error {
	e.data = source
	b, err := json.Marshal(source)
	if err != nil {
		return err
	}
	var data []map[string]interface{}

	err = json.Unmarshal(b, &data)
	if err != nil {
		return err
	}

	return e.FromSliceMap(data)
}

func (e *Excel) FromSliceMap(source []map[string]interface{}) error {
	if len(source) == 0 {
		return errors.New("source data len is 0")
	}

	e.data = source
	e.File = excelize.NewFile()

	head := make(map[string]interface{})

	for i := 0; i < len(source); i++ {
		j := 0
		for k, v := range source[i] {

			// 加入不存在表头
			if _, ok := head[k]; !ok {
				head[k] = cellName[len(head)]
				e.File.SetCellValue("Sheet1", fmt.Sprintf("%s%d", head[k], 1), k)
			}

			e.File.SetCellValue("Sheet1", fmt.Sprintf("%s%d", head[k], i+2), v)
			j++
		}
	}
	return nil
}

func (e *Excel) SaveAs(name string) error {
	return e.File.SaveAs(name)
}
