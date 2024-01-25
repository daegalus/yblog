---
author: Yulian Kuncheff
date: 2012-09-29T19:00:00Z
draft: false
slug: language-bloat
title: Language Bloat
cover_alt: A bloated python representing language bloat
type: blog
tags:
  - programming
  - std-library
  - languages
---

<img src="/images/${slug}/cover.png" transform-images="avif jxl webp png" />

*Disclaimer: After reading through this a few times, I realized my train of thought jumps around a lot on this. There is just so much going on in my head on this topic, its hard to not jump around. I apologize for the all over the place article. Hopefully my point is still conveyed.*

> Update: I just want to clarify that when I refer to JVM/CLR I include their primary language and their standard libraries. The JVM and CLR are robust, generally small, and powerful VMs. The primary bloat gain on them comes from the primary language implementation and the massive libraries that come with it. Java is bloated, with it 100s of libraries and standard libraries were made, which in turn makes them bloated and in turn bloats the whole entire platform. If for example, we took the JVM, rewrote all the libraries in Scala, and used Scala, it would be much much lighter. The problem is Java and the bloat it has introduced into libraries and programs written in it on the JVM. Most languages on the JVM use the Java standard library as a base, so the bloat comes with it.

As of recently, I have been having a hard time enjoying the languages I program in. At first I was afraid it might be burnout, but it isn't. I still thoroughly enjoy programming, making things, and bringing my ideas to reality. So I sat down and just asked myself what I didn't like about them. It always came to the same point: Bloat.

Lately, all the fun and interesting languages that I enjoy using tend to be on either the .NET CLR or the JVM. That is all fine and dandy, but anytime I see that, I get this nagging feeling of distaste. I was curious on what was causing this when I really liked the language semantics and syntax.

All these languages had something in common, very bloated standard libraries and sdks. For anything on the JVM, I need to get the JDK which is a few hundred megs, then if its a language other than Java, I need to get that language's SDK also. For .NET, its all the Mono libraries or all of .NET and Visual Studio.

An example, Grails. I use it a bit at work, along with Groovy. So I need to get the JDK, Groovy SDK, and the Grails SDK. Of the top of my head, I think this was 200 megs in packages, libraries, and other crap, 90% of which I didn't need.

But then I realized I was having problems with Python and Ruby, which are much lighter than the above said setups. I thoroughly enjoy using Python, and it is much lighter. But why does it still feel bloated? For these two languages, it came down to the library/module systems. Python Eggs and Ruby Gems along with their package managers.

One of the many reasons I dislike Ruby (actually Ruby has been growing on my a bit, primarily because of Groovy, I still hate Rails, Sinatra or Flask for Python seem more my thing), is the Gem package system. It feels messy, and bloated. Python and Ruby by themselves aren't all that bloated, they are kind of midway between bloated like Java/.NET and what I consider a light platform (I will get to these in a bit). But their module systems make the whole thing a mess.

These factors tend to kinda put me off when I think about sharing my apps, programs, and tools with others, or if I consider deploying said server-side software and managing it. Java, Ruby, and Python tend to come pre-installed on most unix systems, or very easily installed with the respective system's package manager, but with it tends to come way too many supporting libraries and dependencies.

I have been apart of the Node.js craze, and I still think Node is a solid platform and would use it for various things, but my problems with it stem in the actual Javascript language, all its pitfalls and things I have to constantly work around, I really just don't want to deal with it. But the platform seems somehow lightweight and clean. Node itself is only a few megs. It has a basic standard library, and leaves the rest to third-party modules, or even if first party, they aren't part of Node. The modules themselves feel light aswell. Some are a mess and can be bloated, but overall, things seem less of a problem compared to Jars, Gems, or Eggs (I feel like I am talking about some scenery in Aladdin).

I also am a fan of Dart. Google might an insane uphill battle to even get other browsers to consider it, but the language and DartVM are a pleasure to work with, and it has the same light and fast feel of Javascript and Node.js. Dart is still in very early development, and I wouldn't use it for anything serious, but it is a pleasure to work with. If it doesn't catch on in the browser, I will probably use it as a general purpose language using the DartVM like many use Node.js.

But this brings me back to the other languages/platforms. I feel they could be implemented in much lighter and robust packages. And their module/library systems can be streamlined and lighter. I know Ruby has a package manager called Bundler, but I never used it, and I hear mixed reviews about it. For Python, I use pip, and it seems lighter than easy_install.

Maybe I am just delusional and this is all in my head, but there are so many fun languages, yet they suffer not because the language itself, but because of the platforms under them. I tend to enjoy more Python-like languages, like Boo, Groovy, Python itself, Genie.

Here is a list of languages I enjoy that feel ruined because of whats under them:

* JVM
  * Groovy (I use Java and Groovy on Grails at work)
  * Scala (Has a CLR emitter too)
  * Clojure (I know of ClojureScript, but then i deal with Javascript)
* CLR
  * C#
  * F#
  * Boo
* Python (Also the JVM/CLR versions and I currently tend to use Python for most of my personal work)
* Ruby (Also the JVM/CLR Versions)

Now, I have considered compiled languages, and I do want to use them, but then there aren't many interesting compiled languages. C/C++ I do not like working with. D is ok, but something about it made me feel meh about it. Vala and Genie seem awesome, but they are quite dependent on GObject and GLIB, so portability to Windows and partially to OSX is messy and a headache. (in all honesty, Genie could probably be for me). I have been trying Go lately, and I have a love/hate opinion on it. It feels awesome, but at other times it just feels odd.

I guess my ideal language is the lightness of Dart/Node and their module system, with the syntax of Python/Groovy/Ruby, and the popularity of Java/Ruby/Node. Wishful thinking.

You might be saying I should try to make my own and see it isn't so easy, and maybe I should write my own, but then it would be horrible and I know I wouldn't be able to make it as light and efficient as I want it to be. (I am taking a Compiler Design course atm, and a Programming paradigms course. One is having me do Pascal interpreter in Java that can also compile to C using JavaCC, and another is having me write a Scheme Interpreter in Javascript using JS/CC).

There is just this bloated feel I get from languages nowadays. Its the best way I can describe it. I guess I can always just make my won, but I don't have experience in that, and my distaste for C/C++ would hinder me from making it really efficient. I can make a Shim language around Dart, but i need to wait for it to mature some more. I could make my compiler/interpreter in Go. That might be an interesting project.

P.S. I need a new blog system, this one is nice, but dev on it seems dead, and this thing seems bloated, and its in PHP, blah.
