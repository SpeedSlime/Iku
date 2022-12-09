package hotupdate

import (
	"fmt"
	"net/http"

	db "github.com/SpeedSlime/Asahi/database"
	"github.com/SpeedSlime/Asahi"
	"github.com/SpeedSlime/Asahi/reply"
)

func HotUpdateVersionGetRoute(w http.ResponseWriter, r *http.Request) {
	v := XVersion{Platform: asahi.Parameter(r, "device")}
	if has, _ := db.Select(&v); !has { 
		asahi.Handle(reply.RespondWithResult(w, http.StatusNotFound, ""), "HotUpdateVersionGetRoute"); return
	}
	asahi.Handle(reply.RespondWithJSON(w, http.StatusOK, v), "HotUpdateVersionGetRoute")
}

func HotUpdateDownloadGetRoute(w http.ResponseWriter, r *http.Request) {
	asahi.Handle(reply.RespondWithFile(w, http.StatusOK, fmt.Sprintf("%s/%s", asahi.Parameter(r, "device"), asahi.Parameter(r, "file"))), "HotUpdateDownloadGetRoute")
}