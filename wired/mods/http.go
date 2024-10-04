package mods

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/kercre123/wire-os/wired/vars"
)

func ImplHTTP() {
	http.HandleFunc("/api/mods/modify/", ModHTTPHandler)
	http.HandleFunc("/api/mods/current/", ModHTTPHandler)
	http.HandleFunc("/api/mods/needsrestart/", ModHTTPHandler)
	http.HandleFunc("/api/restartvic", ModHTTPHandler)

	http.HandleFunc("/api/mods/custom/TestBootAnim", BootAnim_HTTP)
	http.HandleFunc("/api/mods/wakeword/", WakeWord_HTTP)
}

func ModHTTPHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case strings.HasPrefix(r.URL.Path, "/api/mods/modify/"):
		modFromURL := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, "/api/mods/modify/"))
		jsonFromReq, err := io.ReadAll(r.Body)
		if err != nil {
			vars.HTTPError(w, r, "error reading request body: "+err.Error())
			return
		}
		mod, err := vars.FindMod(modFromURL)
		if err != nil {
			vars.HTTPError(w, r, "mod does not exist")
			return
		}
		err = mod.Do(vars.SysRoot, string(jsonFromReq))
		if err != nil {
			vars.HTTPError(w, r, "error running mod: "+err.Error())
			return
		}
		vars.HTTPSuccess(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/mods/current/"):
		modFromURL := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, "/api/mods/current/"))
		mod, err := vars.FindMod(modFromURL)
		if err != nil {
			vars.HTTPError(w, r, "mod does not exist")
			return
		}
		fmt.Fprint(w, mod.Current())
	case strings.HasPrefix(r.URL.Path, "/api/mods/needsrestart/"):
		modFromURL := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, "/api/mods/needsrestart/"))
		mod, err := vars.FindMod(modFromURL)
		if err != nil {
			vars.HTTPError(w, r, "mod does not exist")
			return
		}
		fmt.Fprint(w, mod.RestartRequired())
	case strings.HasPrefix(r.URL.Path, "/api/restartvic"):
		vars.RestartVic()
		vars.HTTPSuccess(w, r)
	default:
		vars.HTTPError(w, r, "404 not found")
	}
}
