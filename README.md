#   v-switch

V-switch is an encrypted virtual switch. Following the concept of Tinc (part of it) and SDN (Sowftare Defined Networks) it creates a virtual interface
which peers with other daemons around the internet. All the machines running the same daemon, with the same encryption key
will have a device configured, which behaves like it was cabled to the same physical switch.  

The aim is to create a LAN across the internet or inside the cloud, where the machine just appears to be connected
each others on layer 2. Adding a new machine to the switch will advertise each other machine _without the need to provision them_ , **unlike Tinc**. Once you have the same key, you are in.

Encryption is using AES256 + the GPG container, meaning the key MUST be as long as the MTU. This makes the flow relatively safe.

V-Switch will take care to self-configure the interface and to keep the ARP table updated on your kernel using the linux netlink interface (ip).

Ready for testing (it works with 2 nodes) TODO:

1. CONFIGURATION GUIDE
2. LINT code.

