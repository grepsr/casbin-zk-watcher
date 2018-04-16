package zkwatcher

import (
	"testing"
	"time"

	"github.com/casbin/casbin"
	"github.com/casbin/casbin/util"
)

func updateCallback(rev string) {
	util.LogPrint("New revision detected:", rev)
}

func TestWatcher(t *testing.T) {
	// updater represents the Casbin enforcer instance that changes the policy in DB.
	updater := NewWatcher("localhost:2181", "/casbin")

	// listener represents any other Casbin enforcer instance that watches the change of policy in DB.
	listener := NewWatcher("localhost:2181", "/casbin")

	// listener should set a callback that gets called when policy changes.
	listener.SetUpdateCallback(updateCallback)

	// updater changes the policy, and sends the notifications.
	err := updater.Update()
	if err != nil {
		panic(err)
	}

	// Now the listener's callback updateCallback() should be called,
	// because it receives the notification of policy update.
	// You should see "[New revision detected: X]" in the log.

	// Add delay so that the callbacks get called before the program exits.
	time.Sleep(time.Second * 1)
}

func TestWithEnforcer(t *testing.T) {
	// Initialize the watcher.
	w := NewWatcher("localhost:2181")

	// Initialize the enforcer.
	e := casbin.NewEnforcer("examples/rbac_model.conf", "examples/rbac_policy.csv")

	// Set the watcher for the enforcer.
	e.SetWatcher(w)

	// By default, the watcher's callback is automatically set to the
	// enforcer's LoadPolicy() in the SetWatcher() call.
	// We can change it by explicitly setting a callback.
	w.SetUpdateCallback(updateCallback)

	// Update the policy to test the effect.
	// You should see "[New revision detected: X]" in the log.
	e.SavePolicy()

	// Add delay so that the callbacks get called before the program exits.
	time.Sleep(time.Second * 1)
}
