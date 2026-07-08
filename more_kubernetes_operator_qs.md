# Kubernetes Operator Interview Q&A

## 1. Explain how Kubernetes controllers work internally

**Weak answer:**

> "Controllers watch resources and run reconciliation."

**Stronger answer:**

> "Controllers use informers backed by client-go. The informer maintains a local cache populated by a watch stream from the API server. Events are converted into reconcile requests and placed into a rate-limited work queue. The reconciler reads desired state from the CR spec and current state from the cache, then performs idempotent actions to converge the system."

Then add:

> "The important design property is eventual consistency. Controllers don't react to events directly; events are only hints that something may have changed. The reconciler always evaluates the current state."

That last sentence will stand out.

## 2. Why does Kubernetes use reconciliation instead of imperative commands?

**Bad:**

> "Because Kubernetes likes declarative YAML."

**Better:**

> "Imperative systems encode a sequence of actions, which makes recovery difficult. Declarative reconciliation means the controller continuously compares desired state with observed state and can recover from partial failures, restarts, and drift."

**Example:**

> "If a node dies and three Pods disappear, the Deployment controller doesn't need to remember what happened. It just sees desired replicas=3 and actual replicas=1 and creates replacements."

## 3. How would you design a cluster lifecycle management system?

He may ask:

> "How would you manage thousands of Kubernetes clusters?"

Think like a platform architect. Answer structure:

```
API layer
    │
Desired state store
    │
Controller/reconciliation engine
    │
Provisioning workers
    │
Cloud APIs
```

Mention:

- Asynchronous workflows
- Idempotency
- Retries
- State machines
- Events
- Audit trails

**Example answer:**

> "I would avoid a synchronous API that provisions clusters in one request. Cluster creation is a long-running distributed workflow. I would represent lifecycle state as a resource, persist workflow state, and use controllers that reconcile cluster intent against reality."

That is exactly the Kubernetes mindset.

## 4. Why would you use PostgreSQL + Pub/Sub instead of etcd?

Very relevant.

**A good answer:**

> "etcd is excellent for strongly consistent Kubernetes API state, but it isn't designed as a general workflow database at massive scale. For fleet management, I might separate concerns: PostgreSQL for durable workflow state and relational queries, Pub/Sub for event distribution, and Kubernetes controllers for reconciliation."

Mention:

- Don't put everything in the Kubernetes API
- Separate control plane state from workload state

## 5. What problems happen when managing many clusters?

**Control plane scaling problems:**

- API server pressure
- etcd growth
- Watch explosion
- Reconciliation storms

**Solutions:**

- Sharding controllers
- Caching
- Rate limiting
- Batching
- Hierarchical controllers

**Example:**

> "At fleet scale, the challenge isn't just creating clusters. It is maintaining thousands of independent control loops without creating correlated failures."

That is a very senior observation.

## 6. Explain Kubernetes networking

Know:

```
Pod
 │
 ▼
CNI plugin
 │
 ▼
Linux networking
 │
 ▼
Node network
```

Understand:

- CNI
- kube-proxy
- Services
- ClusterIP
- Ingress
- Network policies

**A strong statement:**

> "Kubernetes does not implement networking itself. It defines the networking model, and CNI implementations provide the actual dataplane."

## 7. What happens when you run `kubectl apply`?

Excellent interview question.

```
kubectl
   │
  HTTPS
   │
API server
   │
authentication
   │
authorization
   │
admission controllers
   │
  etcd
   │
watch events
   │
controllers reconcile
```

Mention:

- API server is the gateway
- etcd stores desired state
- Controllers create reality

## 8. How would you debug a Kubernetes production issue?

**Good framework — first, ask "what changed?"**

Check:

```
kubectl get events
kubectl describe
kubectl logs
```

**Then work through the layers:**

```
Application
    │
   Pod
    │
Container runtime
    │
  Node
    │
Network
    │
Control plane
```

Avoid random debugging. Say:

> "I start from symptoms and narrow the failure domain rather than immediately changing configuration."

## 9. What is GitOps?

Don't answer:

> "Using Git for Kubernetes YAML."

**Better:**

> "GitOps treats Git as the desired state source, with automated reconciliation between Git state and cluster state. The important part is not Git itself, but having an auditable desired state and automated convergence."

Mention:

- ArgoCD
- Drift detection
- Rollback
- Auditability

## 10. What is your approach to reliability?

Reliability is designed, not accidental:

- Automation
- Observability
- Failure testing
- Safe deployments
- Operational simplicity

**Good phrase:**

> "A system is not reliable because failures don't happen. It is reliable because failures are expected and recovery is automated."

---

## Questions YOU can ask him (these will impress him)

Senior engineers judge candidates partly by their questions.

1. **"At your scale, what were the biggest bottlenecks: API server throughput, etcd, reconciliation storms, or operational workflows?"**
   This shows you understand real problems.

2. **"How did you decide what belongs inside Kubernetes controllers versus external workflow systems?"**
   Very relevant.

3. **"When managing many clusters, where did you draw the boundary between cluster-local control loops and fleet-level orchestration?"**
   Excellent platform engineering question.

4. **"How do you handle failure domains when a fleet management system itself becomes a dependency?"**
   Very senior question.

5. **"What parts of Kubernetes abstractions stop working well at hyperscale?"**
   This invites a deep discussion.
