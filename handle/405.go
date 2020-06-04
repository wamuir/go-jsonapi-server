package handle

import (
	"net/http"

	"github.com/wamuir/go-jsonapi-server/model"
)

func (env *Environment) Handle405(w http.ResponseWriter, r *http.Request) {

	e := model.MakeError(http.StatusMethodNotAllowed)
	e.Code = "f7519b"
	env.Fail(w, r, e)
	return
}
