package commoninterface

type ParamOption func(ParamOptions)

type ParamOptions interface {
	Apply(...ParamOption)
}
