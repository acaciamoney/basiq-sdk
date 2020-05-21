package utilities

import "strings"

type FilterBuilder struct {
	filter []string
}

func (fb *FilterBuilder) Eq(field string, value string) *FilterBuilder {
	fb.filter = append(fb.filter, field+".eq('"+value+"')")

	return fb
}

func (fb *FilterBuilder) Gt(field string, value string) *FilterBuilder {
	fb.filter = append(fb.filter, field+".gt('"+value+"')")

	return fb
}

func (fb *FilterBuilder) Gteq(field string, value string) *FilterBuilder {
	fb.filter = append(fb.filter, field+".gteq('"+value+"')")

	return fb
}

func (fb *FilterBuilder) Lt(field string, value string) *FilterBuilder {
	fb.filter = append(fb.filter, field+".lt('"+value+"')")

	return fb
}

func (fb *FilterBuilder) Lteq(field string, value string) *FilterBuilder {
	fb.filter = append(fb.filter, field+".lteq('"+value+"')")

	return fb
}

func (fb *FilterBuilder) Bt(field string, valueOne string, valueTwo string) *FilterBuilder {
	fb.filter = append(fb.filter, field+".bt('"+valueOne+"', '"+valueTwo+"')")

	return fb
}

func (fb *FilterBuilder) ToString() string {
	return strings.Join(fb.filter, ",")
}

func (fb *FilterBuilder) GetFilter() string {
	return "filter=" + strings.Join(fb.filter, ",")
}

func (fb *FilterBuilder) SetFilter(filters []string) *FilterBuilder {
	fb.filter = filters

	return fb
}
