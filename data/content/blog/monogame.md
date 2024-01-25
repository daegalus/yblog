---
author: "Yulian Kuncheff"
date: 2012-04-30T19:00:00Z
description: ""
draft: false
slug: "monogame"
title: "Monogame"
cover_alt: "Monogame Logo"
type: "blog"
tags:
  - gamedev
  - programming
  - frameworks
  - engines
---

<img src="/images/${slug}/cover.png" transform-images="avif jxl webp png" />

*Disclaimer: This is all opinion and personal preference. A lot of the frameworks and game engines I talk about are great, and its my own choices that make them subpar for me. Please take this as an opinion peace, I am not trying to sway anyone in any way, it is just an explanation on why I chose what I did in hopes it might give someone else some ideas.*

Ok, so this really isn't much a rant, but more of a why I chose to go this approach over some other approach. I recently wanted to actually "make" games. Reason I have it in quotes, is because I always have an idea, I do the research, write the first 20 lines, and get a new idea and switch to it. Nothing ever got done.

I finally found a major reason (but not the only reason) on why I do this. Out of personal choice and preference, I want any game I make to be as cross-platform as possible. This puts a real damper on things because there are very few frameworks or platforms that are cross-platform enough. There is always an important platform missing, be it Linux, Windows Phone 7, or a combination of many missing ones. I have gone through many frameworks, including Haxe, Corona, Monkey, Unity, and even less frameworky setups. I finally ended up with Mono and MonoGame, and I will detail why I chose this after I lay out what I have tried and why it didn't work.

Second, it is most programmers' dream to write code once, and release to multiple platforms. Even with minor modifications per platform. So I am naively still chasing that idea, as many others are too.

Third, I want to be able to prototype and build the game as quickly as possible. I have too many ideas, and I tend to switch too often and never get anything done. So taking the time to write my own Engine, Library, or Framework would really slow me down and I would probably lose interest. BUT, I might make it a project to write my own Javascript Game Library/Framework. It sounds like a fun project in of itself.

## Mobile Only Frameworks

I am going to preface this by stating that I understand that WP7 support for many frameworks is not possible because it doesn't allow for custom runtimes to run on the phone mostly because it doesn't allow low level access to the system or C/C++ access. So WP7 is really not something I will analyze or really mention often.

There are many frameworks that are made specifically for only mobile targets. This is fine, as long as I can find a complimentary Desktop library if the game is geared to Desktop play alongside mobile play. If its mobile only thats fine, but even then it has problems. One of the major problems is frameworks tend to be geared for 1 platform. One popular framework is Corona. It looks like a superb platform, but its missing some key features and has been missing them for as long as I have known about Corona. One being Shaders. I have always waited for them to add shader support, but it is constantly getting other (in my opinion) less important features.

But there are so many mobile frameworks out right now. It is getting a bit ridiculous. During my search I found at least 10 different frameworks that do almost the exact same thing, targeting iOS or Android. Some throw in Blackberry, Symbian, and HTML5. But in the long run, they don't provide a benefit over each other and tend provide lackluster feature support of the target platforms.

So all those get scratched off the list right off the bat.

## Desktop/Mobile Frameworks

So now we go into the realm where 1 framework can target multiple form factors, paradigms, and platforms. There is always a need to understand that a certain game might not be playable or as fun on a certain form factor or control scheme. So you will always need to evaluate which platform to push to. But being able to make a game that fits all the form factors, and play well is something that is possible, and seems to be a growing trend.

So I first started out with some fan favorites like Unity3D. First thing, as the name suggests, it has really no support for 2D game making. You can emulate 2/2.5D games if you fix camera positions and do various other tricks, but you still need to have a solid understanding of 3D game programming. But I had a few 3D game ideas, and I was willing to make some of my 2D games into 3D. Now, I am fully aware that this might be due to my lack of experience in game making, let along 3D game making, but I find Unity3D horribly confusing to use. Maybe its the different structure, but the way the scripts meshed with the scene, and the files. I really had a tough time figuring out how to access and control objects in the scene from the code. So Unity3D got a huge pass from me because, I just couldn't use it. I will revisit it in the future when I am more mature in this, and it might make perfect sense, but right now, I feel really out of place using it.

There are 2 other frameworks that offered the same choice. One was DeltaEngine, and another was ShiVa3D. Both are great in their own respect, and DeltaEngine even uses C# (one of my favorite languages, but more on that a bit later) but it is still new and by the time I considered it and ShiVa, I had figured I was too inexperienced to jump straight into 3D, at least through a framework (writing 3D from scratch doesn't seem all that daunting, even though it probably is). But another killer for me on ShiVa3D and Unity is now WP7 support. Though, due to WP7 limitations, I understand that it might not be even a valid choice for them as they can't port the runtime or the C/C++ code. Maybe in WP8.

## HTML5/Javascript

After searching countless hours, I decided I should just go the Javascript/HTML5 route. I greatly enjoy non-Website Javascript programming and I found plenty of frameworks that can package it all up into a nice package for every platform. Desktops were a bit harder to find, but there are solutions like TideSDK and Air. Along with writing my own wrapper setup with Chromium Embedded. Mobile is not an issue as there are things like PhoneGap that will allow me to do just that on most mobile platforms. If the platform doesn't have a packager, all the important mobile platforms have HTML5/JS supporting browsers, so it can run in the browser.

Javascript really provides the true Cross-platform support, and with Node.js (my server side language of choice), I can really write everything in 1 language and 1 codebase.

But I hit a stumbling block. Most engines, frameworks, or even rendering libraries are in very early stages of development. They work, but not all the way. And I have seen a lot of problematic support with the packagers. I can always write most from scratch, and even release it open source as a library for others to use. But Javascript is still newer to me. I didn't really join the bandwagon until Node.js (and I joined that around version 4.10ish). I was never fond of the language until I saw its real power outside the browser. This is probably all just excuses, and laziness on my part, but I really did not feel comfortable rewriting Canvas drawing libraries. I have looked into CAAT, Craftyjs, Impact3D, GameQuery, and a few others. Impact3D was a no go because it requires purchase before usage. I will not buy a framework that I don't even know is the right fit for me. GameQuery is still in early development and I don't think it is ready for primetime. CAAT and Craftyjs are the most mature, and are very solid frameworks. I just found their API a bit limiting. But don't get me wrong, these are great frameworks/engines. I just doesn't have the tools I am looking for. But if I had to do it again. I would probably go with Craftyjs.

I am sure I am forgetting some good ones, but I can't remember them off the top of my head.

## So what was left?

After lots of searching, I re-stumbled on something that I had seen long ago, but did not think much of it at the time. It was something I saw after the Windows 8 development info came out and the future of XNA came into question. I was mostly worried about WP7 support at the time. The solution I found was MonoGame, an open source re-implementation of XNA with works with .NET or Mono, which is an open source implementation of .NET and a C# compiler.

So why did I end up choosing it? Well firstly, I really enjoy C# as a programming language. It just clicks for me and has a lot of features like Lambda Expressions, Anonymous functions, etc. that I find very useful and allows for the language to be used functionally. It is a language I have grown to like a lot, and I am glad there is a project like Mono that has most of .NET 4.0 ported, and I think now a good chunk of 4.5 in the dev branches. Plus the .NET Standard Libraries are really well done, and I find much greater stability in .NET than I do in Java's Standard Libraries. Plus all the dev tools are top notch. (I hate Eclipse with a fiery passion)

Then comes MonoGame. Which takes XNA, a very easy to use and really well supported framework. Since it is a direct compatible API, XNA examples and test games actually just work in my own tests with the beauty of compiling them and running them on OSX and Linux.

But the real kicker, is they are using SharpDevelop to add support for Windows 8 and Windows 8 metro. So I can take my game, compile it for Metro and put it on the up and coming Metro App store. This breathes a lot of life into a framework that was frankensteined apart and put into WinRT. And allows it for a path to grow on its own and add more DX10/11 options that were not possible before. Also means, Win8 Arm Tablet support!

Another benefit is MonoDevelop, which allows me 1 IDE for all the platforms and allows me to work on my Macbook Air without the need to Bootcamp it and waste space installing Visual Studio.

Then came some amazing bonuses. Xamarin's efforts in getting the Mono runtime to work on Android and iOS. This allows for MonoGame to just run directly ontop of those Mono implementations. If I remember correctly, Xamarin manages the OSX binaries and development too.

So we have a mature and supported language with a solid standard library. A pretty much game making framework with a big community and samples. Support for all the important platforms (Windows XP - 8/Metro, OSX, Linux, Android, iOS, WP7, XBOX, or anything Mono runs on), Multi-platform development tools, and future-proofing. Also Playstation Support is planned.

So what more can it really give? ... HTML5/WebGL. Using a tool I found just recently, you can take the CIL bytecode and compile it into pretty readable Javascript. The tool is called JSIL and it is still in early development but the website has some really good examples and ports to show off. Now, WebGL is a bit limiting, but most browsers support it (with the exception of IE), no mobile browser support, but the can be worked on too..

Also, I believe there might be a NaCl target some day. Bastion was able to make their Chrome App Store version using MonoGame. This adds yet another option for in-browser gaming. So that adds another platform to the Package. I also believe Bastion uses MonoGame for their recent OSX release.

## This is all great, but what is the catch?

There is no real catch, more like caveats. They need to be considered when porting things between platforms and a few other tidbits.

First, 3D support in the current MonoGame 2.5 does not exist. But fret not. the current main development branch is very active on 3D support, and there are already test projects running 3D XNA. So it won't be long now when MonoGame 3.0 comes out and we get full 3D support. Just a bit more waiting.

Second, the Content Pipeline is not implemented in MonoGame. At least not in the full sense. It still has the content manager for loading files. But it doesn't use XNBs and it just loads them directly form their normal format. Second, Mono and MonoDevelop do not support the Content Project (.contentproj).

So this issue cropped up for me. I had started my project on OSX in MonoDevelop. All was fine in dandy and I had a major chunk working. Then when I tried to run it on Windows, it wouldn't compile. I found out that XNA requires a Content Project with a Content Reference in the Main project. No leeway. So after frankensteining the XML in the project files, and doing some pathing fun. I got my project to use the Content folder in my Main project as a Content Project with a reference to it. A bit confusing, but it works, and it didn't break the MonoDevelop on other platforms or my other Solutions for the other platforms. But, it did break MonoDevelop on Windows. Without knowing how to work with the Content Project file, it just would not load the content project, since it couldn't pull the Content. So I am forced to use Visual Studio on Windows, until they add better support for Content Projects into Mono/MonoDevelop/MonoGame.

But, once I got all that taken care of, it ran beautifully on all platforms. Oh and another note. You need to make a Solution file for every platform. So you can set the proper settings for each platform and to allow you to use Pre-Processor If/Endif statements to add platform conditional stuff.

So while not ideal, Mono and MonoGame has actually provided a way to write 1 code, and port it out to 9 platforms (8 if you don't count Win8 Metro as a separate platform) and add 1 for planned future support for PLaystation Suite.

Lets hope it works out and stays good.

Links:

* Mono - [www.mono-project.org](https://www.mono-project.org)
* MonoGame - [www.monogame.net](https://www.monogame.net)
* JSIL - [www.jsil.org](https://www.jsil.org)
* CraftyJS - [www.craftyjs.com](https://www.craftyjs.com)
* GameQuery - [www.gamequeryjs.com](https://www.gamequeryjs.com)
* CAAT - [labs.hyperandroid.com/animation](https://labs.hyperandroid.com/animation)
* ShiVa3D - [www.stonetrip.com](https://www.stonetrip.com)
* Unity3D - [www.unity3d.com](https://www.unity3d.com)
* DeltaEngine - [www.deltaengine.com](https://www.deltaengine.com)
* Corona - [www.anscamobile.com](https://www.anscamobile.com)
