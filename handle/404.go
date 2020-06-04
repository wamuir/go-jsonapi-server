package handle

import (
	"net/http"

	"github.com/wamuir/go-jsonapi-server/model"
)

func (env *Environment) Handle404(w http.ResponseWriter, r *http.Request) {

	e := model.MakeError(http.StatusNotFound)
	e.Code = "f7519b"
	env.Fail(w, r, e)
	return
}
