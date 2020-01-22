Hacking on Tweed
================

Hi there.

This document exists to help you, an aspiring Tweed developer, to
get up and running with a development environment that you can
hack in.

Day-to-Day Code Hacking
-----------------------

Tweed has a lot of moving parts, and all of those moving parts
need to be spinning to properly run unit tests, perform
interactive testing and exploration, etc.  To make all of this
easier, we've built a custom Docker image that bundles all of
these dependencies into an easy-to-use container.

The definition of this image is contained in `Dockerfile.test`.

Unit tests are run _inside_ of the test container, under the
watchful eye of Ginkgo + inotify -- as you change files in the
checked out local copy, the container will notice and re-run
affected test suites, without intervention.  In order to make this
work, the test container needs the working copy directory mounted
mounted into the container.

For this reason, you will most likely need to do development /
testing _on the box_ that runs the Docker daemon, so that file
systems can be bind-mounted into the container name space properly.

To run the unit tests in the background, run:

    make bg-tests

Now, you can hack on the source code and flip back to the tests
every once in a while to see if everything is still okay, from a
unit test perspective.
