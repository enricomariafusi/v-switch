#   v-switch

V-switch is an encrypted virtual switch. Following the concept of Tinc (part of it) and SDN (Sowftare Defined Networks) it creates a virtual interface
which peers with other daemons around the internet. All the machines running the same daemon, with the same encryption key
will have a device configured, which behaves like it was cabled to the same physical switch.  

The aim is to create a LAN across the internet or inside the cloud, where the machine just appears to be connected
each others on layer 2. Adding a new machine to the switch will advertise each other machine _without the need to provision them_ , **unlike Tinc**. This is to be able to use it inside a cloud while autoscaling: no provisioning is needed, there is only one key _per virtual switch_.

Ideally, when the cloud will autoscale, all the bundles with v-switch installed will "automagically" plug into the virtual switch. (then, to run a dhcp daemon on a port would allow to give the IP address, like it was a real switch). 


Ready for testing. TODO:

1. CONFIGURATION GUIDE
2. LINT code.

