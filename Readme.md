
## Admission Webhook. 


## References
https://github.com/douglasmakey/admissioncontroller

https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/

## What is going to happen?

Below is all scripted out in install.sh but the overview is

* Create certificates.

* Import tls secrets into admission

* install webhooks

* Deploy test pod

### Couple Todo's:

* Setup this install script to pull from an env var
* Setup mutate and verify to be better examples

## Why?

For the primary reason I like to understand how something underneath works. Having a good understanding how something works behind the currents providers the oportunity for better troubleshooting skills. 

For the second reason, I have a another project that this became handy. This repo was created a while ago to have the understanding and utilize webadmission hooks for it.

### What is your other projects?

In a very high detailed overview I am using an webadmission webhook to do authentication with another project and injects secrets into the pod. Could it have been done differently? Yes, I decided to go this route. 


### Full documentation on what is happening?

Coming soon: [will update when I have the blog post finished.]


## How to get started?

* Spin up a kubernetes cluster. I recomend K3d.
* Build the docker file, which will compile the code as well.
    ```
    docker build . --tag  [your docker repo]
    ```
* Update deployment/deployment.yaml image to you rrepo.
* Run deployment.install.sh
* Deploy test pod manifest.
    ```
    kubectl -n test apply -f deployment/pod-test.yaml
    ```

** The code will fail on when the name is **testname** and   **admission-webhook** label is set enable.

** The code will mutate the pod and set an env var value when name is **pod1**


### Can this code be used to make cool stuff?

Yes, This is a good very basic base to get you started.