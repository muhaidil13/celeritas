package celeritas

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/muhaidil13/celeritas/render"
)

const version = "1.0.0"

// celeritas type for my app
type Celeritas struct {
	AppName  string
	Debug    bool
	Version  string
	Routers  *chi.Mux
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	RootPath string
	config   config
	JetViews *jet.Set
	Render   *render.Render
}

// for myapp configuration
type config struct {
	port     string
	renderer string
}

func (c *Celeritas) New(rootPath string) error {

	// init for create path
	pathConfig := initPaths{
		rootPath:    rootPath,
		folderNames: []string{"handlers", "migrations", "views", "data", "public", "tmp", "logs", "middleware"},
	}
	err := c.Init(pathConfig)
	if err != nil {
		return err
	}

	// cek .env
	err = c.checkDotEnv(rootPath)
	if err != nil {
		return err
	}

	// read .env for setting file
	err = godotenv.Load(rootPath + "/.env")
	if err != nil {
		return err
	}

	// Create Loggers
	infolog, errorlog := c.startLoggers()
	c.InfoLog = infolog
	c.ErrorLog = errorlog
	// parsing string to bolean
	c.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))

	c.Version = version
	c.RootPath = rootPath
	c.Routers = c.routes().(*chi.Mux)
	// config
	c.config = config{
		port:     os.Getenv("PORT"),
		renderer: os.Getenv("RENDERER"),
	}

	var views = jet.NewSet(
		jet.NewOSFileSystemLoader(fmt.Sprintf("%s/views", rootPath)),
		jet.InDevelopmentMode(),
	)

	c.JetViews = views

	// choosing type template
	c.CreateRenderer()

	return nil
}

// create folder
func (c *Celeritas) Init(p initPaths) error {
	root := p.rootPath
	for _, path := range p.folderNames {
		err := c.CreateDirIfNotExists(root + "/" + path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Celeritas) ListenAndServe() {
	srv := http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		ErrorLog:     c.ErrorLog,
		Handler:      c.Routers,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	c.InfoLog.Println("Listening on Port %d ", os.Getenv("PORT"))
	err := srv.ListenAndServe()
	c.ErrorLog.Fatal(err)
}

// check .enf func
func (c *Celeritas) checkDotEnv(path string) error {
	err := c.CreateFileIfNotExists(fmt.Sprintf("%s/.env", path))
	if err != nil {
		return err
	}
	return nil
}

// startLoggers function
func (c *Celeritas) startLoggers() (*log.Logger, *log.Logger) {
	var infoLog *log.Logger
	var errorLog *log.Logger
	infoLog = log.New(os.Stdout, "INFO \t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stdout, "ERROR \t", log.Ldate|log.Ltime|log.Lshortfile)

	return infoLog, errorLog
}

func (c *Celeritas) CreateRenderer() {
	myRenderer := render.Render{
		Renderer: c.config.renderer,
		RootPath: c.RootPath,
		Port:     c.config.port,
		JetViews: c.JetViews,
	}
	c.Render = &myRenderer
}
