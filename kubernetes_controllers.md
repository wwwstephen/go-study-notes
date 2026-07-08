# Kubebuilder Patterns

Kubebuilder is built around a small set of recurring patterns. If you understand these, you'll recognize almost every Kubernetes operator written with Kubebuilder. The most important concepts are:

- Controller Pattern
- Reconciliation Loop
- Informer/Watch Pattern
- Work Queue Pattern
- Desired State Pattern
- Finalizer Pattern

## 1. The Reconciliation Pattern (the core of Kubebuilder)

Everything revolves around a `Reconcile()` function:

```
Current Cluster State
        ↓
      Compare
        ↓
   Desired State
        ↓
   Take Actions
        ↓
      Return
```

A reconciler is not event-driven in the traditional sense. Instead:

1. Something changes
2. Kubernetes tells the controller
3. The controller asks: "What should the world look like?"
4. It keeps making changes until reality matches the desired state

**Example** — suppose your CRD is:

```yaml
apiVersion: apps.example.com/v1
kind: Database
spec:
  replicas: 3
```

Your reconcile loop may observe:

```
Desired replicas = 3
Actual replicas  = 2
```

It creates one more Pod. Next reconciliation:

```
Desired = 3
Actual  = 3
```

Do nothing. This makes reconciliation **idempotent**.

## 2. Desired State Pattern

Users only describe `spec`. Controllers update `status`.

```yaml
spec:
  replicas: 5
```

Controller updates:

```yaml
status:
  readyReplicas: 5
  phase: Running
```

This separation is fundamental throughout Kubernetes.

## 3. Watch Pattern

Kubebuilder controllers watch resources:

```go
ctrl.NewControllerManagedBy(mgr).
    For(&myv1.Database{}).
    Complete(r)
```

This tells controller-runtime to watch `Database` objects.

You can also watch other resources:

```go
Owns(&appsv1.Deployment{})
```

meaning:

```
Database
    ↓
 creates
    ↓
Deployment
```

If the Deployment changes → enqueue Database → reconcile Database.

## 4. Owner Reference Pattern

Controllers usually set:

```yaml
metadata:
  ownerReferences: ...
```

This creates an ownership tree:

```
Database
   ├── Deployment
   ├── Service
   └── Secret
```

Benefits:

- Automatic garbage collection
- Clear ownership
- Watches work correctly

Kubebuilder provides `controllerutil.SetControllerReference()`.

## 5. Finalizer Pattern

Deleting isn't immediate.

```
kubectl delete database mydb
```

becomes `metadata.deletionTimestamp != nil`. The controller notices → cleanup → remove finalizer → resource disappears.

Typical cleanup:

- Delete cloud database
- Revoke credentials
- Remove DNS
- Backup data

## 6. Status Pattern

Controller should never write into `Spec`. Instead update:

```go
instance.Status.Phase = "Running"
```

then:

```go
r.Status().Update(ctx, instance)
```

**Status** is owned by the controller. **Spec** is owned by users.

## 7. Event Recording Pattern

Instead of only logging `Deployment created`, controllers often emit Kubernetes Events:

```
Normal  Created      Deployment
Warning FailedCreate ...
```

These appear in `kubectl describe`.

## How Kubebuilder Monitors CRDs

This is probably the most interesting part. Kubebuilder itself does not watch the Kubernetes API directly — it uses the controller-runtime library, which in turn builds on client-go.

The architecture looks like this:

```
API Server
     │
     ▼
 Watch API
     │
     ▼
Shared Informer
     │
     ▼
 Local Cache
     │
     ▼
Event Handler
     │
     ▼
 Work Queue
     │
     ▼
 Reconciler
```

This is the same architecture used by many Kubernetes controllers.

### Step 1. Watch

When you write:

```go
ctrl.NewControllerManagedBy(mgr).
    For(&myv1.Database{}).
    Complete(r)
```

Kubebuilder registers a watch. Internally it becomes something conceptually like:

```
WATCH /apis/example.com/v1/databases
```

The API server sends a stream of events: `ADDED`, `MODIFIED`, `DELETED`. The controller receives them continuously.

### Step 2. Shared Informer

Instead of every controller calling the API repeatedly (`GET`, `GET`, `GET`, `GET`, ...), controller-runtime starts a shared informer.

The informer:

- Watches the API server once
- Maintains a local cache
- Broadcasts changes to interested controllers

Think of it as:

```
API Server
      │
      ▼
Informer Cache
      │
 ┌────┼────┐
 │    │    │
 A    B    C
```

One watch can feed many controllers.

### Step 3. Local Cache

The reconciler typically reads from the cache, not directly from the API server. When you write:

```go
r.Get(ctx, req.NamespacedName, obj)
```

the client usually retrieves the object from the informer cache.

Advantages:

- Much faster
- Fewer API requests
- Less API server load

The cache is kept up to date by the watch stream.

### Step 4. Event Handler

When an object changes (`Database updated`), the event handler enqueues a reconcile request:

```go
Request{
    Namespace: "default",
    Name:      "db1",
}
```

It does not invoke `Reconcile()` directly. Instead, it pushes a key onto a queue.

### Step 5. Work Queue

The queue stores keys like:

```
default/db1
default/db2
default/db3
```

Workers pull one item at a time:

```
Worker 1 → Reconcile(db1) → Done → Next
```

The work queue provides:

- Retries on errors
- Rate limiting
- Deduplication (multiple rapid updates often collapse into a single queued key)
- Configurable concurrency

### Step 6. Reconcile

The controller receives:

```go
func (r *DatabaseReconciler) Reconcile(
    ctx context.Context,
    req ctrl.Request,
)
```

It then:

```
Get object
    ↓
Inspect spec
    ↓
Inspect current resources
    ↓
Compare
    ↓
Update cluster
    ↓
Update status
    ↓
Return
```

The reconcile function is intentionally written to tolerate repeated execution.

## Watches on Other Resources

Kubebuilder can also watch resources related to your CRD:

```go
ctrl.NewControllerManagedBy(mgr).
    For(&myv1.Database{}).
    Owns(&appsv1.Deployment{}).
    Owns(&corev1.Service{}).
    Complete(r)
```

This means:

```
Database changed  → Reconcile Database

Deployment changed → Find owning Database → Reconcile Database

Service changed    → Find owning Database → Reconcile Database
```

The owner reference set on managed objects lets the controller map those changes back to the correct `Database` resource.

## Putting It All Together

A complete Kubebuilder controller follows this flow:

```
User creates a CRD instance
          │
          ▼
API Server stores the object
          │
          ▼
Watch receives an ADDED event
          │
          ▼
Shared Informer updates the local cache
          │
          ▼
Event Handler enqueues the object's namespace/name
          │
          ▼
Work Queue dispatches the request to a worker
          │
          ▼
Reconcile() reads the object from the cache
          │
          ▼
Controller compares desired state (spec) with actual state
          │
          ▼
Creates, updates, or deletes dependent resources
          │
          ▼
Updates the CRD's status
          │
          ▼
Returns; future changes trigger reconciliation again
```

These patterns are the foundation of Kubebuilder and Kubernetes operators in general. Once you're comfortable with reconciliation, shared informers, caching, watches, work queues, owner references, finalizers, and the separation of spec and status, you'll be able to understand the structure of most operator implementations built with Kubebuilder.
