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



Interactive Testing via Kubernetes
----------------------------------

Because of all of Tweed's moving parts, one way to spin it up for
interactive use, testing, and exploration is to actually deploy it
to Kubernetes.

To do this with the latest official Tweed images, target a
Kubernetes cluster and then run:

    make deploy
    source env/dev/envrc

This creates a `tweed` namespace and deploys the broker and all of
its pieces (a credentials vault, a database, etc.) into that.
To use an alternate namespace, set the `$NAMESPACE` environment
variable:

    NAMESPACE=dev-tweed make deploy
    source env/dev/envrc

The second command updates your current shell with the `TWEED_*`
environment variables that you will need to seamlessly interact
with the broker via the `tweed` CLI, without having to specify a
bunch of arguments every time.

The details of the Kubernetes deployment can be found in the
`env/dev/k8s.yml`.  If you are interested solely in the
contents of the manifest, as rendered (given the `$IMAGE`,
`$NAMESPACE` and `$VERSION` environment variables), run:

    make deploy.yml

If you are making changes to Tweed itself (as opposed to
developing a new stencil), you'll want to override the broker
image that gets deployed by setting the `$IMAGE` environment
variable:

    IMAGE=filefrog/tweed make deploy

To build (and push) your alternate image:

    IMAGE=filefrog/tweed make docker
    docker push filefrog/tweed

Note: you do **not** want to use `make push` to push your Docker
images upstream; that Makefile target does a bunch of semantic
versioning re-tagging, making it a bad idea for one-off dev tags.

When you're all done and want to tear it down, you can either just
delete the Kubernetes namespace ("tweed", by default), or run:

    make retire

Remember: you can chain Makefile targets in a single invocation of
`make`, allowing you to delete and redeploy with just:

    make retire deploy
