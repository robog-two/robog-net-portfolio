---
title: Designing a Human App in an Artificial World
layout: blog.liquid
date: 2026-07-23
tags: post
description: yes
---

<section id="slideViewer">
<img src=""></img>
</section>


![How to make a very human app in a very artificial world](./media/human-app-ai-world/User%20Research%20for%20Local%20Music%20App/User%20Research%20for%20Local%20Music%20App.001.png)


Hi, my name is Sam Knight. My recent research project is titled "Designing an App for Recommending Local Music." Or, to put it another way, I'm figuring out how to make a very human app in a very artificial world.

This work was done with my faculty advisor, Dr. Doug Turnbull, and with assistance from Fisher Griesel and Abe Manfra.

Before I get started, I'd like to address the *artificial intelephant* in the room.

![Let's address the artificial intelephant in the room](./media/human-app-ai-world/User%20Research%20for%20Local%20Music%20App/User%20Research%20for%20Local%20Music%20App.002.png)

This pie chart is from Quinnipiac University, and it represents whether people believe AI is more beneficial or more harmful to society.

![Artificial Intelligence Opinion Survey](./media/human-app-ai-world/User%20Research%20for%20Local%20Music%20App/User%20Research%20for%20Local%20Music%20App.003.png)

You see this larger blob, which is 55% of respondents saying AI does more harm than good overall. Then you have this smaller circle, which is 33%, saying AI causes more benefit than harm. And of course you have this undecided circle. (12%)

So which one am I?

![Which one am I?](./media/human-app-ai-world/User%20Research%20for%20Local%20Music%20App/User%20Research%20for%20Local%20Music%20App.004.png)

I would say let's challenge the premise.

![Challenge the premise](./media/human-app-ai-world/User%20Research%20for%20Local%20Music%20App/User%20Research%20for%20Local%20Music%20App.005.png)

I have a different view, because I think the more you research and the more you learn about this new technology, the more it becomes clear that it's too early at this time to form an opinion.

I think we should be this tiny sliver: people who aren't looking at it from an angle of more harm than good, nor more benefit than harm, nor undecided.

![Let's be this tiny sliver](./media/human-app-ai-world/User%20Research%20for%20Local%20Music%20App/User%20Research%20for%20Local%20Music%20App.006.png)

We should be informed, but we should look at it on a case-by-case basis. For the sake of this presentation, I invite you to leave your existing opinions at the door and look at this as a completely neutral technology. Instead, let's judge the inputs and outputs. Let's judge all of the resources needed and potential harms of using them as well as the outputs themselves including their potential adverse effects or how useful or beneficial they are.

With that out of the way, I'd like to go into a history of local music apps that my team and I have been working on.

![Local Music Apps: A History](./media/human-app-ai-world/User%20Research%20for%20Local%20Music%20App/User%20Research%20for%20Local%20Music%20App.007.png)

I joined this project only a year ago, and it had been around for several years before me, and there's been a lot of research in this field in general. I work for a nonprofit that creates an app that recommends local music to people. It kind of depends on what your definition of "local" is, but for me personally, local music is the kind of bands you can go see on the weekend that are just in your town. It's very accessible, and oftentimes it's people you know. It's part of your community.

![Localify.org and PorQuest](./media/human-app-ai-world/User%20Research%20for%20Local%20Music%20App/User%20Research%20for%20Local%20Music%20App.008.png)

Our previous app, Localify.org, provided a website that looks familiar if you've seen other music streaming services: you pick some artists you like, and it gives you recommendations for new ones you might want to check out. A classic feed of "for you."

Here's what the architecture formerly looked like.

![Architecture of Localify.org](./media/human-app-ai-world/User%20Research%20for%20Local%20Music%20App/User%20Research%20for%20Local%20Music%20App.009.png)

You can see at the top there are basically three or four different websites where we got information from. These were hardcoded. We didn't have anything that wasn't an API we had applied for or paid for. And the data goes to very specific places within our application: you see this big box in the middle, and all of that data from these APIs just goes to specific points in the process. From Wikipedia, we'd get data about where artists are from. From BandsInTown, we'd get information about their events. From Spotify, we'd get information about the artists themselves, like their picture or an audio preview.

Ultimately, we were limited by the services we were using, and the few ways in which we used them. For example, a lot of artists have trouble getting onto Spotify, or have principled objections to it. Even if we could replace Spotify with another music service, but the data would still be incomplete. There's no way you're going to find one music service with all of the local artists on it, because local artists are just a people in your town—they don't make music full-time, so they usually don't have the resources to blast their music everywhere the way you can if you're signed to a label.

We also had issues with event data. The event data from BandsInTown is very focused on places that sell tickets, because BandsInTown works with venue providers and ticketing services to make money as well as feed their algorithm. Smaller events that don't sell tickets, or charge cash at the door, or are free events for college students, maybe run by local community organizers and funded publicly or through donations... Those just don't show up on BandsInTown as much as we would like. So even though we developed this complex system for recommending artists and learning what artists you might be interested in, then translating that to local artists, the problem was we just didn't have enough local artists to recommend you. We needed more.

What if we could take this concept and generalize it, create a "magic black box," distilling all data sources and recommendation channels into one place?

![But what if?](./media/human-app-ai-world/User%20Research%20for%20Local%20Music%20App/User%20Research%20for%20Local%20Music%20App.010.png)

We could pull in data from everywhere, and I mean literally everywhere: community calendars, local radio stations, posters on a lamppost, on a street corner, social media like an artist's Instagram or newsletters from local fans, news websites, the artist's website directly, and of course every streaming service under the sun, YouTube, Deezer, Bandcamp, Spotify, SoundCloud. We could combine all of that and cross-reference everything to get the data we actually want. We would know that Artist X is from Place Y (they're local), they play music that sounds like Z, and because they play music that sounds like Z and you like Z, you'd probably like them a lot.

Here are some examples of the kind of data it was basically my dream to import into our system.

![Examples of data sources](./media/human-app-ai-world/User%20Research%20for%20Local%20Music%20App/User%20Research%20for%20Local%20Music%20App.011.png)

Through some technological advances and a lot of work — some of which my colleagues have presented posters about, like processing show posters automatically to intake event data — we've created this magic black box.

So what's inside our magic black box?

![What's in there?](./media/human-app-ai-world/User%20Research%20for%20Local%20Music%20App/User%20Research%20for%20Local%20Music%20App.012.png)

It's a large language model...

![It's a large language model](./media/human-app-ai-world/User%20Research%20for%20Local%20Music%20App/User%20Research%20for%20Local%20Music%20App.014.png)

...and a lot of SQL queries.

![And a lot of SQL queries](./media/human-app-ai-world/User%20Research%20for%20Local%20Music%20App/User%20Research%20for%20Local%20Music%20App.015.png)

So, is this evil?

![Is this evil?](./media/human-app-ai-world/User%20Research%20for%20Local%20Music%20App/User%20Research%20for%20Local%20Music%20App.016.png)

If you're someone deep in machine learning research, you're probably thinking, "Of course that's not evil! That's a normal, acceptable use of this technology." The problem is, our ideal user base is not machine learning researchers. It's local music fans and musicians. These are people who are very entrenched in the scene. They often play music or make art themselves, and they don't like what they're hearing about AI using their art to train models to replace them, or environmental harms happening where prominent AI companies with venture capital funding are building these hyperscale datacenters. In the minds of our users, AI models taking advantage of the art that they're producing, and come with a cost that they don't see as justifiable. So they're understandably skeptical of this type of application.

For example, we had one of our users send us a signed letter about refusing to be a part of a platform that is directly connected to Spotify.

![Wikipedia: Criticism of Spotify](./media/human-app-ai-world/User%20Research%20for%20Local%20Music%20App/User%20Research%20for%20Local%20Music%20App.017.png)

This is the Wikipedia page for criticism of Spotify, and you'll notice it's gigantic — because, like I said, these artists are very entrenched in the art world. They're very favorable toward copyright, and very protective of their work and of the things they spend their lives producing.

We also sometimes hear environmental arguments from the AI crowd.

![Environmental concerns](./media/human-app-ai-world/User%20Research%20for%20Local%20Music%20App/User%20Research%20for%20Local%20Music%20App.018.png)

On one hand, people say AI uses a lot of water; on the other hand, people who create AI models often say data-center pollution is necessary because AI is going to cure cancer and create huge scientific benefits that outweigh the costs. But the problem is we haven't seen it outweigh the cost so far. We haven't seen AI cure cancer. We haven't seen AI music that is incredibly interesting and captivating and tells a story that otherwise couldn't be told by a human. We haven't seen AI art created by a generative model that has the same emotional connection and transformative power as something produced by humans.

So I think it's definitely a balance.

![Ethical considerations for our magic black box](./media/human-app-ai-world/User%20Research%20for%20Local%20Music%20App/User%20Research%20for%20Local%20Music%20App.019.png)

I like to think of it as reducing the human and environmental costs while increasing the community benefit, because I think this service is only possible at this scale due to advances in AI and large language models. We spent a lot of time carefully thinking about how to balance this new technology, which lets us do things that weren't previously possible. We're able to intake more organic data and combine it in intuitive, language-based ways that weren't possible before. But how do we do this more efficiently? How do we do this in a way that is better for the environment? How do we use and deploy it in a manner that is ethical?

We do a number of things to address this. First and foremost, our entire business model is about supporting independent artists. We don't display LLM-generated text in the app. They're only used to collect or aggregate data behind the scenes, not to generate new things. We use paid APIs, so we're not just scraping the entire internet. We do use some scraping, but for most services we've talked to people directly about what we do and many websites in the music-data industry have been really excited about our project and offered us great APIs to use. We use LLMs on the backend to combine everything together into a complete picture of an artist from all these different data sources, which means that they are given maximal context and multiple sources to increase accuracy. We also use human art and design in our app, and we regularly audit the data the models output. We have our own internal dashboard where we can review everything the models produce and make sure it's accurate, since most of our team is involved in the local scene and can spot where things are right or wrong. We also run open-weight models locally, on hardware we're able to host  largely on carbon-neutral energy thanks to the Ithaca College's solar farm and other sustainability initiatives. So we do all of these things to reduce the human and environmental cost while increasing the community benefit.

That brings me to the results of the user interviews we did, which has informed a lot of what I've been saying, and is how we went about transforming our old model into this new AI-powered system without losing sight of our values.

![User research word cloud](./media/human-app-ai-world/User%20Research%20for%20Local%20Music%20App/User%20Research%20for%20Local%20Music%20App.020.png)

What was most unexpected was our user's high standards of safety and privacy. These are people who care about being protected from others interested in their data. For example, if you know someone is going to a local venue at a certain time, you now know exactly where they'll be. That could be dangerous if a bad actor were trying to find this person in real life. So even though it seems innocuous, something like "I'm going to be at the show on Friday" can be very private information depending on who it's shared with.

These people are also very frugal, creative, adaptable, and pragmatic. So we're looking at how to make local music more accessible and more affordable — how to make sure you know the price of a show, since sometimes you show up and have no idea what it costs. We've been working on incorporating different APIs and data sources so we can crowdsource that kind of pricing information from people who've already been to the show.

I've been talking about leveraging large language models within our system, but I haven't talked much about using them to create the systems with agentic programming. Given all this user research, you might think we could simply plug it into a coding agent and have it build the app. Right?

![OK, make the robots do it!](./media/human-app-ai-world/User%20Research%20for%20Local%20Music%20App/User%20Research%20for%20Local%20Music%20App.021.png)

Well, we've used agentic coding tools throughout this process to develop software more quickly, since our team is so small. It's very useful for rapidly iterating on prototypes.

![2 hours later](./media/human-app-ai-world/User%20Research%20for%20Local%20Music%20App/User%20Research%20for%20Local%20Music%20App.022.png)

But the problem is, if you just plug everything into an LLM, you get something like this: very basic wireframes that just reuse existing metaphors and affordances from other apps. That's fine, but it's not very interesting or artsy, and we're held to a higher standard. We want our app to feel unique. We're grounded in this local culture built on show posters, T-shirts, and album covers. There's so much art, design, and craft in these scenes. An app that's bland and similar to every other app on the market feels corporate, and feels like it doesn't match what our users want, since these are people who are creative and interested in seeing what other people and artists are doing.

So we decided to have humans do all of the design work upfront.

![Instead, make the humans do it FIRST](./media/human-app-ai-world/User%20Research%20for%20Local%20Music%20App/User%20Research%20for%20Local%20Music%20App.023.png)

We created mood boards, art assets, and different components. We created visual specifications for dynamic elements.

![Now let the robots do it](./media/human-app-ai-world/User%20Research%20for%20Local%20Music%20App/User%20Research%20for%20Local%20Music%20App.024.png)

We incorporated all of this with agentic programming into something creative, interesting, and fun to look at—because that's really important to us! We want people to feel the brush behind the pixels. We want something that feels creative, friendly, and like it could have come from one of these local community artists, because it came from *our* community of artists and designers. We're building from the aesthetic our users associate with, because we want them to feel welcome in the app and drawn to discovering other local artists. If our app looks like every other app, they'll simply close it and go to another one built by a corporate team with a bigger budget and more harmful economic incentives. So we have to take user feedback seriously and craft every inch of the app around our users, which is why user studies like this matter so much.

Of course, I'm not showing an app today. I'm giving a presentation because our app isn't ready yet.

![In the near future...](./media/human-app-ai-world/User%20Research%20for%20Local%20Music%20App/User%20Research%20for%20Local%20Music%20App.025.png)

It takes a lot more than just designing an app to make it actually work. We've done a lot of work behind the scenes updating our older backend to support new types and sources of data. We're working on this "magic black box," since we now have so many data sources and need to determine which are higher or lower quality, which are more trustworthy, and new features that bring user feedback and crowdsourced data into the very same model.

On top of that, we're planning new social features: seeing where your friends are going on a given day, meeting up with them, and sending and receiving notifications. One thing we found in our user studies was that a big reason people don't go to shows is they have no one to go with. We've tried to bridge that gap with social features layered on top of our existing music recommendation features. It's a one-click process once you say "this looks interesting," you can easily ask everyone on your friends list if they want to go. It's fine if only a couple of people respond, since you'll know who's going right in the app. With one tap, you can bring together a group of people to go see live music.

So, this has all been a big shift for us in terms of the data we're storing and how we handle it, along with the privacy concerns I mentioned earlier. It also means adding wholly new ways to store data to our existing app: previously our main endpoints were around searching for artists and events, sorting them, and recommending them based on what you like. Now we also have notification models, like pub/sub, and graph networks for friends. We also have to use technologies within the app itself, not just on our backend servers—things like Bluetooth Low Energy for connecting with nearby friends, and native notification systems on Android and iOS, while figuring out how much code we can share between the two platforms. This is a more integrated, system-level task than our previous app had to handle. So that's our work for the future.

Finally, let's go back to our first question:

![Back to our first question](./media/human-app-ai-world/User%20Research%20for%20Local%20Music%20App/User%20Research%20for%20Local%20Music%20App.026.png)

How do you build a very human app like this in a very artificial world?

![100 percent human, 100 percent art](./media/human-app-ai-world/User%20Research%20for%20Local%20Music%20App/User%20Research%20for%20Local%20Music%20App.027.png)

I think the answer is simple, even though it takes a lot of work to get right. At the end of the day, it's about combining human design, human art, human direction, and human creativity, and using LLMs to implement everyday tasks. It's about expanding the amount of work a small team can get done, but keeping lots of checks and balances in place. You must constantly review the environmental impact, the correctness, the quality of the code. You want to make sure people can see the hands behind the pixels you're displaying. As long as you have real oversight, you can combine human work and human creativity with the power and adaptability of large language models, but the results will always stay 100% human, and 100% art.

Thank you.

![How to make a very human app in a very artificial world](./media/human-app-ai-world/User%20Research%20for%20Local%20Music%20App/User%20Research%20for%20Local%20Music%20App.028.png)

*This research was funded by Ithaca College as part of the Summer Scholars program.*
*An early draft of this post was created with an LLM from an audio recording of a live presentation.*
*Scroll animations were developed with the assistance of Claude Sonnet 5.*

<script src="/_scripts/slide-viewer.js" defer></script>
<link rel="stylesheet" href="/_styles/slide-viewer.css">
