# Common Bugs with Kubernetes Operators

The main failure areas tend to cluster around reconciliation and
infrastructure provisioning: storage failures where PVCs fail to bind or
volumes fail to mount, scheduling failures when requested CPU/memory can't
be satisfied, serialization issues between the CLI and operator, resource
creation failures, secret/configuration issues, and reconciliation bugs
where desired state isn't correctly applied or recovered after a partial
failure.

## Storage / Volume Issues

- **Storage class does not exist**
  → Create the missing StorageClass or update the PVC to reference a valid one.

- **No available PV capacity**
  → Increase cluster storage capacity or provision more persistent volumes.

- **Wrong access mode (ReadWriteOnce vs ReadWriteMany)**
  → Change the PVC access mode or use a storage backend that supports the required mode.

- **Cloud provider volume provisioning failure**
  → Check cloud provider logs, IAM permissions, quotas, and retry volume provisioning.

- **Volume attached to the wrong node**
  → Verify node affinity and let Kubernetes detach/reattach the volume correctly.

- **Incorrect mount path**
  → Fix the volume mount path in the Deployment configuration and redeploy.

- **PVC stuck in `Terminating`**
  → Usually a leftover finalizer (e.g. `kubernetes.io/pvc-protection` waiting
  on a pod, or a CSI finalizer waiting on the driver). Check for pods still
  referencing the PVC, then confirm the CSI driver is healthy before
  removing finalizers manually.

## Kubernetes Scheduling Issues

- **Insufficient CPU/memory for pod scheduling**
  → Add more cluster capacity or reduce the pod resource requests.

- **Pod stuck in `Pending` state**
  → Check pod events (`kubectl describe pod`) to identify scheduling constraints.

- **Taints/tolerations or node affinity mismatch**
  → Verify the pod's tolerations and affinity rules match the intended nodes; `kubectl describe pod` shows unmatched taints in the scheduling events.

- **Image pull failures (`ImagePullBackOff`)**
  → Check the image tag/registry, `imagePullSecrets`, and registry rate limits or auth expiry.

## Operator Issues

- **Operator does not notice changes**
  → Fix the watch/reconciliation logic and ensure the CRD events trigger reconciliation.

- **Operator updates the wrong resource**
  → Fix resource selectors, labels, ownership references, or reconciliation logic.

- **Infinite reconciliation loop**
  → Ensure reconciliation is idempotent and only updates resources when changes are needed (e.g. avoid re-writing a status field with a new timestamp every reconcile, which re-triggers the watch).

- **Fails to recover after partial failure**
  → Add retry logic, proper error handling, and make reconciliation able to rebuild missing resources.

- **Stale CRD schema / conversion webhook mismatch**
  → When bumping a CRD version, ensure the conversion webhook (or `none` strategy) correctly translates old stored objects; otherwise existing resources fail to decode.

- **RBAC denies the operator's service account**
  → Check the operator's `ClusterRole`/`Role` covers every verb and resource it touches (including subresources like `status` and `finalizers`); `kubectl auth can-i --as=system:serviceaccount:<ns>:<sa> <verb> <resource>` is the fastest way to confirm.

- **Leader election flapping (multiple replicas fighting for control)**
  → Check lease renewal timing/network latency between replicas; a slow API server or short lease duration can cause repeated leader changes and duplicate reconciles.

- **Finalizer stuck, resource won't delete**
  → The controller crashed or errored before removing its finalizer. Check operator logs for the failed cleanup step, fix it, and let reconciliation retry — avoid manually stripping finalizers unless the underlying resource is confirmed gone.

- **Webhook (validating/mutating) unavailable blocks all applies**
  → If the webhook service/pod is down, the API server rejects matching requests. Check `failurePolicy` (fail-open vs fail-closed) and the webhook pod's health.

## WordPress Application Issues

- **Database connection failure**
  → Verify database availability, networking, credentials, and connection settings.

- **Wrong environment variables**
  → Fix the ConfigMap/Secret values and restart the affected pods.

- **Missing secrets**
  → Ensure the operator creates and mounts required Kubernetes Secrets.

- **Incorrect PHP configuration**
  → Update PHP settings via ConfigMaps or container configuration.

- **Plugin incompatibility**
  → Disable the problematic plugin or pin compatible plugin/WordPress versions.

## Database Issues

- **Database pod fails**
  → Check container logs, resource limits, storage, and database configuration.

- **Credentials mismatch**
  → Rotate/update credentials and ensure WordPress references the correct Secret.

- **Database storage unavailable**
  → Fix PVC/storage issues and restore database connectivity.

- **Database migrations fail**
  → Review migration errors, fix schema conflicts, and retry migrations.

- **Connection limits exceeded**
  → Increase database connection limits or add connection pooling.

## What an Operator Team Gets Paged For

These are the incidents that tend to actually wake someone up — usually
because they're cluster-wide, self-inflicted by the operator, or actively
getting worse the longer they run.

- **Operator crashlooping**
  → Symptom: `CrashLoopBackOff` on the operator pod itself, so nothing it
  manages gets reconciled. Check `kubectl logs --previous` for the panic/error,
  and roll back to the last known-good operator image if it's a bad deploy.

- **Reconcile storm hammering the API server**
  → A bug causes every reconcile to requeue immediately (e.g. status update
  triggers its own watch event), spiking API server load and starving other
  controllers. Symptom: API server latency/5xx spikes correlating with the
  operator's request rate. Mitigate by scaling the operator to 0 or adding
  rate limiting, then fix the requeue logic.

- **Operator managing thousands of CRs falls behind (queue depth growing)**
  → Reconcile latency or workqueue depth metrics climb unbounded. Often a
  slow external dependency (cloud API, database) blocking the reconcile
  loop. Check for missing timeouts/context cancellation and whether
  concurrency (`MaxConcurrentReconciles`) is tuned too low.

- **Mass resource deletion from a bad reconcile or owner-reference bug**
  → A logic error causes the operator to delete resources it shouldn't
  (e.g. garbage-collecting on a stale view of desired state, or an incorrect
  `ownerReference` causing cascading deletes). This is the highest-severity
  page — first action is to scale the operator down to stop the bleeding,
  then restore from backups/GitOps state.

- **CRD deletion hangs, blocking `kubectl delete namespace` cluster-wide**
  → Custom resources with finalizers can't be cleaned up because the owning
  operator is down or erroring, which blocks namespace deletion entirely.
  Bring the operator back up first; only strip finalizers manually as a last
  resort once you've confirmed no orphaned cloud resources will result.

- **Leader election lost cluster-wide control (split-brain)**
  → Two operator replicas both believe they're leader (e.g. after an etcd/API
  server blip) and reconcile concurrently, causing conflicting writes.
  Symptom: resource flapping between two different states. Check lease
  objects (`kubectl get lease -n <ns>`) and restart the non-leader replica.

- **Webhook outage blocks all deploys/applies cluster-wide**
  → If the operator's admission webhook pod is unreachable and
  `failurePolicy: Fail`, *every* matching resource in the cluster fails to
  apply, not just the operator's own resources. This pages fast because it
  looks like "the whole cluster is broken." Check webhook pod health and
  consider `failurePolicy: Ignore` for non-critical webhooks as a safety net.

- **Secret/cert rotation breaks every managed workload at once**
  → A rotated credential or webhook TLS cert isn't picked up correctly,
  and every workload the operator manages starts failing auth
  simultaneously. Check cert-manager/secret rotation logs and whether the
  operator (and its webhook server) reloaded the new cert rather than
  serving a cached/expired one.

- **Cloud provider API throttling cascades into reconcile failures**
  → The operator hits cloud API rate limits (e.g. provisioning volumes or
  load balancers) and every reconcile for every resource starts failing,
  which often triggers *more* retries and makes throttling worse. Add
  exponential backoff and check for a retry storm before assuming the cloud
  provider is down.

- **Out-of-memory operator pod (OOMKilled)**
  → Usually an unbounded cache (e.g. watching all resources in a large
  cluster without field/label selectors) or a memory leak in a long-lived
  reconcile loop. Check memory usage trend over time, not just at the OOM
  event, to distinguish a leak from a genuine capacity issue.

## Useful Debugging Commands

```sh
kubectl describe pod <pod>              # scheduling/mount/pull events
kubectl get events --sort-by=.lastTimestamp -n <ns>
kubectl logs <operator-pod> -n <ns>     # reconcile errors
kubectl get pvc,pv -n <ns>              # storage binding state
kubectl auth can-i <verb> <resource> --as=system:serviceaccount:<ns>:<sa>
```
