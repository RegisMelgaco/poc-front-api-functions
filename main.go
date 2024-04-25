package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"rogchap.com/v8go"
)

type config []moduleConfig

type moduleConfig struct {
	Module       string       `json:"module"`
	Owner        string       `json:"owner"`
	Contact      string       `json:"contact"`
	MaxErrorRate string       `json:"max_error_rate"`
	Slack        string       `json:"slack"`
	Files        []fileConfig `json:"files"`
}

type fileConfig struct {
	Path      string          `json:"path"`
	Functions []functioConfig `json:"functions"`
}

type functioConfig struct {
	Name        string `json:"name"`
	RoutePath   string `json:"route_path"`
	RouteMethod string `json:"route_method"`
}

func main() {
	cfgFile, err := os.ReadFile("config.json")
	if err != nil {
		fmt.Println(err)

		return
	}

	var cfg config
	if err := json.Unmarshal(cfgFile, &cfg); err != nil {
		fmt.Println(err)

		return
	}

	mux := http.NewServeMux()

	for _, modCfg := range cfg {
		iso := v8go.NewIsolate() // creates a new JavaScript VM
		v8Ctx := v8go.NewContext(iso)

		for _, fileCfg := range modCfg.Files {
			filePath := fmt.Sprintf("funcs/%s/%s", modCfg.Module, fileCfg.Path)

			file, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Println(err)

				return
			}

			script, err := iso.CompileUnboundScript(string(file), filePath, v8go.CompileOptions{}) // compile script to get cached data
			if err != nil {
				fmt.Println(err)

				return
			}

			val, err := script.Run(v8Ctx)
			if err != nil {
				fmt.Println(err)

				return
			}

			fmt.Println(val)

			state := ""

			for _, funCfg := range fileCfg.Functions {
				route := fmt.Sprintf("%s /%s%s", funCfg.RouteMethod, modCfg.Module, funCfg.RoutePath)

				fmt.Printf("registering function=%s at pattern=%s\n", funCfg.Name, route)
				mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
					body := []byte("null")
					if b, err := io.ReadAll(r.Body); err == nil && len(b) != 0 {
						body = b
					}

					s := state
					s = strings.TrimLeft(s, `"`)
					s = strings.TrimRight(s, `"`)
					s = strings.ReplaceAll(s, "\\\"", "\"")

					funCall := fmt.Sprintf(`%s({"body":%s},%s);`, funCfg.Name, string(body), s)

					fmt.Println(funCall)

					val, err := v8Ctx.RunScript(funCall, "")
					if err != nil {
						fmt.Println(err)
						w.WriteHeader(http.StatusInternalServerError)

						return
					}

					state = val.String()

					w.Header().Add("Content-Type", "application/json")

					fmt.Println(state)

					w.Write([]byte(state))
				})
			}
		}
	}

	fmt.Println("start listening on :3000")
	err = http.ListenAndServe(":3000", mux)
	if err != nil {
		fmt.Println(err)

		return
	}
}
