package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/kovey/cli-go/env"
	"github.com/kovey/debug-go/debug"
	"github.com/kovey/debug-go/run"
	"github.com/kovey/kom/internal/html"
)

var serv = &http.Server{}

func Shutdown() {
	serv.Shutdown(context.Background())
}

func Run(s *Serv) {
	port, _ := env.GetInt("SERV_PORT")
	http.HandleFunc("/ko/rpc", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				run.Panic(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
			}
		}()
		r.ParseForm()
		servName := r.Form.Get("q")
		tmp := s
		if servName != "" {
			tmp = &Serv{Version: s.Version}
			for _, str := range s.Structs {
				if strings.Contains(strings.ToLower(str.Name), strings.ToLower(servName)) {
					tmp.Structs = append(tmp.Structs, str)
				}
			}
		}
		w.WriteHeader(http.StatusOK)
		tpl := template.Must(template.New("list").Parse(html.List_Html))
		if err := tpl.Execute(w, tmp); err != nil {
			debug.Erro(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
			return
		}
	})
	http.HandleFunc("/ko/rpc/interface", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				run.Panic(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
			}
		}()
		r.ParseForm()
		servName := r.Form.Get("q")
		str := s.Get(servName)
		if str == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf("service[%s] not found", servName)))
			return
		}
		funcName := r.Form.Get("m")
		method := str.Get(funcName)
		if method == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf("method[%s] of service[%s] not found", funcName, servName)))
			return
		}

		tpl := template.Must(template.New("interface").Parse(html.Interface_Html))
		w.WriteHeader(http.StatusOK)
		if err := tpl.Execute(w, method.DoHtml(servName, funcName, s.Version)); err != nil {
			debug.Erro(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
			return
		}
	})
	http.HandleFunc("/ko/rpc/do", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				run.Panic(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
			}
		}()
		r.ParseForm()
		servName := r.Form.Get("q")
		str := s.Get(servName)
		if str == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf("service[%s] not found", servName)))
			return
		}
		funcName := r.Form.Get("m")
		method := str.Get(funcName)
		if method == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf("method[%s] of service[%s] not found", funcName, servName)))
			return
		}

		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(err.Error()))
			return
		}

		var data map[string]string
		if err := json.Unmarshal(body, &data); err != nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(err.Error()))
			return
		}

		args, err := method.ParseArgs(data)
		if err != nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(err.Error()))
			return
		}

		out, err := method.Call(args)
		w.WriteHeader(http.StatusOK)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		for index, o := range out {
			if index > 0 {
				w.Write([]byte("<br>"))
			}

			if content, err := json.Marshal(o); err == nil {
				w.Write(content)
			}
		}
	})
	serv.Addr = fmt.Sprintf("%s:%d", os.Getenv("SERV_HOST"), port+10000)
	serv.ListenAndServe()
}
