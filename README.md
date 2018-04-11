## Golang Pod Counter

A Simple Golang app to count the number of pods in a Kubernetes cluster.




## Dev Environment

* Docker for Mac; Edge channel.
* Kubertenes version 1.9.6.
* Golang version 1.9.4




## Quickstart Guide

Run the following command to deploy the app to a cluster

`kubectl create -f https://raw.githubusercontent.com/murali44/PodCounter/master/src/deploy.yml --namespace=podcount`

Get the IP address and Port to access the service

`kubectl describe services podcount --namespace=podcount`

Here's the output on my dev machine

```Name:                     podcount
Namespace:                podcount
Labels:                   run=podcounter
Annotations:              <none>
Selector:                 run=podcounter
Type:                     NodePort
IP:                       10.111.89.246
LoadBalancer Ingress:     localhost
Port:                     <unset>  8080/TCP
TargetPort:               8080/TCP
NodePort:                 <unset>  30813/TCP
Endpoints:                10.1.0.47:8080,10.1.0.48:8080
Session Affinity:         None
External Traffic Policy:  Cluster
Events:                   <none>
```

Access the app using the following service values

`curl <LoadBalancer Ingress>:<NodePort>`

From the example above, the app URL is

`curl localhost:30813`




## Solution Design

**Assumptions**: To keep things simple, I'm not creating a load balanced service. We can deploy 
multiple replicas of the app. However, to access it, we'll need to hit one of the pods. Depending 
upon how your cluster is configured, you can 

The basic idea for the solution is to leverage the Kubernetes APIServer to get the number of pods.
We know that every pod in the cluster is automatically injected with a service account 
(Client creds, Token, Certificate, etc), which can be used to authenticate against the APIServer 
on the master node. Kubernetes provides a client library which does the heavylifting of establishing 
a secure connection to the APIServer. See here for more details: [kubernetes/client-go](https://github.com/kubernetes/client-go)


![PodCounterArchitecture](https://github.com/murali44/PodCounter/blob/master/PodCounter.jpg)

To restrict the pod count to the current namespace, I had to inject the namespace into the container 
as an environment variable. I did this in the deployment configuration file (deploy.yml)




## How to build and Package the app

* Make sure you are in the source directory.
    `cd PodCounter/src`
* Get all dependencies.
    `go get`
* Build the app binary.
    `GOOS=linux go build -o ./app .`
* Build the docker image.
    `docker build -t murali44/podcounter .`
* Push the image to docker hub.
    `docker push murali44/podcounter`

**Note**: Don't forget to tag the image with your own docker hub repo and update the deploy.yml file to use your own container image.