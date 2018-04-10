## Golang Pod Counter

A Simple Golang app to count the number of pods in a Kubernetes cluster.


## Dev Environment

* Docker for Mac; Edge channel.
* Kubertenes version 1.9.6.
* Golang version 1.9.4


## How to deploy

Run the following command

`kubectl create -f https://raw.githubusercontent.com/murali44/PodCounter/master/src/deploy.yml --namespace=podcount`


## Solution Design

The basic idea for the solution is to leverage the Kubernetes APIServer to get the number of pods.
We know that every pod in the cluster is automatically injected with a service account (Client creds, Token, Certificate, etc), which can be used to authenticate against the APIServer on the master node. 

![PodCounterArchitecture](https://github.com/murali44/PodCounter/blob/master/PodCounter.jpg)

I found a Golang Kubernetes client library which does the heavylifting of establishing a secure connection to the APIServer. See here for more details: [kubernetes/client-go](https://github.com/kubernetes/client-go)

## How to build and Package the app

* Make sure you are in the source directory.
    `cd PodCounter/src`
* Get all dependencies.
    `go get`
* Build the app.
    `GOOS=linux go build -o ./app .`
* Build the docker image.
    `docker build -t murali44/podcounter .`
* Push the image to docker hub.
    `docker push murali44/podcounter`

Note: Don't forget to update the deploy.yml file to use your own container image.