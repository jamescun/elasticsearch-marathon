package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"syscall"
)

func main() {
	marathonAddr := os.Getenv("MARATHON_ADDR")
	marathonAppId := os.Getenv("MARATHON_APP_ID")

	if marathonAddr == "" {
		log.Fatalln("fatal: marathon: $MARATHON_ADDR is required")
	} else if marathonAppId == "" {
		log.Fatalln("fatal: marathon: $MARATHON_APP_ID (automatic) is not set")
	}

	args := env2args(os.Environ(), "ELASTICSEARCH_", nil)

	appsUrl := marathonAddr + "/v2/apps?embed=apps.tasks&id="
	appsPath := path.Dir(marathonAppId)

	if appsPath == "/" {
		appsUrl += marathonAppId
	} else {
		appsUrl += appsPath
	}
	log.Printf("info: marathon: query url=%s\n", appsUrl)

	apps, err := GetMarathonApps(appsUrl)
	if err != nil {
		log.Fatalln("fatal: marathon:", err)
	}

	if len(apps) > 0 {
		hosts := []string{}

		for _, app := range apps {
			for _, task := range app.Tasks {
				if task.State == "TASK_RUNNING" && len(task.Ports) >= 2 {
					addr := task.Addr(1)
					hosts = append(hosts, addr)

					log.Printf("info: marathon: node discovered appId=%s taskId=%s addr=%s\n", app.Id, task.Id, addr)
				}
			}
		}

		if len(hosts) > 0 {
			args = append(args, "--discovery.zen.ping.unicast.hosts="+strings.Join(hosts, ","))
		}
	}

	args = append(os.Args[1:], args...)
	log.Println("info: exec: /docker-entrypoint.sh", args)

	err = syscall.Exec("/docker-entrypoint.sh", args, os.Environ())
	if err != nil {
		log.Fatalln("fatal: exec:", err)
	}
}

type HTTPError int

func (he HTTPError) Error() string {
	return fmt.Sprintf("HTTP Error %d", he)
}

type App struct {
	Id string `json:"id"`

	Tasks []Task `json:"tasks"`
}

type Task struct {
	Id    string `json:"id"`
	AppId string `json:"appId"`
	State string `json:"state"`

	Host  string `json:"host"`
	Ports []int  `json:"ports"`
}

func (t Task) Addr(port int) string {
	if port > (len(t.Ports) - 1) {
		panic("invalid ports index")
	}

	return net.JoinHostPort(t.Host, strconv.Itoa(t.Ports[port]))
}

func GetMarathonApps(url string) ([]App, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, HTTPError(res.StatusCode)
	}

	var apps struct {
		Apps []App `json:"apps"`
	}
	err = json.NewDecoder(res.Body).Decode(&apps)
	if err != nil {
		return nil, err
	}

	return apps.Apps, nil
}

// env2args converts a set of prefixed environment variables (in KEY=VALUE format) in to
// command line arguments, optionally transformed with a custom function.
// if the transform function is nil, the default GNU style transform will be applied.
func env2args(environ []string, prefix string, transformFn func(key string) string) []string {
	var res []string

	if transformFn == nil {
		transformFn = defaultTransformFn
	}

	for _, env := range environ {
		if strings.HasPrefix(env, prefix) {
			env = strings.TrimPrefix(env, prefix)

			i := strings.IndexFunc(env, isEqual)
			if i > -1 {
				k, v := env[:i], env[i+1:]
				k = transformFn(k)

				res = append(res, k+v)
			}
		}
	}

	return res
}

func defaultTransformFn(key string) string {
	key = strings.ToLower(key)
	key = strings.Replace(key, "_", ".", -1)

	return "--" + key + "="
}

func isEqual(r rune) bool {
	return r == '='
}
