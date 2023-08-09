# VoIP - Studio Quality Audio Call

### This is an example for a contractor grade audio call package implemented in Golang and JavaScript

```
Licensed under Creative Commons CC0.

To the extent possible under law, the author(s) have dedicated all copyright and related and
neighboring rights to this software to the public domain worldwide.
This software is distributed without any warranty.
You should have received a copy of the CC0 Public Domain Dedication along with this software.
If not, see <https:#creativecommons.org/publicdomain/zero/1.0/legalcode>.
```

## Abstract

This is an example audio call package implemented in golang and html/js.
It allows a peer to peer connection between two parties.

You run a simple container yourself in your cloud or on premises environment.
You can also rewrite and reuse parts as you wish due to the most permissive creative commons zero licensing as of today.

It reduces the risks to privacy suggesting an end-to-end encrypted framework.
The keys need to be shared through external channels different from the routing domain.

Previous generation communication tools like Whatsapp or Zoom share keys within their tls domain.
Attackers compromising client trust certificates can impersonate such sites without the site owner's knowledge.

## Try it

### Trial service

[Use our ad supported pilot](https://l.eper.io)

### Deploy yourself

[![Deploy to DO](https://www.deploytodo.com/do-btn-white-ghost.svg)](https://cloud.digitalocean.com/apps/new?repo=https://github.com/eper.io/tree/main&refcode=48f147bd7dcd)

Make sure you update the environment settings for a production site.
1. Go to `App/Settings and App-Level Environment Variables` to set a unique `APIKEY`.
2. Set a `CNAME DNS record` to the app name (optional) i.e `voip.example.com`
3. Go to `App/Settings and App-Level Environment Variables` to set a unique `SITEURL` to either `https://voip.example.com` or the DO assigment like `https://voip-zhx34.ondigitalocean.app`.

## Security

The connection is carried out over https websockets.

There are two layers of security. We provide a private key called leaf that is known by the two clients only.

We also allow you to set your own tls certificate and an api key for the servers called the root.
This secures the client server communication ensuring you paid for, or you are authorized to use the servers.

The administrator sets this root apikey for the line with the APIKEY environment variable.

The client encryption leaf key is a hash search tag in the url that is generated upon the first download instructed by `#generate_leaf`.

```
https://line.example.com:49999/line.html?apikey=ZXNTYC...IBR#generate_leaf
```

The generated end-to-end encryption key embedded into the hash search tag never reaches the server.
Browsers typically cut it off being an instruction to the browser not the backend.

It can be shared through a separate channel deemed to be private with the entire url.
The privacy of this separate channel sets adequate the level of privacy, not our service.
Our service is merely a router.

Example url with leaf generated:
```
https://line.example.com:49999/line.html?apikey=ZXNTYC...IBR#leaf_AAMVfpBo2J...MOjej3BZw
```

You can use NGINX to secure the socket. https://www.nginx.com/blog/using-free-ssltls-certificates-from-lets-encrypt-with-nginx
You can also use the DigitalOcean example above how to leverage your cloud provider to perform the tls closure.

This is important.
We have a central way to manage TLS parameters.
You do not expose yourself to upgrade issues, if microservices need security updates.
Only nginx has read permission to your private keys. This is an important best practice.
You need to bring your own domain and tls certificates.

This is a contractor grade solution.
We advise you to follow your organization's security best practices.
You can secure it with your own standards and cloud infrastructure.

Still you can try our integrated solution live.
[Hop](https://l.eper.io)

## Protocol

We do not use webrtc to give you full control over the packets.
We simply pass a 200kHz Pulse Code Modulated audio of 16 bits per sample.
The total bandwidth requirement is about 7Mbps full duplex.

**We suggest 20Mbps connections, if the channel is shared.**
This reduces the probability of any glitches or echo.

The power lies in the simplicity.
You do not need to worry about supporting different hardware, etc.
This protocol is always supported everywhere.

Do you have issues with bandwidth?
You can just buy a better internet connection.
Contact us for suggestions. [Contact](mailto:hq@schmied.us?subject=Fiber%20Sales)

We do not use variable bitrates, because we learned that available bandwidth changes over time.
Carriers can do this, if their backbone networks merge traffic.
It is easier to communicate a maximum bandwidth that is signed off by the carrier.
This reduces our support costs.

You can set &mobile=1 in the root key to limit the bandwidth need for suburban mobile channels.

## Design considerations

The VoIP audio call package follows these design principles.
- Simplicity supports the use of microcontrollers on one or both sides. **The client code is 200 JavaScript lines**
- Remove any special codec and hardware requirements for the broadest support. We support plain old pulse code modulation digitally sampling the audio levels.
- Use studio quality 200 kHz and 16 bits for each pulse sample. Studios originally used 192 kHz audio sampling and 24 bits.
- Simplify the codebase, so that you can run the server on your own.
- Let your operations team deal with security. This helps to follow organizational standards,
and it reduces support costs by eliminating the need for any extra security reviews.
- Authorization is established by an `apikey` that is shared by participants privately over their own environment (email, corporate chat).
- API key should be at least 72 characters containing a-zA-Z. This helps to prevent brute force attacks.
- High bandwidth connections over 20Mbps can easily handle pulse code modulated audio with the lowest latency. There is no extra buffering needed for convolutional codecs like MP3 ensuring the lowest codec latency possible. (Martucci, 1994)
- We have a limitation using Websockets that we may change to udp in the future to achieve even lower latencies.
The websocket solution allows us and you, the contractor to reach a broad set of leads of the customer base. 
- We implement websockets privately to avoid security risks of third party components. There are **no dependencies**.
- No dependencies eliminate a large attach surface of changing versions and the resulting code reviews.

## Requirements

Please use a connection of at least 7 megabits per second on each client and the server to run this example.
However, shared channels may require more than 20 megabits per second.
A one gigabit connection is good to present the full quality safely to a prospective opportunity.

## Usage

2. Establish at a TLS 1.3 or better closure using the documentation below.

```dockerfile
https://www.nginx.com/blog/using-free-ssltls-certificates-from-lets-encrypt-with-nginx
```

2. Build te docker file.
```
docker build . -t line.example.com/line
```

3. Run the server side in docker. Set APIKEY to a unique cryptographically strong random value.

```
docker run -d --rm --name line -p 7777:443 -v /etc/letsencrypt/live/line.example.com/fullchain.pem:/tmp/fullchain.pem:ro -v /etc/letsencrypt/live/line.example.com/privkey.pem:/tmp/privkey.pem:ro line.example.com/line
```

4. Connect to the server from the peer, who hosts the meeting. This will generate a leaf key for end to end encryption

```dockerfile
https://line.example.com/?apikey=8f711a8f43a3d6fbf9c367a8cd1b68b14db2781273c56e206040ecd31761f9d3#generate_leaf
```

5. Connect the second peer using the same URL.

```dockerfile
https://line.example.com/?apikey=8f711a8f43a3d6fbf9c367a8cd1b68b14db2781273c56e206040ecd31761f9d3#leaf_sywICvC1tnE7D5IOdPahv0NzQC-C3mpNLhluhEEW0vA
```

6. Example run logs

```
connection received
connection received
```

## References

Martucci, S. A. (May 1994). "Symmetric convolution and the discrete sine and cosine transforms". IEEE Transactions on Signal Processing.

## Contact and support

Do you need servicing from us?

- support agreement
- warranty, and reliability guarantee
- patent licensing and third party patent licensing
- security certification
- legal and privacy compliance certificates
- integration (Kubernetes, Helm, etc.)

**Please contact hq@schmied.us at Schmied Enterprises LLC.**
