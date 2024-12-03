package configs

import "strings"

type TaskSection struct {
	Id      string
	Name    string
	Env     map[string]*ExprStringItem
	Timeout *ExprIntItem
	Use     string
	If      *ExprStringItem
	With    map[string]*ExprStringItem
	Force   *ExprBoolItem
	Run     *ExprStringItem
}

type TaskDirectiveElement struct {
	Ref  string
	Task *TaskSection
}

type JobSection struct {
	Id      string
	Name    string
	Env     map[string]*ExprIntItem
	Timeout *ExprIntItem
	Force   *ExprBoolItem
	Tasks   []TaskDirectiveElement
}

type ExprValueItem struct {
	Value  interface{}
	isExpr *bool
	Raw    string
	Kind   string
}

func (e *ExprValueItem) HasValue() bool {
	return e.Value != nil
}

func (e *ExprValueItem) IsString() bool {
	return e.Kind == "string"
}

func (e *ExprValueItem) IsInt() bool {
	return e.Kind == "int"
}

func (e *ExprValueItem) IsBool() bool {
	return e.Kind == "bool"
}

func (e *ExprValueItem) IsExpr() bool {
	if e.isExpr == nil {
		if (len(e.Raw) > 0) && (strings.Index(e.Raw, "${{") == 0) {
			b := true
			e.isExpr = &b
		} else {
			b := false
			e.isExpr = &b
		}
	}

	return *e.isExpr
}

type ExprBoolItem struct {
	ExprValueItem
	Evalutated bool
}

type ExprIntItem struct {
	ExprValueItem
	Evaluated int
}

type ExprStringItem struct {
	ExprValueItem
	Evaluated string
}
