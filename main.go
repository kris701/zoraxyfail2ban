package main

import (
	"embed"
	_ "embed"
	"fmt"
	"net/http"
	"strconv"
	"os"
	"os/exec"
	"io"

	plugin "github.com/kris701/zoraxyfail2ban/mod/zoraxy_plugin"
)

const (
	PLUGIN_ID = "zoraxyfail2ban"
	UI_PATH   = "/"
	WEB_ROOT  = "/www"
)

//go:embed www/*
var content embed.FS

func main() {
	runtimeCfg, err := plugin.ServeAndRecvSpec(&plugin.IntroSpect{
		ID:            "zoraxyfail2ban",
		Name:          "Fail2Ban",
		Author:        "Kristian Skov Johansen",
		AuthorContact: "kris701kj@gmail.com",
		Description:   "A plugin to interact with fail2ban in Zoraxy",
		URL:           "https://github.com/kris701/zoraxyfail2ban",
		Type:          plugin.PluginType_Utilities,
		VersionMajor:  0,
		VersionMinor:  1,
		VersionPatch:  1,
		UIPath: UI_PATH,
	})
	if err != nil {
		panic(err)
	}

	embedWebRouter := plugin.NewPluginEmbedUIRouter(PLUGIN_ID, &content, WEB_ROOT, UI_PATH)
	embedWebRouter.RegisterTerminateHandler(func() {
		fmt.Println("Fail2Ban Plugin Exited")
	}, nil)
	
	embedWebRouter.HandleFunc("/api/filter", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			cmd := exec.Command("cat", "/etc/fail2ban/filter.d/zoraxy.conf")
			stdout, _ := cmd.Output()
			w.Header().Set("Content-Type", "text/html")
			response := string(stdout)
			w.Write([]byte(response))
			return 
		} else if r.Method == http.MethodPost {
			bytedata, _ := io.ReadAll(r.Body)
			reqBodyString := string(bytedata)
			fmt.Println("Updating fail2ban filter config.")
			cmd2 := exec.Command("echo", string(reqBodyString))
			outfile, err := os.Create("/etc/fail2ban/filter.d/zoraxy.conf")
			defer outfile.Close()
			cmd2.Stdout = outfile
			err = cmd2.Start(); if err != nil {
				panic(err)
			}
			cmd2.Wait()
			cmd3 := exec.Command("fail2ban-client", "reload", "--restart")
			cmd3.Run();
			cmd3.Wait();
			return
		}
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}, nil)

	embedWebRouter.HandleFunc("/api/jail", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			cmd := exec.Command("cat", "/etc/fail2ban/jail.local")
			stdout, _ := cmd.Output()
			w.Header().Set("Content-Type", "text/html")
			response := string(stdout)
			w.Write([]byte(response))
			return 
		} else if r.Method == http.MethodPost {
			bytedata, _ := io.ReadAll(r.Body)
			reqBodyString := string(bytedata)
			fmt.Println("Updating fail2ban jail config.")
			cmd2 := exec.Command("echo", string(reqBodyString))
			outfile, err := os.Create("/etc/fail2ban/jail.local")
			defer outfile.Close()
			cmd2.Stdout = outfile
			err = cmd2.Start(); if err != nil {
				panic(err)
			}
			cmd2.Wait()
			cmd3 := exec.Command("fail2ban-client", "reload", "--restart")
			cmd3.Run();
			cmd3.Wait();
			return
		}
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}, nil)
	
	embedWebRouter.HandleFunc("/api/getstatus", func(w http.ResponseWriter, r *http.Request) {
		cmd := exec.Command("fail2ban-client", "status", "zoraxy")
		stdout, err := cmd.Output()
		if err != nil {
			fmt.Println(err.Error())
		}
		w.Header().Set("Content-Type", "text/html")
		response := string(stdout)
		w.Write([]byte(response))
	}, nil)

	embedWebRouter.HandleFunc("/api/ban", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			bytedata, _ := io.ReadAll(r.Body)
			reqBodyString := string(bytedata)
			fmt.Println("Manually banning IP:" + reqBodyString)
			cmd := exec.Command("fail2ban-client", "set", "zoraxy", "banip", reqBodyString)
			cmd.Run();
			cmd.Wait();
			return
		}
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}, nil)

	embedWebRouter.HandleFunc("/api/unban", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			bytedata, _ := io.ReadAll(r.Body)
			reqBodyString := string(bytedata)
			fmt.Println("Manually unbanning IP:" + reqBodyString)
			cmd := exec.Command("fail2ban-client", "set", "zoraxy", "unbanip", reqBodyString)
			cmd.Run();
			cmd.Wait();
			return
		}
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}, nil)

	http.Handle(UI_PATH, embedWebRouter.Handler())
	fmt.Println("Fail2Ban started at http://127.0.0.1:" + strconv.Itoa(runtimeCfg.Port))
	err = http.ListenAndServe("127.0.0.1:"+strconv.Itoa(runtimeCfg.Port), nil)
	if err != nil {
		panic(err)
	}
}
