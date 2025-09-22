---
slug: no-cloud
title: You can make it rain without a cloud
date: 2025-09-22
author:
  name: Jon Wear
summary: You can run software services at high scale on really cheap hardware.
tags:
  - name: nocloud
    slug: nocloud
  - name: localfirst
    slug: localfirst
hero:
  src: /images/hero_article_1.png
  alt: Corridor lights
css: /css/retro-sci-fi.css
draft: true
reading_time_min: 2
---

### Background : Game Launch

In 2023 I launched a backend service that provided a feature for the biggest selling single player video game in 2023.  The game was Hogwarts Legacy.  The service helped connect a players gaming profile with their official profile at wizardingworld.com (now harrypotter.com).  It wasn't highly technical with regard to what it did, push data here, pull data there, make sure a user can only get to their data.  Basic stuff.  But it had to be fast and scalable and getting it right the next day wasn't going to work.  When a game launches, stuff needs to work _at launch_, not days later after you understand load patterns.  The whole thing ran in the AWS cloud.  I used a few Lambdas, API Gateway, DynamoDB, SQS and implemented some webhooks.  I did the server to server authentication using JWT & RSA keypairs.  It worked great.  It held up fine at launch and as far as I know the service is still in use.  But it wasn't cheap.  Since the load was high for the first few weeks the cost was high with all of those API Gateway and Lambda calls.  Of course, the beauty of Lambda is that when traffic falls off, cost falls off which it did after the initial crush of the launch.

But one thing I've thought about since then is, "Was all that cloud infrastructure necessary?"  Building something that scales, and I mean _really_ scales can be difficult.  But, I think we just assume that we need all the cloud stuff when we might not.  In larger companies the true cost of the cloud is often abstracted away from the developers.  There's just always more infrastructure to be had, so why not spin up another EC2 instance, or docker container, or kube cluster or whatever it is and just go to town?

### The Experiment

So I did an experiment.  I bought a P4 Light Gaming Mini PC(32GB of RAM and 1TB of storage) for $330.  I formatted the drive, installed Ubuntu server and Docker and provisioned a container for mySQL, RabbitMQ and Redis.  I wrote a simple web server in Go and made it accessible to the outside web using Cloudflare Tunnels.  I also put Tailscale on the server so I could work on it from anywhere (whether the cloudflare tunnel was up or not).  

This service exposed a REST API that accepted POSTS, DELETES and GETS.  A POST wrote to rabbitMQ, golang workers processed the queue and wrote everything to the database.  GETs would pull first from the database but also throw the result into Redis so that future calls were cached (with a sensible TTL).  Finally, a Delete would write to a delete queue, where workers would then pick up the queue item and delete the record from the database and remove it from the cache. I got all all this _totally free software_ running and decided to have some fun.

I ran some load tests.

### Results

At 2K req/s it ran just fine.  The queue didn't back up, the responses times were well under 100ms.  So I cranked it to 3K req/s.  Had to add a few more workers, response times crept up but it all ran fine.  At 5K req/s it still ran fine but the response times were creeping into the 500ms range.  That's a bit slow for client facing stuff but could work for server to server.  Anyway, I let it run at 5K req/s for hours and it just...ran.  Very rarely I'd get back a 404 error because a GET tried to run before the queue item from the POST had been processed.  But these things happen and further optimizations can take care of it.


### Limitations

Now, would you want to run a AAA game backend service from a basement?  Probably not although you could maybe get away with it (not infrastucture advise!).  A production service needs backups, it needs a recovery plan it might need a failover infrastructure.  Yup, it needs all those things.  The cloud gives you that.  But a nocloud setup can give you that as well.  Back up cron jobs are easy to setup.  Off site storage is easy to setup.  That could even be your one cloud touch point, upload the backups to S3.  You can setup mySQL and Postgres to have read replicas that can be used in a failover situation.   "But wait!" you say.  "What about the costs of having some operations engineer babysit all this?"  My answer, what about it?  I've yet to see a cloud deployment that didn't have _team_ of Devops people checking out cloud dashboards and waking up at 3am on Christmas Eve because a certain data center in Virginia went down again.  You have to pay sys admin types anyway.  If you don't have sys admin types you have your developers doing the baby sitting and I know from experience that they would much rather be writing code(or at least review PR's from AI agents) than writing terraform scripts and figuring out every changing cloud documentation.


### Takeaways

I'm not saying we should all close our AWS accounts and run things locally.  But I am saying there is probably _some_ amount of infrastructure that you could run outside of the cloud.  A small production service.  Your Dev and QA environments.  So much work has been done on open source projects.  We have amazingly powerful open source web servers, queueing system, database servers and nosql data stores like redis.  Much of the power of those _free_ software resources are abstracted away behind cloud services because we think we have to run it all in the cloud or we'll get laughed at.

Try it out.  Throw Ubuntu on an old laptop or get a mini-pc of some kind and see what you can do with it.  You may be surprised at what you can do outside of the clouds.


