---
author: Yulian Kuncheff
date: 2013-02-21T20:00:00Z
draft: false
slug: why-are-there-no-true-cross-platform-filesystems
title: Why are there no true cross-platform filesystems?
cover_alt: A table of filesystems
type: blog
tags:
  - filesystems
  - operating-systems
  - cross-platform
---

@![A table of filesystems](/images/{{\(index .Posts 0\).FrontMatter.Slug}}/cover.png "A table of filesystems"){"avif,jxl,webp"}

As someone that has recently installed multiple operating systems on my desktop, I find a major problem that I think shouldn't really exist, but does.

Over the last few months I have made it so that my PC desktop runs Windows 8, Linux Mint 14, and OSX Mountain Lion. This was a long endeavor, especially OSX. But everything works nicely and I have the development environments I need on all 3 platforms and all the additional chat, social, and browsing software I like to have running. So that no matter which platform I am on, I just work as normal.

But there is one massive thorn in my entire setup. Filesystems. After so many years of multiple OSes, I am surprised that there isn't a single filesystem that is truely cross-platform and modern. The keyword is modern. FAT32 is still the only candidate for a true cross-platform Filesystem, but its dated, it has very hard limitations on file-sizes that makes it hard if you want to store 4gig+ movies, games, and other data. Also running a Plex server using data on the drive, torrenting, and games requires fast reads and writes.

I will outline a few of the downsides to many of the primary filesystems.

### EXT4 (& 2/3)

So the first step was the possibility of just adding EXT support to Windows, at the time thinking OSX had out of the box support for EXT4.

But after countless searching, all I found was old and no longer maintained drivers that supported EXT2 and parts of 3 (usually no journalling support). After a while, I decided I should get OSX setup, and then come back to this, thinking OSX would be just as easy.

To my surprise, OSX had the same problem. There were a bunch of FUSE projects for OSX, most out of date, OSXFUSE seemed updated, but then I couldn't find a viable EXT3/4 fuse driver. They all suffered the same problems as in Windows.

I was really surprised that there wasn't any good free EXT4 support on OSX. There is stuff like Paragon EXT and so on, but those cost money, and I don't know how good they are, because they seem to be just paid FUSE implementations.

So EXT was out of the picture, because on a massive secondary drive with lots of big data, modern features of EXT4 are very beneficial (like extents and journalling) and I would prefer them.

### NTFS

This I thought would be a decent alternative. I know there is ntfs-3g for write support in Linux, and I got pissed off one day, and just baught Tuxera NTFS for OSX for my Laptop, so I used that on the OSX partition.

Read performance was generally good, but one thing that was a huge blocker was write performance. Both ntfs-3g and Tuxera NTFS had poor write speeds. I compared them using the same file, starting from the native filesystem of the OS to the NTFS drive. Windows was obviously the fastest, ntfs-3g was ok, but took about 10 mins longer for the same file, and similar poor results on OSX with Tuxera.

Forget about downloading any torrents to the drive from a non-Windows OS. Takes way too long for it to actually write anything, and lots of corrupted blocks.

Not to mention the third-party implementations can be buggy, cause corruption, and other things. I didn't experience them, but the possibility is much higher with them. Being as all of them are reverse-engineered.

So that kinda fell through. Even though its still what I have the drive set to, but I am not happy with it at all.

### HFS

This is fine on OSX and Linux, but forget about it on Windows.

Linux support was almost there, but with limitations. Anything over 2TB would corrupt, and no journalling support. But it was better than nothing.

Windows though, only read-only implementations. Even with an official driver from Apple for read-only in bootcamp. Without write support its useless.

### ExFAT

I thought this could be the saving grace for me, then I found about the so-so Linux support, and it doesn't have any modern filesystem features. Just mostly an enhanced FAT. You need to use a fuse driver for Linux, which isn't too bad, but not ideal. Patent encumberance doesn't allow it to be included into the Kernel.

Not to mention it is a closed spec, so all third-party implementations are reverse-engineering and potentially have major problems.

### FAT32

This just is not suitable for huge drives and large data, its old, very restricted, and somewhat useless for anything outside small flash drives, and possible legacy support.

### UDF

I actually strongly considered using Plain UDF 2.01. Its supported by all the OSes. Plain version works well on hard drives. Linux doesn't support writing on anything above 2.01. Maybe we should fix that.

Maybe a few mounting issues on Windows at first, but it seems like a good solution.

Just this is a new idea to me, and I want to look into it more before I try. At first glance, it seems a lot like FAT and doesn't have a lot of big data fancy features. Though it is used for DVDs and Blu-rays.

Maybe a Spared version could be better. I need to research more into this, but this might be the only good solution.

### Btrfs

This has only Linux support atm, but as a filesystem still in development, I would love to see this get fully featured Windows and OSX drivers by the time it hits 1.0. This would be amazing, and would be the first open source, free, modern filesystem to do it.

But I am not holding my breath, OSX maybe, but there seems to be some severe hate for Windows by Linux diehards, so I don't think it will happen. But its nice to dream.

### May other Filesystems only supported by Linux

There are a slew of filesystems out there, a majority are only supported by Linux or *nix derivatives. Some have read support on OSX through MacFuse and OSXFuse, but only read support.

And forget Windows, nothing there.

### Conclusion

There isn't a single quality modern filesystem that is cross-platform. NTFS is the closest, but it suffers from performance problems on non-Windows systems.

Maybe I should just buckle-down for the next 2-3 years, and make it myself. I am just more of a high-level developer. I live in languages that are interpreted or on VMs. I know enough C/C++ to get by, but not enough to make efficient OS drivers.

I think this is an issue that needs to be remedied. At least OSX and Linux. Windows can come after.

For now, I might give UDF a try, not sure yet. Guess I will stick to the slow NTFS on OSX/Linux.

> [Discuss on HN](https://news.ycombinator.com/item?id=5272960)
