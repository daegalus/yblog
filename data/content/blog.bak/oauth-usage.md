---
author: Yulian Kuncheff
date: 2012-04-04T19:00:00Z
draft: false
slug: oauth-usage
title: OAuth Usage
type: blog
tags:
  - oauth
  - programming
  - selectivism
---

![oauth2](/images/2014/Jul/oauth2.png#c)

One thing that really bugs me lately is the selective use of OAuth providers. Specifically, Facebook. It is one of my pet peeves when I come to a site that I want to try out, but they require Facebook in order to login. I rarely if ever use Facebook, I personally hate using Facebook (rant for another day), but I have one, for some basic connections and logging into places if I really have to.

This normally wouldn't be an issue if developers would add more sources. At least the more popular ones like Twitter, Google, Github, Windows Live, whatever. And the way OAuth is made, its not very hard to do this. You just need to add a bit of services specific code to handle the keys, if that. There are even libraries for most languages and platforms that do this all for you.

So I really don't get why developers don't do this. 95% of the usage I have seen is to pull your name, email, picture, and maybe your bio. This can be done from all of these services and its not hard to do. Even with my hatred for Facebook, and that I prefer Bitbucket over Github. I would add support for these OAuth providers if I ever write any login-based web apps. I would much prefer to login with Twitter for everything, but that choice is never given to me, its always Facebook or Normal Login, or just Facebook. Rarely do I see more than that.

But this also goes for others too. For example Geekli.st. I know they are trying to emulate Twitter more than anything, but nothing is stopping them from using the other OAuth providers to provide a login to their Geekli.st persona. The same info that they pull from Twitter can be pulled from other services, and then once internal, they can work like they always do.

A couple of examples of places that could easily implement more, but fail to do so, some of these are sites I frequent and love, some are not. These are off the top of my head, so I can't name many, but I know I have stumbled on many, throughout surfing the web.

* The Verge - www.theverge.com - Only uses Facebook or Normal Login
* Geekli.st - www.geekli.st - Only Twitter
* Yahoo - www.yahoo.com - Only Windows Live & Google or Normal Login
* Pinterest - www.pinterest.com - Facebook and Twitter, (was excited when I saw twitter) but they could still add Google and maybe Live.
* Minus - www.minus.com - Facebook, Twitter, and Normal Login. Could add Google and Live.
* Flixster - www.flixster.com - Only Facebook and Normal Login
* Many others, I tend to not remember/go to sites that force me to use Facebook Only.

Numerous apps too:

* Draw Something - Facebook or Normal Login
* Zynga Games - Facebook, but this one I can understand due to how tied Zynga is to Facebook still.
* Plenty others.

I can only think of this one off the top of my head for providing lots of option, I know there are more, just can't remember them all:

* Disqus - Facebook, Twitter, Google, OpenID, Normal Login - Lots of choice, bound to hit something a user prefers.

On a side note, I do wish more services provided OAuth Logins in general. I prefer to just hit Login with Twitter, and since I already stay logged into Twitter all the time, I just click Accept for new apps, or it will just do a couple redirects and I am in.

I just feel developers need to not lock themselves into once service. One because no one knows how long a service will exist or be ok with a service's policies. Second, it lowers the amount of people that use your service. There are many that hate one service or the other, or hate Social Networks in general, but the more options you provide, the bigger the market you will capture. I know numerous websites that I would love to use, but don't want to login with Facebook, so I just never use them and/or find alternatives.
