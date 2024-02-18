package format

import (
	"github.com/henges/trackrouter/model"
	"github.com/samber/lo"
	"strings"
)

func LinksMatchResult(l *model.LinksMatchResult) string {

	return strings.Join(lo.Values(l.Links), "\n")
}
