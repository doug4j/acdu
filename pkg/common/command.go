package common

import (
	"fmt"
	"os"
	"reflect"

	"github.com/spf13/cobra"
	validator "gopkg.in/go-playground/validator.v9"
)

func MarkCommand(cmd *cobra.Command, parm interface{}, parmMap map[string]string, validator *validator.Validate) {

	parmType := reflect.TypeOf(parm)
	// if parmType.Kind() != reflect.Struct {
	// 	LogExit(fmt.Sprintf("Provided parameter is not Struct it is %v", parmType))
	// }

	structReflectValue := reflect.ValueOf(parm)
	//structReflectValue := reflect.ValueOf(&structReflectValue0).Elem()

	LogInfo(fmt.Sprintf("field count: %v", parmType.NumField()))

	len := structReflectValue.NumField()
	// len := parmType.FieldAlign()
	for i := 0; i < len; i++ {
		structField := parmType.Field(i)
		// LogInfo(fmt.Sprintf("structReflectValue:%#v", structReflectValue))

		f := structReflectValue.Field(i)
		theType := f.Type()
		theValue := f.Interface()

		//fieldType := reflect.TypeOf(structField)
		LogInfo(fmt.Sprintf("%v : fieldName:%v, fieldType:%v, value:[%v]", i, structField.Name, theType.Name(), theValue))

		LogInfo(fmt.Sprintf("Package: [%v]", parmType.PkgPath()))
		itm := reflect.ValueOf(parmType.PkgPath() + "/Arg" + structField.Name)
		LogInfo(fmt.Sprintf("itm:%v", itm.String()))

		arg, err := ParseArgFromField(structField)
		if err != nil {
			LogExit(fmt.Sprintf("Parse Arg from Type '%v' Field error: %v", parmType.PkgPath()+"/"+parmType.Name(), err))
		}
		help := structField.Tag.Get("help")

		// LogInfo(fmt.Sprintf("%#v", arg))

		switch structField.Type.Kind() {
		case reflect.String:
			parmMap[structField.Name] = ""

			var strVal string
			// strVal = f.String()
			// strVal, ok := theValue.(string)
			// if !ok {
			// 	LogExit(f err.Error())
			// 	os.Exit(1)
			// }
			//cmd.Flags().StringToStringP()
			cmd.Flags().StringVarP(&strVal, arg.LongName, arg.ShortName, "", help)
		case reflect.Bool:
			var boolVal bool
			boolVal = f.Bool()
			// boolVal, ok := theValue.(bool)
			// if !ok {
			// 	LogExit(fmt.Sprintf("%v : fieldType:%v, value:%v cannot be converted to a string", i, theType.Name(), theValue))
			// 	LogError(err.Error())
			// 	os.Exit(1)
			// }
			cmd.Flags().BoolVarP(&boolVal, arg.LongName, arg.ShortName, false, help)
		case reflect.Int:
			var intVal int64
			intVal = f.Int()
			// intVal, ok := theValue.(int)
			// if !ok {
			// 	LogExit(fmt.Sprintf("%v : fieldType:%v, value:%v cannot be converted to a string", i, theType.Name(), theValue))
			// 	LogError(err.Error())
			// 	os.Exit(1)
			// }
			cmd.Flags().Int64VarP(&intVal, arg.LongName, arg.ShortName, 0, help)
			//cmd.Flags().IntVarP(&intVal, arg.LongName, arg.ShortName, 0, help)

			//LogInfo(fmt.Sprintf("field %v: structField=%v type=%v", i, structField, structField.Type))
		}
		if arg.Required {
			if err := cmd.MarkFlagRequired(arg.LongName); err != nil {
				LogError(err.Error())
				os.Exit(1)
			}
		}
	}

}
