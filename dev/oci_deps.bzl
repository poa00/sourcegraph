load("@rules_oci//oci:pull.bzl", "oci_pull")

def oci_deps():
    oci_pull(
        name = "wolfi_base",
        digest = "sha256:79d9c1e76beff31da0f182f30a2664dace9d9153cad8cbde7dba5edcef9e564d",
        image = "us.gcr.io/sourcegraph-dev/wolfi-sourcegraph-dev-base",
    )

    oci_pull(
        name = "wolfi_symbols_base",
        digest = "sha256:67731f797ebf8f6f1dcd08a5f4804adeee9ba2b7600d5ba8ed2329b64becd59a",
        image = "us.gcr.io/sourcegraph-dev/wolfi-symbols-base",
    )

    oci_pull(
        name = "wolfi_server_base",
        digest = "sha256:faa898afeba5f277cad13a16cd7c6baa304de389956433eb6d42fbb02a027ff9",
        image = "us.gcr.io/sourcegraph-dev/wolfi-server-base",
    )
