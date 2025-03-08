---
author: Yulian Kuncheff
date: 2013-03-21T19:00:00Z
draft: false
slug: using-intellijwebstorm-to-debug-web-applications
title: Using IntelliJ/WebStorm to debug Dart Web Applications
type: blog
tags:
  - programming
  - dart
  - jetbrains
  - webstorm
  - intellij
  - web
---

@![Jetbrains IntelliJ/Webstorm debug console](/images/{{\(index .Posts 0\).FrontMatter.Slug}}/debug_console.png "Jetbrains IntelliJ/Webstorm debug console"){"avif,jxl,webp"}

Since I feel that the Dart Plugin in IntelliJ is at a point where it can do most, if not all things that the Dart Editor can do. I thought I would show how to setup debugging of Web Applications.
Primarily because a certain few on IRC wanted a step-by-step guide because they were too lazy to figure it out (you know who you are..... Don). Seeing how I have been Evangelizing IntelliJ/Webstorm on
IRC, I decided to just help out the community with this mini-guide. To save on typing IntelliJ/Webstorm everywhere, I am going to us IW as a short for refering to both.

I will be using Spectre as my test project, because I do not have any of my own personal projects that are Web Apps, I mostly do backend stuff. Plus this also shows that this can handle large complex
applications.

Caveats:

* Due to a bug in IW, launching the Dartium browser while running any other Chrome version, will just open a tab in that version. So you need to shutdown any open Chrome browsers beforehand. ([Bug #1](https://youtrack.jetbrains.com/issue/WEB-1561 "WEB-1561") & [Bug #2](https://youtrack.jetbrains.com/issue/WEB-6695 "WEB-6695"))

### Setting up the Project

Setting up the project is pretty straightforward. If you are starting a new project, just make a new "Dart Web Application" and you are mostly set.

If you already have a working project, and you just want to switch IDEs, I find this approach the quickest and least error prone. Hit new Project and select "Dart Web Application", and set the Project
and Module settings to the same root folder (makes life easier later). Make sure the Dart SDK path is right. If you have at any point added a Dart SDK anywhere, it should auto-fill for you.

@![Creating a dart web application in Jetbrains IntelliJ/Webstorm](/images/{{\(index .Posts 0\).FrontMatter.Slug}}/create_project.png "Creating a dart web application in Jetbrains IntelliJ/Webstorm"){"avif,jxl,webp"}

@![Jetbrains IntelliJ/Webstorm debug console](/images/{{\(index .Posts 0\).FrontMatter.Slug}}/debug_console.png "Jetbrains IntelliJ/Webstorm debug console"){"avif,jxl,webp"}

### Browser Setup

This is just to make Dart Debuggin easier, but this is global, so you are welcome to set this up to your preference in any other way. In my settings below, I override the default browser with Dartium,
and also the Chrome location. This will make sure no matter how it wants to open Chrome, it will be Dartium. This setting should be easily found in both IDE versions.

@![Configuring web browsers in Jetbrains IntelliJ/WebStorm](/images/{{\(index .Posts 0\).FrontMatter.Slug}}/browser_setup.png "Configuring web browsers in Jetbrains IntelliJ/WebStorm"){"avif,jxl,webp"}

### Debugging

At this point you should be ready to go. The final step is to debug the app. Find the primary HTML file in you folder structure, right-click it, and hit 'Debug "x.html"'

@![Menu entry for debugging a dart web app](/images/{{\(index .Posts 0\).FrontMatter.Slug}}/debug_menu.png "Menu entry for debugging a dart web app"){"avif,jxl,webp"}

It will open up, but it seems like nothing worked right (this happens only the first time). The reason being is that IntelliJ uses its own Plugin to support their debugging and hook into the browser. It will redirect you to the download page, and ask you to install it. You need to install this.

@![Installing the Jetbrains IDE extension in the browser](/images/{{\(index .Posts 0\).FrontMatter.Slug}}/jetbrains_extension.png "Installing the Jetbrains IDE extension in the browser"){"avif,jxl,webp"}

Once this is installed, it will automatically continue where it left off, load the page, and hook in the debugger.

@![Application running in the browser](/images/{{\(index .Posts 0\).FrontMatter.Slug}}/browser_running.png "Application running in the browser"){"avif,jxl,webp"}

And now if we switch back to the IDE we can see the debugger hooked and running.

@![Jetbrains IntelliJ/Webstorm debug console](/images/{{\(index .Posts 0\).FrontMatter.Slug}}/debug_console.png "Jetbrains IntelliJ/Webstorm debug console"){"avif,jxl,webp"}

You can see console output in the console tab. Debugger tab will have all your watches and object views when you hit a breakpoint, and the Scripts tab I assume is for SourceMaps and other stuff for debugging.

Breakpointing and stopping works fully as long as the debug connection is open.

### Debugging stuff like WebGL

(WebStorm Only as of writing. Explained below.)

So some cases require you to connect to a server address, even if its localhost. For example WebGL. You can't run WebGL stuff from the filesystem, needs to go through a server.

Well, WebStorm (and soon IntelliJ) has an awesome feature where it already is running a Web Server. And you can access it from any browser for any project you have.

***https://localhost:63342/[projectName]/path/to/file***

In spectre's case in my project, it would be: ***https://localhost:63342/spectre/web/asset_pack/asset_pack.html***

This will load the html file through the web server that is built in to WebStorm. So WebGL works.

But how do you make access this easy? and have the debugger turned on?

Well, we need to make a new custom configuration for the project, and use that. In order to create a custom configuration, there is a dropdown box in the toolbar at the top, right before the Run and Debug buttons. Open the dropdown and hit "Edit Configuration".

@![Edit Configuration dropdown in Jetbrains IntelliJ/WebStorm](/images/{{\(index .Posts 0\).FrontMatter.Slug}}/edit_conf_dropdown.png "Edit Configuration dropdown in Jetbrains IntelliJ/WebStorm"){"avif,jxl,webp"}

Once in the configuration screen, hit the **+** button at the top left, hover over "JavaScript Debug", and select Remote from the popout list.

@![Selecting a new entry type to add.](/images/{{\(index .Posts 0\).FrontMatter.Slug}}/edit_conf_plus.png "Selecting a new entry type to add."){"avif,jxl,webp"}

Now set it up for Chrome, with the localhost URL I posted above, and adjust it to your project (You can test in a normal browser before-hand). Once its to your liking, you can ignore the mapping, hit apply, ok. And you are set.

@![Setting up the remote configuration.](/images/{{\(index .Posts 0\).FrontMatter.Slug}}/edit_conf_remote.png "Setting up the remote configuration."){"avif,jxl,webp"}

Now here is the difference. You can no longer use the right-click menu to launch, as that always uses the "Local" configuration. You need to select your configuration in the dropdown in the toolbar, and hit the "Debug" button next to it. It should launch, with fulld debugging enabled.

***IntelliJ Note:*** IntelliJ Idea will get this a bit later. Features from the specialied IDEs eventually trickle back up to the main IDE, but it takes a bit of time. Usually EAP releases get them quick, but very buggy. At the moment, I can access the URL using IntelliJ, but it always returns a custom 404 page. So its not usable. It will be here soon. There is an open bug about it on their tracker.

## Full Steam Ahead

Hopefully this helps someone switch over to IntelliJ/Webstorm or if you are a current user, to use your environment to enjoy Dart.

The Dart Plugin works for IntelliJ Ultimate, IntelliJ Community, and Webstorm. So you can get IntelliJ CE for free, install the plugin, and you are ready to go. And if you want more features, you can always upgrade to the paid version or Webstorm. (Webstorm is just a subset of IntelliJ, everything that goes into the subset IDEs, gets a Plugin/Integration into IntelliJ, so IntelliJ gets everything, but if you want a focused (and cheaper) IDE, the subeset ones are a good buy)

### Links

* [IntelliJ Idea](https://www.jetbrains.com/idea/ "IntelliJ Idea")
* [WebStorm](https://www.jetbrains.com/webstorm/ "WebStorm")
