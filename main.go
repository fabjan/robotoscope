package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/fabjan/robotoscope/html"
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

func showIndex(w http.ResponseWriter, r *http.Request) {
	data := html.Page{
		Title:    "Robotoscope",
		Robots:   []html.RobotInfo{},
		Cheaters: []html.RobotInfo{},
	}

	robots.lock.RLock()
	defer robots.lock.RUnlock()
	for robot, count := range robots.data {
		info := html.RobotInfo{
			Seen:      count,
			UserAgent: robot,
		}
		data.Robots = append(data.Robots, info)
	}

	cheaters.lock.RLock()
	defer cheaters.lock.RUnlock()
	for robot, count := range cheaters.data {
		info := html.RobotInfo{
			Seen:      count,
			UserAgent: robot,
		}
		data.Cheaters = append(data.Cheaters, info)
	}

	var b strings.Builder
	err := html.Render(&b, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "cannot render page")
	} else {
		w.Write([]byte(b.String()))
	}
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
	r.HandleFunc(regexp.MustCompile("/"), showIndex)
	http.Handle("/", &r)

	log.Fatal(http.ListenAndServe(addr, nil))
}
