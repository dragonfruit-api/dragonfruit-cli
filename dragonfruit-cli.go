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

const (
	VERSION = "0.5.2"
)

type cnf struct {
	config         dragonfruit.Conf
	addInteractive bool
	serve          bool
	version        bool
	resourcetype   string
	resourcefile   string
}

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
	fmt.Println("\n\n\033[31m~~~~~Dragon\033[32mFruit~~~~~\033[0m")
	cnf := parseFlags()
	fmt.Print("\033[1mVersion", VERSION, " \033[0m\n\n\n")
	if cnf.version {
		return
	}

	// dragonfruit setup
	d := backend_couchdb.Db_backend_couch{}
	d.Connect("http://" + cnf.config.DbServer + ":" + cnf.config.DbPort)

	if (cnf.resourcefile != "" && cnf.resourcetype == "") ||
		(cnf.resourcefile == "" && cnf.resourcetype != "") {
		fmt.Println("\033[31;1mYou must enter both a resource file and a resource type if you pass a resource from the command line.\033[0m")
		return
	}

	if cnf.addInteractive {
		addResouce(&d, cnf.config)
	} else if cnf.resourcefile != "" {
		addResourceFromFile(&d, cnf.config, cnf.resourcetype, cnf.resourcefile)
	}

	if cnf.serve {
		launchServer(&d, cnf.config)
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

	dragonfruit.ServeDocSet(m, d, cnf)

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

/* addResourceFromFile adds a new resource directly from a file, with the
standard naming conventions */
func addResourceFromFile(d dragonfruit.Db_backend,
	cnf dragonfruit.Conf,
	resourceType string,
	fname string) {

	byt, err := ioutil.ReadFile(fname)
	if err != nil {
		panic(err)
	}

	err = dragonfruit.RegisterType(d, byt, cnf, resourceType, "")
	if err != nil {
		panic(err)
	}

}

func addResouce(d dragonfruit.Db_backend, cnf dragonfruit.Conf) {

	// start a scanner to read from the command line
	scanner := bufio.NewScanner(os.Stdin)

	// set the resource type name
	fmt.Print("\033[1mWhat is the base model that this API returns?\033[0m ")
	scanner.Scan()

	resourceType := inflector.Singularize(scanner.Text())
	path := inflector.Pluralize(resourceType)

	fmt.Print("\033[1mWhat is the path for APIs for this model?\033[0m press [enter] for \"", path, "\" ")
	scanner.Scan()
	tmpPath := scanner.Text()
	if tmpPath != "" {
		path = tmpPath
	}

	fmt.Print("\033[1mEnter a path to a sample data file:\033[0m ")
	scanner.Scan()
	fname := scanner.Text()
	byt, err := ioutil.ReadFile(fname)
	if err != nil {
		fmt.Println(err)
	}

	err = dragonfruit.RegisterType(d, byt, cnf, resourceType, path)

	if err != nil {
		panic(err)
	}

	fmt.Println("Done!")

}

func parseFlags() cnf {

	// set up a config object
	dfcnf := dragonfruit.Conf{}

	/* should we start a server? */
	var serve = flag.Bool("serve", true, "Start a server after running")

	/* should we start a server? */
	var version = flag.Bool("version", false, "Display the version and terminate.")

	/* should we try to parse a resource? */
	var addresource = flag.Bool("add", false, "Add a new resource (interactive mode).")

	/* should we try to parse a specific file? */
	var resourcefile = flag.String("file", "", "Load and parse a resource file (with standard naming).")

	/* If we do parse a file, what is the resource type for that file? */
	var resourcetype = flag.String("type", "", "The resource type for the file.")

	var conflocation = flag.String("conf", "", "Path to a config file.")

	flag.Parse()
	baseConf, err := ioutil.ReadFile("/usr/local/etc/dragonfruit.conf")
	if err != nil {
		panic("base conf file missing")
	}

	err = json.Unmarshal(baseConf, &dfcnf)

	if err != nil {
		panic(err)
	}

	if *conflocation != "" {
		out, err := ioutil.ReadFile(*conflocation)
		if err != nil {
			panic("cannot find file " + *conflocation)
		}
		err = json.Unmarshal(out, &dfcnf)
		if err != nil {
			panic(err)
		}
	}

	outconfig := cnf{
		config:         dfcnf,
		addInteractive: *addresource,
		serve:          *serve,
		version:        *version,
		resourcetype:   *resourcetype,
		resourcefile:   *resourcefile,
	}

	return outconfig
}
