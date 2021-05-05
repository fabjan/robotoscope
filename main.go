package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/fabjan/robotoscope/router"
)

type robotMap struct {
	lock sync.RWMutex
	data map[string]int
}

func newRobotMap() robotMap {
	return robotMap{
		data: make(map[string]int),
	}
}

func (m *robotMap) inc(name string) int {
	m.lock.Lock()
	defer m.lock.Unlock()
	count := m.data[name]
	count++
	m.data[name] = count
	return count
}

var robots = newRobotMap()
var cheaters = newRobotMap()

func count(r *http.Request, m *robotMap) {
	ua := r.Header.Get("User-Agent")
	if ua != "" {
		m.inc(ua)
	}
}

func collectRobot(w http.ResponseWriter, r *http.Request) {
	count(r, &robots)
	fmt.Fprintln(w, "User-agent: *")
	fmt.Fprintln(w, "Disallow: /secret/")
}

func reportCheater(w http.ResponseWriter, r *http.Request) {
	count(r, &cheaters)
	w.WriteHeader(http.StatusPaymentRequired)
}

func list(w http.ResponseWriter, m *robotMap) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	var b strings.Builder
	for robot, count := range m.data {
		fmt.Fprintf(&b, "%3v: %q\n", count, robot)
	}
	fmt.Fprint(w, b.String())
}

func showRobots(w http.ResponseWriter, r *http.Request) {
	list(w, &robots)
}

func showCheaters(w http.ResponseWriter, r *http.Request) {
	list(w, &cheaters)
}

func main() {
	addr := ":5000"
	if os.Getenv("PORT") != "" {
		addr = ":" + os.Getenv("PORT")
	}

	var r router.RegexpRouter

	r.HandleFunc(regexp.MustCompile("/robots.txt"), collectRobot)
	r.HandleFunc(regexp.MustCompile("/secret/*"), reportCheater)
	r.HandleFunc(regexp.MustCompile("/list.txt"), showRobots)
	r.HandleFunc(regexp.MustCompile("/cheaters.txt"), showCheaters)
	http.Handle("/", &r)

	log.Fatal(http.ListenAndServe(addr, nil))
}
