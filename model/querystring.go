package model

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/wamuir/go-jsonapi-core"
)

type QueryParams struct {
	Limit   int64
	Offset  int64
	Include KeyRing
	Sort    string
}

type Parameters map[string]Parameter

type Parameter struct {
	Allowed bool
	Default int64
	Minimum int64
	Maximum int64
}

type KeyRing [][]string

func (ring KeyRing) IsValidAgainst(keys []string) bool {

	for _, k := range ring {

		if len(k) > 0 && !stringInSlice(k[0], keys) {
			return false
		}

	}

	return true
}

func (ring KeyRing) SplitOn(key string) KeyRing {

	var r KeyRing

	for _, k := range ring {

		if len(k) == 0 {
			continue
		}

		part, remainder := k[0], k[1:]
		if part == key {
			r = append(r, remainder)
		}
	}

	return r
}

func (ring KeyRing) Requests(key string) bool {

	for _, k := range ring {
		if len(k) == 0 {
			continue
		}

		if k[0] == key {
			return true
		}
	}
	return false
}

func ParseIntegerParameter(parameter string, q url.Values, params Parameters) (int64, *core.Error) {

	var (
		value  int64
		err    error
		errObj *core.Error
	)

	switch entries := q[parameter]; {

	case len(entries) == 0:

		value = params[parameter].Default

	case len(entries) == 1:

		value, err = strconv.ParseInt(entries[0], 10, 64)
		if err != nil {
			errObj = core.MakeError(http.StatusBadRequest)
			errObj.Code = "ab1fb9"
			errObj.Title = "Invalid query string"
			errObj.Detail = fmt.Sprintf(
				"Unable to parse parameter %s as an integer",
				parameter,
			)
			return value, errObj

		} else if value < params[parameter].Minimum {
			errObj = core.MakeError(http.StatusBadRequest)
			errObj.Code = "ee05db"
			errObj.Title = "Invalid query string"
			errObj.Detail = fmt.Sprintf(
				"Value of parameter %s less than minimum %d",
				parameter,
				params[parameter].Minimum,
			)
			return value, errObj

		} else if value > params[parameter].Maximum {
			errObj = core.MakeError(http.StatusBadRequest)
			errObj.Code = "2f6967"
			errObj.Title = "Invalid query string"
			errObj.Detail = fmt.Sprintf(
				"Value of parameter %s greater than maximum %d",
				parameter,
				params[parameter].Maximum,
			)
			return value, errObj

		}

	case len(entries) > 1:

		errObj = core.MakeError(http.StatusBadRequest)
		errObj.Code = "d5ea49"
		errObj.Title = "Invalid query string"
		errObj.Detail = fmt.Sprintf(
			"Superfluous parameter: %s",
			parameter,
		)
		return value, errObj

	}

	return value, nil

}

func ParseQueryString(u *url.URL, params Parameters) (QueryParams, *core.Error) {

	var (
		queryParams QueryParams
		errObj      *core.Error
	)

	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		errObj = core.MakeError(http.StatusBadRequest)
		errObj.Code = "fc6f74"
		errObj.Title = "Invalid query string"
		errObj.Detail = err.Error()
		return queryParams, errObj
	}

	for par := range q {
		if !params[par].Allowed {
			errObj = core.MakeError(http.StatusBadRequest)
			errObj.Code = "6d01c2"
			errObj.Title = "Invalid query string"
			errObj.Detail = fmt.Sprintf("Parameter not allowed: %s", par)
			return queryParams, errObj
		}
	}

	// Parse page[limit]
	queryParams.Limit, errObj = ParseIntegerParameter(
		"page[limit]", q, params,
	)
	if errObj != nil {
		return queryParams, errObj
	}

	// Parse page[offset]
	queryParams.Offset, errObj = ParseIntegerParameter(
		"page[offset]", q, params,
	)
	if errObj != nil {
		return queryParams, errObj
	}

	// Parse sort
	queryParams.Sort = q.Get("sort")

	// Parse include
	set := map[string]bool{}

	for _, i := range q["include"] {
		for _, j := range strings.Split(i, ",") {
			set[j] = true
		}
	}

	for key := range set {
		teeth := strings.Split(key, ".")
		if int64(len(teeth)) > params["include"].Maximum {
			errObj = core.MakeError(http.StatusBadRequest)
			errObj.Code = "8bfeb8"
			errObj.Title = "Invalid query string"
			errObj.Detail = "Request for included resources exceeds maximum traversal depth"
			return queryParams, errObj
		}
		queryParams.Include = append(queryParams.Include, teeth)
	}

	return queryParams, nil
}

func stringInSlice(a string, b []string) bool {
	for _, c := range b {
		if c == a {
			return true
		}
	}
	return false
}
