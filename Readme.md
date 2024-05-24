
## Admission Webhook. 


## References
https://github.com/douglasmakey/admissioncontroller

https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/

## How to get started?

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


