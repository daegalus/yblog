---
author: Yulian Kuncheff
date: 2022-08-22T23:12:00Z
draft: false
slug: tils-1
title: TILs 8/22/22
type: blog
tags:
  - til
  - gke
  - gcp
  - glb
  - neg
  - nginx
  - iap
  - headers
  - basicauth
  - jenkins
---
So TILs didn't work out for me, I prefer longer form writing. But I might gather up TILs like this and add them in groups to a post like this.

## TIL 2 *2021-02-16*

**Update 8/24/2022**: This is no longer needed, there is now official programmatic way to do this using a header GCP providers for IAP using the `Proxy-Authorization` header.
[Official IAP Documentation](https://cloud.google.com/iap/docs/authentication-howto#authenticating_from_proxy-authorization_header)
The original post is at the more link.

<!--more-->
If you need to use basic auth when behind IAP, you can work around it by setting a custom header, and then remapping it after you get past IAP. Specifically when running NGINX Ingress on GKE.

My example is for allowing Basic Auth to Jenkins while running behind IAP on GKE. I created a custom header called `X-Jenkins-Authorization`. (This can be anything, and for any service, not just jenkins).

Some official examples:
[https://github.com/kubernetes/ingress-nginx/tree/master/docs/examples/customization/custom-headers](https://github.com/kubernetes/ingress-nginx/tree/master/docs/examples/customization/custom-headers)

Kubernetes Nginx Ingress configmap changes

```yaml
...
data:
    proxy-set-headers: `somens/nginx-proxy-headers`
...
```

Proxy Headers config map:

```yaml
apiVersion: v1
data:
  Authorization: ${http_x_jenkins_authorization}
kind: ConfigMap
metadata:
  name: nginx-proxy-headers
  namespace: somens
```

Apply those accordingly, and now you can basic auth through IAP, you just need to set your basic auth in the appropriate header.

## TIL 1 *2020-05-16*

I had third party software that relied on NGINX ingress and wouldn't work with GKE Ingress. After lots of digging around and piecing together some info, I found I can attach expose an NEG directly to the ingress controller and route into the cluster that way.

Just append the following to your annotations.

```yaml
cloud.google.com/neg: '{"exposed_ports": {"80":{}, "443":{}}}'
```

If you want to use HTTPS load balancer, only 80 and 443 will be usable. if you add more ports, more things will be accessible through the NEGs it creates. 1 NEG for each.
