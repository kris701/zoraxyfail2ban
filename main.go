package main

import (
	"embed"
	_ "embed"
	"fmt"
	"net/http"
	"strconv"
	"encoding/json"
	"os/exec"

	plugin "github.com/kris701/zoraxyfail2ban/mod/zoraxy_plugin"
)

const (
	PLUGIN_ID = "zoraxyfail2ban"
	UI_PATH   = "/"
	WEB_ROOT  = "/www"
)

var content embed.FS

func main() {
	runtimeCfg, err := plugin.ServeAndRecvSpec(&plugin.IntroSpect{
		ID:            "zoraxyfail2ban",
		Name:          "Fail2Ban",
		Author:        "Kristian Skov Johansen",
		AuthorContact: "kris701kj@gmail.com",
		Description:   "A plugin to interact with fail2ban",
		URL:           "https://github.com/kris701/zoraxyfail2ban",
		Type:          plugin.PluginType_Utilities,
		VersionMajor:  0,
		VersionMinor:  1,
		VersionPatch:  0,
		UIPath: UI_PATH,
	})
	if err != nil {
		panic(err)
	}

	embedWebRouter := plugin.NewPluginEmbedUIRouter(PLUGIN_ID, &content, WEB_ROOT, UI_PATH)
	embedWebRouter.RegisterTerminateHandler(func() {
		fmt.Println("Fail2Ban Plugin Exited")
	}, nil)
	
	embedWebRouter.HandleFunc("/api/getstatus", func(w http.ResponseWriter, r *http.Request) {
		cmd := exec.Command("bash", "-c", "fail2ban-client status zoraxy")
		stdout, err := cmd.Output()
		if err != nil {
			fmt.Println(err.Error())
		}
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{"message": string(stdout)}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}, nil)

	http.Handle(UI_PATH, embedWebRouter.Handler())
	fmt.Println("Fail2Ban started at http://127.0.0.1:" + strconv.Itoa(runtimeCfg.Port))
	err = http.ListenAndServe("127.0.0.1:"+strconv.Itoa(runtimeCfg.Port), nil)
	if err != nil {
		panic(err)
	}

}
