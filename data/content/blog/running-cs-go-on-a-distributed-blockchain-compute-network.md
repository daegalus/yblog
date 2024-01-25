---
author: Yulian Kuncheff
date: 2017-09-08T16:38:54Z
draft: false
slug: running-cs-go-on-a-distributed-blockchain-compute-network
title: Running CS:GO on a Distributed Blockchain Compute Network
type: blog
tags:
  - programming
  - blockchain
  - somn
  - fog-computing
  - docker
---

Many that conversed with me in the cryptocurrency world know that I am a huge fan and supporter of SONM. After reading their whitepaper, I quickly saw that they had a very adoptable design to their tech. Now, its not purely blockchain. It uses the blockchain to distributed tasks, send messaging, and other things like that, the actual tasks get run on standard compute devices using Docker.

So to put it long story short, SONM is a distributed docker runner with blockchain to handle distribution. This is a fantastic middle ground setup, and poises SONM as being able to run anything that Docker can run, quickly beating out competitors like Golem (which can only do image rendering), iExec, and Elastic.

When I got the opportunity to setup a CS:GO server for the Stratis community and as I have the week off between my job change, I decided it would be a fantastic project to not only get a CS:GO server in Docker, but also get it working on the SONM platform. This sounded like a fun staycation project.

Now, because SONM runs docker images, I first had to get a docker image of the CS:GO Server. I could have started from scratch, but being the lazy developer I am, I looked at what I could find that was already working. I quickly found [this](https://hub.docker.com/r/austinsaintaubin/docker-csgoserver/) and [this](https://hub.docker.com/r/johnjelinek/csgoserver/). These docker images were a great starting point, but neither has been updated in almost 2 years. Firstly they were running on Ubuntu 14.04, which while still in LTS, will very soon not be. Second, the [script/tool](https://gameservermanagers.com/lgsm/csgoserver/) they are using internally to setup the server has been updated significantly over the 2 years, and no longer works right inside the containers.

Next order of business, was where to run it. I normally run most of my things in DigitalOcean. As a DevOps engineer for a living, I spend a LOT of time in AWS land, and I honestly don't want to fall into the trap of AWS. So I spun up a $5 Droplet, and started working.

I setup SONM first, seeing as this was the newest tech, I figured I would see if it even works before investing time working towards it. Plus besides previous mini-tests, this is the first time I setup SONM for actual work.

<img src="/images/${slug}/somn-setup" alt="A picture of the terminal with the setup of the SOMN miner" transform-images="avif jxl webp png" />

As you can see from the picture above, it all went pretty smoothly. My droplet has an external IP and I didn't have to change any of the default config settings. I just ran the hub and miner, and it all worked out. That was absolutely painless. I used the precompiled Linux binaries for this and it all worked swimmingly.

Next was getting a workign Docker image. This is what took 90% of the 16 hours I spent on this. Figuring out why the my variables weren't being injected properly (LGSM moved the config settings to a different file), figuring out how to get the docker container from not terminating at the end (csgoserver start runs in the background and just returns, ending the docker session, had to use csgoserver debug to get it runnning in the foreground).

But the biggest problem I had initially, was hard disk space. GS:GO Server is 15 gigs of data, it pretty much downloads the full game. Seeing as how a $5 droplet only has 20 gigs of space, and a lot of it is quickly taken up, this started throwing errors from the SteamCMD downloader getting the game.

So I started over on a new $5 droplet, but in SFO2 region so I could add additional volumes. For another $10, I added a 100gig drive. And with a little googling, I found a few instructions on which folder I needed to move over and how to setup the symlinks.

Finally, after tons of tinkering with the image, I finally succeeded, and actually published the image to dockerhub. [All 16 gigs of glory](https://hub.docker.com/r/daegalus/docker-csgo/)(Well 9gigs compressed). I originally wasn't going to push it to dockerhub, but SONM requires the image be on a registry for the task to work, and I didn't want to setup my own local registry. Plus I felt others might want to try this out too. The script will update the game on launch, so no need to worry about it becoming out of date either.

And once I got the image up and running, I was able to log into the game and play agains the bots. Did a little more tinkering with the image to allow for the user to pass in a GSLT token to so the server didn't start in LAN Only mode.

<img src="/images/${slug}/csgo-server-running" alt="Logs from he CS GO docker server running." transform-images="avif jxl webp png" />

After all that, you would think that getting SONM to run it would be even easier now that I have an image it can wrong. But thats where my assumptions fell. Because I needed to pass in that GSLT token to run a basic server, I needed to find a way to pass it in. After much searching through the SONM code, I found that there is no way to actually pass in environment variables from the task.yaml file that is used by the hub and miner. But I did find a test that shows you can pass in a `command` parameter with the startup command.

With that in mind, I modified the start script to take an optional parameter. Starting the docker server with `docker run -i -t daegalus/docker-csgo:latest /home/csgoserver/start.sh <token>` actually worked and let me pass in the token directly. Without any further thought, I quickly added to the parameter to my task.yaml and the appropriate command. And sent the task to the miner.

Nothing, it started in LAN only mode. Turns out, the `command` parameter is either no longer used, or never got hooked up, as it doesn't do anything.

<img src="/images/${slug}/csgo-lanonly-somn-task" alt="Logs from SOMN task running CS GO in LanOnly mode." transform-images="avif jxl webp png" />

The picture above shows a working CS:GO server in LAN Only mode, running in a Docker container, launched as a task sent to the SONM hub, and is being run by the SONM Miner. It works, I was impressed.

Now, I wouldn't knock SONM for not having these parameters. This task.yaml and hub/miner setup wasn't even in the original Alpha that was put out during the ICO, this is Alpha 2. Also, its alpha, they might have not gotten around to hooking everything up. Just needs a bit more time to get all the features and config parameters in. There are a lot of `TODO` in the code still, which is fine for something in early Alpha. There are bugs, needs for more configuration options, but what they have is solid.I look forward to future alpha releases.

I have opened a [Github Issue](https://github.com/sonm-io/core/issues/114) for this, and hopefully it will keep it in mind when they work on adding features/fixes in future alphas and betas.

I do have to say, I am very impressed with SONM and what they have accomplished so far. Getting a SONM stack up was very simple and easy to do, and if this is only Alpha, I can only imagine how great it can be as it gets closer to finished product. They have blown their competitors out of the water.

I leave you with this final image of the SONM Docker server showing up and joinable from the CS:GO server browser. I did launch the server directly, so the docker container is up and running, and people can play on the server. The IP in the images will work, and if the server is up (should be at the time of this writing), you can hop in and play.

<img src="/images/${slug}/csgo-server-browser" alt="Ingame CSGO server browser showing our browser." transform-images="avif jxl webp png" />

*PS: I know this isn't a very technical blog post, but I wanted to just talk about the experience of working with all this.*

*NANO: `nano_1syg4hkhe7b4mkuy33mbaabwbzete3fo1zmih31ukeizkti77bp9tgm5iupe`*
*ETH/ERC20: `0xD0a54F1614F6373f55E68c6C34CeD127aff8b05E`*
