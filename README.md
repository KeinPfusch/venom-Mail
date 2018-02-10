# venom-Mail

Venom Mail is a gateway between SMTP/POP3 and the TOX protocol ( https://tox.chat/) 

It allows you to use a normal email client with the TOX Protocol, which is normally used for chat. 

Email addresses will look like:

"Max Mustermann" \<1B4CD4F8D84EE4540B2ABC97777448AEC5DA4073613DD6A61B1AF2CEA07D1B70895A4DC0DFB8@venom\>

The idea is:

Joe and Jill have a normal email client, capable of SMTP and POP3.

So when Joe wants to write to Jill, the flow is:

JoeCLient -> SMTP -> Venom on localhost -> TOX DHT -> Venom on localhost -> POP3 -> JillClient

of course both Venom daemons are supposed to run 24/7, so I am adding UPnP in order to run it on a raspi inside
the home network.
