package clause

import (
	"io"

	"github.com/stephenafamo/bob"
)

type IWindow interface {
	SetFrom(string)
	AddPartitionBy(...any)
	AddOrderBy(...any)
	SetMode(string)
	SetStart(any)
	SetEnd(any)
	SetExclusion(string)
}

type WindowDef struct {
	From        string // an existing window name
	orderBy     []any
	partitionBy []any
	Frame
}

func (wi WindowDef) Valid() bool {
	if wi.From != "" {
		return true
	}

	if len(wi.orderBy) != 0 {
		return true
	}

	if len(wi.partitionBy) != 0 {
		return true
	}

	if wi.Frame.Defined {
		return true
	}

	return false
}

func (wi *WindowDef) SetFrom(from string) {
	wi.From = from
}

func (wi *WindowDef) AddPartitionBy(condition ...any) {
	wi.partitionBy = append(wi.partitionBy, condition...)
}

func (wi *WindowDef) AddOrderBy(order ...any) {
	wi.orderBy = append(wi.orderBy, order...)
}

func (wi WindowDef) WriteSQL(w io.Writer, d bob.Dialect, start int) ([]any, error) {
	if wi.From != "" {
		w.Write([]byte(wi.From))
		w.Write([]byte(" "))
	}

	args, err := bob.ExpressSlice(w, d, start, wi.partitionBy, "PARTITION BY ", ", ", " ")
	if err != nil {
		return nil, err
	}

	orderArgs, err := bob.ExpressSlice(w, d, start, wi.orderBy, "ORDER BY ", ", ", "")
	if err != nil {
		return nil, err
	}
	args = append(args, orderArgs...)

	frameArgs, err := bob.ExpressIf(w, d, start, wi.Frame, wi.Frame.Defined, " ", "")
	if err != nil {
		return nil, err
	}
	args = append(args, frameArgs...)

	return args, nil
}

type NamedWindow struct {
	Name       string
	Definition any
}

func (n NamedWindow) WriteSQL(w io.Writer, d bob.Dialect, start int) ([]any, error) {
	w.Write([]byte(n.Name))
	w.Write([]byte(" AS ("))
	args, err := bob.Express(w, d, start, n.Definition)
	w.Write([]byte(")"))

	return args, err
}

type Windows struct {
	Windows []NamedWindow
}

func (wi *Windows) AppendWindow(w NamedWindow) {
	wi.Windows = append(wi.Windows, w)
}

func (wi Windows) WriteSQL(w io.Writer, d bob.Dialect, start int) ([]any, error) {
	return bob.ExpressSlice(w, d, start, wi.Windows, "WINDOW ", ", ", "")
}
