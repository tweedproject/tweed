Tweed
=====

Service Instance Lifecycle
--------------------------

If we remove the data migration responsibilities, the service
instance orchestrator only has four jobs, which correspond quite
closely to the OSB API lifecycle actions:

  1. Provision - Build new thing, given inputs.
  2. Bind - (Optionally) create and communicate access creds.
  3. Unbind - (Optionally) deactivate creds.
  4. Deprovision - Tear thing down.

These can be neatly wrapped up in four different external
programs, provided by service instance specialists, for the
Tweed broker to leverage in order.

A Stencil has the following on-disk filesystem layout:

  - `$ROOT/lifecycle/provision` - Provisioning executable
  - `$ROOT/lifecycle/bind` - Binding executable
  - `$ROOT/lifecycle/unbind` - Unbinding executable
  - `$ROOT/lifecycle/deprovision` - Deprovisioning executable
  - `$ROOT/lifecycle/viable` - Viability detection executable
  - `$ROOT/lifecycle/files` - Interesting files manifest
  - `$ROOT/lifecycle/*` - Reserved for future use by Tweed
  - `$ROOT/*` - Reserved for use by the Stencil author

The `lifecycle/viable` executable is run to ensure that the local
environment (installed utilities, safe versions, etc.) are
suitable for the proper execution of the four lifecycle hooks.

The `lifecycle/files` executable is run on-demand to pull back any
interesting files (including their content) for review by
Tweed operators.  This might include things like Kubernetes
resource YAML files, BOSH manifests, BOSH task logs, etc.

The lifecycle hooks execute authenticated to a Vault
(via a `~/.saferc`), and with the following binaries already in
the $PATH:

  - All of a POSIX userland
  - `safe` for accessing Vault
  - `vault` for advanced `safe vault` calls
  - Utilities required for deployment, based on Tweed
    configuration (i.e. `bosh`, `kubectl`, etc.)
