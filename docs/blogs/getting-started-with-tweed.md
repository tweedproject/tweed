Getting Started with Tweed
==========================

Have you ever wanted to effortlessly deploy data services?  Large
ones?  Small ones?  Single-node?  Clustered?  Mesh?  I've got a
system for you; it's called **Tweed**.

Tweed is an _on-demand_ _service broker_ that provisions data
servces for both _shared_ (containerized) and _dedicated_
(VM-ified) workloads.

Don't worry, we'll unpack that mouthful as we go along.  But
first, let's get started with a Kubernetes cluster.  If you
already have one stood up and ready to `kubectl` away, great!  If
not, I highly recommend (INSERT LINK TO BLOG ARTICLE HERE).

Got a Kubernetes?  Awesome.

Tweed comes with a deployment you can use for evaluation.
It's hosted up on GitHub.

    $ kubectl apply -f (GITHUB URL)
    (OUTPUT)

That spins up a new namespace called `tweed` and deploys all
the bits and pieces that Tweed needs to properly deploy
services to the cluster you put it on.

    $ kubectl -n tweed deployment,po
    (OUTPUT)

Once the pod is in the `running` state, we can start poking at it
with the `tweed` CLI, and we'll provision our first service.

The technically correct way of running the `tweed` CLI is
this:

    $ kubectl exec -it -n tweed -c broker \
      $(kubectl -n tweed get po -l app=tweed -o name) \
      -- tweed

But that's way too much to remember, let alone type.  So we're
going to alias it:

    $ alias tweed='kubectl exec -it -n tweed -c broker $(kubectl -n tweed get po -l app=tweed -o name) -- tweed'

Let's give it a whirl, shall we?  To find out what a Tweed
can do, we need to take a look at its _catalog_:

    $ tweed catalog
    Service     Plan   #    Free?  Tags      Description
    =======     ====   =    =====  ====      ===========
    PostgreSQL                               A standalone, single-node PostgreSQL RDBMS
    [postgres]
                v9     0/2  no     postgres  PostgreSQL version 9.x
                [v9]               psql
                                   pg
                                   shared

                v10    0/1  no     postgres  PostgreSQL version 10.x
                [v10]              psql
                                   pg
                                   shared

                v11    0/1  no     postgres  PostgreSQL version 11.x
                [v11]              psql
                                   pg
                                   shared

The demo Tweed can deploy PostgreSQL (a relational SQL
database system) in three different versions: v9.x, v10.x, and
v11.x.

To spin one of these up, we ask Tweed to _provision_ a
service instance, and tell it which plan we want:

    $ tweed provision postgres/v9
    service instance scheduled for provisioning; thank you for your patience.
    instance i-30043cb6da5e04 is still provisioning; sleeping for another 2 seconds...
    instance i-30043cb6da5e04 is still provisioning; sleeping for another 2 seconds...
    instance i-30043cb6da5e04 is still provisioning; sleeping for another 2 seconds...
    instance i-30043cb6da5e04 is still provisioning; sleeping for another 2 seconds...
    instance i-30043cb6da5e04 is still provisioning; sleeping for another 2 seconds...
    instance i-30043cb6da5e04 is no longer provisioning.
    instance: i-30043cb6da5e04

    run tweed instance i-30043cb6da5e04 for more details.
    run tweed bind i-30043cb6da5e04 to get some credentials.

(Your instance identifier will be different than mine.)

We now have a new PostgreSQL pod spinning.  Tweed wraps all
of it up in a Kubernetes namespace, named after the instance ID:

    $ kubectl get -n i-30043cb6da5e04 all
    NAME                           READY   STATUS    RESTARTS   AGE
    pod/postgres-876cfbc4d-kz2lj   1/1     Running   0          2m26s

    NAME               TYPE       CLUSTER-IP   EXTERNAL-IP   PORT(S)          AGE
    service/postgres   NodePort   10.245.0.3   <none>        5432:31357/TCP   2m26s

    NAME                       READY   UP-TO-DATE   AVAILABLE   AGE
    deployment.apps/postgres   1/1     1            1           2m26s

    NAME                                 DESIRED   CURRENT   READY   AGE
    replicaset.apps/postgres-876cfbc4d   1         1         1       2m26s

To get access to the database server, we next _bind_ the service:

    $ tweed bind i-30043cb6da5e04
    service instance bind operation scheduled; thank you for your patience.
    instance i-30043cb6da5e04 is still binding; sleeping for another 2 seconds...
    instance i-30043cb6da5e04 is no longer binding.
    binding: b-84fc0f491a81d7

    run tweed instance i-30043cb6da5e04 for more details.
    run tweed bindings i-30043cb6da5e04 to show all bound credentials.

Tweed went out and created a new (randomized) user account,
inside of PostgreSQL, gave it a (randomized) password, and granted
it full access to the system.  We can see those details by
checking the _bindings_:

    $ tweed bindings i-30043cb6da5e04
    Binding           Credentials
    =======           ===========
    b-84fc0f491a81d7  {
                        "database": "pg1",
                        "host": "10.0.0.16",
                        "hosts": [
                          "10.0.0.16",
                          "10.0.0.17",
                          "10.0.0.18",
                        ],
                        "password": "s/7db2c8bf-c18d-4ca9-9c4d-93bdc6dde688",
                        "port": 31357,
                        "tryit": "PGPASSWORD=s/7db2c8bf-c18d-4ca9-9c4d-93bdc6dde688 psql -h 10.0.0.16 -p 31357 -U u2246d3ca pg1",
                        "username": "u2246d3ca",
                        "version": 9
                      }

(Again, your output will be different, but structurally similar).

To test drive PostgreSQL, we can use the `tryit` credential that
Tweed helpfully gave us (assuming you have `psql` installed
locally):

    $ PGPASSWORD=s/7db2c8bf-c18d-4ca9-9c4d-93bdc6dde688 psql -h 10.0.0.16 -p 31357 -U u2246d3ca pg1
    psql (11.5, server 9.6.16)
    Type "help" for help.

    pg1=>

Tada!  A PostgreSQL instance.

Let's build another one!

    $ tweed provision postgres/v9
    service instance scheduled for provisioning; thank you for your patience.
    instance i-e5beca94ab1788 is still provisioning; sleeping for another 2 seconds...
    instance i-e5beca94ab1788 is still provisioning; sleeping for another 2 seconds...
    instance i-e5beca94ab1788 is still provisioning; sleeping for another 2 seconds...
    instance i-e5beca94ab1788 is still provisioning; sleeping for another 2 seconds...
    instance i-e5beca94ab1788 is still provisioning; sleeping for another 2 seconds...
    instance i-e5beca94ab1788 is still provisioning; sleeping for another 2 seconds...
    instance i-e5beca94ab1788 is still provisioning; sleeping for another 2 seconds...
    instance i-e5beca94ab1788 is no longer provisioning.
    instance: i-e5beca94ab1788

    run tweed instance i-e5beca94ab1788 for more details.
    run tweed bind i-e5beca94ab1788 to get some credentials.


    $ tweed bind i-e5beca94ab1788
    service instance bind operation scheduled; thank you for your patience.
    instance i-e5beca94ab1788 is still binding; sleeping for another 2 seconds...
    instance i-e5beca94ab1788 is no longer binding.
    binding: b-3989300556956a

    run tweed instance i-e5beca94ab1788 for more details.
    run tweed bindings i-e5beca94ab1788 to show all bound credentials.


    $ tweed ls
    ID                State  Service   Plan
    ==                =====  =======   ====
    i-30043cb6da5e04  quiet  postgres  v9
    i-e5beca94ab1788  quiet  postgres  v9

That was fun!  Let's do it again!

    $ tweed provision postgres/v9
    unable to provision a postgres / v9 service instance: too many instances of this service have been provisioned

Uh-oh!  Looks like we've built too many PostgreSQL instances!
Indeed, if we look back at the catalog again:

Service     Plan   #    Free?  Tags      Description
=======     ====   =    =====  ====      ===========
PostgreSQL                               A standalone, single-node PostgreSQL RDBMS
[postgres]
            v9     2/2  no     postgres  PostgreSQL version 9.x
            [v9]               psql
                               pg
                               shared

            v10    0/1  no     postgres  PostgreSQL version 10.x
            [v10]              psql
                               pg
                               shared

            v11    1/1  no     postgres  PostgreSQL version 11.x
            [v11]              psql
                               pg
                               shared

You'll notice that the `postgres` / `v9` plan is currently at
`2/2` allowed instances (the `#` column).  Tweed enforces
limits and quotas on every service plan, to aid capacity planning
and help with infrastructure spend -- this gets even _more_
helpful when we start talking about deploying VMs for our service
instances!

To tear down some of these instances, ask Tweed to
_deprovision_ them:

    $ tweed deprovision i-30043cb6da5e04
    service instance scheduled for deprovisioning; thank you for your patience.
    instance i-30043cb6da5e04 is not yet gone; sleeping for another 2 seconds...
    instance i-30043cb6da5e04 is not yet gone; sleeping for another 2 seconds...
    instance i-30043cb6da5e04 is now gone.

    run tweed instance i-30043cb6da5e04 for more historical information.
    run tweed purge i-30043cb6da5e04 to remove all trace of this service instance.



Everything's Coming Up Services
-------------------------------

Now that you've seen Tweed, you might be wondering what kind
of fun and mischief you can get up to with it.  Here's some
interesting ideas:

  - On-demand Redis cache / session storage
  - Clustered PostgreSQL databases
  - MongoDB cluster(s)
  - Deploy Kubernetes and lather, rinse, repeat!
  - Per-application credentials Vaults
  - A Ci/CD instance for every delivery team

Stay tuned, as we've got a whole series of guides, howtos,
tutorials, and more, all about Tweed and modern data
services.
