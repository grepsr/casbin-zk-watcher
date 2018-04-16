package zkwatcher

import (
	"strconv"
	"strings"
	"time"

	"github.com/casbin/casbin/persist"
	"github.com/casbin/casbin/util"
	"github.com/samuel/go-zookeeper/zk"
)

// Watcher represents a zk watcher.
type Watcher struct {
	hosts    string
	path     string
	conn     *zk.Conn
	running  bool
	callback func(string)
}

// finalizer is the destructor for Watcher.
func finalizer(w *Watcher) {
	w.running = false
}

// NewWatcher is the constructor for Watcher.
// hosts is the comma-separated URLs to zookeeper hosts.
// path is the path which will be watched. If not provided, it defaults to "/casbin".
func NewWatcher(hosts string, path ...string) persist.Watcher {
	w := &Watcher{}
	w.hosts = hosts
	w.running = true
	w.callback = nil

	if len(path) == 1 {
		w.path = path[0]
	} else {
		w.path = "/casbin"
	}

	w.createConnection()

	go w.startWatch()

	return w
}

// createConnection creates a new connection to Zookeeper.
func (w *Watcher) createConnection() error {
	hostSlice := strings.Split(w.hosts, ",")
	c, _, err := zk.Connect(hostSlice, time.Second)
	if err != nil {
		return err
	}
	w.conn = c
	return nil
}

// SetUpdateCallback sets the callback function which will be called
// by the watcher when the policy is changed in the DB by other instances.
func (w *Watcher) SetUpdateCallback(callback func(string)) error {
	w.callback = callback
	return nil
}

// Update calls the update callback of other instances to synchronize
// their policy.
func (w *Watcher) Update() error {
	rev := 0

	data, stat, err := w.conn.Get(w.path)
	if err != nil {
		return err
	}
	rev, err = strconv.Atoi(string(data))
	if err != nil {
		return err
	}
	util.LogPrint("Get revision: ", rev)

	rev++
	newRev := strconv.Itoa(rev)

	util.LogPrint("Set revision: ", newRev)
	_, err = w.conn.Set(w.path, []byte(newRev), stat.Version)
	return err

}

// startWatch is a goroutine that watches for policy changes.
func (w *Watcher) startWatch() error {
	if !w.running {
		return nil
	}

	data, errors := w.watchPath()

	for {
		select {
		case d := <-data:
			if w.callback != nil {
				w.callback(d)
			}
		case err := <-errors:
			return err
		}
	}
}

// watchPath is a function which continuously watches
// a given path for changes.
func (w *Watcher) watchPath() (chan string, chan error) {
	data := make(chan string)
	errors := make(chan error)

	go func() {
		for {
			_, _, events, err := w.conn.GetW(w.path)
			if err != nil {
				errors <- err
				return
			}

			evt := <-events

			if evt.Err != nil {
				errors <- evt.Err
				return
			}

			d, _, err := w.conn.Get(w.path)
			if err != nil {
				errors <- err
			}
			data <- string(d)

		}
	}()

	return data, errors
}
