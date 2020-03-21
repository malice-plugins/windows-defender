package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/gorilla/mux"
	"github.com/malice-plugins/pkgs/clitable"
	"github.com/malice-plugins/pkgs/database"
	"github.com/malice-plugins/pkgs/database/elasticsearch"
	"github.com/malice-plugins/pkgs/utils"
	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	name     = "windows-defender"
	category = "av"
)

var (
	// Version stores the plugin's version
	Version string
	// BuildTime stores the plugin's build time
	BuildTime string

	path string

	// es is the elasticsearch database object
	es elasticsearch.Database
)

type pluginResults struct {
	ID   string      `json:"id" structs:"id,omitempty"`
	Data ResultsData `json:"windows_defender" structs:"windows_defender"`
}

// WindowsDefender json object
type WindowsDefender struct {
	Results ResultsData `json:"windows_defender" structs:"windows_defender"`
}

// ResultsData json object
type ResultsData struct {
	Infected bool   `json:"infected" structs:"infected"`
	Result   string `json:"result" structs:"result"`
	Engine   string `json:"engine" structs:"engine"`
	Updated  string `json:"updated" structs:"updated"`
	MarkDown string `json:"markdown,omitempty" structs:"markdown,omitempty"`
}

func assert(err error) {
	if err != nil {
		log.WithFields(log.Fields{
			"plugin":   strings.Replace(name, "-", "_", -1),
			"category": category,
			"path":     path,
		}).Fatal(err)
	}
}

// RunCommand runs cmd on file
func RunCommand(ctx context.Context, cmd string, args ...string) (string, error) {

	var c *exec.Cmd

	if ctx != nil {
		c = exec.CommandContext(ctx, cmd, args...)
	} else {
		c = exec.Command(cmd, args...)
	}

	output, err := c.CombinedOutput()
	if err != nil {
		return string(output), err
	}

	// check for exec context timeout
	if ctx != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("command %s timed out", cmd)
		}
	}

	return string(output), nil
}

// AvScan performs antivirus scan
func AvScan(timeout int) WindowsDefender {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	// needs to be run from the /loadlibrary folder
	if err := os.Chdir("/loadlibrary"); err != nil {
		assert(err)
	}
	// will change back to the /malware folder when func returns
	defer os.Chdir("/malware")

	results, err := ParseWinDefOutput(RunCommand(ctx, "./mpclient", path))
	assert(err)

	return WindowsDefender{
		Results: results,
	}
}

// ParseWinDefOutput convert windef output into ResultsData struct
func ParseWinDefOutput(windefout string, err error) (ResultsData, error) {

	// root@d9f8dca1d59e:/loadlibrary# ./mpclient /malware/EICAR
	// main(): The map file wasn't found, symbols wont be available
	// main(): Scanning /malware/EICAR...
	// EngineScanCallback(): Scanning input
	// EngineScanCallback(): Threat Virus:DOS/EICAR_Test_File identified.

	if err != nil {
		return ResultsData{}, err
	}

	windef := ResultsData{Infected: false, Engine: getWinDefVersion()}

	log.WithFields(log.Fields{
		"plugin":   strings.Replace(name, "-", "_", -1),
		"category": category,
		"path":     path,
	}).Debug("Windows Defender Output: ", windefout)

	lines := strings.Split(windefout, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Scanning input") {
			continue
		}
		if strings.Contains(line, "EngineScanCallback") {
			threat := strings.TrimPrefix(strings.TrimSpace(line), "EngineScanCallback():")
			if len(threat) > 0 {
				windef.Infected = true
				threat = strings.TrimSpace(threat)
				threat = strings.TrimPrefix(threat, "Threat")
				threat = strings.TrimSuffix(threat, "identified.")
				windef.Result = strings.TrimSpace(threat)
			} else {
				log.Errorf("Umm... len(threat)=%d, threat=%v", len(threat), threat)
			}
		}
	}
	windef.Updated = getUpdatedDate()
	return windef, nil
}

func getWinDefVersion() string {

	versionOut, err := utils.RunCommand(nil, "/usr/bin/exiftool", "/loadlibrary/engine/mpengine.dll")
	assert(err)

	log.Debug("Windows Defender Version: ", versionOut)
	for _, line := range strings.Split(versionOut, "\n") {
		if len(line) != 0 {
			if strings.Contains(line, "Product Version Number          :") {
				return strings.TrimSpace(strings.TrimPrefix(line, "Product Version Number          :"))
			}
		}
	}
	return "error"
}

func getUpdatedDate() string {
	if _, err := os.Stat("/opt/malice/UPDATED"); os.IsNotExist(err) {
		return BuildTime
	}
	updated, err := ioutil.ReadFile("/opt/malice/UPDATED")
	assert(err)
	return string(updated)
}

func parseUpdatedDate(date string) string {
	layout := "200601021504"
	t, _ := time.Parse(layout, date)
	return fmt.Sprintf("%d%02d%02d", t.Year(), t.Month(), t.Day())
}

func updateAV(ctx context.Context) error {
	fmt.Println("Updating Windows Defender...")
	fmt.Println(utils.RunCommand(ctx, "/opt/malice/update"))
	// Update UPDATED file
	t := time.Now().Format("20060102")
	err := ioutil.WriteFile("/opt/malice/UPDATED", []byte(t), 0644)
	return err
}

func generateMarkDownTable(w WindowsDefender) string {
	var tplOut bytes.Buffer

	t := template.Must(template.New("windef").Parse(tpl))

	err := t.Execute(&tplOut, w)
	if err != nil {
		log.Println("executing template:", err)
	}

	return tplOut.String()
}

func printMarkDownTable(windef WindowsDefender) {

	fmt.Println("#### Windows Defender")
	table := clitable.New([]string{"Infected", "Result", "Engine", "Updated"})
	table.AddRow(map[string]interface{}{
		"Infected": windef.Results.Infected,
		"Result":   windef.Results.Result,
		"Engine":   windef.Results.Engine,
		"Updated":  windef.Results.Updated,
	})
	table.Markdown = true
	table.Print()
}

func printStatus(resp gorequest.Response, body string, errs []error) {
	fmt.Println(resp.Status)
}

func webService() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/scan", webAvScan).Methods("POST")
	log.Info("web service listening on port :3993")
	log.Fatal(http.ListenAndServe(":3993", router))
}

func webAvScan(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(32 << 20)
	file, header, err := r.FormFile("malware")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Please supply a valid file to scan.")
		log.Error(err)
	}
	defer file.Close()

	log.Debug("Uploaded fileName: ", header.Filename)

	tmpfile, err := ioutil.TempFile("/malware", "web_")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	data, _ := ioutil.ReadAll(file)

	if _, err = tmpfile.Write(data); err != nil {
		log.Fatal(err)
	}
	if err = tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

	// Do AV scan
	path = tmpfile.Name()
	windef := AvScan(60)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(windef); err != nil {
		log.Fatal(err)
	}
}

func main() {

	cli.AppHelpTemplate = utils.AppHelpTemplate
	app := cli.NewApp()

	app.Name = name
	app.Author = "blacktop"
	app.Email = "https://github.com/blacktop"
	app.Version = Version + ", BuildTime: " + BuildTime
	app.Compiled, _ = time.Parse("20060102", BuildTime)
	app.Usage = "Malice Windows Defender AntiVirus Plugin"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose, V",
			Usage: "verbose output",
		},
		cli.BoolFlag{
			Name:  "table, t",
			Usage: "output as Markdown table",
		},
		cli.BoolFlag{
			Name:   "callback, c",
			Usage:  "POST results to Malice webhook",
			EnvVar: "MALICE_ENDPOINT",
		},
		cli.BoolFlag{
			Name:   "proxy, x",
			Usage:  "proxy settings for Malice webhook endpoint",
			EnvVar: "MALICE_PROXY",
		},
		cli.StringFlag{
			Name:        "elasticsearch",
			Value:       "",
			Usage:       "elasticsearch url for Malice to store results",
			EnvVar:      "MALICE_ELASTICSEARCH_URL",
			Destination: &es.URL,
		},
		cli.IntFlag{
			Name:   "timeout",
			Value:  60,
			Usage:  "malice plugin timeout (in seconds)",
			EnvVar: "MALICE_TIMEOUT",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "update",
			Aliases: []string{"u"},
			Usage:   "Update virus definitions",
			Action: func(c *cli.Context) error {
				return updateAV(nil)
			},
		},
		{
			Name:  "web",
			Usage: "Create a Windows Defender scan web service",
			Action: func(c *cli.Context) error {
				webService()
				return nil
			},
		},
	}
	app.Action = func(c *cli.Context) error {

		var err error

		if c.Bool("verbose") {
			log.SetLevel(log.DebugLevel)
		}

		if c.Args().Present() {
			path, err = filepath.Abs(c.Args().First())
			assert(err)

			if _, err = os.Stat(path); os.IsNotExist(err) {
				assert(err)
			}

			windef := AvScan(c.Int("timeout"))
			windef.Results.MarkDown = generateMarkDownTable(windef)

			// upsert into Database
			if len(c.String("elasticsearch")) > 0 {
				err := es.Init()
				if err != nil {
					return errors.Wrap(err, "failed to initalize elasticsearch")
				}
				err = es.StorePluginResults(database.PluginResults{
					ID:       utils.Getopt("MALICE_SCANID", utils.GetSHA256(path)),
					Name:     strings.Replace(name, "-", "_", -1),
					Category: category,
					Data:     structs.Map(windef.Results),
				})
				if err != nil {
					return errors.Wrapf(err, "failed to index malice/%s results", name)
				}
			}

			if c.Bool("table") {
				fmt.Printf(windef.Results.MarkDown)
			} else {
				windef.Results.MarkDown = ""
				windefJSON, err := json.Marshal(windef)
				assert(err)
				if c.Bool("post") {
					request := gorequest.New()
					if c.Bool("proxy") {
						request = gorequest.New().Proxy(os.Getenv("MALICE_PROXY"))
					}
					request.Post(os.Getenv("MALICE_ENDPOINT")).
						Set("X-Malice-ID", utils.Getopt("MALICE_SCANID", utils.GetSHA256(path))).
						Send(string(windefJSON)).
						End(printStatus)

					return nil
				}
				fmt.Println(string(windefJSON))
			}
		} else {
			log.Fatal(fmt.Errorf("Please supply a file to scan with malice/windef"))
		}
		return nil
	}

	err := app.Run(os.Args)
	assert(err)
}
