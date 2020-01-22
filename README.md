Tweed - A Data Services Broker
==============================

Tweed is an [Open Service Broker][osbapi] that can spin up
arbitrary data services on top of Kubernetes and BOSH, for shared
and dedicated use cases.

[osbapi]: https://www.openservicebrokerapi.org/



Current Status: Verrrry Alpha
-----------------------------

Tweed is very much alpha software.  We're currently working
through some changes to the original MVP to allow for better
operator experience.



Installing and Running Tweed
----------------------------

Yeah, so about that.  Tweed is very hard to run by hand.  We're
working on that.  [Hopefully we'll have a Helm chart by the end of
February!][gh7]

[gh1]: https://github.com/tweedproject/tweed/issues/7



Hacking on Tweed
----------------

If you want to try to hack on Tweed and get us closer to a beta
release, check out the [HACKING.md][dev] guide, and bring your
Docker daemon and/or Kubernetes cluster along with.

[dev]: HACKING.md



Contributing
------------

This project falls under the [governance model][gov] of the larger
[Tweed Project][gh], and contibutions are governed by our [Code of
Conduct][coc] and [CONTRIBUTING.md][contrib] guide.

[gov]:     https://github.com/tweedproject/governance
[gh]:      https://github.com/tweedproject
[coc]:     CONDUCT.md
[contrib]: CONTRIBUTING.md
