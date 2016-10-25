// Code generated by "enumer -type OptionType,NodeType -output node_enumer.go"; DO NOT EDIT

package volgen

import "fmt"

const _OptionType_name = "OPT_NONEOPT_STRINGOPT_INTOPT_DOUBLEOPT_BOOLOPT_MAX"

var _OptionType_index = [...]uint8{0, 8, 18, 25, 35, 43, 50}

func (i OptionType) String() string {
	if i < 0 || i >= OptionType(len(_OptionType_index)-1) {
		return fmt.Sprintf("OptionType(%d)", i)
	}
	return _OptionType_name[_OptionType_index[i]:_OptionType_index[i+1]]
}

var _OptionTypeNameToValue_map = map[string]OptionType{
	_OptionType_name[0:8]:   0,
	_OptionType_name[8:18]:  1,
	_OptionType_name[18:25]: 2,
	_OptionType_name[25:35]: 3,
	_OptionType_name[35:43]: 4,
	_OptionType_name[43:50]: 5,
}

func OptionTypeString(s string) (OptionType, error) {
	if val, ok := _OptionTypeNameToValue_map[s]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to OptionType values", s)
}

const _NodeType_name = "TYPE_NONETYPE_XLATORTYPE_TARGETTYPE_MAX"

var _NodeType_index = [...]uint8{0, 9, 20, 31, 39}

func (i NodeType) String() string {
	if i < 0 || i >= NodeType(len(_NodeType_index)-1) {
		return fmt.Sprintf("NodeType(%d)", i)
	}
	return _NodeType_name[_NodeType_index[i]:_NodeType_index[i+1]]
}

var _NodeTypeNameToValue_map = map[string]NodeType{
	_NodeType_name[0:9]:   0,
	_NodeType_name[9:20]:  1,
	_NodeType_name[20:31]: 2,
	_NodeType_name[31:39]: 3,
}

func NodeTypeString(s string) (NodeType, error) {
	if val, ok := _NodeTypeNameToValue_map[s]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to NodeType values", s)
}