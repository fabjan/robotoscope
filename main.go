package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/fabjan/robotoscope/core"
	"github.com/fabjan/robotoscope/database"
	"github.com/fabjan/robotoscope/html"
	"github.com/fabjan/robotoscope/router"
)

// RobotStore can track and list robots
type RobotStore interface {
	Count(name string) error
	List() ([]core.RobotInfo, error)
}

type closable interface {
	Close() error
}

type server struct {
	robots   RobotStore
	cheaters RobotStore
	closer   closable
}

func count(r *http.Request, s RobotStore) {
	ua := r.Header.Get("User-Agent")
	if ua != "" {
		err := s.Count(ua)
		if err != nil {
			log.Printf("ERROR: store error when counting (%v)", err)
		}
	}
}

func (s server) collectRobot(w http.ResponseWriter, r *http.Request) {
	count(r, s.robots)
	fmt.Fprintln(w, "User-agent: *")
	fmt.Fprintln(w, "Disallow: /secret/")
}

func (s server) reportCheater(w http.ResponseWriter, r *http.Request) {
	count(r, s.cheaters)
	w.WriteHeader(http.StatusPaymentRequired)
}

func list(w http.ResponseWriter, m RobotStore) {
	var b strings.Builder
	infos, err := m.List()
	if err != nil {
		log.Printf("ERROR: store error when listing (%v)", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for _, info := range infos {
		fmt.Fprintf(&b, "%3v: %q\n", info.Seen, info.UserAgent)
	}
	fmt.Fprint(w, b.String())
}

func (s server) showRobots(w http.ResponseWriter, r *http.Request) {
	list(w, s.robots)
}

func (s server) showCheaters(w http.ResponseWriter, r *http.Request) {
	list(w, s.cheaters)
}

func (s server) showIndex(w http.ResponseWriter, r *http.Request) {
	rInfos, err := s.robots.List()
	if err != nil {
		log.Printf("ERROR: store error when listing robots (%v)", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	cInfos, err := s.cheaters.List()
	if err != nil {
		log.Printf("ERROR: store error when listing cheaters (%v)", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := html.Page{
		Title:    "Robotoscope",
		Robots:   rInfos,
		Cheaters: cInfos,
	}

	var b strings.Builder
	err = html.Render(&b, data)
	if err != nil {
		log.Printf("ERROR: render error when listing (%v)", err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write([]byte(b.String()))
	}
}

func (s *server) connectPostgres(rawURL string) {
	db, err := database.OpenPg(rawURL)
	if err != nil {
		log.Fatalf("cannot open database connection: %v\n", err)
	}

	s.closer = db

	robots, err := database.NewPgStore(db, "robots")
	if err != nil {
		log.Fatalf("cannot initialize robots store: %v\n", err)
	}

	cheaters, err := database.NewPgStore(db, "cheaters")
	if err != nil {
		log.Fatalf("cannot initialize robots store: %v\n", err)
	}

	s.robots = robots
	s.cheaters = cheaters
}

func (s *server) connectRedis(rawURL string) {
	c := database.OpenRedis(rawURL)
	s.closer = c
	s.robots = database.NewRedisStore(c, "robots")
	s.cheaters = database.NewRedisStore(c, "cheaters")
}

func (s *server) useInMemoryStores() {
	r := database.NewRobotMap()
	c := database.NewRobotMap()
	s.robots = &r
	s.cheaters = &c
}

func main() {

	s := server{}

	redisURL := os.Getenv("REDIS_URL")
	dbURL := os.Getenv("DATABASE_URL")

	if redisURL != "" {
		log.Println("INFO: REDIS_URL is set, connecting to Redis")
		s.connectRedis(redisURL)
	} else if dbURL != "" {
		log.Println("INFO: DATABASE_URL is set, connecting to Postgres")
		s.connectPostgres(dbURL)
	} else {
		log.Println("WARN: using in-memory store, set REDIS_URL or DATABASE_URL in the environment for persistence")
		s.useInMemoryStores()
	}

	defer func() {
		if s.closer != nil {
			s.closer.Close()
		}
	}()

	addr := ":5000"
	if os.Getenv("PORT") != "" {
		addr = ":" + os.Getenv("PORT")
	}

	var r router.RegexpRouter

	r.HandleFunc(regexp.MustCompile("/robots.txt"), s.collectRobot)
	r.HandleFunc(regexp.MustCompile("/secret/*"), s.reportCheater)
	r.HandleFunc(regexp.MustCompile("/list.txt"), s.showRobots)
	r.HandleFunc(regexp.MustCompile("/cheaters.txt"), s.showCheaters)
	r.HandleFunc(regexp.MustCompile("/"), s.showIndex)
	http.Handle("/", &r)

	log.Printf("INFO: ready, serving at %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
