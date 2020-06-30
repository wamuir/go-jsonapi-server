package model

import (
	"net/url"
	"strconv"

	"github.com/wamuir/go-jsonapi-core"
)

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func makePaginationLink(base url.URL, ref *url.URL, limit, offset int64) string {

	v := url.Values{}
	v.Set("page[limit]", strconv.FormatInt(limit, 10))
	v.Set("page[offset]", strconv.FormatInt(offset, 10))

	ref.RawQuery = v.Encode()

	return base.ResolveReference(ref).String()
}

func paginate(collection core.Collection, base url.URL, ref *url.URL, limit, offset, count int64) (core.Collection, core.LinksObject) {

	links := make(core.LinksObject)
	links["first"] = makePaginationLink(base, ref, limit, 0)
	links["last"] = makePaginationLink(base, ref, limit, maxInt64(0, count-(count-offset%limit)%limit))
	if offset > 0 {
		links["prev"] = makePaginationLink(base, ref, limit, maxInt64(0, offset-limit))
	}
	if (offset + limit) < count {
		links["next"] = makePaginationLink(base, ref, limit, offset+limit)
	}
	links["self"] = makePaginationLink(base, ref, limit, offset)

	return collection, links

}
