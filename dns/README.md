DNS
=====

Python DNS implementation of [SCION](http://www.netsec.ethz.ch/research/SCION)


authoritative_resolver.py
--------------

Contains the resolver for the *authoritative server*.
The authoritative resolver is responsible for his zone and answers queries
with the matching data.

dummy_client.py
--------------
A client sending queries to the server and printing the answers.


dnslib9.0.1
--------------
[DNSlib](https://pypi.python.org/pypi/dnslib) is a simple library allowing to encode/decode DNS wire-format packets.
It supports 3.2+ python.


* [doc](https://github.com/netsec-ethz/scion/tree/master/doc) contains documentation and material to present SCION
* [infrastructure](https://github.com/netsec-ethz/scion/tree/master/infrastructure)
* [lib](https://github.com/netsec-ethz/scion/tree/master/lib) contains the most relevant SCION libraries
* [topology](https://github.com/netsec-ethz/scion/tree/servers/topology) contains the scripts to generate the SCION configuration and topology files, as well as the certificates and ROT files

