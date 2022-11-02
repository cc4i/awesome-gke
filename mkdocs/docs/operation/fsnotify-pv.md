#

## Description

Sometimes monitoring the status of a file is important for an application, such as triggering a hot reloading, related actions, etc. You can get a notification through 'Inotify' in Linux when operating a file, it also works for Pod mounted with Cloud Filestore or persistem storage. Following demo tried to show you how it works. 

However it's not working for object filesystem, which is not operating through the traditional filesystem API.


## Guide
```sh
git clone https://github.com/cc4i/multi-k8s.git
cd multi-k8s/asset/inotify

# 1. Create a storage class & PVC, which takes a while to provision (minutes)
kubectl apply -f manifests/pv.yaml

# 2. Run Pod
kubectl apply -f manifests/golang-deploy.yaml

# 3. Copy the testing program
kubectl get pods
# Example output:
# NAME                            READY   STATUS    RESTARTS   AGE
# golang-1-6f65477fc4-n4q98       1/1     Running   0          23h
# golang-2-5598d8996f-jhnvx       1/1     Running   0          23h

# Replace the pod name
kubectl cp ../inotify golang-1-6f65477fc4-n4q98:/go/inotify
kubectl exec -it pods/golang-1-6f65477fc4-n4q98 -- bash
cd /go/inotify
go mod tidy
go run main.go /tmp1

# 4. Open other termial, and validating Inotify. You should be able to see operations after type 'touch cc'
kubectl exec -it pods/golang-2-5598d8996f-jhnvx -- bash
cd /tmp1
touch cc

```
## Limitation

> Inotify reports only events that a user-space program triggers through the filesystem API.  As a result, it does not catch remote events that occur on network filesystems.  (Applications must fall back to polling the filesystem to catch such events.) Furthermore, various pseudo-filesystems such as /proc, /sys, and /dev/pts are not monitorable with inotify.

## References

- https://github.com/fsnotify/fsnotify
- https://github.com/ofek/csi-gcs
- https://man7.org/linux/man-pages/man7/inotify.7.html#:~:text=When%20a%20directory%20is%20monitored,referring%20to%20the%20inotify%20instance.