package common

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var VerboseLogging bool

type Argument struct {
	Required        bool
	LongName        string
	ShortName       string
	HasDefaultValue bool
	DefaultValue    string
}

func ParseArgFromField(field reflect.StructField) (argument Argument, err error) {
	if len(field.Name) < 2 {
		return argument, fmt.Errorf("Arg item field [%v] has a name that is less than 2, this is illegal", field.Name)
	}

	defaultName := strings.ToLower(field.Name[0:1]) + field.Name[1:]
	argument.LongName = defaultName
	rawArgStr := field.Tag.Get("arg")
	argItems := strings.Split(rawArgStr, ",")
	for index, argItem := range argItems {
		nameValue := strings.Split(argItem, "=")
		if len(nameValue) != 2 {
			return argument, fmt.Errorf("Arg item at %v index for field '%v' is not a single '='", index, field.Name)
		}
		tagName := strings.ToLower(nameValue[0])
		tagValue := nameValue[1]
		err = processArg(&argument, field.Name, tagName, tagValue)
		if err != nil {
			return argument, err
		}
	}
	return argument, nil
}

func processArg(argument *Argument, fieldName, tagName, tagValue string) error {
	switch tagName {
	case "required":
		required, err := strconv.ParseBool(tagValue)
		if err != nil {
			return fmt.Errorf("Arg field %v for 'required' field is not a boolean, it's name/value %v/[%v]", fieldName, tagName, tagValue)
		}
		argument.Required = required
	case "longname":
		if len(tagValue) > 0 {
			argument.LongName = tagValue
		}
	case "defaultvalue":
		argument.DefaultValue = tagValue
		argument.HasDefaultValue = true
	case "shortname":
		if len(tagValue) > 1 {
			return fmt.Errorf("Arg field %v for 'shortname' field's value is greater than 1 character, it's name/value %v/[%v]", fieldName, tagName, tagValue)
		}
		argument.ShortName = strings.ToLower(tagValue)
	}
	return nil
}

//AttachStringArg uses reflection to read the provided struct to determine the arguments. otherArgs has the first argument is the defaultDefault value that overrides anything defined in the struct argument tag.
func AttachStringArg(cmd *cobra.Command, parmType reflect.Type, variableName string, variableValue *string, otherArgs ...string) {
	arg, rawHelp := parseArg(parmType, variableName)
	var defaultValue string
	if arg.HasDefaultValue {
		defaultValue = arg.DefaultValue
	}
	if len(otherArgs) > 0 {
		defaultValue = otherArgs[0]
	}
	cmd.Flags().StringVarP(variableValue, arg.LongName, arg.ShortName, defaultValue, rationalizeHelp(arg, rawHelp))
	processRequiredArg(cmd, arg)
}

func AttachBoolArg(cmd *cobra.Command, parmType reflect.Type, variableName string, variableValue *bool) {
	arg, rawHelp := parseArg(parmType, variableName)
	var err error
	var defaultValue bool
	if arg.HasDefaultValue {
		defaultValue, err = strconv.ParseBool(arg.DefaultValue)
		if err != nil {
			LogExit(fmt.Sprintf("Fatal mis-configuration. Field %v could not process default value: %v", parmType.Name()+"."+variableName, arg.DefaultValue))
		}
	}
	cmd.Flags().BoolVarP(variableValue, arg.LongName, arg.ShortName, defaultValue, rationalizeHelp(arg, rawHelp))
	processRequiredArg(cmd, arg)
}

func AttachIntArg(cmd *cobra.Command, parmType reflect.Type, variableName string, variableValue *int) {
	arg, rawHelp := parseArg(parmType, variableName)
	var err error
	var defaultValue int
	if arg.HasDefaultValue {
		defaultValue, err = strconv.Atoi(arg.DefaultValue)
		if err != nil {
			LogExit(fmt.Sprintf("Fatal mis-configuration. Field %v could not process default value: %v", parmType.Name()+"."+variableName, arg.DefaultValue))
		}
	}
	cmd.Flags().IntVarP(variableValue, arg.LongName, arg.ShortName, defaultValue, rationalizeHelp(arg, rawHelp))
	processRequiredArg(cmd, arg)
}

func rationalizeHelp(arg Argument, rawHelp string) (help string) {
	if arg.Required {
		help = "MANDATORY: "
	} else {
		help = "optional: "
	}
	help = help + rawHelp
	return help
}

func parseArg(parmType reflect.Type, variableName string) (arg Argument, rawHelp string) {
	var field reflect.StructField
	var has bool
	var err error
	if field, has = parmType.FieldByName(variableName); !has {
		LogExit(fmt.Sprintf("Fatal mis-configuration by the variable [%v]", variableName))
	}
	arg, err = ParseArgFromField(field)
	if err != nil {
		LogExit(fmt.Sprintf("Fatal mis-configuration, could not get argements from field [%v]", field))
	}
	rawHelp = field.Tag.Get("help")
	return arg, rawHelp
}

func processRequiredArg(cmd *cobra.Command, arg Argument) {
	if arg.Required {
		if err := cmd.MarkFlagRequired(arg.LongName); err != nil {
			LogExit(fmt.Sprintf("Fatal mis-configuration, could not mark required field: %v", err.Error()))
		}
	}
}
