package cache

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
)

// This file provides test-only access to engine internals through the
// clusterCache facade. The lifecycle methods live on *legacyEngine; rather
// than thread the engine type through every white-box test that drives a
// single lifecycle step, we forward here. Production code routes through
// c.engine (EnsureSynced, Invalidate) or the engine calls its own methods.

// legacyEngineOf returns the active legacy engine, panicking if the cache is
// not running one (a test setup bug).
func legacyEngineOf(c *clusterCache) *legacyEngine {
	return c.engine.(*legacyEngine)
}

func (c *clusterCache) sync() error {
	return c.engine.sync()
}

func (c *clusterCache) startMissingWatches() error {
	return c.engine.startMissingWatches()
}

// stopWatching is not on the syncEngine interface (its namespace semantics are
// engine-specific — see the note in sync_engine.go), so dispatch on the
// concrete type here.
func (c *clusterCache) stopWatching(gk schema.GroupKind, ns string) {
	legacyEngineOf(c).stopWatching(gk, ns)
}

func (c *clusterCache) recordEvent(event watch.EventType, un *unstructured.Unstructured) {
	legacyEngineOf(c).recordEvent(event, un)
}
