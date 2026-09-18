package cache

// syncEngine encapsulates the mode-specific cluster cache lifecycle: how
// resources are initially listed, how watches are started for newly discovered
// APIs, and how per-invalidate runtime state is dropped.
//
// Everything shared — the resource index (store.resources / store.nsIndex /
// store.parentUIDToChildren), the cross-GK index maintenance, the read paths
// (FindResources, IterateHierarchyV2, GetManagedLiveObjs), CRD handling, and
// the discovery primitives (processApi, listResources, checkPermission) —
// lives on *store. An engine holds a *store and reaches shared state through
// it; it never references the outer clusterCache facade.
//
// legacyEngine is currently the only implementation. The interface exists to
// make the shared/lifecycle boundary explicit and enforced by the compiler, so
// an alternative lifecycle (e.g. an informer-based one, see issue #19199) can
// be added without another pass over the shared layer.
type syncEngine interface {
	// sync performs a full (re)synchronization of the cluster cache: fetch
	// cluster metadata, discover APIs, and start watching every monitored
	// resource. On return the cache is populated and watches are running.
	//
	// Locking: called with store.lock held (EnsureSynced), and the legacy
	// engine holds it for the entire call. An implementation that needs to
	// release store.lock mid-flight would make store state mutable across a
	// sync() call, so every call site would have to tolerate that; none does
	// today.
	sync() error
	// startMissingWatches discovers the current API surface and starts a
	// watch for every GroupKind not already being watched. Invoked on the
	// CRD add/change path (handleCRDEvent) so newly registered custom
	// resources begin syncing without a full re-sync. Caller holds store.lock.
	startMissingWatches() error
	// onInvalidate drops the engine's per-invalidate runtime state. Called by
	// clusterCache.Invalidate under store.lock so the facade never has to know
	// which engine-private fields exist (legacy retires its batched-event
	// channel).
	onInvalidate()
}

// Note: stopWatching (tearing down the watch set for a GroupKind) is
// deliberately NOT part of the interface. Its semantics are engine-specific —
// the legacy engine tears down a single (GroupKind, namespace) watch, while an
// informer-based engine would have to purge every namespace for the GroupKind
// because its informers share one watch context — and every production caller
// is the owning engine itself, so exposing it polymorphically would only
// invite callers to rely on whichever semantics they happened to test against.
