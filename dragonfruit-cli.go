package main

import (
	"encoding/json"
	"fmt"

	"github.com/go-martini/martini"
	"github.com/ideo/dragonfruit"

	"bufio"
	"flag"
	"github.com/gedex/inflector"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ideo/dragonfruit/backends/backend_couchdb"
	"github.com/martini-contrib/gzip"
	"github.com/martini-contrib/oauth2"
	"github.com/martini-contrib/sessions"
	goauth2 "golang.org/x/oauth2"
)

type loginresult struct {
	User       gplususer `json:"user"`
	IsLoggedIn bool      `json:"isLoggedIn"`
}

type gplususer struct {
	Emails []struct {
		Value string `json:"value"`
		Type  string `json:"type"`
	} `json:"emails"`
	DisplayName string `json:"displayName"`
	Image       struct {
		Url string `json:"url"`
	} `json:"image"`
	Expires time.Time `json:"expires"`
}

func main() {
	fmt.Println("\n\n\033[31m~~~~~Dragon\033[32mFruit~~~~~\033[0m\n\n")
	cnf, serve, add := parseFlags()

	// dragonfruit setup
	d := backend_couchdb.Db_backend_couch{}
	d.Connect("http://" + cnf.DbServer + ":" + cnf.DbPort)

	if add {
		addResouce(&d, cnf)
	}

	if serve {
		launchServer(&d, cnf)
	}

}

func launchServer(d dragonfruit.Db_backend, cnf dragonfruit.Conf) *martini.ClassicMartini {
	wd, _ := os.Getwd()
	st_opts := martini.StaticOptions{}
	st_opts.Prefix = wd
	m := dragonfruit.GetMartiniInstance(cnf)

	for _, dir := range cnf.StaticDirs {

		m.Use(martini.Static(dir))
	}

	dragonfruit.ServeDocSet(d, cnf)

	m.Use(sessions.Sessions("my_session", sessions.NewCookieStore([]byte("secret123"))))

	m.Use(gzip.All())

	oauthConf := &goauth2.Config{
		ClientID:     "288198830216-s4klktd4qm3asq72sm7acifugcumdseq.apps.googleusercontent.com",
		ClientSecret: "zzKEh5RPFhaEeJOsQi-D34Rc",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email"},
		RedirectURL: "http://" + cnf.Host + ":" + cnf.Port + "/oauth2callback",
	}
	m.Use(oauth2.Google(oauthConf))
	m.Post("/checkauth", oauth2.LoginRequired, func(tokens oauth2.Tokens, res http.ResponseWriter, s sessions.Session) (int, string) {
		h := res.Header()
		h.Add("Content-Type", "application/json; charset=utf-8")
		code := 200

		result := loginresult{}

		if tokens.Expired() {
			result.IsLoggedIn = false
			code = 403
		}
		test, user := checkAuth(oauthConf, tokens)

		if !test {
			result.IsLoggedIn = false
			code = 403
		} else {
			result.IsLoggedIn = true
			result.User = user
		}
		out, err := json.Marshal(result)
		if err != nil {
			code = 500
		}
		return code, string(out)
	})

	m.RunOnAddr(cnf.Host + ":" + cnf.Port)
	return m
}

func returnSuccess(res http.ResponseWriter) (int, string) {
	h := res.Header()

	h.Add("Content-Type", "text/plain")
	//res.Write([]byte("Content-Type: text/plain"))
	return 200, "auth"
}

func checkAuth(conf *goauth2.Config, tokens oauth2.Tokens) (bool, gplususer) {
	var user gplususer
	//fmt.Println(tokens)
	tok := &goauth2.Token{
		AccessToken:  tokens.Access(),
		Expiry:       tokens.ExpiryTime(),
		RefreshToken: tokens.Refresh(),
	}

	client := conf.Client(goauth2.NoContext, tok)
	resp, err := client.Get("https://www.googleapis.com/plus/v1/people/me")
	defer resp.Body.Close()

	if err != nil {
		return false, user
	}

	out, err := ioutil.ReadAll(resp.Body)

	json.Unmarshal(out, &user)

	for _, v := range user.Emails {
		if strings.Contains(v.Value, "ideo.com") {
			user.Expires = tokens.ExpiryTime()
			return true, user
		}
	}

	return false, user
}

/* shamelessly cut and paste from the dragonfruit api builder */
func addResouce(d dragonfruit.Db_backend, cnf dragonfruit.Conf) {
	rd, err := dragonfruit.LoadDescriptionFromDb(d, cnf)
	res := cnf.ResourceTemplate

	scanner := bufio.NewScanner(os.Stdin)

	res.BasePath = "/api"

	fmt.Print("\033[1mEnter a base path for all APIs\033[0m (press [enter] for \"" + res.BasePath + "\"):")
	scanner.Scan()

	tmpBasepath := scanner.Text()
	if tmpBasepath != "" {
		res.BasePath = tmpBasepath
	}

	fmt.Print("\033[1mEnter the resource type for this API:\033[0m ")

	// resource type
	scanner.Scan()
	resourceType := inflector.Singularize(scanner.Text())

	path := inflector.Pluralize(resourceType)
	if err != nil {
		panic(err)
	}

	defaultText := "Describes operations on " + path
	// description
	fmt.Print("\033[1mDescribe what the ", path, " service will do\033[0m (press [enter] for \""+defaultText+"\"):")
	scanner.Scan()
	resourceDescription := scanner.Text()
	if resourceDescription == "" {
		resourceDescription = defaultText
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	ok := true
	for _, v := range rd.APIs {
		if v.Path == "/"+path {
			ok = false
			break
		}
	}

	if ok {
		tmp := dragonfruit.ResourceSummary{
			Path:        "/" + path,
			Description: resourceDescription,
		}
		rd.APIs = append(rd.APIs, &tmp)
	}

	fmt.Print("\033[1mWhat is the base model that this API returns?\033[0m (press [enter] for \"", resourceType, "\"):")
	res.ResourcePath = "/" + path
	scanner.Scan()
	modelType := scanner.Text()
	if modelType == "" {
		modelType = resourceType
	}

	fmt.Print("\033[1mEnter a path to some sample data:\033[0m ")
	scanner.Scan()
	fname := scanner.Text()
	byt, err := ioutil.ReadFile(fname)
	if err != nil {
		fmt.Println(err)
	}

	modelMap, err := dragonfruit.Decompose(byt, modelType, cnf)

	res.Models = modelMap
	upstreamParams := make([]*dragonfruit.Property, 0)
	apis := dragonfruit.MakeCommonAPIs("",
		path,
		strings.Title(resourceType),
		modelMap,
		upstreamParams,
		cnf)

	res.Apis = append(res.Apis, apis...)

	rd.Save(d)
	res.Save(d)
	preperror := d.Prep(path, res)
	if preperror != nil {
		panic(preperror)
	}

	fmt.Println("Done!")

}

func parseFlags() (dragonfruit.Conf, bool, bool) {
	// set up a config object
	cnf := dragonfruit.Conf{}

	/* should we start a server? */
	var serve = flag.Bool("serve", true, "Start a server after running")
	/* should we try to parse a resource? */
	var addresource = flag.Bool("add", false, "Add a new resource")

	var conflocation = flag.String("conf", "/usr/local/etc/dragonfruit.conf", "Path to a config file.")

	flag.Parse()

	out, err := ioutil.ReadFile(*conflocation)
	if err != nil {
		panic("cannot find file " + *conflocation)
	}

	err = json.Unmarshal(out, &cnf)
	if err != nil {
		panic(err)
	}

	return cnf, *serve, *addresource
}
