/**
 * Created with IntelliJ IDEA.
 * User: mario
 * Date: 13.12.12
 * Time: 10:44
 * To change this template use File | Settings | File Templates.
 */
package web

import "io/ioutil"
import "net/http"
import "fmt"
import "github.com/kless/goconfig/config"
import "svn"
import "log"
import "encoding/json"

type PageManager struct {
	Config *config.Config
}

type Page struct {
	Title string
	Body  []byte
}

func (pm *PageManager) StartServer(config *config.Config) {
	pm.Config = config
	port, _ := pm.Config.String("DAEMON","port")

	log.Print("Starting server")

	svnHandle := new(svn.Svn)
	svnHandle.Config = config

	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
			p, _ := pm.loadPage("index")
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, string(p.Body))
		})

	// TODO: Refactor and move the main logic into the svn module
	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
			cmd := svnHandle.SvnLocalInfo()
			out, oerr := cmd.CombinedOutput()

			if oerr != nil {
				log.Fatal(oerr)
			}
			w.Header().Set("Content-Type", "text/plain")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Write(out)
		})

	http.HandleFunc("/switch", func(w http.ResponseWriter, r *http.Request) {
			var cmd string

			branchName := r.URL.Query().Get("branch")
			tagName := r.URL.Query().Get("tag")

			// Crazy comparison to prevent empty string allocation
			if branchName + "FOOBAR" != "FOOBAR" {
				if branchName == "trunk" {
					cmd = svnHandle.SwitchTrunk()
				}

				if branchName != "trunk" {
					cmd = svnHandle.SwitchBranch(branchName)
				}
			}
			if tagName + "FOOBAR" != "FOOBAR" {
				cmd = svnHandle.SwitchTag(tagName)
			}
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(cmd))
		})

	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
			output := svnHandle.UpdateWorkingCopy()
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(output))
		})

	http.HandleFunc("/list/tags", func(w http.ResponseWriter, r *http.Request) {
			encoder := json.NewEncoder(w)
			branches := svnHandle.GetTagList()
			w.Header().Set("Content-Type", "application/json")
			encoder.Encode(branches)
		})

	http.HandleFunc("/list/branches", func(w http.ResponseWriter, r *http.Request) {
			encoder := json.NewEncoder(w)
			tags := svnHandle.GetBranchList()
			w.Header().Set("Content-Type", "application/json")
			encoder.Encode(tags)
		})
	http.ListenAndServe(":" + port, nil)
}

func (pm *PageManager) loadPage(title string) (*Page, error) {
	basedir,_ := pm.Config.String("DAEMON", "static_path")

	if title == "" {
		title = "index"
	}

	filename := basedir + "/" + title + ".html"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}
