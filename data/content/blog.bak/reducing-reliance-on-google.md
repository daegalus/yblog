---
author: Yulian Kuncheff
date: 2022-09-06T08:57:00Z
draft: false
slug: reducing-reliance-on-google
title: Reducing Ones Reliance on Google
type: blog
tags:
  - google
  - alternatives
  - self-host
---

There are many efforts around De-Googling all over the internet, but I personally find those are not realisitic for the majority of the people. Many people have no problems sacrificing privacy for convenience. Or they aren't tech savvy enough to find stuff out themselves.

What many in communities for degoogling, or the ones that will read this blog post forget sometimes, is the masses are very technically illiterate still. Many of them don't know the difference between the Browser, the Internet, and their Operating System. Now while I agree that companies like Google exploit this fact, we also need to realize we can't take the all or nothing approach that is popular around the internet.

Removing Google from your life might be a goal, but for some, just reducing their reliance on Google, or having a minimum a backup plan, might be sufficient.

So I will be talking about how to reduce your reliance on Google by offering what I have found after a week of research and experimenting with alternatives, that might help you form a backup plan or find new alternatives. Many of these are listed on other sites and lists, but many are avoided due to an anti-corporation, pro-privacy stance.

As a brief aside, I have no problems with giving my information to Google, Microsoft, Amazon, etc, especially since it improves my conveniences. Likewise, I am very supportive and respectful of the efforts by the community to create privacy-oriented alternatives. They just aren't there for me or I don't want to bother self-hosting.

This is also a US-centric list. It is where I live, so it is hard for me ot offer EU, or other region offerings. The only real reason I mention Yandex, is I am part Russian, and I hear about it a lot from family, I don't use it at all myself.

## Why Not Self-host?

I am very tech savvy and I am a DevOps engineer. I can self-host something with my eyes closed, and even I don't want to bother. Maintaining servers and automation is my day job, and doing it at home is exhausting. Yes, self-hosting is low effort for the few services, but if you do it at home, what do you do if your internet goes out? What if I am running services for my whole family? Do they just no have access to things while my internet is out? Power outages? UPS systems only last so long.

I personally run Miniflux, FoundryVTT, Seafile, and Caddy in docker, on a couple of Raspberry Pis at home for myself. But I would never rely on them for anything important, especially during winters where I live, we get power outages that can last 2 days.

Hosting in the cloud is better, but then you have to learn about VPS providers, get something launched, and be comfortable with a Linux command-line.

Now I will have a section at the end about Nextcloud and how running the Hetzner ShareStorage might be a good option as a middle-ground.

But long-story short, self-hosting isn't a viable option for most. People need services that are hosted and managed by others they can just install and use.

## Alternatives

One thing I want to preface the following list of things, is that I will be suggesting other corporations that might not have any better of a track record for privacy than Google. Again, my goal is to create an backup/alternative plan, not to degoogle my life. This also will be a fairly Android focused article, but I will mention a few Apple options. I also run Linux myself, so most of these solutions have Linux options or web access.

### Google Mail/Calendar/Contacts

This one is fairly easy, there are a slew of great alternatives out there, from corporations to small companies, and everything inbetween, I will list a couple, but many more exist.

* [Outlook.com](https://outlook.com) (Microsoft)
  * This is a good alternative, and is already familiar to many that work in companies that use Microsoft products.
  * It can also act as a client to your Google account, so it could be a great way to test it before you switch.
* [ProtonMail](https://proton.me)
  * This is an awesome alternative that is privacy-focused, independent company.
* [Fastmail](https://www.fastmail.com)
  * A great private service that focuses on features and being better than Gmail.
* [StartMail](https://startmail.com)
  * Also a great privacy-focused alternative from a small independent company, that also runs a search engine.
* [Hey](https://hey.com) (no calendar)
  * A unique alternative that has a different take on email, but is definitely an option.
* [Yandex Mail](https://mail.yandex.com)
  * Might not be the best choice due to current events, but its still an option.
* [Apple iCloud](https://icloud.com) (Apple/iOS only)
  
And there are many more I am not listing like Posteo, Mailfence, etc. Email is an easy jump.

My biggest suggestion is to get your own domain, then signup for a service like ImprovMX (has a super generous free tier), and just route your custom domain email to any provider. You can also pay for the upgraded version that gives you access to their SMTP servers so you can send emails through them, and no thave the `on behalf of` show up on your emails when sending through Google for example or something.

This will let you use any provider, for example I use a plain gmail email behind my personal email, but my mother for example uses outlook.com. If i need to switch, I just go into the admin panel, and change it from the `@gmail.com` email to say an `protonmail.com` email, without affecting anyone else.

But regardless, do yourself a favor, get your own domain, ideally from a good trusted provider like [Namecheap.com](https://namecheap.com) or [Cloudflare.com](https://cloudflare.com), but there are also others.

### Chrome

Also fairly easy with how many forks of Chrome there are. I will just list them without too much description because they generally offer similar experiences.

* [Vivaldi](https://vivaldi.com) (what I use currently)
* [Edge](https://microsoftedge.com)
* [Brave](https://brave.com)
  * has some cryptocurrency stuff built-in, but it can be disabled.
* [Firefox](https://firefox.com)
* [Bromite](https://www.bromite.org)
* [Ungoogled Chromium](https://github.com/ungoogled-software/ungoogled-chromium)

Honestly, to use the modern web, you need something that uses Blink, Webkit, or Gecko. There aren't any other engines out there that are up to snuff, and even Gecko/Firefox are having trouble keeping up with Blink.

### Home

To be honest, there is only really 1 ecosystem that matches this, and its Amazon Alexa and the Echo devices. They work fine, Ive had an Echo, but I prefer Google's ecosystem personally.

HomeAssistant is a great self-host option, and since it rarely needs access outside the network, its one self-host option that really works well, and if the power is out, you can't do much anyway. I use this to manage a bunch of Zigbee motion sensors, and ESP32s running ESPHome for some faerie lights controlled by the motion sensors.

### Google Assistant

This also has few options, but there are some.

* Alexa (Amazon)
  * The primary competitor
  * Does not need Google Play Services to function
* Bixby (Samsung)
  * This is a good option if you have a Samsung phone, I don't find it as good, but it's there.
* Siri (Apple only)
  * If you are part of the Apple walled garden, this is your only option, but if you are ok with switching to Apple, its an option.

Microsoft's Cortana doesn't really exist for consumers anymore.

There is also Soundhound's Hound, but I haven't heard much from them in years, and not sure how good it is, but the Play store reviews aren't promising.

### Docs

Here we have a few more options. Some are less feature-filled than others, but still an option.

* [Office 365](https://office.com) (Microsoft)
  * The major competitor and again has a huge foothold in enterprise.
* [Zoho](https://zoho.com)
  * This is one many forget about but Zoho has a LOT of options for fairly cheap prices ranging from Email, Docs, CRM, and more.
* [Cryptee](https://crypt.ee)
  * This is one I am watching very very closely. I love what they are doing and going and hope they grow the products they offer. Only reason I don't use this myself is their Photos offering doesn't have an auto-upload for phones. It is a PWA so its easy to get on any phone or device though.
* [Yandex Disk](https://disk.yandex.com)
  * Disk has document editors.
  * Again, an option for those that are ok with it.

### Drive

This is a growing space and lots of competition, which is great for the user. I will keep the descriptions minimal, as overall they are pretty similar. It is all around what you feel is good for you, pricing-wise. Many of these offer a Docs suite also or are part of one.

* [Syncthing](https://syncthing.net)
  * This is less of a cloud hosting, but more of a sync solution between computers, regardless, it is something I use a lot to keep files in sync between my 2 laptops, desktop, work computer, and phone. Its highly configurable. I recommend [Syncthing-Fork](https://github.com/Catfriend1/syncthing-android) for Android, on Play Store and FDroid (fdroid recommended).
* [Backblaze](https://backblaze.com)
* [Mega](https://mega.nz)
* [OneDrive](https://onedrive.com) (Microsoft)
* [Tresorit](https://tresorit.com)
* [Zoho WorkDrive](https://zoho.com/workdrive)
* [Dropbox](https://dropbox.com)
* [Cryptee](https://crypt.ee)
* [NordLocker](https://nordlocker.com)
* [Yandex Disk](https://disk.yandex.com)
* [iCloud](https://icloud.com) (apple only)
  
### Authenticator

Another quick and simple one. Google Authenticator is a fairly standard RFC6238 app. There are tons of good alternatives, I will list just a few

* [Bitwarden](https://bitwarden.com)
  * I believe its part of the paid features, its what I personally use.
* [Aegis Authenticator](https://getaegis.app)
* [Microsoft Authenticator](https://www.microsoft.com/en-us/security/mobile-authenticator-app)
  * Really good if you are using Microsoft products, and mandatory if you set your Microsoft account to Passwordless
* [Authy](https://authy.com)

### Google Chat

Honestly, Google Chat isn't as popular as many of the alternatives, so its not that huge of a change for something better, but since its kind of jammed into th GMail app, I am sure some use it.

* [Slack](https://slack.com)
* [Teams](https://www.microsoft.com/en-us/microsoft-teams/group-chat-software) (Microsoft)
* [Element](https://element.io) (Matrix-protocol)
* [Telegram](https://telegram.org)
* [Discord](https://discord.com)
* [Yandex Messenger](https://yandex.com/chat)
  
### Google Meet (Duo)

Since Duo was merged into Google Meet, this one is also fairly easy to provide replacements for.

* [Zoom](https://zoom.com)
* [Teams](https://www.microsoft.com/en-us/microsoft-teams/group-chat-software) (Microsoft)
* [Slack Huddle](https://slack.com)
* [Telegram](https://telegram.org)
* [Discord](https://discord.com)
* [Element](https://element.io) (Matrix-protocol)
* [Jami](https://jami.net)
* [Jitsi](https://jitsi.org)
* [Facebook Messenger](https://messenger.com)
* [Skype](https://skype.com) (Microsoft)
* [Yandex Telemost](https://telemost.yandex.ru) (Russia only it seems.)
* [Facetime](https://www.apple.com/mac/facetime/index.html) (Apple Only)

### Google Pay

This is where it gets tricky. There aren't a whole lot of options, and many are region specific.

* [Samsung Wallet](https://samsung.com/us/samsung-wallet/)
  * Mostly for Samsung phones, but I am fairly certain you can use it on any android phone. Samsung Pay got merged into Samsung Wallet. Also on Play Store.
* [Paypal](https://paypal.com)
  * Used to have direct NFC support, but now it goes through Google Pay. You can still use Paypal QR Code payments at supporting retailers.
* [YooMoney](https://yoomoney.ru) (Yandex)
  * Only works in Russia and former Soviet territories, or allied ones.
* [LINE Pay](https://pay.line.me), [Kakao Pay](https://kakaocorp.com/page/service/service/KakaoPay?lang=en), [WeChat Pay](https://pay.weixin.qq.com/index.php/public/wechatpay), etc
  * This one is mostly in Japan, Korea, China, etc.
* [Apple Pay](https://apple.com/apple-pay/) (apple only)

### Google Photos

This is also a tough one, few if any options have the same AI capabilities that Google does. But for Photo storage, there are a few options

* [OneDrive Photos](https://microsoft.com/en-us/microsoft-365/onedrive/online-photo-storage) (Microsoft)
* [Cryptee](https://crypt.ee)
  * I just wish they had an auto-upload app. I would switch now.
* [iCloud](https://icloud.com) (apple only)

You can also use any of the Drive alternatives.

### Messages (RCS)

If you want to use RCS, there are no other options. If its just SMS/MMS there are many other SMS apps, too many to list.

Now, if you can move off of SMS/MMS or convince the people you talk to there, there are better things to use.

* [Facebook Messenger](https://messenger.com)
* [WhatsApp](https://whatsapp.com)
* [Telegram](https://telegram.org)
* [Signal](https://signal.org)
* [Wire](https://wire.com)
* [Threema](https://threema.ch)
* [Session](https://getsession.org)
* [FluffyChat](https://fluffychat.im) (Matrix Protocl)
* [DeltaChat](https://delta.chat) (onto of email)
* [Pony Messenger](https://ponymessenger.com)
  * This one is nice for slowing down. You only receive messages once a day, regardless of when they were sent.

### Maps

This is a tough one. Google Maps is one of the BEST options for points of interest and related data. But there are a few that I managed to check that had decently updated information. There are a few locations near me that recently opened, or opened during the pandemic, and very few maps have it updated.

* [Here WeGo](https://wego.here.com)
  * Managed by a group of Automobile manufacturers like BMW, Daimler, etc. It is fairly updated, and had a Habit Burger location that opened up only 2 months ago in my area.
* [CityMapper](https://citymapper.com)
  * Focuses on Walking/Cycling directions, Uses a collection of sources, including Apple, Google, OpenStreetMap, Foursquare, Yelp, etc. This was up to date also and had lots of features
* [OpenStreetMap](https://openstreetmap.org)
  * Lots of apps for it, including the following
    * [OsmAnd](https://osmand.net)
    * [OrganicMaps](https://organicmaps.app)
    * [Maps.me](https://maps.me)
  * This is crowdsourced map system. Directions are good, and fairly updated, it had everythin but the Habit Burger on there, but the great thing was, I spent 5 mins adding it, and it will probably show up in all OpenStreetMap based software soon. (I also removed the old Jack In the Box that Habit replaced, before we knew what was replacing it.)
* [Bing Maps](https://bing.com/maps) (Microsoft)
  * This is in partnership with Here WeGo (same data), and they no longer have a dedicated android app, so only browser access.
* [Apple Maps](https://www.apple.com/maps/) (Apple only, mostly)
  * [DuckDuckGo Maps](https://help.duckduckgo.com/duckduckgo-help-pages/features/maps/) uses Apple Maps for showing maps of things, locations, ratings, etc. But you can't navigate or get directions directly, it will use Google Maps or your local maps app for it on Android. Might be some contract with apple that might not allow it, or a limitation of APIs.

Other than that, many other solutions were very very out of date. like Sysig. And there is a chance OpenStreetMap apps are out of date in your area, if no one has done the work to update/add things.

### Translate

There are some AI based solutions popping up and good alternatives. Not much more to say.

* [DeepL](https://deepl.com)
  * Probably the best alternative and works really well with the languages it has, but doesn't have as many as Google. Uses AI and Machine Learning to do the translations. Had pretty great results with Bulgarian.
* [Lingva](https://lingva.ml)
  * Also a machine learning based translator, does fairly well in my tests with Bulgarian.
* [Bing Translator](https://www.bing.com/translator)
  * Haven't used this much, but it's an option.
  * The app for phones is called Microsoft Translator

### Google Voice

There are few services that offer similar, most focused on businesses. But I did unearth some options.

* [Skype](https://skype.com) (Microsoft)
  * Skype has a few features that can match Google Voice. First [Skype Number](https://www.skype.com/en/skype-number/) (I assume this used to be called SkypeIn), which lets you have a phone number people can call, that will ring you in Skype. Another is the ability to call phones, with very nice rates for [international](https://www.skype.com/en/international-calls/). This used to be called SkypeOut.
* SIP Provider
  * There is a roundabout way to sign up for a SIP number, and hook it up to your phone, as Android supports SIP numbers in the phone app on most phones.

### Youtube

There really is no alternative. You can go into Decentralized services like PeerTube, but quality of the videos and potentially the content material might not be what you are interested in. Regardless, decentralized or blockchain based ones might be the future, so I will list a few things here, and a few apps to look into as alternative clients.

* [Freetube](https://freetubeapp.io)
  * Alternative frontend for using Youtube on Desktop
* [Newspipe](https://newpipe.net)
  * Alternative app for using Youtube on Android
  * There is also Newspipe Sponsorblock which has sponsorblock builtin
* [LBRY](https://lbry.com)/[Odysee](https://odysee.com)
  * A growing platform built ontop of blockchain that is gaining some traction and content. Odysee is a web frontend for LBRY

### Android

Google forces most phone makers to include a bunch of Google services pre-installed. If you want to minimize these, you need to use a Custom Rom, or something like [LineageOS](https://lineageos.org). With those, you can then install a barebones Google Play Service install, which just adds the absolute bare essentials for apps to work properly on your phone and the Google Play Store to get more stuff. I personally have tested [NikGApps](https://nikgapps.com)' Core package, which does just that, and everything I wanted to work, just works. NikGApps also has a bunch of versions with more stuff included, but that is up to you to decide.

If you want to completely go away from Google, you can try [MicroG](https://microg.org) as an alternative. Or switch ot Apple.

If you need to replace Google Play Store, [Aurora Store](https://auroraoss.com) is a great option (Can be downloaded on F-Droid also).

Another option if you want to bail from Android altogether and have a supported device is [SailfishOS](https://sailfishos.org), it is linux-based but can run Android apps.

## Final Thoughts

Overall, if you want that 1 stop shop experience, you will mostly have to switch to Microsoft services. Amazon doesn't compete, and there aren't many smaller options. There is /e/ Foundation's cloud services, but they are mostly a NextCloud or other open source software as a solution. There is nothing wrong with that, and it is a viable solution, it just needs a bit more time to mature and work out the kinks.

I think overall we really need a competitor that focuses on consumer tech like this, because Microsoft is honestly more focused on Enterprise than everyday consumers, as their dwindling consumer apps show (many closed due to Microsoft not able to gain traction). Zoho used to be a nice place for individuals or families, but they have pivoted hard to small business. You can still use them for Email, Drive, Docs, Password, etc. But it's been refocused to business needs.

If I were to lose Google tomorrow, as some have (like the recent NYTimes articles) for crazy reasons, I would probably use mostly Microsoft stuff, along with maybe Cryptee for Docs, Drive, and Photos if they get an auto-uploader up. Probably with Here WeGo as my maps app, while having OsmAnd on standby. I would finally try to get all the people I communicate with on Telegram or FB messenger, maybe whats app, and offer DeltaChat as an option. For storage, I already have Office 365 through my dad who has a subscription for 5 people and have 1TB of storage there. But I could use NordLocker fairly easily if needed or Backblaze if they would ever release a Linux app.

On the Android front, I already tested much of this on my old OnePlus 6T (primary is a Pixel 6 Pro), and using LineageOS + NikGApps Core, let me do everything I need, including adding F-Droid and getting some great apps from there too. Most of the alternatives above, worked great and flawlessly with minimal to no Google involvement when using Aurora Store. Didn't even log into my Google Account on the device or any other app.

I only use Linux on my machines, no Windows, but I can get most Microsoft stuff to work there anyway, especially the web-based stuff.

## Nextcloud

So the reason I left this for last, is because while it is a Self-hosted solution, there are ways to pay for a managed solution that you just administer.

This is a great middle-ground between shirking all corporations, but having a managed experience.

I personally tested [Hetzner's StorageShare](https://www.hetzner.com/storage/storage-share), which is a managed Nextcloud instance by Hetzner. I tried thier $5 base option, which includes 1TB of storage. it goes up to $15 for 5TB, and $27 for 10TB. I set it up and it worked great if you want to go this route.

Email will still need to be setup with a different provider, but you can use Nextcloud as a client. It also offers a Google Integration to import your google stuff into Nextcloud, and many other features.

This solution is honestly one I might use more and more over time, the concern is convincing my family and my wife to use it, which isn't always easy.

## Not the US

Another thing is in many other countries that might have better alternatives. I know in Japan, LINE has a huge slew of services. South Korea has Kakao, and China has WeChat, maybe something from Tencent I think. India has other services, same for Europe and different European countries. Russia has Yandex, and a few other competitors.

I just didn't mention them as I have little exprience or knowledge of them.
