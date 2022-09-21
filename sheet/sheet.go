package sheet

import (
	"github.com/d5/tengo/v2"
	"github.com/ipsusila/tnglib"
	"github.com/xuri/excelize/v2"
)

// modules name
var (
	Name = "sheet"
)

var (
	// Module registered here
	sheetModule = map[string]tengo.Object{
		"st_cell_formula_type_array": &tengo.String{Value: excelize.STCellFormulaTypeArray},
	}
)

func init() {
	// register module
	tnglib.RegisterModule(Name, sheetModule)
}
