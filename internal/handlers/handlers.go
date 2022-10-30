package handlers

import (
	"net/http"
	"githob.com/CloudyKit/jet/v6"
) 

var views = jet.NewSet(
	jet.NewOsFileSystemLoader("./html"),
	jet.InDevelopmentMode()
)

func Home(w http.ResponseWriter, r *http.Request)  {
	
}