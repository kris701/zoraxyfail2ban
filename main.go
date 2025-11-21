package main

import (
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

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
		VersionMajor:  1,
		VersionMinor:  1,
		VersionPatch:  1,
		UIPath:        UI_PATH,
	})
	if err != nil {
		panic(err)
	}

	embedWebRouter := plugin.NewPluginEmbedUIRouter(PLUGIN_ID, &content, WEB_ROOT, UI_PATH)
	embedWebRouter.RegisterTerminateHandler(func() {
		fmt.Println("Fail2Ban Plugin Exited")
	}, nil)

	embedWebRouter.HandleFunc("/api/filter", processFilter, nil)
	embedWebRouter.HandleFunc("/api/jail", processJail, nil)
	embedWebRouter.HandleFunc("/api/getstatus", getStatus, nil)
	embedWebRouter.HandleFunc("/api/ban", banIp, nil)
	embedWebRouter.HandleFunc("/api/unban", unBanIp, nil)

	http.Handle(UI_PATH, embedWebRouter.Handler())
	fmt.Println("Fail2Ban started at http://127.0.0.1:" + strconv.Itoa(runtimeCfg.Port))
	err = http.ListenAndServe("127.0.0.1:"+strconv.Itoa(runtimeCfg.Port), nil)
	if err != nil {
		panic(err)
	}
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("fail2ban-client", "status", "zoraxy")
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	w.Header().Set("Content-Type", "application/json")

	currentlyFailedRegex := regexp.MustCompile(`(?m)Currently failed:\s*(.*)`)
	totalFailedRegex := regexp.MustCompile(`(?m)Total failed:\s*(.*)`)
	currentlyBannedRegex := regexp.MustCompile(`(?m)Currently banned:\s*(.*)`)
	totalBannedRegex := regexp.MustCompile(`(?m)Total banned:\s*(.*)`)
	bannedIPListRegex := regexp.MustCompile(`(?m)Banned IP list:\s*(.*)`)

	currentlyFailedMatches := currentlyFailedRegex.FindAllStringSubmatch(string(stdout), -1)
	totalFailedMatches := totalFailedRegex.FindAllStringSubmatch(string(stdout), -1)
	currentlyBannedMatches := currentlyBannedRegex.FindAllStringSubmatch(string(stdout), -1)
	totalBannedMatches := totalBannedRegex.FindAllStringSubmatch(string(stdout), -1)
	bannedIPListMatches := bannedIPListRegex.FindAllStringSubmatch(string(stdout), -1)

	response := map[string]any{}

	if len(currentlyFailedMatches) > 0 {
		response["currentlyFailed"] = currentlyFailedMatches[0][1]
	}
	if len(totalFailedMatches) > 0 {
		response["totalFailed"] = totalFailedMatches[0][1]
	}
	if len(currentlyBannedMatches) > 0 {
		response["currentlyBanned"] = currentlyBannedMatches[0][1]
	}
	if len(totalBannedMatches) > 0 {
		response["totalBanned"] = totalBannedMatches[0][1]
	}
	if len(bannedIPListMatches) > 0 {
		response["bannedIps"] = strings.Split(bannedIPListMatches[0][1], " ")
	}

	if err2 := json.NewEncoder(w).Encode(response); err2 != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func banIp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		bytedata, _ := io.ReadAll(r.Body)
		reqBodyString := string(bytedata)
		fmt.Println("Manually banning IP:" + reqBodyString)
		cmd := exec.Command("fail2ban-client", "set", "zoraxy", "banip", reqBodyString)
		cmd.Run()
		cmd.Wait()
		return
	}
	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}

func unBanIp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		bytedata, _ := io.ReadAll(r.Body)
		reqBodyString := string(bytedata)
		fmt.Println("Manually unbanning IP:" + reqBodyString)
		cmd := exec.Command("fail2ban-client", "set", "zoraxy", "unbanip", reqBodyString)
		cmd.Run()
		cmd.Wait()
		return
	}
	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}

func processFilter(w http.ResponseWriter, r *http.Request) {
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
		err = cmd2.Start()
		if err != nil {
			panic(err)
		}
		cmd2.Wait()
		cmd3 := exec.Command("fail2ban-client", "reload", "--restart")
		cmd3.Run()
		cmd3.Wait()
		return
	}
	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}

func processJail(w http.ResponseWriter, r *http.Request) {
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
		err = cmd2.Start()
		if err != nil {
			panic(err)
		}
		cmd2.Wait()
		cmd3 := exec.Command("fail2ban-client", "reload", "--restart")
		cmd3.Run()
		cmd3.Wait()
		return
	}
	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}
