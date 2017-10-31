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

	"github.com/ideo/dragonfruit/backends/backend_couchdb"
	"github.com/martini-contrib/gzip"
	"github.com/martini-contrib/sessions"
)

const (
	VERSION                = "0.5.5"
	PATH_TO_DEFAULT_CONFIG = "/usr/local/etc/dragonfruit.conf"
)

type cnf struct {
	config         dragonfruit.Conf
	addInteractive bool
	serve          bool
	version        bool
	resourcetype   string
	resourcefile   string
}

/* main executable function */
func main() {
	fmt.Println("\n\n\033[31m~~~~~Dragon\033[32mFruit~~~~~\033[0m")
	cnf := parseFlags()
	fmt.Print("\033[1mVersion ", VERSION, " \033[0m\n\n\n")
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

	m.RunOnAddr(cnf.Host + ":" + cnf.Port)
	return m
}

func returnSuccess(res http.ResponseWriter) (int, string) {
	h := res.Header()

	h.Add("Content-Type", "text/plain")
	//res.Write([]byte("Content-Type: text/plain"))
	return 200, "auth"
}

/* addResourceFromFile adds a new resource directly from a file, with
standard naming conventions for resources and API endpoints */
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

/* addResource handles interactive mode parsing and naming of resources */
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

/* parseFlags parses the command-line flags passed to the CLI */
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

	/* If not using the default config file, path to a config file */
	var conflocation = flag.String("conf", "", "Path to a config file.")

	flag.Parse()
	baseConf, err := ioutil.ReadFile(PATH_TO_DEFAULT_CONFIG)
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
